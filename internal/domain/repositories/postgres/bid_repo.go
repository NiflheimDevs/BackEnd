package repositories

import (
	"time"

	"github.com/niflheimdevs/backend/internal/domain/models"
)

type BidRepo interface {
	PutBid(teamID int, projectID int, value int64, expected_time time.Time) int
	GetBidOfProject(projectID int) []models.BidModel
}
