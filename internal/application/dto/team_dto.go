package dto

import (
	"github.com/niflheimdevs/backend/internal/domain/enums"
	"github.com/niflheimdevs/backend/internal/domain/models"
)

type TeamCreateDto struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
	Members     []int  `json:"members"`
}

// type TeamPreviewDto struct {
// 	ID          int64  `json:"id"`
// 	Title       string `json:"title"`
// 	Description string `json:"description"`
// 	Position    string `json:"position"`
// }

type GetInternalTeamInfo struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Type        int    `json:"type"`
	Profile     string `json:"profile"`
}

type GetTeamPreviewDto struct {
	ID          int64          `json:"id"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Position    string         `json:"position"`
	Profile     string         `json:"profile"`
	OwnerInfo   *MemberInfoDto `json:"owner"`
}

type UpdateTeamInfoDto struct {
	Title       string `json:"title" validate:"required,lt=25"`
	Description string `json:"description" validate:"required,lt=4000"`
	ID          int64  `json:"id" validate:"required,numeric"`
}

type UpdateMemberRoleDto struct {
	UserID int    `json:"user_id" validate:"required,numeric"`
	TeamID int64  `json:"team_id" validate:"required,numeric"`
	Role   string `json:"role" validate:"required"`
}

type UpdateMemberPositionDto struct {
	Userid      int    `json:"user_id" validate:"required,numeric"`
	NewPosition string `json:"position" validate:"required,lt=20"`
	Teamid      int64  `json:"team_id" validate:"required,numeric"`
}

type MemberInfoDto struct {
	Userid    int    `json:"userid"`
	Username  string `json:"username"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Position  string `json:"position"`
}

type ReadMemberDto struct {
	Info *MemberInfoDto
	Role enums.RoleType
}

type SendMemberDto struct {
	Info    *MemberInfoDto `json:"member_info"`
	Profile string         `json:"profile"`
	Role    string         `json:"role"`
	// Permissions []string `json:"permissions"`
}

type GetTeamDto struct {
	Info        *models.TeamModel `json:"team_info"`
	Members     []SendMemberDto   `json:"members"`
	UserID      int               `json:"user_id"`
	Permissions []string          `json:"permissions"`
}
