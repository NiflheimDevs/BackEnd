package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/niflheimdevs/backend/bootstrap"
	"github.com/niflheimdevs/backend/internal/application/dto"
	"github.com/niflheimdevs/backend/internal/application/services"
	"github.com/niflheimdevs/backend/internal/domain/exceptions"
	elasticmodel "github.com/niflheimdevs/backend/internal/domain/models/elastic"
	"github.com/niflheimdevs/backend/internal/utils"
)

type ProjectHandler struct {
	Constants      *bootstrap.Constants
	ProjectService services.ProjectService
	UserService    services.UserService
	LabelService   services.LabelService
	BidService     services.BidService
	JWTService     services.JWT
	Validator      *validator.Validate
}

func NewProjectHandler(
	constants *bootstrap.Constants,
	projectService services.ProjectService,
	userService services.UserService,
	labelService services.LabelService,
	bidService services.BidService,
	jwtService services.JWT,
	validator *validator.Validate,
) *ProjectHandler {
	return &ProjectHandler{
		Constants:      constants,
		ProjectService: projectService,
		UserService:    userService,
		LabelService:   labelService,
		JWTService:     jwtService,
		BidService:     bidService,
		Validator:      validator,
	}
}

func (projectHandler *ProjectHandler) LandingProps(w http.ResponseWriter, r *http.Request) {
	projects := projectHandler.ProjectService.LandingProps()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(projects); err != nil {
		panic(exceptions.Exception{
			Tag:    exceptions.INTERNAL_ERROR,
			Errors: []exceptions.SpecificError{exceptions.CAST_ERROR},
		})
	}
}

func (projectHandler *ProjectHandler) GetUserProject(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(projectHandler.Constants.Context.UserID).(int)
	targetUserIDString := chi.URLParam(r, "user_id")
	targetUserID, _ := strconv.Atoi(targetUserIDString)

	query := r.URL.Query()
	offset, err := strconv.Atoi(query.Get("offset"))
	if err != nil {
		offset = projectHandler.Constants.Pagination.Offset
	}
	limit, err := strconv.Atoi(query.Get("limit"))
	if err != nil {
		limit = projectHandler.Constants.Pagination.Limit
	}

	projects, count := projectHandler.ProjectService.GetUserProjects(userID, targetUserID, offset, limit)

	var projectsDTO dto.UserProject

	for _, project := range projects {
		label := projectHandler.LabelService.GetLabelInfo(project.Label)
		projectsDTO.Projects = append(projectsDTO.Projects, dto.Project{
			ProjectID:   project.ID,
			OwnerID:     project.OwnerID,
			Title:       project.Title,
			Description: project.Description,
			Label:       *label,
			SelectedBid: project.SelectedBid,
			Status:      project.State,
			Tags:        project.Tags,
			Comment:     project.Comment,
			Duration:    project.Duration,
			StartTime:   project.StartTime,
			EndTime:     project.EndTime,
		})
	}

	projectsDTO.Count = count

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(projectsDTO); err != nil {
		panic(exceptions.Exception{
			Tag:    exceptions.INTERNAL_ERROR,
			Errors: []exceptions.SpecificError{exceptions.CAST_ERROR},
		})
	}
}

func (projectHandler *ProjectHandler) GetSpeceficProject(w http.ResponseWriter, r *http.Request) {
	projectIDString := chi.URLParam(r, "project_id")
	projectID, _ := strconv.Atoi(projectIDString)

	project := projectHandler.ProjectService.GetProject(projectID)
	userInfo := projectHandler.UserService.GetUserInfo(project.OwnerID, 0)
	label := projectHandler.LabelService.GetLabelInfo(project.Label)
	bids := projectHandler.BidService.GetPublicProjectBids(projectID)

	projectDTO := dto.Project{
		ProjectID:   project.ID,
		OwnerID:     project.OwnerID,
		Title:       project.Title,
		Description: project.Description,
		Label:       *label,
		Bids:        bids,
		SelectedBid: project.SelectedBid,
		Status:      project.State,
		FirstName:   userInfo.FirstName,
		LastName:    userInfo.LastName,
		Username:    userInfo.Username,
		Tags:        project.Tags,
		Comment:     project.Comment,
		Duration:    project.Duration,
		StartTime:   project.StartTime,
		EndTime:     project.EndTime,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(projectDTO); err != nil {
		panic(exceptions.Exception{
			Tag:    exceptions.INTERNAL_ERROR,
			Errors: []exceptions.SpecificError{exceptions.CAST_ERROR},
		})
	}
}

func (projectHandler *ProjectHandler) CreateProject(w http.ResponseWriter, r *http.Request) {
	type createProjectParams struct {
		Title       string `json:"title" validate:"required"`
		Description string `json:"description" validate:"required"`
		Tags        []int  `json:"tags" validate:"required"`
		Label       *int   `json:"label" validate:"required,gt=-1,lt=3"`
		Duration    int    `json:"duration" validate:"required"`
	}

	params := Validated[createProjectParams](projectHandler.Validator, r)

	userID := r.Context().Value(projectHandler.Constants.Context.UserID).(int)

	price := projectHandler.LabelService.GetLabelInfo(*params.Label).Price

	project := projectHandler.ProjectService.CreateProject(userID, params.Title, params.Description, *params.Label, params.Duration, price, params.Tags)

	dto := dto.CreateProjectDTO{
		ProjectID: project,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(dto); err != nil {
		panic(exceptions.Exception{
			Tag:    exceptions.INTERNAL_ERROR,
			Errors: []exceptions.SpecificError{exceptions.CAST_ERROR},
		})
	}
}

func (projectHandler *ProjectHandler) UpdateProject(w http.ResponseWriter, r *http.Request) {
	type updateProjectParams struct {
		Title       string `json:"title" validate:"required,max=50"`
		Description string `json:"description" validate:"required"`
		Tags        []int  `json:"tags" validate:"required"`
		Label       *int   `json:"label" validate:"required,gt=-1,lt=3"`
	}

	params := Validated[updateProjectParams](projectHandler.Validator, r)

	userID := r.Context().Value(projectHandler.Constants.Context.UserID).(int)

	projectIDString := chi.URLParam(r, "project_id")
	projectID, _ := strconv.Atoi(projectIDString)

	price := projectHandler.LabelService.GetLabelInfo(*params.Label).Price

	projectHandler.ProjectService.UpdateProject(projectID, userID, params.Title, params.Description, *params.Label, price, params.Tags)

	w.WriteHeader(http.StatusNoContent)
}

func (projectHandler *ProjectHandler) DeleteProject(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(projectHandler.Constants.Context.UserID).(int)

	projectIDString := chi.URLParam(r, "project_id")
	projectID, _ := strconv.Atoi(projectIDString)

	projectHandler.ProjectService.DeleteProject(userID, projectID)

	w.WriteHeader(http.StatusNoContent)
}

func (projectHandler *ProjectHandler) EndProject(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(projectHandler.Constants.Context.UserID).(int)

	projectIDString := chi.URLParam(r, "project_id")
	projectID, _ := strconv.Atoi(projectIDString)

	projectHandler.ProjectService.EndOfProject(userID, projectID)

	w.WriteHeader(http.StatusNoContent)
}

func (projectHandler *ProjectHandler) GetTeamProjects(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(projectHandler.Constants.Context.UserID).(int)

	teamIDString := chi.URLParam(r, "team_id")
	teamID, _ := strconv.Atoi(teamIDString)

	projects := projectHandler.ProjectService.GetTeamProjects(userID, int64(teamID))

	var projectsDTO dto.UserProject

	for _, project := range projects {
		label := projectHandler.LabelService.GetLabelInfo(project.Label)
		projectsDTO.Projects = append(projectsDTO.Projects, dto.Project{
			ProjectID:   project.ID,
			OwnerID:     project.OwnerID,
			Title:       project.Title,
			Description: project.Description,
			Label:       *label,
			SelectedBid: project.SelectedBid,
			Status:      project.State,
			Tags:        project.Tags,
			Comment:     project.Comment,
			Duration:    project.Duration,
			StartTime:   project.StartTime,
			EndTime:     project.EndTime,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(projectsDTO); err != nil {
		panic(exceptions.Exception{
			Tag:    exceptions.INTERNAL_ERROR,
			Errors: []exceptions.SpecificError{exceptions.CAST_ERROR},
		})
	}
}

// func (projectHandler *ProjectHandler) GetOneManTeamProjects(w http.ResponseWriter, r *http.Request) {
// 	userID := r.Context().Value(projectHandler.Constants.Context.UserID).(int)

// 	userIDString := chi.URLParam(r, "id")
// 	targetUserID, _ := strconv.Atoi(userIDString)

// 	if targetUserID != 0 {
// 		userID = targetUserID
// 	}

// 	projects := projectHandler.ProjectService.GetOneManTeamProjects(userID)

// 	var projectsDTO dto.UserProject

// 	for _, project := range projects {
// 		label := projectHandler.LabelService.GetLabelInfo(project.Label)
// 		projectsDTO.Projects = append(projectsDTO.Projects, dto.Project{
// 			ProjectID:   project.ID,
// 			OwnerID:     project.OwnerID,
// 			Title:       project.Title,
// 			Description: project.Description,
// 			Label:       *label,
// 			SelectedBid: project.SelectedBid,
// 			Status:      project.State,
// 			Tags:        project.Tags,
// 			Duration:    project.Duration,
// 			StartTime:   project.StartTime,
// 			EndTime:     project.EndTime,
// 		})
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	if err := json.NewEncoder(w).Encode(projectsDTO); err != nil {
// 		panic(exceptions.Exception{
// 			Tag:    exceptions.INTERNAL_ERROR,
// 			Errors: []exceptions.SpecificError{exceptions.CAST_ERROR},
// 		})
// 	}
// }

func (ph *ProjectHandler) GetParticipatedProjectsForUser(w http.ResponseWriter, r *http.Request) {
	commanderid := r.Context().Value(ph.Constants.Context.UserID).(int)
	userIdString := chi.URLParam(r, "user_id")
	userid, _ := strconv.Atoi(userIdString)

	if userid == 0 {
		userid = commanderid
	}

	includes := r.URL.Query()["include"]

	response := make(map[string]interface{})
	var projectsTeam, projectsUser []dto.Project

	if utils.Contains(includes, "team") {
		projects := ph.ProjectService.GetParticipatedProjectsForUser(userid)

		for _, project := range projects {
			label := ph.LabelService.GetLabelInfo(project.Label)
			projectsTeam = append(projectsTeam, dto.Project{
				ProjectID:   project.ID,
				OwnerID:     project.OwnerID,
				Title:       project.Title,
				Description: project.Description,
				Label:       *label,
				SelectedBid: project.SelectedBid,
				Status:      project.State,
				Tags:        project.Tags,
				Duration:    project.Duration,
				StartTime:   project.StartTime,
				EndTime:     project.EndTime,
			})
		}
		response["team"] = projectsTeam
	}
	log.Println("Awsd")
	if utils.Contains(includes, "user") {
		projects := ph.ProjectService.GetOneManTeamProjects(userid)

		for _, project := range projects {
			label := ph.LabelService.GetLabelInfo(project.Label)
			projectsUser = append(projectsUser, dto.Project{
				ProjectID:   project.ID,
				OwnerID:     project.OwnerID,
				Title:       project.Title,
				Description: project.Description,
				Label:       *label,
				SelectedBid: project.SelectedBid,
				Status:      project.State,
				Tags:        project.Tags,
				Duration:    project.Duration,
				StartTime:   project.StartTime,
				EndTime:     project.EndTime,
			})
		}
		response["user"] = projectsUser
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		panic(exceptions.Exception{
			Tag:    exceptions.INTERNAL_ERROR,
			Errors: []exceptions.SpecificError{exceptions.CAST_ERROR},
		})
	}
}

func (ph *ProjectHandler) SearchProjects(w http.ResponseWriter, r *http.Request) {

	params := Validated[elasticmodel.QueryAndTagSearchReqDto](ph.Validator, r)

	if params.Order != "" && params.Order != "desc" && params.Order != "asc" {
		panic(exceptions.Exception{
			Tag: exceptions.BAD_REQUEST,
		})
	}

	res := ph.UserService.SearchUsers(&params)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}
