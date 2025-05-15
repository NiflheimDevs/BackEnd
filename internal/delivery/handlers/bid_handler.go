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
)

type BidHandler struct {
	Validator   *validator.Validate
	Constants   *bootstrap.Constants
	BidService  services.BidService
	TeamService services.TeamService
}

func NewBidHandler(
	constants *bootstrap.Constants,
	validator *validator.Validate,
	bidService services.BidService,
	teamService services.TeamService,
) *BidHandler {
	return &BidHandler{
		Validator:   validator,
		Constants:   constants,
		BidService:  bidService,
		TeamService: teamService,
	}
}

func (bh *BidHandler) PutBid(w http.ResponseWriter, r *http.Request) {
	params := Validated[dto.BidInfo](bh.Validator, r)

	userID := r.Context().Value(bh.Constants.Context.UserID).(int)

	params.UserID = userID

	bidid := bh.BidService.PutBidOnProject(params)

	id := dto.PutBid{
		BidID: bidid,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(id)
}

func (bh *BidHandler) AcceptBid(w http.ResponseWriter, r *http.Request) {
	type AcceptBidParam struct {
		ProjectID int `json:"project_id" validate:"required"`
		BidID     int `json:"bid_id" validate:""`
	}

	params := Validated[AcceptBidParam](bh.Validator, r)

	bidIDString := chi.URLParam(r, "id")
	bidID, _ := strconv.Atoi(bidIDString)

	params.BidID = bidID

	userid := r.Context().Value(bh.Constants.Context.UserID).(int)

	bh.BidService.AcceptBid(userid, params.BidID, params.ProjectID)

	w.WriteHeader(http.StatusNoContent)
}

func (bh *BidHandler) UpdateBid(w http.ResponseWriter, r *http.Request) {
	params := Validated[dto.BidInfo](bh.Validator, r)

	bidIDString := chi.URLParam(r, "id")
	bidID, _ := strconv.Atoi(bidIDString)

	params.BidID = bidID

	userID := r.Context().Value(bh.Constants.Context.UserID).(int)

	params.UserID = userID

	bh.BidService.UpdateBid(params)

	w.WriteHeader(http.StatusNoContent)
}

func (bh *BidHandler) GetProjectBids(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(bh.Constants.Context.UserID).(int)

	projectIDString := chi.URLParam(r, "project_id")
	projectID, _ := strconv.Atoi(projectIDString)

	bids := bh.BidService.GetPublicProjectBids(projectID)

	teamIDs := bh.TeamService.GetAllTeamsIDs(userID)

	result := dto.ProjectBidInfo{
		Bids:    bids,
		TeamIDs: teamIDs,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func (bh *BidHandler) ViewBidsOfTheProject(w http.ResponseWriter, r *http.Request) {
	projectIDString := chi.URLParam(r, "project_id")
	projectID, _ := strconv.Atoi(projectIDString)
	userID := r.Context().Value(bh.Constants.Context.UserID).(int)

	result := bh.BidService.GetPrivateProjectBids(userID, projectID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func (bh *BidHandler) GetTeamBids(w http.ResponseWriter, r *http.Request) {
	teamIDString := chi.URLParam(r, "team_id")
	teamID, _ := strconv.Atoi(teamIDString)

	userID := r.Context().Value(bh.Constants.Context.UserID).(int)

	result := bh.BidService.GetTeamBids(userID, int64(teamID))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}
