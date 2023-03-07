package resources

import (
	"context"
	"github.com/ionos-cloud/ionosctl/v6/pkg/config"

	"github.com/fatih/structs"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
)

type FlowLog struct {
	ionoscloud.FlowLog
}

type FlowLogPut struct {
	ionoscloud.FlowLogPut
}

type FlowLogProperties struct {
	ionoscloud.FlowLogProperties
}

type FlowLogs struct {
	ionoscloud.FlowLogs
}

// FlowLogsService is a wrapper around ionoscloud.FlowLog
type FlowLogsService interface {
	List(datacenterId, serverId, nicId string, params ListQueryParams) (FlowLogs, *Response, error)
	Get(datacenterId, serverId, nicId, flowLogId string, params QueryParams) (*FlowLog, *Response, error)
	Create(datacenterId, serverId, nicId string, input FlowLog, params QueryParams) (*FlowLog, *Response, error)
	Update(datacenterId, serverId, nicId, flowlogId string, input FlowLogPut, params QueryParams) (*FlowLog, *Response, error)
	Delete(datacenterId, serverId, nicId, flowLogId string, params QueryParams) (*Response, error)
}

type flowLogsService struct {
	client  *ionoscloud.APIClient
	context context.Context
}

var _ FlowLogsService = &flowLogsService{}

func NewFlowLogService(client *config.Client, ctx context.Context) FlowLogsService {
	return &flowLogsService{
		client:  client.CloudClient,
		context: ctx,
	}
}

func (svc *flowLogsService) List(datacenterId, serverId, nicId string, params ListQueryParams) (FlowLogs, *Response, error) {
	req := svc.client.FlowLogsApi.DatacentersServersNicsFlowlogsGet(svc.context, datacenterId, serverId, nicId)
	if !structs.IsZero(params) {
		if params.Filters != nil {
			for k, v := range *params.Filters {
				req = req.Filter(k, v)
			}
		}
		if params.OrderBy != nil {
			req = req.OrderBy(*params.OrderBy)
		}
		if params.MaxResults != nil {
			req = req.MaxResults(*params.MaxResults)
		}
		if !structs.IsZero(params.QueryParams) {
			if params.QueryParams.Depth != nil {
				req = req.Depth(*params.QueryParams.Depth)
			}
			if params.QueryParams.Pretty != nil {
				// Currently not implemented
				req = req.Pretty(*params.QueryParams.Pretty)
			}
		}
	}
	flowlogs, resp, err := svc.client.FlowLogsApi.DatacentersServersNicsFlowlogsGetExecute(req)
	return FlowLogs{flowlogs}, &Response{*resp}, err
}

func (svc *flowLogsService) Get(datacenterId, serverId, nicId, flowLogId string, params QueryParams) (*FlowLog, *Response, error) {
	req := svc.client.FlowLogsApi.DatacentersServersNicsFlowlogsFindById(svc.context, datacenterId, serverId, nicId, flowLogId)
	flowlog, resp, err := svc.client.FlowLogsApi.DatacentersServersNicsFlowlogsFindByIdExecute(req)
	return &FlowLog{flowlog}, &Response{*resp}, err
}

func (svc *flowLogsService) Create(datacenterId, serverId, nicId string, input FlowLog, params QueryParams) (*FlowLog, *Response, error) {
	req := svc.client.FlowLogsApi.DatacentersServersNicsFlowlogsPost(svc.context, datacenterId, serverId, nicId).Flowlog(input.FlowLog)
	flowlog, resp, err := svc.client.FlowLogsApi.DatacentersServersNicsFlowlogsPostExecute(req)
	return &FlowLog{flowlog}, &Response{*resp}, err
}

func (svc *flowLogsService) Update(datacenterId, serverId, nicId, flowlogId string, input FlowLogPut, params QueryParams) (*FlowLog, *Response, error) {
	req := svc.client.FlowLogsApi.DatacentersServersNicsFlowlogsPut(svc.context, datacenterId, serverId, nicId, flowlogId).Flowlog(input.FlowLogPut)
	flowlog, resp, err := svc.client.FlowLogsApi.DatacentersServersNicsFlowlogsPutExecute(req)
	return &FlowLog{flowlog}, &Response{*resp}, err
}

func (svc *flowLogsService) Delete(datacenterId, serverId, nicId, flowLogId string, params QueryParams) (*Response, error) {
	req := svc.client.FlowLogsApi.DatacentersServersNicsFlowlogsDelete(svc.context, datacenterId, serverId, nicId, flowLogId)
	resp, err := svc.client.FlowLogsApi.DatacentersServersNicsFlowlogsDeleteExecute(req)
	return &Response{*resp}, err
}
