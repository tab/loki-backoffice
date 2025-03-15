package app

import (
	"context"
	"net/http"
	"time"

	"go.uber.org/fx"

	"loki/internal/app/controllers"
	"loki/internal/app/repositories"
	"loki/internal/app/services"
	"loki/internal/config"
	"loki/internal/config/middlewares"
	"loki/internal/config/router"
	"loki/internal/config/server"
	"loki/internal/config/telemetry"
	"loki/pkg/logger"
)

var Module = fx.Options(
	logger.Module,

	controllers.Module,
	repositories.Module,
	services.Module,

	middlewares.Module,

	server.Module,
	router.Module,
	telemetry.Module,

	fx.Invoke(registerHooks),
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
