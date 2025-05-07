package servicesimpl

import (
	"context"
	"time"

	"github.com/niflheimdevs/backend/internal/application/dto"
	"github.com/niflheimdevs/backend/internal/domain/exceptions"
	repositories "github.com/niflheimdevs/backend/internal/domain/repositories/postgres"
	"github.com/niflheimdevs/backend/internal/domain/repositories/postgres/transaction"
)

type PaymentService struct {
	PaymentRepo repositories.PaymentRepo
	TxManager   transaction.TxManager
}

func NewPaymentService(
	paymentRepo repositories.PaymentRepo,
	txManager transaction.TxManager,
) *PaymentService {
	return &PaymentService{
		PaymentRepo: paymentRepo,
		TxManager:   txManager,
	}
}

func (paymentService *PaymentService) GetUserBalance(userID int) int64 {
	if userID == -1 || userID == -2 {
		panic(exceptions.Exception{
			Tag: exceptions.UNAUTHORIZED,
			Errors: []exceptions.SpecificError{
				exceptions.AUTH_ACCESS_DENIED,
			},
		})
	}

	balance, err := paymentService.PaymentRepo.GetBalance(userID)

	if err != nil {
		panic(exceptions.Exception{
			Tag: exceptions.NOT_FOUND,
			Errors: []exceptions.SpecificError{
				exceptions.USER_NOT_FOUND,
			},
		})
	}

	return balance
}

func (paymentService *PaymentService) GetUserTransactions(userID int, offset, limit int, sortBy, order string) (int, []dto.UserTransactionDTO) {
	if userID == -1 || userID == -2 {
		panic(exceptions.Exception{
			Tag: exceptions.UNAUTHORIZED,
			Errors: []exceptions.SpecificError{
				exceptions.AUTH_ACCESS_DENIED,
			},
		})
	}

	count := paymentService.PaymentRepo.GetTransactionCount(userID)

	transactions, err := paymentService.PaymentRepo.GetUserTransactions(userID, offset, limit, sortBy, order)

	if err != nil {
		panic(exceptions.Exception{
			Tag: exceptions.NOT_FOUND,
			Errors: []exceptions.SpecificError{
				exceptions.USER_NOT_FOUND,
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
				Amount:      transaction.Amount,
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

	return count, output
}

func (paymentService *PaymentService) ProjectPayment(ctx context.Context, tx transaction.Tx, userID int, amount int64) {
	if userID == -1 || userID == -2 {
		panic(exceptions.Exception{
			Tag: exceptions.UNAUTHORIZED,
			Errors: []exceptions.SpecificError{
				exceptions.AUTH_ACCESS_DENIED,
			},
		})
	}

	balance, err := paymentService.PaymentRepo.GetBalance(userID)

	if err != nil {
		panic(exceptions.Exception{
			Tag: exceptions.NOT_FOUND,
			Errors: []exceptions.SpecificError{
				exceptions.USER_NOT_FOUND,
			},
		})
	}

	if balance < amount {
		panic(exceptions.Exception{
			Tag: exceptions.FORBIDDEN,
			Errors: []exceptions.SpecificError{
				exceptions.INSUFFICIENT_BALANCE,
			},
		})
	}

	paymentService.PaymentRepo.Withdraw(ctx, tx, userID, amount, "")
	paymentService.PaymentRepo.UpdateWallet(ctx, tx, userID, -1*amount)
}

func (paymentService *PaymentService) Deposit(userID int, amount int64, description string) {
	if userID == -1 || userID == -2 {
		panic(exceptions.Exception{
			Tag: exceptions.UNAUTHORIZED,
			Errors: []exceptions.SpecificError{
				exceptions.AUTH_ACCESS_DENIED,
			},
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tx, err := paymentService.TxManager.Begin(ctx)
	if err != nil {
		panic(exceptions.Exception{
			Tag: exceptions.INTERNAL_ERROR,
			Errors: []exceptions.SpecificError{
				exceptions.DATABASE_ERROR,
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
			Tag: exceptions.NOT_FOUND,
			Errors: []exceptions.SpecificError{
				exceptions.USER_NOT_FOUND,
			},
		})
	}

	paymentService.PaymentRepo.Deposit(ctx, tx, userID, amount, description)

	if err := tx.Commit(ctx); err != nil {
		panic(exceptions.Exception{
			Tag: exceptions.INTERNAL_ERROR,
			Errors: []exceptions.SpecificError{
				exceptions.DATABASE_ERROR,
			},
		})
	}
}

func (paymentService *PaymentService) Withdraw(userID int, amount int64, description string) {
	if userID == -1 || userID == -2 {
		panic(exceptions.Exception{
			Tag: exceptions.UNAUTHORIZED,
			Errors: []exceptions.SpecificError{
				exceptions.AUTH_ACCESS_DENIED,
			},
		})
	}

	balance, err := paymentService.PaymentRepo.GetBalance(userID)

	if err != nil {
		panic(exceptions.Exception{
			Tag: exceptions.NOT_FOUND,
			Errors: []exceptions.SpecificError{
				exceptions.USER_NOT_FOUND,
			},
		})
	}

	if balance < amount {
		panic(exceptions.Exception{
			Tag: exceptions.FORBIDDEN,
			Errors: []exceptions.SpecificError{
				exceptions.INSUFFICIENT_BALANCE,
			},
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tx, err := paymentService.TxManager.Begin(ctx)
	if err != nil {
		panic(exceptions.Exception{
			Tag: exceptions.INTERNAL_ERROR,
			Errors: []exceptions.SpecificError{
				exceptions.DATABASE_ERROR,
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
			Tag: exceptions.INTERNAL_ERROR,
			Errors: []exceptions.SpecificError{
				exceptions.DATABASE_ERROR,
			},
		})
	}
}

func (paymentService *PaymentService) TransferMoney(ctx context.Context, tx transaction.Tx, fromUserID, toUserID int, amount int64, description string) {
	balance, _ := paymentService.PaymentRepo.GetBalance(fromUserID)

	if balance < amount {
		panic(exceptions.Exception{
			Tag: exceptions.FORBIDDEN,
			Errors: []exceptions.SpecificError{
				exceptions.INSUFFICIENT_BALANCE,
			},
		})
	}

	paymentService.PaymentRepo.UpdateWallet(ctx, tx, fromUserID, -1*amount)
	paymentService.PaymentRepo.UpdateWallet(ctx, tx, toUserID, amount)
	paymentService.PaymentRepo.TransferMoney(ctx, tx, fromUserID, toUserID, amount, description)
}
