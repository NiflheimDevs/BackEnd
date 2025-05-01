package repositoriesimpl

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/niflheimdevs/backend/internal/domain/exceptions"
	"github.com/niflheimdevs/backend/internal/domain/models"
	"github.com/niflheimdevs/backend/internal/domain/repositories/postgres/transaction"
)

type PaymentRepo struct {
	PG *pgxpool.Pool
}

func NewPaymentRepo(PG *pgxpool.Pool) *PaymentRepo {
	return &PaymentRepo{
		PG: PG,
	}
}

func (paymentRepo *PaymentRepo) GetTransactionCount(userID int) int {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var count int

	query := "SELECT COUNT(*) FROM transaction WHERE (from_user_id = $1 OR to_user_id = $1)"

	err := paymentRepo.PG.QueryRow(ctx, query, userID).Scan(&count)

	if err != nil {
		panic(exceptions.Exception{
			Tag: exceptions.INTERNAL_ERROR,
			Errors: []exceptions.SpecificError{
				exceptions.DATABASE_ERROR,
			},
		})
	}

	return count
}

func (paymentRepo *PaymentRepo) GetUserTransactions(userID int, offset, limit int, sortBy, order string) ([]models.TransactionModel, error) {
	var transactions []models.TransactionModel

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if sortBy == "" {
		sortBy = "date"
	}
	if order == "" {
		order = "desc"
	}

	query := "SELECT * FROM transaction WHERE (from_user_id = $1 OR to_user_id = $1)"
	query += " ORDER BY " + sortBy + " " + order + " OFFSET $2 LIMIT $3"

	result, err := paymentRepo.PG.Query(ctx, query, userID, offset, limit)

	if err == pgx.ErrNoRows {
		return nil, err
	}

	if err != nil {
		panic(exceptions.Exception{
			Tag: exceptions.INTERNAL_ERROR,
			Errors: []exceptions.SpecificError{
				exceptions.DATABASE_ERROR,
			},
		})
	}

	defer result.Close()

	for result.Next() {
		var transaction models.TransactionModel
		var date time.Time
		if err := result.Scan(&transaction.ID, &transaction.FromUser, &transaction.ToUser, &transaction.Amount, &date, &transaction.Description); err != nil {
			panic(exceptions.Exception{
				Tag: exceptions.INTERNAL_ERROR,
				Errors: []exceptions.SpecificError{
					exceptions.DATABASE_ERROR,
				},
			})
		}
		transaction.Date = date.Format("2006-01-02 15:04:05")
		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func (paymentRepo *PaymentRepo) GetBalance(userID int) (int64, error) {
	var balance int64

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "Select wallet from users where id = $1"
	err := paymentRepo.PG.QueryRow(ctx, query, userID).Scan(&balance)

	if err == pgx.ErrNoRows {
		return 0, err
	}

	if err != nil {
		panic(exceptions.Exception{
			Tag: exceptions.INTERNAL_ERROR,
			Errors: []exceptions.SpecificError{
				exceptions.DATABASE_ERROR,
			},
		})
	}

	return balance, nil
}

func (paymentRepo *PaymentRepo) Withdraw(ctx context.Context, tx transaction.Tx, userID int, amount int64, description string) {

	now := time.Now().Format("2006-01-02 15:04:05")

	query := "INSERT INTO transaction (from_user_id, to_user_id, amount, date, description) VALUES ($1, $2, $3, $4, $5)"
	_, err := tx.Exec(ctx, query, userID, 1, amount, now, description)

	if err != nil {
		panic(exceptions.Exception{
			Tag: exceptions.INTERNAL_ERROR,
			Errors: []exceptions.SpecificError{
				exceptions.DATABASE_ERROR,
			},
		})
	}
}

func (paymentRepo *PaymentRepo) Deposit(ctx context.Context, tx transaction.Tx, userID int, amount int64, description string) {
	now := time.Now().Format("2006-01-02 15:04:05")

	query := "INSERT INTO transaction (from_user_id, to_user_id, amount, date, description) VALUES ($1, $2, $3, $4, $5)"
	_, err := tx.Exec(ctx, query, 1, userID, amount, now, description)

	if err != nil {
		panic(exceptions.Exception{
			Tag: exceptions.INTERNAL_ERROR,
			Errors: []exceptions.SpecificError{
				exceptions.DATABASE_ERROR,
			},
		})
	}
}

func (paymentRepo *PaymentRepo) UpdateWallet(ctx context.Context, tx transaction.Tx, userID int, amount int64) error {

	query := "UPDATE users SET wallet = wallet + $1 WHERE id = $2"
	_, err := tx.Exec(ctx, query, amount, userID)

	if err == pgx.ErrNoRows {
		return err
	}

	if err != nil {
		panic(exceptions.Exception{
			Tag: exceptions.INTERNAL_ERROR,
			Errors: []exceptions.SpecificError{
				exceptions.DATABASE_ERROR,
			},
		})
	}

	return nil
}
