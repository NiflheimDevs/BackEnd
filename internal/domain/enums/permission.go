package enums

type PermType uint

/*
ALTER SEQUENCE permission_id_seq RESTART WITH 0;

INSERT INTO "permission" ("name" , "description") VALUES
("ADD_MEMBER", "adds member"),
("REMOVE_MEMEBER", "removes a member"),
("EDIT_INFO", "edit title, bio and ..."),
("BIDDER", "the one who bids"),
("EDIT_NICKNAME", "for teams, it works for positions. for groups and etc for nickname"),
("EDIT_ROLE", "able to change the roles");
*/

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
