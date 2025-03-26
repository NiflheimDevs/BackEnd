package services

import (
	"fmt"

	"github.com/niflheimdevs/backend/internal/enums"
	"github.com/niflheimdevs/backend/internal/exceptions"
	"github.com/niflheimdevs/backend/internal/repositories"
)

type PaymentService struct {
	PaymentRepo *repositories.PaymentRepo
}

func NewPaymentService(
	paymentRepo *repositories.PaymentRepo,
) *PaymentService {
	return &PaymentService{
		PaymentRepo: paymentRepo,
	}
}

func (paymentService *PaymentService) ProjectPayment(userID, projectID int, amount int64) {
	if userID == -1 || userID == -2 {
		panic(exceptions.Exception{
			Tag: enums.UNAUTHORIZED,
			Errors: []enums.SpecificError{
				enums.AUTH_ACCESS_DENIED,
			},
		})
	}

	balance, err := paymentService.PaymentRepo.GetBalance(userID)

	if err != nil {
		panic(exceptions.Exception{
			Tag: enums.NOT_FOUND,
			Errors: []enums.SpecificError{
				enums.USER_NOT_FOUND,
			},
		})
	}

	if balance < amount {
		panic(exceptions.Exception{
			Tag: enums.FORBIDDEN,
			Errors: []enums.SpecificError{
				enums.INSUFFICIENT_BALANCE,
			},
		})
	}

	description := fmt.Sprintf("Payment for Project %d", projectID)

	paymentService.PaymentRepo.AdminTransaction(userID, amount, description)
}
