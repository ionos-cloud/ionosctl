package resources

import (
	"context"
	"time"

	sdkgo "github.com/ionos-cloud/sdk-go-dbaas-postgres"
)

type Cluster struct {
	sdkgo.ClusterList
}

type ClusterList struct {
	sdkgo.ClusterListAllOf
}

type CreateClusterRequest struct {
	sdkgo.CreateClusterRequest
}

type PatchClusterRequest struct {
	sdkgo.PatchClusterRequest
}

type Response struct {
	sdkgo.APIResponse
}

// ClustersService is a wrapper around ionoscloud.Cluster
type ClustersService interface {
	List(filterName string) (ClusterList, *Response, error)
	Get(clusterId string) (*Cluster, *Response, error)
	Create(input CreateClusterRequest, backupId string, recoveryTargetTime time.Time) (*Cluster, *Response, error)
	Update(clusterId string, input PatchClusterRequest) (*Cluster, *Response, error)
	Delete(clusterId string) (*Response, error)
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

func (svc *clustersService) List(filterName string) (ClusterList, *Response, error) {
	req := svc.client.ClustersApi.ClustersGet(svc.context)
	if filterName != "" {
		req = req.FilterName(filterName)
	}
	clusterList, res, err := svc.client.ClustersApi.ClustersGetExecute(req)
	return ClusterList{clusterList}, &Response{*res}, err
}

func (svc *clustersService) Get(clusterId string) (*Cluster, *Response, error) {
	req := svc.client.ClustersApi.ClustersFindById(svc.context, clusterId)
	cluster, res, err := svc.client.ClustersApi.ClustersFindByIdExecute(req)
	return &Cluster{cluster}, &Response{*res}, err
}

func (svc *clustersService) Create(input CreateClusterRequest, backupId string, recoveryTargetTime time.Time) (*Cluster, *Response, error) {
	req := svc.client.ClustersApi.ClustersPost(svc.context).CreateClusterRequest(input.CreateClusterRequest)
	if backupId != "" {
		// Create Cluster from a specified Backup
		req = req.FromBackup(backupId)
	}
	if !recoveryTargetTime.IsZero() {
		// Create Cluster from a specified Backup from a specific timestamp
		req = req.FromRecoveryTargetTime(recoveryTargetTime)
	}
	cluster, res, err := svc.client.ClustersApi.ClustersPostExecute(req)
	return &Cluster{cluster}, &Response{*res}, err
}

func (svc *clustersService) Update(clusterId string, input PatchClusterRequest) (*Cluster, *Response, error) {
	req := svc.client.ClustersApi.ClustersPatch(svc.context, clusterId).Cluster(input.PatchClusterRequest)
	cluster, res, err := svc.client.ClustersApi.ClustersPatchExecute(req)
	return &Cluster{cluster}, &Response{*res}, err
}

func (svc *clustersService) Delete(clusterId string) (*Response, error) {
	req := svc.client.ClustersApi.ClustersDelete(svc.context, clusterId)
	_, res, err := svc.client.ClustersApi.ClustersDeleteExecute(req)
	return &Response{*res}, err
}
