package dto

import "github.com/niflheimdevs/backend/internal/models"

type Project struct {
	ProjectID   int    `json:"project_id"`
	OwnerID     int    `json:"Owner_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	//temp
	Label     int               `json:"label"`
	Tags      []models.TagModel `json:"tags"`
	FirstName string            `json:"first_name"`
	LastName  string            `json:"last_name"`
	Username  string            `json:"username"`
	Duration  string            `json:"duration"`
}

type CreateProjectDTO struct {
	ProjectID int
}
