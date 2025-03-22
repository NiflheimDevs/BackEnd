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
	"github.com/niflheimdevs/backend/internal/repositories/storage"
	"github.com/niflheimdevs/backend/internal/utils"

	"github.com/niflheimdevs/backend/internal/services"
	"github.com/redis/go-redis/v9"
)

var RepoProviderSet = wire.NewSet(
	wire.Struct(new(repositories.UserRepo), "*"),
	wire.Struct(new(repositories.ProjectRepo), "*"),
	wire.Struct(new(repositories.GeneralRepo), "*"),
	wire.Struct(new(R.UserCache), "*"),
	wire.Struct(new(storage.FileStorage), "*"),
)

var ServiceProviderSet = wire.NewSet(
	wire.Struct(new(services.UserService), "*"),
	wire.Struct(new(services.FileService), "*"),
	wire.Struct(new(services.ProjectService), "*"),
	wire.Struct(new(services.GeneralService), "*"),
	services.NewJWT,
	ProvideConstants,
)

var HandlerProviderSet = wire.NewSet(
	wire.Struct(new(handlers.UserHandler), "*"),
	wire.Struct(new(handlers.FileHandler), "*"),
	wire.Struct(new(handlers.ErrorHandler), "*"),
	wire.Struct(new(handlers.ProjectHandler), "*"),
	wire.Struct(new(handlers.GeneralHandler), "*"),
	handlers.NewValidator,
	utils.NewUtils,
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
	FileHandler    *handlers.FileHandler
	UserHandler    *handlers.UserHandler
	ErrorHandler   *handlers.ErrorHandler
	ProjectHandler *handlers.ProjectHandler
	GeneralHandler *handlers.GeneralHandler
	Middlewares    *Middlewares
}

func InitializeApplication(container *bootstrap.Di, db *pgxpool.Pool, myRedis *redis.Client) (*Application, error) {
	wire.Build(
		ProviderSet,
		wire.Struct(new(Application), "*"),
	)
	return &Application{}, nil
}
