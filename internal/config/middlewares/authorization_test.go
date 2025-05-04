package middlewares

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"loki-backoffice/internal/config"
	"loki-backoffice/internal/config/logger"
	"loki-backoffice/pkg/jwt"
)

func Test_NewAuthorizationMiddleware(t *testing.T) {
	cfg := &config.Config{
		AppEnv:   "test",
		AppAddr:  "localhost:8080",
		LogLevel: "info",
	}
	log := logger.NewLogger(cfg)

	middleware := NewAuthorizationMiddleware(log)
	assert.NotNil(t, middleware)
}

func Test_AuthorizationMiddleware_Check(t *testing.T) {
	cfg := &config.Config{
		AppEnv:   "test",
		AppAddr:  "localhost:8080",
		LogLevel: "info",
	}
	log := logger.NewLogger(cfg)
	middleware := NewAuthorizationMiddleware(log)

	type request struct {
		claim      *jwt.Payload
		permission string
	}

	tests := []struct {
		name     string
		request  request
		expected int
	}{
		{
			name: "Success",
			request: request{
				claim: &jwt.Payload{
					ID:          "test-user",
					Permissions: []string{"read:users", "write:users"},
				},
				permission: "read:users",
			},
			expected: http.StatusOK,
		},
		{
			name: "User does not have permission",
			request: request{
				claim: &jwt.Payload{
					ID:          "test-user",
					Permissions: []string{"read:users"},
				},
				permission: "write:users",
			},
			expected: http.StatusForbidden,
		},
		{
			name: "No claims in context",
			request: request{
				claim:      nil,
				permission: "read:users",
			},
			expected: http.StatusUnauthorized,
		},
		{
			name: "Empty permissions",
			request: request{
				claim: &jwt.Payload{
					ID:          "test-user",
					Permissions: []string{},
				},
				permission: "read:users",
			},
			expected: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte("Success"))
			})

			req, err := http.NewRequest("GET", "/test", nil)
			assert.NoError(t, err)

			if tt.request.claim != nil {
				ctx := context.WithValue(req.Context(), Claim{}, tt.request.claim)
				req = req.WithContext(ctx)
			}

			rr := httptest.NewRecorder()

			middleware.Check(tt.request.permission)(handler).ServeHTTP(rr, req)
			assert.Equal(t, tt.expected, rr.Code)
		})
	}
}
