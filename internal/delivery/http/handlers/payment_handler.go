package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/niflheimdevs/backend/bootstrap"
	"github.com/niflheimdevs/backend/internal/application/services"
	"github.com/niflheimdevs/backend/internal/exceptions"
)

type PaymentHandler struct {
	PaymentService services.PaymentService
	Constants      *bootstrap.Constants
	Validator      *validator.Validate
}

func NewPaymentHandler(
	paymentService services.PaymentService,
	constants *bootstrap.Constants,
	validator *validator.Validate,
) *PaymentHandler {
	return &PaymentHandler{
		PaymentService: paymentService,
		Constants:      constants,
		Validator:      validator,
	}
}

func (paymentHandler *PaymentHandler) GetUserTransactions(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(paymentHandler.Constants.Context.UserID).(int)

	query := r.URL.Query()
	offset, _ := strconv.Atoi(query.Get("offset"))
	limit, _ := strconv.Atoi(query.Get("limit"))

	transactions := paymentHandler.PaymentService.GetUserTransactions(userID, offset, limit)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(transactions); err != nil {
		panic(exceptions.Exception{
			Tag:    exceptions.INTERNAL_ERROR,
			Errors: []exceptions.SpecificError{exceptions.CAST_ERROR},
		})
	}
}

func (paymentHandler *PaymentHandler) GetUserBalance(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(paymentHandler.Constants.Context.UserID).(int)

	balance := paymentHandler.PaymentService.GetUserBalance(userID)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(strconv.FormatInt(balance, 10)))
}

func (paymentHandler *PaymentHandler) Deposit(w http.ResponseWriter, r *http.Request) {
	type DepositParam struct {
		Amount      int64  `json:"amount" validate:"required"`
		Description string `json:"description"`
	}

	params := Validated[DepositParam](paymentHandler.Validator, r)

	userID := r.Context().Value(paymentHandler.Constants.Context.UserID).(int)

	paymentHandler.PaymentService.Deposit(userID, params.Amount, params.Description)

	w.WriteHeader(http.StatusNoContent)
}

func (paymentHandler *PaymentHandler) Withdraw(w http.ResponseWriter, r *http.Request) {
	type DepositParam struct {
		Amount      int64  `json:"amount" validate:"required"`
		Description string `json:"description"`
	}

	params := Validated[DepositParam](paymentHandler.Validator, r)

	userID := r.Context().Value(paymentHandler.Constants.Context.UserID).(int)

	paymentHandler.PaymentService.Withdraw(userID, params.Amount, params.Description)

	w.WriteHeader(http.StatusNoContent)
}
