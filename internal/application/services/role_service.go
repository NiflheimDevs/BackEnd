package services

type RoleService interface {
	GetPermissionsForRole(role string) []string
	GetAllRolesForTeam() []string
}
