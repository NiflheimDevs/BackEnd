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
		user, err = userService.UserRepo.FindUserByPhone(identifier)
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

func (userService *UserService) ChangePassword(user_id int, old_password string, new_password string) {
	user, err := userService.UserRepo.FindUserByID(user_id)

	if err != nil {
		panic(err)
	}

	err = bcrypt.CompareHashAndPassword(user.Password, []byte(old_password))

	if err != nil {
		panic(err)
	}

	password, err := bcrypt.GenerateFromPassword([]byte(new_password), 15)

	if err != nil {
		panic(err)
	}

	err = userService.UserRepo.UpdateUserPassword(user_id, password)

	if err != nil {
		panic(err)
	}
}
