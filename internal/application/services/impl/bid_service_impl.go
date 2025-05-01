package servicesimpl

import (
	"context"
	"time"

	"github.com/niflheimdevs/backend/internal/application/dto"
	"github.com/niflheimdevs/backend/internal/application/services"
	"github.com/niflheimdevs/backend/internal/domain/exceptions"
	"github.com/niflheimdevs/backend/internal/domain/models"
	repositories "github.com/niflheimdevs/backend/internal/domain/repositories/postgres"
	"github.com/niflheimdevs/backend/internal/domain/repositories/postgres/transaction"
)

type BidService struct {
	ProjectService services.ProjectService
	PaymentService services.PaymentService
	TxManager      transaction.TxManager
	BidRepo        repositories.BidRepo
}

func NewBidService(
	projectservice services.ProjectService,
	paymentservice services.PaymentService,
	bidRepo repositories.BidRepo,
	txManager transaction.TxManager,
) *BidService {
	return &BidService{
		ProjectService: projectservice,
		PaymentService: paymentservice,
		BidRepo:        bidRepo,
		TxManager:      txManager,
	}
}

func (bs BidService) GetBidsOfProject(projectID int) []models.BidModel {
	// check that who is the owner of project
	bids := bs.BidRepo.GetBidOfProject(projectID)

	return bids
}

func (bs BidService) PutBidOnProject(info dto.BidInfo) int {
	if info.UserID == -1 || info.UserID == -2 {
		panic(exceptions.Exception{
			Tag: exceptions.UNAUTHORIZED,
			Errors: []exceptions.SpecificError{
				exceptions.AUTH_ACCESS_DENIED,
			},
		})
	}

	//check permission

	bs.ProjectService.GetProject(info.ProjectID)

	bidID := bs.BidRepo.PutBid(info)

	return bidID
}

func (bs BidService) AcceptBid(userID int, bidID int, projectID int) {
	if userID == -1 || userID == -2 {
		panic(exceptions.Exception{
			Tag: exceptions.UNAUTHORIZED,
			Errors: []exceptions.SpecificError{
				exceptions.AUTH_ACCESS_DENIED,
			},
		})
	}

	project := bs.ProjectService.GetProject(projectID)

	bid, err := bs.BidRepo.GetBidInfo(bidID)

	if err != nil {
		panic(exceptions.Exception{
			Tag:    exceptions.NOT_FOUND,
			Errors: []exceptions.SpecificError{exceptions.BID_NOT_FOUND},
		})
	}

	if project.OwnerID != userID {
		panic(exceptions.Exception{
			Tag:    exceptions.FORBIDDEN,
			Errors: []exceptions.SpecificError{exceptions.USER_NOT_OWNER},
		})
	}

	if project.SelectedBid != 0 {
		panic(exceptions.Exception{
			Tag:    exceptions.FORBIDDEN,
			Errors: []exceptions.SpecificError{exceptions.ALREADY_HAS_BID},
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tx, err := bs.TxManager.Begin(ctx)
	if err != nil {
		panic(exceptions.Exception{
			Tag: exceptions.INTERNAL_ERROR,
			Errors: []exceptions.SpecificError{
				exceptions.DATABASE_ERROR,
			},
		})
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(ctx)
			panic(p)
		}
	}()

	bs.PaymentService.ProjectPayment(ctx, tx, userID, bid.PrePayment)

	bs.BidRepo.AcceptBid(ctx, tx, bidID, projectID)

	if err := tx.Commit(ctx); err != nil {
		panic(exceptions.Exception{
			Tag: exceptions.INTERNAL_ERROR,
			Errors: []exceptions.SpecificError{
				exceptions.DATABASE_ERROR,
			},
		})
	}
}

func (bs BidService) UpdateBid(bidID int, info dto.BidInfo) {
	bid, err := bs.BidRepo.GetBidInfo(bidID)
	if err != nil {
		panic(exceptions.Exception{
			Tag: exceptions.NOT_FOUND,
			Errors: []exceptions.SpecificError{
				exceptions.BID_NOT_FOUND,
			},
		})
	}

	if bid.TeamID != info.TeamID {
		panic(exceptions.Exception{
			Tag: exceptions.FORBIDDEN,
			Errors: []exceptions.SpecificError{
				exceptions.USER_NOT_OWNER,
			},
		})
	}

	//check permission

	if info.PP+info.Total > bid.PrePayment+bid.Total {
		panic(exceptions.Exception{
			Tag: exceptions.FORBIDDEN,
			Errors: []exceptions.SpecificError{
				exceptions.INCREASE_NOT_ALLOWED,
			},
		})
	}

	bs.BidRepo.UpdateBid(bidID, info)
}
