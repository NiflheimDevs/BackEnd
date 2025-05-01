package handlers

import (
	"encoding/json"
	"net/http"

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
