package controllers

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewHealthController),
	fx.Provide(NewPermissionsController),
	fx.Provide(NewRolesController),
	fx.Provide(NewScopesController),
	fx.Provide(NewTokensController),
	fx.Provide(NewUsersController),
)
