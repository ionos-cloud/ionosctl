// Code generated by MockGen. DO NOT EDIT.
// Source: k8s.go

// Package mock_resources is a generated GoMock package.
package mock_resources

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	resources "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
)

// MockK8sService is a mock of K8sService interface.
type MockK8sService struct {
	ctrl     *gomock.Controller
	recorder *MockK8sServiceMockRecorder
}

// MockK8sServiceMockRecorder is the mock recorder for MockK8sService.
type MockK8sServiceMockRecorder struct {
	mock *MockK8sService
}

// NewMockK8sService creates a new mock instance.
func NewMockK8sService(ctrl *gomock.Controller) *MockK8sService {
	mock := &MockK8sService{ctrl: ctrl}
	mock.recorder = &MockK8sServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockK8sService) EXPECT() *MockK8sServiceMockRecorder {
	return m.recorder
}

// CreateCluster mocks base method.
func (m *MockK8sService) CreateCluster(u resources.K8sClusterForPost, params resources.QueryParams) (*resources.K8sCluster, *resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateCluster", u, params)
	ret0, _ := ret[0].(*resources.K8sCluster)
	ret1, _ := ret[1].(*resources.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// CreateCluster indicates an expected call of CreateCluster.
func (mr *MockK8sServiceMockRecorder) CreateCluster(u, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateCluster", reflect.TypeOf((*MockK8sService)(nil).CreateCluster), u, params)
}

// CreateNodePool mocks base method.
func (m *MockK8sService) CreateNodePool(clusterId string, nodepool resources.K8sNodePoolForPost, params resources.QueryParams) (*resources.K8sNodePool, *resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateNodePool", clusterId, nodepool, params)
	ret0, _ := ret[0].(*resources.K8sNodePool)
	ret1, _ := ret[1].(*resources.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// CreateNodePool indicates an expected call of CreateNodePool.
func (mr *MockK8sServiceMockRecorder) CreateNodePool(clusterId, nodepool, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateNodePool", reflect.TypeOf((*MockK8sService)(nil).CreateNodePool), clusterId, nodepool, params)
}

// DeleteCluster mocks base method.
func (m *MockK8sService) DeleteCluster(clusterId string, params resources.QueryParams) (*resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteCluster", clusterId, params)
	ret0, _ := ret[0].(*resources.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteCluster indicates an expected call of DeleteCluster.
func (mr *MockK8sServiceMockRecorder) DeleteCluster(clusterId, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteCluster", reflect.TypeOf((*MockK8sService)(nil).DeleteCluster), clusterId, params)
}

// DeleteNode mocks base method.
func (m *MockK8sService) DeleteNode(clusterId, nodepoolId, nodeId string, params resources.QueryParams) (*resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteNode", clusterId, nodepoolId, nodeId, params)
	ret0, _ := ret[0].(*resources.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteNode indicates an expected call of DeleteNode.
func (mr *MockK8sServiceMockRecorder) DeleteNode(clusterId, nodepoolId, nodeId, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteNode", reflect.TypeOf((*MockK8sService)(nil).DeleteNode), clusterId, nodepoolId, nodeId, params)
}

// DeleteNodePool mocks base method.
func (m *MockK8sService) DeleteNodePool(clusterId, nodepoolId string, params resources.QueryParams) (*resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteNodePool", clusterId, nodepoolId, params)
	ret0, _ := ret[0].(*resources.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteNodePool indicates an expected call of DeleteNodePool.
func (mr *MockK8sServiceMockRecorder) DeleteNodePool(clusterId, nodepoolId, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteNodePool", reflect.TypeOf((*MockK8sService)(nil).DeleteNodePool), clusterId, nodepoolId, params)
}

// GetCluster mocks base method.
func (m *MockK8sService) GetCluster(clusterId string, params resources.QueryParams) (*resources.K8sCluster, *resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCluster", clusterId, params)
	ret0, _ := ret[0].(*resources.K8sCluster)
	ret1, _ := ret[1].(*resources.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetCluster indicates an expected call of GetCluster.
func (mr *MockK8sServiceMockRecorder) GetCluster(clusterId, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCluster", reflect.TypeOf((*MockK8sService)(nil).GetCluster), clusterId, params)
}

// GetNode mocks base method.
func (m *MockK8sService) GetNode(clusterId, nodepoolId, nodeId string, params resources.QueryParams) (*resources.K8sNode, *resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNode", clusterId, nodepoolId, nodeId, params)
	ret0, _ := ret[0].(*resources.K8sNode)
	ret1, _ := ret[1].(*resources.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetNode indicates an expected call of GetNode.
func (mr *MockK8sServiceMockRecorder) GetNode(clusterId, nodepoolId, nodeId, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNode", reflect.TypeOf((*MockK8sService)(nil).GetNode), clusterId, nodepoolId, nodeId, params)
}

// GetNodePool mocks base method.
func (m *MockK8sService) GetNodePool(clusterId, nodepoolId string, params resources.QueryParams) (*resources.K8sNodePool, *resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNodePool", clusterId, nodepoolId, params)
	ret0, _ := ret[0].(*resources.K8sNodePool)
	ret1, _ := ret[1].(*resources.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetNodePool indicates an expected call of GetNodePool.
func (mr *MockK8sServiceMockRecorder) GetNodePool(clusterId, nodepoolId, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNodePool", reflect.TypeOf((*MockK8sService)(nil).GetNodePool), clusterId, nodepoolId, params)
}

// GetVersion mocks base method.
func (m *MockK8sService) GetVersion() (string, *resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetVersion")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(*resources.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetVersion indicates an expected call of GetVersion.
func (mr *MockK8sServiceMockRecorder) GetVersion() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetVersion", reflect.TypeOf((*MockK8sService)(nil).GetVersion))
}

// ListClusters mocks base method.
func (m *MockK8sService) ListClusters(params resources.ListQueryParams) (resources.K8sClusters, *resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListClusters", params)
	ret0, _ := ret[0].(resources.K8sClusters)
	ret1, _ := ret[1].(*resources.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ListClusters indicates an expected call of ListClusters.
func (mr *MockK8sServiceMockRecorder) ListClusters(params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListClusters", reflect.TypeOf((*MockK8sService)(nil).ListClusters), params)
}

// ListNodePools mocks base method.
func (m *MockK8sService) ListNodePools(clusterId string, params resources.ListQueryParams) (resources.K8sNodePools, *resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListNodePools", clusterId, params)
	ret0, _ := ret[0].(resources.K8sNodePools)
	ret1, _ := ret[1].(*resources.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ListNodePools indicates an expected call of ListNodePools.
func (mr *MockK8sServiceMockRecorder) ListNodePools(clusterId, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListNodePools", reflect.TypeOf((*MockK8sService)(nil).ListNodePools), clusterId, params)
}

// ListNodes mocks base method.
func (m *MockK8sService) ListNodes(clusterId, nodepoolId string, params resources.ListQueryParams) (resources.K8sNodes, *resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListNodes", clusterId, nodepoolId, params)
	ret0, _ := ret[0].(resources.K8sNodes)
	ret1, _ := ret[1].(*resources.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ListNodes indicates an expected call of ListNodes.
func (mr *MockK8sServiceMockRecorder) ListNodes(clusterId, nodepoolId, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListNodes", reflect.TypeOf((*MockK8sService)(nil).ListNodes), clusterId, nodepoolId, params)
}

// ListVersions mocks base method.
func (m *MockK8sService) ListVersions() ([]string, *resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListVersions")
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(*resources.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ListVersions indicates an expected call of ListVersions.
func (mr *MockK8sServiceMockRecorder) ListVersions() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListVersions", reflect.TypeOf((*MockK8sService)(nil).ListVersions))
}

// ReadKubeConfig mocks base method.
func (m *MockK8sService) ReadKubeConfig(clusterId string) (string, *resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadKubeConfig", clusterId)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(*resources.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ReadKubeConfig indicates an expected call of ReadKubeConfig.
func (mr *MockK8sServiceMockRecorder) ReadKubeConfig(clusterId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadKubeConfig", reflect.TypeOf((*MockK8sService)(nil).ReadKubeConfig), clusterId)
}

// RecreateNode mocks base method.
func (m *MockK8sService) RecreateNode(clusterId, nodepoolId, nodeId string, params resources.QueryParams) (*resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RecreateNode", clusterId, nodepoolId, nodeId, params)
	ret0, _ := ret[0].(*resources.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RecreateNode indicates an expected call of RecreateNode.
func (mr *MockK8sServiceMockRecorder) RecreateNode(clusterId, nodepoolId, nodeId, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RecreateNode", reflect.TypeOf((*MockK8sService)(nil).RecreateNode), clusterId, nodepoolId, nodeId, params)
}

// UpdateCluster mocks base method.
func (m *MockK8sService) UpdateCluster(clusterId string, input resources.K8sClusterForPut, params resources.QueryParams) (*resources.K8sCluster, *resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateCluster", clusterId, input, params)
	ret0, _ := ret[0].(*resources.K8sCluster)
	ret1, _ := ret[1].(*resources.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// UpdateCluster indicates an expected call of UpdateCluster.
func (mr *MockK8sServiceMockRecorder) UpdateCluster(clusterId, input, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateCluster", reflect.TypeOf((*MockK8sService)(nil).UpdateCluster), clusterId, input, params)
}

// UpdateNodePool mocks base method.
func (m *MockK8sService) UpdateNodePool(clusterId, nodepoolId string, nodepool resources.K8sNodePoolForPut, params resources.QueryParams) (*resources.K8sNodePool, *resources.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateNodePool", clusterId, nodepoolId, nodepool, params)
	ret0, _ := ret[0].(*resources.K8sNodePool)
	ret1, _ := ret[1].(*resources.Response)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// UpdateNodePool indicates an expected call of UpdateNodePool.
func (mr *MockK8sServiceMockRecorder) UpdateNodePool(clusterId, nodepoolId, nodepool, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateNodePool", reflect.TypeOf((*MockK8sService)(nil).UpdateNodePool), clusterId, nodepoolId, nodepool, params)
}
