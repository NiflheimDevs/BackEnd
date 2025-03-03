//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/google/wire"
	"github.com/niflheimdevs/backend/internal/handlers"
	"github.com/niflheimdevs/backend/internal/services"
)

var DatabaseProviderSet = wire.NewSet()

var ServiceProviderSet = wire.NewSet(
	services.NewUserService,
)

var HandlerProviderSet = wire.NewSet(
	handlers.NewUserHandler,
)

var ProviderSet = wire.NewSet(
	DatabaseProviderSet,
	ServiceProviderSet,
	HandlerProviderSet,
)

type Application struct {
}

func InitializeApplication() (*Application, error) {
	wire.Build(
		ProviderSet,
		wire.Struct(new(Application), "*"),
	)
	return &Application{}, nil
}
