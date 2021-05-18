// Code generated by MockGen. DO NOT EDIT.
// Source: server.go

// Package mock_resources is a generated GoMock package.
package mock_resources

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	resources "github.com/ionos-cloud/ionosctl/pkg/resources"
)

// MockServersService is a mock of ServersService interface.
type MockServersService struct {
	ctrl     *gomock.Controller
	recorder *MockServersServiceMockRecorder
}

// MockServersServiceMockRecorder is the mock recorder for MockServersService.
type MockServersServiceMockRecorder struct {
	mock *MockServersService
}

// NewMockServersService creates a new mock instance.
func NewMockServersService(ctrl *gomock.Controller) *MockServersService {
	mock := &MockServersService{ctrl: ctrl}
	mock.recorder = &MockServersServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockServersService) EXPECT() *MockServersServiceMockRecorder {
	return m.recorder
}

// AttachCdrom mocks base method.
func (m *MockServersService) AttachCdrom(datacenterId, serverId string, cdrom resources.Image) (*resources.Image, *resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AttachCdrom", datacenterId, serverId, cdrom)
	ret0, _ := ret[0].(*resources.Image)
	ret1, _ := ret[1].(*resources.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// AttachCdrom indicates an expected call of AttachCdrom.
func (mr *MockServersServiceMockRecorder) AttachCdrom(datacenterId, serverId, cdrom interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AttachCdrom", reflect.TypeOf((*MockServersService)(nil).AttachCdrom), datacenterId, serverId, cdrom)
}

// AttachVolume mocks base method.
func (m *MockServersService) AttachVolume(datacenterId, serverId, volumeId string) (*resources.Volume, *resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AttachVolume", datacenterId, serverId, volumeId)
	ret0, _ := ret[0].(*resources.Volume)
	ret1, _ := ret[1].(*resources.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// AttachVolume indicates an expected call of AttachVolume.
func (mr *MockServersServiceMockRecorder) AttachVolume(datacenterId, serverId, volumeId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AttachVolume", reflect.TypeOf((*MockServersService)(nil).AttachVolume), datacenterId, serverId, volumeId)
}

// Create mocks base method.
func (m *MockServersService) Create(name, cpufamily, datacenterId, zone string, cores, ram int32) (*resources.Server, *resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", name, cpufamily, datacenterId, zone, cores, ram)
	ret0, _ := ret[0].(*resources.Server)
	ret1, _ := ret[1].(*resources.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Create indicates an expected call of Create.
func (mr *MockServersServiceMockRecorder) Create(name, cpufamily, datacenterId, zone, cores, ram interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockServersService)(nil).Create), name, cpufamily, datacenterId, zone, cores, ram)
}

// Delete mocks base method.
func (m *MockServersService) Delete(datacenterId, serverId string) (*resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", datacenterId, serverId)
	ret0, _ := ret[0].(*resources.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Delete indicates an expected call of Delete.
func (mr *MockServersServiceMockRecorder) Delete(datacenterId, serverId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockServersService)(nil).Delete), datacenterId, serverId)
}

// DetachCdrom mocks base method.
func (m *MockServersService) DetachCdrom(datacenterId, serverId, cdromId string) (*resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DetachCdrom", datacenterId, serverId, cdromId)
	ret0, _ := ret[0].(*resources.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DetachCdrom indicates an expected call of DetachCdrom.
func (mr *MockServersServiceMockRecorder) DetachCdrom(datacenterId, serverId, cdromId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DetachCdrom", reflect.TypeOf((*MockServersService)(nil).DetachCdrom), datacenterId, serverId, cdromId)
}

// DetachVolume mocks base method.
func (m *MockServersService) DetachVolume(datacenterId, serverId, volumeId string) (*resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DetachVolume", datacenterId, serverId, volumeId)
	ret0, _ := ret[0].(*resources.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DetachVolume indicates an expected call of DetachVolume.
func (mr *MockServersServiceMockRecorder) DetachVolume(datacenterId, serverId, volumeId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DetachVolume", reflect.TypeOf((*MockServersService)(nil).DetachVolume), datacenterId, serverId, volumeId)
}

// Get mocks base method.
func (m *MockServersService) Get(datacenterId, serverId string) (*resources.Server, *resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", datacenterId, serverId)
	ret0, _ := ret[0].(*resources.Server)
	ret1, _ := ret[1].(*resources.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Get indicates an expected call of Get.
func (mr *MockServersServiceMockRecorder) Get(datacenterId, serverId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockServersService)(nil).Get), datacenterId, serverId)
}

// GetCdrom mocks base method.
func (m *MockServersService) GetCdrom(datacenterId, serverId, cdromId string) (*resources.Image, *resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCdrom", datacenterId, serverId, cdromId)
	ret0, _ := ret[0].(*resources.Image)
	ret1, _ := ret[1].(*resources.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetCdrom indicates an expected call of GetCdrom.
func (mr *MockServersServiceMockRecorder) GetCdrom(datacenterId, serverId, cdromId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCdrom", reflect.TypeOf((*MockServersService)(nil).GetCdrom), datacenterId, serverId, cdromId)
}

// GetVolume mocks base method.
func (m *MockServersService) GetVolume(datacenterId, serverId, volumeId string) (*resources.Volume, *resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetVolume", datacenterId, serverId, volumeId)
	ret0, _ := ret[0].(*resources.Volume)
	ret1, _ := ret[1].(*resources.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetVolume indicates an expected call of GetVolume.
func (mr *MockServersServiceMockRecorder) GetVolume(datacenterId, serverId, volumeId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetVolume", reflect.TypeOf((*MockServersService)(nil).GetVolume), datacenterId, serverId, volumeId)
}

// List mocks base method.
func (m *MockServersService) List(datacenterId string) (resources.Servers, *resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", datacenterId)
	ret0, _ := ret[0].(resources.Servers)
	ret1, _ := ret[1].(*resources.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// List indicates an expected call of List.
func (mr *MockServersServiceMockRecorder) List(datacenterId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockServersService)(nil).List), datacenterId)
}

// ListCdroms mocks base method.
func (m *MockServersService) ListCdroms(datacenterId, serverId string) (resources.Cdroms, *resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListCdroms", datacenterId, serverId)
	ret0, _ := ret[0].(resources.Cdroms)
	ret1, _ := ret[1].(*resources.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ListCdroms indicates an expected call of ListCdroms.
func (mr *MockServersServiceMockRecorder) ListCdroms(datacenterId, serverId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListCdroms", reflect.TypeOf((*MockServersService)(nil).ListCdroms), datacenterId, serverId)
}

// ListVolumes mocks base method.
func (m *MockServersService) ListVolumes(datacenterId, serverId string) (resources.AttachedVolumes, *resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListVolumes", datacenterId, serverId)
	ret0, _ := ret[0].(resources.AttachedVolumes)
	ret1, _ := ret[1].(*resources.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ListVolumes indicates an expected call of ListVolumes.
func (mr *MockServersServiceMockRecorder) ListVolumes(datacenterId, serverId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListVolumes", reflect.TypeOf((*MockServersService)(nil).ListVolumes), datacenterId, serverId)
}

// Reboot mocks base method.
func (m *MockServersService) Reboot(datacenterId, serverId string) (*resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Reboot", datacenterId, serverId)
	ret0, _ := ret[0].(*resources.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Reboot indicates an expected call of Reboot.
func (mr *MockServersServiceMockRecorder) Reboot(datacenterId, serverId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Reboot", reflect.TypeOf((*MockServersService)(nil).Reboot), datacenterId, serverId)
}

// Start mocks base method.
func (m *MockServersService) Start(datacenterId, serverId string) (*resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Start", datacenterId, serverId)
	ret0, _ := ret[0].(*resources.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Start indicates an expected call of Start.
func (mr *MockServersServiceMockRecorder) Start(datacenterId, serverId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Start", reflect.TypeOf((*MockServersService)(nil).Start), datacenterId, serverId)
}

// Stop mocks base method.
func (m *MockServersService) Stop(datacenterId, serverId string) (*resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Stop", datacenterId, serverId)
	ret0, _ := ret[0].(*resources.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Stop indicates an expected call of Stop.
func (mr *MockServersServiceMockRecorder) Stop(datacenterId, serverId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stop", reflect.TypeOf((*MockServersService)(nil).Stop), datacenterId, serverId)
}

// Update mocks base method.
func (m *MockServersService) Update(datacenterId, serverId string, input resources.ServerProperties) (*resources.Server, *resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", datacenterId, serverId, input)
	ret0, _ := ret[0].(*resources.Server)
	ret1, _ := ret[1].(*resources.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Update indicates an expected call of Update.
func (mr *MockServersServiceMockRecorder) Update(datacenterId, serverId, input interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockServersService)(nil).Update), datacenterId, serverId, input)
}
