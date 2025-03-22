package services

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/google/uuid"

	dto "github.com/niflheimdevs/backend/internal/dto/users"
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

// signup stage. checks if username or phonenumber is already taken.
func (us *UserService) CheckAvailabilityForSignup(phonenumber string, username string) string {
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

	sess1, err := us.CacheRepo.FindByUsername(username)
	if err == nil {
		exc.AddError(enums.USERNAME_TAKEN)
	}

	sess2, err := us.CacheRepo.FindByPhone(phonenumber)
	if err == nil {
		exc.AddError(enums.PHONE_TAKEN)

	}
	if sess1 == sess2 {
		return sess1
	}

	if len(exc.Errors) > 0 {
		panic(exc)
	}
	return ""
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

	userdata := models.UserCacheData{
		Phone:    phonenumber,
		Username: username,
		Password: hashedPass,
		OTP:      otp,
	}

	us.CacheRepo.PostUserCreds(session, &userdata)

	return session
}

var develop_mode = true

// checks the code with cache and returns the data. returns phonenumber, username, password
func (us *UserService) ValidateOTP(session string, otp string) (string, string, []byte) {
	val, err := us.CacheRepo.FindBySession(session)
	if err != nil {
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

	us.CacheRepo.DeleteRedis(session)

	return userdata.Phone, userdata.Username, userdata.Password
}

// last stage of signup
func (us *UserService) Register(phonenumber string, username string, password []byte, session string) int {
	//? before creation in main?
	us.CacheRepo.ClearUserCreds(phonenumber, username)

	userid, err := us.UserRepo.PostUser(phonenumber, username, password)
	if err != nil {
		panic(exceptions.Exception{
			Tag:    enums.INTERNAL_ERROR,
			Errors: []enums.SpecificError{enums.DATABASE_ERROR},
		})
	}
	return userid
}

func (userService *UserService) AuthenticateUser(identifier string, password string) *models.UserModel {
	user, err := userService.UserRepo.FindUserByUsername(identifier)
	if err != nil {
		user, err = userService.UserRepo.FindUserByPhone(identifier)
		if err != nil {
			panic(exceptions.Exception{
				Tag: enums.NOT_FOUND,
				Errors: []enums.SpecificError{
					enums.USERNAME_PASSWORD_WRONG,
				},
			})
		}
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

func (userService *UserService) ChangePasswordValidate(user_id int, old_password string) {
	if user_id == -1 {
		panic(exceptions.Exception{
			Tag: enums.AUTHENTICATION_ERROR,
			Errors: []enums.SpecificError{
				enums.AUTH_INVALID_CREDENTIALS,
			},
		})
	}

	user, err := userService.UserRepo.FindUserByID(user_id)

	if err != nil {
		panic(exceptions.Exception{
			Tag: enums.INTERNAL_ERROR,
			Errors: []enums.SpecificError{
				enums.DATABASE_ERROR,
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
}
func (userService *UserService) ChangePassword(user_id int, new_password string) {
	password, err := bcrypt.GenerateFromPassword([]byte(new_password), bcrypt.DefaultCost)

	if err != nil {
		panic(exceptions.Exception{
			Tag: enums.INTERNAL_ERROR,
			Errors: []enums.SpecificError{
				enums.SERVICE_UNAVAILABLE,
			},
		})
	}

	err = userService.UserRepo.UpdateUserPassword(user_id, password)

	if err != nil {
		panic(exceptions.Exception{
			Tag: enums.INTERNAL_ERROR,
			Errors: []enums.SpecificError{
				enums.DATABASE_ERROR,
			},
		})
	}
}

// for forget password: creates uuid and otp and caches otp
func (us *UserService) SetupOTP(phonenumber string, code string) string {
	userdata, err := us.UserRepo.FindUserByPhone(phonenumber)
	if err != nil {
		panic(exceptions.Exception{
			Tag:    enums.UNPROCESSABLE,
			Errors: []enums.SpecificError{enums.USER_NOT_FOUND},
		})
	}

	session := uuid.New().String()
	val := models.UserCacheData{
		Phone: fmt.Sprintf("%d", userdata.ID),
		OTP:   code,
	}

	us.CacheRepo.PostSessionOTP(session, &val)

	return session
}

func (us *UserService) SetForgetPasswordFlag(userID string) string {
	session := uuid.New().String()
	us.CacheRepo.PostSessionFlag(session, userID)
	return session
}

func (us *UserService) CheckFlagForPasswordReset(session string) int {
	val, err := us.CacheRepo.GetSessionFlag(session)
	if err != nil {
		panic(exceptions.Exception{
			Tag:    enums.BAD_REQUEST,
			Errors: []enums.SpecificError{enums.BAD_SESSION},
		})
	}
	userID, _ := strconv.Atoi(val)
	return userID
}

func (us *UserService) UpdateUserData(userData *dto.UpdateUserDTO) {
	if userData.ID < 0 {
		panic(exceptions.Exception{
			Tag: enums.UNAUTHORIZED,
			Errors: []enums.SpecificError{
				enums.AUTH_ACCESS_DENIED,
			},
		})
	}

	user, err := us.UserRepo.FindUserByID(userData.ID)

	if err != nil {
		panic(exceptions.Exception{
			Tag:    enums.UNAUTHORIZED,
			Errors: []enums.SpecificError{enums.USER_NOT_FOUND},
		})
	}
	if !user.Is_verified && user.Email != userData.Email {
		panic(exceptions.Exception{
			Tag:    enums.VALIDATION_ERROR,
			Errors: []enums.SpecificError{enums.EMAIL_NOT_VERIFIED},
		})
	}

	_, err = us.UserRepo.UpdateUserData(userData)

	if err == nil {
		return
	}

	validErrs := exceptions.Exception{
		Tag: enums.VALIDATION_ERROR,
	}
	_, err = us.UserRepo.FindUserByEmail(userData.Email)
	if err != nil {
		validErrs.AddError(enums.EMAIL_TAKEN)
	}
	_, err = us.UserRepo.FindUserByUsername(userData.Username)
	if err != nil {
		validErrs.AddError(enums.USERNAME_TAKEN)
	}
	if len(validErrs.Errors) > 0 {
		panic(validErrs)
	}

	panic(exceptions.Exception{
		Tag: enums.INTERNAL_ERROR,
	})
}

func (us *UserService) UpdatePhoneSendOTP(phone string, userid int, code string) string {
	if userid < 0 {
		panic(exceptions.Exception{
			Tag:    enums.UNAUTHORIZED,
			Errors: []enums.SpecificError{enums.AUTH_ACCESS_DENIED},
		})
	}
	_, err := us.CacheRepo.FindByPhone(phone)
	if err == nil {
		panic(exceptions.Exception{
			Tag:    enums.VALIDATION_ERROR,
			Errors: []enums.SpecificError{enums.PHONE_TAKEN},
		})
	}

	_, err = us.UserRepo.FindUserByPhone(phone)
	if err == nil {
		panic(exceptions.Exception{
			Tag:    enums.VALIDATION_ERROR,
			Errors: []enums.SpecificError{enums.PHONE_TAKEN},
		})
	}

	session := uuid.New().String()
	val := models.UserCacheData{
		Phone:    phone,
		Username: fmt.Sprintf("%d", userid),
		OTP:      code,
	}

	us.CacheRepo.PostPhone(phone, session)
	us.CacheRepo.PostSessionOTP(session, &val)

	return session
}

func (us *UserService) UpdatePhone(phone string, userid string) {
	us.UserRepo.UpdatePhone(phone, userid)
}
