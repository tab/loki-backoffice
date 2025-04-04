package interceptors

import (
	"context"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"loki-backoffice/internal/config/middlewares"
)

const (
	RequestId = "X-Request-Id"
	TraceId   = "X-Trace-Id"
)

type TraceInterceptor interface {
	Trace() grpc.UnaryClientInterceptor
}

type traceInterceptor struct{}

func NewTraceInterceptor() TraceInterceptor {
	return &traceInterceptor{}
}

func (i *traceInterceptor) Trace() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		requestId := middleware.GetReqID(ctx)

		if requestId == "" {
			requestId = uuid.New().String()
		}

		traceId, ok := ctx.Value(middlewares.TraceId{}).(string)
		if !ok || traceId == "" {
			traceId = uuid.New().String()
		}

		ctx = metadata.AppendToOutgoingContext(ctx,
			RequestId, requestId,
			TraceId, traceId,
			middlewares.AuthenticationTraceKey, traceId)

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
