package resources

import (
	"context"

	sdkgo "github.com/ionos-cloud/sdk-go-autoscaling"
)

type ClusterLogs struct {
	sdkgo.ClusterLogs
}

// LogsService is a wrapper around ionoscloud.ClusterLogs
type LogsService interface {
	Get(clusterId string, start, end string, limit int32) (ClusterLogs, *Response, error)
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

func (svc *logsService) Get(clusterId string, start, end string, limit int32) (ClusterLogs, *Response, error) {
	req := svc.client.LogsApi.ClusterLogsGet(svc.context, clusterId)
	if start != "" {
		req = req.Start(start)
	}
	if end != "" {
		req = req.End(end)
	}
	if limit != 0 {
		req = req.Limit(limit)
	}
	logs, res, err := svc.client.LogsApi.ClusterLogsGetExecute(req)
	return ClusterLogs{logs}, &Response{*res}, err
}
