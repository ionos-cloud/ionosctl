package resources

import (
	"context"

	client2 "github.com/ionos-cloud/ionosctl/v6/internal/client"
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

func NewLocationsService(client *client2.Client, ctx context.Context) LocationsService {
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
