package resources

import (
	"context"
	"time"

	sdkgo "github.com/ionos-cloud/sdk-go-autoscaling"
)

type ClusterLogs struct {
	sdkgo.ClusterLogs
}

// LogsService is a wrapper around ionoscloud.ClusterLogs
type LogsService interface {
	Get(clusterId string, start, end string, limit int32) (*ClusterLogs, *Response, error)
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

func (svc *logsService) Get(clusterId string, start, end string, limit int32) (*ClusterLogs, *Response, error) {
	req := svc.client.LogsApi.ClusterLogsGet(svc.context, clusterId)
	if start != "" {
		startFormat, err := time.Parse(time.RFC3339, start)
		if err != nil {
			return nil, nil, err
		}
		req = req.Start(startFormat)
	}
	if end != "" {
		endFormat, err := time.Parse(time.RFC3339, end)
		if err != nil {
			return nil, nil, err
		}
		req = req.End(endFormat)
	}
	if limit != 0 {
		req = req.Limit(limit)
	}
	logs, res, err := svc.client.LogsApi.ClusterLogsGetExecute(req)
	return &ClusterLogs{logs}, &Response{*res}, err
}
