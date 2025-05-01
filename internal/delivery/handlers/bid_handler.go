package handlers

import (
	"encoding/json"
	"net/http"
	"time"

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
	type BidParams struct {
		TeamID       int       `json:"team_id" validate:"required"`
		ProjectID    int       `json:"project_id" validate:"required"`
		PrePayment   int64     `json:"pre_payment" validate:"required"`
		Total        int64     `json:"total" validate:"required"`
		Description  string    `json:"description" validate:""`
		ExpectedTime time.Time `json:"expected_time" validate:"required"`
	}

	params := Validated[BidParams](bh.Validator, r)

	userid := r.Context().Value(bh.Constants.Context.UserID).(int)

	bidid := bh.BidService.PutBidOnProject(userid, params.TeamID, params.ProjectID, params.PrePayment, params.Total, params.Description, params.ExpectedTime)

	id := dto.PutBid{
		BidID: bidid,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(id)
}
