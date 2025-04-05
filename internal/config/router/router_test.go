package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"loki-backoffice/internal/app/controllers"
	"loki-backoffice/internal/config"
	"loki-backoffice/internal/config/middlewares"
)

func Test_HealthCheck(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		AppEnv:  "test",
		AppAddr: "localhost:8080",
	}

	mockAuthenticationMiddleware := middlewares.NewMockAuthenticationMiddleware(ctrl)
	mockAuthorizationMiddleware := middlewares.NewMockAuthorizationMiddleware(ctrl)
	mockTelemetryMiddleware := middlewares.NewMockTelemetryMiddleware(ctrl)

	mockHealthController := controllers.NewMockHealthController(ctrl)
	mockPermissionsController := controllers.NewMockPermissionsController(ctrl)
	mockRolesController := controllers.NewMockRolesController(ctrl)
	mockScopesController := controllers.NewMockScopesController(ctrl)
	mockTokensController := controllers.NewMockTokensController(ctrl)
	mockUsersController := controllers.NewMockUsersController(ctrl)

	mockAuthenticationMiddleware.EXPECT().
		Authenticate(gomock.Any()).
		AnyTimes().
		DoAndReturn(func(next http.Handler) http.Handler {
			return next
		})
	mockAuthorizationMiddleware.EXPECT().
		Check(gomock.Any()).
		AnyTimes().
		DoAndReturn(func(permission string) func(http.Handler) http.Handler {
			return func(next http.Handler) http.Handler {
				return next
			}
		})
	mockTelemetryMiddleware.EXPECT().
		Trace(gomock.Any()).
		AnyTimes().
		DoAndReturn(func(next http.Handler) http.Handler {
			return next
		})

	router := NewRouter(
		cfg,
		mockAuthenticationMiddleware,
		mockAuthorizationMiddleware,
		mockTelemetryMiddleware,
		mockHealthController,
		mockPermissionsController,
		mockRolesController,
		mockScopesController,
		mockTokensController,
		mockUsersController,
	)

	req := httptest.NewRequest(http.MethodHead, "/health", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
