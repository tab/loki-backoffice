package server

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"loki-backoffice/internal/app/controllers"
	"loki-backoffice/internal/config"
	"loki-backoffice/internal/config/middlewares"
	"loki-backoffice/internal/config/router"
)

func Test_NewServer(t *testing.T) {
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
	mockScopesController := controllers.NewMockScopesController(ctrl)

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

	appRouter := router.NewRouter(
		cfg,
		mockAuthenticationMiddleware,
		mockAuthorizationMiddleware,
		mockTelemetryMiddleware,
		mockHealthController,
		mockPermissionsController,
		mockScopesController,
	)

	srv := NewServer(cfg, appRouter)
	assert.NotNil(t, srv)

	s, ok := srv.(*server)
	assert.True(t, ok)

	assert.Equal(t, cfg.AppAddr, s.httpServer.Addr)
	assert.Equal(t, appRouter, s.httpServer.Handler)
	assert.Equal(t, 5*time.Second, s.httpServer.ReadTimeout)
	assert.Equal(t, 10*time.Second, s.httpServer.WriteTimeout)
	assert.Equal(t, 120*time.Second, s.httpServer.IdleTimeout)
}

func Test_Server_RunAndShutdown(t *testing.T) {
	cfg := &config.Config{
		AppEnv:  "test",
		AppAddr: "localhost:5000",
	}
	handler := http.NewServeMux()
	srv := NewServer(cfg, handler)

	runErrCh := make(chan error, 1)
	go func() {
		err := srv.Run()
		if err != nil && err != http.ErrServerClosed {
			runErrCh <- err
		}
		close(runErrCh)
	}()

	time.Sleep(100 * time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err := srv.Shutdown(ctx)
	assert.NoError(t, err)

	err = <-runErrCh
	assert.NoError(t, err)
}
