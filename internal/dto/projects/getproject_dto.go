package dto

import "github.com/niflheimdevs/backend/internal/models"

type Project struct {
	ProjectID   int               `json:"project_id"`
	OwnerID     int               `json:"Owner_id"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Tags        []models.TagModel `json:"tags"`
}
