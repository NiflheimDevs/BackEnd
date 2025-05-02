package handlers

import (
	"github.com/go-playground/validator/v10"
    "github.com/niflheimdevs/backend/bootstrap"
	"github.com/niflheimdevs/backend/internal/application/services"
)

type RoleHandler struct {
	RoleService services.RoleService
	Constants      *bootstrap.Constants
	Validator      *validator.Validate
}

func NewRoleHandler(
	roleService services.RoleService,
	constants *bootstrap.Constants,
	validator *validator.Validate,
) *RoleHandler {
	return &RoleHandler{
		RoleService: roleService,
		Constants:      constants,
		Validator:      validator,
	}
}

