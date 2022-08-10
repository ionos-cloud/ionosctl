package resources

import (
	"context"
)

type VersionsList []string

// VersionsService is a wrapper around dp.Cluster
type VersionsService interface {
	List() ([]string, *Response, error)
}

type versionsService struct {
	client  *Client
	context context.Context
}

var _ VersionsService = &versionsService{}

func NewVersionsService(client *Client, ctx context.Context) VersionsService {
	return &versionsService{
		client:  client,
		context: ctx,
	}
}

func (svc *versionsService) List() ([]string, *Response, error) {
	req := svc.client.DataPlatformMetaDataApi.VersionsGet(svc.context)
	versionsResponse, res, err := svc.client.DataPlatformMetaDataApi.VersionsGetExecute(req)
	return versionsResponse, &Response{*res}, err
}
