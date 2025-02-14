package resources

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	compute "github.com/ionos-cloud/sdk-go/v6"

	"github.com/fatih/structs"
)

type Request struct {
	compute.Request
}

type RequestStatus struct {
	compute.RequestStatus
}

type Requests struct {
	compute.Requests
}

// RequestsService is a wrapper around compute.Request
type RequestsService interface {
	List(params ListQueryParams) (Requests, *Response, error)
	Get(requestId string, params QueryParams) (*Request, *Response, error)
	GetStatus(requestId string) (*RequestStatus, *Response, error)
	Wait(requestId string) (*Response, error)
}

type requestsService struct {
	client  *compute.APIClient
	context context.Context
}

var _ RequestsService = &requestsService{}

func NewRequestService(client *client.Client, ctx context.Context) RequestsService {
	return &requestsService{
		client:  client.CloudClient,
		context: ctx,
	}
}

func (rs *requestsService) List(params ListQueryParams) (Requests, *Response, error) {
	req := rs.client.RequestsApi.RequestsGet(rs.context)
	if !structs.IsZero(params) {
		if params.Filters != nil {
			for k, v := range *params.Filters {
				for _, val := range v {
					req = req.Filter(k, val)
				}
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
	reqs, res, err := rs.client.RequestsApi.RequestsGetExecute(req)
	return Requests{reqs}, &Response{*res}, err
}

func (rs *requestsService) Get(requestId string, params QueryParams) (*Request, *Response, error) {
	req := rs.client.RequestsApi.RequestsFindById(rs.context, requestId)
	if !structs.IsZero(params) {
		if params.Depth != nil {
			req = req.Depth(*params.Depth)
		}
		if params.Pretty != nil {
			// Currently not implemented
			req = req.Pretty(*params.Pretty)
		}
	}
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
