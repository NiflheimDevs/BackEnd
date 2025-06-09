package repositories

import "github.com/niflheimdevs/backend/internal/application/dto"

type CommentRepo interface {
	AddComment(projectID, bidID int, content string, star int) int
	GetStar(userID int) (float32, error)
	GetUserComments(userID int) ([]dto.CommentDTO, error)
	GetCommentInfo(id int) dto.CommentDTO
}
