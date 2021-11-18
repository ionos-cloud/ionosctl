package resources

import (
	"context"

	"github.com/fatih/structs"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
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
	List(params ListQueryParams) (Datacenters, *Response, error)
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

func (ds *dataCentersService) List(params ListQueryParams) (Datacenters, *Response, error) {
	req := ds.client.DataCenterApi.DatacentersGet(ds.context)
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
	}
	dcs, res, err := ds.client.DataCenterApi.DatacentersGetExecute(req)
	return Datacenters{dcs}, &Response{*res}, err
}

func (ds *dataCentersService) Get(datacenterId string) (*Datacenter, *Response, error) {
	req := ds.client.DataCenterApi.DatacentersFindById(ds.context, datacenterId)
	datacenter, res, err := ds.client.DataCenterApi.DatacentersFindByIdExecute(req)
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
	req := ds.client.DataCenterApi.DatacentersPost(ds.context).Datacenter(dc)
	datacenter, res, err := ds.client.DataCenterApi.DatacentersPostExecute(req)
	return &Datacenter{datacenter}, &Response{*res}, err
}

func (ds *dataCentersService) Update(datacenterId string, input DatacenterProperties) (*Datacenter, *Response, error) {
	req := ds.client.DataCenterApi.DatacentersPatch(ds.context, datacenterId).Datacenter(input.DatacenterProperties)
	datacenter, res, err := ds.client.DataCenterApi.DatacentersPatchExecute(req)
	return &Datacenter{datacenter}, &Response{*res}, err
}

func (ds *dataCentersService) Delete(datacenterId string) (*Response, error) {
	req := ds.client.DataCenterApi.DatacentersDelete(context.Background(), datacenterId)
	_, res, err := ds.client.DataCenterApi.DatacentersDeleteExecute(req)
	return &Response{*res}, err
}
