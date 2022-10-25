package resources

import (
	"context"
	sdkgo "github.com/ionos-cloud/sdk-go-dbaas-mongo"
)

type ClusterResponse struct {
	sdkgo.ClusterResponse
}

type ClusterList struct {
	sdkgo.ClusterList
}

type CreateClusterRequest struct {
	sdkgo.CreateClusterRequest
}

type CreateClusterProperties struct {
	sdkgo.CreateClusterProperties
}

//type PatchClusterRequest struct {
//	sdkgo.PatchClusterRequest
//}
//
//type PatchClusterProperties struct {
//	sdkgo.PatchClusterProperties
//}

type Response struct {
	sdkgo.APIResponse
}

// ClustersService is a wrapper around ionoscloud.Cluster
type ClustersService interface {
	List(filterName string) (sdkgo.ClusterList, *sdkgo.APIResponse, error)
	Get(clusterId string) (sdkgo.ClusterResponse, *sdkgo.APIResponse, error)
	Create(input sdkgo.CreateClusterRequest) (sdkgo.ClusterResponse, *sdkgo.APIResponse, error)
	Delete(clusterId string) (*sdkgo.APIResponse, error)
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

func (svc *clustersService) List(filterName string) (sdkgo.ClusterList, *sdkgo.APIResponse, error) {
	req := svc.client.ClustersApi.ClustersGet(svc.context)
	if filterName != "" {
		req = req.FilterName(filterName)
	}
	clusterList, res, err := svc.client.ClustersApi.ClustersGetExecute(req)
	return clusterList, res, err
}

func (svc *clustersService) Get(clusterId string) (sdkgo.ClusterResponse, *sdkgo.APIResponse, error) {
	req := svc.client.ClustersApi.ClustersFindById(svc.context, clusterId)
	cluster, res, err := svc.client.ClustersApi.ClustersFindByIdExecute(req)
	return cluster, res, err
}

func (svc *clustersService) Create(input sdkgo.CreateClusterRequest) (sdkgo.ClusterResponse, *sdkgo.APIResponse, error) {
	req := svc.client.ClustersApi.ClustersPost(svc.context).CreateClusterRequest(input)
	cluster, res, err := svc.client.ClustersApi.ClustersPostExecute(req)
	return cluster, res, err
}

func (svc *clustersService) Delete(clusterId string) (*sdkgo.APIResponse, error) {
	req := svc.client.ClustersApi.ClustersDelete(svc.context, clusterId)
	_, res, err := svc.client.ClustersApi.ClustersDeleteExecute(req)
	return res, err
}
