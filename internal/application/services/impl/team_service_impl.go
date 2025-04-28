package servicesimpl

import (
	"context"
	"log"
	"time"

	"github.com/niflheimdevs/backend/internal/application/dto"
	"github.com/niflheimdevs/backend/internal/application/services"
	"github.com/niflheimdevs/backend/internal/domain/enums"
	"github.com/niflheimdevs/backend/internal/domain/exceptions"
	repositories "github.com/niflheimdevs/backend/internal/domain/repositories/postgres"
	"github.com/niflheimdevs/backend/internal/domain/repositories/postgres/transaction"
	"github.com/niflheimdevs/backend/internal/utils"
)

type TeamService struct {
	TeamRepo           repositories.TeamRepo
	RoleRepo           repositories.RoleRepo
	TransactionManager transaction.TxManager
	FileService        services.FileService
}

func NewTeamService(
	teamRepo repositories.TeamRepo,
	roleRepo repositories.RoleRepo,
	tManager transaction.TxManager,
	fileService services.FileService,
) *TeamService {
	return &TeamService{
		TeamRepo:           teamRepo,
		RoleRepo:           roleRepo,
		TransactionManager: tManager,
		FileService:        fileService,
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

	err = ts.TeamRepo.AddMemberWithTx(ctx, tx, userid, teamid, "master", enums.TEAM_OWNER)
	if err != nil {
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
	addMembersFunc(teamid, teamInfo.Members, ts)

	return teamid
}

func (ts *TeamService) GetTeamsForUser(userid int) []dto.GetTeamPreviewDto {
	return ts.TeamRepo.GetTeamsForUser(userid)
}

func (ts *TeamService) GetTeam(commanderid int, teamid int64) *dto.GetTeamDto {
	var res dto.GetTeamDto

	res.Info = ts.TeamRepo.GetTeam(teamid)

	members := ts.TeamRepo.GetMembersForTeam(teamid)
	res.Members = make([]dto.SendMemberDto, len(members))
	for i := 0; i < len(members); i++ {
		res.Members[i].Info = members[i].Info
		res.Members[i].Role = members[i].Role.String()
		//TODO: profile
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

func (ts *TeamService) UpdateTeamInfo(userid int, info *dto.UpdateTeamInfoDto) {

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

}

// func (ts *TeamService) GetTeamWithRole(userid int, teamid int64) *dto.GetTeamPreviewDto {
// 	var res dto.GetTeamDto

// 	res.Info = ts.TeamRepo.GetTeamForUser(userid, teamid)
// 	member := ts.TeamRepo.GetMemberForTeam(teamid, userid)
// 	if member != nil {
// 		res.Role = member.Role.String()
// 	}

// 	return &res
// }

// TODO: email? some sort of request must be sent and then when it is accepted, the member gets added
// ! this version is naive
func (ts *TeamService) AddMembers(userid int, teamid int64, members []int) {
	member := ts.TeamRepo.GetMemberForTeam(teamid, userid)
	if member == nil {
		panic(exceptions.Exception{
			Tag: exceptions.FORBIDDEN,
		})
	}

	//TODO: if super admin can create roles then it should be a query instead of static call
	if !utils.Contains(member.Role.GetPermissionsForRole(), enums.ADD_MEMBER) {
		panic(exceptions.Exception{
			Tag:    exceptions.FORBIDDEN,
			Errors: []exceptions.SpecificError{exceptions.LACKS_PERMISSION},
		})
	}

	addMembersFunc(teamid, members, ts)

}

func addMembersFunc(teamid int64, members []int, ts *TeamService) {
	for _, id := range members {
		ts.TeamRepo.AddMember(id, teamid, "", enums.TEAM_NEWBIE)
	}
}

func (ts *TeamService) LeaveTeam(userid int, teamid int64) {
	err := ts.TeamRepo.RemoveMember(userid, teamid)
	if err != nil {
		panic(exceptions.Exception{
			Tag: exceptions.CONFLICT_ERROR,
		})
	}
}

func (ts *TeamService) KickMemebr(commanderid int, poorGuyid int, teamid int64) {
	if commanderid == poorGuyid {
		ts.LeaveTeam(commanderid, teamid)
		return
	}

	member := ts.TeamRepo.GetMemberForTeam(teamid, commanderid)
	if member == nil {
		panic(exceptions.Exception{
			Tag: exceptions.FORBIDDEN,
		})
	}

	if !utils.Contains(member.Role.GetPermissionsForRole(), enums.REMOVE_MEMEBER) {
		panic(exceptions.Exception{
			Tag:    exceptions.FORBIDDEN,
			Errors: []exceptions.SpecificError{exceptions.LACKS_PERMISSION},
		})
	}

	err := ts.TeamRepo.RemoveMember(poorGuyid, teamid)
	if err != nil {
		panic(exceptions.Exception{
			Tag: exceptions.INTERNAL_ERROR,
		})
	}

}

func (ts *TeamService) UpdateMemeberRole(commanderid int, info *dto.UpdateMemberRoleDto) {
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

	err := ts.TeamRepo.UpdateMemberRole(role, info.UserID, info.TeamID)
	if err != nil {
		panic(exceptions.Exception{
			Tag: exceptions.CONFLICT_ERROR,
		})
	}
}

func (ts *TeamService) UpdatePosition(commanderid int, req *dto.UpdateMemberPositionDto) {
	member := ts.TeamRepo.GetMemberForTeam(req.Teamid, commanderid)
	if member == nil {
		panic(exceptions.Exception{
			Tag: exceptions.FORBIDDEN,
		})
	}

	if !utils.Contains(member.Role.GetPermissionsForRole(), enums.EDIT_NICKNAME) {
		panic(exceptions.Exception{
			Tag:    exceptions.FORBIDDEN,
			Errors: []exceptions.SpecificError{exceptions.LACKS_PERMISSION},
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
				Profile: ts.FileService.GetUserProfileName(member.Info.Userid, false),
			})
		}
	} else {
		for _, member := range members {
			res = append(res, dto.SendMemberDto{
				Info:    member.Info,
				Profile: ts.FileService.GetUserProfileName(member.Info.Userid, false),
				Role:    member.Role.String(),
			})
		}
	}

	return res
}
