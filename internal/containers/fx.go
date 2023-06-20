package containers

import (
	"context"
	authInterceptor "github.com/018bf/companies/internal/auth/interceptor"
	authRepository "github.com/018bf/companies/internal/auth/repository/jwt"
	authService "github.com/018bf/companies/internal/auth/service"
	companyGrpc "github.com/018bf/companies/internal/company/grpc"
	companyInterceptor "github.com/018bf/companies/internal/company/interceptor"
	companyRepository "github.com/018bf/companies/internal/company/repository/postgres"
	companyService "github.com/018bf/companies/internal/company/service"
	eventService "github.com/018bf/companies/internal/event/service"
	companiespb "github.com/018bf/companies/pkg/companiespb/v1"
	"github.com/Shopify/sarama"

	"github.com/018bf/companies/internal/configs"
	eventRepository "github.com/018bf/companies/internal/event/repositories/kafka"
	grpcInterface "github.com/018bf/companies/internal/interfaces/grpc"
	kafkaInterface "github.com/018bf/companies/internal/interfaces/kafka"
	postgresInterface "github.com/018bf/companies/internal/interfaces/postgres"
	restInterface "github.com/018bf/companies/internal/interfaces/rest"
	"github.com/018bf/companies/pkg/clock"
	"github.com/018bf/companies/pkg/log"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

var FXModule = fx.Options(
	fx.WithLogger(func(logger log.Logger) fxevent.Logger {
		return logger
	}),
	fx.Provide(
		func(config *configs.Config) (log.Logger, error) {
			return log.NewLog(config.LogLevel)
		},
		context.Background,
		configs.ParseConfig,
		clock.NewRealClock,
		postgresInterface.NewDatabase,
		postgresInterface.NewMigrateManager,
		kafkaInterface.NewProducer,
		grpcInterface.NewServer,
		grpcInterface.NewRequestIDMiddleware,
		func(authInterceptor *authInterceptor.AuthInterceptor, logger log.Logger, config *configs.Config) *grpcInterface.AuthMiddleware {
			return grpcInterface.NewAuthMiddleware(authInterceptor, logger, config)
		},
		restInterface.NewServer,
		restInterface.NewAuthMiddleware,

		authRepository.NewAuthRepository,
		func(authRepository *authRepository.AuthRepository, logger log.Logger) *authService.AuthService {
			return authService.NewAuthService(authRepository, logger)
		},
		func(
			authService *authService.AuthService,
			clock clock.Clock,
			logger log.Logger,
		) *authInterceptor.AuthInterceptor {
			return authInterceptor.NewAuthInterceptor(authService, clock, logger)
		},

		eventRepository.NewEventRepository,
		func(eventRepository *eventRepository.EventRepository, logger log.Logger) *eventService.EventService {
			return eventService.NewEventService(eventRepository, logger)
		},

		companyRepository.NewCompanyRepository,
		func(
			companyRepository *companyRepository.CompanyRepository,
			clock clock.Clock,
			logger log.Logger,
		) *companyService.CompanyService {
			return companyService.NewCompanyService(companyRepository, clock, logger)
		},
		func(
			companyService *companyService.CompanyService,
			authService *authService.AuthService,
			eventService *eventService.EventService,
			clock clock.Clock,
			logger log.Logger,
		) *companyInterceptor.CompanyInterceptor {
			return companyInterceptor.NewCompanyInterceptor(companyService, authService, eventService, logger)
		},
		fx.Annotate(
			func(companyInterceptor *companyInterceptor.CompanyInterceptor, logger log.Logger) *companyGrpc.CompanyServiceServer {
				return companyGrpc.NewCompanyServiceServer(companyInterceptor, logger)
			},
			fx.As(new(companiespb.CompanyServiceServer)),
		),
	),
	fx.Invoke(func(
		lifecycle fx.Lifecycle,
		logger log.Logger,
		producer sarama.SyncProducer,
		shutdowner fx.Shutdowner,
	) {
		lifecycle.Append(fx.Hook{
			OnStop: func(_ context.Context) error {
				return producer.Close()
			},
		})
	}),
)

func NewMigrateContainer(config string) *fx.App {
	app := fx.New(
		fx.Provide(func() string {
			return config
		}),
		FXModule,
		fx.Invoke(func(
			lifecycle fx.Lifecycle,
			logger log.Logger,
			manager *postgresInterface.MigrateManager,
			shutdowner fx.Shutdowner,
		) {
			lifecycle.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					err := manager.Up(ctx)
					if err != nil {
						logger.Error("shutdown", log.Any("error", err))
					}
					return shutdowner.Shutdown(fx.ExitCode(0))
				},
				OnStop: nil,
			})
		}),
	)
	return app
}
func NewGRPCContainer(config string) *fx.App {
	app := fx.New(
		fx.Provide(func() string {
			return config
		}),
		FXModule,
		fx.Invoke(func(
			lifecycle fx.Lifecycle,
			logger log.Logger,
			server *grpcInterface.Server,
			shutdowner fx.Shutdowner,
		) {
			lifecycle.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					go func() {
						err := server.Start(ctx)
						if err != nil {
							logger.Error("shutdown", log.Any("error", err))
							_ = shutdowner.Shutdown()
						}
					}()
					return nil
				},
				OnStop: server.Stop,
			})
		}),
	)
	return app
}

func NewRESTContainer(config string) *fx.App {
	app := fx.New(
		fx.Provide(func() string {
			return config
		}),
		FXModule,
		fx.Invoke(func(
			lifecycle fx.Lifecycle,
			logger log.Logger,
			server *restInterface.Server,
			shutdowner fx.Shutdowner,
		) {
			lifecycle.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					go func() {
						err := server.Start(ctx)
						if err != nil {
							logger.Error("shutdown", log.Any("error", err))
							_ = shutdowner.Shutdown()
						}
					}()
					return nil
				},
				OnStop: server.Stop,
			})
		}),
	)
	return app
}
