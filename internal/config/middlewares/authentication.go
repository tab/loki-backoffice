package middlewares

import (
	"encoding/json"
	"net/http"
	"strings"

	"loki-backoffice/internal/app/serializers"
	"loki-backoffice/pkg/jwt"
	"loki-backoffice/pkg/logger"
)

type AuthenticationMiddleware interface {
	Authenticate(next http.Handler) http.Handler
}

type authenticationMiddleware struct {
	jwt jwt.Jwt
	log *logger.Logger
}

func NewAuthenticationMiddleware(jwt jwt.Jwt, log *logger.Logger) AuthenticationMiddleware {
	return &authenticationMiddleware{
		jwt: jwt,
		log: log,
	}
}

func (m *authenticationMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, ok := extractBearerToken(r)
		if !ok {
			m.log.Error().Msg("Invalid authorization header")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		claims, err := m.jwt.Decode(token)
		if err != nil {
			m.log.Error().Err(err).Msg("Failed to decode token")
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(serializers.ErrorSerializer{Error: err.Error()})
			return
		}

		ctx := NewContextModifier(r.Context()).
			WithClaim(claims).
			WithToken(token).
			Context()

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func extractBearerToken(r *http.Request) (string, bool) {
	authHeader := r.Header.Get(Authorization)
	if authHeader == "" {
		return "", false
	}

	if len(authHeader) < len(bearerScheme) || !strings.EqualFold(authHeader[:len(bearerScheme)], bearerScheme) {
		return "", false
	}

	token := authHeader[len(bearerScheme):]
	if token == "" {
		return "", false
	}

	return token, true
}
