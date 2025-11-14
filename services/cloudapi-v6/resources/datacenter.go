package resources

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"

	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
)

type Datacenter struct {
	ionoscloud.Datacenter
}

type DatacenterProperties struct {
	ionoscloud.DatacenterProperties
}

type DatacenterPropertiesPut struct {
	ionoscloud.DatacenterPropertiesPut
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
	Get(datacenterId string, queryParams QueryParams) (*Datacenter, *Response, error)
	Create(name, description, region string, queryParams QueryParams) (*Datacenter, *Response, error)
	Update(datacenterId string, input DatacenterPropertiesPut, queryParams QueryParams) (*Datacenter, *Response, error)
	Delete(datacenterId string, queryParams QueryParams) (*Response, error)
}

type dataCentersService struct {
	client  *ionoscloud.APIClient
	context context.Context
}

var _ DatacentersService = &dataCentersService{}

func NewDataCenterService(client *client.Client, ctx context.Context) DatacentersService {
	return &dataCentersService{
		client:  client.CloudClient,
		context: ctx,
	}
}

// func NewDataCenterServices(client *client2.Client, ctx context.Context) DatacentersService {
//	return &dataCentersService{
//		client:  client,
//		context: ctx,
//	}
// }

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
	dc := ionoscloud.DatacenterPost{
		Properties: &ionoscloud.DatacenterPropertiesPost{
			Name:        &name,
			Description: &description,
			Location:    &region,
		},
	}
	req := ds.client.DataCentersApi.DatacentersPost(ds.context).Datacenter(dc)
	datacenter, res, err := ds.client.DataCentersApi.DatacentersPostExecute(req)
	return &Datacenter{datacenter}, &Response{*res}, err
}

func (ds *dataCentersService) Update(datacenterId string, input DatacenterPropertiesPut) (*Datacenter, *Response, error) {
	req := ds.client.DataCentersApi.DatacentersPatch(ds.context, datacenterId).Datacenter(input.DatacenterPropertiesPut)
	datacenter, res, err := ds.client.DataCentersApi.DatacentersPatchExecute(req)
	return &Datacenter{datacenter}, &Response{*res}, err
}

func (ds *dataCentersService) Delete(datacenterId string) (*Response, error) {
	req := ds.client.DataCentersApi.DatacentersDelete(context.Background(), datacenterId)
	res, err := ds.client.DataCentersApi.DatacentersDeleteExecute(req)
	return &Response{*res}, err
}
