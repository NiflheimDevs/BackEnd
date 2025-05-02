package repositories

import (
	"context"

	"github.com/niflheimdevs/backend/internal/application/dto"
	"github.com/niflheimdevs/backend/internal/domain/models"
	"github.com/niflheimdevs/backend/internal/domain/repositories/postgres/transaction"
)

type UserRepo interface {
	UpdateUserPassword(id int, password []byte) error
	FindUserByPhone(phone string) (*models.UserModel, error)
	FindUserByID(id int) (*models.UserModel, error)
	FindUserByUsername(username string) (*models.UserModel, error)
	FindUserByEmail(email string) (*models.UserModel, error)
	PostUser(ctx context.Context, tx transaction.Tx, phonenumber string, username string, password []byte) (int, error)
	UpdateUserData(userid int, userData *dto.UpdateUserDTO) (int64, error)
	UpdateUsername(userid int, newUsername string) (int64, error)
	UpdateEmail(userid int, email string) (int64, error)
	UpdatePhone(phonenumber string, userid string)
	DeleteUser(userid int)
}
