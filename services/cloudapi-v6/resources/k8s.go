package resources

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	compute "github.com/ionos-cloud/sdk-go/v6"

	"github.com/fatih/structs"
)

type K8sCluster struct {
	compute.KubernetesCluster
}

type K8sClusterForPost struct {
	compute.KubernetesClusterForPost
}

type K8sClusterForPut struct {
	compute.KubernetesClusterForPut
}

type K8sClusterProperties struct {
	compute.KubernetesClusterProperties
}

type K8sClusterPropertiesForPut struct {
	compute.KubernetesClusterPropertiesForPut
}

type K8sClusterPropertiesForPost struct {
	compute.KubernetesClusterPropertiesForPost
}

type K8sClusters struct {
	compute.KubernetesClusters
}

type K8sNodePool struct {
	compute.KubernetesNodePool
}

type K8sNodePoolLan struct {
	compute.KubernetesNodePoolLan
}

type K8sNodePoolProperties struct {
	compute.KubernetesNodePoolProperties
}

type K8sNodePoolPropertiesForPut struct {
	compute.KubernetesNodePoolPropertiesForPut
}

type K8sNodePoolPropertiesForPost struct {
	compute.KubernetesNodePoolPropertiesForPost
}

type K8sNodePoolForPost struct {
	compute.KubernetesNodePoolForPost
}

type K8sNodePoolForPut struct {
	compute.KubernetesNodePoolForPut
}

type K8sNodePools struct {
	compute.KubernetesNodePools
}

type K8sNode struct {
	compute.KubernetesNode
}

type K8sNodeProperties struct {
	compute.KubernetesNodeProperties
}

type K8sNodes struct {
	compute.KubernetesNodes
}

type K8sMaintenanceWindow struct {
	compute.KubernetesMaintenanceWindow
}

// K8sService is a wrapper around compute.K8s
type K8sService interface {
	ListClusters(params ListQueryParams) (K8sClusters, *Response, error)
	GetCluster(clusterId string, params QueryParams) (*K8sCluster, *Response, error)
	// IsPublicCluster(clusterId string) (bool, error)
	CreateCluster(u K8sClusterForPost, params QueryParams) (*K8sCluster, *Response, error)
	UpdateCluster(clusterId string, input K8sClusterForPut, params QueryParams) (*K8sCluster, *Response, error)
	DeleteCluster(clusterId string, params QueryParams) (*Response, error)
	ReadKubeConfig(clusterId string) (string, *Response, error)
	ListNodePools(clusterId string, params ListQueryParams) (K8sNodePools, *Response, error)
	GetNodePool(clusterId, nodepoolId string, params QueryParams) (*K8sNodePool, *Response, error)
	CreateNodePool(clusterId string, nodepool K8sNodePoolForPost, params QueryParams) (*K8sNodePool, *Response, error)
	UpdateNodePool(clusterId, nodepoolId string, nodepool K8sNodePoolForPut, params QueryParams) (*K8sNodePool, *Response, error)
	DeleteNodePool(clusterId, nodepoolId string, params QueryParams) (*Response, error)
	DeleteNode(clusterId, nodepoolId, nodeId string, params QueryParams) (*Response, error)
	RecreateNode(clusterId, nodepoolId, nodeId string, params QueryParams) (*Response, error)
	GetNode(clusterId, nodepoolId, nodeId string, params QueryParams) (*K8sNode, *Response, error)
	ListNodes(clusterId, nodepoolId string, params ListQueryParams) (K8sNodes, *Response, error)
	ListVersions() ([]string, *Response, error)
	GetVersion() (string, *Response, error)
}

type k8sService struct {
	client  *compute.APIClient
	context context.Context
}

var _ K8sService = &k8sService{}

func NewK8sService(client *client.Client, ctx context.Context) K8sService {
	return &k8sService{
		client:  client.CloudClient,
		context: ctx,
	}
}

func (s *k8sService) ListClusters(params ListQueryParams) (K8sClusters, *Response, error) {
	req := s.client.KubernetesApi.K8sGet(s.context)
	if !structs.IsZero(params) {
		if params.Filters != nil {
			for k, v := range *params.Filters {
				for _, val := range v {
					req = req.Filter(k, val)
				}
			}
		}
		if params.OrderBy != nil {
			req = req.OrderBy(*params.OrderBy)
		}
		if params.MaxResults != nil {
			req = req.MaxResults(*params.MaxResults)
		}
		if !structs.IsZero(params.QueryParams) {
			if params.QueryParams.Depth != nil {
				req = req.Depth(*params.QueryParams.Depth)
			}
			if params.QueryParams.Pretty != nil {
				// Currently not implemented
				req = req.Pretty(*params.QueryParams.Pretty)
			}
		}
	}
	dcs, res, err := s.client.KubernetesApi.K8sGetExecute(req)
	return K8sClusters{dcs}, &Response{*res}, err
}

func (s *k8sService) GetCluster(clusterId string, params QueryParams) (*K8sCluster, *Response, error) {
	req := s.client.KubernetesApi.K8sFindByClusterId(s.context, clusterId)
	if !structs.IsZero(params) {
		if params.Depth != nil {
			req = req.Depth(*params.Depth)
		}
		if params.Pretty != nil {
			// Currently not implemented
			req = req.Pretty(*params.Pretty)
		}
	}
	k8sCluster, res, err := s.client.KubernetesApi.K8sFindByClusterIdExecute(req)
	return &K8sCluster{k8sCluster}, &Response{*res}, err
}

func (s *k8sService) CreateCluster(u K8sClusterForPost, params QueryParams) (*K8sCluster, *Response, error) {
	req := s.client.KubernetesApi.K8sPost(s.context).KubernetesCluster(u.KubernetesClusterForPost)
	if !structs.IsZero(params) {
		if params.Depth != nil {
			req = req.Depth(*params.Depth)
		}
		if params.Pretty != nil {
			// Currently not implemented
			req = req.Pretty(*params.Pretty)
		}
	}
	k8sCluster, res, err := s.client.KubernetesApi.K8sPostExecute(req)
	return &K8sCluster{k8sCluster}, &Response{*res}, err
}

func (s *k8sService) UpdateCluster(clusterId string, input K8sClusterForPut, params QueryParams) (*K8sCluster, *Response, error) {
	req := s.client.KubernetesApi.K8sPut(s.context, clusterId).KubernetesCluster(input.KubernetesClusterForPut)
	if !structs.IsZero(params) {
		if params.Depth != nil {
			req = req.Depth(*params.Depth)
		}
		if params.Pretty != nil {
			// Currently not implemented
			req = req.Pretty(*params.Pretty)
		}
	}
	k8sCluster, res, err := s.client.KubernetesApi.K8sPutExecute(req)
	return &K8sCluster{k8sCluster}, &Response{*res}, err
}

func (s *k8sService) DeleteCluster(clusterId string, params QueryParams) (*Response, error) {
	req := s.client.KubernetesApi.K8sDelete(s.context, clusterId)
	if !structs.IsZero(params) {
		if params.Depth != nil {
			req = req.Depth(*params.Depth)
		}
		if params.Pretty != nil {
			// Currently not implemented
			req = req.Pretty(*params.Pretty)
		}
	}
	res, err := s.client.KubernetesApi.K8sDeleteExecute(req)
	return &Response{*res}, err
}

func (s *k8sService) ReadKubeConfig(clusterId string) (string, *Response, error) {
	req := s.client.KubernetesApi.K8sKubeconfigGet(s.context, clusterId)
	file, res, err := s.client.KubernetesApi.K8sKubeconfigGetExecute(req)
	return file, &Response{*res}, err
}

func (s *k8sService) ListNodePools(clusterId string, params ListQueryParams) (K8sNodePools, *Response, error) {
	req := s.client.KubernetesApi.K8sNodepoolsGet(s.context, clusterId)
	if !structs.IsZero(params) {
		if params.Filters != nil {
			for k, v := range *params.Filters {
				for _, val := range v {
					req = req.Filter(k, val)
				}
			}
		}
		if params.OrderBy != nil {
			req = req.OrderBy(*params.OrderBy)
		}
		if params.MaxResults != nil {
			req = req.MaxResults(*params.MaxResults)
		}
		if !structs.IsZero(params.QueryParams) {
			if params.QueryParams.Depth != nil {
				req = req.Depth(*params.QueryParams.Depth)
			}
			if params.QueryParams.Pretty != nil {
				// Currently not implemented
				req = req.Pretty(*params.QueryParams.Pretty)
			}
		}
	}
	ns, res, err := s.client.KubernetesApi.K8sNodepoolsGetExecute(req)
	return K8sNodePools{ns}, &Response{*res}, err
}

func (s *k8sService) GetNodePool(clusterId, nodepoolId string, params QueryParams) (*K8sNodePool, *Response, error) {
	req := s.client.KubernetesApi.K8sNodepoolsFindById(s.context, clusterId, nodepoolId)
	if !structs.IsZero(params) {
		if params.Depth != nil {
			req = req.Depth(*params.Depth)
		}
		if params.Pretty != nil {
			// Currently not implemented
			req = req.Pretty(*params.Pretty)
		}
	}
	k8sNodePool, res, err := s.client.KubernetesApi.K8sNodepoolsFindByIdExecute(req)
	return &K8sNodePool{k8sNodePool}, &Response{*res}, err
}

func (s *k8sService) CreateNodePool(clusterId string, nodepool K8sNodePoolForPost, params QueryParams) (*K8sNodePool, *Response, error) {
	req := s.client.KubernetesApi.K8sNodepoolsPost(s.context, clusterId).KubernetesNodePool(nodepool.KubernetesNodePoolForPost)
	if !structs.IsZero(params) {
		if params.Depth != nil {
			req = req.Depth(*params.Depth)
		}
		if params.Pretty != nil {
			// Currently not implemented
			req = req.Pretty(*params.Pretty)
		}
	}
	k8sNodePool, res, err := s.client.KubernetesApi.K8sNodepoolsPostExecute(req)
	return &K8sNodePool{k8sNodePool}, &Response{*res}, err
}

func (s *k8sService) UpdateNodePool(clusterId, nodepoolId string, nodepool K8sNodePoolForPut, params QueryParams) (*K8sNodePool, *Response, error) {
	req := s.client.KubernetesApi.K8sNodepoolsPut(s.context, clusterId, nodepoolId).KubernetesNodePool(nodepool.KubernetesNodePoolForPut)
	if !structs.IsZero(params) {
		if params.Depth != nil {
			req = req.Depth(*params.Depth)
		}
		if params.Pretty != nil {
			// Currently not implemented
			req = req.Pretty(*params.Pretty)
		}
	}
	k8sNodePool, res, err := s.client.KubernetesApi.K8sNodepoolsPutExecute(req)
	return &K8sNodePool{k8sNodePool}, &Response{*res}, err
}

func (s *k8sService) DeleteNodePool(clusterId, nodepoolId string, params QueryParams) (*Response, error) {
	req := s.client.KubernetesApi.K8sNodepoolsDelete(s.context, clusterId, nodepoolId)
	if !structs.IsZero(params) {
		if params.Depth != nil {
			req = req.Depth(*params.Depth)
		}
		if params.Pretty != nil {
			// Currently not implemented
			req = req.Pretty(*params.Pretty)
		}
	}
	res, err := s.client.KubernetesApi.K8sNodepoolsDeleteExecute(req)
	return &Response{*res}, err
}

func (s *k8sService) DeleteNode(clusterId, nodepoolId, nodeId string, params QueryParams) (*Response, error) {
	req := s.client.KubernetesApi.K8sNodepoolsNodesDelete(s.context, clusterId, nodepoolId, nodeId)
	if !structs.IsZero(params) {
		if params.Depth != nil {
			req = req.Depth(*params.Depth)
		}
		if params.Pretty != nil {
			// Currently not implemented
			req = req.Pretty(*params.Pretty)
		}
	}
	res, err := s.client.KubernetesApi.K8sNodepoolsNodesDeleteExecute(req)
	return &Response{*res}, err
}

func (s *k8sService) RecreateNode(clusterId, nodepoolId, nodeId string, params QueryParams) (*Response, error) {
	req := s.client.KubernetesApi.K8sNodepoolsNodesReplacePost(s.context, clusterId, nodepoolId, nodeId)
	if !structs.IsZero(params) {
		if params.Depth != nil {
			req = req.Depth(*params.Depth)
		}
		if params.Pretty != nil {
			// Currently not implemented
			req = req.Pretty(*params.Pretty)
		}
	}
	res, err := s.client.KubernetesApi.K8sNodepoolsNodesReplacePostExecute(req)
	return &Response{*res}, err
}

func (s *k8sService) GetNode(clusterId, nodepoolId, nodeId string, params QueryParams) (*K8sNode, *Response, error) {
	req := s.client.KubernetesApi.K8sNodepoolsNodesFindById(s.context, clusterId, nodepoolId, nodeId)
	if !structs.IsZero(params) {
		if params.Depth != nil {
			req = req.Depth(*params.Depth)
		}
		if params.Pretty != nil {
			// Currently not implemented
			req = req.Pretty(*params.Pretty)
		}
	}
	k8sNode, res, err := s.client.KubernetesApi.K8sNodepoolsNodesFindByIdExecute(req)
	return &K8sNode{k8sNode}, &Response{*res}, err
}

func (s *k8sService) ListNodes(clusterId, nodepoolId string, params ListQueryParams) (K8sNodes, *Response, error) {
	req := s.client.KubernetesApi.K8sNodepoolsNodesGet(s.context, clusterId, nodepoolId)
	if !structs.IsZero(params) {
		if params.Filters != nil {
			for k, v := range *params.Filters {
				for _, val := range v {
					req = req.Filter(k, val)
				}
			}
		}
		if params.OrderBy != nil {
			req = req.OrderBy(*params.OrderBy)
		}
		if params.MaxResults != nil {
			req = req.MaxResults(*params.MaxResults)
		}
		if !structs.IsZero(params.QueryParams) {
			if params.QueryParams.Depth != nil {
				req = req.Depth(*params.QueryParams.Depth)
			}
			if params.QueryParams.Pretty != nil {
				// Currently not implemented
				req = req.Pretty(*params.QueryParams.Pretty)
			}
		}
	}
	k8sNodes, res, err := s.client.KubernetesApi.K8sNodepoolsNodesGetExecute(req)
	return K8sNodes{k8sNodes}, &Response{*res}, err
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
