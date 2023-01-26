package resources

import (
	"context"
	"github.com/ionos-cloud/ionosctl/pkg/config"

	"github.com/fatih/structs"
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

func NewLocationService(client *config.Client, ctx context.Context) LocationsService {
	return &locationsService{
		client:  client.CloudClient,
		context: ctx,
	}
}

func (s *locationsService) List(params ListQueryParams) (Locations, *Response, error) {
	req := s.client.LocationsApi.LocationsGet(s.context)
	if !structs.IsZero(params) {
		if params.Filters != nil {
			for k, v := range *params.Filters {
				req = req.Filter(k, v)
			}
		}
		if params.OrderBy != nil {
			req = req.OrderBy(*params.OrderBy)
		}
		if params.MaxResults != nil {
			req = req.MaxResults(*params.MaxResults)
		}
		if !structs.IsZero(params.QueryParams) {
			if params.QueryParams.Depth != nil {
				req = req.Depth(*params.QueryParams.Depth)
			}
			if params.QueryParams.Pretty != nil {
				// Currently not implemented
				req = req.Pretty(*params.QueryParams.Pretty)
			}
		}
	}
	locations, resp, err := s.client.LocationsApi.LocationsGetExecute(req)
	return Locations{locations}, &Response{*resp}, err
}

func (s *locationsService) GetByRegionAndLocationId(regionId, locationId string, params QueryParams) (*Location, *Response, error) {
	req := s.client.LocationsApi.LocationsFindByRegionIdAndId(s.context, regionId, locationId)
	if !structs.IsZero(params) {
		if params.Depth != nil {
			req = req.Depth(*params.Depth)
		}
		if params.Pretty != nil {
			// Currently not implemented
			req = req.Pretty(*params.Pretty)
		}
	}
	loc, resp, err := s.client.LocationsApi.LocationsFindByRegionIdAndIdExecute(req)
	return &Location{loc}, &Response{*resp}, err
}

func (s *locationsService) GetByRegionId(regionId string, params QueryParams) (Locations, *Response, error) {
	req := s.client.LocationsApi.LocationsFindByRegionId(s.context, regionId)
	if !structs.IsZero(params) {
		if params.Depth != nil {
			req = req.Depth(*params.Depth)
		}
		if params.Pretty != nil {
			// Currently not implemented
			req = req.Pretty(*params.Pretty)
		}
	}
	locs, resp, err := s.client.LocationsApi.LocationsFindByRegionIdExecute(req)
	return Locations{locs}, &Response{*resp}, err
}
