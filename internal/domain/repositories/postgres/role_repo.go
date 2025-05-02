package repositories

import (
	"context"

	"github.com/niflheimdevs/backend/internal/domain/models"
	"github.com/niflheimdevs/backend/internal/domain/repositories/postgres/transaction"
)

type RoleRepo interface {
	AddRole(userid int, originid int, type_ int, roleid uint) error
	AddRoleWithTx(ctx context.Context, tx transaction.Tx, userid int, originid int, type_ int, roleid uint) error
	UpdateRole(userid int, originid int, type_ int, roleid uint) error
	DeleteRole(userid int, originid int, type_ int) error
	GetRoleModel(userid int, originid int, type_ int) *models.RoleModel
	GetRoleId(name string) (uint, error)
	GetRoleName(id uint) (string, error)
	GetPermissionsForRole(roleid uint) []models.PermissionModel
}
