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

	query := "SELECT id,username,password,firstname,lastname FROM users WHERE username = $1"

	err := repo.PG.QueryRow(ctx, query, username).Scan(&user.ID, &user.Username, &user.Password, &user.FirstName, &user.LastName)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (repo *UserRepo) FindUserByPhone(phonenumber string) (*models.UserModel, error) {
	var user models.UserModel

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT id,username,password,firstname,lastname FROM users WHERE phone = $1`

	err := repo.PG.QueryRow(ctx, query, phonenumber).Scan(&user.ID, &user.Username, &user.Password, &user.FirstName, &user.LastName)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (repo *UserRepo) PostUser(phonenumber string, username string, password []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	INSERT INTO users 
	(phone , username, password)
	VALUES ($1 , $2 , $3)`

	_, err := repo.PG.Exec(ctx, query, phonenumber, username, password)
	return err
}
