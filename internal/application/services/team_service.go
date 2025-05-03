package services

import (
	"context"

	"github.com/niflheimdevs/backend/internal/application/dto"
	"github.com/niflheimdevs/backend/internal/domain/repositories/postgres/transaction"
)

type TeamService interface {
	BehindCurtainTeam(ctx context.Context, tx transaction.Tx, userid int)
	CreateTeam(userid int, teamInfo *dto.TeamCreateDto) int64
	GetTeamsForUser(userid int) []dto.GetTeamPreviewDto
	GetTeam(commanderid int, teamid int64) *dto.GetTeamDto
	UpdateTeamInfo(userid int, info *dto.UpdateTeamInfoDto)
	UpdatePosition(commanderid int, req *dto.UpdateMemberPositionDto)
	UpdateMemeberRole(commanderid int, info *dto.UpdateMemberRoleDto)
	DeleteTeam(commanderid int, teamid int64)
	// GetTeamWithRole(userid int, teamid int64) *dto.GetTeamDto
	AddMembers(userid int, teamid int64, members []int)
	LeaveTeam(userid int, teamid int64)
	KickMemebr(commanderid int, poorGuysid []int, teamid int64)
	GetMembers(commanderid int, teamid int64) []dto.SendMemberDto
	DeleteTeamProfile(commanderid int, teamid int64)
	UpdateTeamProfile(commanderid int, teamid int64, data []byte)
}
