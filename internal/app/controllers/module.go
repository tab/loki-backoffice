package controllers

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewHealthController),
	fx.Provide(NewPermissionsController),
	fx.Provide(NewScopesController),
)
