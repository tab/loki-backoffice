package interceptors

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc"

	"loki-backoffice/internal/config"
	"loki-backoffice/internal/config/logger"
)

func Test_LoggerInterceptor_Log(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		AppEnv:   "test",
		AppAddr:  "localhost:8080",
		LogLevel: "info",
	}
	log := logger.NewLogger(cfg)
	interceptorInstance := NewLoggerInterceptor(log)

	mockInvoker := func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
		return nil
	}

	tests := []struct {
		name   string
		method string
		error  bool
	}{
		{
			name:   "Success",
			method: "/sso.v1.PermissionService/List",
			error:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			interceptor := interceptorInstance.Log()

			err := interceptor(context.Background(), tt.method, nil, nil, nil, mockInvoker)
			assert.NoError(t, err)
		})
	}
}
