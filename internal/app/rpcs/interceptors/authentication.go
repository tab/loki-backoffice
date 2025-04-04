package interceptors

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"loki-backoffice/internal/config/middlewares"
	"loki-backoffice/pkg/logger"
)

const (
	Authorization = "authorization"
	bearerScheme  = "Bearer"
)

type AuthenticationInterceptor interface {
	Authenticate() grpc.UnaryClientInterceptor
}

type authenticationInterceptor struct {
	log *logger.Logger
}

func NewAuthenticationInterceptor(log *logger.Logger) AuthenticationInterceptor {
	return &authenticationInterceptor{
		log: log,
	}
}

func (i *authenticationInterceptor) Authenticate() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		token, ok := extractBearerToken(ctx)
		if !ok {
			i.log.Warn().Msg("No authentication token in context")
			return invoker(ctx, method, req, reply, cc, opts...)
		}

		md := metadata.Pairs(Authorization, fmt.Sprintf("%s %s", bearerScheme, token))
		ctx = metadata.NewOutgoingContext(ctx, md)

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func extractBearerToken(ctx context.Context) (string, bool) {
	t, ok := ctx.Value(middlewares.Token{}).(string)
	return t, ok
}
