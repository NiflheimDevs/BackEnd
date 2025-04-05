package repositories

import (
	"context"
	"database/sql"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/niflheimdevs/backend/internal/dto"
	"github.com/niflheimdevs/backend/internal/enums"
	"github.com/niflheimdevs/backend/internal/exceptions"
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

func fillUserModel(row pgx.Row) (*models.UserModel, error) {
	var user models.UserModel

	var firstname, lastname, bio, email sql.NullString
	err := row.Scan(&user.ID, &user.Username, &user.Password, &firstname, &lastname, &bio, &email, &user.Is_verified, &user.Phone, &user.Wallet)
	if err != nil {
		return nil, err
	}
	if firstname.Valid {
		user.FirstName = firstname.String
	}
	if lastname.Valid {
		user.LastName = lastname.String
	}
	if bio.Valid {
		user.Bio = bio.String
	}
	if email.Valid {
		user.Email = email.String
	}
	return &user, nil
}

func (repo *UserRepo) UpdateUserPassword(id int, password []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "UPDATE users SET password=$1 WHERE id = $2"

	_, err := repo.PG.Exec(ctx, query, password, id)

	return err
}

func (repo *UserRepo) FindUserByPhone(phone string) (*models.UserModel, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "SELECT id,username,password,firstname,lastname,bio,email,is_verified,phone,wallet FROM users WHERE phone = $1"

	row := repo.PG.QueryRow(ctx, query, phone)

	return fillUserModel(row)
}

func (repo *UserRepo) FindUserByID(id int) (*models.UserModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := "SELECT id,username,password,firstname,lastname,bio,email,is_verified,phone,wallet FROM users WHERE id = $1"

	row := repo.PG.QueryRow(ctx, query, id)

	return fillUserModel(row)

}

func (repo *UserRepo) FindUserByUsername(username string) (*models.UserModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "SELECT id,username,password,firstname,lastname,bio,email,is_verified,phone,wallet FROM users WHERE username = $1"

	row := repo.PG.QueryRow(ctx, query, username)

	return fillUserModel(row)
}

func (repo *UserRepo) FindUserByEmail(email string) (*models.UserModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "SELECT id,username,password,firstname,lastname,bio,email,is_verified,phone,wallet FROM users WHERE email = $1"

	row := repo.PG.QueryRow(ctx, query, email)

	return fillUserModel(row)
}

func (repo *UserRepo) PostUser(phonenumber string, username string, password []byte) (int, error) {
	var userid int

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	INSERT INTO users 
	(phone , username, password)
	VALUES ($1 , $2 , $3)
	RETURNING id`

	err := repo.PG.QueryRow(ctx, query, phonenumber, username, password).Scan(&userid)
	return userid, err
}

func (repo *UserRepo) UpdateUserData(userid int, userData *dto.UpdateUserDTO) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	UPDATE users
	SET firstname = $2 , lastname = $3, bio = $4
	WHERE id = $1`

	res, err := repo.PG.Exec(ctx, query, userid, userData.FirstName, userData.LastName, userData.Bio)
	return res.RowsAffected(), err
}

func (repo *UserRepo) UpdateUsername(userid int, newUsername string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	UPDATE users
	SET username = $2
	WHERE id = $1`

	res, err := repo.PG.Exec(ctx, query, userid, newUsername)
	return res.RowsAffected(), err
}

func (repo *UserRepo) UpdateEmail(userid int, email string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	UPDATE users
	SET email = $2, is_verified = false
	WHERE id = $1 AND is_verified = true`

	res, err := repo.PG.Exec(ctx, query, userid, email)
	//? better error handling for internal errors?
	return res.RowsAffected(), err
}

func (repo *UserRepo) UpdatePhone(phonenumber string, userid string) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	UPDATE users
	SET phone = $2
	WHERE id = $1`

	_, err := repo.PG.Exec(ctx, query, userid, phonenumber)

	if err != nil {
		panic(exceptions.Exception{
			Tag: enums.INTERNAL_ERROR,
		})
	}
}

func (repo *UserRepo) DeleteUser(userid int) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	DELETE FROM users
	WHERE id = $1`
	_, err := repo.PG.Exec(ctx, query, userid)

	if err != nil {
		panic(exceptions.Exception{
			Tag: enums.INTERNAL_ERROR,
		})
	}
}
