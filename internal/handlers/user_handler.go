package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

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
		Identifier string `json:"identifier" validate:"required"`
		Password   string `json:"password" validate:"required"`
	}

	params := Validated[loginParams](r)

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
		NewPassword string `json:"new_password" validate:"required"`
	}

	params := Validated[changePasswordParam](r)

	userID := r.Context().Value("userID").(int)

	userHandler.UserService.ChangePassword(userID, params.OldPassword, params.NewPassword)

	log.Printf("%d has changed his password", userID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
