package services

import (
	"github.com/niflheimdevs/backend/internal/domain/models"
)

type ProjectService interface {
	GetProject(projectID int) *models.ProjectModel
	GetUserProjects(userID, offset, limit int) []models.ProjectModel
	CreateProject(userID int, title, description string, label int, price int64, tags []int) int
	UpdateProject(projectID, userID int, title, description string, label int, price int64, tags []int)
	DeleteProject(userID, projectID int)
}
