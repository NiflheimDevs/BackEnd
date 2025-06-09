package services

import (
	"github.com/niflheimdevs/backend/internal/application/dto"
)

type CommentService interface {
	PutComment(userID, ProjectID int, content string, star int) int
	GetUserStar(userID int) float64
	GetUserComments(userID int) []dto.CommentWithUserDTO
	GetCommentInfo(id int) *dto.CommentWithUserDTO
	GetCommentOfProject(projectID int) *dto.CommentDTO
}
