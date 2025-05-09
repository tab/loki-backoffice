// Code generated by MockGen. DO NOT EDIT.
// Source: internal/app/services/roles.go
//
// Generated by this command:
//
//	mockgen -source=internal/app/services/roles.go -destination=internal/app/services/roles_mock.go -package=services
//

// Package services is a generated GoMock package.
package services

import (
	context "context"
	models "loki-backoffice/internal/app/models"
	reflect "reflect"

	uuid "github.com/google/uuid"
	gomock "go.uber.org/mock/gomock"
)

// MockRoles is a mock of Roles interface.
type MockRoles struct {
	ctrl     *gomock.Controller
	recorder *MockRolesMockRecorder
	isgomock struct{}
}

// MockRolesMockRecorder is the mock recorder for MockRoles.
type MockRolesMockRecorder struct {
	mock *MockRoles
}

// NewMockRoles creates a new mock instance.
func NewMockRoles(ctrl *gomock.Controller) *MockRoles {
	mock := &MockRoles{ctrl: ctrl}
	mock.recorder = &MockRolesMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRoles) EXPECT() *MockRolesMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockRoles) Create(ctx context.Context, params *models.Role) (*models.Role, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, params)
	ret0, _ := ret[0].(*models.Role)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockRolesMockRecorder) Create(ctx, params any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockRoles)(nil).Create), ctx, params)
}

// Delete mocks base method.
func (m *MockRoles) Delete(ctx context.Context, id uuid.UUID) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, id)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Delete indicates an expected call of Delete.
func (mr *MockRolesMockRecorder) Delete(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockRoles)(nil).Delete), ctx, id)
}

// FindById mocks base method.
func (m *MockRoles) FindById(ctx context.Context, id uuid.UUID) (*models.Role, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindById", ctx, id)
	ret0, _ := ret[0].(*models.Role)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindById indicates an expected call of FindById.
func (mr *MockRolesMockRecorder) FindById(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindById", reflect.TypeOf((*MockRoles)(nil).FindById), ctx, id)
}

// List mocks base method.
func (m *MockRoles) List(ctx context.Context, pagination *Pagination) ([]models.Role, uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx, pagination)
	ret0, _ := ret[0].([]models.Role)
	ret1, _ := ret[1].(uint64)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// List indicates an expected call of List.
func (mr *MockRolesMockRecorder) List(ctx, pagination any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockRoles)(nil).List), ctx, pagination)
}

// Update mocks base method.
func (m *MockRoles) Update(ctx context.Context, params *models.Role) (*models.Role, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, params)
	ret0, _ := ret[0].(*models.Role)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockRolesMockRecorder) Update(ctx, params any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockRoles)(nil).Update), ctx, params)
}
