package services

import (
	"github.com/niflheimdevs/backend/internal/application/dto"
	"github.com/niflheimdevs/backend/internal/domain/models"
)

type BidService interface {
	PutBidOnProject(info dto.BidInfo) int
	GetBidsOfProject(projectID int) []models.BidModel
	AcceptBid(userID int, bidID int, projectID int)
}
