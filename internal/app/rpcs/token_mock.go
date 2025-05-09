// Code generated by MockGen. DO NOT EDIT.
// Source: internal/app/rpcs/proto/sso/v1/token_grpc.pb.go
//
// Generated by this command:
//
//	mockgen -source=internal/app/rpcs/proto/sso/v1/token_grpc.pb.go -destination=internal/app/rpcs/token_mock.go -package=rpcs
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

// MockTokenServiceClient is a mock of TokenServiceClient interface.
type MockTokenServiceClient struct {
	ctrl     *gomock.Controller
	recorder *MockTokenServiceClientMockRecorder
	isgomock struct{}
}

// MockTokenServiceClientMockRecorder is the mock recorder for MockTokenServiceClient.
type MockTokenServiceClientMockRecorder struct {
	mock *MockTokenServiceClient
}

// NewMockTokenServiceClient creates a new mock instance.
func NewMockTokenServiceClient(ctrl *gomock.Controller) *MockTokenServiceClient {
	mock := &MockTokenServiceClient{ctrl: ctrl}
	mock.recorder = &MockTokenServiceClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTokenServiceClient) EXPECT() *MockTokenServiceClientMockRecorder {
	return m.recorder
}

// Delete mocks base method.
func (m *MockTokenServiceClient) Delete(ctx context.Context, in *ssov1.DeleteTokenRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
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
func (mr *MockTokenServiceClientMockRecorder) Delete(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockTokenServiceClient)(nil).Delete), varargs...)
}

// List mocks base method.
func (m *MockTokenServiceClient) List(ctx context.Context, in *ssov1.PaginatedListRequest, opts ...grpc.CallOption) (*ssov1.ListTokensResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "List", varargs...)
	ret0, _ := ret[0].(*ssov1.ListTokensResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockTokenServiceClientMockRecorder) List(ctx, in any, opts ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockTokenServiceClient)(nil).List), varargs...)
}

// MockTokenServiceServer is a mock of TokenServiceServer interface.
type MockTokenServiceServer struct {
	ctrl     *gomock.Controller
	recorder *MockTokenServiceServerMockRecorder
	isgomock struct{}
}

// MockTokenServiceServerMockRecorder is the mock recorder for MockTokenServiceServer.
type MockTokenServiceServerMockRecorder struct {
	mock *MockTokenServiceServer
}

// NewMockTokenServiceServer creates a new mock instance.
func NewMockTokenServiceServer(ctrl *gomock.Controller) *MockTokenServiceServer {
	mock := &MockTokenServiceServer{ctrl: ctrl}
	mock.recorder = &MockTokenServiceServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTokenServiceServer) EXPECT() *MockTokenServiceServerMockRecorder {
	return m.recorder
}

// Delete mocks base method.
func (m *MockTokenServiceServer) Delete(arg0 context.Context, arg1 *ssov1.DeleteTokenRequest) (*emptypb.Empty, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0, arg1)
	ret0, _ := ret[0].(*emptypb.Empty)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Delete indicates an expected call of Delete.
func (mr *MockTokenServiceServerMockRecorder) Delete(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockTokenServiceServer)(nil).Delete), arg0, arg1)
}

// List mocks base method.
func (m *MockTokenServiceServer) List(arg0 context.Context, arg1 *ssov1.PaginatedListRequest) (*ssov1.ListTokensResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", arg0, arg1)
	ret0, _ := ret[0].(*ssov1.ListTokensResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockTokenServiceServerMockRecorder) List(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockTokenServiceServer)(nil).List), arg0, arg1)
}

// mustEmbedUnimplementedTokenServiceServer mocks base method.
func (m *MockTokenServiceServer) mustEmbedUnimplementedTokenServiceServer() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "mustEmbedUnimplementedTokenServiceServer")
}

// mustEmbedUnimplementedTokenServiceServer indicates an expected call of mustEmbedUnimplementedTokenServiceServer.
func (mr *MockTokenServiceServerMockRecorder) mustEmbedUnimplementedTokenServiceServer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "mustEmbedUnimplementedTokenServiceServer", reflect.TypeOf((*MockTokenServiceServer)(nil).mustEmbedUnimplementedTokenServiceServer))
}

// MockUnsafeTokenServiceServer is a mock of UnsafeTokenServiceServer interface.
type MockUnsafeTokenServiceServer struct {
	ctrl     *gomock.Controller
	recorder *MockUnsafeTokenServiceServerMockRecorder
	isgomock struct{}
}

// MockUnsafeTokenServiceServerMockRecorder is the mock recorder for MockUnsafeTokenServiceServer.
type MockUnsafeTokenServiceServerMockRecorder struct {
	mock *MockUnsafeTokenServiceServer
}

// NewMockUnsafeTokenServiceServer creates a new mock instance.
func NewMockUnsafeTokenServiceServer(ctrl *gomock.Controller) *MockUnsafeTokenServiceServer {
	mock := &MockUnsafeTokenServiceServer{ctrl: ctrl}
	mock.recorder = &MockUnsafeTokenServiceServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUnsafeTokenServiceServer) EXPECT() *MockUnsafeTokenServiceServerMockRecorder {
	return m.recorder
}

// mustEmbedUnimplementedTokenServiceServer mocks base method.
func (m *MockUnsafeTokenServiceServer) mustEmbedUnimplementedTokenServiceServer() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "mustEmbedUnimplementedTokenServiceServer")
}

// mustEmbedUnimplementedTokenServiceServer indicates an expected call of mustEmbedUnimplementedTokenServiceServer.
func (mr *MockUnsafeTokenServiceServerMockRecorder) mustEmbedUnimplementedTokenServiceServer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "mustEmbedUnimplementedTokenServiceServer", reflect.TypeOf((*MockUnsafeTokenServiceServer)(nil).mustEmbedUnimplementedTokenServiceServer))
}
