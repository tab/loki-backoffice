package interceptors

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"loki-backoffice/internal/config/logger"
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

		traceId := extractTraceId(ctx)
		requestId := extractRequestId(ctx)

		err := invoker(ctx, method, req, reply, cc, opts...)

		code := status.Code(err).String()
		duration := time.Since(startTime)

		reqLogger := i.log.
			WithComponent("gRPC").
			WithRequestId(requestId).
			WithTraceId(traceId)

		reqLogger.Info().
			Str("method", method).
			Str("status", code).
			Dur("duration", duration).
			Msgf("%s - %s in %s", method, code, duration)

		return err
	}
}

func extractTraceId(ctx context.Context) string {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		return ""
	}

	traceId := md.Get(TraceId)
	if len(traceId) == 0 {
		return ""
	}

	return traceId[0]
}

func extractRequestId(ctx context.Context) string {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		return ""
	}

	requestId := md.Get(RequestId)
	if len(requestId) == 0 {
		return ""
	}

	return requestId[0]
}
