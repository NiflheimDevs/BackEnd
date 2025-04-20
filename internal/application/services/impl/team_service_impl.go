package servicesimpl

import (
	"context"

	"github.com/niflheimdevs/backend/internal/domain/exceptions"
	repositories "github.com/niflheimdevs/backend/internal/domain/repositories/postgres"
	"github.com/niflheimdevs/backend/internal/domain/repositories/postgres/transaction"
)

type TeamService struct {
	TeamRepo repositories.TeamRepo
}

func NewTeamService(
	teamRepo repositories.TeamRepo,
) *TeamService {
	return &TeamService{
		TeamRepo: teamRepo,
	}
}

func (ts *TeamService) BehindCurtainTeam(ctx context.Context, tx transaction.Tx, userid int) {

	teamid := ts.TeamRepo.CreateOneManTeam(ctx, tx, userid)

	err := ts.TeamRepo.AddMemberWithTx(ctx, tx, userid, teamid, "")
	if err != nil {
		panic(exceptions.Exception{
			Tag:    exceptions.CONFLICT_ERROR,
			Errors: []exceptions.SpecificError{exceptions.ALREADY_A_MEMBER},
		})
	}
}
