package resources

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"

	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
)

type Lan struct {
	ionoscloud.Lan
}

type IpFailover struct {
	ionoscloud.IPFailover
}

type LanProperties struct {
	ionoscloud.LanProperties
}

type LanPost struct {
	ionoscloud.Lan
}

type Lans struct {
	ionoscloud.Lans
}

// LansService is a wrapper around ionoscloud.Lan
type LansService interface {
	List(datacenterId string, params ListQueryParams) (Lans, *Response, error)
	Get(datacenterId, lanId string) (*Lan, *Response, error)
	Create(datacenterId string, input LanPost) (*LanPost, *Response, error)
	Update(datacenterId, lanId string, input LanProperties) (*Lan, *Response, error)
	Delete(datacenterId, lanId string) (*Response, error)
}

type lansService struct {
	client  *ionoscloud.APIClient
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
	lans, resp, err := ls.client.LANsApi.DatacentersLansGetExecute(req)
	return Lans{lans}, &Response{*resp}, err
}

func (ls *lansService) Get(datacenterId, lanId string) (*Lan, *Response, error) {
	req := ls.client.LANsApi.DatacentersLansFindById(ls.context, datacenterId, lanId)
	lan, resp, err := ls.client.LANsApi.DatacentersLansFindByIdExecute(req)
	return &Lan{lan}, &Response{*resp}, err
}

func (ls *lansService) Create(datacenterId string, input LanPost) (*LanPost, *Response, error) {
	req := ls.client.LANsApi.DatacentersLansPost(ls.context, datacenterId).Lan(input.Lan)
	lan, resp, err := ls.client.LANsApi.DatacentersLansPostExecute(req)
	return &LanPost{lan}, &Response{*resp}, err
}

func (ls *lansService) Update(datacenterId, lanId string, input LanProperties) (*Lan, *Response, error) {
	req := ls.client.LANsApi.DatacentersLansPatch(ls.context, datacenterId, lanId).Lan(input.LanProperties)
	lan, resp, err := ls.client.LANsApi.DatacentersLansPatchExecute(req)
	return &Lan{lan}, &Response{*resp}, err
}

func (ls *lansService) Delete(datacenterId, lanId string) (*Response, error) {
	req := ls.client.LANsApi.DatacentersLansDelete(ls.context, datacenterId, lanId)
	resp, err := ls.client.LANsApi.DatacentersLansDeleteExecute(req)
	return &Response{*resp}, err
}
