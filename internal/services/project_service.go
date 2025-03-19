package services

import (
	"time"

	"github.com/niflheimdevs/backend/internal/bootstrap"
	"github.com/niflheimdevs/backend/internal/enums"
	"github.com/niflheimdevs/backend/internal/exceptions"
	"github.com/niflheimdevs/backend/internal/models"
	"github.com/niflheimdevs/backend/internal/repositories"
)

type ProjectService struct {
	ProjectRepo *repositories.ProjectRepo
	Constants   *bootstrap.Constants
}

func NewProjectService(
	projectRepo *repositories.ProjectRepo,
	constants *bootstrap.Constants,
) *ProjectService {
	return &ProjectService{
		ProjectRepo: projectRepo,
		Constants:   constants,
	}
}

func (projectService *ProjectService) GetProject(projectID int) *models.ProjectModel {
	project, err := projectService.ProjectRepo.GetProject(projectID)

	if err != nil {
		panic(exceptions.Exception{
			Tag: enums.NOT_FOUND,
			Errors: []enums.SpecificError{
				enums.PROJECT_NOT_FOUND,
			},
		})
	}

	tags := projectService.ProjectRepo.GetProjectTag(projectID)

	project.Tags = tags

	return project
}

func (projectService *ProjectService) CreateProject(userID int, title, description string, tag []int) int {
	if userID == -1 || userID == -2 {
		panic(exceptions.Exception{
			Tag: enums.AUTHENTICATION_ERROR,
			Errors: []enums.SpecificError{
				enums.AUTH_ACCESS_DENIED,
			},
		})
	}

	duration := time.Now().Add(projectService.Constants.Project.LastTime).Format("2006-01-02 15:04:05")
	project_id := projectService.ProjectRepo.CreateProject(userID, title, description, duration)

	for _, tag_id := range tag {
		projectService.ProjectRepo.AddProjectTag(project_id, tag_id)
	}

	return project_id
}

func (projectService *ProjectService) UpdateProject(projectID, UserID int, title, description string, tag []int) {
	if UserID == -1 || UserID == -2 {
		panic(exceptions.Exception{
			Tag: enums.AUTHENTICATION_ERROR,
			Errors: []enums.SpecificError{
				enums.AUTH_ACCESS_DENIED,
			},
		})
	}
	//err := projectService.ProjectRepo.UpdateProject(projectID, UserID, title, description)

	projectService.ProjectRepo.DeleteProjectTags(projectID)

	for _, tag_id := range tag {
		projectService.ProjectRepo.AddProjectTag(projectID, tag_id)
	}
}
