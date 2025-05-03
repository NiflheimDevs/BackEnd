package services

import (
	"github.com/niflheimdevs/backend/internal/application/dto"
	"github.com/niflheimdevs/backend/internal/domain/models"
)

type BidService interface {
	GetPublicProjectBids(projectID int) []dto.PublicProjectBidInfo
	GetPrivateProjectBids(userID int, projectID int) []dto.PrivateProjectBidInfo
	GetTeamBid(teamID int64) []models.BidModel
	PutBidOnProject(info dto.BidInfo) int
	AcceptBid(userID int, bidID int, projectID int)
	UpdateBid(info dto.BidInfo)
}
