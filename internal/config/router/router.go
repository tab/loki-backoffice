package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"loki-backoffice/internal/app/controllers"
	"loki-backoffice/internal/config"
	"loki-backoffice/internal/config/middlewares"
)

func NewRouter(
	cfg *config.Config,
	telemetry middlewares.TelemetryMiddleware,
	health controllers.HealthController,
) http.Handler {
	r := chi.NewRouter()

	r.Use(telemetry.Trace)
	r.Use(middleware.Logger)
	r.Use(middleware.Compress(5))
	r.Use(middleware.Heartbeat("/health"))
	r.Use(
		cors.Handler(cors.Options{
			AllowedOrigins: []string{"http://*", cfg.ClientURL},
			AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
			AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-Trace-ID"},
			MaxAge:         300,
		}),
	)

	r.Get("/live", health.HandleLiveness)
	r.Get("/ready", health.HandleReadiness)

	return r
}
