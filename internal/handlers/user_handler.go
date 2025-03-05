package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/niflheimdevs/backend/internal/bootstrap"
	dto "github.com/niflheimdevs/backend/internal/dto/users"
	jwt_keys "github.com/niflheimdevs/backend/internal/jwt"
	"github.com/niflheimdevs/backend/internal/services"
)

type UserHandler struct {
	Constants   *bootstrap.Constants
	UserService services.UserService
	JWTService  services.JWTToken
}

func NewUserHandler(
	Constants *bootstrap.Constants,
	userService services.UserService,
	jwtService services.JWTToken,
) *UserHandler {
	return &UserHandler{
		Constants:   Constants,
		UserService: userService,
		JWTService:  jwtService,
	}
}

func (userHandler *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	type loginParams struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	params := Validated[loginParams](r)

	user := userHandler.UserService.AuthenticateUser(params.Username, params.Password)

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
