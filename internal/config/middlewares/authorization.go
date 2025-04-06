package middlewares

import (
	"encoding/json"
	"net/http"

	"loki-backoffice/internal/app/errors"
	"loki-backoffice/internal/app/serializers"
	"loki-backoffice/internal/config/logger"
	"loki-backoffice/pkg/rbac"
)

type AuthorizationMiddleware interface {
	Check(permission string) func(http.Handler) http.Handler
}

type authorizationMiddleware struct {
	log *logger.Logger
}

func NewAuthorizationMiddleware(log *logger.Logger) AuthorizationMiddleware {
	return &authorizationMiddleware{
		log: log,
	}
}

func (m *authorizationMiddleware) Check(permission string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claim, ok := CurrentClaimFromContext(r.Context())
			if !ok {
				m.log.Error().Msg("No claims found in context")
				w.WriteHeader(http.StatusUnauthorized)
				_ = json.NewEncoder(w).Encode(serializers.ErrorSerializer{Error: errors.ErrUnauthorized.Error()})
				return
			}

			if !rbac.HasPermission(claim.Permissions, permission) {
				m.log.Warn().Msgf("User %s does not have required permission: %s", claim.ID, permission)
				w.WriteHeader(http.StatusForbidden)
				_ = json.NewEncoder(w).Encode(serializers.ErrorSerializer{Error: errors.ErrForbidden.Error()})
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
