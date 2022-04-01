package resources

import (
	"context"
	"strings"
	"time"

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
	Get(clusterId string, queryParams LogsQueryParams) (*ClusterLogs, *Response, error)
}

type logsService struct {
	client  *Client
	context context.Context
}

var _ LogsService = &logsService{}

func NewLogsService(client *Client, ctx context.Context) LogsService {
	return &logsService{
		client:  client,
		context: ctx,
	}
}

func (svc *logsService) Get(clusterId string, queryParams LogsQueryParams) (*ClusterLogs, *Response, error) {
	req := svc.client.LogsApi.ClusterLogsGet(svc.context, clusterId)
	if queryParams.Direction != "" {
		req = req.Direction(strings.ToUpper(queryParams.Direction))
	}
	if !queryParams.StartTime.IsZero() {
		req = req.Start(queryParams.StartTime)
	}
	if !queryParams.EndTime.IsZero() {
		req = req.End(queryParams.EndTime)
	}
	if queryParams.Limit != 0 {
		req = req.Limit(queryParams.Limit)
	}
	logs, res, err := svc.client.LogsApi.ClusterLogsGetExecute(req)
	return &ClusterLogs{logs}, &Response{*res}, err
}
