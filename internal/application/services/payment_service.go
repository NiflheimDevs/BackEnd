package services

import (
	"context"

	"github.com/niflheimdevs/backend/internal/application/dto"
	"github.com/niflheimdevs/backend/internal/domain/repositories/postgres/transaction"
)

type PaymentService interface {
	GetUserBalance(userID int) int64
	GetUserTransactions(userID int, offset, limit int, sortBy, order string) (int, []dto.UserTransactionDTO)
	ProjectPayment(ctx context.Context, tx transaction.Tx, userID int, amount int64)
	Deposit(userID int, amount int64, description string)
	Withdraw(userID int, amount int64, description string)
}
