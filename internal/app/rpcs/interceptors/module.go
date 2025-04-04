package interceptors

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewAuthenticationInterceptor),
	fx.Provide(NewLoggerInterceptor),
	fx.Provide(NewTraceInterceptor),
)
