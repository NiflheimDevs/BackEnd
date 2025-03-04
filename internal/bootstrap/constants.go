package bootstrap

import "time"

type Constants struct {
	Database DBConst
}

type DBConst struct {
	MaxOpenDbConn int
	MaxIdleDbConn time.Duration
	MaxDbLifeTime time.Duration
}

func NewConstant() *Constants {
	return &Constants{
		DBConst{
			MaxOpenDbConn: 10,
			MaxIdleDbConn: 5 * time.Minute,
			MaxDbLifeTime: 5 * time.Minute,
		},
	}
}
