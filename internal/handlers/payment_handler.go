package handlers

import (
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
