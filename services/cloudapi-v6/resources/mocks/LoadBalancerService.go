// Code generated by MockGen. DO NOT EDIT.
// Source: loadbalancer.go

// Package mock_resources is a generated GoMock package.
package mock_resources

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	resources "github.com/ionos-cloud/ionosctl/services/cloudapi-v6/resources"
)

// MockLoadbalancersService is a mock of LoadbalancersService interface.
type MockLoadbalancersService struct {
	ctrl     *gomock.Controller
	recorder *MockLoadbalancersServiceMockRecorder
}

// MockLoadbalancersServiceMockRecorder is the mock recorder for MockLoadbalancersService.
type MockLoadbalancersServiceMockRecorder struct {
	mock *MockLoadbalancersService
}

// NewMockLoadbalancersService creates a new mock instance.
func NewMockLoadbalancersService(ctrl *gomock.Controller) *MockLoadbalancersService {
	mock := &MockLoadbalancersService{ctrl: ctrl}
	mock.recorder = &MockLoadbalancersServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLoadbalancersService) EXPECT() *MockLoadbalancersServiceMockRecorder {
	return m.recorder
}

// AttachNic mocks base method.
func (m *MockLoadbalancersService) AttachNic(datacenterId, loadbalancerId, nicId string) (*resources.Nic, *resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AttachNic", datacenterId, loadbalancerId, nicId)
	ret0, _ := ret[0].(*resources.Nic)
	ret1, _ := ret[1].(*resources.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// AttachNic indicates an expected call of AttachNic.
func (mr *MockLoadbalancersServiceMockRecorder) AttachNic(datacenterId, loadbalancerId, nicId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AttachNic", reflect.TypeOf((*MockLoadbalancersService)(nil).AttachNic), datacenterId, loadbalancerId, nicId)
}

// Create mocks base method.
func (m *MockLoadbalancersService) Create(datacenterId, name string, dhcp bool) (*resources.Loadbalancer, *resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", datacenterId, name, dhcp)
	ret0, _ := ret[0].(*resources.Loadbalancer)
	ret1, _ := ret[1].(*resources.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Create indicates an expected call of Create.
func (mr *MockLoadbalancersServiceMockRecorder) Create(datacenterId, name, dhcp interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockLoadbalancersService)(nil).Create), datacenterId, name, dhcp)
}

// Delete mocks base method.
func (m *MockLoadbalancersService) Delete(datacenterId, loadbalancerId string) (*resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", datacenterId, loadbalancerId)
	ret0, _ := ret[0].(*resources.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Delete indicates an expected call of Delete.
func (mr *MockLoadbalancersServiceMockRecorder) Delete(datacenterId, loadbalancerId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockLoadbalancersService)(nil).Delete), datacenterId, loadbalancerId)
}

// DetachNic mocks base method.
func (m *MockLoadbalancersService) DetachNic(datacenterId, loadbalancerId, nicId string) (*resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DetachNic", datacenterId, loadbalancerId, nicId)
	ret0, _ := ret[0].(*resources.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DetachNic indicates an expected call of DetachNic.
func (mr *MockLoadbalancersServiceMockRecorder) DetachNic(datacenterId, loadbalancerId, nicId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DetachNic", reflect.TypeOf((*MockLoadbalancersService)(nil).DetachNic), datacenterId, loadbalancerId, nicId)
}

// Get mocks base method.
func (m *MockLoadbalancersService) Get(datacenterId, loadbalancerId string) (*resources.Loadbalancer, *resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", datacenterId, loadbalancerId)
	ret0, _ := ret[0].(*resources.Loadbalancer)
	ret1, _ := ret[1].(*resources.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Get indicates an expected call of Get.
func (mr *MockLoadbalancersServiceMockRecorder) Get(datacenterId, loadbalancerId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockLoadbalancersService)(nil).Get), datacenterId, loadbalancerId)
}

// GetNic mocks base method.
func (m *MockLoadbalancersService) GetNic(datacenterId, loadbalancerId, nicId string) (*resources.Nic, *resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNic", datacenterId, loadbalancerId, nicId)
	ret0, _ := ret[0].(*resources.Nic)
	ret1, _ := ret[1].(*resources.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetNic indicates an expected call of GetNic.
func (mr *MockLoadbalancersServiceMockRecorder) GetNic(datacenterId, loadbalancerId, nicId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNic", reflect.TypeOf((*MockLoadbalancersService)(nil).GetNic), datacenterId, loadbalancerId, nicId)
}

// List mocks base method.
func (m *MockLoadbalancersService) List(datacenterId string, params resources.ListQueryParams) (resources.Loadbalancers, *resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", datacenterId, params)
	ret0, _ := ret[0].(resources.Loadbalancers)
	ret1, _ := ret[1].(*resources.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// List indicates an expected call of List.
func (mr *MockLoadbalancersServiceMockRecorder) List(datacenterId, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockLoadbalancersService)(nil).List), datacenterId, params)
}

// ListNics mocks base method.
func (m *MockLoadbalancersService) ListNics(datacenterId, loadbalancerId string, params resources.ListQueryParams) (resources.BalancedNics, *resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListNics", datacenterId, loadbalancerId, params)
	ret0, _ := ret[0].(resources.BalancedNics)
	ret1, _ := ret[1].(*resources.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ListNics indicates an expected call of ListNics.
func (mr *MockLoadbalancersServiceMockRecorder) ListNics(datacenterId, loadbalancerId, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListNics", reflect.TypeOf((*MockLoadbalancersService)(nil).ListNics), datacenterId, loadbalancerId, params)
}

// Update mocks base method.
func (m *MockLoadbalancersService) Update(datacenterId, loadbalancerId string, input resources.LoadbalancerProperties) (*resources.Loadbalancer, *resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", datacenterId, loadbalancerId, input)
	ret0, _ := ret[0].(*resources.Loadbalancer)
	ret1, _ := ret[1].(*resources.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Update indicates an expected call of Update.
func (mr *MockLoadbalancersServiceMockRecorder) Update(datacenterId, loadbalancerId, input interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockLoadbalancersService)(nil).Update), datacenterId, loadbalancerId, input)
}