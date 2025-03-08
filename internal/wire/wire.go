//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/google/wire"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/niflheimdevs/backend/internal/handlers"
	panicwall "github.com/niflheimdevs/backend/internal/middlewares/exceptions"
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
)

var HandlerProviderSet = wire.NewSet(
	wire.Struct(new(handlers.UserHandler), "*"),
	handlers.NewValidator,
)

var ProviderSet = wire.NewSet(
	RepoProviderSet,
	ServiceProviderSet,
	HandlerProviderSet,
	wire.Struct(new(panicwall.PanicWall), "*"),
)

type Application struct {
	UserHandler *handlers.UserHandler
	Recovery    *panicwall.PanicWall
}

func InitializeApplication(db *pgxpool.Pool, myRedis *redis.Client) (*Application, error) {
	wire.Build(
		ProviderSet,
		wire.Struct(new(Application), "*"),
	)
	return &Application{}, nil
}
