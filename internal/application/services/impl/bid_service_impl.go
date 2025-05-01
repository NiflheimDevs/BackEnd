package servicesimpl

import (
	"time"

	"github.com/niflheimdevs/backend/internal/application/services"
	"github.com/niflheimdevs/backend/internal/domain/exceptions"
	"github.com/niflheimdevs/backend/internal/domain/models"
	repositories "github.com/niflheimdevs/backend/internal/domain/repositories/postgres"
)

type BidService struct {
	ProjectService services.ProjectService
	BidRepo        repositories.BidRepo
}

func NewBidService(
	projectservice services.ProjectService,
	bidRepo repositories.BidRepo,
) *BidService {
	return &BidService{
		ProjectService: projectservice,
		BidRepo:        bidRepo,
	}
}

func (bs BidService) GetBidsOfProject(projectID int) []models.BidModel {
	// check that who is the owner of project
	bids := bs.BidRepo.GetBidOfProject(projectID)

	return bids
}

func (bs BidService) PutBidOnProject(userID int, teamID int, projectID int, pp int64, total int64, description string, expected_time time.Time) int {
	if userID == -1 || userID == -2 {
		panic(exceptions.Exception{
			Tag: exceptions.UNAUTHORIZED,
			Errors: []exceptions.SpecificError{
				exceptions.AUTH_ACCESS_DENIED,
			},
		})
	}

	//check permission

	bidID := bs.BidRepo.PutBid(teamID, projectID, pp, total, description, expected_time)

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

	bs.BidRepo.AcceptBid(bidID, projectID)
}
