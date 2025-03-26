package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/niflheimdevs/backend/internal/bootstrap"
	"github.com/niflheimdevs/backend/internal/dto"
	"github.com/niflheimdevs/backend/internal/enums"
	"github.com/niflheimdevs/backend/internal/exceptions"
	"github.com/niflheimdevs/backend/internal/services"
)

type ProjectHandler struct {
	Constants      *bootstrap.Constants
	ProjectService *services.ProjectService
	UserService    *services.UserService
	JWTService     *services.JWT
	Validator      *validator.Validate
}

func NewProjectHandler(
	constants *bootstrap.Constants,
	projectService *services.ProjectService,
	userService *services.UserService,
	jwtService *services.JWT,
	validator *validator.Validate,
) *ProjectHandler {
	return &ProjectHandler{
		Constants:      constants,
		ProjectService: projectService,
		UserService:    userService,
		JWTService:     jwtService,
		Validator:      validator,
	}
}

func (projectHandler *ProjectHandler) GetUserProject(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(projectHandler.Constants.Context.UserID).(int)

	query := r.URL.Query()
	offset, _ := strconv.Atoi(query.Get("offset"))
	limit, _ := strconv.Atoi(query.Get("limit"))

	projects := projectHandler.ProjectService.GetUserProjects(userID, offset, limit)

	var projectsDTO []dto.Project

	for _, project := range projects {
		projectsDTO = append(projectsDTO, dto.Project{
			ProjectID:   project.ID,
			OwnerID:     project.OwnerID,
			Title:       project.Title,
			Description: project.Description,
			Label:       project.Label,
			Tags:        project.Tags,
			Duration:    project.Duration,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(projectsDTO); err != nil {
		panic(exceptions.Exception{
			Tag:    enums.INTERNAL_ERROR,
			Errors: []enums.SpecificError{enums.CAST_ERROR},
		})
	}
}

func (projectHandler *ProjectHandler) GetSpeceficProject(w http.ResponseWriter, r *http.Request) {
	projectIDString := chi.URLParam(r, "project_id")
	projectID, _ := strconv.Atoi(projectIDString)

	project := projectHandler.ProjectService.GetProject(projectID)
	userInfo := projectHandler.UserService.GetUserInfo(project.OwnerID, 0)

	projectDTO := dto.Project{
		ProjectID:   project.ID,
		OwnerID:     project.OwnerID,
		Title:       project.Title,
		Description: project.Description,
		Label:       project.Label,
		FirstName:   userInfo.FirstName,
		LastName:    userInfo.LastName,
		Username:    userInfo.Username,
		Tags:        project.Tags,
		Duration:    project.Duration,
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
		Label       int    `json:"label" validate:"required"`
	}

	params := Validated[createProjectParams](projectHandler.Validator, r)

	userID := r.Context().Value(projectHandler.Constants.Context.UserID).(int)

	project := projectHandler.ProjectService.CreateProject(userID, params.Title, params.Description, params.Label, params.Tags)

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
		Title       string `json:"title" validate:"required,max=50"`
		Description string `json:"description" validate:"required"`
		Tags        []int  `json:"tags" validate:"required"`
		Label       int    `json:"label" validate:"required"`
	}

	params := Validated[updateProjectParams](projectHandler.Validator, r)

	userID := r.Context().Value(projectHandler.Constants.Context.UserID).(int)

	projectIDString := chi.URLParam(r, "project_id")
	projectID, _ := strconv.Atoi(projectIDString)

	projectHandler.ProjectService.UpdateProject(projectID, userID, params.Title, params.Description, params.Label, params.Tags)

	w.WriteHeader(http.StatusNoContent)
}

func (projectHandler *ProjectHandler) DeleteProject(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(projectHandler.Constants.Context.UserID).(int)

	projectIDString := chi.URLParam(r, "project_id")
	projectID, _ := strconv.Atoi(projectIDString)

	projectHandler.ProjectService.DeleteProject(userID, projectID)

	w.WriteHeader(http.StatusNoContent)
}
