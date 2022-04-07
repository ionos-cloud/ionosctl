// Code generated by MockGen. DO NOT EDIT.
// Source: restore.go

// Package mock_resources is a generated GoMock package.
package mock_resources

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	resources "github.com/ionos-cloud/ionosctl/services/dbaas-postgres/resources"
)

// MockRestoresService is a mock of RestoresService interface.
type MockRestoresService struct {
	ctrl     *gomock.Controller
	recorder *MockRestoresServiceMockRecorder
}

// MockRestoresServiceMockRecorder is the mock recorder for MockRestoresService.
type MockRestoresServiceMockRecorder struct {
	mock *MockRestoresService
}

// NewMockRestoresService creates a new mock instance.
func NewMockRestoresService(ctrl *gomock.Controller) *MockRestoresService {
	mock := &MockRestoresService{ctrl: ctrl}
	mock.recorder = &MockRestoresServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRestoresService) EXPECT() *MockRestoresServiceMockRecorder {
	return m.recorder
}

// Restore mocks base method.
func (m *MockRestoresService) Restore(clusterId string, input resources.CreateRestoreRequest) (*resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Restore", clusterId, input)
	ret0, _ := ret[0].(*resources.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Restore indicates an expected call of Restore.
func (mr *MockRestoresServiceMockRecorder) Restore(clusterId, input interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Restore", reflect.TypeOf((*MockRestoresService)(nil).Restore), clusterId, input)
}
