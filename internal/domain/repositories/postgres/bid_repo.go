package repositories

import (
	"time"

	"github.com/niflheimdevs/backend/internal/domain/models"
)

type BidRepo interface {
	PutBid(teamID int, projectID int, pp int64, total int64, description string, expected_time time.Time) int
	GetBidOfProject(projectID int) []models.BidModel
	AcceptBid(bidID int, projectID int)
}
