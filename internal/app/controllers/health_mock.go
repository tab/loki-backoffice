// Code generated by MockGen. DO NOT EDIT.
// Source: internal/app/controllers/health.go
//
// Generated by this command:
//
//	mockgen -source=internal/app/controllers/health.go -destination=internal/app/controllers/health_mock.go -package=controllers
//

// Package controllers is a generated GoMock package.
package controllers

import (
	http "net/http"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockHealthController is a mock of HealthController interface.
type MockHealthController struct {
	ctrl     *gomock.Controller
	recorder *MockHealthControllerMockRecorder
	isgomock struct{}
}

// MockHealthControllerMockRecorder is the mock recorder for MockHealthController.
type MockHealthControllerMockRecorder struct {
	mock *MockHealthController
}

// NewMockHealthController creates a new mock instance.
func NewMockHealthController(ctrl *gomock.Controller) *MockHealthController {
	mock := &MockHealthController{ctrl: ctrl}
	mock.recorder = &MockHealthControllerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockHealthController) EXPECT() *MockHealthControllerMockRecorder {
	return m.recorder
}

// HandleLiveness mocks base method.
func (m *MockHealthController) HandleLiveness(w http.ResponseWriter, r *http.Request) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "HandleLiveness", w, r)
}

// HandleLiveness indicates an expected call of HandleLiveness.
func (mr *MockHealthControllerMockRecorder) HandleLiveness(w, r any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HandleLiveness", reflect.TypeOf((*MockHealthController)(nil).HandleLiveness), w, r)
}

// HandleReadiness mocks base method.
func (m *MockHealthController) HandleReadiness(w http.ResponseWriter, r *http.Request) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "HandleReadiness", w, r)
}

// HandleReadiness indicates an expected call of HandleReadiness.
func (mr *MockHealthControllerMockRecorder) HandleReadiness(w, r any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HandleReadiness", reflect.TypeOf((*MockHealthController)(nil).HandleReadiness), w, r)
}
