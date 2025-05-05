package bootstrap

import "time"

type Constants struct {
	Database      DBConst
	JWTKeysPath   string
	SSLKeysPath   string
	Context       Context
	Project       Project
	Pagination    Pagination
	Port          string
	DevelopMode   bool
	MaxPhotoSize  int64
	MaxResumeSize int64
	RateLimiter   RateLimiter
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

type RateLimiter struct {
	Limit float64
	Burst int
}

func NewConstant() *Constants {
	return &Constants{
		Database: DBConst{
			MaxOpenDbConn: 10,
			MaxIdleDbConn: 5 * time.Minute,
			MaxDbLifeTime: 5 * time.Minute,
		},
		JWTKeysPath: "./internal/jwt",
		SSLKeysPath: "./SSL",
		Context: Context{
			UserID: "userID",
		},
		Project: Project{
			LastTime: 7 * 24 * time.Hour,
		},
		Pagination: Pagination{
			Offset: 0,
			Limit:  10,
		},
		MaxPhotoSize:  5000000,
		MaxResumeSize: 10000000,
		Port:          ":8080",
		DevelopMode:   true,
		RateLimiter: RateLimiter{
			Limit: (6.0 / 60),
			Burst: 20,
		},
	}
}
