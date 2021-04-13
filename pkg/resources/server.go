package resources

import (
	"context"

	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
)

type Server struct {
	ionoscloud.Server
}

type ServerProperties struct {
	ionoscloud.ServerProperties
}

type Servers struct {
	ionoscloud.Servers
}

// ServersService is a wrapper around ionoscloud.Server
type ServersService interface {
	List(datacenterId string) (Servers, *Response, error)
	Get(datacenterId, serverId string) (*Server, *Response, error)
	Create(name, cpufamily, datacenterId, zone string, cores, ram int32) (*Server, *Response, error)
	Update(datacenterId, serverId string, input ServerProperties) (*Server, *Response, error)
	Delete(datacenterId, serverId string) (*Response, error)
	Start(datacenterId, serverId string) (*Response, error)
	Stop(datacenterId, serverId string) (*Response, error)
	Reboot(datacenterId, serverId string) (*Response, error)
	// Volume Actions
	AttachVolume(datacenterId, serverId, volumeId string) (*Volume, *Response, error)
	DetachVolume(datacenterId, serverId, volumeId string) (*Response, error)
	ListVolumes(datacenterId, serverId string) (AttachedVolumes, *Response, error)
	GetVolume(datacenterId, serverId, volumeId string) (*Volume, *Response, error)
}

type serversService struct {
	client  *Client
	context context.Context
}

var _ ServersService = &serversService{}

func NewServerService(client *Client, ctx context.Context) ServersService {
	return &serversService{
		client:  client,
		context: ctx,
	}
}

func (ss *serversService) List(datacenterId string) (Servers, *Response, error) {
	req := ss.client.ServerApi.DatacentersServersGet(ss.context, datacenterId)
	s, res, err := ss.client.ServerApi.DatacentersServersGetExecute(req)
	return Servers{s}, &Response{*res}, err
}

func (ss *serversService) Get(datacenterId, serverId string) (*Server, *Response, error) {
	req := ss.client.ServerApi.DatacentersServersFindById(ss.context, datacenterId, serverId)
	s, res, err := ss.client.ServerApi.DatacentersServersFindByIdExecute(req)
	return &Server{s}, &Response{*res}, err
}

func (ss *serversService) Create(name, cpufamily, datacenterId, zone string, cores, ram int32) (*Server, *Response, error) {
	s := ionoscloud.Server{
		Properties: &ionoscloud.ServerProperties{
			Name:             &name,
			AvailabilityZone: &zone,
			Cores:            &cores,
			Ram:              &ram,
			CpuFamily:        &cpufamily,
		},
	}
	req := ss.client.ServerApi.DatacentersServersPost(ss.context, datacenterId).Server(s)
	server, res, err := ss.client.ServerApi.DatacentersServersPostExecute(req)
	return &Server{server}, &Response{*res}, err
}

func (ss *serversService) Update(datacenterId, serverId string, input ServerProperties) (*Server, *Response, error) {
	req := ss.client.ServerApi.DatacentersServersPatch(ss.context, datacenterId, serverId).Server(input.ServerProperties)
	server, resp, err := ss.client.ServerApi.DatacentersServersPatchExecute(req)
	return &Server{server}, &Response{*resp}, err
}

func (ss *serversService) Delete(datacenterId, serverId string) (*Response, error) {
	req := ss.client.ServerApi.DatacentersServersDelete(ss.context, datacenterId, serverId)
	_, res, err := ss.client.ServerApi.DatacentersServersDeleteExecute(req)
	return &Response{*res}, err
}

func (ss *serversService) Start(datacenterId, serverId string) (*Response, error) {
	req := ss.client.ServerApi.DatacentersServersStartPost(ss.context, datacenterId, serverId)
	_, res, err := ss.client.ServerApi.DatacentersServersStartPostExecute(req)
	return &Response{*res}, err
}

func (ss *serversService) Stop(datacenterId, serverId string) (*Response, error) {
	req := ss.client.ServerApi.DatacentersServersStopPost(ss.context, datacenterId, serverId)
	_, res, err := ss.client.ServerApi.DatacentersServersStopPostExecute(req)
	return &Response{*res}, err
}

func (ss *serversService) Reboot(datacenterId, serverId string) (*Response, error) {
	req := ss.client.ServerApi.DatacentersServersRebootPost(ss.context, datacenterId, serverId)
	_, res, err := ss.client.ServerApi.DatacentersServersRebootPostExecute(req)
	return &Response{*res}, err
}

func (vs *serversService) ListVolumes(datacenterId, serverId string) (AttachedVolumes, *Response, error) {
	req := vs.client.ServerApi.DatacentersServersVolumesGet(vs.context, datacenterId, serverId)
	vols, res, err := vs.client.ServerApi.DatacentersServersVolumesGetExecute(req)
	return AttachedVolumes{vols}, &Response{*res}, err
}

func (vs *serversService) AttachVolume(datacenterId, serverId, volumeId string) (*Volume, *Response, error) {
	req := vs.client.ServerApi.DatacentersServersVolumesPost(vs.context, datacenterId, serverId)
	req = req.Volume(ionoscloud.Volume{Id: &volumeId})
	vol, res, err := vs.client.ServerApi.DatacentersServersVolumesPostExecute(req)
	return &Volume{vol}, &Response{*res}, err
}

func (vs *serversService) GetVolume(datacenterId, serverId, volumeId string) (*Volume, *Response, error) {
	req := vs.client.ServerApi.DatacentersServersVolumesFindById(vs.context, datacenterId, serverId, volumeId)
	vol, res, err := vs.client.ServerApi.DatacentersServersVolumesFindByIdExecute(req)
	return &Volume{vol}, &Response{*res}, err
}

func (vs *serversService) DetachVolume(datacenterId, serverId, volumeId string) (*Response, error) {
	req := vs.client.ServerApi.DatacentersServersVolumesDelete(vs.context, datacenterId, serverId, volumeId)
	_, res, err := vs.client.ServerApi.DatacentersServersVolumesDeleteExecute(req)
	return &Response{*res}, err
}
