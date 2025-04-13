//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/google/wire"
	"github.com/niflheimdevs/backend/bootstrap"
	"github.com/niflheimdevs/backend/internal/delivery/handlers"
	midauth "github.com/niflheimdevs/backend/internal/delivery/middlewares/authentication"
	panicwall "github.com/niflheimdevs/backend/internal/delivery/middlewares/exceptions"
	midratelimit "github.com/niflheimdevs/backend/internal/delivery/middlewares/ratelimit"
	"github.com/niflheimdevs/backend/internal/infrastructure/db/driver"
	db "github.com/niflheimdevs/backend/internal/infrastructure/db/transaction"
	repositoriesimpl "github.com/niflheimdevs/backend/internal/infrastructure/repositories/postgres"
	redisimpl "github.com/niflheimdevs/backend/internal/infrastructure/repositories/redis"
	storageimpl "github.com/niflheimdevs/backend/internal/infrastructure/repositories/storage"
	"github.com/niflheimdevs/backend/pkg"

	repositries "github.com/niflheimdevs/backend/internal/domain/repositories/postgres"
	"github.com/niflheimdevs/backend/internal/domain/repositories/postgres/transaction"
	"github.com/niflheimdevs/backend/internal/domain/repositories/redis"
	"github.com/niflheimdevs/backend/internal/domain/repositories/storage"

	"github.com/niflheimdevs/backend/internal/application/services"
	servicesimpl "github.com/niflheimdevs/backend/internal/application/services/impl"
)

var DatabaseProviderSet = wire.NewSet(
	driver.ConnectSQL,
	driver.ConncetRedis,

	db.NewTxManager,
	wire.Bind(new(transaction.TxManager), new(*db.PgxTxManager)),
)

var PkgProviderSet = wire.NewSet(
	pkg.NewValidator,
	pkg.NewSecretSauce,
)

var RepoProviderSet = wire.NewSet(
	repositoriesimpl.NewUserRepo,
	repositoriesimpl.NewProjectRepo,
	repositoriesimpl.NewCareerRepo,
	repositoriesimpl.NewTagRepo,
	repositoriesimpl.NewLabelRepo,
	repositoriesimpl.NewPaymentRepo,
	storageimpl.NewFileStorage,
	redisimpl.NewUserCache,
	wire.Bind(new(repositries.UserRepo), new(*repositoriesimpl.UserRepo)),
	wire.Bind(new(repositries.TagRepo), new(*repositoriesimpl.TagRepo)),
	wire.Bind(new(repositries.CareerRepo), new(*repositoriesimpl.CareerRepo)),
	wire.Bind(new(repositries.LabelRepo), new(*repositoriesimpl.LabelRepo)),
	wire.Bind(new(repositries.ProjectRepo), new(*repositoriesimpl.ProjectRepo)),
	wire.Bind(new(repositries.PaymentRepo), new(*repositoriesimpl.PaymentRepo)),
	wire.Bind(new(storage.FileStorage), new(*storageimpl.FileStorage)),
	wire.Bind(new(redis.UserCache), new(*redisimpl.UserCache)),
)

var FileServiceProviderSet = wire.NewSet(
	servicesimpl.NewFileService,
	wire.Bind(new(services.FileService), new(*servicesimpl.FileService)),
)

var ServiceProviderSet = wire.NewSet(
	servicesimpl.NewUserService,
	servicesimpl.NewTagService,
	servicesimpl.NewCareerService,
	servicesimpl.NewLabelService,
	servicesimpl.NewProjectService,
	servicesimpl.NewPaymentService,
	servicesimpl.NewSmsService,
	servicesimpl.NewJWT,

	wire.Bind(new(services.UserService), new(*servicesimpl.UserService)),
	wire.Bind(new(services.TagService), new(*servicesimpl.TagService)),
	wire.Bind(new(services.CareerService), new(*servicesimpl.CareerService)),
	wire.Bind(new(services.LabelService), new(*servicesimpl.LabelService)),
	wire.Bind(new(services.ProjectService), new(*servicesimpl.ProjectService)),
	wire.Bind(new(services.PaymentService), new(*servicesimpl.PaymentService)),
	wire.Bind(new(services.SmsService), new(*servicesimpl.SmsService)),
	wire.Bind(new(services.JWT), new(*servicesimpl.JWT)),

	ProvideConstants,
	ProvideEnv,
)

var HandlerProviderSet = wire.NewSet(
	handlers.NewFileHandler,
	handlers.NewUserHandler,
	handlers.NewErrorHandler,
	handlers.NewProjectHandler,
	handlers.NewGeneralHandler,
	handlers.NewPaymentHandler,
	wire.Struct(new(Handlers), "*"),
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

func ProvideEnv(container *bootstrap.Di) *bootstrap.Env {
	return container.Env
}

var ProviderSet = wire.NewSet(
	DatabaseProviderSet,
	PkgProviderSet,
	RepoProviderSet,
	FileServiceProviderSet,
	ServiceProviderSet,
	HandlerProviderSet,
	MiddlewareProviderSet,
)

type Middlewares struct {
	Recovery       *panicwall.PanicWall
	RateLimit      *midratelimit.RateLimit
	Authentication *midauth.Authentication
}

type Handlers struct {
	FileHandler    *handlers.FileHandler
	UserHandler    *handlers.UserHandler
	ErrorHandler   *handlers.ErrorHandler
	ProjectHandler *handlers.ProjectHandler
	GeneralHandler *handlers.GeneralHandler
	PaymentHandler *handlers.PaymentHandler
}

type Application struct {
	Handlers    *Handlers
	Middlewares *Middlewares
}

func InitializeApplication(container *bootstrap.Di) (*Application, error) {
	wire.Build(
		ProviderSet,
		wire.Struct(new(Application), "*"),
	)
	return &Application{}, nil
}
