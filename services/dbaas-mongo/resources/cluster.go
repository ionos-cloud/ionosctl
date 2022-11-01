package resources

import (
	"context"

	sdkgo "github.com/ionos-cloud/sdk-go-dbaas-mongo"
)

type ClustersService interface {
	List(filterName string) (sdkgo.ClusterList, *sdkgo.APIResponse, error)
	Get(clusterId string) (sdkgo.ClusterResponse, *sdkgo.APIResponse, error)
	Create(input sdkgo.CreateClusterRequest) (sdkgo.ClusterResponse, *sdkgo.APIResponse, error)
	Delete(clusterId string) (*sdkgo.APIResponse, error)
	Restore(clusterId, snapshotId string) (*sdkgo.APIResponse, error)
	SnapshotsList(clusterId string) (sdkgo.SnapshotList, *sdkgo.APIResponse, error)
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

func (svc *clustersService) Restore(clusterId, snapshotId string) (*sdkgo.APIResponse, error) {
	req := svc.client.RestoresApi.ClustersRestorePost(svc.context, clusterId)
	req.CreateRestoreRequest(sdkgo.CreateRestoreRequest{SnapshotId: &snapshotId})
	res, err := svc.client.RestoresApi.ClustersRestorePostExecute(req)
	return res, err
}

func (svc *clustersService) SnapshotsList(clusterId string) (sdkgo.SnapshotList, *sdkgo.APIResponse, error) {
	req := svc.client.SnapshotsApi.ClustersSnapshotsGet(svc.context, clusterId)
	snapshots, res, err := svc.client.SnapshotsApi.ClustersSnapshotsGetExecute(req)
	return snapshots, res, err
}
