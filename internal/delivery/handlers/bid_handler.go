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
	Validator  *validator.Validate
	Constants  *bootstrap.Constants
	BidService services.BidService
}

func NewBidHandler(
	constants *bootstrap.Constants,
	validator *validator.Validate,
	bidService services.BidService,
) *BidHandler {
	return &BidHandler{
		Validator:  validator,
		Constants:  constants,
		BidService: bidService,
	}
}

func (bh *BidHandler) PutBid(w http.ResponseWriter, r *http.Request) {
	params := Validated[dto.BidInfo](bh.Validator, r)

	userid := r.Context().Value(bh.Constants.Context.UserID).(int)

	params.UserID = userid

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

	bh.BidService.UpdateBid(params)

	w.WriteHeader(http.StatusNoContent)
}
