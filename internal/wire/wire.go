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
