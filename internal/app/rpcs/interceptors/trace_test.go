package interceptors

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc"

	"loki-backoffice/internal/config/middlewares"
)

func TestTraceInterceptor_Trace(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	interceptorInstance := NewTraceInterceptor()
	mockInvoker := func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
		return nil
	}

	tests := []struct {
		name  string
		ctx   context.Context
		error bool
	}{
		{
			name:  "Success",
			ctx:   context.WithValue(context.Background(), middlewares.TraceId{}, "test-trace-id"),
			error: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			interceptor := interceptorInstance.Trace()
			err := interceptor(tt.ctx, "test-method", nil, nil, nil, mockInvoker)
			assert.NoError(t, err)
		})
	}
}
