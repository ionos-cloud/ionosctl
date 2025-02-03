// Code generated by MockGen. DO NOT EDIT.
// Source: ./token.go

// Package mock_resources is a generated GoMock package.
package mock_resources

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	resources "github.com/ionos-cloud/ionosctl/v6/services/auth-v1/resources"
)

// MockTokensService is a mock of TokensService interface.
type MockTokensService struct {
	ctrl     *gomock.Controller
	recorder *MockTokensServiceMockRecorder
}

// MockTokensServiceMockRecorder is the mock recorder for MockTokensService.
type MockTokensServiceMockRecorder struct {
	mock *MockTokensService
}

// NewMockTokensService creates a new mock instance.
func NewMockTokensService(ctrl *gomock.Controller) *MockTokensService {
	mock := &MockTokensService{ctrl: ctrl}
	mock.recorder = &MockTokensServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTokensService) EXPECT() *MockTokensServiceMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockTokensService) Create(contractNumber, ttl int32) (*resources.Jwt, *resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", contractNumber, ttl)
	ret0, _ := ret[0].(*resources.Jwt)
	ret1, _ := ret[1].(*resources.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Create indicates an expected call of Create.
func (mr *MockTokensServiceMockRecorder) Create(contractNumber, ttl interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockTokensService)(nil).Create), contractNumber, ttl)
}

// DeleteByCriteria mocks base method.
func (m *MockTokensService) DeleteByCriteria(criteria string, contractNumber int32) (*resources.DeleteResponse, *resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteByCriteria", criteria, contractNumber)
	ret0, _ := ret[0].(*resources.DeleteResponse)
	ret1, _ := ret[1].(*resources.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// DeleteByCriteria indicates an expected call of DeleteByCriteria.
func (mr *MockTokensServiceMockRecorder) DeleteByCriteria(criteria, contractNumber interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteByCriteria", reflect.TypeOf((*MockTokensService)(nil).DeleteByCriteria), criteria, contractNumber)
}

// DeleteByID mocks base method.
func (m *MockTokensService) DeleteByID(tokenId string, contractNumber int32) (*resources.DeleteResponse, *resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteByID", tokenId, contractNumber)
	ret0, _ := ret[0].(*resources.DeleteResponse)
	ret1, _ := ret[1].(*resources.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// DeleteByID indicates an expected call of DeleteByID.
func (mr *MockTokensServiceMockRecorder) DeleteByID(tokenId, contractNumber interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteByID", reflect.TypeOf((*MockTokensService)(nil).DeleteByID), tokenId, contractNumber)
}

// Get mocks base method.
func (m *MockTokensService) Get(tokenId string, contractNumber int32) (*resources.Token, *resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", tokenId, contractNumber)
	ret0, _ := ret[0].(*resources.Token)
	ret1, _ := ret[1].(*resources.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Get indicates an expected call of Get.
func (mr *MockTokensServiceMockRecorder) Get(tokenId, contractNumber interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockTokensService)(nil).Get), tokenId, contractNumber)
}

// List mocks base method.
func (m *MockTokensService) List(contractNumber int32) (resources.Tokens, *resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", contractNumber)
	ret0, _ := ret[0].(resources.Tokens)
	ret1, _ := ret[1].(*resources.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// List indicates an expected call of List.
func (mr *MockTokensServiceMockRecorder) List(contractNumber interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockTokensService)(nil).List), contractNumber)
}
