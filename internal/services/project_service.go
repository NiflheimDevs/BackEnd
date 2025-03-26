package services

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
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

	return projects
}

func (projectService *ProjectService) CreateProject(userID int, title, description string, label int, tags []int) int {
	if userID == -1 || userID == -2 {
		panic(exceptions.Exception{
			Tag: enums.UNAUTHORIZED,
			Errors: []enums.SpecificError{
				enums.AUTH_ACCESS_DENIED,
			},
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tx, err := projectService.ProjectRepo.PG.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		panic(exceptions.Exception{
			Tag: enums.INTERNAL_ERROR,
			Errors: []enums.SpecificError{
				enums.DATABASE_ERROR,
			},
		})
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(ctx)
			panic(p)
		}
	}()

	duration := time.Now().Add(projectService.Constants.Project.LastTime).Format("2006-01-02 15:04:05")

	projectID := projectService.ProjectRepo.CreateProject(ctx, tx, userID, label, title, description, duration)

	for _, tagID := range tags {
		projectService.ProjectRepo.AddProjectTag(ctx, tx, projectID, tagID)
	}

	if err := tx.Commit(ctx); err != nil {
		panic(exceptions.Exception{
			Tag: enums.INTERNAL_ERROR,
			Errors: []enums.SpecificError{
				enums.DATABASE_ERROR,
			},
		})
	}

	return projectID
}

func (projectService *ProjectService) UpdateProject(projectID, userID int, title, description string, label int, tags []int) {
	if userID == -1 || userID == -2 {
		panic(exceptions.Exception{
			Tag: enums.UNAUTHORIZED,
			Errors: []enums.SpecificError{
				enums.AUTH_ACCESS_DENIED,
			},
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tx, err := projectService.ProjectRepo.PG.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		panic(exceptions.Exception{
			Tag: enums.INTERNAL_ERROR,
			Errors: []enums.SpecificError{
				enums.DATABASE_ERROR,
			},
		})
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(ctx)
			panic(p)
		}
	}()

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

	projectService.ProjectRepo.UpdateProject(ctx, tx, projectID, userID, label, title, description)

	existingTags := projectService.ProjectRepo.GetProjectTag(projectID)

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
			projectService.ProjectRepo.AddProjectTag(ctx, tx, projectID, tag)
		}
	}

	for _, tag := range existingTags {
		if !newTagSet[tag.ID] {
			projectService.ProjectRepo.DeleteProjectTags(ctx, tx, projectID, tag.ID)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		panic(exceptions.Exception{
			Tag: enums.INTERNAL_ERROR,
			Errors: []enums.SpecificError{
				enums.DATABASE_ERROR,
			},
		})
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

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tx, err := projectService.ProjectRepo.PG.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		panic(exceptions.Exception{
			Tag: enums.INTERNAL_ERROR,
			Errors: []enums.SpecificError{
				enums.DATABASE_ERROR,
			},
		})
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(ctx)
			panic(p)
		}
	}()

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
		projectService.ProjectRepo.DeleteProjectTags(ctx, tx, projectID, tag.ID)
	}

	projectService.ProjectRepo.DeleteProject(ctx, tx, projectID)

	if err := tx.Commit(ctx); err != nil {
		panic(exceptions.Exception{
			Tag: enums.INTERNAL_ERROR,
			Errors: []enums.SpecificError{
				enums.DATABASE_ERROR,
			},
		})
	}
}
