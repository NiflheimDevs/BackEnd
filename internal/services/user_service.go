package services

import (
	"github.com/niflheimdevs/backend/internal/models"
	"github.com/niflheimdevs/backend/internal/repositories"
)

type UserService struct {
	UserRepo repositories.UserRepo
}

func NewUserService(userRepo repositories.UserRepo) *UserService {
	return &UserService{
		UserRepo: userRepo,
	}
}

func (userService *UserService) AuthenticateUser(username string, password string) (user *models.UserModel) {

}
