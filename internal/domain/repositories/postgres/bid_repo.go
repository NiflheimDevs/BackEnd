package repositories

import (
	"context"

	"github.com/niflheimdevs/backend/internal/application/dto"
	"github.com/niflheimdevs/backend/internal/domain/models"
	"github.com/niflheimdevs/backend/internal/domain/repositories/postgres/transaction"
)

type BidRepo interface {
	GetBidInfo(bidID int) (*models.BidModel, error)
	GetTeamBids(teamID int64) []models.BidModel
	GetBidOfProject(projectID int) []models.BidModel
	PutBid(info dto.BidInfo) int
	AcceptBid(ctx context.Context, tx transaction.Tx, bidID int, projectID int)
	UpdateBid(info dto.BidInfo)
}
