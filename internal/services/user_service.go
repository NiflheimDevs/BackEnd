package services

import (
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/niflheimdevs/backend/internal/models"
	"github.com/niflheimdevs/backend/internal/repositories"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	UserRepo repositories.UserRepo
}

func NewUserService(userRepo repositories.UserRepo) *UserService {
	return &UserService{
		UserRepo: userRepo,
	}
}

func (userService *UserService) AuthenticateUser(identifier string, password string) *models.UserModel {
	user, err := userService.UserRepo.FindUserByUsername(identifier)
	if err == pgx.ErrNoRows {
		user, err = userService.UserRepo.FindUserByEmail(identifier)
		if err == pgx.ErrNoRows {
			//username/email not found
			log.Print("username/email not found")
			return &models.UserModel{}
		}
	}

	if err != nil {
		panic(err)
	}

	err = bcrypt.CompareHashAndPassword(user.Password, []byte(password))

	if err != nil {
		panic(err)
	}

	return user
}
