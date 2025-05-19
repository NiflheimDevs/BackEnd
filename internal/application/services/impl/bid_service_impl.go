package servicesimpl

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/niflheimdevs/backend/internal/application/dto"
	"github.com/niflheimdevs/backend/internal/application/services"
	"github.com/niflheimdevs/backend/internal/domain/enums"
	"github.com/niflheimdevs/backend/internal/domain/exceptions"
	repositories "github.com/niflheimdevs/backend/internal/domain/repositories/postgres"
	"github.com/niflheimdevs/backend/internal/domain/repositories/postgres/transaction"
	"github.com/niflheimdevs/backend/internal/utils"
)

type BidService struct {
	ProjectService services.ProjectService
	PaymentService services.PaymentService
	TeamService    services.TeamService
	TxManager      transaction.TxManager
	BidRepo        repositories.BidRepo
	TeamRepo       repositories.TeamRepo
}

func NewBidService(
	projectservice services.ProjectService,
	paymentservice services.PaymentService,
	teamservice services.TeamService,
	bidRepo repositories.BidRepo,
	teamRepo repositories.TeamRepo,
	txManager transaction.TxManager,
) *BidService {
	return &BidService{
		ProjectService: projectservice,
		PaymentService: paymentservice,
		TeamService:    teamservice,
		BidRepo:        bidRepo,
		TeamRepo:       teamRepo,
		TxManager:      txManager,
	}
}

func (bs BidService) GetPublicProjectBids(projectID int) []dto.PublicProjectBidInfo {
	bs.ProjectService.GetProject(projectID)

	var bidInfos []dto.PublicProjectBidInfo
	bids := bs.BidRepo.GetBidOfProject(projectID)

	for _, bid := range bids {
		var bidInfo dto.PublicProjectBidInfo
		teamInfo := bs.TeamService.GetInternalTeamInfo(bid.TeamID)
		bidInfo.BidID = bid.ID
		bidInfo.TeamInfo = teamInfo
		bidInfo.PrePayment = bid.PrePayment
		bidInfo.Total = bid.Total
		bidInfo.Description = bid.Description
		bidInfo.ExpectedTime = bid.ExpectedTime
		bidInfos = append(bidInfos, bidInfo)
	}
	return bidInfos
}

func (bs BidService) GetPrivateProjectBids(userID int, projectID int) []dto.PrivateProjectBidInfo {
	if userID == -1 || userID == -2 {
		panic(exceptions.Exception{
			Tag:    exceptions.UNAUTHORIZED,
			Errors: []exceptions.SpecificError{exceptions.AUTH_ACCESS_DENIED},
		})
	}

	project := bs.ProjectService.GetProject(projectID)

	if project.OwnerID != userID {
		panic(exceptions.Exception{
			Tag: exceptions.FORBIDDEN,
		})
	}

	var bidInfos []dto.PrivateProjectBidInfo
	bids := bs.BidRepo.GetBidOfProject(projectID)

	for _, bid := range bids {
		var bidInfo dto.PrivateProjectBidInfo
		teamInfo := bs.TeamService.GetInternalTeamInfo(bid.TeamID)
		bidInfo.BidID = bid.ID
		bidInfo.TeamInfo = teamInfo
		bidInfo.PrePayment = bid.PrePayment
		bidInfo.Total = bid.Total
		bidInfo.ExpectedTime = bid.ExpectedTime
		bidInfos = append(bidInfos, bidInfo)
	}

	return bidInfos
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

	teamInfo := bs.TeamRepo.GetEveryTeamInfo(info.TeamID)

	if teamInfo == nil {
		panic(exceptions.Exception{
			Tag: exceptions.NOT_FOUND,
		})
	}

	if strconv.Itoa(info.UserID) != teamInfo.Title {
		member := bs.TeamRepo.GetMemberForTeam(info.TeamID, info.UserID)
		if member == nil {
			panic(exceptions.Exception{
				Tag:    exceptions.FORBIDDEN,
				Errors: []exceptions.SpecificError{exceptions.NOT_A_MEMBER},
			})
		}

		if !utils.Contains(member.Role.GetPermissionsForRole(), enums.BIDDER) {
			panic(exceptions.Exception{
				Tag:    exceptions.FORBIDDEN,
				Errors: []exceptions.SpecificError{exceptions.LACKS_PERMISSION},
			})
		}
	}

	if info.PP > info.Total {
		panic(exceptions.Exception{
			Tag: exceptions.FORBIDDEN,
		})
	}

	project := bs.ProjectService.GetProject(info.ProjectID)

	if project.State > 1 {
		panic(exceptions.Exception{
			Tag:    exceptions.FORBIDDEN,
			Errors: []exceptions.SpecificError{exceptions.NOT_PROPER_PROJECT_STATE},
		})
	}

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

	ownerID := bs.TeamService.GetInternalTeamInfo(bid.TeamID).OwnerID

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
			Errors: []exceptions.SpecificError{exceptions.NOT_PROPER_PROJECT_STATE},
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

	description := fmt.Sprintf("Pre Payment For Project %d", projectID)

	bs.PaymentService.TransferMoney(ctx, tx, userID, ownerID, bid.PrePayment, description)

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
	if info.UserID == -1 || info.UserID == -2 {
		panic(exceptions.Exception{
			Tag: exceptions.UNAUTHORIZED,
			Errors: []exceptions.SpecificError{
				exceptions.AUTH_ACCESS_DENIED,
			},
		})
	}

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

	teamInfo := bs.TeamRepo.GetEveryTeamInfo(info.TeamID)

	if teamInfo == nil {
		panic(exceptions.Exception{
			Tag: exceptions.NOT_FOUND,
		})
	}

	if strconv.Itoa(info.UserID) != teamInfo.Title {
		member := bs.TeamRepo.GetMemberForTeam(info.TeamID, info.UserID)
		if member == nil {
			panic(exceptions.Exception{
				Tag:    exceptions.FORBIDDEN,
				Errors: []exceptions.SpecificError{exceptions.NOT_A_MEMBER},
			})
		}

		if !utils.Contains(member.Role.GetPermissionsForRole(), enums.BIDDER) {
			panic(exceptions.Exception{
				Tag:    exceptions.FORBIDDEN,
				Errors: []exceptions.SpecificError{exceptions.LACKS_PERMISSION},
			})
		}
	}

	if info.PP > info.Total {
		panic(exceptions.Exception{
			Tag: exceptions.FORBIDDEN,
		})
	}

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
			Tag:    exceptions.FORBIDDEN,
			Errors: []exceptions.SpecificError{exceptions.NOT_PROPER_PROJECT_STATE},
		})
	}

	bs.BidRepo.UpdateBid(info)
}

func (bs BidService) GetTeamBids(userID int, teamID int64) []dto.BidInfo {
	if userID == -1 || userID == -2 {
		panic(exceptions.Exception{
			Tag: exceptions.UNAUTHORIZED,
			Errors: []exceptions.SpecificError{
				exceptions.AUTH_ACCESS_DENIED,
			},
		})
	}

	teamInfo := bs.TeamRepo.GetEveryTeamInfo(teamID)

	if teamInfo == nil {
		panic(exceptions.Exception{
			Tag: exceptions.NOT_FOUND,
		})
	}

	if strconv.Itoa(userID) != teamInfo.Title {
		member := bs.TeamRepo.GetMemberForTeam(teamID, userID)
		if member == nil {
			panic(exceptions.Exception{
				Tag:    exceptions.FORBIDDEN,
				Errors: []exceptions.SpecificError{exceptions.NOT_A_MEMBER},
			})
		}

		if !utils.Contains(member.Role.GetPermissionsForRole(), enums.BIDDER) {
			panic(exceptions.Exception{
				Tag:    exceptions.FORBIDDEN,
				Errors: []exceptions.SpecificError{exceptions.LACKS_PERMISSION},
			})
		}
	}

	bids := bs.BidRepo.GetTeamBids(teamID)

	var bidDTOs []dto.BidInfo

	for _, bid := range bids {
		project := bs.ProjectService.GetProject(bid.ProjectID)
		var status int
		if project.State == 1 {
			status = 1
		} else if project.State == 2 && project.SelectedBid == bid.ID {
			status = 2
		} else {
			status = 3
		}
		bidDTO := dto.BidInfo{
			BidID:        bid.ID,
			TeamID:       teamID,
			ProjectID:    bid.ProjectID,
			PP:           bid.PrePayment,
			Total:        bid.Total,
			Status:       status,
			Description:  bid.Description,
			ExpectedTime: bid.ExpectedTime,
		}
		bidDTOs = append(bidDTOs, bidDTO)
	}

	return bidDTOs
}
