// Code generated by MockGen. DO NOT EDIT.
// Source: natgateway.go

// Package mock_resources is a generated GoMock package.
package mock_resources

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	resources "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
)

// MockNatGatewaysService is a mock of NatGatewaysService interface.
type MockNatGatewaysService struct {
	ctrl     *gomock.Controller
	recorder *MockNatGatewaysServiceMockRecorder
}

// MockNatGatewaysServiceMockRecorder is the mock recorder for MockNatGatewaysService.
type MockNatGatewaysServiceMockRecorder struct {
	mock *MockNatGatewaysService
}

// NewMockNatGatewaysService creates a new mock instance.
func NewMockNatGatewaysService(ctrl *gomock.Controller) *MockNatGatewaysService {
	mock := &MockNatGatewaysService{ctrl: ctrl}
	mock.recorder = &MockNatGatewaysServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockNatGatewaysService) EXPECT() *MockNatGatewaysServiceMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockNatGatewaysService) Create(datacenterId string, input resources.NatGateway, params resources.QueryParams) (*resources.NatGateway, *resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", datacenterId, input, params)
	ret0, _ := ret[0].(*resources.NatGateway)
	ret1, _ := ret[1].(*resources.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Create indicates an expected call of Create.
func (mr *MockNatGatewaysServiceMockRecorder) Create(datacenterId, input, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockNatGatewaysService)(nil).Create), datacenterId, input, params)
}

// CreateFlowLog mocks base method.
func (m *MockNatGatewaysService) CreateFlowLog(datacenterId, natGatewayId string, input resources.FlowLog, params resources.QueryParams) (*resources.FlowLog, *resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateFlowLog", datacenterId, natGatewayId, input, params)
	ret0, _ := ret[0].(*resources.FlowLog)
	ret1, _ := ret[1].(*resources.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// CreateFlowLog indicates an expected call of CreateFlowLog.
func (mr *MockNatGatewaysServiceMockRecorder) CreateFlowLog(datacenterId, natGatewayId, input, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateFlowLog", reflect.TypeOf((*MockNatGatewaysService)(nil).CreateFlowLog), datacenterId, natGatewayId, input, params)
}

// CreateRule mocks base method.
func (m *MockNatGatewaysService) CreateRule(datacenterId, natGatewayId string, input resources.NatGatewayRule, params resources.QueryParams) (*resources.NatGatewayRule, *resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateRule", datacenterId, natGatewayId, input, params)
	ret0, _ := ret[0].(*resources.NatGatewayRule)
	ret1, _ := ret[1].(*resources.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// CreateRule indicates an expected call of CreateRule.
func (mr *MockNatGatewaysServiceMockRecorder) CreateRule(datacenterId, natGatewayId, input, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateRule", reflect.TypeOf((*MockNatGatewaysService)(nil).CreateRule), datacenterId, natGatewayId, input, params)
}

// Delete mocks base method.
func (m *MockNatGatewaysService) Delete(datacenterId, natGatewayId string, params resources.QueryParams) (*resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", datacenterId, natGatewayId, params)
	ret0, _ := ret[0].(*resources.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Delete indicates an expected call of Delete.
func (mr *MockNatGatewaysServiceMockRecorder) Delete(datacenterId, natGatewayId, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockNatGatewaysService)(nil).Delete), datacenterId, natGatewayId, params)
}

// DeleteFlowLog mocks base method.
func (m *MockNatGatewaysService) DeleteFlowLog(datacenterId, natGatewayId, flowlogId string, params resources.QueryParams) (*resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteFlowLog", datacenterId, natGatewayId, flowlogId, params)
	ret0, _ := ret[0].(*resources.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteFlowLog indicates an expected call of DeleteFlowLog.
func (mr *MockNatGatewaysServiceMockRecorder) DeleteFlowLog(datacenterId, natGatewayId, flowlogId, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteFlowLog", reflect.TypeOf((*MockNatGatewaysService)(nil).DeleteFlowLog), datacenterId, natGatewayId, flowlogId, params)
}

// DeleteRule mocks base method.
func (m *MockNatGatewaysService) DeleteRule(datacenterId, natGatewayId, ruleId string, params resources.QueryParams) (*resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteRule", datacenterId, natGatewayId, ruleId, params)
	ret0, _ := ret[0].(*resources.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteRule indicates an expected call of DeleteRule.
func (mr *MockNatGatewaysServiceMockRecorder) DeleteRule(datacenterId, natGatewayId, ruleId, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteRule", reflect.TypeOf((*MockNatGatewaysService)(nil).DeleteRule), datacenterId, natGatewayId, ruleId, params)
}

// Get mocks base method.
func (m *MockNatGatewaysService) Get(datacenterId, natGatewayId string, params resources.QueryParams) (*resources.NatGateway, *resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", datacenterId, natGatewayId, params)
	ret0, _ := ret[0].(*resources.NatGateway)
	ret1, _ := ret[1].(*resources.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Get indicates an expected call of Get.
func (mr *MockNatGatewaysServiceMockRecorder) Get(datacenterId, natGatewayId, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockNatGatewaysService)(nil).Get), datacenterId, natGatewayId, params)
}

// GetFlowLog mocks base method.
func (m *MockNatGatewaysService) GetFlowLog(datacenterId, natGatewayId, flowlogId string, params resources.QueryParams) (*resources.FlowLog, *resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFlowLog", datacenterId, natGatewayId, flowlogId, params)
	ret0, _ := ret[0].(*resources.FlowLog)
	ret1, _ := ret[1].(*resources.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetFlowLog indicates an expected call of GetFlowLog.
func (mr *MockNatGatewaysServiceMockRecorder) GetFlowLog(datacenterId, natGatewayId, flowlogId, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFlowLog", reflect.TypeOf((*MockNatGatewaysService)(nil).GetFlowLog), datacenterId, natGatewayId, flowlogId, params)
}

// GetRule mocks base method.
func (m *MockNatGatewaysService) GetRule(datacenterId, natGatewayId, ruleId string, params resources.QueryParams) (*resources.NatGatewayRule, *resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRule", datacenterId, natGatewayId, ruleId, params)
	ret0, _ := ret[0].(*resources.NatGatewayRule)
	ret1, _ := ret[1].(*resources.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetRule indicates an expected call of GetRule.
func (mr *MockNatGatewaysServiceMockRecorder) GetRule(datacenterId, natGatewayId, ruleId, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRule", reflect.TypeOf((*MockNatGatewaysService)(nil).GetRule), datacenterId, natGatewayId, ruleId, params)
}

// List mocks base method.
func (m *MockNatGatewaysService) List(datacenterId string, params resources.ListQueryParams) (resources.NatGateways, *resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", datacenterId, params)
	ret0, _ := ret[0].(resources.NatGateways)
	ret1, _ := ret[1].(*resources.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// List indicates an expected call of List.
func (mr *MockNatGatewaysServiceMockRecorder) List(datacenterId, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockNatGatewaysService)(nil).List), datacenterId, params)
}

// ListFlowLogs mocks base method.
func (m *MockNatGatewaysService) ListFlowLogs(datacenterId, natGatewayId string, params resources.ListQueryParams) (resources.FlowLogs, *resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListFlowLogs", datacenterId, natGatewayId, params)
	ret0, _ := ret[0].(resources.FlowLogs)
	ret1, _ := ret[1].(*resources.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ListFlowLogs indicates an expected call of ListFlowLogs.
func (mr *MockNatGatewaysServiceMockRecorder) ListFlowLogs(datacenterId, natGatewayId, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListFlowLogs", reflect.TypeOf((*MockNatGatewaysService)(nil).ListFlowLogs), datacenterId, natGatewayId, params)
}

// ListRules mocks base method.
func (m *MockNatGatewaysService) ListRules(datacenterId, natGatewayId string, params resources.ListQueryParams) (resources.NatGatewayRules, *resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListRules", datacenterId, natGatewayId, params)
	ret0, _ := ret[0].(resources.NatGatewayRules)
	ret1, _ := ret[1].(*resources.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ListRules indicates an expected call of ListRules.
func (mr *MockNatGatewaysServiceMockRecorder) ListRules(datacenterId, natGatewayId, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListRules", reflect.TypeOf((*MockNatGatewaysService)(nil).ListRules), datacenterId, natGatewayId, params)
}

// Update mocks base method.
func (m *MockNatGatewaysService) Update(datacenterId, natGatewayId string, input resources.NatGatewayProperties, params resources.QueryParams) (*resources.NatGateway, *resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", datacenterId, natGatewayId, input, params)
	ret0, _ := ret[0].(*resources.NatGateway)
	ret1, _ := ret[1].(*resources.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Update indicates an expected call of Update.
func (mr *MockNatGatewaysServiceMockRecorder) Update(datacenterId, natGatewayId, input, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockNatGatewaysService)(nil).Update), datacenterId, natGatewayId, input, params)
}

// UpdateFlowLog mocks base method.
func (m *MockNatGatewaysService) UpdateFlowLog(datacenterId, natGatewayId, flowlogId string, input *resources.FlowLogProperties, params resources.QueryParams) (*resources.FlowLog, *resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateFlowLog", datacenterId, natGatewayId, flowlogId, input, params)
	ret0, _ := ret[0].(*resources.FlowLog)
	ret1, _ := ret[1].(*resources.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// UpdateFlowLog indicates an expected call of UpdateFlowLog.
func (mr *MockNatGatewaysServiceMockRecorder) UpdateFlowLog(datacenterId, natGatewayId, flowlogId, input, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateFlowLog", reflect.TypeOf((*MockNatGatewaysService)(nil).UpdateFlowLog), datacenterId, natGatewayId, flowlogId, input, params)
}

// UpdateRule mocks base method.
func (m *MockNatGatewaysService) UpdateRule(datacenterId, natGatewayId, ruleId string, input resources.NatGatewayRuleProperties, params resources.QueryParams) (*resources.NatGatewayRule, *resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateRule", datacenterId, natGatewayId, ruleId, input, params)
	ret0, _ := ret[0].(*resources.NatGatewayRule)
	ret1, _ := ret[1].(*resources.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// UpdateRule indicates an expected call of UpdateRule.
func (mr *MockNatGatewaysServiceMockRecorder) UpdateRule(datacenterId, natGatewayId, ruleId, input, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateRule", reflect.TypeOf((*MockNatGatewaysService)(nil).UpdateRule), datacenterId, natGatewayId, ruleId, input, params)
}
