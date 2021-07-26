package v5

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
	Create(datacenterId string, input Volume) (*Volume, *Response, error)
	Update(datacenterId, volumeId string, input VolumeProperties) (*Volume, *Response, error)
	Delete(datacenterId, volumeId string) (*Response, error)
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

func (vs *volumesService) Create(datacenterId string, input Volume) (*Volume, *Response, error) {
	req := vs.client.VolumeApi.DatacentersVolumesPost(vs.context, datacenterId).Volume(input.Volume)
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
