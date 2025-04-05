package rpcs

import (
	proto "loki-backoffice/internal/app/rpcs/proto/sso/v1"
)

type Registry struct {
	permissionClient proto.PermissionServiceClient
	scopeClient      proto.ScopeServiceClient
}

func NewRegistry(client Client) *Registry {
	return &Registry{
		permissionClient: proto.NewPermissionServiceClient(client.Connection()),
		scopeClient:      proto.NewScopeServiceClient(client.Connection()),
	}
}

func (r *Registry) GetPermissionClient() proto.PermissionServiceClient {
	return r.permissionClient
}

func (r *Registry) GetScopeClient() proto.ScopeServiceClient {
	return r.scopeClient
}
