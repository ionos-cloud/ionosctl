package resources

import (
	"context"
	"github.com/ionos-cloud/ionosctl/v6/pkg/config"
	sdkgo "github.com/ionos-cloud/sdk-go-container-registry"
)

type LocationsService interface {
	Get() (sdkgo.LocationsResponse, *sdkgo.APIResponse, error)
}

type locationsService struct {
	client  *sdkgo.APIClient
	context context.Context
}

var _ LocationsService = &locationsService{}

func NewLocationsService(client *config.Client, ctx context.Context) LocationsService {
	return &locationsService{
		client:  client.RegistryClient,
		context: ctx,
	}
}

func (svc *locationsService) Get() (sdkgo.LocationsResponse, *sdkgo.APIResponse, error) {
	req := svc.client.LocationsApi.LocationsGet(svc.context)
	loc, res, err := svc.client.LocationsApi.LocationsGetExecute(req)
	return loc, res, err
}
