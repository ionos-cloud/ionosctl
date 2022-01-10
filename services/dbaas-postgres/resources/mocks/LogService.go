// Code generated by MockGen. DO NOT EDIT.
// Source: logs.go

// Package mock_resources is a generated GoMock package.
package mock_resources

import (
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
	resources "github.com/ionos-cloud/ionosctl/services/dbaas-postgres/resources"
)

// MockLogsService is a mock of LogsService interface.
type MockLogsService struct {
	ctrl     *gomock.Controller
	recorder *MockLogsServiceMockRecorder
}

// MockLogsServiceMockRecorder is the mock recorder for MockLogsService.
type MockLogsServiceMockRecorder struct {
	mock *MockLogsService
}

// NewMockLogsService creates a new mock instance.
func NewMockLogsService(ctrl *gomock.Controller) *MockLogsService {
	mock := &MockLogsService{ctrl: ctrl}
	mock.recorder = &MockLogsServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLogsService) EXPECT() *MockLogsServiceMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockLogsService) Get(clusterId string, limit int32, startTime, endTime time.Time) (*resources.ClusterLogs, *resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", clusterId, limit, startTime, endTime)
	ret0, _ := ret[0].(*resources.ClusterLogs)
	ret1, _ := ret[1].(*resources.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Get indicates an expected call of Get.
func (mr *MockLogsServiceMockRecorder) Get(clusterId, limit, startTime, endTime interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockLogsService)(nil).Get), clusterId, limit, startTime, endTime)
}