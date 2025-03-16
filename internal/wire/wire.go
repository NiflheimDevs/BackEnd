//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/google/wire"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/niflheimdevs/backend/internal/bootstrap"
	"github.com/niflheimdevs/backend/internal/handlers"
	midauth "github.com/niflheimdevs/backend/internal/middlewares/authentication"
	panicwall "github.com/niflheimdevs/backend/internal/middlewares/exceptions"
	midratelimit "github.com/niflheimdevs/backend/internal/middlewares/ratelimit"
	"github.com/niflheimdevs/backend/internal/repositories"
	R "github.com/niflheimdevs/backend/internal/repositories/redis"

	"github.com/niflheimdevs/backend/internal/services"
	"github.com/redis/go-redis/v9"
)

var RepoProviderSet = wire.NewSet(
	wire.Struct(new(repositories.UserRepo), "*"),
	wire.Struct(new(R.UserCache), "*"),
)

var ServiceProviderSet = wire.NewSet(
	wire.Struct(new(services.UserService), "*"),
	services.NewJWT,
	ProvideConstants,
	// wire.Struct(new(services.JWT), "*"),
)

var HandlerProviderSet = wire.NewSet(
	wire.Struct(new(handlers.UserHandler), "*"),
	wire.Struct(new(handlers.ErrorHandler), "*"),
	handlers.NewValidator,
	// services.NewJWT,
)

var MiddlewareProviderSet = wire.NewSet(
	midratelimit.NewRateLimit,
	midauth.NewAuth,
	panicwall.NewPanicWall,
	wire.Struct(new(Middlewares), "*"),
)

func ProvideConstants(container *bootstrap.Di) *bootstrap.Constants {
	return container.Const
}

var ProviderSet = wire.NewSet(
	RepoProviderSet,
	ServiceProviderSet,
	HandlerProviderSet,
	MiddlewareProviderSet,
)

type Middlewares struct {
	Recovery       *panicwall.PanicWall
	RateLimit      *midratelimit.RateLimit
	Authentication *midauth.Authentication
}

type Application struct {
	UserHandler  *handlers.UserHandler
	ErrorHandler *handlers.ErrorHandler
	Middlewares  *Middlewares
}

func InitializeApplication(container *bootstrap.Di, db *pgxpool.Pool, myRedis *redis.Client) (*Application, error) {
	wire.Build(
		ProviderSet,
		wire.Struct(new(Application), "*"),
	)
	return &Application{}, nil
}
