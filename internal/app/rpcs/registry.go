package rpcs

import (
	proto "loki-backoffice/internal/app/rpcs/proto/sso/v1"
)

type Registry struct {
	permissionClient proto.PermissionServiceClient
	roleClient       proto.RoleServiceClient
	scopeClient      proto.ScopeServiceClient
	tokenClient      proto.TokenServiceClient
	userClient       proto.UserServiceClient
}

func NewRegistry(client Client) *Registry {
	return &Registry{
		permissionClient: proto.NewPermissionServiceClient(client.Connection()),
		roleClient:       proto.NewRoleServiceClient(client.Connection()),
		scopeClient:      proto.NewScopeServiceClient(client.Connection()),
		tokenClient:      proto.NewTokenServiceClient(client.Connection()),
		userClient:       proto.NewUserServiceClient(client.Connection()),
	}
}

func (r *Registry) GetPermissionClient() proto.PermissionServiceClient {
	return r.permissionClient
}

func (r *Registry) GetRoleClient() proto.RoleServiceClient {
	return r.roleClient
}

func (r *Registry) GetScopeClient() proto.ScopeServiceClient {
	return r.scopeClient
}

func (r *Registry) GetTokenClient() proto.TokenServiceClient {
	return r.tokenClient
}

func (r *Registry) GetUserClient() proto.UserServiceClient {
	return r.userClient
}
