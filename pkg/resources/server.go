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

type Cdroms struct {
	ionoscloud.Cdroms
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
	AttachVolume(datacenterId, serverId, volumeId string) (*Volume, *Response, error)
	DetachVolume(datacenterId, serverId, volumeId string) (*Response, error)
	ListVolumes(datacenterId, serverId string) (AttachedVolumes, *Response, error)
	GetVolume(datacenterId, serverId, volumeId string) (*Volume, *Response, error)
	ListCdroms(datacenterId, serverId string) (Cdroms, *Response, error)
	AttachCdrom(datacenterId, serverId, cdromId string) (*Image, *Response, error)
	GetCdrom(datacenterId, serverId, cdromId string) (*Image, *Response, error)
	DetachCdrom(datacenterId, serverId, cdromId string) (*Response, error)
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

func (ss *serversService) ListVolumes(datacenterId, serverId string) (AttachedVolumes, *Response, error) {
	req := ss.client.ServerApi.DatacentersServersVolumesGet(ss.context, datacenterId, serverId)
	vols, res, err := ss.client.ServerApi.DatacentersServersVolumesGetExecute(req)
	return AttachedVolumes{vols}, &Response{*res}, err
}

func (ss *serversService) AttachVolume(datacenterId, serverId, volumeId string) (*Volume, *Response, error) {
	req := ss.client.ServerApi.DatacentersServersVolumesPost(ss.context, datacenterId, serverId)
	req = req.Volume(ionoscloud.Volume{Id: &volumeId})
	vol, res, err := ss.client.ServerApi.DatacentersServersVolumesPostExecute(req)
	return &Volume{vol}, &Response{*res}, err
}

func (ss *serversService) GetVolume(datacenterId, serverId, volumeId string) (*Volume, *Response, error) {
	req := ss.client.ServerApi.DatacentersServersVolumesFindById(ss.context, datacenterId, serverId, volumeId)
	vol, res, err := ss.client.ServerApi.DatacentersServersVolumesFindByIdExecute(req)
	return &Volume{vol}, &Response{*res}, err
}

func (ss *serversService) DetachVolume(datacenterId, serverId, volumeId string) (*Response, error) {
	req := ss.client.ServerApi.DatacentersServersVolumesDelete(ss.context, datacenterId, serverId, volumeId)
	_, res, err := ss.client.ServerApi.DatacentersServersVolumesDeleteExecute(req)
	return &Response{*res}, err
}

func (ss *serversService) ListCdroms(datacenterId, serverId string) (Cdroms, *Response, error) {
	req := ss.client.ServerApi.DatacentersServersCdromsGet(ss.context, datacenterId, serverId)
	imgs, res, err := ss.client.ServerApi.DatacentersServersCdromsGetExecute(req)
	return Cdroms{imgs}, &Response{*res}, err
}

func (ss *serversService) AttachCdrom(datacenterId, serverId, cdromId string) (*Image, *Response, error) {
	req := ss.client.ServerApi.DatacentersServersCdromsPost(ss.context, datacenterId, serverId).Cdrom(ionoscloud.Image{Id: &cdromId})
	img, res, err := ss.client.ServerApi.DatacentersServersCdromsPostExecute(req)
	return &Image{img}, &Response{*res}, err
}

func (ss *serversService) GetCdrom(datacenterId, serverId, cdromId string) (*Image, *Response, error) {
	req := ss.client.ServerApi.DatacentersServersCdromsFindById(ss.context, datacenterId, serverId, cdromId)
	img, res, err := ss.client.ServerApi.DatacentersServersCdromsFindByIdExecute(req)
	return &Image{img}, &Response{*res}, err
}

func (ss *serversService) DetachCdrom(datacenterId, serverId, cdromId string) (*Response, error) {
	req := ss.client.ServerApi.DatacentersServersCdromsDelete(ss.context, datacenterId, serverId, cdromId)
	_, res, err := ss.client.ServerApi.DatacentersServersCdromsDeleteExecute(req)
	return &Response{*res}, err
}
