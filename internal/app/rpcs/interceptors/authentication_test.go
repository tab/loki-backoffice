package interceptors

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc"

	"loki-backoffice/internal/config/middlewares"
	"loki-backoffice/pkg/logger"
)

func Test_AuthenticationInterceptor_Authenticate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	log := logger.NewLogger()
	authInterceptor := NewAuthenticationInterceptor(log)

	mockInvoker := func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
		return nil
	}

	tests := []struct {
		name     string
		ctx      context.Context
		expected string
	}{
		{
			name:     "Success",
			ctx:      context.WithValue(context.Background(), middlewares.Token{}, "test-token"),
			expected: "test-token",
		},
		{
			name:     "No token in context",
			ctx:      context.Background(),
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			interceptor := authInterceptor.Authenticate()
			err := interceptor(tt.ctx, "test-method", nil, nil, nil, mockInvoker)
			assert.NoError(t, err)

			token, _ := extractBearerToken(tt.ctx)
			assert.Equal(t, tt.expected, token)
		})
	}
}
