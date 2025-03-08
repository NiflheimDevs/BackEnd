package services

import (
	"github.com/jackc/pgx/v5"
	"github.com/niflheimdevs/backend/internal/enums"
	"github.com/niflheimdevs/backend/internal/exceptions"
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
			panic(exceptions.Exception{
				Tag: enums.NOT_FOUND,
				Errors: []enums.SpecificError{
					enums.USERNAME_PASSWORD_WRONG,
				},
			})
		}
	}

	if err != nil {
		panic(exceptions.Exception{
			Tag: enums.NOT_FOUND,
			Errors: []enums.SpecificError{
				enums.USERNAME_PASSWORD_WRONG,
			},
		})
	}

	err = bcrypt.CompareHashAndPassword(user.Password, []byte(password))

	if err != nil {
		panic(exceptions.Exception{
			Tag: enums.NOT_FOUND,
			Errors: []enums.SpecificError{
				enums.USERNAME_PASSWORD_WRONG,
			},
		})
	}

	return user
}

func (userService *UserService) ChangePassword(user_id int, old_password string, new_password string) {
	user, err := userService.UserRepo.FindUserByID(user_id)

	if err != nil {
		panic(exceptions.Exception{
			Tag: enums.INTERNAL_ERROR,
			Errors: []enums.SpecificError{
				enums.SYS_DATABASE_ERROR,
			},
		})
	}

	err = bcrypt.CompareHashAndPassword(user.Password, []byte(old_password))

	if err != nil {
		panic(exceptions.Exception{
			Tag: enums.VALIDATION_ERROR,
			Errors: []enums.SpecificError{
				enums.PASSWORD_INVALID,
			},
		})
	}

	password, err := bcrypt.GenerateFromPassword([]byte(new_password), bcrypt.DefaultCost)

	if err != nil {
		panic(exceptions.Exception{
			Tag: enums.INTERNAL_ERROR,
			Errors: []enums.SpecificError{
				enums.SYS_SERVICE_UNAVAILABLE,
			},
		})
	}

	err = userService.UserRepo.UpdateUserPassword(user_id, password)

	if err != nil {
		panic(exceptions.Exception{
			Tag: enums.INTERNAL_ERROR,
			Errors: []enums.SpecificError{
				enums.SYS_DATABASE_ERROR,
			},
		})
	}
}
