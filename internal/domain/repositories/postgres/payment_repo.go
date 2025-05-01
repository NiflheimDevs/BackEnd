package repositories

import (
	"context"

	"github.com/niflheimdevs/backend/internal/domain/models"
	"github.com/niflheimdevs/backend/internal/domain/repositories/postgres/transaction"
)

type PaymentRepo interface {
	GetUserTransactions(userID int, offset, limit int, sortBy, order string) ([]models.TransactionModel, error)
	GetBalance(userID int) (int64, error)
	Withdraw(ctx context.Context, tx transaction.Tx, userID int, amount int64, description string)
	Deposit(ctx context.Context, tx transaction.Tx, userID int, amount int64, description string)
	UpdateWallet(ctx context.Context, tx transaction.Tx, userID int, amount int64) error
	GetTransactionCount(userID int) int
}
