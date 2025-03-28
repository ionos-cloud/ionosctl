package resources

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/sdk-go-bundle/products/compute/v2"

	"github.com/fatih/structs"
)

type Server struct {
	compute.Server
}

type ServerProperties struct {
	compute.ServerProperties
}

type Servers struct {
	compute.Servers
}

type Cdroms struct {
	compute.Cdroms
}

type Token struct {
	compute.Token
}

type RemoteConsoleUrl struct {
	compute.RemoteConsoleUrl
}

// ServersService is a wrapper around compute.Server
type ServersService interface {
	List(datacenterId string, params ListQueryParams) (Servers, *Response, error)
	Get(datacenterId, serverId string, params QueryParams) (*Server, *Response, error)
	Create(datacenterId string, input Server, params QueryParams) (*Server, *Response, error)
	Update(datacenterId, serverId string, input ServerProperties, params QueryParams) (*Server, *Response, error)
	Delete(datacenterId, serverId string, params QueryParams) (*Response, error)
	Start(datacenterId, serverId string, params QueryParams) (*Response, error)
	Stop(datacenterId, serverId string, params QueryParams) (*Response, error)
	Reboot(datacenterId, serverId string, params QueryParams) (*Response, error)
	Suspend(datacenterId, serverId string, params QueryParams) (*Response, error)
	Resume(datacenterId, serverId string, params QueryParams) (*Response, error)
	GetToken(datacenterId, serverId string) (Token, *Response, error)
	GetRemoteConsoleUrl(datacenterId, serverId string) (RemoteConsoleUrl, *Response, error)
	AttachVolume(datacenterId, serverId, volumeId string, params QueryParams) (*Volume, *Response, error)
	DetachVolume(datacenterId, serverId, volumeId string, params QueryParams) (*Response, error)
	ListVolumes(datacenterId, serverId string, params ListQueryParams) (AttachedVolumes, *Response, error)
	GetVolume(datacenterId, serverId, volumeId string, params QueryParams) (*Volume, *Response, error)
	ListCdroms(datacenterId, serverId string, params ListQueryParams) (Cdroms, *Response, error)
	AttachCdrom(datacenterId, serverId, cdromId string, params QueryParams) (*Image, *Response, error)
	GetCdrom(datacenterId, serverId, cdromId string, params QueryParams) (*Image, *Response, error)
	DetachCdrom(datacenterId, serverId, cdromId string, params QueryParams) (*Response, error)
}

type serversService struct {
	client  *compute.APIClient
	context context.Context
}

var _ ServersService = &serversService{}

func NewServerService(client *client.Client, ctx context.Context) ServersService {
	return &serversService{
		client:  client.CloudClient,
		context: ctx,
	}
}

func (ss *serversService) List(datacenterId string, params ListQueryParams) (Servers, *Response, error) {
	req := ss.client.ServersApi.DatacentersServersGet(ss.context, datacenterId)
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
	s, res, err := ss.client.ServersApi.DatacentersServersGetExecute(req)
	return Servers{s}, &Response{*res}, err
}

func (ss *serversService) Get(datacenterId, serverId string, params QueryParams) (*Server, *Response, error) {
	req := ss.client.ServersApi.DatacentersServersFindById(ss.context, datacenterId, serverId)
	if !structs.IsZero(params) {
		if params.Depth != nil {
			req = req.Depth(*params.Depth)
		}
		if params.Pretty != nil {
			// Currently not implemented
			req = req.Pretty(*params.Pretty)
		}
	}
	s, res, err := ss.client.ServersApi.DatacentersServersFindByIdExecute(req)
	return &Server{s}, &Response{*res}, err
}

func (ss *serversService) Create(datacenterId string, input Server, params QueryParams) (*Server, *Response, error) {
	req := ss.client.ServersApi.DatacentersServersPost(ss.context, datacenterId).Server(input.Server)
	if !structs.IsZero(params) {
		if params.Depth != nil {
			req = req.Depth(*params.Depth)
		}
		if params.Pretty != nil {
			// Currently not implemented
			req = req.Pretty(*params.Pretty)
		}
	}
	server, res, err := ss.client.ServersApi.DatacentersServersPostExecute(req)
	return &Server{server}, &Response{*res}, err
}

func (ss *serversService) Update(datacenterId, serverId string, input ServerProperties, params QueryParams) (*Server, *Response, error) {
	req := ss.client.ServersApi.DatacentersServersPatch(ss.context, datacenterId, serverId).Server(input.ServerProperties)
	if !structs.IsZero(params) {
		if params.Depth != nil {
			req = req.Depth(*params.Depth)
		}
		if params.Pretty != nil {
			// Currently not implemented
			req = req.Pretty(*params.Pretty)
		}
	}
	server, resp, err := ss.client.ServersApi.DatacentersServersPatchExecute(req)
	return &Server{server}, &Response{*resp}, err
}

func (ss *serversService) Delete(datacenterId, serverId string, params QueryParams) (*Response, error) {
	req := ss.client.ServersApi.DatacentersServersDelete(ss.context, datacenterId, serverId)
	if !structs.IsZero(params) {
		if params.Depth != nil {
			req = req.Depth(*params.Depth)
		}
		if params.Pretty != nil {
			// Currently not implemented
			req = req.Pretty(*params.Pretty)
		}
	}
	res, err := ss.client.ServersApi.DatacentersServersDeleteExecute(req)
	return &Response{*res}, err
}

func (ss *serversService) Start(datacenterId, serverId string, params QueryParams) (*Response, error) {
	req := ss.client.ServersApi.DatacentersServersStartPost(ss.context, datacenterId, serverId)
	if !structs.IsZero(params) {
		if params.Depth != nil {
			req = req.Depth(*params.Depth)
		}
		if params.Pretty != nil {
			// Currently not implemented
			req = req.Pretty(*params.Pretty)
		}
	}
	res, err := ss.client.ServersApi.DatacentersServersStartPostExecute(req)
	return &Response{*res}, err
}

func (ss *serversService) Stop(datacenterId, serverId string, params QueryParams) (*Response, error) {
	req := ss.client.ServersApi.DatacentersServersStopPost(ss.context, datacenterId, serverId)
	if !structs.IsZero(params) {
		if params.Depth != nil {
			req = req.Depth(*params.Depth)
		}
		if params.Pretty != nil {
			// Currently not implemented
			req = req.Pretty(*params.Pretty)
		}
	}
	res, err := ss.client.ServersApi.DatacentersServersStopPostExecute(req)
	return &Response{*res}, err
}

func (ss *serversService) Reboot(datacenterId, serverId string, params QueryParams) (*Response, error) {
	req := ss.client.ServersApi.DatacentersServersRebootPost(ss.context, datacenterId, serverId)
	if !structs.IsZero(params) {
		if params.Depth != nil {
			req = req.Depth(*params.Depth)
		}
		if params.Pretty != nil {
			// Currently not implemented
			req = req.Pretty(*params.Pretty)
		}
	}
	res, err := ss.client.ServersApi.DatacentersServersRebootPostExecute(req)
	return &Response{*res}, err
}

func (ss *serversService) Suspend(datacenterId, serverId string, params QueryParams) (*Response, error) {
	req := ss.client.ServersApi.DatacentersServersSuspendPost(ss.context, datacenterId, serverId)
	if !structs.IsZero(params) {
		if params.Depth != nil {
			req = req.Depth(*params.Depth)
		}
		if params.Pretty != nil {
			// Currently not implemented
			req = req.Pretty(*params.Pretty)
		}
	}
	res, err := ss.client.ServersApi.DatacentersServersSuspendPostExecute(req)
	return &Response{*res}, err
}

func (ss *serversService) Resume(datacenterId, serverId string, params QueryParams) (*Response, error) {
	req := ss.client.ServersApi.DatacentersServersResumePost(ss.context, datacenterId, serverId)
	if !structs.IsZero(params) {
		if params.Depth != nil {
			req = req.Depth(*params.Depth)
		}
		if params.Pretty != nil {
			// Currently not implemented
			req = req.Pretty(*params.Pretty)
		}
	}
	res, err := ss.client.ServersApi.DatacentersServersResumePostExecute(req)
	return &Response{*res}, err
}

func (ss *serversService) GetToken(datacenterId, serverId string) (Token, *Response, error) {
	req := ss.client.ServersApi.DatacentersServersTokenGet(ss.context, datacenterId, serverId)
	token, res, err := ss.client.ServersApi.DatacentersServersTokenGetExecute(req)
	return Token{token}, &Response{*res}, err
}

func (ss *serversService) GetRemoteConsoleUrl(datacenterId, serverId string) (RemoteConsoleUrl, *Response, error) {
	req := ss.client.ServersApi.DatacentersServersRemoteConsoleGet(ss.context, datacenterId, serverId)
	url, res, err := ss.client.ServersApi.DatacentersServersRemoteConsoleGetExecute(req)
	return RemoteConsoleUrl{url}, &Response{*res}, err
}

func (ss *serversService) ListVolumes(datacenterId, serverId string, params ListQueryParams) (AttachedVolumes, *Response, error) {
	req := ss.client.ServersApi.DatacentersServersVolumesGet(ss.context, datacenterId, serverId)
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
	vols, res, err := ss.client.ServersApi.DatacentersServersVolumesGetExecute(req)
	return AttachedVolumes{vols}, &Response{*res}, err
}

func (ss *serversService) AttachVolume(datacenterId, serverId, volumeId string, params QueryParams) (*Volume, *Response, error) {
	req := ss.client.ServersApi.DatacentersServersVolumesPost(ss.context, datacenterId, serverId)
	if !structs.IsZero(params) {
		if params.Depth != nil {
			req = req.Depth(*params.Depth)
		}
		if params.Pretty != nil {
			// Currently not implemented
			req = req.Pretty(*params.Pretty)
		}
	}
	req = req.Volume(compute.Volume{Id: &volumeId})
	vol, res, err := ss.client.ServersApi.DatacentersServersVolumesPostExecute(req)
	return &Volume{vol}, &Response{*res}, err
}

func (ss *serversService) GetVolume(datacenterId, serverId, volumeId string, params QueryParams) (*Volume, *Response, error) {
	req := ss.client.ServersApi.DatacentersServersVolumesFindById(ss.context, datacenterId, serverId, volumeId)
	if !structs.IsZero(params) {
		if params.Depth != nil {
			req = req.Depth(*params.Depth)
		}
		if params.Pretty != nil {
			// Currently not implemented
			req = req.Pretty(*params.Pretty)
		}
	}
	vol, res, err := ss.client.ServersApi.DatacentersServersVolumesFindByIdExecute(req)
	return &Volume{vol}, &Response{*res}, err
}

func (ss *serversService) DetachVolume(datacenterId, serverId, volumeId string, params QueryParams) (*Response, error) {
	req := ss.client.ServersApi.DatacentersServersVolumesDelete(ss.context, datacenterId, serverId, volumeId)
	if !structs.IsZero(params) {
		if params.Depth != nil {
			req = req.Depth(*params.Depth)
		}
		if params.Pretty != nil {
			// Currently not implemented
			req = req.Pretty(*params.Pretty)
		}
	}
	res, err := ss.client.ServersApi.DatacentersServersVolumesDeleteExecute(req)
	return &Response{*res}, err
}

func (ss *serversService) ListCdroms(datacenterId, serverId string, params ListQueryParams) (Cdroms, *Response, error) {
	req := ss.client.ServersApi.DatacentersServersCdromsGet(ss.context, datacenterId, serverId)
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
	imgs, res, err := ss.client.ServersApi.DatacentersServersCdromsGetExecute(req)
	return Cdroms{imgs}, &Response{*res}, err
}

func (ss *serversService) AttachCdrom(datacenterId, serverId, cdromId string, params QueryParams) (*Image, *Response, error) {
	req := ss.client.ServersApi.DatacentersServersCdromsPost(ss.context, datacenterId, serverId).Cdrom(compute.Image{Id: &cdromId})
	if !structs.IsZero(params) {
		if params.Depth != nil {
			req = req.Depth(*params.Depth)
		}
		if params.Pretty != nil {
			// Currently not implemented
			req = req.Pretty(*params.Pretty)
		}
	}
	img, res, err := ss.client.ServersApi.DatacentersServersCdromsPostExecute(req)
	return &Image{img}, &Response{*res}, err
}

func (ss *serversService) GetCdrom(datacenterId, serverId, cdromId string, params QueryParams) (*Image, *Response, error) {
	req := ss.client.ServersApi.DatacentersServersCdromsFindById(ss.context, datacenterId, serverId, cdromId)
	if !structs.IsZero(params) {
		if params.Depth != nil {
			req = req.Depth(*params.Depth)
		}
		if params.Pretty != nil {
			// Currently not implemented
			req = req.Pretty(*params.Pretty)
		}
	}
	img, res, err := ss.client.ServersApi.DatacentersServersCdromsFindByIdExecute(req)
	return &Image{img}, &Response{*res}, err
}

func (ss *serversService) DetachCdrom(datacenterId, serverId, cdromId string, params QueryParams) (*Response, error) {
	req := ss.client.ServersApi.DatacentersServersCdromsDelete(ss.context, datacenterId, serverId, cdromId)
	if !structs.IsZero(params) {
		if params.Depth != nil {
			req = req.Depth(*params.Depth)
		}
		if params.Pretty != nil {
			// Currently not implemented
			req = req.Pretty(*params.Pretty)
		}
	}
	res, err := ss.client.ServersApi.DatacentersServersCdromsDeleteExecute(req)
	return &Response{*res}, err
}
