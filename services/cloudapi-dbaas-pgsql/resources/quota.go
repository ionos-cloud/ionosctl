package resources

import (
	"context"
	sdkgo "github.com/ionos-cloud/sdk-go-autoscaling"
)

type QuotaList struct {
	sdkgo.QuotaList
}

// QuotasService is a wrapper around ionoscloud.QuotaList
type QuotasService interface {
	Get() (QuotaList, *Response, error)
}

type quotasService struct {
	client  *Client
	context context.Context
}

var _ QuotasService = &quotasService{}

func NewQuotasService(client *Client, ctx context.Context) QuotasService {
	return &quotasService{
		client:  client,
		context: ctx,
	}
}

func (svc *quotasService) Get() (QuotaList, *Response, error) {
	req := svc.client.QuotaApi.QuotaGet(svc.context)
	quotas, res, err := svc.client.QuotaApi.QuotaGetExecute(req)
	return QuotaList{quotas}, &Response{*res}, err
}
