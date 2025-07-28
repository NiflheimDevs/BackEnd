//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/google/wire"
	"github.com/niflheimdevs/backend/bootstrap"
	"github.com/niflheimdevs/backend/internal/delivery/consumer"
	"github.com/niflheimdevs/backend/internal/delivery/handlers"
	midauth "github.com/niflheimdevs/backend/internal/delivery/middlewares/authentication"
	panicwall "github.com/niflheimdevs/backend/internal/delivery/middlewares/exceptions"
	midratelimit "github.com/niflheimdevs/backend/internal/delivery/middlewares/ratelimit"
	midupgrader "github.com/niflheimdevs/backend/internal/delivery/middlewares/upgrader"
	"github.com/niflheimdevs/backend/internal/infrastructure/db/driver"
	"github.com/niflheimdevs/backend/internal/infrastructure/db/seed"
	db "github.com/niflheimdevs/backend/internal/infrastructure/db/transaction"
	elasticimpl "github.com/niflheimdevs/backend/internal/infrastructure/repositories/elastic"
	repositoriesimpl "github.com/niflheimdevs/backend/internal/infrastructure/repositories/postgres"
	redisimpl "github.com/niflheimdevs/backend/internal/infrastructure/repositories/redis"
	storageimpl "github.com/niflheimdevs/backend/internal/infrastructure/repositories/storage"
	"github.com/niflheimdevs/backend/internal/infrastructure/websocket"
	"github.com/niflheimdevs/backend/pkg"

	"github.com/niflheimdevs/backend/internal/domain/repositories/elastic"
	repositories "github.com/niflheimdevs/backend/internal/domain/repositories/postgres"
	"github.com/niflheimdevs/backend/internal/domain/repositories/postgres/transaction"
	"github.com/niflheimdevs/backend/internal/domain/repositories/redis"
	"github.com/niflheimdevs/backend/internal/domain/repositories/storage"

	cdcservices "github.com/niflheimdevs/backend/internal/application/cdc/services"
	cdcservicesimpl "github.com/niflheimdevs/backend/internal/application/cdc/services/impl"
	"github.com/niflheimdevs/backend/internal/application/services"
	servicesimpl "github.com/niflheimdevs/backend/internal/application/services/impl"
)

var DatabaseProviderSet = wire.NewSet(
	driver.ConnectSQL,
	driver.ConncetRedis,
	driver.ConnectElastic,

	db.NewTxManager,
	wire.Bind(new(transaction.TxManager), new(*db.PgxTxManager)),
)

var ElasticAndDatabaseProviderSet = wire.NewSet(
	driver.ConnectElastic,
	driver.ConnectSQL,
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
	repositoriesimpl.NewTeamRepo,
	repositoriesimpl.NewRoleRepo,
	repositoriesimpl.NewBidRepo,
	repositoriesimpl.NewChatRepo,
	repositoriesimpl.NewUrlTokenRepo,
	elasticimpl.NewSearchElastic,
	repositoriesimpl.NewCommentRepo,
	storageimpl.NewS3Storage,
	redisimpl.NewUserCache,
	wire.Bind(new(repositories.UserRepo), new(*repositoriesimpl.UserRepo)),
	wire.Bind(new(repositories.TagRepo), new(*repositoriesimpl.TagRepo)),
	wire.Bind(new(repositories.CareerRepo), new(*repositoriesimpl.CareerRepo)),
	wire.Bind(new(repositories.LabelRepo), new(*repositoriesimpl.LabelRepo)),
	wire.Bind(new(repositories.ProjectRepo), new(*repositoriesimpl.ProjectRepo)),
	wire.Bind(new(repositories.PaymentRepo), new(*repositoriesimpl.PaymentRepo)),
	wire.Bind(new(repositories.BidRepo), new(*repositoriesimpl.BidRepo)),
	wire.Bind(new(repositories.TeamRepo), new(*repositoriesimpl.TeamRepo)),
	wire.Bind(new(repositories.RoleRepo), new(*repositoriesimpl.RoleRepo)),
	wire.Bind(new(repositories.ChatRepo), new(*repositoriesimpl.ChatRepo)),
	wire.Bind(new(repositories.UrlTokenRepo), new(*repositoriesimpl.UrlTokenRepo)),
	wire.Bind(new(elastic.SearchRepo), new(*elasticimpl.SearchRepo)),
	wire.Bind(new(repositories.CommentRepo), new(*repositoriesimpl.CommentRepo)),
	wire.Bind(new(storage.S3Storage), new(*storageimpl.S3Storage)),
	wire.Bind(new(redis.UserCache), new(*redisimpl.UserCache)),
)

var ElRepoProviderSet = wire.NewSet(
	repositoriesimpl.NewTagRepo,
	elasticimpl.NewProjectElastic,
	elasticimpl.NewUserElastic,
	elasticimpl.NewTeamElastic,

	wire.Bind(new(repositories.TagRepo), new(*repositoriesimpl.TagRepo)),
	wire.Bind(new(elastic.ProjectElastic), new(*elasticimpl.ProjectElastic)),
	wire.Bind(new(elastic.UserElastic), new(*elasticimpl.UserElastic)),
	wire.Bind(new(elastic.TeamElastic), new(*elasticimpl.TeamElastic)),
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
	servicesimpl.NewTeamService,
	servicesimpl.NewBidService,
	servicesimpl.NewSmsService,
	servicesimpl.NewJWT,
	servicesimpl.NewRoleService,
	servicesimpl.NewChatService,
	servicesimpl.NewSherlockService,
	servicesimpl.NewCommentService,
	servicesimpl.NewEmailService,
	servicesimpl.NewUrlTokenService,
	wire.Bind(new(services.UserService), new(*servicesimpl.UserService)),
	wire.Bind(new(services.TagService), new(*servicesimpl.TagService)),
	wire.Bind(new(services.CareerService), new(*servicesimpl.CareerService)),
	wire.Bind(new(services.LabelService), new(*servicesimpl.LabelService)),
	wire.Bind(new(services.ProjectService), new(*servicesimpl.ProjectService)),
	wire.Bind(new(services.PaymentService), new(*servicesimpl.PaymentService)),
	wire.Bind(new(services.TeamService), new(*servicesimpl.TeamService)),
	wire.Bind(new(services.BidService), new(*servicesimpl.BidService)),
	wire.Bind(new(services.SmsService), new(*servicesimpl.SmsService)),
	wire.Bind(new(services.JWT), new(*servicesimpl.JWT)),
	wire.Bind(new(services.ChatService), new(*servicesimpl.ChatService)),
	wire.Bind(new(services.RoleService), new(*servicesimpl.RoleService)),
	wire.Bind(new(services.SherlockService), new(*servicesimpl.SherlockService)),
	wire.Bind(new(services.CommentService), new(*servicesimpl.CommentService)),
	wire.Bind(new(services.EmailService), new(*servicesimpl.EmailService)),
	wire.Bind(new(services.UrlTokenService), new(*servicesimpl.UrlTokenService)),
	ProvideConstants,
	ProvideEnv,
	ProvideS3,
)

var CdcServiceProviderSet = wire.NewSet(
	cdcservicesimpl.NewProjectCdc,
	cdcservicesimpl.NewUserCdc,
	cdcservicesimpl.NewTeamCdc,
	wire.Bind(new(cdcservices.ProjectCdc), new(*cdcservicesimpl.ProjectCdc)),
	wire.Bind(new(cdcservices.UserCdc), new(*cdcservicesimpl.UserCdc)),
	wire.Bind(new(cdcservices.TeamCdc), new(*cdcservicesimpl.TeamCdc)),
	ProvideConstants,
	ProvideEnv,
)

var HandlerProviderSet = wire.NewSet(
	handlers.NewFileHandler,
	handlers.NewUserHandler,
	handlers.NewProjectHandler,
	handlers.NewGeneralHandler,
	handlers.NewPaymentHandler,
	handlers.NewBidHandler,
	handlers.NewTeamHandler,
	handlers.NewRoleHandler,
	handlers.NewChatHandler,
	handlers.NewSherlockHandler,
	handlers.NewCommentHandler,
	wire.Struct(new(Handlers), "*"),
)

var KafkaCdcProviderSet = wire.NewSet(
	consumer.NewKafkaCdc,
)

var MiddlewareProviderSet = wire.NewSet(
	midratelimit.NewRateLimit,
	midauth.NewAuth,
	panicwall.NewPanicWall,
	midupgrader.NewWebSocketUpgrader,
	wire.Struct(new(Middlewares), "*"),
)

func ProvideConstants(container *bootstrap.Di) *bootstrap.Constants {
	return container.Const
}

func ProvideEnv(container *bootstrap.Di) *bootstrap.Env {
	return container.Env
}

func ProvideS3(container *bootstrap.Di) *bootstrap.S3 {
	return &container.Env.Storage
}

var KafkaElasticProviderSet = wire.NewSet(
	ElasticAndDatabaseProviderSet,
	ElRepoProviderSet,
	CdcServiceProviderSet,
	KafkaCdcProviderSet,
)

var ProviderSet = wire.NewSet(
	DatabaseProviderSet,
	PkgProviderSet,
	RepoProviderSet,
	FileServiceProviderSet,
	ServiceProviderSet,
	HandlerProviderSet,
	MiddlewareProviderSet,
	seed.NewSeeder,
)

type Middlewares struct {
	Recovery       *panicwall.PanicWall
	RateLimit      *midratelimit.RateLimit
	Authentication *midauth.Authentication
	Upgrader       *midupgrader.WebSocketUpgrader
}

type Handlers struct {
	FileHandler     *handlers.FileHandler
	UserHandler     *handlers.UserHandler
	ProjectHandler  *handlers.ProjectHandler
	GeneralHandler  *handlers.GeneralHandler
	PaymentHandler  *handlers.PaymentHandler
	BidHandler      *handlers.BidHandler
	TeamHandler     *handlers.TeamHandler
	RoleHandler     *handlers.RoleHandler
	ChatHandler     *handlers.ChatHandler
	SherlockHandler *handlers.SherlockHandler
	CommentHandler  *handlers.CommentHandler
}

type Application struct {
	Handlers    *Handlers
	Middlewares *Middlewares
	Seeder      *seed.Seeder
}

type ElasticApp struct {
	KafkaCdc *consumer.KafkaCdc
}

func InitializeApplication(container *bootstrap.Di, hub *websocket.Hub) (*Application, error) {
	wire.Build(
		ProviderSet,
		wire.Struct(new(Application), "*"),
	)
	return &Application{}, nil
}

func InitializeCdcConsumer(container *bootstrap.Di) (*ElasticApp, error) {
	wire.Build(
		KafkaElasticProviderSet,
		wire.Struct(new(ElasticApp), "*"),
	)
	return &ElasticApp{}, nil
}
