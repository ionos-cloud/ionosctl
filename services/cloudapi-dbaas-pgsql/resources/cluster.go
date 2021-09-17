package resources

import (
	"context"

	sdkgo "github.com/ionos-cloud/sdk-go-autoscaling"
)

type Cluster struct {
	sdkgo.Cluster
}

type ClusterList struct {
	sdkgo.ClusterList
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
	List() (ClusterList, *Response, error)
	Get(clusterId string) (*Cluster, *Response, error)
	Create(input CreateClusterRequest, backupId string, recoveryTargetTime string) (*Cluster, *Response, error)
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

func (cs *clustersService) List() (ClusterList, *Response, error) {
	req := cs.client.ClustersApi.ClustersGet(cs.context)
	clusterList, res, err := cs.client.ClustersApi.ClustersGetExecute(req)
	return ClusterList{clusterList}, &Response{*res}, err
}

func (cs *clustersService) Get(clusterId string) (*Cluster, *Response, error) {
	req := cs.client.ClustersApi.ClustersFindById(cs.context, clusterId)
	cluster, res, err := cs.client.ClustersApi.ClustersFindByIdExecute(req)
	return &Cluster{cluster}, &Response{*res}, err
}

func (cs *clustersService) Create(input CreateClusterRequest, backupId string, recoveryTargetTime string) (*Cluster, *Response, error) {
	req := cs.client.ClustersApi.ClustersPost(cs.context).Cluster(input.CreateClusterRequest)
	if backupId != "" {
		req = req.FromBackup(backupId)
	}
	if recoveryTargetTime != "" {
		req = req.FromRecoveryTargetTime(recoveryTargetTime)
	}
	cluster, res, err := cs.client.ClustersApi.ClustersPostExecute(req)
	return &Cluster{cluster}, &Response{*res}, err
}

func (cs *clustersService) Update(clusterId string, input PatchClusterRequest) (*Cluster, *Response, error) {
	req := cs.client.ClustersApi.ClustersPatch(cs.context, clusterId).Cluster(input.PatchClusterRequest)
	cluster, res, err := cs.client.ClustersApi.ClustersPatchExecute(req)
	return &Cluster{cluster}, &Response{*res}, err
}

func (cs *clustersService) Delete(clusterId string) (*Response, error) {
	req := cs.client.ClustersApi.ClustersDelete(context.Background(), clusterId)
	_, res, err := cs.client.ClustersApi.ClustersDeleteExecute(req)
	return &Response{*res}, err
}
