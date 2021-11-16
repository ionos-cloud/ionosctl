// Code generated by MockGen. DO NOT EDIT.
// Source: networkloadbalancer.go

// Package mock_resources is a generated GoMock package.
package mock_resources

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	resources "github.com/ionos-cloud/ionosctl/services/cloudapi-v6/resources"
)

// MockNetworkLoadBalancersService is a mock of NetworkLoadBalancersService interface.
type MockNetworkLoadBalancersService struct {
	ctrl     *gomock.Controller
	recorder *MockNetworkLoadBalancersServiceMockRecorder
}

// MockNetworkLoadBalancersServiceMockRecorder is the mock recorder for MockNetworkLoadBalancersService.
type MockNetworkLoadBalancersServiceMockRecorder struct {
	mock *MockNetworkLoadBalancersService
}

// NewMockNetworkLoadBalancersService creates a new mock instance.
func NewMockNetworkLoadBalancersService(ctrl *gomock.Controller) *MockNetworkLoadBalancersService {
	mock := &MockNetworkLoadBalancersService{ctrl: ctrl}
	mock.recorder = &MockNetworkLoadBalancersServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockNetworkLoadBalancersService) EXPECT() *MockNetworkLoadBalancersServiceMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockNetworkLoadBalancersService) Create(datacenterId string, input resources.NetworkLoadBalancer) (*resources.NetworkLoadBalancer, *resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", datacenterId, input)
	ret0, _ := ret[0].(*resources.NetworkLoadBalancer)
	ret1, _ := ret[1].(*resources.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Create indicates an expected call of Create.
func (mr *MockNetworkLoadBalancersServiceMockRecorder) Create(datacenterId, input interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockNetworkLoadBalancersService)(nil).Create), datacenterId, input)
}

// CreateFlowLog mocks base method.
func (m *MockNetworkLoadBalancersService) CreateFlowLog(datacenterId, networkLoadBalancerId string, input resources.FlowLog) (*resources.FlowLog, *resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateFlowLog", datacenterId, networkLoadBalancerId, input)
	ret0, _ := ret[0].(*resources.FlowLog)
	ret1, _ := ret[1].(*resources.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// CreateFlowLog indicates an expected call of CreateFlowLog.
func (mr *MockNetworkLoadBalancersServiceMockRecorder) CreateFlowLog(datacenterId, networkLoadBalancerId, input interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateFlowLog", reflect.TypeOf((*MockNetworkLoadBalancersService)(nil).CreateFlowLog), datacenterId, networkLoadBalancerId, input)
}

// CreateForwardingRule mocks base method.
func (m *MockNetworkLoadBalancersService) CreateForwardingRule(datacenterId, networkLoadBalancerId string, input resources.NetworkLoadBalancerForwardingRule) (*resources.NetworkLoadBalancerForwardingRule, *resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateForwardingRule", datacenterId, networkLoadBalancerId, input)
	ret0, _ := ret[0].(*resources.NetworkLoadBalancerForwardingRule)
	ret1, _ := ret[1].(*resources.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// CreateForwardingRule indicates an expected call of CreateForwardingRule.
func (mr *MockNetworkLoadBalancersServiceMockRecorder) CreateForwardingRule(datacenterId, networkLoadBalancerId, input interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateForwardingRule", reflect.TypeOf((*MockNetworkLoadBalancersService)(nil).CreateForwardingRule), datacenterId, networkLoadBalancerId, input)
}

// Delete mocks base method.
func (m *MockNetworkLoadBalancersService) Delete(datacenterId, networkLoadBalancerId string) (*resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", datacenterId, networkLoadBalancerId)
	ret0, _ := ret[0].(*resources.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Delete indicates an expected call of Delete.
func (mr *MockNetworkLoadBalancersServiceMockRecorder) Delete(datacenterId, networkLoadBalancerId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockNetworkLoadBalancersService)(nil).Delete), datacenterId, networkLoadBalancerId)
}

// DeleteFlowLog mocks base method.
func (m *MockNetworkLoadBalancersService) DeleteFlowLog(datacenterId, networkLoadBalancerId, flowLogId string) (*resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteFlowLog", datacenterId, networkLoadBalancerId, flowLogId)
	ret0, _ := ret[0].(*resources.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteFlowLog indicates an expected call of DeleteFlowLog.
func (mr *MockNetworkLoadBalancersServiceMockRecorder) DeleteFlowLog(datacenterId, networkLoadBalancerId, flowLogId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteFlowLog", reflect.TypeOf((*MockNetworkLoadBalancersService)(nil).DeleteFlowLog), datacenterId, networkLoadBalancerId, flowLogId)
}

// DeleteForwardingRule mocks base method.
func (m *MockNetworkLoadBalancersService) DeleteForwardingRule(datacenterId, networkLoadBalancerId, forwardingRuleId string) (*resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteForwardingRule", datacenterId, networkLoadBalancerId, forwardingRuleId)
	ret0, _ := ret[0].(*resources.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteForwardingRule indicates an expected call of DeleteForwardingRule.
func (mr *MockNetworkLoadBalancersServiceMockRecorder) DeleteForwardingRule(datacenterId, networkLoadBalancerId, forwardingRuleId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteForwardingRule", reflect.TypeOf((*MockNetworkLoadBalancersService)(nil).DeleteForwardingRule), datacenterId, networkLoadBalancerId, forwardingRuleId)
}

// Get mocks base method.
func (m *MockNetworkLoadBalancersService) Get(datacenterId, networkLoadBalancerId string) (*resources.NetworkLoadBalancer, *resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", datacenterId, networkLoadBalancerId)
	ret0, _ := ret[0].(*resources.NetworkLoadBalancer)
	ret1, _ := ret[1].(*resources.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Get indicates an expected call of Get.
func (mr *MockNetworkLoadBalancersServiceMockRecorder) Get(datacenterId, networkLoadBalancerId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockNetworkLoadBalancersService)(nil).Get), datacenterId, networkLoadBalancerId)
}

// GetFlowLog mocks base method.
func (m *MockNetworkLoadBalancersService) GetFlowLog(datacenterId, networkLoadBalancerId, flowLogId string) (*resources.FlowLog, *resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFlowLog", datacenterId, networkLoadBalancerId, flowLogId)
	ret0, _ := ret[0].(*resources.FlowLog)
	ret1, _ := ret[1].(*resources.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetFlowLog indicates an expected call of GetFlowLog.
func (mr *MockNetworkLoadBalancersServiceMockRecorder) GetFlowLog(datacenterId, networkLoadBalancerId, flowLogId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFlowLog", reflect.TypeOf((*MockNetworkLoadBalancersService)(nil).GetFlowLog), datacenterId, networkLoadBalancerId, flowLogId)
}

// GetForwardingRule mocks base method.
func (m *MockNetworkLoadBalancersService) GetForwardingRule(datacenterId, networkLoadBalancerId, forwardingRuleId string) (*resources.NetworkLoadBalancerForwardingRule, *resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetForwardingRule", datacenterId, networkLoadBalancerId, forwardingRuleId)
	ret0, _ := ret[0].(*resources.NetworkLoadBalancerForwardingRule)
	ret1, _ := ret[1].(*resources.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetForwardingRule indicates an expected call of GetForwardingRule.
func (mr *MockNetworkLoadBalancersServiceMockRecorder) GetForwardingRule(datacenterId, networkLoadBalancerId, forwardingRuleId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetForwardingRule", reflect.TypeOf((*MockNetworkLoadBalancersService)(nil).GetForwardingRule), datacenterId, networkLoadBalancerId, forwardingRuleId)
}

// List mocks base method.
func (m *MockNetworkLoadBalancersService) List(datacenterId string, params resources.ListQueryParams) (resources.NetworkLoadBalancers, *resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", datacenterId, params)
	ret0, _ := ret[0].(resources.NetworkLoadBalancers)
	ret1, _ := ret[1].(*resources.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// List indicates an expected call of List.
func (mr *MockNetworkLoadBalancersServiceMockRecorder) List(datacenterId, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockNetworkLoadBalancersService)(nil).List), datacenterId, params)
}

// ListFlowLogs mocks base method.
func (m *MockNetworkLoadBalancersService) ListFlowLogs(datacenterId, networkLoadBalancerId string, params resources.ListQueryParams) (resources.FlowLogs, *resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListFlowLogs", datacenterId, networkLoadBalancerId, params)
	ret0, _ := ret[0].(resources.FlowLogs)
	ret1, _ := ret[1].(*resources.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ListFlowLogs indicates an expected call of ListFlowLogs.
func (mr *MockNetworkLoadBalancersServiceMockRecorder) ListFlowLogs(datacenterId, networkLoadBalancerId, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListFlowLogs", reflect.TypeOf((*MockNetworkLoadBalancersService)(nil).ListFlowLogs), datacenterId, networkLoadBalancerId, params)
}

// ListForwardingRules mocks base method.
func (m *MockNetworkLoadBalancersService) ListForwardingRules(datacenterId, networkLoadBalancerId string, params resources.ListQueryParams) (resources.NetworkLoadBalancerForwardingRules, *resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListForwardingRules", datacenterId, networkLoadBalancerId, params)
	ret0, _ := ret[0].(resources.NetworkLoadBalancerForwardingRules)
	ret1, _ := ret[1].(*resources.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ListForwardingRules indicates an expected call of ListForwardingRules.
func (mr *MockNetworkLoadBalancersServiceMockRecorder) ListForwardingRules(datacenterId, networkLoadBalancerId, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListForwardingRules", reflect.TypeOf((*MockNetworkLoadBalancersService)(nil).ListForwardingRules), datacenterId, networkLoadBalancerId, params)
}

// Update mocks base method.
func (m *MockNetworkLoadBalancersService) Update(datacenterId, networkLoadBalancerId string, input resources.NetworkLoadBalancerProperties) (*resources.NetworkLoadBalancer, *resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", datacenterId, networkLoadBalancerId, input)
	ret0, _ := ret[0].(*resources.NetworkLoadBalancer)
	ret1, _ := ret[1].(*resources.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Update indicates an expected call of Update.
func (mr *MockNetworkLoadBalancersServiceMockRecorder) Update(datacenterId, networkLoadBalancerId, input interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockNetworkLoadBalancersService)(nil).Update), datacenterId, networkLoadBalancerId, input)
}

// UpdateFlowLog mocks base method.
func (m *MockNetworkLoadBalancersService) UpdateFlowLog(datacenterId, networkLoadBalancerId, flowLogId string, input *resources.FlowLogProperties) (*resources.FlowLog, *resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateFlowLog", datacenterId, networkLoadBalancerId, flowLogId, input)
	ret0, _ := ret[0].(*resources.FlowLog)
	ret1, _ := ret[1].(*resources.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// UpdateFlowLog indicates an expected call of UpdateFlowLog.
func (mr *MockNetworkLoadBalancersServiceMockRecorder) UpdateFlowLog(datacenterId, networkLoadBalancerId, flowLogId, input interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateFlowLog", reflect.TypeOf((*MockNetworkLoadBalancersService)(nil).UpdateFlowLog), datacenterId, networkLoadBalancerId, flowLogId, input)
}

// UpdateForwardingRule mocks base method.
func (m *MockNetworkLoadBalancersService) UpdateForwardingRule(datacenterId, networkLoadBalancerId, forwardingRuleId string, input *resources.NetworkLoadBalancerForwardingRuleProperties) (*resources.NetworkLoadBalancerForwardingRule, *resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateForwardingRule", datacenterId, networkLoadBalancerId, forwardingRuleId, input)
	ret0, _ := ret[0].(*resources.NetworkLoadBalancerForwardingRule)
	ret1, _ := ret[1].(*resources.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// UpdateForwardingRule indicates an expected call of UpdateForwardingRule.
func (mr *MockNetworkLoadBalancersServiceMockRecorder) UpdateForwardingRule(datacenterId, networkLoadBalancerId, forwardingRuleId, input interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateForwardingRule", reflect.TypeOf((*MockNetworkLoadBalancersService)(nil).UpdateForwardingRule), datacenterId, networkLoadBalancerId, forwardingRuleId, input)
}
