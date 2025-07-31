package repositories

import (
	"github.com/niflheimdevs/backend/internal/application/dto"
	"github.com/niflheimdevs/backend/internal/domain/models"
)

type CommentRepo interface {
	AddComment(projectID, bidID int, content string, star int) int
	GetStar(userID int) (float64, error)
	GetCommentCount(userID int) (int, error)
	GetUserComments(userID int) ([]dto.CommentWithUserDTO, error)
	GetCommentInfo(id int) dto.CommentWithUserDTO
	GetCommentOfProject(projectID int) (*models.CommentModel, error)
}
