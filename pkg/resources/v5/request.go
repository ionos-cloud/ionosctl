package v5

import (
	"context"

	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
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
	List() (Requests, *Response, error)
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

func (rs *requestsService) List() (Requests, *Response, error) {
	req := rs.client.RequestApi.RequestsGet(rs.context)
	reqs, res, err := rs.client.RequestApi.RequestsGetExecute(req)
	return Requests{reqs}, &Response{*res}, err
}

func (rs *requestsService) Get(requestId string) (*Request, *Response, error) {
	req := rs.client.RequestApi.RequestsFindById(rs.context, requestId)
	reqs, res, err := rs.client.RequestApi.RequestsFindByIdExecute(req)
	return &Request{reqs}, &Response{*res}, err
}

func (rs *requestsService) GetStatus(requestId string) (*RequestStatus, *Response, error) {
	req := rs.client.RequestApi.RequestsStatusGet(rs.context, requestId)
	reqs, res, err := rs.client.RequestApi.RequestsStatusGetExecute(req)
	return &RequestStatus{reqs}, &Response{*res}, err
}

func (rs *requestsService) Wait(path string) (*Response, error) {
	res, err := rs.client.WaitForRequest(rs.context, path)
	return &Response{*res}, err
}
