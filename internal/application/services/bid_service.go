package services

import (
	"time"

	"github.com/niflheimdevs/backend/internal/domain/models"
)

type BidService interface {
	PutBidOnProject(userID int, teamID int, projectID int, value int64, expected_time time.Time) int
	GetBidsOfProject(projectID int) []models.BidModel
}
