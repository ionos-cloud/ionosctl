// Code generated by MockGen. DO NOT EDIT.
// Source: volume.go

// Package mock_resources is a generated GoMock package.
package mock_resources

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	resources "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
)

// MockVolumesService is a mock of VolumesService interface.
type MockVolumesService struct {
	ctrl     *gomock.Controller
	recorder *MockVolumesServiceMockRecorder
}

// MockVolumesServiceMockRecorder is the mock recorder for MockVolumesService.
type MockVolumesServiceMockRecorder struct {
	mock *MockVolumesService
}

// NewMockVolumesService creates a new mock instance.
func NewMockVolumesService(ctrl *gomock.Controller) *MockVolumesService {
	mock := &MockVolumesService{ctrl: ctrl}
	mock.recorder = &MockVolumesServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockVolumesService) EXPECT() *MockVolumesServiceMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockVolumesService) Create(datacenterId string, input resources.Volume, params resources.QueryParams) (*resources.Volume, *resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", datacenterId, input, params)
	ret0, _ := ret[0].(*resources.Volume)
	ret1, _ := ret[1].(*resources.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Create indicates an expected call of Create.
func (mr *MockVolumesServiceMockRecorder) Create(datacenterId, input, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockVolumesService)(nil).Create), datacenterId, input, params)
}

// Delete mocks base method.
func (m *MockVolumesService) Delete(datacenterId, volumeId string, params resources.QueryParams) (*resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", datacenterId, volumeId, params)
	ret0, _ := ret[0].(*resources.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Delete indicates an expected call of Delete.
func (mr *MockVolumesServiceMockRecorder) Delete(datacenterId, volumeId, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockVolumesService)(nil).Delete), datacenterId, volumeId, params)
}

// Get mocks base method.
func (m *MockVolumesService) Get(datacenterId, volumeId string, params resources.QueryParams) (*resources.Volume, *resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", datacenterId, volumeId, params)
	ret0, _ := ret[0].(*resources.Volume)
	ret1, _ := ret[1].(*resources.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Get indicates an expected call of Get.
func (mr *MockVolumesServiceMockRecorder) Get(datacenterId, volumeId, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockVolumesService)(nil).Get), datacenterId, volumeId, params)
}

// List mocks base method.
func (m *MockVolumesService) List(datacenterId string, params resources.ListQueryParams) (resources.Volumes, *resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", datacenterId, params)
	ret0, _ := ret[0].(resources.Volumes)
	ret1, _ := ret[1].(*resources.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// List indicates an expected call of List.
func (mr *MockVolumesServiceMockRecorder) List(datacenterId, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockVolumesService)(nil).List), datacenterId, params)
}

// Update mocks base method.
func (m *MockVolumesService) Update(datacenterId, volumeId string, input resources.VolumeProperties, params resources.QueryParams) (*resources.Volume, *resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", datacenterId, volumeId, input, params)
	ret0, _ := ret[0].(*resources.Volume)
	ret1, _ := ret[1].(*resources.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Update indicates an expected call of Update.
func (mr *MockVolumesServiceMockRecorder) Update(datacenterId, volumeId, input, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockVolumesService)(nil).Update), datacenterId, volumeId, input, params)
}
