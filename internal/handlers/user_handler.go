package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/niflheimdevs/backend/internal/bootstrap"
	dto "github.com/niflheimdevs/backend/internal/dto/users"
	jwt_keys "github.com/niflheimdevs/backend/internal/jwt"
	"github.com/niflheimdevs/backend/internal/services"
	"github.com/niflheimdevs/backend/internal/services/communications/sms"
)

type UserHandler struct {
	Constants   *bootstrap.Constants
	UserService *services.UserService
	JWTService  *services.JWTToken
	Validator   *validator.Validate
}

func NewUserHandler(
	Constants *bootstrap.Constants,
	userService *services.UserService,
	jwtService *services.JWTToken,
	validator *validator.Validate,
) *UserHandler {
	return &UserHandler{
		Constants:   Constants,
		UserService: userService,
		Validator:   validator,
		JWTService:  jwtService,
	}
}

func (userHandler *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	type loginParams struct {
		Identifier string `json:"identifier" validate:"required"`
		Password   string `json:"password" validate:"required"`
	}

	params := Validated[loginParams](userHandler.Validator, r)

	params.Identifier = strings.ToLower(params.Identifier)

	user := userHandler.UserService.AuthenticateUser(params.Identifier, params.Password)

	jwt_keys.SetupJWTKeys(userHandler.Constants.JWTKeysPath)

	accessToken, refreshToken := userHandler.JWTService.GenerateToken(user.ID)

	userDTO := dto.LoginDTO{
		Username:     user.Username,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(userDTO); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (userHandler *UserHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	type changePasswordParam struct {
		OldPassword string `json:"old_password" validate:"required"`
		NewPassword string `json:"new_password" validate:"required,password"`
	}

	params := Validated[changePasswordParam](userHandler.Validator, r)

	userID := r.Context().Value("userID").(int)

	userHandler.UserService.ChangePasswordValidate(userID, params.OldPassword)
	userHandler.UserService.ChangePassword(userID, params.NewPassword)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (uh *UserHandler) ReserveInfo(w http.ResponseWriter, r *http.Request) {
	type Info struct {
		Phonenumber string `json:"phonenumber" validate:"required,phone"`
		Username    string `json:"username" validate:"required,username"`
		Password    string `json:"password" validate:"required,password"`
	}
	info := Validated[Info](uh.Validator, r)

	info.Username = strings.ToLower(info.Username)
	session := uh.UserService.CheckAvailabilityForSignup(info.Phonenumber, info.Username)

	if session == "" {
		code := sms.GenerateOTP()
		session = uh.UserService.CacheUserInfo(info.Phonenumber, info.Username, info.Password, code)
		// ? placement
		sms.SendOTP(info.Phonenumber, code)
	}

	w.Write([]byte(session))
}

func (uh *UserHandler) SignupWithOtp(w http.ResponseWriter, r *http.Request) {
	type Params struct {
		Code      string `json:"code" validate:"required,len=5,numeric"`
		SessionID string `json:"sessionid" validate:"required"`
	}
	params := Validated[Params](uh.Validator, r)
	phonenumber, username, password := uh.UserService.ValidateOTP(params.SessionID, params.Code)
	userid := uh.UserService.Register(phonenumber, username, password, params.SessionID)

	type Tokens struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}
	var tokens Tokens
	tokens.AccessToken, tokens.RefreshToken = uh.JWTService.GenerateToken(userid)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tokens)

}

func (uh *UserHandler) SendOTP(w http.ResponseWriter, r *http.Request) {
	type Params struct {
		Phonenumber string `json:"phonenumber" validate:"required,phone"`
	}
	params := Validated[Params](uh.Validator, r)
	code := sms.GenerateOTP()
	session := uh.UserService.SetupOTP(params.Phonenumber, code)
	// ? placement
	sms.SendOTP(params.Phonenumber, code)

	w.Write([]byte(session))
}

func (uh *UserHandler) VerifyOTP(w http.ResponseWriter, r *http.Request) {
	type Params struct {
		Code      string `json:"code" validate:"required,len=5,numeric"`
		SessionID string `json:"sessionid" validate:"required"`
	}
	params := Validated[Params](uh.Validator, r)
	userID, _, _ := uh.UserService.ValidateOTP(params.SessionID, params.Code)
	session := uh.UserService.SetForgetPasswordFlag(userID)

	w.Write([]byte(session))
}

func (uh *UserHandler) ForgetPassword(w http.ResponseWriter, r *http.Request) {
	type ChangePassword struct {
		NewPassword string `json:"new_password" validate:"required,password"`
		SessionID   string `json:"sessionid" validate:"required"`
	}

	changePassword := Validated[ChangePassword](uh.Validator, r)

	userID := uh.UserService.CheckFlagForPasswordReset(changePassword.SessionID)

	uh.UserService.ChangePassword(userID, changePassword.NewPassword)
}
