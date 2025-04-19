package repositories

import (
	"context"

	"github.com/niflheimdevs/backend/internal/application/dto"
	"github.com/niflheimdevs/backend/internal/domain/models"
	"github.com/niflheimdevs/backend/internal/domain/repositories/postgres/transaction"
)

type ProjectRepo interface {
	LandingProps() []dto.ProjectLanding
	GetProject(projectID int) (*models.ProjectModel, error)
	GetUserProject(userID, offset, limit int) []models.ProjectModel
	CreateProject(ctx context.Context, tx transaction.Tx, userID, label int, title, description, duration string) int
	UpdateProject(projectID, UserID int, title, description string)
	DeleteProject(ctx context.Context, tx transaction.Tx, projectID int)
}
