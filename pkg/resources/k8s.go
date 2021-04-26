package resources

import (
	"context"

	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
)

type K8sCluster struct {
	ionoscloud.KubernetesCluster
}

type K8sKubeconfig struct {
	ionoscloud.KubernetesConfig
}

type K8sClusterProperties struct {
	ionoscloud.KubernetesClusterProperties
}

type K8sClusters struct {
	ionoscloud.KubernetesClusters
}

type K8sNodePool struct {
	ionoscloud.KubernetesNodePool
}

type K8sNodePoolProperties struct {
	ionoscloud.KubernetesNodePoolProperties
}

type K8sNodePoolUpdated struct {
	ionoscloud.KubernetesNodePoolForPut
}

type K8sNodePools struct {
	ionoscloud.KubernetesNodePools
}

type K8sNode struct {
	ionoscloud.KubernetesNode
}

type K8sNodeProperties struct {
	ionoscloud.KubernetesNodeProperties
}

type K8sNodes struct {
	ionoscloud.KubernetesNodes
}

// K8sService is a wrapper around ionoscloud.K8s
type K8sService interface {
	ListClusters() (K8sClusters, *Response, error)
	GetCluster(clusterId string) (*K8sCluster, *Response, error)
	CreateCluster(u K8sCluster) (*K8sCluster, *Response, error)
	UpdateCluster(clusterId string, input K8sCluster) (*K8sCluster, *Response, error)
	DeleteCluster(clusterId string) (*Response, error)
	ReadKubeConfig(clusterId string) (K8sKubeconfig, *Response, error)
	// Node Pools
	ListNodePools(clusterId string) (K8sNodePools, *Response, error)
	GetNodePool(clusterId, nodepoolId string) (*K8sNodePool, *Response, error)
	CreateNodePool(clusterId string, nodepool K8sNodePool) (*K8sNodePool, *Response, error)
	UpdateNodePool(clusterId, nodepoolId string, nodepool K8sNodePool) (*K8sNodePoolUpdated, *Response, error)
	DeleteNodePool(clusterId, nodepoolId string) (*Response, error)
	// Node
	DeleteNode(clusterId, nodepoolId, nodeId string) (*Response, error)
	RecreateNode(clusterId, nodepoolId, nodeId string) (*Response, error)
	GetNode(clusterId, nodepoolId, nodeId string) (*K8sNode, *Response, error)
	ListNodes(clusterId, nodepoolId string) (K8sNodes, *Response, error)
	// Versions
	ListVersions() ([]string, *Response, error)
	GetVersion() (string, *Response, error)
	ListCompatibleVersions(clusterVersion string) ([]string, *Response, error)
}

type k8sService struct {
	client  *Client
	context context.Context
}

var _ K8sService = &k8sService{}

func NewK8sService(client *Client, ctx context.Context) K8sService {
	return &k8sService{
		client:  client,
		context: ctx,
	}
}

func (s *k8sService) ListClusters() (K8sClusters, *Response, error) {
	req := s.client.KubernetesApi.K8sGet(s.context)
	dcs, res, err := s.client.KubernetesApi.K8sGetExecute(req)
	return K8sClusters{dcs}, &Response{*res}, err
}

func (s *k8sService) GetCluster(clusterId string) (*K8sCluster, *Response, error) {
	req := s.client.KubernetesApi.K8sFindBySClusterId(s.context, clusterId)
	user, res, err := s.client.KubernetesApi.K8sFindBySClusterIdExecute(req)
	return &K8sCluster{user}, &Response{*res}, err
}

func (s *k8sService) CreateCluster(u K8sCluster) (*K8sCluster, *Response, error) {
	req := s.client.KubernetesApi.K8sPost(s.context).KubernetesCluster(u.KubernetesCluster)
	user, res, err := s.client.KubernetesApi.K8sPostExecute(req)
	return &K8sCluster{user}, &Response{*res}, err
}

func (s *k8sService) UpdateCluster(clusterId string, input K8sCluster) (*K8sCluster, *Response, error) {
	req := s.client.KubernetesApi.K8sPut(s.context, clusterId).KubernetesCluster(input.KubernetesCluster)
	user, res, err := s.client.KubernetesApi.K8sPutExecute(req)
	return &K8sCluster{user}, &Response{*res}, err
}

func (s *k8sService) DeleteCluster(clusterId string) (*Response, error) {
	req := s.client.KubernetesApi.K8sDelete(s.context, clusterId)
	_, res, err := s.client.KubernetesApi.K8sDeleteExecute(req)
	return &Response{*res}, err
}

func (s *k8sService) ReadKubeConfig(clusterId string) (K8sKubeconfig, *Response, error) {
	req := s.client.KubernetesApi.K8sKubeconfigGet(s.context, clusterId)
	file, res, err := s.client.KubernetesApi.K8sKubeconfigGetExecute(req)
	return K8sKubeconfig{file}, &Response{*res}, err
}

func (s *k8sService) ListNodePools(clusterId string) (K8sNodePools, *Response, error) {
	req := s.client.KubernetesApi.K8sNodepoolsGet(s.context, clusterId)
	ns, res, err := s.client.KubernetesApi.K8sNodepoolsGetExecute(req)
	return K8sNodePools{ns}, &Response{*res}, err
}

func (s *k8sService) GetNodePool(clusterId, nodepoolId string) (*K8sNodePool, *Response, error) {
	req := s.client.KubernetesApi.K8sNodepoolsFindById(s.context, clusterId, nodepoolId)
	ns, res, err := s.client.KubernetesApi.K8sNodepoolsFindByIdExecute(req)
	return &K8sNodePool{ns}, &Response{*res}, err
}

func (s *k8sService) CreateNodePool(clusterId string, nodepool K8sNodePool) (*K8sNodePool, *Response, error) {
	req := s.client.KubernetesApi.K8sNodepoolsPost(s.context, clusterId).KubernetesNodePool(nodepool.KubernetesNodePool)
	ns, res, err := s.client.KubernetesApi.K8sNodepoolsPostExecute(req)
	return &K8sNodePool{ns}, &Response{*res}, err
}

func (s *k8sService) UpdateNodePool(clusterId, nodepoolId string, nodepool K8sNodePool) (*K8sNodePoolUpdated, *Response, error) {
	req := s.client.KubernetesApi.K8sNodepoolsPut(s.context, clusterId, nodepoolId).KubernetesNodePool(nodepool.KubernetesNodePool)
	ns, res, err := s.client.KubernetesApi.K8sNodepoolsPutExecute(req)
	return &K8sNodePoolUpdated{ns}, &Response{*res}, err
}

func (s *k8sService) DeleteNodePool(clusterId, nodepoolId string) (*Response, error) {
	req := s.client.KubernetesApi.K8sNodepoolsDelete(s.context, clusterId, nodepoolId)
	_, res, err := s.client.KubernetesApi.K8sNodepoolsDeleteExecute(req)
	return &Response{*res}, err
}

func (s *k8sService) DeleteNode(clusterId, nodepoolId, nodeId string) (*Response, error) {
	req := s.client.KubernetesApi.K8sNodepoolsNodesDelete(s.context, clusterId, nodepoolId, nodeId)
	_, res, err := s.client.KubernetesApi.K8sNodepoolsNodesDeleteExecute(req)
	return &Response{*res}, err
}

func (s *k8sService) RecreateNode(clusterId, nodepoolId, nodeId string) (*Response, error) {
	req := s.client.KubernetesApi.K8sNodepoolsNodesReplacePost(s.context, clusterId, nodepoolId, nodeId)
	_, res, err := s.client.KubernetesApi.K8sNodepoolsNodesReplacePostExecute(req)
	return &Response{*res}, err
}

func (s *k8sService) GetNode(clusterId, nodepoolId, nodeId string) (*K8sNode, *Response, error) {
	req := s.client.KubernetesApi.K8sNodepoolsNodesFindById(s.context, clusterId, nodepoolId, nodeId)
	n, res, err := s.client.KubernetesApi.K8sNodepoolsNodesFindByIdExecute(req)
	return &K8sNode{n}, &Response{*res}, err
}

func (s *k8sService) ListNodes(clusterId, nodepoolId string) (K8sNodes, *Response, error) {
	req := s.client.KubernetesApi.K8sNodepoolsNodesGet(s.context, clusterId, nodepoolId)
	ns, res, err := s.client.KubernetesApi.K8sNodepoolsNodesGetExecute(req)
	return K8sNodes{ns}, &Response{*res}, err
}

func (s *k8sService) ListVersions() ([]string, *Response, error) {
	req := s.client.KubernetesApi.K8sVersionsGet(s.context)
	vs, res, err := s.client.KubernetesApi.K8sVersionsGetExecute(req)
	return vs, &Response{*res}, err
}

func (s *k8sService) GetVersion() (string, *Response, error) {
	req := s.client.KubernetesApi.K8sVersionsDefaultGet(s.context)
	v, res, err := s.client.KubernetesApi.K8sVersionsDefaultGetExecute(req)
	return v, &Response{*res}, err
}

func (s *k8sService) ListCompatibleVersions(clusterVersion string) ([]string, *Response, error) {
	req := s.client.KubernetesApi.K8sVersionsCompatibilitiesGet(s.context, clusterVersion)
	vs, res, err := s.client.KubernetesApi.K8sVersionsCompatibilitiesGetExecute(req)
	return vs, &Response{*res}, err
}
