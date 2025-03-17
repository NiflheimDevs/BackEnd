package services

import (
	"time"

	"github.com/niflheimdevs/backend/internal/bootstrap"
	"github.com/niflheimdevs/backend/internal/enums"
	"github.com/niflheimdevs/backend/internal/exceptions"
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

func (projectService *ProjectService) CreateProject(userID int, title, description string, tag []int) int {
	duration := time.Now().Add(projectService.Constants.Project.LastTime).Format("2006-01-02")
	project_id, err := projectService.ProjectRepo.CreateProject(userID, title, description, duration)
	if err != nil {
		panic(exceptions.Exception{
			Tag: enums.INTERNAL_ERROR,
			Errors: []enums.SpecificError{
				enums.DATABASE_ERROR,
			},
		})
	}

	for _, tag_id := range tag {
		err = projectService.ProjectRepo.AddProjectTag(project_id, tag_id)
		if err != nil {
			panic(exceptions.Exception{
				Tag: enums.INTERNAL_ERROR,
				Errors: []enums.SpecificError{
					enums.DATABASE_ERROR,
				},
			})
		}
	}

	return project_id
}
