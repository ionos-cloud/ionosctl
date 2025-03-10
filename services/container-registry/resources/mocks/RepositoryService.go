// Code generated by MockGen. DO NOT EDIT.
// Source: ./repository.go

// Package mock_resources is a generated GoMock package.
package mock_resources

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	shared "github.com/ionos-cloud/sdk-go-bundle/shared"
)

// MockRepositoryService is a mock of RepositoryService interface.
type MockRepositoryService struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryServiceMockRecorder
}

// MockRepositoryServiceMockRecorder is the mock recorder for MockRepositoryService.
type MockRepositoryServiceMockRecorder struct {
	mock *MockRepositoryService
}

// NewMockRepositoryService creates a new mock instance.
func NewMockRepositoryService(ctrl *gomock.Controller) *MockRepositoryService {
	mock := &MockRepositoryService{ctrl: ctrl}
	mock.recorder = &MockRepositoryServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepositoryService) EXPECT() *MockRepositoryServiceMockRecorder {
	return m.recorder
}

// Delete mocks base method.
func (m *MockRepositoryService) Delete(regId, name string) (*shared.APIResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", regId, name)
	ret0, _ := ret[0].(*shared.APIResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Delete indicates an expected call of Delete.
func (mr *MockRepositoryServiceMockRecorder) Delete(regId, name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockRepositoryService)(nil).Delete), regId, name)
}
