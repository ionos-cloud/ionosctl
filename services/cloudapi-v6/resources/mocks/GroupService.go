// Code generated by MockGen. DO NOT EDIT.
// Source: ./group.go

// Package mock_resources is a generated GoMock package.
package mock_resources

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	resources "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
)

// MockGroupsService is a mock of GroupsService interface.
type MockGroupsService struct {
	ctrl     *gomock.Controller
	recorder *MockGroupsServiceMockRecorder
}

// MockGroupsServiceMockRecorder is the mock recorder for MockGroupsService.
type MockGroupsServiceMockRecorder struct {
	mock *MockGroupsService
}

// NewMockGroupsService creates a new mock instance.
func NewMockGroupsService(ctrl *gomock.Controller) *MockGroupsService {
	mock := &MockGroupsService{ctrl: ctrl}
	mock.recorder = &MockGroupsServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGroupsService) EXPECT() *MockGroupsServiceMockRecorder {
	return m.recorder
}

// AddShare mocks base method.
func (m *MockGroupsService) AddShare(groupId, resourceId string, input resources.GroupShare, params resources.QueryParams) (*resources.GroupShare, *resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddShare", groupId, resourceId, input, params)
	ret0, _ := ret[0].(*resources.GroupShare)
	ret1, _ := ret[1].(*resources.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// AddShare indicates an expected call of AddShare.
func (mr *MockGroupsServiceMockRecorder) AddShare(groupId, resourceId, input, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddShare", reflect.TypeOf((*MockGroupsService)(nil).AddShare), groupId, resourceId, input, params)
}

// AddUser mocks base method.
func (m *MockGroupsService) AddUser(groupId string, input ionoscloud.UserGroupPost, params resources.QueryParams) (*resources.User, *resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddUser", groupId, input, params)
	ret0, _ := ret[0].(*resources.User)
	ret1, _ := ret[1].(*resources.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// AddUser indicates an expected call of AddUser.
func (mr *MockGroupsServiceMockRecorder) AddUser(groupId, input, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddUser", reflect.TypeOf((*MockGroupsService)(nil).AddUser), groupId, input, params)
}

// Create mocks base method.
func (m *MockGroupsService) Create(u resources.Group, params resources.QueryParams) (*resources.Group, *resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", u, params)
	ret0, _ := ret[0].(*resources.Group)
	ret1, _ := ret[1].(*resources.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Create indicates an expected call of Create.
func (mr *MockGroupsServiceMockRecorder) Create(u, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockGroupsService)(nil).Create), u, params)
}

// Delete mocks base method.
func (m *MockGroupsService) Delete(groupId string, params resources.QueryParams) (*resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", groupId, params)
	ret0, _ := ret[0].(*resources.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Delete indicates an expected call of Delete.
func (mr *MockGroupsServiceMockRecorder) Delete(groupId, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockGroupsService)(nil).Delete), groupId, params)
}

// Get mocks base method.
func (m *MockGroupsService) Get(groupId string, params resources.QueryParams) (*resources.Group, *resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", groupId, params)
	ret0, _ := ret[0].(*resources.Group)
	ret1, _ := ret[1].(*resources.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Get indicates an expected call of Get.
func (mr *MockGroupsServiceMockRecorder) Get(groupId, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockGroupsService)(nil).Get), groupId, params)
}

// GetShare mocks base method.
func (m *MockGroupsService) GetShare(groupId, resourceId string, params resources.QueryParams) (*resources.GroupShare, *resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetShare", groupId, resourceId, params)
	ret0, _ := ret[0].(*resources.GroupShare)
	ret1, _ := ret[1].(*resources.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetShare indicates an expected call of GetShare.
func (mr *MockGroupsServiceMockRecorder) GetShare(groupId, resourceId, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetShare", reflect.TypeOf((*MockGroupsService)(nil).GetShare), groupId, resourceId, params)
}

// List mocks base method.
func (m *MockGroupsService) List(params resources.ListQueryParams) (resources.Groups, *resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", params)
	ret0, _ := ret[0].(resources.Groups)
	ret1, _ := ret[1].(*resources.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// List indicates an expected call of List.
func (mr *MockGroupsServiceMockRecorder) List(params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockGroupsService)(nil).List), params)
}

// ListResources mocks base method.
func (m *MockGroupsService) ListResources(groupId string, params resources.ListQueryParams) (resources.ResourceGroups, *resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListResources", groupId, params)
	ret0, _ := ret[0].(resources.ResourceGroups)
	ret1, _ := ret[1].(*resources.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ListResources indicates an expected call of ListResources.
func (mr *MockGroupsServiceMockRecorder) ListResources(groupId, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListResources", reflect.TypeOf((*MockGroupsService)(nil).ListResources), groupId, params)
}

// ListShares mocks base method.
func (m *MockGroupsService) ListShares(groupId string, params resources.ListQueryParams) (resources.GroupShares, *resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListShares", groupId, params)
	ret0, _ := ret[0].(resources.GroupShares)
	ret1, _ := ret[1].(*resources.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ListShares indicates an expected call of ListShares.
func (mr *MockGroupsServiceMockRecorder) ListShares(groupId, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListShares", reflect.TypeOf((*MockGroupsService)(nil).ListShares), groupId, params)
}

// ListUsers mocks base method.
func (m *MockGroupsService) ListUsers(groupId string, params resources.ListQueryParams) (resources.GroupMembers, *resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListUsers", groupId, params)
	ret0, _ := ret[0].(resources.GroupMembers)
	ret1, _ := ret[1].(*resources.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ListUsers indicates an expected call of ListUsers.
func (mr *MockGroupsServiceMockRecorder) ListUsers(groupId, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListUsers", reflect.TypeOf((*MockGroupsService)(nil).ListUsers), groupId, params)
}

// RemoveShare mocks base method.
func (m *MockGroupsService) RemoveShare(groupId, resourceId string, params resources.QueryParams) (*resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveShare", groupId, resourceId, params)
	ret0, _ := ret[0].(*resources.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RemoveShare indicates an expected call of RemoveShare.
func (mr *MockGroupsServiceMockRecorder) RemoveShare(groupId, resourceId, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveShare", reflect.TypeOf((*MockGroupsService)(nil).RemoveShare), groupId, resourceId, params)
}

// RemoveUser mocks base method.
func (m *MockGroupsService) RemoveUser(groupId, userId string, params resources.QueryParams) (*resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveUser", groupId, userId, params)
	ret0, _ := ret[0].(*resources.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RemoveUser indicates an expected call of RemoveUser.
func (mr *MockGroupsServiceMockRecorder) RemoveUser(groupId, userId, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveUser", reflect.TypeOf((*MockGroupsService)(nil).RemoveUser), groupId, userId, params)
}

// Update mocks base method.
func (m *MockGroupsService) Update(groupId string, input resources.Group, params resources.QueryParams) (*resources.Group, *resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", groupId, input, params)
	ret0, _ := ret[0].(*resources.Group)
	ret1, _ := ret[1].(*resources.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Update indicates an expected call of Update.
func (mr *MockGroupsServiceMockRecorder) Update(groupId, input, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockGroupsService)(nil).Update), groupId, input, params)
}

// UpdateShare mocks base method.
func (m *MockGroupsService) UpdateShare(groupId, resourceId string, input resources.GroupShare, params resources.QueryParams) (*resources.GroupShare, *resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateShare", groupId, resourceId, input, params)
	ret0, _ := ret[0].(*resources.GroupShare)
	ret1, _ := ret[1].(*resources.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// UpdateShare indicates an expected call of UpdateShare.
func (mr *MockGroupsServiceMockRecorder) UpdateShare(groupId, resourceId, input, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateShare", reflect.TypeOf((*MockGroupsService)(nil).UpdateShare), groupId, resourceId, input, params)
}
