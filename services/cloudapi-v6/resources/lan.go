package resources

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/sdk-go-bundle/products/compute/v2"

	"github.com/fatih/structs"
)

type Lan struct {
	compute.Lan
}

type IpFailover struct {
	compute.IPFailover
}

type LanProperties struct {
	compute.LanProperties
}

type LanPost struct {
	compute.LanPost
}

type Lans struct {
	compute.Lans
}

// LansService is a wrapper around compute.Lan
type LansService interface {
	List(datacenterId string, params ListQueryParams) (Lans, *Response, error)
	Get(datacenterId, lanId string, params QueryParams) (*Lan, *Response, error)
	Create(datacenterId string, input LanPost, params QueryParams) (*LanPost, *Response, error)
	Update(datacenterId, lanId string, input LanProperties, params QueryParams) (*Lan, *Response, error)
	Delete(datacenterId, lanId string, params QueryParams) (*Response, error)
}

type lansService struct {
	client  *compute.APIClient
	context context.Context
}

var _ LansService = &lansService{}

func NewLanService(client *client.Client, ctx context.Context) LansService {
	return &lansService{
		client:  client.CloudClient,
		context: ctx,
	}
}

func (ls *lansService) List(datacenterId string, params ListQueryParams) (Lans, *Response, error) {
	req := ls.client.LANsApi.DatacentersLansGet(ls.context, datacenterId)
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
	}
	lans, resp, err := ls.client.LANsApi.DatacentersLansGetExecute(req)
	return Lans{lans}, &Response{*resp}, err
}

func (ls *lansService) Get(datacenterId, lanId string, params QueryParams) (*Lan, *Response, error) {
	req := ls.client.LANsApi.DatacentersLansFindById(ls.context, datacenterId, lanId)
	if !structs.IsZero(params) {
		if params.Depth != nil {
			req = req.Depth(*params.Depth)
		}
		if params.Pretty != nil {
			// Currently not implemented
			req = req.Pretty(*params.Pretty)
		}
	}
	lan, resp, err := ls.client.LANsApi.DatacentersLansFindByIdExecute(req)
	return &Lan{lan}, &Response{*resp}, err
}

func (ls *lansService) Create(datacenterId string, input LanPost, params QueryParams) (*LanPost, *Response, error) {
	req := ls.client.LANsApi.DatacentersLansPost(ls.context, datacenterId).Lan(input.LanPost)
	if !structs.IsZero(params) {
		if params.Depth != nil {
			req = req.Depth(*params.Depth)
		}
		if params.Pretty != nil {
			// Currently not implemented
			req = req.Pretty(*params.Pretty)
		}
	}
	lan, resp, err := ls.client.LANsApi.DatacentersLansPostExecute(req)
	return &LanPost{lan}, &Response{*resp}, err
}

func (ls *lansService) Update(datacenterId, lanId string, input LanProperties, params QueryParams) (*Lan, *Response, error) {
	req := ls.client.LANsApi.DatacentersLansPatch(ls.context, datacenterId, lanId).Lan(input.LanProperties)
	if !structs.IsZero(params) {
		if params.Depth != nil {
			req = req.Depth(*params.Depth)
		}
		if params.Pretty != nil {
			// Currently not implemented
			req = req.Pretty(*params.Pretty)
		}
	}
	lan, resp, err := ls.client.LANsApi.DatacentersLansPatchExecute(req)
	return &Lan{lan}, &Response{*resp}, err
}

func (ls *lansService) Delete(datacenterId, lanId string, params QueryParams) (*Response, error) {
	req := ls.client.LANsApi.DatacentersLansDelete(ls.context, datacenterId, lanId)
	if !structs.IsZero(params) {
		if params.Depth != nil {
			req = req.Depth(*params.Depth)
		}
		if params.Pretty != nil {
			// Currently not implemented
			req = req.Pretty(*params.Pretty)
		}
	}
	resp, err := ls.client.LANsApi.DatacentersLansDeleteExecute(req)
	return &Response{*resp}, err
}
