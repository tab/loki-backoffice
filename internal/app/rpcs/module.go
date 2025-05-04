package rpcs

import (
	"go.uber.org/fx"

	"loki-backoffice/internal/app/rpcs/interceptors"
)

var Module = fx.Options(
	interceptors.Module,

	fx.Provide(NewClient),
	fx.Provide(NewRegistry),
)
