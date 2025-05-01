package repositories

import (
	"context"

	"github.com/niflheimdevs/backend/internal/application/dto"
	"github.com/niflheimdevs/backend/internal/domain/models"
	"github.com/niflheimdevs/backend/internal/domain/repositories/postgres/transaction"
)

type BidRepo interface {
	GetBidInfo(bidID int) (*models.BidModel, error)
	PutBid(info dto.BidInfo) int
	GetBidOfProject(projectID int) []models.BidModel
	AcceptBid(ctx context.Context, tx transaction.Tx, bidID int, projectID int)
	UpdateBid(bidID int, info dto.BidInfo)
}
