package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
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

func (projectHandler *ProjectHandler) GetProject(w http.ResponseWriter, r *http.Request) {
	ProjectIDString := chi.URLParam(r, "project_id")
	ProjectID, _ := strconv.Atoi(ProjectIDString)
	project := projectHandler.ProjectService.GetProject(ProjectID)

	projectDTO := dto.Project{
		ProjectID:   project.ID,
		OwnerID:     project.OwnerID,
		Title:       project.Title,
		Description: project.Description,
		Tags:        project.Tags,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(projectDTO); err != nil {
		panic(exceptions.Exception{
			Tag:    enums.INTERNAL_ERROR,
			Errors: []enums.SpecificError{enums.CAST_ERROR},
		})
	}
}

func (projectHandler *ProjectHandler) CreateProject(w http.ResponseWriter, r *http.Request) {
	type createProjectParams struct {
		Title       string `json:"title" validate:"required"`
		Description string `json:"description" validate:"required"`
		Tags        []int  `json:"tags" validate:"required"`
	}

	params := Validated[createProjectParams](projectHandler.Validator, r)

	UserID := r.Context().Value("userID").(int)

	project := projectHandler.ProjectService.CreateProject(UserID, params.Title, params.Description, params.Tags)

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

	params := Validated[updateProjectParams](projectHandler.Validator, r)

	UserID := r.Context().Value("UserID").(int)

	ProjectID := r.Context().Value("ProjectID").(int)

	var tags []int

	json.Unmarshal([]byte(params.Tag), &tags)

	projectHandler.ProjectService.UpdateProject(ProjectID, UserID, params.Title, params.Description, tags)

	w.WriteHeader(http.StatusNoContent)
}
