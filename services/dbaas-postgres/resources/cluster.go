package resources

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"

	"github.com/ionos-cloud/sdk-go-bundle/products/dbaas/psql/v2"
)

type ClusterResponse struct {
	psql.ClusterResponse
}

type ClusterList struct {
	psql.ClusterList
}

type CreateClusterRequest struct {
	psql.CreateClusterRequest
}

type CreateClusterProperties struct {
	psql.CreateClusterProperties
}

type PatchClusterRequest struct {
	psql.PatchClusterRequest
}

type PatchClusterProperties struct {
	psql.PatchClusterProperties
}

type Response struct {
	psql.APIResponse
}

// ClustersService is a wrapper around ionoscloud.Cluster
type ClustersService interface {
	List(filterName string) (ClusterList, *Response, error)
	Get(clusterId string) (*ClusterResponse, *Response, error)
	Create(input CreateClusterRequest) (*ClusterResponse, *Response, error)
	Update(clusterId string, input PatchClusterRequest) (*ClusterResponse, *Response, error)
	Delete(clusterId string) (*Response, error)
}

type clustersService struct {
	client  *psql.APIClient
	context context.Context
}

var _ ClustersService = &clustersService{}

func NewClustersService(client *client.Client, ctx context.Context) ClustersService {
	return &clustersService{
		client:  client.PostgresClient,
		context: ctx,
	}
}

func (svc *clustersService) List(filterName string) (ClusterList, *Response, error) {
	req := svc.client.ClustersApi.ClustersGet(svc.context)
	if filterName != "" {
		req = req.FilterName(filterName)
	}
	clusterList, res, err := svc.client.ClustersApi.ClustersGetExecute(req)
	return ClusterList{clusterList}, &Response{*res}, err
}

func (svc *clustersService) Get(clusterId string) (*ClusterResponse, *Response, error) {
	req := svc.client.ClustersApi.ClustersFindById(svc.context, clusterId)
	cluster, res, err := svc.client.ClustersApi.ClustersFindByIdExecute(req)
	return &ClusterResponse{cluster}, &Response{*res}, err
}

func (svc *clustersService) Create(input CreateClusterRequest) (*ClusterResponse, *Response, error) {
	req := svc.client.ClustersApi.ClustersPost(svc.context).CreateClusterRequest(input.CreateClusterRequest)
	cluster, res, err := svc.client.ClustersApi.ClustersPostExecute(req)
	return &ClusterResponse{cluster}, &Response{*res}, err
}

func (svc *clustersService) Update(clusterId string, input PatchClusterRequest) (*ClusterResponse, *Response, error) {
	req := svc.client.ClustersApi.ClustersPatch(svc.context, clusterId).PatchClusterRequest(input.PatchClusterRequest)
	cluster, res, err := svc.client.ClustersApi.ClustersPatchExecute(req)
	return &ClusterResponse{cluster}, &Response{*res}, err
}

func (svc *clustersService) Delete(clusterId string) (*Response, error) {
	req := svc.client.ClustersApi.ClustersDelete(svc.context, clusterId)
	_, res, err := svc.client.ClustersApi.ClustersDeleteExecute(req)
	return &Response{*res}, err
}
