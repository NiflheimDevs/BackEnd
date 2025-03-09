package services

import (
	"encoding/json"
	"github.com/google/uuid"
  
	"github.com/jackc/pgx/v5"
  
	"github.com/niflheimdevs/backend/internal/enums"
	"github.com/niflheimdevs/backend/internal/exceptions"
	"github.com/niflheimdevs/backend/internal/models"
	"github.com/niflheimdevs/backend/internal/repositories"
	"github.com/niflheimdevs/backend/internal/repositories/redis"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	UserRepo  *repositories.UserRepo
	CacheRepo *redis.UserCache
}

func NewUserService(userRepo *repositories.UserRepo, cacheRepo *redis.UserCache) *UserService {
	return &UserService{
		UserRepo:  userRepo,
		CacheRepo: cacheRepo,
	}
}

func (userService *UserService) AuthenticateUser(username string, password string) (user *models.UserModel) {
	return nil
}

// signup stage. checks if username or phonenumber is already taken.
func (us *UserService) CheckAvailabilityForSignup(phonenumber string, username string) {
	var exc = exceptions.Exception{
		Tag: enums.VALIDATION_ERROR,
	}
	_, err := us.UserRepo.FindUserByPhone(phonenumber)
	if err == nil {
		exc.AddError(enums.PHONE_TAKEN)
	}
	_, err = us.UserRepo.FindUserByUsername(username)
	if err == nil {
		exc.AddError(enums.USERNAME_TAKEN)
	}

	if len(exc.Errors) > 0 {
		panic(exc)
	}

	//TODO: handle sending val somehow
	_, err = us.CacheRepo.FindByUsername(username)
	if err == nil {
		// panic(val)
		exc.AddError(enums.USERNAME_TAKEN)
	}

	_, err = us.CacheRepo.FindByPhone(phonenumber)
	if err == nil {
		exc.AddError(enums.PHONE_TAKEN)
	}

	if len(exc.Errors) > 0 {
		panic(exc)
	}
}

// caches the data until otp is expired or entered.
func (us *UserService) CacheUserInfo(phonenumber string, username string, password string, otp string) string {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(exceptions.Exception{
			Tag: enums.INTERNAL_ERROR,
		})
	}

	session := uuid.New().String()

	err = us.CacheRepo.PostUserCreds(phonenumber, username, hashedPass, session, otp)
	if err != nil {
		panic(exceptions.Exception{
			Tag:    enums.INTERNAL_ERROR,
			Errors: []enums.SpecificError{enums.CACHE_ERROR},
		})
	}

	return session
}

var develop_mode = true

// checks the code with cache and returns the data. returns phonenumber, username, password
func (us *UserService) ValidateOTP(session string, otp string) (string, string, []byte) {
	val := us.CacheRepo.GetUserCreds(session)
	if val == "" {
		panic(exceptions.Exception{
			Tag:    enums.VALIDATION_ERROR,
			Errors: []enums.SpecificError{enums.OTP_EXPIRED_OR_BAD_SESSION},
		})
	}

	var userdata models.UserCacheData
	json.Unmarshal([]byte(val), &userdata)

	if otp != userdata.OTP && !develop_mode {
		panic(exceptions.Exception{
			Tag:    enums.VALIDATION_ERROR,
			Errors: []enums.SpecificError{enums.OTP_INVALID},
		})
	}
	return userdata.Phone, userdata.Username, userdata.Password
}

// last stage of signup
func (us *UserService) Register(phonenumber string, username string, password []byte) {
	us.CacheRepo.ClearUserCreds(phonenumber, username)

	err := us.UserRepo.PostUser(phonenumber, username, password)
	if err != nil {
		panic(exceptions.Exception{
			Tag:    enums.INTERNAL_ERROR,
			Errors: []enums.SpecificError{enums.DATABASE_ERROR},
      })
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
