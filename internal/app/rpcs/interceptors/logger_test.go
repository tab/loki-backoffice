package interceptors

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc"

	"loki-backoffice/pkg/logger"
)

func Test_LoggerInterceptor_Log(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	log := logger.NewLogger()
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
