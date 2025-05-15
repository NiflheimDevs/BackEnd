package repositories

import (
	"context"

	"github.com/niflheimdevs/backend/internal/application/dto"
	"github.com/niflheimdevs/backend/internal/domain/enums"
	"github.com/niflheimdevs/backend/internal/domain/models"
	"github.com/niflheimdevs/backend/internal/domain/repositories/postgres/transaction"
)

type TeamRepo interface {
	CreateOneManTeam(ctx context.Context, tx transaction.Tx, userid int) int64
	AddMember(userid int, teamid int64, position string, roleid enums.RoleType) error
	AddMemberWithTx(ctx context.Context, tx transaction.Tx, userid int, teamid int64, position string, roleid enums.RoleType) error
	RemoveMember(userid int, teamid int64) error
	UpdateMemberPosition(info *dto.UpdateMemberPositionDto) error
	UpdateMemberRole(newRole enums.RoleType, userid int, teamid int64) error
	CreateTeam(ctx context.Context, tx transaction.Tx, title string, description string) int64
	UpdateTeam(teamid int64, title string, description string) error
	DeleteTeam(teamid int64) error
	GetTeamForUser(userid int, teamid int64) *models.TeamModel
	GetTeam(teamid int64) *models.TeamModel
	GetEveryTeamInfo(teamid int64) *models.TeamModel
	GetTeamsForUser(userid int) []dto.GetTeamPreviewDto
	GetTeamsForUserWithRole(userid int) []dto.GetTeamWithRole
	GetMembersForTeamFilterdByRole(teamid int64, roleid enums.RoleType) []dto.ReadMemberDto
	GetMembersForTeam(teamid int64) []dto.ReadMemberDto
	GetMemberForTeam(teamid int64, userid int) *dto.ReadMemberDto
	GetTeamInfo(teamid int64) (*dto.GetInternalTeamInfo, error)
	GetOneManTeamInfo(teamid int64) (*dto.GetInternalTeamInfo, error)
	GetOneManTeamID(userid int) (int64, error)
	GetEveryTeamID(userID int) []int64
}
