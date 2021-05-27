package resources

import (
	"context"

	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
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

type Token struct {
	ionoscloud.Token
}

type RemoteConsoleUrl struct {
	ionoscloud.RemoteConsoleUrl
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
	Suspend(datacenterId, serverId string) (*Response, error)
	Resume(datacenterId, serverId string) (*Response, error)
	GetToken(datacenterId, serverId string) (Token, *Response, error)
	GetRemoteConsoleUrl(datacenterId, serverId string) (RemoteConsoleUrl, *Response, error)
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
	req := ss.client.ServersApi.DatacentersServersGet(ss.context, datacenterId)
	s, res, err := ss.client.ServersApi.DatacentersServersGetExecute(req)
	return Servers{s}, &Response{*res}, err
}

func (ss *serversService) Get(datacenterId, serverId string) (*Server, *Response, error) {
	req := ss.client.ServersApi.DatacentersServersFindById(ss.context, datacenterId, serverId)
	s, res, err := ss.client.ServersApi.DatacentersServersFindByIdExecute(req)
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
	req := ss.client.ServersApi.DatacentersServersPost(ss.context, datacenterId).Server(s)
	server, res, err := ss.client.ServersApi.DatacentersServersPostExecute(req)
	return &Server{server}, &Response{*res}, err
}

func (ss *serversService) Update(datacenterId, serverId string, input ServerProperties) (*Server, *Response, error) {
	req := ss.client.ServersApi.DatacentersServersPatch(ss.context, datacenterId, serverId).Server(input.ServerProperties)
	server, resp, err := ss.client.ServersApi.DatacentersServersPatchExecute(req)
	return &Server{server}, &Response{*resp}, err
}

func (ss *serversService) Delete(datacenterId, serverId string) (*Response, error) {
	req := ss.client.ServersApi.DatacentersServersDelete(ss.context, datacenterId, serverId)
	_, res, err := ss.client.ServersApi.DatacentersServersDeleteExecute(req)
	return &Response{*res}, err
}

func (ss *serversService) Start(datacenterId, serverId string) (*Response, error) {
	req := ss.client.ServersApi.DatacentersServersStartPost(ss.context, datacenterId, serverId)
	_, res, err := ss.client.ServersApi.DatacentersServersStartPostExecute(req)
	return &Response{*res}, err
}

func (ss *serversService) Stop(datacenterId, serverId string) (*Response, error) {
	req := ss.client.ServersApi.DatacentersServersStopPost(ss.context, datacenterId, serverId)
	_, res, err := ss.client.ServersApi.DatacentersServersStopPostExecute(req)
	return &Response{*res}, err
}

func (ss *serversService) Reboot(datacenterId, serverId string) (*Response, error) {
	req := ss.client.ServersApi.DatacentersServersRebootPost(ss.context, datacenterId, serverId)
	_, res, err := ss.client.ServersApi.DatacentersServersRebootPostExecute(req)
	return &Response{*res}, err
}

func (ss *serversService) Suspend(datacenterId, serverId string) (*Response, error) {
	req := ss.client.ServersApi.DatacentersServersSuspendPost(ss.context, datacenterId, serverId)
	_, res, err := ss.client.ServersApi.DatacentersServersSuspendPostExecute(req)
	return &Response{*res}, err
}

func (ss *serversService) Resume(datacenterId, serverId string) (*Response, error) {
	req := ss.client.ServersApi.DatacentersServersResumePost(ss.context, datacenterId, serverId)
	_, res, err := ss.client.ServersApi.DatacentersServersResumePostExecute(req)
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

func (ss *serversService) ListVolumes(datacenterId, serverId string) (AttachedVolumes, *Response, error) {
	req := ss.client.ServersApi.DatacentersServersVolumesGet(ss.context, datacenterId, serverId)
	vols, res, err := ss.client.ServersApi.DatacentersServersVolumesGetExecute(req)
	return AttachedVolumes{vols}, &Response{*res}, err
}

func (ss *serversService) AttachVolume(datacenterId, serverId, volumeId string) (*Volume, *Response, error) {
	req := ss.client.ServersApi.DatacentersServersVolumesPost(ss.context, datacenterId, serverId)
	req = req.Volume(ionoscloud.Volume{Id: &volumeId})
	vol, res, err := ss.client.ServersApi.DatacentersServersVolumesPostExecute(req)
	return &Volume{vol}, &Response{*res}, err
}

func (ss *serversService) GetVolume(datacenterId, serverId, volumeId string) (*Volume, *Response, error) {
	req := ss.client.ServersApi.DatacentersServersVolumesFindById(ss.context, datacenterId, serverId, volumeId)
	vol, res, err := ss.client.ServersApi.DatacentersServersVolumesFindByIdExecute(req)
	return &Volume{vol}, &Response{*res}, err
}

func (ss *serversService) DetachVolume(datacenterId, serverId, volumeId string) (*Response, error) {
	req := ss.client.ServersApi.DatacentersServersVolumesDelete(ss.context, datacenterId, serverId, volumeId)
	_, res, err := ss.client.ServersApi.DatacentersServersVolumesDeleteExecute(req)
	return &Response{*res}, err
}

func (ss *serversService) ListCdroms(datacenterId, serverId string) (Cdroms, *Response, error) {
	req := ss.client.ServersApi.DatacentersServersCdromsGet(ss.context, datacenterId, serverId)
	imgs, res, err := ss.client.ServersApi.DatacentersServersCdromsGetExecute(req)
	return Cdroms{imgs}, &Response{*res}, err
}

func (ss *serversService) AttachCdrom(datacenterId, serverId, cdromId string) (*Image, *Response, error) {
	req := ss.client.ServersApi.DatacentersServersCdromsPost(ss.context, datacenterId, serverId).Cdrom(ionoscloud.Image{Id: &cdromId})
	img, res, err := ss.client.ServersApi.DatacentersServersCdromsPostExecute(req)
	return &Image{img}, &Response{*res}, err
}

func (ss *serversService) GetCdrom(datacenterId, serverId, cdromId string) (*Image, *Response, error) {
	req := ss.client.ServersApi.DatacentersServersCdromsFindById(ss.context, datacenterId, serverId, cdromId)
	img, res, err := ss.client.ServersApi.DatacentersServersCdromsFindByIdExecute(req)
	return &Image{img}, &Response{*res}, err
}

func (ss *serversService) DetachCdrom(datacenterId, serverId, cdromId string) (*Response, error) {
	req := ss.client.ServersApi.DatacentersServersCdromsDelete(ss.context, datacenterId, serverId, cdromId)
	_, res, err := ss.client.ServersApi.DatacentersServersCdromsDeleteExecute(req)
	return &Response{*res}, err
}
