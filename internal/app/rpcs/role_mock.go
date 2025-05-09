// Code generated by MockGen. DO NOT EDIT.
// Source: internal/app/rpcs/proto/sso/v1/role_grpc.pb.go
//
// Generated by this command:
//
//	mockgen -source=internal/app/rpcs/proto/sso/v1/role_grpc.pb.go -destination=internal/app/rpcs/role_mock.go -package=rpcs
//

// Package rpcs is a generated GoMock package.
package rpcs

import (
	context "context"
	ssov1 "loki-backoffice/internal/app/rpcs/proto/sso/v1"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
	grpc "google.golang.org/grpc"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// MockRoleServiceClient is a mock of RoleServiceClient interface.
type MockRoleServiceClient struct {
	ctrl     *gomock.Controller
	recorder *MockRoleServiceClientMockRecorder
	isgomock struct{}
}

// MockRoleServiceClientMockRecorder is the mock recorder for MockRoleServiceClient.
type MockRoleServiceClientMockRecorder struct {
	mock *MockRoleServiceClient
}

// NewMockRoleServiceClient creates a new mock instance.
func NewMockRoleServiceClient(ctrl *gomock.Controller) *MockRoleServiceClient {
	mock := &MockRoleServiceClient{ctrl: ctrl}
	mock.recorder = &MockRoleServiceClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRoleServiceClient) EXPECT() *MockRoleServiceClientMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockRoleServiceClient) Create(ctx context.Context, in *ssov1.CreateRoleRequest, opts ...grpc.CallOption) (*ssov1.CreateRoleResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Create", varargs...)
	ret0, _ := ret[0].(*ssov1.CreateRoleResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockRoleServiceClientMockRecorder) Create(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockRoleServiceClient)(nil).Create), varargs...)
}

// Delete mocks base method.
func (m *MockRoleServiceClient) Delete(ctx context.Context, in *ssov1.DeleteRoleRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Delete", varargs...)
	ret0, _ := ret[0].(*emptypb.Empty)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Delete indicates an expected call of Delete.
func (mr *MockRoleServiceClientMockRecorder) Delete(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockRoleServiceClient)(nil).Delete), varargs...)
}

// Get mocks base method.
func (m *MockRoleServiceClient) Get(ctx context.Context, in *ssov1.GetRoleRequest, opts ...grpc.CallOption) (*ssov1.GetRoleResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Get", varargs...)
	ret0, _ := ret[0].(*ssov1.GetRoleResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockRoleServiceClientMockRecorder) Get(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockRoleServiceClient)(nil).Get), varargs...)
}

// List mocks base method.
func (m *MockRoleServiceClient) List(ctx context.Context, in *ssov1.PaginatedListRequest, opts ...grpc.CallOption) (*ssov1.ListRolesResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "List", varargs...)
	ret0, _ := ret[0].(*ssov1.ListRolesResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockRoleServiceClientMockRecorder) List(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockRoleServiceClient)(nil).List), varargs...)
}

// Update mocks base method.
func (m *MockRoleServiceClient) Update(ctx context.Context, in *ssov1.UpdateRoleRequest, opts ...grpc.CallOption) (*ssov1.UpdateRoleResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Update", varargs...)
	ret0, _ := ret[0].(*ssov1.UpdateRoleResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockRoleServiceClientMockRecorder) Update(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockRoleServiceClient)(nil).Update), varargs...)
}

// MockRoleServiceServer is a mock of RoleServiceServer interface.
type MockRoleServiceServer struct {
	ctrl     *gomock.Controller
	recorder *MockRoleServiceServerMockRecorder
	isgomock struct{}
}

// MockRoleServiceServerMockRecorder is the mock recorder for MockRoleServiceServer.
type MockRoleServiceServerMockRecorder struct {
	mock *MockRoleServiceServer
}

// NewMockRoleServiceServer creates a new mock instance.
func NewMockRoleServiceServer(ctrl *gomock.Controller) *MockRoleServiceServer {
	mock := &MockRoleServiceServer{ctrl: ctrl}
	mock.recorder = &MockRoleServiceServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRoleServiceServer) EXPECT() *MockRoleServiceServerMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockRoleServiceServer) Create(arg0 context.Context, arg1 *ssov1.CreateRoleRequest) (*ssov1.CreateRoleResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1)
	ret0, _ := ret[0].(*ssov1.CreateRoleResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockRoleServiceServerMockRecorder) Create(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockRoleServiceServer)(nil).Create), arg0, arg1)
}

// Delete mocks base method.
func (m *MockRoleServiceServer) Delete(arg0 context.Context, arg1 *ssov1.DeleteRoleRequest) (*emptypb.Empty, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0, arg1)
	ret0, _ := ret[0].(*emptypb.Empty)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Delete indicates an expected call of Delete.
func (mr *MockRoleServiceServerMockRecorder) Delete(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockRoleServiceServer)(nil).Delete), arg0, arg1)
}

// Get mocks base method.
func (m *MockRoleServiceServer) Get(arg0 context.Context, arg1 *ssov1.GetRoleRequest) (*ssov1.GetRoleResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0, arg1)
	ret0, _ := ret[0].(*ssov1.GetRoleResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockRoleServiceServerMockRecorder) Get(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockRoleServiceServer)(nil).Get), arg0, arg1)
}

// List mocks base method.
func (m *MockRoleServiceServer) List(arg0 context.Context, arg1 *ssov1.PaginatedListRequest) (*ssov1.ListRolesResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", arg0, arg1)
	ret0, _ := ret[0].(*ssov1.ListRolesResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockRoleServiceServerMockRecorder) List(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockRoleServiceServer)(nil).List), arg0, arg1)
}

// Update mocks base method.
func (m *MockRoleServiceServer) Update(arg0 context.Context, arg1 *ssov1.UpdateRoleRequest) (*ssov1.UpdateRoleResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0, arg1)
	ret0, _ := ret[0].(*ssov1.UpdateRoleResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockRoleServiceServerMockRecorder) Update(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockRoleServiceServer)(nil).Update), arg0, arg1)
}

// mustEmbedUnimplementedRoleServiceServer mocks base method.
func (m *MockRoleServiceServer) mustEmbedUnimplementedRoleServiceServer() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "mustEmbedUnimplementedRoleServiceServer")
}

// mustEmbedUnimplementedRoleServiceServer indicates an expected call of mustEmbedUnimplementedRoleServiceServer.
func (mr *MockRoleServiceServerMockRecorder) mustEmbedUnimplementedRoleServiceServer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "mustEmbedUnimplementedRoleServiceServer", reflect.TypeOf((*MockRoleServiceServer)(nil).mustEmbedUnimplementedRoleServiceServer))
}

// MockUnsafeRoleServiceServer is a mock of UnsafeRoleServiceServer interface.
type MockUnsafeRoleServiceServer struct {
	ctrl     *gomock.Controller
	recorder *MockUnsafeRoleServiceServerMockRecorder
	isgomock struct{}
}

// MockUnsafeRoleServiceServerMockRecorder is the mock recorder for MockUnsafeRoleServiceServer.
type MockUnsafeRoleServiceServerMockRecorder struct {
	mock *MockUnsafeRoleServiceServer
}

// NewMockUnsafeRoleServiceServer creates a new mock instance.
func NewMockUnsafeRoleServiceServer(ctrl *gomock.Controller) *MockUnsafeRoleServiceServer {
	mock := &MockUnsafeRoleServiceServer{ctrl: ctrl}
	mock.recorder = &MockUnsafeRoleServiceServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUnsafeRoleServiceServer) EXPECT() *MockUnsafeRoleServiceServerMockRecorder {
	return m.recorder
}

// mustEmbedUnimplementedRoleServiceServer mocks base method.
func (m *MockUnsafeRoleServiceServer) mustEmbedUnimplementedRoleServiceServer() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "mustEmbedUnimplementedRoleServiceServer")
}

// mustEmbedUnimplementedRoleServiceServer indicates an expected call of mustEmbedUnimplementedRoleServiceServer.
func (mr *MockUnsafeRoleServiceServerMockRecorder) mustEmbedUnimplementedRoleServiceServer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "mustEmbedUnimplementedRoleServiceServer", reflect.TypeOf((*MockUnsafeRoleServiceServer)(nil).mustEmbedUnimplementedRoleServiceServer))
}
