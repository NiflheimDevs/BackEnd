package services

import "github.com/niflheimdevs/backend/internal/application/dto"

type CommentService interface {
	PutComment(userID, ProjectID int, content string, star int) int
	GetUserStar(userID int) float32
	GetUserComments(userID int) []dto.CommentDTO
	GetCommentInfo(id int) dto.CommentDTO
}
