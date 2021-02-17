package resources

import (
	"context"

	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
)

type Volume struct {
	ionoscloud.Volume
}

type VolumeProperties struct {
	ionoscloud.VolumeProperties
}

type Volumes struct {
	ionoscloud.Volumes
}

type AttachedVolumes struct {
	ionoscloud.AttachedVolumes
}

type VolumesService interface {
	List(datacenterId string) (Volumes, *Response, error)
	Get(datacenterId, volumeId string) (*Volume, *Response, error)
	Create(datacenterId, name, bus, volumetype, licencetype, zone string, size float32) (*Volume, *Response, error)
	Update(datacenterId, volumeId string, input VolumeProperties) (*Volume, *Response, error)
	Delete(datacenterId, volumeId string) (*Response, error)
	// Volume Actions
	Attach(datacenterId, serverId, volumeId string) (*Volume, *Response, error)
	Detach(datacenterId, serverId, volumeId string) (*Response, error)
	ListAttached(datacenterId, serverId string) (AttachedVolumes, *Response, error)
	GetAttached(datacenterId, serverId, volumeId string) (*Volume, *Response, error)
}

type volumesService struct {
	client  *Client
	context context.Context
}

var _ VolumesService = &volumesService{}

func NewVolumeService(client *Client, ctx context.Context) VolumesService {
	return &volumesService{
		client:  client,
		context: ctx,
	}
}

func (vs *volumesService) List(datacenterId string) (Volumes, *Response, error) {
	req := vs.client.VolumeApi.DatacentersVolumesGet(vs.context, datacenterId)
	volumes, res, err := vs.client.VolumeApi.DatacentersVolumesGetExecute(req)
	return Volumes{volumes}, &Response{*res}, err
}

func (vs *volumesService) Get(datacenterId, volumeId string) (*Volume, *Response, error) {
	req := vs.client.VolumeApi.DatacentersVolumesFindById(vs.context, datacenterId, volumeId)
	volume, res, err := vs.client.VolumeApi.DatacentersVolumesFindByIdExecute(req)
	return &Volume{volume}, &Response{*res}, err
}

func (vs *volumesService) Create(datacenterId, name, bus, volumetype, licencetype, zone string, size float32) (*Volume, *Response, error) {
	v := ionoscloud.Volume{
		Metadata: nil,
		Properties: &ionoscloud.VolumeProperties{
			Name:             &name,
			Type:             &volumetype,
			Size:             &size,
			AvailabilityZone: &zone,
			Bus:              &bus,
			LicenceType:      &licencetype,
		},
	}
	req := vs.client.VolumeApi.DatacentersVolumesPost(vs.context, datacenterId).Volume(v)
	volume, res, err := vs.client.VolumeApi.DatacentersVolumesPostExecute(req)
	return &Volume{volume}, &Response{*res}, err
}

func (vs *volumesService) Update(datacenterId, volumeId string, input VolumeProperties) (*Volume, *Response, error) {
	req := vs.client.VolumeApi.DatacentersVolumesPatch(vs.context, datacenterId, volumeId).Volume(input.VolumeProperties)
	volume, res, err := vs.client.VolumeApi.DatacentersVolumesPatchExecute(req)
	return &Volume{volume}, &Response{*res}, err
}

func (vs *volumesService) Delete(datacenterId, volumeId string) (*Response, error) {
	req := vs.client.VolumeApi.DatacentersVolumesDelete(vs.context, datacenterId, volumeId)
	_, res, err := vs.client.VolumeApi.DatacentersVolumesDeleteExecute(req)
	return &Response{*res}, err
}

func (vs *volumesService) ListAttached(datacenterId, serverId string) (AttachedVolumes, *Response, error) {
	req := vs.client.ServerApi.DatacentersServersVolumesGet(vs.context, datacenterId, serverId)
	vols, res, err := vs.client.ServerApi.DatacentersServersVolumesGetExecute(req)
	return AttachedVolumes{vols}, &Response{*res}, err
}

func (vs *volumesService) Attach(datacenterId, serverId, volumeId string) (*Volume, *Response, error) {
	req := vs.client.ServerApi.DatacentersServersVolumesPost(vs.context, datacenterId, serverId)
	req = req.Volume(ionoscloud.Volume{Id: &volumeId})
	vol, res, err := vs.client.ServerApi.DatacentersServersVolumesPostExecute(req)
	return &Volume{vol}, &Response{*res}, err
}

func (vs *volumesService) GetAttached(datacenterId, serverId, volumeId string) (*Volume, *Response, error) {
	req := vs.client.ServerApi.DatacentersServersVolumesFindById(vs.context, datacenterId, serverId, volumeId)
	vol, res, err := vs.client.ServerApi.DatacentersServersVolumesFindByIdExecute(req)
	return &Volume{vol}, &Response{*res}, err
}

func (vs *volumesService) Detach(datacenterId, serverId, volumeId string) (*Response, error) {
	req := vs.client.ServerApi.DatacentersServersVolumesDelete(vs.context, datacenterId, serverId, volumeId)
	_, res, err := vs.client.ServerApi.DatacentersServersVolumesDeleteExecute(req)
	return &Response{*res}, err
}
