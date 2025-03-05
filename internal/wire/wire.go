//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/google/wire"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/niflheimdevs/backend/internal/handlers"
	"github.com/niflheimdevs/backend/internal/repositories"
	"github.com/niflheimdevs/backend/internal/services"
)

var RepoProviderSet = wire.NewSet(
	// repositories.NewUserRepo,
	wire.Struct(new(repositories.UserRepo), "*"),
)

var ServiceProviderSet = wire.NewSet(
	// services.NewUserService,
	wire.Struct(new(services.UserService), "*"),
)

var HandlerProviderSet = wire.NewSet(
	// handlerss.NewUserHandler,
	wire.Struct(new(handlers.UserHandler), "*"),
)

var ProviderSet = wire.NewSet(
	RepoProviderSet,
	ServiceProviderSet,
	HandlerProviderSet,
)

type Application struct {
	UserHandler *handlers.UserHandler
}

func InitializeApplication(db *pgxpool.Pool) (*Application, error) {
	wire.Build(
		ProviderSet,
		wire.Struct(new(Application), "*"),
	)
	return &Application{}, nil
}
