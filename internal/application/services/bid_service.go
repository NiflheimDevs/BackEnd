package services

import (
	"github.com/niflheimdevs/backend/internal/application/dto"
)

type BidService interface {
	PutBidOnProject(info dto.BidInfo) int
	AcceptBid(userID int, bidID int, projectID int)
	UpdateBid(info dto.BidInfo)
}
