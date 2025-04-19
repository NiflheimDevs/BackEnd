package servicesimpl

import (
	"context"
	"time"

	"github.com/niflheimdevs/backend/bootstrap"
	"github.com/niflheimdevs/backend/internal/application/dto"
	"github.com/niflheimdevs/backend/internal/application/services"
	"github.com/niflheimdevs/backend/internal/domain/exceptions"
	"github.com/niflheimdevs/backend/internal/domain/models"
	repositories "github.com/niflheimdevs/backend/internal/domain/repositories/postgres"
	"github.com/niflheimdevs/backend/internal/domain/repositories/postgres/transaction"
)

type ProjectService struct {
	ProjectRepo    repositories.ProjectRepo
	PaymentService services.PaymentService
	Constants      *bootstrap.Constants
	TxManager      transaction.TxManager
	TagRepo        repositories.TagRepo
	TagService     services.TagService
}

func NewProjectService(
	projectRepo repositories.ProjectRepo,
	paymentService services.PaymentService,
	constants *bootstrap.Constants,
	tagRepo repositories.TagRepo,
	tagService services.TagService,
	txManager transaction.TxManager,
) *ProjectService {
	return &ProjectService{
		ProjectRepo:    projectRepo,
		PaymentService: paymentService,
		Constants:      constants,
		TxManager:      txManager,
		TagRepo:        tagRepo,
		TagService:     tagService,
	}
}

func (projectService *ProjectService) LandingProps() []dto.ProjectLanding {
	projects := projectService.ProjectRepo.LandingProps()

	return projects
}

func (projectService *ProjectService) GetProject(projectID int) *models.ProjectModel {
	project, err := projectService.ProjectRepo.GetProject(projectID)

	if err != nil {
		panic(exceptions.Exception{
			Tag: exceptions.NOT_FOUND,
			Errors: []exceptions.SpecificError{
				exceptions.PROJECT_NOT_FOUND,
			},
		})
	}

	tags := projectService.TagRepo.GetProjectTag(projectID)

	project.Tags = tags

	return project
}

func (projectService *ProjectService) GetUserProjects(userID, offset, limit int) []models.ProjectModel {
	if userID == -1 || userID == -2 {
		panic(exceptions.Exception{
			Tag: exceptions.UNAUTHORIZED,
			Errors: []exceptions.SpecificError{
				exceptions.AUTH_ACCESS_DENIED,
			},
		})
	}

	projects := projectService.ProjectRepo.GetUserProject(userID, offset, limit)

	return projects
}

func (projectService *ProjectService) CreateProject(userID int, title, description string, label int, price int64, tags []int) int {
	if userID == -1 || userID == -2 {
		panic(exceptions.Exception{
			Tag: exceptions.UNAUTHORIZED,
			Errors: []exceptions.SpecificError{
				exceptions.AUTH_ACCESS_DENIED,
			},
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tx, err := projectService.TxManager.Begin(ctx)
	if err != nil {
		panic(exceptions.Exception{
			Tag: exceptions.INTERNAL_ERROR,
			Errors: []exceptions.SpecificError{
				exceptions.DATABASE_ERROR,
			},
		})
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(ctx)
			panic(p)
		}
	}()

	projectService.PaymentService.ProjectPayment(ctx, tx, userID, price)

	duration := time.Now().Add(projectService.Constants.Project.LastTime).Format("2006-01-02 15:04:05")

	projectID := projectService.ProjectRepo.CreateProject(ctx, tx, userID, label, title, description, duration)

	for _, tagID := range tags {
		projectService.TagRepo.AddProjectTagWithTx(ctx, tx, projectID, tagID)
	}

	if err := tx.Commit(ctx); err != nil {
		panic(exceptions.Exception{
			Tag: exceptions.INTERNAL_ERROR,
			Errors: []exceptions.SpecificError{
				exceptions.DATABASE_ERROR,
			},
		})
	}

	return projectID
}

func (projectService *ProjectService) UpdateProject(projectID, userID int, title, description string, label int, price int64, tags []int) {
	if userID == -1 || userID == -2 {
		panic(exceptions.Exception{
			Tag: exceptions.UNAUTHORIZED,
			Errors: []exceptions.SpecificError{
				exceptions.AUTH_ACCESS_DENIED,
			},
		})
	}

	project, err := projectService.ProjectRepo.GetProject(projectID)

	if err != nil {
		panic(exceptions.Exception{
			Tag: exceptions.NOT_FOUND,
			Errors: []exceptions.SpecificError{
				exceptions.PROJECT_NOT_FOUND,
			},
		})
	}

	if project.OwnerID != userID {
		panic(exceptions.Exception{
			Tag: exceptions.BAD_REQUEST,
			Errors: []exceptions.SpecificError{
				exceptions.USER_NOT_OWNER,
			},
		})
	}

	projectService.ProjectRepo.UpdateProject(projectID, userID, title, description)

	projectService.TagService.UpdateTagsForProject(projectID, tags) // !
}

func (projectService *ProjectService) DeleteProject(userID, projectID int) {
	if userID == -1 || userID == -2 {
		panic(exceptions.Exception{
			Tag: exceptions.UNAUTHORIZED,
			Errors: []exceptions.SpecificError{
				exceptions.AUTH_ACCESS_DENIED,
			},
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tx, err := projectService.TxManager.Begin(ctx)

	if err != nil {
		panic(exceptions.Exception{
			Tag: exceptions.INTERNAL_ERROR,
			Errors: []exceptions.SpecificError{
				exceptions.DATABASE_ERROR,
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
			Tag: exceptions.NOT_FOUND,
			Errors: []exceptions.SpecificError{
				exceptions.PROJECT_NOT_FOUND,
			},
		})
	}

	if project.OwnerID != userID {
		panic(exceptions.Exception{
			Tag: exceptions.BAD_REQUEST,
			Errors: []exceptions.SpecificError{
				exceptions.USER_NOT_OWNER,
			},
		})
	}

	tags := projectService.TagRepo.GetProjectTag(projectID)

	for _, tag := range tags {
		projectService.TagRepo.DeleteProjectTagWithTx(ctx, tx, projectID, tag.ID)
	}

	projectService.ProjectRepo.DeleteProject(ctx, tx, projectID)

	if err := tx.Commit(ctx); err != nil {
		panic(exceptions.Exception{
			Tag: exceptions.INTERNAL_ERROR,
			Errors: []exceptions.SpecificError{
				exceptions.DATABASE_ERROR,
			},
		})
	}
}
