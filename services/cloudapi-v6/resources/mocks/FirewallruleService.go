// Code generated by MockGen. DO NOT EDIT.
// Source: ./firewallrule.go

// Package mock_resources is a generated GoMock package.
package mock_resources

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	resources "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
)

// MockFirewallRulesService is a mock of FirewallRulesService interface.
type MockFirewallRulesService struct {
	ctrl     *gomock.Controller
	recorder *MockFirewallRulesServiceMockRecorder
}

// MockFirewallRulesServiceMockRecorder is the mock recorder for MockFirewallRulesService.
type MockFirewallRulesServiceMockRecorder struct {
	mock *MockFirewallRulesService
}

// NewMockFirewallRulesService creates a new mock instance.
func NewMockFirewallRulesService(ctrl *gomock.Controller) *MockFirewallRulesService {
	mock := &MockFirewallRulesService{ctrl: ctrl}
	mock.recorder = &MockFirewallRulesServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFirewallRulesService) EXPECT() *MockFirewallRulesServiceMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockFirewallRulesService) Create(datacenterId, serverId, nicId string, input resources.FirewallRule, params resources.QueryParams) (*resources.FirewallRule, *resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", datacenterId, serverId, nicId, input, params)
	ret0, _ := ret[0].(*resources.FirewallRule)
	ret1, _ := ret[1].(*resources.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Create indicates an expected call of Create.
func (mr *MockFirewallRulesServiceMockRecorder) Create(datacenterId, serverId, nicId, input, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockFirewallRulesService)(nil).Create), datacenterId, serverId, nicId, input, params)
}

// Delete mocks base method.
func (m *MockFirewallRulesService) Delete(datacenterId, serverId, nicId, firewallRuleId string, params resources.QueryParams) (*resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", datacenterId, serverId, nicId, firewallRuleId, params)
	ret0, _ := ret[0].(*resources.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Delete indicates an expected call of Delete.
func (mr *MockFirewallRulesServiceMockRecorder) Delete(datacenterId, serverId, nicId, firewallRuleId, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockFirewallRulesService)(nil).Delete), datacenterId, serverId, nicId, firewallRuleId, params)
}

// Get mocks base method.
func (m *MockFirewallRulesService) Get(datacenterId, serverId, nicId, firewallRuleId string, params resources.QueryParams) (*resources.FirewallRule, *resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", datacenterId, serverId, nicId, firewallRuleId, params)
	ret0, _ := ret[0].(*resources.FirewallRule)
	ret1, _ := ret[1].(*resources.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Get indicates an expected call of Get.
func (mr *MockFirewallRulesServiceMockRecorder) Get(datacenterId, serverId, nicId, firewallRuleId, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockFirewallRulesService)(nil).Get), datacenterId, serverId, nicId, firewallRuleId, params)
}

// List mocks base method.
func (m *MockFirewallRulesService) List(datacenterId, serverId, nicId string, params resources.ListQueryParams) (resources.FirewallRules, *resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", datacenterId, serverId, nicId, params)
	ret0, _ := ret[0].(resources.FirewallRules)
	ret1, _ := ret[1].(*resources.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// List indicates an expected call of List.
func (mr *MockFirewallRulesServiceMockRecorder) List(datacenterId, serverId, nicId, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockFirewallRulesService)(nil).List), datacenterId, serverId, nicId, params)
}

// Update mocks base method.
func (m *MockFirewallRulesService) Update(datacenterId, serverId, nicId, firewallRuleId string, input resources.FirewallRuleProperties, params resources.QueryParams) (*resources.FirewallRule, *resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", datacenterId, serverId, nicId, firewallRuleId, input, params)
	ret0, _ := ret[0].(*resources.FirewallRule)
	ret1, _ := ret[1].(*resources.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Update indicates an expected call of Update.
func (mr *MockFirewallRulesServiceMockRecorder) Update(datacenterId, serverId, nicId, firewallRuleId, input, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockFirewallRulesService)(nil).Update), datacenterId, serverId, nicId, firewallRuleId, input, params)
}