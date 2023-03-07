// Code generated by MockGen. DO NOT EDIT.
// Source: label.go

// Package mock_resources is a generated GoMock package.
package mock_resources

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	resources "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
)

// MockLabelResourcesService is a mock of LabelResourcesService interface.
type MockLabelResourcesService struct {
	ctrl     *gomock.Controller
	recorder *MockLabelResourcesServiceMockRecorder
}

// MockLabelResourcesServiceMockRecorder is the mock recorder for MockLabelResourcesService.
type MockLabelResourcesServiceMockRecorder struct {
	mock *MockLabelResourcesService
}

// NewMockLabelResourcesService creates a new mock instance.
func NewMockLabelResourcesService(ctrl *gomock.Controller) *MockLabelResourcesService {
	mock := &MockLabelResourcesService{ctrl: ctrl}
	mock.recorder = &MockLabelResourcesServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLabelResourcesService) EXPECT() *MockLabelResourcesServiceMockRecorder {
	return m.recorder
}

// DatacenterCreate mocks base method.
func (m *MockLabelResourcesService) DatacenterCreate(datacenterId, key, value string) (*resources.LabelResource, *resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DatacenterCreate", datacenterId, key, value)
	ret0, _ := ret[0].(*resources.LabelResource)
	ret1, _ := ret[1].(*resources.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// DatacenterCreate indicates an expected call of DatacenterCreate.
func (mr *MockLabelResourcesServiceMockRecorder) DatacenterCreate(datacenterId, key, value interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DatacenterCreate", reflect.TypeOf((*MockLabelResourcesService)(nil).DatacenterCreate), datacenterId, key, value)
}

// DatacenterDelete mocks base method.
func (m *MockLabelResourcesService) DatacenterDelete(datacenterId, key string) (*resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DatacenterDelete", datacenterId, key)
	ret0, _ := ret[0].(*resources.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DatacenterDelete indicates an expected call of DatacenterDelete.
func (mr *MockLabelResourcesServiceMockRecorder) DatacenterDelete(datacenterId, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DatacenterDelete", reflect.TypeOf((*MockLabelResourcesService)(nil).DatacenterDelete), datacenterId, key)
}

// DatacenterGet mocks base method.
func (m *MockLabelResourcesService) DatacenterGet(datacenterId, key string) (*resources.LabelResource, *resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DatacenterGet", datacenterId, key)
	ret0, _ := ret[0].(*resources.LabelResource)
	ret1, _ := ret[1].(*resources.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// DatacenterGet indicates an expected call of DatacenterGet.
func (mr *MockLabelResourcesServiceMockRecorder) DatacenterGet(datacenterId, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DatacenterGet", reflect.TypeOf((*MockLabelResourcesService)(nil).DatacenterGet), datacenterId, key)
}

// DatacenterList mocks base method.
func (m *MockLabelResourcesService) DatacenterList(params resources.ListQueryParams, datacenterId string) (resources.LabelResources, *resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DatacenterList", params, datacenterId)
	ret0, _ := ret[0].(resources.LabelResources)
	ret1, _ := ret[1].(*resources.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// DatacenterList indicates an expected call of DatacenterList.
func (mr *MockLabelResourcesServiceMockRecorder) DatacenterList(params, datacenterId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DatacenterList", reflect.TypeOf((*MockLabelResourcesService)(nil).DatacenterList), params, datacenterId)
}

// GetByUrn mocks base method.
func (m *MockLabelResourcesService) GetByUrn(labelurn string) (*resources.Label, *resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByUrn", labelurn)
	ret0, _ := ret[0].(*resources.Label)
	ret1, _ := ret[1].(*resources.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetByUrn indicates an expected call of GetByUrn.
func (mr *MockLabelResourcesServiceMockRecorder) GetByUrn(labelurn interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByUrn", reflect.TypeOf((*MockLabelResourcesService)(nil).GetByUrn), labelurn)
}

// IpBlockCreate mocks base method.
func (m *MockLabelResourcesService) IpBlockCreate(ipblockId, key, value string) (*resources.LabelResource, *resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IpBlockCreate", ipblockId, key, value)
	ret0, _ := ret[0].(*resources.LabelResource)
	ret1, _ := ret[1].(*resources.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// IpBlockCreate indicates an expected call of IpBlockCreate.
func (mr *MockLabelResourcesServiceMockRecorder) IpBlockCreate(ipblockId, key, value interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IpBlockCreate", reflect.TypeOf((*MockLabelResourcesService)(nil).IpBlockCreate), ipblockId, key, value)
}

// IpBlockDelete mocks base method.
func (m *MockLabelResourcesService) IpBlockDelete(ipblockId, key string) (*resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IpBlockDelete", ipblockId, key)
	ret0, _ := ret[0].(*resources.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IpBlockDelete indicates an expected call of IpBlockDelete.
func (mr *MockLabelResourcesServiceMockRecorder) IpBlockDelete(ipblockId, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IpBlockDelete", reflect.TypeOf((*MockLabelResourcesService)(nil).IpBlockDelete), ipblockId, key)
}

// IpBlockGet mocks base method.
func (m *MockLabelResourcesService) IpBlockGet(ipblockId, key string) (*resources.LabelResource, *resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IpBlockGet", ipblockId, key)
	ret0, _ := ret[0].(*resources.LabelResource)
	ret1, _ := ret[1].(*resources.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// IpBlockGet indicates an expected call of IpBlockGet.
func (mr *MockLabelResourcesServiceMockRecorder) IpBlockGet(ipblockId, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IpBlockGet", reflect.TypeOf((*MockLabelResourcesService)(nil).IpBlockGet), ipblockId, key)
}

// IpBlockList mocks base method.
func (m *MockLabelResourcesService) IpBlockList(params resources.ListQueryParams, ipblockId string) (resources.LabelResources, *resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IpBlockList", params, ipblockId)
	ret0, _ := ret[0].(resources.LabelResources)
	ret1, _ := ret[1].(*resources.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// IpBlockList indicates an expected call of IpBlockList.
func (mr *MockLabelResourcesServiceMockRecorder) IpBlockList(params, ipblockId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IpBlockList", reflect.TypeOf((*MockLabelResourcesService)(nil).IpBlockList), params, ipblockId)
}

// List mocks base method.
func (m *MockLabelResourcesService) List(params resources.ListQueryParams) (resources.Labels, *resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", params)
	ret0, _ := ret[0].(resources.Labels)
	ret1, _ := ret[1].(*resources.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// List indicates an expected call of List.
func (mr *MockLabelResourcesServiceMockRecorder) List(params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockLabelResourcesService)(nil).List), params)
}

// ServerCreate mocks base method.
func (m *MockLabelResourcesService) ServerCreate(datacenterId, serverId, key, value string) (*resources.LabelResource, *resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ServerCreate", datacenterId, serverId, key, value)
	ret0, _ := ret[0].(*resources.LabelResource)
	ret1, _ := ret[1].(*resources.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ServerCreate indicates an expected call of ServerCreate.
func (mr *MockLabelResourcesServiceMockRecorder) ServerCreate(datacenterId, serverId, key, value interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ServerCreate", reflect.TypeOf((*MockLabelResourcesService)(nil).ServerCreate), datacenterId, serverId, key, value)
}

// ServerDelete mocks base method.
func (m *MockLabelResourcesService) ServerDelete(datacenterId, serverId, key string) (*resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ServerDelete", datacenterId, serverId, key)
	ret0, _ := ret[0].(*resources.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ServerDelete indicates an expected call of ServerDelete.
func (mr *MockLabelResourcesServiceMockRecorder) ServerDelete(datacenterId, serverId, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ServerDelete", reflect.TypeOf((*MockLabelResourcesService)(nil).ServerDelete), datacenterId, serverId, key)
}

// ServerGet mocks base method.
func (m *MockLabelResourcesService) ServerGet(datacenterId, serverId, key string) (*resources.LabelResource, *resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ServerGet", datacenterId, serverId, key)
	ret0, _ := ret[0].(*resources.LabelResource)
	ret1, _ := ret[1].(*resources.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ServerGet indicates an expected call of ServerGet.
func (mr *MockLabelResourcesServiceMockRecorder) ServerGet(datacenterId, serverId, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ServerGet", reflect.TypeOf((*MockLabelResourcesService)(nil).ServerGet), datacenterId, serverId, key)
}

// ServerList mocks base method.
func (m *MockLabelResourcesService) ServerList(params resources.ListQueryParams, datacenterId, serverId string) (resources.LabelResources, *resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ServerList", params, datacenterId, serverId)
	ret0, _ := ret[0].(resources.LabelResources)
	ret1, _ := ret[1].(*resources.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ServerList indicates an expected call of ServerList.
func (mr *MockLabelResourcesServiceMockRecorder) ServerList(params, datacenterId, serverId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ServerList", reflect.TypeOf((*MockLabelResourcesService)(nil).ServerList), params, datacenterId, serverId)
}

// SnapshotCreate mocks base method.
func (m *MockLabelResourcesService) SnapshotCreate(snapshotId, key, value string) (*resources.LabelResource, *resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SnapshotCreate", snapshotId, key, value)
	ret0, _ := ret[0].(*resources.LabelResource)
	ret1, _ := ret[1].(*resources.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// SnapshotCreate indicates an expected call of SnapshotCreate.
func (mr *MockLabelResourcesServiceMockRecorder) SnapshotCreate(snapshotId, key, value interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SnapshotCreate", reflect.TypeOf((*MockLabelResourcesService)(nil).SnapshotCreate), snapshotId, key, value)
}

// SnapshotDelete mocks base method.
func (m *MockLabelResourcesService) SnapshotDelete(snapshotId, key string) (*resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SnapshotDelete", snapshotId, key)
	ret0, _ := ret[0].(*resources.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SnapshotDelete indicates an expected call of SnapshotDelete.
func (mr *MockLabelResourcesServiceMockRecorder) SnapshotDelete(snapshotId, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SnapshotDelete", reflect.TypeOf((*MockLabelResourcesService)(nil).SnapshotDelete), snapshotId, key)
}

// SnapshotGet mocks base method.
func (m *MockLabelResourcesService) SnapshotGet(snapshotId, key string) (*resources.LabelResource, *resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SnapshotGet", snapshotId, key)
	ret0, _ := ret[0].(*resources.LabelResource)
	ret1, _ := ret[1].(*resources.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// SnapshotGet indicates an expected call of SnapshotGet.
func (mr *MockLabelResourcesServiceMockRecorder) SnapshotGet(snapshotId, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SnapshotGet", reflect.TypeOf((*MockLabelResourcesService)(nil).SnapshotGet), snapshotId, key)
}

// SnapshotList mocks base method.
func (m *MockLabelResourcesService) SnapshotList(params resources.ListQueryParams, snapshotId string) (resources.LabelResources, *resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SnapshotList", params, snapshotId)
	ret0, _ := ret[0].(resources.LabelResources)
	ret1, _ := ret[1].(*resources.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// SnapshotList indicates an expected call of SnapshotList.
func (mr *MockLabelResourcesServiceMockRecorder) SnapshotList(params, snapshotId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SnapshotList", reflect.TypeOf((*MockLabelResourcesService)(nil).SnapshotList), params, snapshotId)
}

// VolumeCreate mocks base method.
func (m *MockLabelResourcesService) VolumeCreate(datacenterId, serverId, key, value string) (*resources.LabelResource, *resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VolumeCreate", datacenterId, serverId, key, value)
	ret0, _ := ret[0].(*resources.LabelResource)
	ret1, _ := ret[1].(*resources.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// VolumeCreate indicates an expected call of VolumeCreate.
func (mr *MockLabelResourcesServiceMockRecorder) VolumeCreate(datacenterId, serverId, key, value interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VolumeCreate", reflect.TypeOf((*MockLabelResourcesService)(nil).VolumeCreate), datacenterId, serverId, key, value)
}

// VolumeDelete mocks base method.
func (m *MockLabelResourcesService) VolumeDelete(datacenterId, serverId, key string) (*resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VolumeDelete", datacenterId, serverId, key)
	ret0, _ := ret[0].(*resources.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// VolumeDelete indicates an expected call of VolumeDelete.
func (mr *MockLabelResourcesServiceMockRecorder) VolumeDelete(datacenterId, serverId, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VolumeDelete", reflect.TypeOf((*MockLabelResourcesService)(nil).VolumeDelete), datacenterId, serverId, key)
}

// VolumeGet mocks base method.
func (m *MockLabelResourcesService) VolumeGet(datacenterId, serverId, key string) (*resources.LabelResource, *resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VolumeGet", datacenterId, serverId, key)
	ret0, _ := ret[0].(*resources.LabelResource)
	ret1, _ := ret[1].(*resources.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// VolumeGet indicates an expected call of VolumeGet.
func (mr *MockLabelResourcesServiceMockRecorder) VolumeGet(datacenterId, serverId, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VolumeGet", reflect.TypeOf((*MockLabelResourcesService)(nil).VolumeGet), datacenterId, serverId, key)
}

// VolumeList mocks base method.
func (m *MockLabelResourcesService) VolumeList(params resources.ListQueryParams, datacenterId, serverId string) (resources.LabelResources, *resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VolumeList", params, datacenterId, serverId)
	ret0, _ := ret[0].(resources.LabelResources)
	ret1, _ := ret[1].(*resources.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// VolumeList indicates an expected call of VolumeList.
func (mr *MockLabelResourcesServiceMockRecorder) VolumeList(params, datacenterId, serverId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VolumeList", reflect.TypeOf((*MockLabelResourcesService)(nil).VolumeList), params, datacenterId, serverId)
}
