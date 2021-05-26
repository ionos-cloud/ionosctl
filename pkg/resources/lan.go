package resources

import (
	"context"

	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
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
	ionoscloud.LanPost
}

type Lans struct {
	ionoscloud.Lans
}

// LansService is a wrapper around ionoscloud.Lan
type LansService interface {
	List(datacenterId string) (Lans, *Response, error)
	Get(datacenterId, lanId string) (*Lan, *Response, error)
	Create(datacenterId string, input LanPost) (*LanPost, *Response, error)
	Update(datacenterId, lanId string, input LanProperties) (*Lan, *Response, error)
	Delete(datacenterId, lanId string) (*Response, error)
}

type lansService struct {
	client  *Client
	context context.Context
}

var _ LansService = &lansService{}

func NewLanService(client *Client, ctx context.Context) LansService {
	return &lansService{
		client:  client,
		context: ctx,
	}
}

func (ls *lansService) List(datacenterId string) (Lans, *Response, error) {
	req := ls.client.LansApi.DatacentersLansGet(ls.context, datacenterId)
	lans, resp, err := ls.client.LansApi.DatacentersLansGetExecute(req)
	return Lans{lans}, &Response{*resp}, err
}

func (ls *lansService) Get(datacenterId, lanId string) (*Lan, *Response, error) {
	req := ls.client.LansApi.DatacentersLansFindById(ls.context, datacenterId, lanId)
	lan, resp, err := ls.client.LansApi.DatacentersLansFindByIdExecute(req)
	return &Lan{lan}, &Response{*resp}, err
}

func (ls *lansService) Create(datacenterId string, input LanPost) (*LanPost, *Response, error) {
	req := ls.client.LansApi.DatacentersLansPost(ls.context, datacenterId).Lan(input.LanPost)
	lan, resp, err := ls.client.LansApi.DatacentersLansPostExecute(req)
	return &LanPost{lan}, &Response{*resp}, err
}

func (ls *lansService) Update(datacenterId, lanId string, input LanProperties) (*Lan, *Response, error) {
	req := ls.client.LansApi.DatacentersLansPatch(ls.context, datacenterId, lanId).Lan(input.LanProperties)
	lan, resp, err := ls.client.LansApi.DatacentersLansPatchExecute(req)
	return &Lan{lan}, &Response{*resp}, err
}

func (ls *lansService) Delete(datacenterId, lanId string) (*Response, error) {
	req := ls.client.LansApi.DatacentersLansDelete(ls.context, datacenterId, lanId)
	_, resp, err := ls.client.LansApi.DatacentersLansDeleteExecute(req)
	return &Response{*resp}, err
}
