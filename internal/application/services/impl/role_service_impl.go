package servicesimpl

import (
	"github.com/niflheimdevs/backend/internal/domain/enums"
	"github.com/niflheimdevs/backend/internal/domain/exceptions"
	repositories "github.com/niflheimdevs/backend/internal/domain/repositories/postgres"
)

type RoleService struct {
	RoleRepo repositories.RoleRepo
}

func NewRoleService(
// roleRepo repositories.RoleRepo,
) *RoleService {
	return &RoleService{
		// RoleRepo: roleRepo,
	}
}

func (rs *RoleService) GetAllRolesForTeam() []string {
	teamRoles := enums.GetTeamRoles()

	res := make([]string, len(teamRoles))

	for i := 0; i < len(teamRoles); i++ {
		res[i] = teamRoles[i].String()
	}
	return res
}

func (rs *RoleService) GetPermissionsForRole(role string) []string {
	roleid := enums.NameToRole(role)
	if roleid == 0 {
		panic(exceptions.Exception{
			Tag: exceptions.NOT_FOUND,
		})
	}

	permissions := roleid.GetPermissionsForRole()

	res := make([]string, len(permissions))

	for i := 0; i < len(permissions); i++ {
		res[i] = permissions[i].String()
	}
	return res
}
