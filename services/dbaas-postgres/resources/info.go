package resources

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"

	"github.com/ionos-cloud/sdk-go-bundle/products/dbaas/psql/v2"
)

type APIVersionList struct {
	Versions []psql.APIVersion
}

type APIVersion struct {
	psql.APIVersion
}

// InfosService is a wrapper around ionoscloud.APIVersion
type InfosService interface {
	List() (APIVersionList, *Response, error)
	Get() (APIVersion, *Response, error)
}

type infosService struct {
	client  *psql.APIClient
	context context.Context
}

var _ InfosService = &infosService{}

func NewInfosService(client *client.Client, ctx context.Context) InfosService {
	return &infosService{
		client:  client.PostgresClient,
		context: ctx,
	}
}

func (svc *infosService) List() (APIVersionList, *Response, error) {
	req := svc.client.MetadataApi.InfosVersionsGet(svc.context)
	versions, res, err := svc.client.MetadataApi.InfosVersionsGetExecute(req)
	return APIVersionList{versions}, &Response{*res}, err
}

func (svc *infosService) Get() (APIVersion, *Response, error) {
	req := svc.client.MetadataApi.InfosVersionGet(svc.context)
	versions, res, err := svc.client.MetadataApi.InfosVersionGetExecute(req)
	return APIVersion{versions}, &Response{*res}, err
}
