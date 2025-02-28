package services

import "github.com/niflheimdevs/backend/internal/repositories"

type UserService struct {
	userRepo repositories.UserRepo
}

func NewUserService(userRepo repositories.UserRepo) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}
