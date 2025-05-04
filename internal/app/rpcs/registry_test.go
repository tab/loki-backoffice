package rpcs

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc"
)

func Test_NewRegistry(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := NewMockClient(ctrl)
	mockConn := &grpc.ClientConn{}
	mockClient.EXPECT().Connection().Return(mockConn).AnyTimes()

	registry := NewRegistry(mockClient)

	assert.NotNil(t, registry)
	assert.NotNil(t, registry.permissionClient)
	assert.NotNil(t, registry.roleClient)
	assert.NotNil(t, registry.scopeClient)

	permissionClient := registry.GetPermissionClient()
	assert.NotNil(t, permissionClient)

	roleClient := registry.GetRoleClient()
	assert.NotNil(t, roleClient)

	scopeClient := registry.GetScopeClient()
	assert.NotNil(t, scopeClient)

	tokenClient := registry.GetTokenClient()
	assert.NotNil(t, tokenClient)

	userClient := registry.GetUserClient()
	assert.NotNil(t, userClient)
}
