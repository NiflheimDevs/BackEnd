//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/google/wire"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/niflheimdevs/backend/internal/bootstrap"
	"github.com/niflheimdevs/backend/internal/handlers"
	midauth "github.com/niflheimdevs/backend/internal/middlewares/authentication"
	midrecovery "github.com/niflheimdevs/backend/internal/middlewares/exceptions"
	midratelimit "github.com/niflheimdevs/backend/internal/middlewares/ratelimit"
	"github.com/niflheimdevs/backend/internal/repositories"
	"github.com/niflheimdevs/backend/internal/services"
)

var RepoProviderSet = wire.NewSet(
	// repositories.NewUserRepo,
	wire.Struct(new(repositories.UserRepo), "*"),
)

var ServiceProviderSet = wire.NewSet(
	wire.Struct(new(services.UserService), "*"),
	wire.Struct(new(services.JWTToken), "*"),
)

var HandlerProviderSet = wire.NewSet(
	// handlerss.NewUserHandler,
	wire.Struct(new(handlers.UserHandler), "*"),
)

var MiddlewareProviderSet = wire.NewSet(
	midratelimit.NewRateLimit,
	midauth.NewAuth,
	midrecovery.NewRecoveryMiddleware,
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
	Recovery       *midrecovery.RecoveryMiddleware
	RateLimit      *midratelimit.RateLimit
	Authentication *midauth.Authentication
}

type Application struct {
	UserHandler *handlers.UserHandler
	Middlewares *Middlewares
}

func InitializeApplication(container *bootstrap.Di, db *pgxpool.Pool) (*Application, error) {
	wire.Build(
		ProvideConstants,
		ProviderSet,
		wire.Struct(new(Application), "*"),
	)
	return &Application{}, nil
}
