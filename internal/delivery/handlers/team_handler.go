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

	th.TeamService.CreateTeam(userid, &params)
	// TODO: send smthg to front. sobhan is pain

	w.WriteHeader(http.StatusNoContent)
}

func (th *TeamHandler) UpdateTeamInfo(w http.ResponseWriter, r *http.Request) {
	params := Validated[dto.UpdateTeamInfoDto](th.Validator, r)

	userid, _ := r.Context().Value(th.Constants.Context.UserID).(int)

	th.TeamService.UpdateTeamInfo(userid, &params)

	// TODO: send smthg to front. sobhan is pain
	w.WriteHeader(http.StatusNoContent)
}

func (th *TeamHandler) DeleteTeam(w http.ResponseWriter, r *http.Request) {
	type Param struct {
		TeamID int64 `json:"team_id" validate:"required"`
	}
	params := Validated[Param](th.Validator, r)

	userid, _ := r.Context().Value(th.Constants.Context.UserID).(int)

	th.TeamService.DeleteTeam(userid, params.TeamID)

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
	commanderid, _ := r.Context().Value(th.Constants.Context.UserID).(int)

	if userid == 0 {
		userid = commanderid
	}

	res := th.TeamService.GetTeamsForUser(userid)

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
		UserID int   `json:"user_id" validate:"required,numeric"`
		TeamID int64 `json:"team_id" validate:"required,numeric"`
	}

	params := Validated[Params](th.Validator, r)

	userid, _ := r.Context().Value(th.Constants.Context.UserID).(int)

	th.TeamService.KickMemebr(userid, params.UserID, params.TeamID)

	w.WriteHeader(http.StatusNoContent)
}

func (th *TeamHandler) LeaveTeam(w http.ResponseWriter, r *http.Request) {
	userid, _ := r.Context().Value(th.Constants.Context.UserID).(int)
	teamidString := chi.URLParam(r, "team-id")
	teamid, err := strconv.Atoi(teamidString)
	if err != nil {
		panic(exceptions.Exception{
			Tag: exceptions.BAD_REQUEST,
		})
	}

	th.TeamService.LeaveTeam(userid, int64(teamid))

	w.WriteHeader(http.StatusNoContent)
}

func (th *TeamHandler) AddMembers(w http.ResponseWriter, r *http.Request) {
	type Param struct {
		TeamID     int64 `json:"team_id" validate:"required,numeric"`
		NewMembers []int `json:"new_members" validate:"required"`
	}

	params := Validated[Param](th.Validator, r)

	userid, _ := r.Context().Value(th.Constants.Context.UserID).(int)

	th.TeamService.AddMembers(userid, params.TeamID, params.NewMembers)

	w.WriteHeader(http.StatusNoContent)
}

func (th *TeamHandler) EditPosition(w http.ResponseWriter, r *http.Request) {

	params := Validated[dto.UpdateMemberPositionDto](th.Validator, r)

	userid, _ := r.Context().Value(th.Constants.Context.UserID).(int)

	th.TeamService.UpdatePosition(userid, &params)

	w.WriteHeader(http.StatusNoContent)
}

func (th *TeamHandler) EditRole(w http.ResponseWriter, r *http.Request) {

	params := Validated[dto.UpdateMemberRoleDto](th.Validator, r)

	userid, _ := r.Context().Value(th.Constants.Context.UserID).(int)

	th.TeamService.UpdateMemeberRole(userid, &params)

	w.WriteHeader(http.StatusNoContent)
}
