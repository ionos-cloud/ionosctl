package resources

import (
	"context"

	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
)

type Request struct {
	ionoscloud.Request
}

type RequestStatus struct {
	ionoscloud.RequestStatus
}

type Requests struct {
	ionoscloud.Requests
}

// RequestsService is a wrapper around ionoscloud.Request
type RequestsService interface {
	List(params ListQueryParams) (Requests, *Response, error)
	Get(requestId string) (*Request, *Response, error)
	GetStatus(requestId string) (*RequestStatus, *Response, error)
	Wait(requestId string) (*Response, error)
}

type requestsService struct {
	client  *Client
	context context.Context
}

var _ RequestsService = &requestsService{}

func NewRequestService(client *Client, ctx context.Context) RequestsService {
	return &requestsService{
		client:  client,
		context: ctx,
	}
}

func (rs *requestsService) List(params ListQueryParams) (Requests, *Response, error) {
	req := rs.client.RequestsApi.RequestsGet(rs.context)
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
	reqs, res, err := rs.client.RequestsApi.RequestsGetExecute(req)
	return Requests{reqs}, &Response{*res}, err
}

func (rs *requestsService) Get(requestId string) (*Request, *Response, error) {
	req := rs.client.RequestsApi.RequestsFindById(rs.context, requestId)
	reqs, res, err := rs.client.RequestsApi.RequestsFindByIdExecute(req)
	return &Request{reqs}, &Response{*res}, err
}

func (rs *requestsService) GetStatus(requestId string) (*RequestStatus, *Response, error) {
	req := rs.client.RequestsApi.RequestsStatusGet(rs.context, requestId)
	reqs, res, err := rs.client.RequestsApi.RequestsStatusGetExecute(req)
	return &RequestStatus{reqs}, &Response{*res}, err
}

func (rs *requestsService) Wait(path string) (*Response, error) {
	res, err := rs.client.WaitForRequest(rs.context, path)
	return &Response{*res}, err
}
