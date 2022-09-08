package resources

import (
	"context"

	dp "github.com/ionos-cloud/sdk-go-dataplatform"
)

type ClusterResponseData struct {
	dp.ClusterResponseData
}

type ClusterListResponseData struct {
	dp.ClusterListResponseData
}

type CreateClusterRequest struct {
	dp.CreateClusterRequest
}

type CreateClusterProperties struct {
	dp.CreateClusterProperties
}

type PatchClusterRequest struct {
	dp.PatchClusterRequest
}

type PatchClusterProperties struct {
	dp.PatchClusterProperties
}

type Response struct {
	dp.APIResponse
}

// ClustersService is a wrapper around dp.Cluster
type ClustersService interface {
	Get(clusterId string) (ClusterResponseData, *Response, error)
	GetKubeConfig(clusterId string) (string, *Response, error)
	List(filterName string) (ClusterListResponseData, *Response, error)
	Create(cluster CreateClusterRequest) (ClusterResponseData, *Response, error)
	Update(clusterId string, cluster PatchClusterRequest) (ClusterResponseData, *Response, error)
	Delete(clusterId string) (ClusterResponseData, *Response, error)
}

type clustersService struct {
	client  *Client
	context context.Context
}

var _ ClustersService = &clustersService{}

func NewClustersService(client *Client, ctx context.Context) ClustersService {
	return &clustersService{
		client:  client,
		context: ctx,
	}
}

func (svc *clustersService) Get(clusterId string) (ClusterResponseData, *Response, error) {
	req := svc.client.DataPlatformClusterApi.GetCluster(svc.context, clusterId)
	clusterResponse, res, err := svc.client.DataPlatformClusterApi.GetClusterExecute(req)
	return ClusterResponseData{clusterResponse}, &Response{*res}, err
}

func (svc *clustersService) GetKubeConfig(clusterId string) (string, *Response, error) {
	req := svc.client.DataPlatformClusterApi.GetClusterKubeconfig(svc.context, clusterId)
	kubeConfigResponse, res, err := svc.client.DataPlatformClusterApi.GetClusterKubeconfigExecute(req)
	return kubeConfigResponse, &Response{*res}, err
}

func (svc *clustersService) List(filterName string) (ClusterListResponseData, *Response, error) {
	request := svc.client.DataPlatformClusterApi.GetClusters(svc.context)
	if filterName != "" {
		request = request.Name(filterName)
	}
	clusterListResponse, res, err := svc.client.DataPlatformClusterApi.GetClustersExecute(request)

	return ClusterListResponseData{clusterListResponse}, &Response{*res}, err
}

func (svc *clustersService) Create(cluster CreateClusterRequest) (ClusterResponseData, *Response, error) {
	req := svc.client.DataPlatformClusterApi.CreateCluster(svc.context).CreateClusterRequest(cluster.CreateClusterRequest)
	clusterResponse, res, err := svc.client.DataPlatformClusterApi.CreateClusterExecute(req)
	return ClusterResponseData{clusterResponse}, &Response{*res}, err
}

func (svc *clustersService) Update(clusterId string, cluster PatchClusterRequest) (ClusterResponseData, *Response, error) {
	req := svc.client.DataPlatformClusterApi.PatchCluster(svc.context, clusterId).PatchClusterRequest(cluster.PatchClusterRequest)
	clusterResponse, res, err := svc.client.DataPlatformClusterApi.PatchClusterExecute(req)
	return ClusterResponseData{clusterResponse}, &Response{*res}, err
}

func (svc *clustersService) Delete(clusterId string) (ClusterResponseData, *Response, error) {
	req := svc.client.DataPlatformClusterApi.DeleteCluster(svc.context, clusterId)
	clusterResponse, res, err := svc.client.DataPlatformClusterApi.DeleteClusterExecute(req)
	return ClusterResponseData{clusterResponse}, &Response{*res}, err
}
