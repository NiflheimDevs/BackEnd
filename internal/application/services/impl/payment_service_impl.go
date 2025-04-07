package servicesimpl

import (
	"context"
	"time"

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

func (paymentService *PaymentService) GetUserBalance(userID int) int64 {
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

	return balance
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

	paymentService.PaymentRepo.Withdraw(ctx, tx, userID, amount, "")
	paymentService.PaymentRepo.UpdateWallet(ctx, tx, userID, -1*amount)
}

func (paymentService *PaymentService) Deposit(userID int, amount int64, description string) {
	if userID == -1 || userID == -2 {
		panic(exceptions.Exception{
			Tag: enums.UNAUTHORIZED,
			Errors: []enums.SpecificError{
				enums.AUTH_ACCESS_DENIED,
			},
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tx, err := paymentService.PaymentRepo.PG.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		panic(exceptions.Exception{
			Tag: enums.INTERNAL_ERROR,
			Errors: []enums.SpecificError{
				enums.DATABASE_ERROR,
			},
		})
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(ctx)
			panic(p)
		}
	}()

	err = paymentService.PaymentRepo.UpdateWallet(ctx, tx, userID, amount)

	if err != nil {
		panic(exceptions.Exception{
			Tag: enums.NOT_FOUND,
			Errors: []enums.SpecificError{
				enums.USER_NOT_FOUND,
			},
		})
	}

	paymentService.PaymentRepo.Deposit(ctx, tx, userID, amount, description)

	if err := tx.Commit(ctx); err != nil {
		panic(exceptions.Exception{
			Tag: enums.INTERNAL_ERROR,
			Errors: []enums.SpecificError{
				enums.DATABASE_ERROR,
			},
		})
	}
}

func (paymentService *PaymentService) Withdraw(userID int, amount int64, description string) {
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

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tx, err := paymentService.PaymentRepo.PG.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		panic(exceptions.Exception{
			Tag: enums.INTERNAL_ERROR,
			Errors: []enums.SpecificError{
				enums.DATABASE_ERROR,
			},
		})
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(ctx)
			panic(p)
		}
	}()

	paymentService.PaymentRepo.Withdraw(ctx, tx, userID, amount, description)
	paymentService.PaymentRepo.UpdateWallet(ctx, tx, userID, -1*amount)

	if err := tx.Commit(ctx); err != nil {
		panic(exceptions.Exception{
			Tag: enums.INTERNAL_ERROR,
			Errors: []enums.SpecificError{
				enums.DATABASE_ERROR,
			},
		})
	}
}
