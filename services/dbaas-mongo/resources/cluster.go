package resources

import (
	"context"
	"errors"
	"fmt"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"time"

	sdkgo "github.com/ionos-cloud/sdk-go-dbaas-mongo"
)

type LogsQueryParams struct {
	Direction          *string
	Limit              *int32
	StartTime, EndTime *time.Time
}

type ClustersService interface {
	List(filterName string) (sdkgo.ClusterList, *sdkgo.APIResponse, error)
	Get(clusterId string) (sdkgo.ClusterResponse, *sdkgo.APIResponse, error)
	Create(input sdkgo.CreateClusterRequest) (sdkgo.ClusterResponse, *sdkgo.APIResponse, error)
	Delete(clusterId string) (*sdkgo.APIResponse, error)
	DeleteAll(name string) (*sdkgo.APIResponse, error)
	Restore(clusterId, snapshotId string) (*sdkgo.APIResponse, error)
	SnapshotsList(clusterId string) (sdkgo.SnapshotList, *sdkgo.APIResponse, error)
	LogsList(clusterId string, logsQueryParams LogsQueryParams) (sdkgo.ClusterLogs, *sdkgo.APIResponse, error)
}

type clustersService struct {
	client  *Client
	context context.Context
}

var _ ClustersService = &clustersService{}

func NewClustersService(client *config.Client, ctx context.Context) ClustersService {
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

func (svc *clustersService) DeleteAll(filterName string) (*sdkgo.APIResponse, error) {
	ls, _, err := svc.List(filterName)

	if err != nil {
		return nil, errors.New("deletion of all clusters failed early: " + err.Error())
	}

	var res *sdkgo.APIResponse
	for _, c := range *ls.GetItems() {
		req := svc.client.ClustersApi.ClustersDelete(svc.context, *c.GetId())
		_, res, err = svc.client.ClustersApi.ClustersDeleteExecute(req)
		if err != nil {
			return nil, errors.New("deletion of all clusters failed early: " + err.Error())
		}
	}

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
	if err != nil {
		fmt.Println(err)
	}
	return snapshots, res, err
}

// Yuck! Reach out if you have a better solution.
func (q LogsQueryParams) applyToRequest(req sdkgo.ApiClustersLogsGetRequest) sdkgo.ApiClustersLogsGetRequest {
	if q.StartTime != nil {
		req = req.Start(*q.StartTime)
	}
	if q.EndTime != nil {
		req = req.Start(*q.EndTime)
	}
	if q.StartTime != nil {
		req = req.Direction(*q.Direction)
	}
	if q.StartTime != nil {
		req = req.Limit(*q.Limit)
	}
	return req
}

func (svc *clustersService) LogsList(clusterId string, q LogsQueryParams) (sdkgo.ClusterLogs, *sdkgo.APIResponse, error) {
	req := svc.client.LogsApi.ClustersLogsGet(svc.context, clusterId)
	logs, res, err := svc.client.LogsApi.ClustersLogsGetExecute(q.applyToRequest(req))
	return logs, res, err
}
