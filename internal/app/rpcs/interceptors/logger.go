package interceptors

import (
	"context"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

	"loki-backoffice/pkg/logger"
)

type LoggerInterceptor interface {
	Log() grpc.UnaryClientInterceptor
}

type loggerInterceptor struct {
	log *logger.Logger
}

func NewLoggerInterceptor(log *logger.Logger) LoggerInterceptor {
	return &loggerInterceptor{
		log: log,
	}
}

func (i *loggerInterceptor) Log() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		startTime := time.Now()

		err := invoker(ctx, method, req, reply, cc, opts...)

		requestId := middleware.GetReqID(ctx)
		code := status.Code(err).String()
		duration := time.Since(startTime)

		i.log.Info().Msgf("[%s] %s - %s in %s", requestId, method, code, duration)

		return err
	}
}
