package bootstrap

import "time"

type Constants struct {
	Database         DBConst
	JWTKeysPath      string
	SSLKeysPath      string
	Context          Context
	Project          Project
	Pagination       Pagination
	Port             string
	DevelopMode      bool
	MaxPhotoSize     int64
	MaxResumeSize    int64
	RateLimiter      RateLimiter
	WebsocketSetting WebsocketSetting
	Kafka            Kafka
	UrlTokenSetting  UrlTokenSetting
}

type Kafka struct {
	Cdctopics         []string
	GroupidForElastic string
}

type DBConst struct {
	MaxOpenDbConn int32
	MaxIdleDbConn time.Duration
	MaxDbLifeTime time.Duration
}

type Context struct {
	UserID              string
	WebSocketConnection string
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

type WebsocketSetting struct {
	WriteTimeout      time.Duration
	ReadTimeout       time.Duration
	PingPeriod        time.Duration
	MaxMessageSize    int
	MessageBufferSize int
}

type UrlTokenSetting struct {
	Length           int
	ExpiryTeamInvite time.Duration
	ExpirtyEmailVer  time.Duration
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
			UserID:              "userID",
			WebSocketConnection: "wsConnection",
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
		Kafka: Kafka{
			GroupidForElastic: "elastic-readers",
			Cdctopics:         []string{"postgres.public.users", "postgres.public.users_career_tag", "postgres.public.project", "postgres.public.project_tag", "postgres.public.team"},
		},
		WebsocketSetting: WebsocketSetting{
			WriteTimeout:      10 * time.Second,
			ReadTimeout:       60 * time.Second,
			PingPeriod:        54 * time.Second,
			MaxMessageSize:    524288,
			MessageBufferSize: 256,
		},
		UrlTokenSetting: UrlTokenSetting{
			Length:           32,
			ExpiryTeamInvite: 24 * 7 * time.Hour,
			ExpirtyEmailVer:  24 * 2 * time.Hour,
		},
	}
}
