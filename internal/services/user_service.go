package services

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/niflheimdevs/backend/internal/models"
	"github.com/niflheimdevs/backend/internal/repositories"
	"github.com/niflheimdevs/backend/internal/repositories/redis"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	UserRepo  repositories.UserRepo
	CacheRepo redis.UserCache
}

func NewUserService(userRepo repositories.UserRepo, cacheRepo redis.UserCache) *UserService {
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
	val, err := us.CacheRepo.FindByPhone(phonenumber)
	//TODO: error
	if err == nil {
		panic(val)
	}
	val, err = us.CacheRepo.FindByUsername(username)
	if err == nil {
		panic(val)
	}

	_, err = us.UserRepo.FindUserByUsername(username)
	if err == nil {
		panic(nil)
	}
	_, err = us.UserRepo.FindUserByPhone(phonenumber)
	if err == nil {
		panic(nil)
	}
}

// caches the data until otp is expired or entered.
func (us *UserService) CacheUserInfo(phonenumber string, username string, password string, otp string) string {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	session := uuid.New().String()

	err = us.CacheRepo.PostUserCreds(phonenumber, username, hashedPass, session, otp)
	if err != nil {
		panic(err)
	}

	return session
}

var develop_mode = true

// checks the code with cache and returns the data. returns phonenumber, username, password
func (us *UserService) ValidateOTP(session string, otp string) (string, string, []byte) {
	val := us.CacheRepo.GetUserCreds(session)
	var userdata models.UserCacheData
	json.Unmarshal([]byte(val), &userdata)

	if otp != userdata.OTP && !develop_mode {
		panic(nil)
	}
	return userdata.Phone, userdata.Username, userdata.Password
}

// last stage of signup
func (us *UserService) Register(phonenumber string, username string, password []byte) {
	err := us.UserRepo.PostUser(phonenumber, username, password)
	if err != nil {
		panic(err)
	}

	us.CacheRepo.ClearUserCreds(phonenumber, username)
}
