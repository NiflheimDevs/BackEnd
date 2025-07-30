package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/niflheimdevs/backend/bootstrap"
	"github.com/niflheimdevs/backend/internal/application/dto"
	"github.com/niflheimdevs/backend/internal/application/services"
	"github.com/niflheimdevs/backend/internal/domain/exceptions"
	elasticmodel "github.com/niflheimdevs/backend/internal/domain/models/elastic"
)

type TeamHandler struct {
	TeamService services.TeamService
	Constants   *bootstrap.Constants
	Validator   *validator.Validate
}

func NewTeamHandler(
	teamService services.TeamService,
	constants *bootstrap.Constants,
	validator *validator.Validate,
) *TeamHandler {
	return &TeamHandler{
		TeamService: teamService,
		Constants:   constants,
		Validator:   validator,
	}
}

func (th *TeamHandler) CreateTeam(w http.ResponseWriter, r *http.Request) {
	params := Validated[dto.TeamCreateDto](th.Validator, r)

	userid, _ := r.Context().Value(th.Constants.Context.UserID).(int)

	type Response struct {
		TeamID int64 `json:"teamid"`
	}

	var res Response

	res.TeamID = th.TeamService.CreateTeam(userid, &params)

	// TODO: send smthg to front. sobhan is pain
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (th *TeamHandler) UpdateTeamInfo(w http.ResponseWriter, r *http.Request) {
	params := Validated[dto.UpdateTeamInfoDto](th.Validator, r)

	userid, _ := r.Context().Value(th.Constants.Context.UserID).(int)

	th.TeamService.UpdateTeamInfo(userid, &params)

	// TODO: send smthg to front. sobhan is pain
	w.WriteHeader(http.StatusNoContent)
}

func (th *TeamHandler) DeleteTeam(w http.ResponseWriter, r *http.Request) {
	teamidString := chi.URLParam(r, "team_id")
	teamid, err := strconv.Atoi(teamidString)

	if err != nil {
		panic(exceptions.Exception{
			Tag: exceptions.BAD_REQUEST,
		})
	}
	userid, _ := r.Context().Value(th.Constants.Context.UserID).(int)

	th.TeamService.DeleteTeam(userid, int64(teamid))

	w.WriteHeader(http.StatusNoContent)
}

func (th *TeamHandler) GetTeamsForUser(w http.ResponseWriter, r *http.Request) {
	useridString := chi.URLParam(r, "user_id")
	userid, err := strconv.Atoi(useridString)

	if err != nil {
		panic(exceptions.Exception{
			Tag: exceptions.BAD_REQUEST,
		})
	}
	query := r.URL.Query()
	active, err := strconv.Atoi(query.Get("active"))
	if err != nil {
		active = 1
	}
	dontCare, err := strconv.Atoi(query.Get("dontcare"))
	if err != nil {
		dontCare = 0
	}
	commanderid, _ := r.Context().Value(th.Constants.Context.UserID).(int)

	if userid == 0 {
		userid = commanderid
	}

	res := th.TeamService.GetTeamsForUser(userid, active, dontCare)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (th *TeamHandler) GetTeamsForBidding(w http.ResponseWriter, r *http.Request) {
	userid, _ := r.Context().Value(th.Constants.Context.UserID).(int)
	res := th.TeamService.GetTeamsForBidding(userid)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (th *TeamHandler) GetTeam(w http.ResponseWriter, r *http.Request) {
	teamidString := chi.URLParam(r, "team_id")
	teamid, err := strconv.Atoi(teamidString)

	if err != nil {
		panic(exceptions.Exception{
			Tag: exceptions.BAD_REQUEST,
		})
	}
	commanderid, _ := r.Context().Value(th.Constants.Context.UserID).(int)

	res := th.TeamService.GetTeam(commanderid, int64(teamid))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (th *TeamHandler) UpdateMemberPosition(w http.ResponseWriter, r *http.Request) {
	params := Validated[dto.UpdateMemberPositionDto](th.Validator, r)

	userid, _ := r.Context().Value(th.Constants.Context.UserID).(int)

	th.TeamService.UpdatePosition(userid, &params)

	w.WriteHeader(http.StatusNoContent)
}

func (th *TeamHandler) RemoveMember(w http.ResponseWriter, r *http.Request) {
	type Params struct {
		Members []int `json:"members" validate:"required"`
		TeamID  int64 `json:"team_id" validate:"required,numeric"`
	}

	params := Validated[Params](th.Validator, r)

	userid, _ := r.Context().Value(th.Constants.Context.UserID).(int)

	for i := 0; i < len(params.Members); i++ {
		if params.Members[i] == 0 {
			params.Members[i] = userid
		}
	}

	// if len(params.Members) == 1 && params.Members[0] == 0 {
	// 	params.Members[0] = userid
	// }

	th.TeamService.KickMemebr(userid, params.Members, params.TeamID)

	w.WriteHeader(http.StatusNoContent)
}
func (th *TeamHandler) AddMembers(w http.ResponseWriter, r *http.Request) {
	type Param struct {
		TeamID     int64 `json:"team_id" validate:"required,numeric"`
		NewMembers []int `json:"members" validate:"required"`
	}

	params := Validated[Param](th.Validator, r)

	userid, _ := r.Context().Value(th.Constants.Context.UserID).(int)

	th.TeamService.InviteMembers(userid, params.TeamID, params.NewMembers)

	w.WriteHeader(http.StatusNoContent)
}

func (th *TeamHandler) UpdateMemberRole(w http.ResponseWriter, r *http.Request) {

	params := Validated[dto.UpdateMemberRoleDto](th.Validator, r)

	userid, _ := r.Context().Value(th.Constants.Context.UserID).(int)

	th.TeamService.UpdateMemeberRole(userid, &params)

	w.WriteHeader(http.StatusNoContent)
}

func (th *TeamHandler) SearchTeams(w http.ResponseWriter, r *http.Request) {

	// params := Validated[elasticmodel.SimpleQuerySearchReqDto](th.Validator, r)
	q := r.URL.Query()

	page, err := strconv.Atoi(q.Get("page"))
	if err != nil {
		panic(exceptions.Exception{
			Tag: exceptions.BAD_REQUEST,
		})
	}

	limit, err := strconv.Atoi(q.Get("limit"))
	if err != nil {
		panic(exceptions.Exception{
			Tag: exceptions.BAD_REQUEST,
		})
	}

	params := &elasticmodel.SimpleQuerySearchReqDto{
		Query:  q.Get("query"),
		Page:   page,
		Limit:  limit,
		SortBy: q.Get("sort_by"),
		Order:  q.Get("order"),
	}

	if params.Order != "" && params.Order != "desc" && params.Order != "asc" {
		panic(exceptions.Exception{
			Tag: exceptions.BAD_REQUEST,
		})
	}

	res := th.TeamService.SearchTeams(params)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (th *TeamHandler) AcceptInvite(w http.ResponseWriter, r *http.Request) {
	type Params struct {
		// idk why not getting userid as well
		TeamID int64  `json:"team_id" validate:"required,numeric"`
		Token  string `json:"token" validate:"required"`
	}

	params := Validated[Params](th.Validator, r)

	th.TeamService.AcceptInvite(params.Token, params.TeamID)

	w.WriteHeader(http.StatusNoContent)
}
