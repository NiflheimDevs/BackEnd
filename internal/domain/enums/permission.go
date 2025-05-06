package enums

type PermType uint

const (
	ADD_MEMBER PermType = iota + 1
	REMOVE_MEMEBER
	EDIT_INFO
	BIDDER
	EDIT_NICKNAME
	EDIT_ROLE
	DELETE_TEAM
)

var permissionNames = map[PermType]string{
	ADD_MEMBER:     "ADD_MEMBER",
	REMOVE_MEMEBER: "REMOVE_MEMEBER",
	EDIT_INFO:      "EDIT_INFO",
	BIDDER:         "BIDDER",
	EDIT_NICKNAME:  "EDIT_NICKNAME",
	EDIT_ROLE:      "EDIT_ROLE",
	DELETE_TEAM:    "DELETE_TEAM",
}

func (r PermType) String() string {
	if name, ok := permissionNames[r]; ok {
		return name
	}
	return ""
}

func PermExists(query string) bool {
	for _, name := range permissionNames {
		if name == query {
			return true
		}
	}
	return false
}

func GetAllPerms() []PermType {
	return []PermType{
		ADD_MEMBER,
		REMOVE_MEMEBER,
		EDIT_INFO,
		BIDDER,
		EDIT_NICKNAME,
		EDIT_ROLE,
		DELETE_TEAM,
	}

}
