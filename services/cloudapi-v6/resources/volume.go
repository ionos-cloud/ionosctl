package resources

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"

	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
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
	List(datacenterId string, params ListQueryParams) (Volumes, *Response, error)
	Get(datacenterId, volumeId string) (*Volume, *Response, error)
	Create(datacenterId string, input Volume) (*Volume, *Response, error)
	Update(datacenterId, volumeId string, input VolumeProperties) (*Volume, *Response, error)
	Delete(datacenterId, volumeId string) (*Response, error)
}

type volumesService struct {
	client  *ionoscloud.APIClient
	context context.Context
}

var _ VolumesService = &volumesService{}

func NewVolumeService(client *client.Client, ctx context.Context) VolumesService {
	return &volumesService{
		client:  client.CloudClient,
		context: ctx,
	}
}

func (vs *volumesService) List(datacenterId string, params ListQueryParams) (Volumes, *Response, error) {
	req := vs.client.VolumesApi.DatacentersVolumesGet(vs.context, datacenterId)
	volumes, res, err := vs.client.VolumesApi.DatacentersVolumesGetExecute(req)
	return Volumes{volumes}, &Response{*res}, err
}

func (vs *volumesService) Get(datacenterId, volumeId string) (*Volume, *Response, error) {
	req := vs.client.VolumesApi.DatacentersVolumesFindById(vs.context, datacenterId, volumeId)
	volume, res, err := vs.client.VolumesApi.DatacentersVolumesFindByIdExecute(req)
	return &Volume{volume}, &Response{*res}, err
}

func (vs *volumesService) Create(datacenterId string, input Volume) (*Volume, *Response, error) {
	req := vs.client.VolumesApi.DatacentersVolumesPost(vs.context, datacenterId).Volume(input.Volume)
	volume, res, err := vs.client.VolumesApi.DatacentersVolumesPostExecute(req)
	return &Volume{volume}, &Response{*res}, err
}

func (vs *volumesService) Update(datacenterId, volumeId string, input VolumeProperties) (*Volume, *Response, error) {
	req := vs.client.VolumesApi.DatacentersVolumesPatch(vs.context, datacenterId, volumeId).Volume(input.VolumeProperties)
	volume, res, err := vs.client.VolumesApi.DatacentersVolumesPatchExecute(req)
	return &Volume{volume}, &Response{*res}, err
}

func (vs *volumesService) Delete(datacenterId, volumeId string) (*Response, error) {
	req := vs.client.VolumesApi.DatacentersVolumesDelete(vs.context, datacenterId, volumeId)
	res, err := vs.client.VolumesApi.DatacentersVolumesDeleteExecute(req)
	return &Response{*res}, err
}
