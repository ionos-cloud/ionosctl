package resources

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"

	sdkgo "github.com/ionos-cloud/sdk-go-dbaas-mongo"
)

type ApiMetadataService interface {
	List() ([]sdkgo.APIVersion, *sdkgo.APIResponse, error)
}

type apiMetadataService struct {
	client  *sdkgo.APIClient
	context context.Context
}

var _ ApiMetadataService = &apiMetadataService{}

func NewApiMetadataService(client *client.Client, ctx context.Context) ApiMetadataService {
	return &apiMetadataService{
		client:  client.MongoClient,
		context: ctx,
	}
}

func (svc apiMetadataService) List() ([]sdkgo.APIVersion, *sdkgo.APIResponse, error) {
	return svc.client.MetadataApi.InfosVersionsGet(svc.context).Execute()
}
