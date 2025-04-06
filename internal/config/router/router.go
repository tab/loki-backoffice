package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"loki-backoffice/internal/app/controllers"
	"loki-backoffice/internal/config"
	"loki-backoffice/internal/config/middlewares"
	"loki-backoffice/pkg/rbac"
)

func NewRouter(
	cfg *config.Config,

	authentication middlewares.AuthenticationMiddleware,
	authorization middlewares.AuthorizationMiddleware,
	telemetry middlewares.TelemetryMiddleware,
	logger middlewares.LoggerMiddleware,

	health controllers.HealthController,

	permissions controllers.PermissionsController,
	roles controllers.RolesController,
	scopes controllers.ScopesController,
	tokens controllers.TokensController,
	users controllers.UsersController,
) http.Handler {
	r := chi.NewRouter()

	r.Use(telemetry.Trace)
	r.Use(middleware.RequestID)
	r.Use(logger.Log)
	r.Use(middleware.Compress(5))
	r.Use(middleware.Heartbeat("/health"))
	r.Use(
		cors.Handler(cors.Options{
			AllowedOrigins: []string{"http://*", cfg.ClientURL},
			AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
			AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-Request-ID", "X-Trace-ID"},
			MaxAge:         300,
		}),
	)

	r.Get("/live", health.HandleLiveness)
	r.Get("/ready", health.HandleReadiness)

	r.Group(func(r chi.Router) {
		r.Use(authentication.Authenticate)
		r.Route("/api/backoffice", func(r chi.Router) {
			r.With(authorization.Check(rbac.ReadPermissions)).Get("/permissions", permissions.List)
			r.With(authorization.Check(rbac.ReadPermissions)).Get("/permissions/{id}", permissions.Get)
			r.With(authorization.Check(rbac.WritePermissions)).Post("/permissions", permissions.Create)
			r.With(authorization.Check(rbac.WritePermissions)).Put("/permissions/{id}", permissions.Update)
			r.With(authorization.Check(rbac.WritePermissions)).Delete("/permissions/{id}", permissions.Delete)

			r.With(authorization.Check(rbac.ReadRoles)).Get("/roles", roles.List)
			r.With(authorization.Check(rbac.ReadRoles)).Get("/roles/{id}", roles.Get)
			r.With(authorization.Check(rbac.WriteRoles)).Post("/roles", roles.Create)
			r.With(authorization.Check(rbac.WriteRoles)).Put("/roles/{id}", roles.Update)
			r.With(authorization.Check(rbac.WriteRoles)).Delete("/roles/{id}", roles.Delete)

			r.With(authorization.Check(rbac.ReadScopes)).Get("/scopes", scopes.List)
			r.With(authorization.Check(rbac.ReadScopes)).Get("/scopes/{id}", scopes.Get)
			r.With(authorization.Check(rbac.WriteScopes)).Post("/scopes", scopes.Create)
			r.With(authorization.Check(rbac.WriteScopes)).Put("/scopes/{id}", scopes.Update)
			r.With(authorization.Check(rbac.WriteScopes)).Delete("/scopes/{id}", scopes.Delete)

			r.With(authorization.Check(rbac.ReadTokens)).Get("/tokens", tokens.List)
			r.With(authorization.Check(rbac.WriteTokens)).Delete("/tokens/{id}", tokens.Delete)

			r.With(authorization.Check(rbac.ReadUsers)).Get("/users", users.List)
			r.With(authorization.Check(rbac.ReadUsers)).Get("/users/{id}", users.Get)
			r.With(authorization.Check(rbac.WriteUsers)).Post("/users", users.Create)
			r.With(authorization.Check(rbac.WriteUsers)).Put("/users/{id}", users.Update)
			r.With(authorization.Check(rbac.WriteUsers)).Delete("/users/{id}", users.Delete)
		})
	})

	return r
}
