package services

import (
	"encoding/json"
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

	var labels []string
	if err := json.Unmarshal([]byte(project.Label[0]), &labels); err != nil {
		panic(exceptions.Exception{
			Tag: enums.INTERNAL_ERROR,
			Errors: []enums.SpecificError{
				enums.CAST_ERROR,
			},
		})
	}
	project.Label = labels

	tags := projectService.ProjectRepo.GetProjectTag(projectID)

	project.Tags = tags

	return project
}

func (projectService *ProjectService) GetUserProjects(userID, offset, limit int) []models.ProjectModel {
	if userID == -1 || userID == -2 {
		panic(exceptions.Exception{
			Tag: enums.UNAUTHORIZED,
			Errors: []enums.SpecificError{
				enums.AUTH_ACCESS_DENIED,
			},
		})
	}

	projects := projectService.ProjectRepo.GetUserProject(userID, offset, limit)

	for i := range projects {
		var labels []string
		if err := json.Unmarshal([]byte(projects[i].Label[0]), &labels); err != nil {
			panic(exceptions.Exception{
				Tag: enums.INTERNAL_ERROR,
				Errors: []enums.SpecificError{
					enums.CAST_ERROR,
				},
			})
		}
		projects[i].Label = labels
	}

	return projects
}

func (projectService *ProjectService) CreateProject(userID int, title, description string, label []string, tags []int) int {
	if userID == -1 || userID == -2 {
		panic(exceptions.Exception{
			Tag: enums.UNAUTHORIZED,
			Errors: []enums.SpecificError{
				enums.AUTH_ACCESS_DENIED,
			},
		})
	}

	labelsJSON, _ := json.Marshal(label)
	labels := string(labelsJSON)

	duration := time.Now().Add(projectService.Constants.Project.LastTime).Format("2006-01-02 15:04:05")
	project_id := projectService.ProjectRepo.CreateProject(userID, title, description, labels, duration)

	for _, tag_id := range tags {
		projectService.ProjectRepo.AddProjectTag(project_id, tag_id)
	}

	return project_id
}

func (projectService *ProjectService) UpdateProject(projectID, userID int, title, description string, label []string, tags []int) {
	if userID == -1 || userID == -2 {
		panic(exceptions.Exception{
			Tag: enums.UNAUTHORIZED,
			Errors: []enums.SpecificError{
				enums.AUTH_ACCESS_DENIED,
			},
		})
	}

	project, err := projectService.ProjectRepo.GetProject(projectID)

	if err != nil {
		panic(exceptions.Exception{
			Tag: enums.NOT_FOUND,
			Errors: []enums.SpecificError{
				enums.PROJECT_NOT_FOUND,
			},
		})
	}

	if project.OwnerID != userID {
		panic(exceptions.Exception{
			Tag: enums.BAD_REQUEST,
			Errors: []enums.SpecificError{
				enums.USER_NOT_OWNER,
			},
		})
	}

	labelsJSON, _ := json.Marshal(label)
	labels := string(labelsJSON)

	projectService.ProjectRepo.UpdateProject(projectID, userID, title, description, labels)

	existingTags := projectService.ProjectRepo.GetProjectTag(projectID)

	// ? extract method? making it util? this code is also needed in in general_service.go
	existingTagSet := make(map[int]bool)
	for _, tag := range existingTags {
		existingTagSet[tag.ID] = true
	}

	newTagSet := make(map[int]bool)
	for _, tag := range tags {
		newTagSet[tag] = true
	}

	for _, tag := range tags {
		if !existingTagSet[tag] {
			projectService.ProjectRepo.AddProjectTag(projectID, tag)
		}
	}

	for _, tag := range existingTags {
		if !newTagSet[tag.ID] {
			projectService.ProjectRepo.DeleteProjectTags(projectID, tag.ID)
		}
	}
}

func (projectService *ProjectService) DeleteProject(userID, projectID int) {
	if userID == -1 || userID == -2 {
		panic(exceptions.Exception{
			Tag: enums.UNAUTHORIZED,
			Errors: []enums.SpecificError{
				enums.AUTH_ACCESS_DENIED,
			},
		})
	}

	project, err := projectService.ProjectRepo.GetProject(projectID)

	if err != nil {
		panic(exceptions.Exception{
			Tag: enums.NOT_FOUND,
			Errors: []enums.SpecificError{
				enums.PROJECT_NOT_FOUND,
			},
		})
	}

	if project.OwnerID != userID {
		panic(exceptions.Exception{
			Tag: enums.BAD_REQUEST,
			Errors: []enums.SpecificError{
				enums.USER_NOT_OWNER,
			},
		})
	}

	tags := projectService.ProjectRepo.GetProjectTag(projectID)

	for _, tag := range tags {
		projectService.ProjectRepo.DeleteProjectTags(projectID, tag.ID)
	}

	projectService.ProjectRepo.DeleteProject(projectID)
}
