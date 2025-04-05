package services

import (
	"go.uber.org/fx"

	"loki-backoffice/internal/app/rpcs"
	proto "loki-backoffice/internal/app/rpcs/proto/sso/v1"
)

var Module = fx.Options(
	fx.Provide(NewHealthChecker),
	fx.Provide(
		func(registry *rpcs.Registry) proto.PermissionServiceClient {
			return registry.GetPermissionClient()
		},
		NewPermissions,
	),
	fx.Provide(
		func(registry *rpcs.Registry) proto.RoleServiceClient {
			return registry.GetRoleClient()
		},
		NewRoles,
	),
	fx.Provide(
		func(registry *rpcs.Registry) proto.ScopeServiceClient {
			return registry.GetScopeClient()
		},
		NewScopes,
	),
	fx.Provide(
		func(registry *rpcs.Registry) proto.TokenServiceClient {
			return registry.GetTokenClient()
		},
		NewTokens,
	),
)
