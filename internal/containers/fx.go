package containers

import (
	"context"
	"github.com/018bf/companies/internal/configs"
	"github.com/018bf/companies/internal/interceptors"
	grpcInterface "github.com/018bf/companies/internal/interfaces/grpc"
	kafkaInterface "github.com/018bf/companies/internal/interfaces/kafka"
	postgresInterface "github.com/018bf/companies/internal/interfaces/postgres"
	restInterface "github.com/018bf/companies/internal/interfaces/rest"
	jwtRepositories "github.com/018bf/companies/internal/repositories/jwt"
	kafkaRepositories "github.com/018bf/companies/internal/repositories/kafka"
	postgresRepositories "github.com/018bf/companies/internal/repositories/postgres"
	"github.com/018bf/companies/internal/usecases"
	"github.com/018bf/companies/pkg/clock"
	"github.com/018bf/companies/pkg/log"
	"github.com/Shopify/sarama"
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
		grpcInterface.NewAuthMiddleware,
		restInterface.NewServer,
		restInterface.NewAuthMiddleware,
		interceptors.NewAuthInterceptor,
		usecases.NewAuthUseCase,
		jwtRepositories.NewAuthRepository,
		grpcInterface.NewCompanyServiceServer,
		restInterface.NewCompanyHandler,
		interceptors.NewCompanyInterceptor,
		usecases.NewCompanyUseCase,
		postgresRepositories.NewCompanyRepository,
		usecases.NewEventUseCase,
		kafkaRepositories.NewEventRepository,
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
