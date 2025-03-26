package repositories

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/niflheimdevs/backend/internal/enums"
	"github.com/niflheimdevs/backend/internal/exceptions"
)

type PaymentRepo struct {
	PG *pgxpool.Pool
}

func NewPaymentRepo(PG *pgxpool.Pool) *PaymentRepo {
	return &PaymentRepo{
		PG: PG,
	}
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
			Tag: enums.INTERNAL_ERROR,
			Errors: []enums.SpecificError{
				enums.DATABASE_ERROR,
			},
		})
	}

	return balance, nil
}

func (paymentRepo *PaymentRepo) AdminTransaction(userID int, amount int64, description string) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	now := time.Now().Format("2006-01-02 15:04:05")

	query := "INSERT INTO transaction (from_user_id, to_user_id, amount, date, description) VALUES ($1, $2, $3, $4, $5)"
	_, err := paymentRepo.PG.Exec(ctx, query, userID, 1, amount, now, description)

	if err != nil {
		panic(exceptions.Exception{
			Tag: enums.INTERNAL_ERROR,
			Errors: []enums.SpecificError{
				enums.DATABASE_ERROR,
			},
		})
	}

	query = "UPDATE users SET wallet = wallet - $1 WHERE id = $2"
	_, err = paymentRepo.PG.Exec(ctx, query, amount, userID)

	if err != nil {
		panic(exceptions.Exception{
			Tag: enums.INTERNAL_ERROR,
			Errors: []enums.SpecificError{
				enums.DATABASE_ERROR,
			},
		})
	}
}
