package dto

import (
	"time"

	"github.com/niflheimdevs/backend/internal/domain/models"
)

type ProjectLanding struct {
	ProjectID   int    `json:"project_id"`
	Title       string `json:"title"`
	Description string `json:"descriptoin"`
	Label       string `json:"label"`
}

type Project struct {
	ProjectID   int                    `json:"project_id"`
	OwnerID     int                    `json:"Owner_id"`
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Label       models.LabelModel      `json:"label"`
	SelectedBid int                    `json:"selected_bid"`
	Status      int                    `json:"status"`
	Tags        []models.TagModel      `json:"tags"`
	Comment     *models.CommentModel   `json:"comment"`
	Bids        []PublicProjectBidInfo `json:"bids"`
	FirstName   string                 `json:"first_name"`
	LastName    string                 `json:"last_name"`
	Username    string                 `json:"username"`
	Duration    time.Time              `json:"duration"`
	StartTime   time.Time              `json:"start_time"`
	EndTime     time.Time              `json:"end_time"`
}

type CreateProjectDTO struct {
	ProjectID int
}

type UserProject struct {
	Projects []Project `json:"projects"`
	Count    int       `json:"count"`
}
