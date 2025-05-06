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
	BidRepo        repositories.BidRepo
	TagService     services.TagService
	TeamService    services.TeamService
}

func NewProjectService(
	projectRepo repositories.ProjectRepo,
	paymentService services.PaymentService,
	constants *bootstrap.Constants,
	tagRepo repositories.TagRepo,
	bidRepo repositories.BidRepo,
	tagService services.TagService,
	teamService services.TeamService,
	txManager transaction.TxManager,
) *ProjectService {
	return &ProjectService{
		ProjectRepo:    projectRepo,
		PaymentService: paymentService,
		Constants:      constants,
		TxManager:      txManager,
		TagRepo:        tagRepo,
		BidRepo:        bidRepo,
		TagService:     tagService,
		TeamService:    teamService,
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

func (projectService *ProjectService) GetUserProjects(userID, targetuserID, offset, limit int) ([]models.ProjectModel, int) {
	var projects []models.ProjectModel
	var count int

	if targetuserID == 0 {
		if userID == -1 || userID == -2 {
			panic(exceptions.Exception{
				Tag: exceptions.UNAUTHORIZED,
				Errors: []exceptions.SpecificError{
					exceptions.AUTH_ACCESS_DENIED,
				},
			})
		}
		projects = projectService.ProjectRepo.GetUserProject(userID, offset, limit)
		count = projectService.GetProjectCount(userID)
		for _, project := range projects {
			if project.State == 1 && project.Duration.Before(time.Now()) {
				projectService.ProjectRepo.UpdateProjectState(project.ID, 2)
			}
			project.State = 2
		}
	} else {
		projects = projectService.ProjectRepo.GetUserProject(targetuserID, offset, limit)
		count = projectService.GetProjectCount(targetuserID)
	}

	return projects, count
}

func (projectService *ProjectService) GetProjectCount(userID int) int {
	count := projectService.ProjectRepo.GetProjectCount(userID)
	return count
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

	if project.State > 1 {
		panic(exceptions.Exception{
			Tag: exceptions.FORBIDDEN,
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

	if project.State > 2 {
		panic(exceptions.Exception{
			Tag: exceptions.FORBIDDEN,
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

func (projectService *ProjectService) EndOfProject(userID, projectID int) {
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

	if project.State != 3 {
		panic(exceptions.Exception{
			Tag: exceptions.FORBIDDEN,
		})
	}

	projectService.ProjectRepo.UpdateProjectState(projectID, 4)
}

func (projectService *ProjectService) GetTeamProjects(userID int, teamID int64) []models.ProjectModel {
	if userID == -1 || userID == -2 {
		panic(exceptions.Exception{
			Tag: exceptions.UNAUTHORIZED,
			Errors: []exceptions.SpecificError{
				exceptions.AUTH_ACCESS_DENIED,
			},
		})
	}

	//check that member is in team

	bids := projectService.BidRepo.GetTeamBids(teamID)
	var projects []models.ProjectModel
	for _, bid := range bids {
		project, _ := projectService.ProjectRepo.GetProject(bid.ProjectID)

		projects = append(projects, *project)
	}
	return projects
}

func (projectService *ProjectService) GetOneManTeamProjects(userID int) []models.ProjectModel {
	if userID == -1 || userID == -2 {
		panic(exceptions.Exception{
			Tag: exceptions.UNAUTHORIZED,
			Errors: []exceptions.SpecificError{
				exceptions.AUTH_ACCESS_DENIED,
			},
		})
	}

	oneManTeamID := projectService.TeamService.GetOneManTeamID(userID)

	var projects []models.ProjectModel

	oneManTeamProjects := projectService.GetTeamProjects(userID, oneManTeamID)
	projects = append(projects, oneManTeamProjects...)

	return projects
}
