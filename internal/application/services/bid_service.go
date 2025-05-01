package services

import (
	"time"

	"github.com/niflheimdevs/backend/internal/domain/models"
)

type BidService interface {
	PutBidOnProject(userID int, teamID int, projectID int, pp int64, total int64, description string, expected_time time.Time) int
	GetBidsOfProject(projectID int) []models.BidModel
	AcceptBid(userID int, bidID int, projectID int)
}
