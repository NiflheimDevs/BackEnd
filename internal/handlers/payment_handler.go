package handlers

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/niflheimdevs/backend/internal/bootstrap"
	"github.com/niflheimdevs/backend/internal/services"
)

type PaymentHandler struct {
	PaymentService *services.PaymentService
	Constants      *bootstrap.Constants
	Validator      *validator.Validate
}

func NewPaymentHandler(
	paymentService *services.PaymentService,
	constants *bootstrap.Constants,
	validator *validator.Validate,
) *PaymentHandler {
	return &PaymentHandler{
		PaymentService: paymentService,
		Constants:      constants,
		Validator:      validator,
	}
}

func (paymentHandler *PaymentHandler) ProjectPayment(w http.ResponseWriter, r *http.Request) {
	type projectPayment struct {
		Amount    int `json:"amount" validate:"required"`
		ProjectID int `json:"project_id" validate:"required"`
	}

	param := Validated[projectPayment](paymentHandler.Validator, r)

	userID := r.Context().Value(paymentHandler.Constants.Context.UserID).(int)

	paymentHandler.PaymentService.ProjectPayment(userID, param.Amount, param.ProjectID)

	w.WriteHeader(http.StatusNoContent)
}
