package resources

import (
	"context"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"time"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"

	sdkgo "github.com/ionos-cloud/sdk-go-dbaas-postgres"
)

type ClusterLogs struct {
	sdkgo.ClusterLogs
}

type LogsQueryParams struct {
	Direction          string
	Limit              int32
	StartTime, EndTime time.Time
}

// LogsService is a wrapper around ionoscloud.ClusterLogs
type LogsService interface {
	Get(clusterId string, queryParams *LogsQueryParams) (*ClusterLogs, *Response, error)
}

type logsService struct {
	client  *sdkgo.APIClient
	context context.Context
}

var _ LogsService = &logsService{}

func NewLogsService(client *client.Client, ctx context.Context) LogsService {
	return &logsService{
		client:  client.PostgresClient,
		context: ctx,
	}
}

func (svc *logsService) Get(clusterId string, queryParams *LogsQueryParams) (*ClusterLogs, *Response, error) {
	req := svc.client.LogsApi.ClusterLogsGet(svc.context, clusterId)
	if queryParams != nil {
		if !queryParams.StartTime.IsZero() {
			req = req.Start(queryParams.StartTime)
		}
		if !queryParams.EndTime.IsZero() {
			req = req.End(queryParams.EndTime)
		}
		if queryParams.Limit != 0 {
			req = req.Limit(queryParams.Limit)
		}
		if queryParams.Direction != "" {
			req = req.Direction(queryParams.Direction)
		}
	}
	logs, res, err := svc.client.LogsApi.ClusterLogsGetExecute(req)
	return &ClusterLogs{logs}, &Response{*res}, err
}
