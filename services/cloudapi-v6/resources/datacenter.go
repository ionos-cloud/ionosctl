package v6

import (
	"context"

	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
)

type Datacenter struct {
	ionoscloud.Datacenter
}

type DatacenterProperties struct {
	ionoscloud.DatacenterProperties
}

type Datacenters struct {
	ionoscloud.Datacenters
}

type Response struct {
	ionoscloud.APIResponse
}

// DatacentersService is a wrapper around ionoscloud.Datacenter
type DatacentersService interface {
	List() (Datacenters, *Response, error)
	Get(datacenterId string) (*Datacenter, *Response, error)
	Create(name, description, region string) (*Datacenter, *Response, error)
	Update(datacenterId string, input DatacenterProperties) (*Datacenter, *Response, error)
	Delete(datacenterId string) (*Response, error)
}

type dataCentersService struct {
	client  *Client
	context context.Context
}

var _ DatacentersService = &dataCentersService{}

func NewDataCenterService(client *Client, ctx context.Context) DatacentersService {
	return &dataCentersService{
		client:  client,
		context: ctx,
	}
}

func (ds *dataCentersService) List() (Datacenters, *Response, error) {
	req := ds.client.DataCentersApi.DatacentersGet(ds.context)
	dcs, res, err := ds.client.DataCentersApi.DatacentersGetExecute(req)
	return Datacenters{dcs}, &Response{*res}, err
}

func (ds *dataCentersService) Get(datacenterId string) (*Datacenter, *Response, error) {
	req := ds.client.DataCentersApi.DatacentersFindById(ds.context, datacenterId)
	datacenter, res, err := ds.client.DataCentersApi.DatacentersFindByIdExecute(req)
	return &Datacenter{datacenter}, &Response{*res}, err
}

func (ds *dataCentersService) Create(name, description, region string) (*Datacenter, *Response, error) {
	dc := ionoscloud.Datacenter{
		Properties: &ionoscloud.DatacenterProperties{
			Name:        &name,
			Description: &description,
			Location:    &region,
		},
	}
	req := ds.client.DataCentersApi.DatacentersPost(ds.context).Datacenter(dc)
	datacenter, res, err := ds.client.DataCentersApi.DatacentersPostExecute(req)
	return &Datacenter{datacenter}, &Response{*res}, err
}

func (ds *dataCentersService) Update(datacenterId string, input DatacenterProperties) (*Datacenter, *Response, error) {
	req := ds.client.DataCentersApi.DatacentersPatch(ds.context, datacenterId).Datacenter(input.DatacenterProperties)
	datacenter, res, err := ds.client.DataCentersApi.DatacentersPatchExecute(req)
	return &Datacenter{datacenter}, &Response{*res}, err
}

func (ds *dataCentersService) Delete(datacenterId string) (*Response, error) {
	req := ds.client.DataCentersApi.DatacentersDelete(context.Background(), datacenterId)
	res, err := ds.client.DataCentersApi.DatacentersDeleteExecute(req)
	return &Response{*res}, err
}
