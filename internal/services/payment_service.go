package services

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/niflheimdevs/backend/internal/dto"
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

func (paymentService *PaymentService) GetUserTransactions(userID int, offset, limit int) []dto.UserTransactionDTO {
	if userID == -1 || userID == -2 {
		panic(exceptions.Exception{
			Tag: enums.UNAUTHORIZED,
			Errors: []enums.SpecificError{
				enums.AUTH_ACCESS_DENIED,
			},
		})
	}

	transactions, err := paymentService.PaymentRepo.GetUserTransactions(userID, offset, limit)

	if err != nil {
		panic(exceptions.Exception{
			Tag: enums.NOT_FOUND,
			Errors: []enums.SpecificError{
				enums.USER_NOT_FOUND,
			},
		})
	}

	var output []dto.UserTransactionDTO

	for _, transaction := range transactions {
		if transaction.FromUser == userID {
			output = append(output, dto.UserTransactionDTO{
				Date:        transaction.Date,
				Type:        1,
				Description: transaction.Description,
				Amount:      -1 * transaction.Amount,
			})
		} else if transaction.ToUser == userID {
			output = append(output, dto.UserTransactionDTO{
				Date:        transaction.Date,
				Type:        2,
				Description: transaction.Description,
				Amount:      transaction.Amount,
			})
		}
	}

	return output
}

func (paymentService *PaymentService) ProjectPayment(ctx context.Context, tx pgx.Tx, userID int, amount int64) {
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

	paymentService.PaymentRepo.AdminTransaction(ctx, tx, userID, amount, "")
}
