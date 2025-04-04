package rpcs

import (
	proto "loki-backoffice/internal/app/rpcs/proto/sso/v1"
)

type Registry struct {
	client proto.PermissionServiceClient
}

func NewRegistry(client Client) *Registry {
	return &Registry{
		client: proto.NewPermissionServiceClient(client.Connection()),
	}
}

func (r *Registry) GetPermissionClient() proto.PermissionServiceClient {
	return r.client
}
