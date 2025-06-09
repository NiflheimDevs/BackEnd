package servicesimpl

import (
	"context"
	"log"
	"time"

	"github.com/niflheimdevs/backend/internal/application/dto"
	"github.com/niflheimdevs/backend/internal/application/services"
	"github.com/niflheimdevs/backend/internal/domain/enums"
	"github.com/niflheimdevs/backend/internal/domain/exceptions"
	elasticmodel "github.com/niflheimdevs/backend/internal/domain/models/elastic"
	"github.com/niflheimdevs/backend/internal/domain/repositories/elastic"
	repositories "github.com/niflheimdevs/backend/internal/domain/repositories/postgres"
	"github.com/niflheimdevs/backend/internal/domain/repositories/postgres/transaction"
	"github.com/niflheimdevs/backend/internal/utils"
)

//? if super admin can create roles then i should do query for permission check

type TeamService struct {
	TeamRepo           repositories.TeamRepo
	UserRepo           repositories.UserRepo
	RoleRepo           repositories.RoleRepo
	TransactionManager transaction.TxManager
	FileService        services.FileService
	SearchRepo         elastic.SearchRepo
}

func NewTeamService(
	teamRepo repositories.TeamRepo,
	userRepo repositories.UserRepo,
	roleRepo repositories.RoleRepo,
	tManager transaction.TxManager,
	fileService services.FileService,
	sp elastic.SearchRepo,
) *TeamService {
	return &TeamService{
		TeamRepo:           teamRepo,
		UserRepo:           userRepo,
		RoleRepo:           roleRepo,
		TransactionManager: tManager,
		FileService:        fileService,
		SearchRepo:         sp,
	}
}

func (ts *TeamService) BehindCurtainTeam(ctx context.Context, tx transaction.Tx, userid int) {

	ts.TeamRepo.CreateOneManTeam(ctx, tx, userid)
	//? OR:
	// teamid := ts.TeamRepo.CreateOneManTeam(ctx, tx, userid)
	// err := ts.TeamRepo.AddMemberWithTx(ctx, tx, userid, teamid, "", enums.TEAM_NEWBIE)
	// if err != nil {
	// 	panic(exceptions.Exception{
	// 		Tag:    exceptions.CONFLICT_ERROR,
	// 		Errors: []exceptions.SpecificError{exceptions.ALREADY_A_MEMBER},
	// 	})
	// }
}

func (ts *TeamService) CreateTeam(userid int, teamInfo *dto.TeamCreateDto) int64 {
	if userid < 0 {
		panic(exceptions.Exception{
			Tag: exceptions.UNAUTHORIZED,
			Errors: []exceptions.SpecificError{
				exceptions.AUTH_ACCESS_DENIED,
			},
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*3)
	defer cancel()

	tx, err := ts.TransactionManager.Begin(ctx)
	if err != nil {
		log.Println("TeamError: transaction creation failed. details:", err)
		panic(exceptions.Exception{
			Tag: exceptions.INTERNAL_ERROR,
		})
	}

	defer func() {
		if rec := recover(); rec != nil {
			tx.Rollback(ctx)
			panic(rec)
		}
	}()

	teamid := ts.TeamRepo.CreateTeam(ctx, tx, teamInfo.Title, teamInfo.Description)

	err = ts.TeamRepo.AddMemberWithTx(ctx, tx, userid, teamid, "", enums.TEAM_OWNER)
	if err != nil {
		log.Println("TeamError: owner couldn't be added. details:", err)
		panic(exceptions.Exception{
			Tag: exceptions.CONFLICT_ERROR,
		})
	}

	err = tx.Commit(ctx)

	if err != nil {
		log.Println("TeamError: no commit. detail:", err)
		panic(exceptions.Exception{
			Tag: exceptions.INTERNAL_ERROR,
		})
	}
	ts.addMembersFunc(teamid, teamInfo.Members)

	return teamid
}

func (ts *TeamService) GetTeamsForUser(userid int, active, dontCare int) []dto.GetTeamPreviewDto {
	res := ts.TeamRepo.GetTeamsForUser(userid, active != 0, dontCare != 0)

	for i := 0; i < len(res); i++ {
		res[i].Profile = ts.FileService.GetTeamProfilePhotoURL(res[i].ID, false)
	}

	var owners []dto.ReadMemberDto
	for i := 0; i < len(res); i++ {
		owners = ts.TeamRepo.GetMembersForTeamFilterdByRole(res[i].ID, enums.TEAM_OWNER)
		if len(owners) == 1 {
			res[i].OwnerInfo.Info = owners[0].Info
			res[i].OwnerInfo.Profile = ts.FileService.GetProfilePhotoURL(owners[0].Info.Userid, false)
		} else {
			log.Println("MemberOwnerError: check owners for team", res[i].ID)
		}
	}
	return res
}

func (ts *TeamService) GetTeam(commanderid int, teamid int64) *dto.GetTeamDto {
	var res dto.GetTeamDto

	res.Info = ts.TeamRepo.GetTeam(teamid)
	if res.Info == nil {
		panic(exceptions.Exception{
			Tag: exceptions.NOT_FOUND,
		})
	}
	res.Profile = ts.FileService.GetTeamProfilePhotoURL(res.Info.ID, true)

	members := ts.TeamRepo.GetMembersForTeam(teamid)
	res.Members = make([]dto.SendMemberDto, len(members))
	for i := 0; i < len(members); i++ {
		res.Members[i].Info = members[i].Info
		res.Members[i].Role = members[i].Role.String()
		// ? takes time! bottleneck
		res.Members[i].Profile = ts.FileService.GetProfilePhotoURL(res.Members[i].Info.Userid, false)
		if members[i].Info.Userid == commanderid {
			res.UserID = commanderid
			perms := members[i].Role.GetPermissionsForRole()
			res.Permissions = make([]string, len(perms))
			for i, p := range perms {
				res.Permissions[i] = p.String()
			}
		}
	}
	return &res
}

func (ts *TeamService) GetInternalTeamInfo(teamid int64) *dto.GetInternalTeamInfo {
	team, err := ts.TeamRepo.GetTeamInfo(teamid)
	if err != nil {
		team, err = ts.TeamRepo.GetOneManTeamInfo(teamid)
		if err != nil {
			log.Println("OneManTeamError: details:", err)
			panic(exceptions.Exception{
				Tag: exceptions.INTERNAL_ERROR,
			})
		}
		team.Type = 2
		team.Profile = ts.FileService.GetProfilePhotoURL(int(team.ID), false)
	} else {
		team.Type = 1
		team.Profile = ts.FileService.GetTeamProfilePhotoURL(teamid, false)
	}

	return team
}

func (ts *TeamService) UpdateTeamInfo(userid int, info *dto.UpdateTeamInfoDto) {
	if userid < 0 {
		panic(exceptions.Exception{
			Tag: exceptions.UNAUTHORIZED,
			Errors: []exceptions.SpecificError{
				exceptions.AUTH_ACCESS_DENIED,
			},
		})
	}

	member := ts.TeamRepo.GetMemberForTeam(info.ID, userid)
	if member == nil {
		panic(exceptions.Exception{
			Tag: exceptions.FORBIDDEN,
		})
	}

	//TODO: if super admin can create roles then it should be a query instead of static call
	if !utils.Contains(member.Role.GetPermissionsForRole(), enums.EDIT_INFO) {
		panic(exceptions.Exception{
			Tag:    exceptions.FORBIDDEN,
			Errors: []exceptions.SpecificError{exceptions.LACKS_PERMISSION},
		})
	}

	err := ts.TeamRepo.UpdateTeam(info.ID, info.Title, info.Description)

	if err != nil {
		panic(exceptions.Exception{
			Tag: exceptions.INTERNAL_ERROR,
		})
	}
}

func (ts *TeamService) DeleteTeam(commanderid int, teamid int64) {
	if commanderid < 0 {
		panic(exceptions.Exception{
			Tag: exceptions.UNAUTHORIZED,
			Errors: []exceptions.SpecificError{
				exceptions.AUTH_ACCESS_DENIED,
			},
		})
	}

	member := ts.TeamRepo.GetMemberForTeam(teamid, commanderid)
	if member == nil {
		panic(exceptions.Exception{
			Tag: exceptions.FORBIDDEN,
		})
	}

	//TODO: if super admin can create roles then it should be a query instead of static call
	if !utils.Contains(member.Role.GetPermissionsForRole(), enums.REMOVE_MEMEBER) {
		panic(exceptions.Exception{
			Tag:    exceptions.FORBIDDEN,
			Errors: []exceptions.SpecificError{exceptions.LACKS_PERMISSION},
		})
	}

	err := ts.TeamRepo.DeleteTeam(teamid)
	if err != nil {
		panic(exceptions.Exception{
			Tag: exceptions.INTERNAL_ERROR,
		})
	}

	go ts.FileService.DeleteTeamProfilePhoto(teamid)

}

// TODO: email? some sort of request must be sent and then when it is accepted, the member gets added
// ! this version is naive
func (ts *TeamService) AddMembers(userid int, teamid int64, members []int) {
	if userid < 0 {
		panic(exceptions.Exception{
			Tag: exceptions.UNAUTHORIZED,
			Errors: []exceptions.SpecificError{
				exceptions.AUTH_ACCESS_DENIED,
			},
		})
	}
	member := ts.TeamRepo.GetMemberForTeam(teamid, userid)
	if member == nil {
		panic(exceptions.Exception{
			Tag: exceptions.FORBIDDEN,
		})
	}

	if !utils.Contains(member.Role.GetPermissionsForRole(), enums.ADD_MEMBER) {
		panic(exceptions.Exception{
			Tag:    exceptions.FORBIDDEN,
			Errors: []exceptions.SpecificError{exceptions.LACKS_PERMISSION},
		})
	}

	ts.addMembersFunc(teamid, members)
}

func (ts *TeamService) addMembersFunc(teamid int64, members []int) {
	var err error
	for _, id := range members {
		err = ts.TeamRepo.AddMember(id, teamid, "", enums.TEAM_NEWBIE)
		if err != nil {
			log.Println("AddMemberError: teamid:", teamid, "member:", id, "datail:", err)
		}
	}
}

func (ts *TeamService) LeaveTeam(userid int, teamid int64) {
	member := ts.TeamRepo.GetMemberForTeam(teamid, userid)

	err := ts.TeamRepo.RemoveMember(userid, teamid)
	if err != nil {
		panic(exceptions.Exception{
			Tag: exceptions.CONFLICT_ERROR,
		})
	}

	if member.Role == enums.TEAM_OWNER {
		ts.TeamRepo.DeleteTeam(teamid)
	}

}

func (ts *TeamService) GetTeamsForBidding(userid int) *dto.TeamListDto {
	if userid < 0 {
		panic(exceptions.Exception{
			Tag:    exceptions.UNAUTHORIZED,
			Errors: []exceptions.SpecificError{exceptions.AUTH_TOKEN_EXPIRED},
		})
	}

	var res dto.TeamListDto

	userModel, err := ts.UserRepo.FindUserByID(userid)
	if err != nil {
		panic(exceptions.Exception{
			Tag:    exceptions.NOT_FOUND,
			Errors: []exceptions.SpecificError{exceptions.USER_NOT_FOUND},
		})
	}
	res.Userid = userModel.ID
	res.FirstName = userModel.FirstName
	res.LastName = userModel.LastName
	res.Username = userModel.Username
	res.UserProfile = ts.FileService.GetProfilePhotoURL(userid, false)
	res.OneManTeamid, err = ts.TeamRepo.GetOneManTeamID(userid)

	if err != nil {
		log.Println(userid, "does not have one man team!!!!!")
		panic(exceptions.Exception{
			Tag: exceptions.INTERNAL_ERROR,
		})
	}

	teams := ts.TeamRepo.GetTeamsForUserWithRole(userid)
	for _, team := range teams {
		res.Teams = append(res.Teams, dto.GetTeamBidDto{
			ID:          team.ID,
			Description: team.Description,
			Title:       team.Title,
			Position:    team.Position,
			CanBid:      utils.Contains(team.RoleId.GetPermissionsForRole(), enums.BIDDER),
			Profile:     ts.FileService.GetTeamProfilePhotoURL(team.ID, false),
		})
	}
	return &res

}

func (ts *TeamService) GetAllTeamsIDs(userID int) []int64 {
	if userID > 0 {
		ids := ts.TeamRepo.GetEveryTeamID(userID)
		oneManID, _ := ts.TeamRepo.GetOneManTeamID(userID)
		ids = append(ids, oneManID)
		return ids
	} else {
		return nil
	}
}

func (ts *TeamService) KickMemebr(commanderid int, poorGuysid []int, teamid int64) {
	if commanderid < 0 {
		panic(exceptions.Exception{
			Tag: exceptions.UNAUTHORIZED,
		})
	}

	var isLeaving bool = false

	member := ts.TeamRepo.GetMemberForTeam(teamid, commanderid)
	if member == nil {
		panic(exceptions.Exception{
			Tag: exceptions.FORBIDDEN,
		})
	}

	if len(poorGuysid) == 1 && poorGuysid[0] == commanderid {
		ts.LeaveTeam(commanderid, teamid)
		return
	}

	if !utils.Contains(member.Role.GetPermissionsForRole(), enums.REMOVE_MEMEBER) {
		panic(exceptions.Exception{
			Tag:    exceptions.FORBIDDEN,
			Errors: []exceptions.SpecificError{exceptions.LACKS_PERMISSION},
		})
	}

	for _, poorGuyid := range poorGuysid {
		if commanderid == poorGuyid {
			isLeaving = true
			continue
		}

		poorMember := ts.TeamRepo.GetMemberForTeam(teamid, poorGuyid)
		if poorMember == nil {
			continue
		}

		if !member.Role.DoesHavePowerOver(poorMember.Role) {
			continue
		}

		err := ts.TeamRepo.RemoveMember(poorGuyid, teamid)
		if err != nil {
			log.Println("KickMemberError: coudldn't kick member", poorGuyid, "by", commanderid, "at team", teamid, "detail:", err)
		}
	}
	if isLeaving {
		ts.LeaveTeam(commanderid, teamid)
	}
}

func (ts *TeamService) UpdateMemeberRole(commanderid int, info *dto.UpdateMemberRoleDto) {
	if commanderid < 0 {
		panic(exceptions.Exception{
			Tag: exceptions.UNAUTHORIZED,
			Errors: []exceptions.SpecificError{
				exceptions.AUTH_ACCESS_DENIED,
			},
		})
	}

	role := enums.NameToRole(info.Role)

	if role == 0 {
		panic(exceptions.Exception{
			Tag: exceptions.BAD_REQUEST,
		})
	} else if role == enums.TEAM_OWNER {
		panic(exceptions.Exception{
			Tag: exceptions.FORBIDDEN,
		})
	}

	member := ts.TeamRepo.GetMemberForTeam(info.TeamID, commanderid)
	if member == nil {

		log.Println("MemberError: commander", commanderid, "is not in team", info.TeamID)
		panic(exceptions.Exception{
			Tag: exceptions.FORBIDDEN,
		})
	}

	if !utils.Contains(member.Role.GetPermissionsForRole(), enums.EDIT_ROLE) {
		panic(exceptions.Exception{
			Tag:    exceptions.FORBIDDEN,
			Errors: []exceptions.SpecificError{exceptions.LACKS_PERMISSION},
		})
	}

	targetMember := ts.TeamRepo.GetMemberForTeam(info.TeamID, info.UserID)
	if targetMember == nil {
		panic(exceptions.Exception{
			Tag: exceptions.CONFLICT_ERROR,
		})
	}
	if !member.Role.DoesHavePowerOver(targetMember.Role) {
		log.Println("RoleError: user", member.Info.Userid, "doesn't have power over user", targetMember.Info.Userid, "in team", info.TeamID)
		panic(exceptions.Exception{
			Tag: exceptions.FORBIDDEN,
		})
	}

	err := ts.TeamRepo.UpdateMemberRole(role, info.UserID, info.TeamID)
	if err != nil {
		panic(exceptions.Exception{
			Tag: exceptions.CONFLICT_ERROR,
		})
	}
}

func (ts *TeamService) UpdatePosition(commanderid int, req *dto.UpdateMemberPositionDto) {
	if commanderid < 0 {
		panic(exceptions.Exception{
			Tag: exceptions.UNAUTHORIZED,
			Errors: []exceptions.SpecificError{
				exceptions.AUTH_ACCESS_DENIED,
			},
		})
	}

	member := ts.TeamRepo.GetMemberForTeam(req.Teamid, commanderid)
	if member == nil {
		panic(exceptions.Exception{
			Tag: exceptions.FORBIDDEN,
			Errors: []exceptions.SpecificError{
				exceptions.AUTH_ACCESS_DENIED,
			},
		})
	}

	if !utils.Contains(member.Role.GetPermissionsForRole(), enums.EDIT_NICKNAME) {
		panic(exceptions.Exception{
			Tag:    exceptions.FORBIDDEN,
			Errors: []exceptions.SpecificError{exceptions.LACKS_PERMISSION},
		})
	}

	targetMember := ts.TeamRepo.GetMemberForTeam(req.Teamid, req.Userid)
	if targetMember == nil {
		panic(exceptions.Exception{
			Tag: exceptions.CONFLICT_ERROR,
		})
	}
	if !member.Role.DoesHavePowerOver(targetMember.Role) {
		panic(exceptions.Exception{
			Tag: exceptions.FORBIDDEN,
		})
	}

	err := ts.TeamRepo.UpdateMemberPosition(req)
	if err != nil {
		panic(exceptions.Exception{
			Tag: exceptions.INTERNAL_ERROR,
		})
	}

}

func (ts *TeamService) GetMembers(commanderid int, teamid int64) []dto.SendMemberDto {
	var res []dto.SendMemberDto

	commanderMember := ts.TeamRepo.GetMemberForTeam(teamid, commanderid)

	members := ts.TeamRepo.GetMembersForTeam(teamid)

	if commanderMember == nil {
		for _, member := range members {
			res = append(res, dto.SendMemberDto{
				Info:    member.Info,
				Profile: ts.FileService.GetProfilePhotoURL(member.Info.Userid, false),
			})
		}
	} else {
		for _, member := range members {
			res = append(res, dto.SendMemberDto{
				Info:    member.Info,
				Profile: ts.FileService.GetProfilePhotoURL(member.Info.Userid, false),
				Role:    member.Role.String(),
			})
		}
	}

	return res
}

func (ts *TeamService) UpdateTeamProfile(commanderid int, teamid int64, data []byte) {
	if commanderid < 0 {
		panic(exceptions.Exception{
			Tag: exceptions.UNAUTHORIZED,
			Errors: []exceptions.SpecificError{
				exceptions.AUTH_ACCESS_DENIED,
			},
		})
	}

	member := ts.TeamRepo.GetMemberForTeam(teamid, commanderid)
	if member == nil {
		panic(exceptions.Exception{
			Tag: exceptions.FORBIDDEN,
			Errors: []exceptions.SpecificError{
				exceptions.AUTH_ACCESS_DENIED,
			},
		})
	}

	if !utils.Contains(member.Role.GetPermissionsForRole(), enums.EDIT_INFO) {
		panic(exceptions.Exception{
			Tag:    exceptions.FORBIDDEN,
			Errors: []exceptions.SpecificError{exceptions.LACKS_PERMISSION},
		})
	}

	ts.FileService.UploadTeamProfilePhoto(data, teamid)
}

func (ts *TeamService) DeleteTeamProfile(commanderid int, teamid int64) {
	if commanderid < 0 {
		panic(exceptions.Exception{
			Tag: exceptions.UNAUTHORIZED,
			Errors: []exceptions.SpecificError{
				exceptions.AUTH_ACCESS_DENIED,
			},
		})
	}

	member := ts.TeamRepo.GetMemberForTeam(teamid, commanderid)
	if member == nil {
		panic(exceptions.Exception{
			Tag: exceptions.FORBIDDEN,
			Errors: []exceptions.SpecificError{
				exceptions.AUTH_ACCESS_DENIED,
			},
		})
	}

	if !utils.Contains(member.Role.GetPermissionsForRole(), enums.EDIT_INFO) {
		panic(exceptions.Exception{
			Tag:    exceptions.FORBIDDEN,
			Errors: []exceptions.SpecificError{exceptions.LACKS_PERMISSION},
		})
	}

	ts.FileService.DeleteTeamProfilePhoto(teamid)
}

func (ts *TeamService) GetOneManTeamID(userid int) int64 {
	teamid, err := ts.TeamRepo.GetOneManTeamID(userid)
	if err != nil {
		panic(exceptions.Exception{
			Tag:    exceptions.NOT_FOUND,
			Errors: []exceptions.SpecificError{exceptions.USER_NOT_FOUND},
		})
	}
	return teamid
}

func (ts *TeamService) SearchTeams(req *elasticmodel.SimpleQuerySearchReqDto) []map[string]any {

	request := elasticmodel.SearchRequest{
		Query: req.Query,
		Limit: req.Limit,
		Page:  req.Page,
		Types: []string{"teams"}, //don't care
	}

	if req.SortBy == "" {
		request.SortBy = "_score"
	} else {
		request.SortBy = req.SortBy
	}
	if req.Order == "" {
		request.Order = "desc"
	} else {
		request.Order = req.Order
	}

	res, err := ts.SearchRepo.SearchTeams(&request)
	if err != nil {
		log.Println(err)
		panic(exceptions.Exception{
			Tag: exceptions.INTERNAL_ERROR,
		})
	}

	for i := 0; i < len(res); i++ {
		res[i]["profile"] = ts.FileService.GetTeamProfilePhotoURL(res[i]["_id"].(int64), false)
	}

	return res
}
