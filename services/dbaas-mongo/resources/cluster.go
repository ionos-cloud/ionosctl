package resources

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ionos-cloud/ionosctl/v6/pkg/config"

	sdkgo "github.com/ionos-cloud/sdk-go-dbaas-mongo"
)

type LogsQueryParams struct {
	Direction          string
	Limit              int32
	StartTime, EndTime time.Time
}

type ClustersService interface {
	List(filterName string, limit, offset *int32) (sdkgo.ClusterList, *sdkgo.APIResponse, error)
	Get(clusterId string) (sdkgo.ClusterResponse, *sdkgo.APIResponse, error)
	Create(input sdkgo.CreateClusterRequest) (sdkgo.ClusterResponse, *sdkgo.APIResponse, error)
	Update(id string, input sdkgo.PatchClusterRequest) (sdkgo.ClusterResponse, *sdkgo.APIResponse, error)
	Delete(clusterId string) (*sdkgo.APIResponse, error)
	DeleteAll(name string) (*sdkgo.APIResponse, error)
	Restore(clusterId, snapshotId string) (*sdkgo.APIResponse, error)
	SnapshotsList(clusterId string, limit, offset *int32) (sdkgo.SnapshotList, *sdkgo.APIResponse, error)
	LogsList(clusterId string, direction *string, limit *int32, start, end *time.Time) (sdkgo.ClusterLogs, *sdkgo.APIResponse, error)
}

type clustersService struct {
	client  *sdkgo.APIClient
	context context.Context
}

var _ ClustersService = &clustersService{}

func NewClustersService(client *config.Client, ctx context.Context) ClustersService {
	return &clustersService{
		client:  client.MongoClient,
		context: ctx,
	}
}

func (svc *clustersService) List(filterName string, limit, offset *int32) (sdkgo.ClusterList, *sdkgo.APIResponse, error) {
	req := svc.client.ClustersApi.ClustersGet(svc.context)
	if offset != nil {
		fmt.Printf("Running with offset: %d\n", *offset)
		req = req.Offset(*offset)
	}
	if limit != nil {
		fmt.Printf("Running with limit: %d\n", *limit)
		req = req.Limit(*limit)
	}
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
	return svc.client.ClustersApi.ClustersPost(svc.context).CreateClusterRequest(input).Execute()
}

func (svc *clustersService) Update(id string, input sdkgo.PatchClusterRequest) (sdkgo.ClusterResponse, *sdkgo.APIResponse, error) {
	return svc.client.ClustersApi.ClustersPatch(svc.context, id).PatchClusterRequest(input).Execute()
}

func (svc *clustersService) Delete(clusterId string) (*sdkgo.APIResponse, error) {
	req := svc.client.ClustersApi.ClustersDelete(svc.context, clusterId)
	_, res, err := svc.client.ClustersApi.ClustersDeleteExecute(req)
	return res, err
}

func (svc *clustersService) DeleteAll(filterName string) (*sdkgo.APIResponse, error) {
	ls, _, err := svc.List(filterName, nil, nil)

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

func (svc *clustersService) SnapshotsList(clusterId string, offset, limit *int32) (sdkgo.SnapshotList, *sdkgo.APIResponse, error) {
	req := svc.client.SnapshotsApi.ClustersSnapshotsGet(svc.context, clusterId)
	if offset != nil {
		req = req.Offset(*offset)
	}
	if limit != nil {
		req = req.Limit(*limit)
	}
	snapshots, res, err := svc.client.SnapshotsApi.ClustersSnapshotsGetExecute(req)
	if err != nil {
		fmt.Println(err)
	}
	return snapshots, res, err
}

func (svc *clustersService) LogsList(clusterId string, direction *string, limit *int32, start, end *time.Time) (sdkgo.ClusterLogs, *sdkgo.APIResponse, error) {
	req := svc.client.LogsApi.ClustersLogsGet(svc.context, clusterId)
	if direction != nil {
		req = req.Direction(*direction)
	}
	if limit != nil {
		req = req.Limit(*limit)
	}
	if start != nil {
		req = req.Start(*start)
	}
	if end != nil {
		req = req.End(*end)
	}
	logs, res, err := svc.client.LogsApi.ClustersLogsGetExecute(req)
	return logs, res, err
}
