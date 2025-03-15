package repositories

import (
	"go.uber.org/fx"

	"loki/internal/app/repositories/postgres"
)

var Module = fx.Options(
	fx.Provide(postgres.NewPostgresClient),
	fx.Provide(NewHealthRepository),
)
