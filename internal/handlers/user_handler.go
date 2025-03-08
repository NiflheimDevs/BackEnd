package handlers

import (
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/niflheimdevs/backend/internal/services"
	"github.com/niflheimdevs/backend/internal/services/communications/sms"
)

type UserHandler struct {
	UserService services.UserService
	Validator   *validator.Validate
}

func NewUserHandler(userService services.UserService, validator *validator.Validate) *UserHandler {
	return &UserHandler{
		UserService: userService,
		Validator:   validator,
	}
}

func (userHandler *UserHandler) Login(w http.ResponseWriter, r *http.Request) {

}

func (uh *UserHandler) ReserveInfo(w http.ResponseWriter, r *http.Request) {
	type Info struct {
		Phonenumber string `json:"phonenumber" validator:"required,phonenumber"`
		Username    string `json:"username" validator:"required,username"`
		Password    string `json:"password" validator:"required,password"`
		// FirstName   string `json:"firstname" validator:"required"`
		// LastName    string `json:"lastname" validator:"required"`
	}

	info := Validated[Info](uh.Validator, r)

	info.Username = strings.ToLower(info.Username)
	uh.UserService.CheckAvailabilityForSignup(info.Phonenumber, info.Username)

	code := sms.GenerateOTP()
	sms.SendOTP(info.Phonenumber)

	session := uh.UserService.CacheUserInfo(info.Phonenumber, info.Username, info.Password, code)

	w.Write([]byte(session))
}

func (uh *UserHandler) VerifyOTP(w http.ResponseWriter, r *http.Request) {
	type Params struct {
		Code      string `json:"code" validator:"required,len=6,numeric"`
		SessionID string `json:"sessionid" validator:"required"`
	}
	params := Validated[Params](uh.Validator, r)
	phoenenumber, username, password := uh.UserService.ValidateOTP(params.SessionID, params.Code)
	uh.UserService.Register(phoenenumber, username, password)

}

func (uh *UserHandler) RedisTest(w http.ResponseWriter, r *http.Request) {
	uh.UserService.CacheRepo.RedisPing()
}
