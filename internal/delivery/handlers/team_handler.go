package handlers

import (
	"github.com/go-playground/validator/v10"
    "github.com/niflheimdevs/backend/bootstrap"
	"github.com/niflheimdevs/backend/internal/application/services"
)

type TeamHandler struct {
	TeamService services.TeamService
	Constants      *bootstrap.Constants
	Validator      *validator.Validate
}

func NewTeamHandler(
	teamService services.TeamService,
	constants *bootstrap.Constants,
	validator *validator.Validate,
) *TeamHandler {
	return &TeamHandler{
		TeamService: teamService,
		Constants:      constants,
		Validator:      validator,
	}
}

