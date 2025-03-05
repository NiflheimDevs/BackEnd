package handlers

import (
	"net/http"

	"github.com/niflheimdevs/backend/internal/services"
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
