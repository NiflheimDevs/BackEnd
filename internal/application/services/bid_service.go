package services

import (
	"github.com/niflheimdevs/backend/internal/application/dto"
)

type BidService interface {
	GetPublicProjectBids(projectID int) []dto.PublicProjectBidInfo
	GetPrivateProjectBids(userID int, projectID int) []dto.PrivateProjectBidInfo
	GetTeamBids(userID int, teamID int64) []dto.BidInfo
	PutBidOnProject(info dto.BidInfo) int
	AcceptBid(userID int, bidID int, projectID int)
	UpdateBid(info dto.BidInfo)
}
