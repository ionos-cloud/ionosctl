package resources

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	compute "github.com/ionos-cloud/sdk-go/v6"

	"github.com/fatih/structs"
)

type Datacenter struct {
	compute.Datacenter
}

type DatacenterProperties struct {
	compute.DatacenterProperties
}

type Datacenters struct {
	compute.Datacenters
}

type Response struct {
	compute.APIResponse
}

// DatacentersService is a wrapper around compute.Datacenter
type DatacentersService interface {
	List(params ListQueryParams) (Datacenters, *Response, error)
	Get(datacenterId string, queryParams QueryParams) (*Datacenter, *Response, error)
	Create(name, description, region string, queryParams QueryParams) (*Datacenter, *Response, error)
	Update(datacenterId string, input DatacenterProperties, queryParams QueryParams) (*Datacenter, *Response, error)
	Delete(datacenterId string, queryParams QueryParams) (*Response, error)
}

type dataCentersService struct {
	client  *compute.APIClient
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

func (ds *dataCentersService) List(params ListQueryParams) (Datacenters, *Response, error) {
	req := ds.client.DataCentersApi.DatacentersGet(ds.context)
	if !structs.IsZero(params) {
		if params.Filters != nil {
			for k, v := range *params.Filters {
				for _, val := range v {
					req = req.Filter(k, val)
				}
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
	dcs, res, err := ds.client.DataCentersApi.DatacentersGetExecute(req)
	return Datacenters{dcs}, &Response{*res}, err
}

func (ds *dataCentersService) Get(datacenterId string, params QueryParams) (*Datacenter, *Response, error) {
	req := ds.client.DataCentersApi.DatacentersFindById(ds.context, datacenterId)
	if !structs.IsZero(params) {
		if params.Depth != nil {
			req = req.Depth(*params.Depth)
		}
		if params.Pretty != nil {
			// Currently not implemented
			req = req.Pretty(*params.Pretty)
		}
	}
	datacenter, res, err := ds.client.DataCentersApi.DatacentersFindByIdExecute(req)
	return &Datacenter{datacenter}, &Response{*res}, err
}

func (ds *dataCentersService) Create(name, description, region string, params QueryParams) (*Datacenter, *Response, error) {
	dc := compute.Datacenter{
		Properties: &compute.DatacenterProperties{
			Name:        &name,
			Description: &description,
			Location:    &region,
		},
	}
	req := ds.client.DataCentersApi.DatacentersPost(ds.context).Datacenter(dc)
	if !structs.IsZero(params) {
		if params.Depth != nil {
			req = req.Depth(*params.Depth)
		}
		if params.Pretty != nil {
			// Currently not implemented
			req = req.Pretty(*params.Pretty)
		}
	}
	datacenter, res, err := ds.client.DataCentersApi.DatacentersPostExecute(req)
	return &Datacenter{datacenter}, &Response{*res}, err
}

func (ds *dataCentersService) Update(datacenterId string, input DatacenterProperties, params QueryParams) (*Datacenter, *Response, error) {
	req := ds.client.DataCentersApi.DatacentersPatch(ds.context, datacenterId).Datacenter(input.DatacenterProperties)
	if !structs.IsZero(params) {
		if params.Depth != nil {
			req = req.Depth(*params.Depth)
		}
		if params.Pretty != nil {
			// Currently not implemented
			req = req.Pretty(*params.Pretty)
		}
	}
	datacenter, res, err := ds.client.DataCentersApi.DatacentersPatchExecute(req)
	return &Datacenter{datacenter}, &Response{*res}, err
}

func (ds *dataCentersService) Delete(datacenterId string, params QueryParams) (*Response, error) {
	req := ds.client.DataCentersApi.DatacentersDelete(context.Background(), datacenterId)
	if !structs.IsZero(params) {
		if params.Depth != nil {
			req = req.Depth(*params.Depth)
		}
		if params.Pretty != nil {
			// Currently not implemented
			req = req.Pretty(*params.Pretty)
		}
	}
	res, err := ds.client.DataCentersApi.DatacentersDeleteExecute(req)
	return &Response{*res}, err
}
