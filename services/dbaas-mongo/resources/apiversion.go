package resources

import (
	"context"
	sdkgo "github.com/ionos-cloud/sdk-go-dbaas-mongo"
)

type ApiMetadataService interface {
	List() ([]sdkgo.APIVersion, *sdkgo.APIResponse, error)
}

type apiMetadataService struct {
	client  *Client
	context context.Context
}

var _ ApiMetadataService = &apiMetadataService{}

func NewApiMetadataService(client *Client, ctx context.Context) ApiMetadataService {
	return &apiMetadataService{
		client:  client,
		context: ctx,
	}
}

func (svc apiMetadataService) List() ([]sdkgo.APIVersion, *sdkgo.APIResponse, error) {
	return svc.client.APIClient.MetadataApi.InfosVersionsGet(svc.context).Execute()
}
