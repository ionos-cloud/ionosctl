package resources

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"

	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
)

type Location struct {
	ionoscloud.Location
}

type CpuArchitectureProperties struct {
	ionoscloud.CpuArchitectureProperties
}

type Locations struct {
	ionoscloud.Locations
}

// LocationsService is a wrapper around ionoscloud.Location
type LocationsService interface {
	List(params ListQueryParams) (Locations, *Response, error)
	GetByRegionAndLocationId(regionId, locationId string, params QueryParams) (*Location, *Response, error)
	GetByRegionId(regionId string, params QueryParams) (Locations, *Response, error)
}

type locationsService struct {
	client  *ionoscloud.APIClient
	context context.Context
}

var _ LocationsService = &locationsService{}

func NewLocationService(client *client.Client, ctx context.Context) LocationsService {
	return &locationsService{
		client:  client.CloudClient,
		context: ctx,
	}
}

func (s *locationsService) List(params ListQueryParams) (Locations, *Response, error) {
	req := s.client.LocationsApi.LocationsGet(s.context)
	locations, resp, err := s.client.LocationsApi.LocationsGetExecute(req)
	return Locations{locations}, &Response{*resp}, err
}

func (s *locationsService) GetByRegionAndLocationId(regionId, locationId string, params QueryParams) (*Location, *Response, error) {
	req := s.client.LocationsApi.LocationsFindByRegionIdAndId(s.context, regionId, locationId)
	loc, resp, err := s.client.LocationsApi.LocationsFindByRegionIdAndIdExecute(req)
	return &Location{loc}, &Response{*resp}, err
}

func (s *locationsService) GetByRegionId(regionId string, params QueryParams) (Locations, *Response, error) {
	req := s.client.LocationsApi.LocationsFindByRegionId(s.context, regionId)
	locs, resp, err := s.client.LocationsApi.LocationsFindByRegionIdExecute(req)
	return Locations{locs}, &Response{*resp}, err
}
