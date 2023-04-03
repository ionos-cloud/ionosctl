package resources

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"

	"github.com/fatih/structs"
	ionoscloud "github.com/ionos-cloud/sdk-go-bundle/products/compute"
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
	Get(datacenterId, volumeId string, params QueryParams) (*Volume, *Response, error)
	Create(datacenterId string, input Volume, params QueryParams) (*Volume, *Response, error)
	Update(datacenterId, volumeId string, input VolumeProperties, params QueryParams) (*Volume, *Response, error)
	Delete(datacenterId, volumeId string, params QueryParams) (*Response, error)
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
	volumes, res, err := vs.client.VolumesApi.DatacentersVolumesGetExecute(req)
	return Volumes{volumes}, &Response{*res}, err
}

func (vs *volumesService) Get(datacenterId, volumeId string, params QueryParams) (*Volume, *Response, error) {
	req := vs.client.VolumesApi.DatacentersVolumesFindById(vs.context, datacenterId, volumeId)
	if !structs.IsZero(params) {
		if params.Depth != nil {
			req = req.Depth(*params.Depth)
		}
		if params.Pretty != nil {
			// Currently not implemented
			req = req.Pretty(*params.Pretty)
		}
	}
	volume, res, err := vs.client.VolumesApi.DatacentersVolumesFindByIdExecute(req)
	return &Volume{volume}, &Response{*res}, err
}

func (vs *volumesService) Create(datacenterId string, input Volume, params QueryParams) (*Volume, *Response, error) {
	req := vs.client.VolumesApi.DatacentersVolumesPost(vs.context, datacenterId).Volume(input.Volume)
	if !structs.IsZero(params) {
		if params.Depth != nil {
			req = req.Depth(*params.Depth)
		}
		if params.Pretty != nil {
			// Currently not implemented
			req = req.Pretty(*params.Pretty)
		}
	}
	volume, res, err := vs.client.VolumesApi.DatacentersVolumesPostExecute(req)
	return &Volume{volume}, &Response{*res}, err
}

func (vs *volumesService) Update(datacenterId, volumeId string, input VolumeProperties, params QueryParams) (*Volume, *Response, error) {
	req := vs.client.VolumesApi.DatacentersVolumesPatch(vs.context, datacenterId, volumeId).Volume(input.VolumeProperties)
	if !structs.IsZero(params) {
		if params.Depth != nil {
			req = req.Depth(*params.Depth)
		}
		if params.Pretty != nil {
			// Currently not implemented
			req = req.Pretty(*params.Pretty)
		}
	}
	volume, res, err := vs.client.VolumesApi.DatacentersVolumesPatchExecute(req)
	return &Volume{volume}, &Response{*res}, err
}

func (vs *volumesService) Delete(datacenterId, volumeId string, params QueryParams) (*Response, error) {
	req := vs.client.VolumesApi.DatacentersVolumesDelete(vs.context, datacenterId, volumeId)
	if !structs.IsZero(params) {
		if params.Depth != nil {
			req = req.Depth(*params.Depth)
		}
		if params.Pretty != nil {
			// Currently not implemented
			req = req.Pretty(*params.Pretty)
		}
	}
	res, err := vs.client.VolumesApi.DatacentersVolumesDeleteExecute(req)
	return &Response{*res}, err
}
