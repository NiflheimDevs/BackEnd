package handlers

import (
	"net/http"

	"github.com/niflheimdevs/backend/internal/services"
	"github.com/niflheimdevs/backend/internal/services/communications/sms"
)

type UserHandler struct {
	UserService services.UserService
}

func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{
		UserService: userService,
	}
}

func (userHandler *UserHandler) Login(w http.ResponseWriter, r *http.Request) {

}

func (uh *UserHandler) ReserveInfo(w http.ResponseWriter, r *http.Request) {
	type Info struct {
		Phonenumber string `json:"phonenumber" validator:"required"`
		Username    string `json:"username" validator:"required,gt=2,lt=20"`
		Password    string `json:"password" validator:"required"`
		// FirstName   string `json:"firstname" validator:"required"`
		// LastName    string `json:"lastname" validator:"required"`
	}

	info := Validated[Info](r)
	uh.UserService.CheckAvailabilityForSignup(info.Phonenumber, info.Username)

	code := sms.GenerateOTP()
	sms.SendOTP(info.Phonenumber)

	session := uh.UserService.CacheUserInfo(info.Phonenumber, info.Username, info.Password, code)

	w.Write([]byte(session))
}

func (uh *UserHandler) VerifyOTP(w http.ResponseWriter, r *http.Request) {
	type Params struct {
		//TODO: validator
		Code      string `json:"code" validator:"required"`
		SessionID string `json:"sessionid" validator:"required"`
	}
	params := Validated[Params](r)
	phoenenumber, username, password := uh.UserService.ValidateOTP(params.SessionID, params.Code)
	uh.UserService.Register(phoenenumber, username, password)

}

func (uh *UserHandler) RedisTest(w http.ResponseWriter, r *http.Request) {
	uh.UserService.CacheRepo.RedisPing()
}
