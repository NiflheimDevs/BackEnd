package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/niflheimdevs/backend/bootstrap"
	"github.com/niflheimdevs/backend/internal/application/services"
	"github.com/niflheimdevs/backend/internal/domain/exceptions"
)

type RoleHandler struct {
	RoleService services.RoleService
	Constants   *bootstrap.Constants
	Validator   *validator.Validate
}

func NewRoleHandler(
	roleService services.RoleService,
	constants *bootstrap.Constants,
	validator *validator.Validate,
) *RoleHandler {
	return &RoleHandler{
		RoleService: roleService,
		Constants:   constants,
		Validator:   validator,
	}
}

func (rh *RoleHandler) GetTeamRoles(w http.ResponseWriter, r *http.Request) {
	res := rh.RoleService.GetAllRolesForTeam()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		panic(exceptions.Exception{
			Tag:    exceptions.INTERNAL_ERROR,
			Errors: []exceptions.SpecificError{exceptions.CAST_ERROR},
		})
	}
}

func (rh *RoleHandler) GetPermissionsForRole(w http.ResponseWriter, r *http.Request) {
	roleString := chi.URLParam(r, "role")

	res := rh.RoleService.GetPermissionsForRole(roleString)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		panic(exceptions.Exception{
			Tag:    exceptions.INTERNAL_ERROR,
			Errors: []exceptions.SpecificError{exceptions.CAST_ERROR},
		})
	}
}
