package resources

import (
	"context"

	client2 "github.com/ionos-cloud/ionosctl/v6/internal/client"
	sdkgo "github.com/ionos-cloud/sdk-go-bundle/products/containerregistry"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
)

type LocationsService interface {
	Get() (sdkgo.LocationsResponse, *shared.APIResponse, error)
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

func (svc *locationsService) Get() (sdkgo.LocationsResponse, *shared.APIResponse, error) {
	req := svc.client.LocationsApi.LocationsGet(svc.context)
	loc, res, err := svc.client.LocationsApi.LocationsGetExecute(req)
	return loc, res, err
}
