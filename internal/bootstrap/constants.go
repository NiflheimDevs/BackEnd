package bootstrap

import "time"

type Constants struct {
	Database    DBConst
	JWTKeysPath string
	Context     Context
}

type DBConst struct {
	MaxOpenDbConn int32
	MaxIdleDbConn time.Duration
	MaxDbLifeTime time.Duration
}

type Context struct {
	UserID string
}

func NewConstant() *Constants {
	return &Constants{
		Database: DBConst{
			MaxOpenDbConn: 10,
			MaxIdleDbConn: 5 * time.Minute,
			MaxDbLifeTime: 5 * time.Minute,
		},
		JWTKeysPath: "../internal/jwt",
		Context: Context{
			UserID: "userID",
		},
	}
}
