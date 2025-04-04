package app

import (
	"context"
	"net/http"
	"time"

	"go.uber.org/fx"

	"loki-backoffice/internal/app/controllers"
	"loki-backoffice/internal/app/repositories"
	"loki-backoffice/internal/app/rpcs"
	"loki-backoffice/internal/app/services"
	"loki-backoffice/internal/config"
	"loki-backoffice/internal/config/middlewares"
	"loki-backoffice/internal/config/router"
	"loki-backoffice/internal/config/server"
	"loki-backoffice/internal/config/telemetry"
	"loki-backoffice/pkg/jwt"
	"loki-backoffice/pkg/logger"
)

var Module = fx.Options(
	logger.Module,
	jwt.Module,

	controllers.Module,
	repositories.Module,
	rpcs.Module,
	services.Module,

	middlewares.Module,

	server.Module,
	router.Module,
	telemetry.Module,

	fx.Invoke(registerHooks),
	fx.Invoke(registerGrpcClient),
	fx.Invoke(registerTelemetry),
)

func registerHooks(
	lifecycle fx.Lifecycle,
	cfg *config.Config,
	server server.Server,
	log *logger.Logger,
) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Info().Msgf("Starting server in %s environment at %s", cfg.AppEnv, cfg.AppAddr)

			go func() {
				if err := server.Run(); err != nil && err != http.ErrServerClosed {
					log.Error().Err(err).Msg("Server failed")
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info().Msg("Shutting down server...")

			shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()

			return server.Shutdown(shutdownCtx)
		},
	})
}

func registerGrpcClient(
	lifecycle fx.Lifecycle,
	client rpcs.Client,
	log *logger.Logger,
) {
	lifecycle.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			log.Info().Msg("Closing gRPC client connection...")
			return client.Close()
		},
	})
}

func registerTelemetry(lifecycle fx.Lifecycle, cfg *config.Config) {
	var ctx, cancel = context.WithCancel(context.Background())
	service, _ := telemetry.NewTelemetry(ctx, cfg)

	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return nil
		},
		OnStop: func(ctx context.Context) error {
			cancel()
			return service.Shutdown(ctx)
		},
	})
}
