// Code generated by MockGen. DO NOT EDIT.
// Source: contract.go

// Package mock_resources is a generated GoMock package.
package mock_resources

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	resources "github.com/ionos-cloud/ionosctl/pkg/resources"
)

// MockContractsService is a mock of ContractsService interface.
type MockContractsService struct {
	ctrl     *gomock.Controller
	recorder *MockContractsServiceMockRecorder
}

// MockContractsServiceMockRecorder is the mock recorder for MockContractsService.
type MockContractsServiceMockRecorder struct {
	mock *MockContractsService
}

// NewMockContractsService creates a new mock instance.
func NewMockContractsService(ctrl *gomock.Controller) *MockContractsService {
	mock := &MockContractsService{ctrl: ctrl}
	mock.recorder = &MockContractsServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockContractsService) EXPECT() *MockContractsServiceMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockContractsService) Get() (resources.Contract, *resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get")
	ret0, _ := ret[0].(resources.Contract)
	ret1, _ := ret[1].(*resources.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Get indicates an expected call of Get.
func (mr *MockContractsServiceMockRecorder) Get() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockContractsService)(nil).Get))
}
