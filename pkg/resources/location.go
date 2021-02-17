package resources

import (
	"context"

	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
)

type Location struct {
	ionoscloud.Location
}

type Locations struct {
	ionoscloud.Locations
}

// LocationsService is a wrapper around ionoscloud.Location
type LocationsService interface {
	List() (Locations, *Response, error)
}

type locationsService struct {
	client  *Client
	context context.Context
}

var _ LocationsService = &locationsService{}

func NewLocationService(client *Client, ctx context.Context) LocationsService {
	return &locationsService{
		client:  client,
		context: ctx,
	}
}

func (s *locationsService) List() (Locations, *Response, error) {
	req := s.client.LocationApi.LocationsGet(s.context)
	locations, resp, err := s.client.LocationApi.LocationsGetExecute(req)
	return Locations{locations}, &Response{*resp}, err
}
