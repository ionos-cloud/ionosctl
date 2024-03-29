// Code generated by MockGen. DO NOT EDIT.
// Source: ./location.go

// Package mock_resources is a generated GoMock package.
package mock_resources

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	ionoscloud "github.com/ionos-cloud/sdk-go-container-registry"
)

// MockLocationsService is a mock of LocationsService interface.
type MockLocationsService struct {
	ctrl     *gomock.Controller
	recorder *MockLocationsServiceMockRecorder
}

// MockLocationsServiceMockRecorder is the mock recorder for MockLocationsService.
type MockLocationsServiceMockRecorder struct {
	mock *MockLocationsService
}

// NewMockLocationsService creates a new mock instance.
func NewMockLocationsService(ctrl *gomock.Controller) *MockLocationsService {
	mock := &MockLocationsService{ctrl: ctrl}
	mock.recorder = &MockLocationsServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLocationsService) EXPECT() *MockLocationsServiceMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockLocationsService) Get() (ionoscloud.LocationsResponse, *ionoscloud.APIResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get")
	ret0, _ := ret[0].(ionoscloud.LocationsResponse)
	ret1, _ := ret[1].(*ionoscloud.APIResponse)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Get indicates an expected call of Get.
func (mr *MockLocationsServiceMockRecorder) Get() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockLocationsService)(nil).Get))
}
