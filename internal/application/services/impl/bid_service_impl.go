package servicesimpl

import (
	"time"

	"github.com/niflheimdevs/backend/internal/domain/exceptions"
	"github.com/niflheimdevs/backend/internal/domain/models"
	repositories "github.com/niflheimdevs/backend/internal/domain/repositories/postgres"
)

type BidService struct {
	BidRepo repositories.BidRepo
}

func NewBidService(bidRepo repositories.BidRepo) *BidService {
	return &BidService{
		BidRepo: bidRepo,
	}
}

func (bs BidService) GetBidsOfProject(projectID int) []models.BidModel {
	// check that who is the owner of project
	bids := bs.BidRepo.GetBidOfProject(projectID)
	return bids
}

func (bs BidService) PutBidOnProject(userID int, teamID int, projectID int, value int64, expected_time time.Time) int {
	if userID == -1 || userID == -2 {
		panic(exceptions.Exception{
			Tag: exceptions.UNAUTHORIZED,
			Errors: []exceptions.SpecificError{
				exceptions.AUTH_ACCESS_DENIED,
			},
		})
	}

	//check permission

	bidid := bs.BidRepo.PutBid(teamID, projectID, value, expected_time)

	return bidid
}

func (bs BidService) AcceptBid() {

}
