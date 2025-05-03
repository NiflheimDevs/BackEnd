package servicesimpl

import (
	"context"
	"log"
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
	TeamService    services.TeamService
	TxManager      transaction.TxManager
	BidRepo        repositories.BidRepo
}

func NewBidService(
	projectservice services.ProjectService,
	paymentservice services.PaymentService,
	teamservice services.TeamService,
	bidRepo repositories.BidRepo,
	txManager transaction.TxManager,
) *BidService {
	return &BidService{
		ProjectService: projectservice,
		PaymentService: paymentservice,
		TeamService:    teamservice,
		BidRepo:        bidRepo,
		TxManager:      txManager,
	}
}

func (bs BidService) GetPublicProjectBids(projectID int) []dto.PublicProjectBidInfo {
	var bidInfos []dto.PublicProjectBidInfo
	bids := bs.BidRepo.GetBidOfProject(projectID)
	log.Println("saman 1")

	for _, bid := range bids {
		var bidInfo dto.PublicProjectBidInfo
		teamInfo := bs.TeamService.GetInternalTeamInfo(bid.TeamID)
		log.Println("saman 2")
		bidInfo.BidID = bid.ID
		bidInfo.TeamInfo = teamInfo
		bidInfo.Total = bid.Total
		bidInfo.ExpectedTime = bid.ExpectedTime
		bidInfos = append(bidInfos, bidInfo)
	}
	log.Println("saman 3")
	return bidInfos
}

func (bs BidService) GetPrivateProjectBids(userID int, projectID int) []dto.PrivateProjectBidInfo {
	if userID == -1 || userID == -2 {
		panic(exceptions.Exception{
			Tag:    exceptions.UNAUTHORIZED,
			Errors: []exceptions.SpecificError{exceptions.AUTH_ACCESS_DENIED},
		})
	}

	var bidInfos []dto.PrivateProjectBidInfo
	bids := bs.BidRepo.GetBidOfProject(projectID)

	for _, bid := range bids {
		var bidInfo dto.PrivateProjectBidInfo
		teamInfo := bs.TeamService.GetInternalTeamInfo(bid.TeamID)
		bidInfo.BidID = bid.ID
		bidInfo.TeamInfo = teamInfo
		bidInfo.Total = bid.Total
		bidInfo.ExpectedTime = bid.ExpectedTime
		bidInfos = append(bidInfos, bidInfo)
	}

	return bidInfos
}

func (bs BidService) GetTeamBid(teamID int64) []models.BidModel {
	bids := bs.BidRepo.GetTeamBids(teamID)

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

	if info.PP > info.Total {
		panic(exceptions.Exception{
			Tag: exceptions.FORBIDDEN,
		})
	}

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

	if project.State != 2 {
		panic(exceptions.Exception{
			Tag:    exceptions.FORBIDDEN,
			Errors: []exceptions.SpecificError{exceptions.ALREADY_HAS_BID},
		})
	}

	//check project state

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

func (bs BidService) UpdateBid(info dto.BidInfo) {
	bid, err := bs.BidRepo.GetBidInfo(info.BidID)
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

	//check permission that user is in group and have permission to edit

	if info.Total > bid.Total {
		panic(exceptions.Exception{
			Tag: exceptions.FORBIDDEN,
			Errors: []exceptions.SpecificError{
				exceptions.INCREASE_NOT_ALLOWED,
			},
		})
	}

	project := bs.ProjectService.GetProject(bid.ProjectID)

	if project.State >= 2 {
		panic(exceptions.Exception{
			Tag: exceptions.FORBIDDEN,
		})
	}

	//check project state

	bs.BidRepo.UpdateBid(info)
}
