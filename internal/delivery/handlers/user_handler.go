package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/niflheimdevs/backend/bootstrap"
	"github.com/niflheimdevs/backend/internal/application/dto"
	"github.com/niflheimdevs/backend/internal/application/services"
	"github.com/niflheimdevs/backend/internal/domain/exceptions"
	elasticmodel "github.com/niflheimdevs/backend/internal/domain/models/elastic"
	"github.com/niflheimdevs/backend/internal/utils"
)

type UserHandler struct {
	Constants     *bootstrap.Constants
	UserService   services.UserService
	JWTService    services.JWT
	Validator     *validator.Validate
	TagService    services.TagService
	CareerService services.CareerService
	SmsService    services.SmsService
}

func NewUserHandler(
	Constants *bootstrap.Constants,
	userService services.UserService,
	jwtService services.JWT,
	validator *validator.Validate,
	tagService services.TagService,
	careerService services.CareerService,
	smsService services.SmsService,
) *UserHandler {
	return &UserHandler{
		Constants:     Constants,
		UserService:   userService,
		Validator:     validator,
		JWTService:    jwtService,
		CareerService: careerService,
		TagService:    tagService,
		SmsService:    smsService,
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
		panic(exceptions.Exception{
			Tag:    exceptions.INTERNAL_ERROR,
			Errors: []exceptions.SpecificError{exceptions.CAST_ERROR},
		})
	}
}

func (userHandler *UserHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	type changePasswordParam struct {
		OldPassword string `json:"old_password" validate:"required"`
		NewPassword string `json:"new_password" validate:"required,password"`
	}

	params := Validated[changePasswordParam](userHandler.Validator, r)

	userID := r.Context().Value(userHandler.Constants.Context.UserID).(int)

	userHandler.UserService.ChangePasswordValidate(userID, params.OldPassword)
	userHandler.UserService.ChangePassword(userID, params.NewPassword)

	w.WriteHeader(http.StatusNoContent)
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
		code := uh.SmsService.GenerateOTP()
		session = uh.UserService.CacheUserInfo(info.Phonenumber, info.Username, info.Password, code)
		// ? placement
		uh.SmsService.SendOTP(info.Phonenumber, code)
	}

	w.WriteHeader(http.StatusOK)
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
	code := uh.SmsService.GenerateOTP()
	session := uh.UserService.SetupOTP(params.Phonenumber, code)
	// ? placement
	uh.SmsService.SendOTP(params.Phonenumber, code)

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

	w.WriteHeader(http.StatusNoContent)
}

func (uh *UserHandler) UpdateUserData(w http.ResponseWriter, r *http.Request) {
	userid, _ := r.Context().Value(uh.Constants.Context.UserID).(int)
	params := Validated[dto.UpdateUserDTO](uh.Validator, r)

	uh.UserService.UpdateUserData(userid, &params)

	w.WriteHeader(http.StatusNoContent)
}

func (uh *UserHandler) UpdateUsername(w http.ResponseWriter, r *http.Request) {
	type Param struct {
		Username string `json:"username" validate:"required,username"`
	}

	param := Validated[Param](uh.Validator, r)

	userid, _ := r.Context().Value(uh.Constants.Context.UserID).(int)

	param.Username = strings.ToLower(param.Username)

	uh.UserService.UpdateUsername(userid, param.Username)

	w.WriteHeader(http.StatusNoContent)
}

func (uh *UserHandler) UpdateEmail(w http.ResponseWriter, r *http.Request) {

	type Param struct {
		Email string `json:"email" validate:"required,email"`
	}

	param := Validated[Param](uh.Validator, r)

	userid, _ := r.Context().Value(uh.Constants.Context.UserID).(int)

	param.Email = strings.ToLower(param.Email)

	uh.UserService.UpdateEmail(userid, param.Email)

	w.WriteHeader(http.StatusNoContent)
}

func (uh *UserHandler) UpdatePhoneSendOTP(w http.ResponseWriter, r *http.Request) {
	type NewPhone struct {
		NewPhone string `json:"phone" validate:"required,phone"`
	}

	userid, _ := r.Context().Value(uh.Constants.Context.UserID).(int)
	params := Validated[NewPhone](uh.Validator, r)

	code := uh.SmsService.GenerateOTP()
	session := uh.UserService.UpdatePhoneSendOTP(params.NewPhone, userid, code)
	uh.SmsService.SendOTP(params.NewPhone, code)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(session))
}

func (uh *UserHandler) UpdatePhoneVerify(w http.ResponseWriter, r *http.Request) {
	type Params struct {
		Code      string `json:"code" validate:"required,len=5,numeric"`
		SessionID string `json:"sessionid" validate:"required"`
	}

	params := Validated[Params](uh.Validator, r)

	phone, userid, _ := uh.UserService.ValidateOTP(params.SessionID, params.Code)

	uh.UserService.UpdatePhone(phone, userid)

	w.WriteHeader(http.StatusNoContent)
}

func (uh *UserHandler) GetUserInfo(w http.ResponseWriter, r *http.Request) {
	userid := r.Context().Value(uh.Constants.Context.UserID).(int)
	useridString := chi.URLParam(r, "id")
	targetUserid, err := strconv.Atoi(useridString)
	if err != nil {
		panic(exceptions.Exception{
			Tag: exceptions.BAD_REQUEST,
		})
	}
	if targetUserid <= 0 {
		targetUserid = userid
	}
	includes := r.URL.Query()["include"]
	response := make(map[string]interface{})

	if utils.Contains(includes, "info") {
		response["info"] = uh.UserService.GetUserInfo(targetUserid, userid)
	}
	if utils.Contains(includes, "career") {
		response["career"] = uh.CareerService.GetCareersForUser(userid, targetUserid)
	}
	if utils.Contains(includes, "tag") {
		response["tag"] = uh.TagService.GetTagsForUser(targetUserid)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (uh *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userid := r.Context().Value(uh.Constants.Context.UserID).(int)

	uh.UserService.DeleteUser(userid)

	w.WriteHeader(http.StatusNoContent)
}

func (uh *UserHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	type RefreshTokenParam struct {
		RefreshToken string `json:"refresh_token" validate:"required"`
	}

	params := Validated[RefreshTokenParam](uh.Validator, r)

	newaccesstoken := uh.JWTService.RefreshToken(params.RefreshToken)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(newaccesstoken))
}

func (uh *UserHandler) SearchUsers(w http.ResponseWriter, r *http.Request) {

	params := Validated[elasticmodel.QueryAndTagSearchReqDto](uh.Validator, r)

	if params.Order != "" && params.Order != "desc" && params.Order != "asc" {
		panic(exceptions.Exception{
			Tag: exceptions.BAD_REQUEST,
		})
	}

	res := uh.UserService.SearchUsers(&params)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}
