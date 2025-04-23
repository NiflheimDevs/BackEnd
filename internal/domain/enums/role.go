package enums

type RoleType uint

const (
	TEAM_OWNER RoleType = iota + 1
	TEAM_ADMIN
	TEAM_CRAWLER
	TEAM_MAINTAINER
)

var roleNames = map[RoleType]string{
	TEAM_OWNER:      "TEAM_OWNER",
	TEAM_ADMIN:      "TEAM_ADMIN",
	TEAM_CRAWLER:    "TEAM_CRAWLER",
	TEAM_MAINTAINER: "TEAM_MAINTAINER",
}

var rolePermissions = map[RoleType][]PermType{
	TEAM_ADMIN:      {ADD_MEMBER, REMOVE_MEMEBER, BIDDER, EDIT_INFO, EDIT_NICKNAME, EDIT_ROLE},
	TEAM_OWNER:      {ADD_MEMBER, REMOVE_MEMEBER, BIDDER, EDIT_INFO, EDIT_NICKNAME, EDIT_ROLE, DELETE_TEAM},
	TEAM_CRAWLER:    {BIDDER},
	TEAM_MAINTAINER: {EDIT_INFO, EDIT_NICKNAME, ADD_MEMBER, REMOVE_MEMEBER},
}

func (r RoleType) String() string {
	if name, ok := roleNames[r]; ok {
		return name
	}
	return ""
}

func (r RoleType) GetPermissionsForRole() []PermType {
	if slice, ok := rolePermissions[r]; ok {
		return slice
	}
	return nil
}

func RoleExists(query string) bool {
	for _, name := range roleNames {
		if name == query {
			return true
		}
	}
	return false
}

func GetAllRoles() []RoleType {
	return []RoleType{
		TEAM_OWNER,
		TEAM_ADMIN,
		TEAM_CRAWLER,
		TEAM_MAINTAINER,
	}

}
