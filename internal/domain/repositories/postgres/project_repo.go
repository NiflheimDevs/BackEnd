package repositories

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/niflheimdevs/backend/internal/domain/models"
)

type ProjectRepo interface {
	GetProject(projectID int) (*models.ProjectModel, error)
	GetUserProject(userID, offset, limit int) []models.ProjectModel
	CreateProject(ctx context.Context, tx pgx.Tx, userID, label int, title, description, duration string) int
	GetProjectTag(projectID int) []models.TagModel
	AddProjectTag(ctx context.Context, tx pgx.Tx, projectID, tagID int)
	DeleteProjectTags(ctx context.Context, tx pgx.Tx, projectID, tagID int)
	UpdateProject(ctx context.Context, tx pgx.Tx, projectID, UserID, label int, title, description string)
	DeleteProject(ctx context.Context, tx pgx.Tx, projectID int)
}
