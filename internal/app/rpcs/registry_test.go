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
	assert.NotNil(t, registry.client)

	permissionClient := registry.GetPermissionClient()
	assert.NotNil(t, permissionClient)
}
