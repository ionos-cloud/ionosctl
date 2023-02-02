package resources

import (
	"context"
	"github.com/ionos-cloud/ionosctl/pkg/config"

	sdkgo "github.com/ionos-cloud/sdk-go-dbaas-postgres"
)

type APIVersionList struct {
	Versions []sdkgo.APIVersion
}

type APIVersion struct {
	sdkgo.APIVersion
}

// InfosService is a wrapper around ionoscloud.APIVersion
type InfosService interface {
	List() (APIVersionList, *Response, error)
	Get() (APIVersion, *Response, error)
}

type infosService struct {
	client  *sdkgo.APIClient
	context context.Context
}

var _ InfosService = &infosService{}

func NewInfosService(client *config.Client, ctx context.Context) InfosService {
	return &infosService{
		client:  client.DbaasClient,
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
