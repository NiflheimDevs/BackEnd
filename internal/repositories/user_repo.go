package repositories

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/niflheimdevs/backend/internal/models"
)

type UserRepo struct {
	PG *pgxpool.Pool
}

func NewUserRepo(PG *pgxpool.Pool) *UserRepo {
	return &UserRepo{
		PG: PG,
	}
}

func (repo *UserRepo) FindUserByUsername(username string) (*models.UserModel, error) {
	var user models.UserModel

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	//query := "SELECT * FROM user WHERE username = $1"

	err := repo.PG.Ping(ctx)
	if err != nil {
		panic(err)
	}
	return &user, nil
}
