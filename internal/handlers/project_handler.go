package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/niflheimdevs/backend/internal/bootstrap"
	dto "github.com/niflheimdevs/backend/internal/dto/projects"
	"github.com/niflheimdevs/backend/internal/enums"
	"github.com/niflheimdevs/backend/internal/exceptions"
	"github.com/niflheimdevs/backend/internal/services"
)

type ProjectHandler struct {
	Constants      *bootstrap.Constants
	ProjectService *services.ProjectService
	JWTService     *services.JWT
	Validator      *validator.Validate
}

func NewProjectHandler(
	Constants *bootstrap.Constants,
	projectService *services.ProjectService,
	jwtService *services.JWT,
	validator *validator.Validate,
) *ProjectHandler {
	return &ProjectHandler{
		Constants:      Constants,
		ProjectService: projectService,
		JWTService:     jwtService,
		Validator:      validator,
	}
}

func (projectHandler *ProjectHandler) CreateProject(w http.ResponseWriter, r *http.Request) {
	type createProjectParams struct {
		Title       string `json:"title" validate:"required"`
		Description string `json:"description" validate:"required"`
		Tag         string `json:"tag" validate:"required"`
	}

	UserID := r.Context().Value("UserID").(int)

	params := Validated[createProjectParams](projectHandler.Validator, r)

	var tags []int

	json.Unmarshal([]byte(params.Tag), &tags)

	project := projectHandler.ProjectService.CreateProject(UserID, params.Title, params.Description, tags)

	dto := dto.CreateProjectDTO{
		ProjectID: project,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(dto); err != nil {
		panic(exceptions.Exception{
			Tag:    enums.INTERNAL_ERROR,
			Errors: []enums.SpecificError{enums.CAST_ERROR},
		})
	}
}

func (projectHandler *ProjectHandler) UpdateProject(w http.ResponseWriter, r *http.Request) {
	type updateProjectParams struct {
		Title       string `json:"title" validate:"required"`
		Description string `json:"description" validate:"required"`
		Tag         string `json:"tag" validate:"required"`
	}

	UserID := r.Context().Value("UserID").(int)

	ProjectID := r.Context().Value("ProjectID").(int)

	params := Validated[updateProjectParams](projectHandler.Validator, r)

	var tags []int

	json.Unmarshal([]byte(params.Tag), &tags)

	projectHandler.ProjectService.UpdateProject(ProjectID, UserID, params.Title, params.Description, tags)

	w.WriteHeader(http.StatusNoContent)
}
