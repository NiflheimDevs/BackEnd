package repositories

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/niflheimdevs/backend/internal/domain/models"
)

type PaymentRepo interface {
	GetUserTransactions(userID int, offset, limit int) ([]models.TransactionModel, error)
	GetBalance(userID int) (int64, error)
	Withdraw(ctx context.Context, tx pgx.Tx, userID int, amount int64, description string)
	Deposit(ctx context.Context, tx pgx.Tx, userID int, amount int64, description string)
	UpdateWallet(ctx context.Context, tx pgx.Tx, userID int, amount int64) error
}
