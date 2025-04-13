package bootstrap

import "time"

type Constants struct {
	Database    DBConst
	JWTKeysPath string
	Context     Context
	StorageDir  string
	Project     Project
	Pagination  Pagination
	Port        string
	DevelopMode bool
}

type DBConst struct {
	MaxOpenDbConn int32
	MaxIdleDbConn time.Duration
	MaxDbLifeTime time.Duration
}

type Context struct {
	UserID string
}

type Project struct {
	LastTime time.Duration
}

type Pagination struct {
	Offset int
	Limit  int
}

func NewConstant() *Constants {
	return &Constants{
		Database: DBConst{
			MaxOpenDbConn: 10,
			MaxIdleDbConn: 5 * time.Minute,
			MaxDbLifeTime: 5 * time.Minute,
		},
		JWTKeysPath: "./internal/jwt",
		Context: Context{
			UserID: "userID",
		},

		StorageDir: "./storage/",
		Project: Project{
			LastTime: 7 * 24 * time.Hour,
		},
		Pagination: Pagination{
			Offset: 0,
			Limit:  10,
		},
		Port:        ":8080",
		DevelopMode: true,
	}
}
