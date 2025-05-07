package services

import (
	"github.com/niflheimdevs/backend/internal/application/dto"
	"github.com/niflheimdevs/backend/internal/domain/models"
)

type ProjectService interface {
	LandingProps() []dto.ProjectLanding
	GetProject(projectID int) *models.ProjectModel
	GetUserProjects(userID, targetuserID, offset, limit int) ([]models.ProjectModel, int)
	GetProjectCount(userID int) int
	CreateProject(userID int, title, description string, label int, price int64, tags []int) int
	UpdateProject(projectID, userID int, title, description string, label int, price int64, tags []int)
	DeleteProject(userID, projectID int)
	EndOfProject(userID, projectID int)
	GetTeamProjects(userID int, teamID int64) []models.ProjectModel
	GetOneManTeamProjects(userID int) []models.ProjectModel
}
