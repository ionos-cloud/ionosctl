package v5

import (
	"context"

	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
)

type Label struct {
	ionoscloud.Label
}

type Labels struct {
	ionoscloud.Labels
}

type LabelResource struct {
	ionoscloud.LabelResource
}

type LabelResources struct {
	ionoscloud.LabelResources
}

// LabelResourcesService is a wrapper around ionoscloud.LabelResource
type LabelResourcesService interface {
	GetByUrn(labelurn string) (*Label, *Response, error)
	List() (Labels, *Response, error)
	DatacenterList(datacenterId string) (LabelResources, *Response, error)
	DatacenterGet(datacenterId, key string) (*LabelResource, *Response, error)
	DatacenterCreate(datacenterId, key, value string) (*LabelResource, *Response, error)
	DatacenterDelete(datacenterId, key string) (*Response, error)
	ServerList(datacenterId, serverId string) (LabelResources, *Response, error)
	ServerGet(datacenterId, serverId, key string) (*LabelResource, *Response, error)
	ServerCreate(datacenterId, serverId, key, value string) (*LabelResource, *Response, error)
	ServerDelete(datacenterId, serverId, key string) (*Response, error)
	VolumeList(datacenterId, serverId string) (LabelResources, *Response, error)
	VolumeGet(datacenterId, serverId, key string) (*LabelResource, *Response, error)
	VolumeCreate(datacenterId, serverId, key, value string) (*LabelResource, *Response, error)
	VolumeDelete(datacenterId, serverId, key string) (*Response, error)
	IpBlockList(ipblockId string) (LabelResources, *Response, error)
	IpBlockGet(ipblockId, key string) (*LabelResource, *Response, error)
	IpBlockCreate(ipblockId, key, value string) (*LabelResource, *Response, error)
	IpBlockDelete(ipblockId, key string) (*Response, error)
	SnapshotList(snapshotId string) (LabelResources, *Response, error)
	SnapshotGet(snapshotId, key string) (*LabelResource, *Response, error)
	SnapshotCreate(snapshotId, key, value string) (*LabelResource, *Response, error)
	SnapshotDelete(snapshotId, key string) (*Response, error)
}

type labelResourcesService struct {
	client  *Client
	context context.Context
}

var _ LabelResourcesService = &labelResourcesService{}

func NewLabelResourceService(client *Client, ctx context.Context) LabelResourcesService {
	return &labelResourcesService{
		client:  client,
		context: ctx,
	}
}

func (svc *labelResourcesService) GetByUrn(labelurn string) (*Label, *Response, error) {
	req := svc.client.LabelApi.LabelsFindByUrn(svc.context, labelurn)
	ls, res, err := svc.client.LabelApi.LabelsFindByUrnExecute(req)
	return &Label{ls}, &Response{*res}, err
}

func (svc *labelResourcesService) List() (Labels, *Response, error) {
	req := svc.client.LabelApi.LabelsGet(svc.context)
	ls, res, err := svc.client.LabelApi.LabelsGetExecute(req)
	return Labels{ls}, &Response{*res}, err
}

func (svc *labelResourcesService) DatacenterList(datacenterId string) (LabelResources, *Response, error) {
	req := svc.client.LabelApi.DatacentersLabelsGet(svc.context, datacenterId)
	ls, res, err := svc.client.LabelApi.DatacentersLabelsGetExecute(req)
	return LabelResources{ls}, &Response{*res}, err
}

func (svc *labelResourcesService) DatacenterGet(datacenterId, key string) (*LabelResource, *Response, error) {
	req := svc.client.LabelApi.DatacentersLabelsFindByKey(svc.context, datacenterId, key)
	ls, res, err := svc.client.LabelApi.DatacentersLabelsFindByKeyExecute(req)
	return &LabelResource{ls}, &Response{*res}, err
}

func (svc *labelResourcesService) DatacenterCreate(datacenterId, key, value string) (*LabelResource, *Response, error) {
	input := ionoscloud.LabelResource{
		Properties: &ionoscloud.LabelResourceProperties{
			Key:   &key,
			Value: &value,
		},
	}
	req := svc.client.LabelApi.DatacentersLabelsPost(svc.context, datacenterId).Label(input)
	ls, res, err := svc.client.LabelApi.DatacentersLabelsPostExecute(req)
	return &LabelResource{ls}, &Response{*res}, err
}

func (svc *labelResourcesService) DatacenterDelete(datacenterId, key string) (*Response, error) {
	req := svc.client.LabelApi.DatacentersLabelsDelete(svc.context, datacenterId, key)
	_, res, err := svc.client.LabelApi.DatacentersLabelsDeleteExecute(req)
	return &Response{*res}, err
}

func (svc *labelResourcesService) ServerList(datacenterId, serverId string) (LabelResources, *Response, error) {
	req := svc.client.LabelApi.DatacentersServersLabelsGet(svc.context, datacenterId, serverId)
	ls, res, err := svc.client.LabelApi.DatacentersServersLabelsGetExecute(req)
	return LabelResources{ls}, &Response{*res}, err
}

func (svc *labelResourcesService) ServerGet(datacenterId, serverId, key string) (*LabelResource, *Response, error) {
	req := svc.client.LabelApi.DatacentersServersLabelsFindByKey(svc.context, datacenterId, serverId, key)
	ls, res, err := svc.client.LabelApi.DatacentersServersLabelsFindByKeyExecute(req)
	return &LabelResource{ls}, &Response{*res}, err
}

func (svc *labelResourcesService) ServerCreate(datacenterId, serverId, key, value string) (*LabelResource, *Response, error) {
	input := ionoscloud.LabelResource{
		Properties: &ionoscloud.LabelResourceProperties{
			Key:   &key,
			Value: &value,
		},
	}
	req := svc.client.LabelApi.DatacentersServersLabelsPost(svc.context, datacenterId, serverId).Label(input)
	ls, res, err := svc.client.LabelApi.DatacentersServersLabelsPostExecute(req)
	return &LabelResource{ls}, &Response{*res}, err
}

func (svc *labelResourcesService) ServerDelete(datacenterId, serverId, key string) (*Response, error) {
	req := svc.client.LabelApi.DatacentersServersLabelsDelete(svc.context, datacenterId, serverId, key)
	_, res, err := svc.client.LabelApi.DatacentersServersLabelsDeleteExecute(req)
	return &Response{*res}, err
}

func (svc *labelResourcesService) VolumeList(datacenterId, volumeId string) (LabelResources, *Response, error) {
	req := svc.client.LabelApi.DatacentersVolumesLabelsGet(svc.context, datacenterId, volumeId)
	ls, res, err := svc.client.LabelApi.DatacentersVolumesLabelsGetExecute(req)
	return LabelResources{ls}, &Response{*res}, err
}

func (svc *labelResourcesService) VolumeGet(datacenterId, volumeId, key string) (*LabelResource, *Response, error) {
	req := svc.client.LabelApi.DatacentersVolumesLabelsFindByKey(svc.context, datacenterId, volumeId, key)
	ls, res, err := svc.client.LabelApi.DatacentersVolumesLabelsFindByKeyExecute(req)
	return &LabelResource{ls}, &Response{*res}, err
}

func (svc *labelResourcesService) VolumeCreate(datacenterId, volumeId, key, value string) (*LabelResource, *Response, error) {
	input := ionoscloud.LabelResource{
		Properties: &ionoscloud.LabelResourceProperties{
			Key:   &key,
			Value: &value,
		},
	}
	req := svc.client.LabelApi.DatacentersVolumesLabelsPost(svc.context, datacenterId, volumeId).Label(input)
	ls, res, err := svc.client.LabelApi.DatacentersVolumesLabelsPostExecute(req)
	return &LabelResource{ls}, &Response{*res}, err
}

func (svc *labelResourcesService) VolumeDelete(datacenterId, volumeId, key string) (*Response, error) {
	req := svc.client.LabelApi.DatacentersVolumesLabelsDelete(svc.context, datacenterId, volumeId, key)
	_, res, err := svc.client.LabelApi.DatacentersVolumesLabelsDeleteExecute(req)
	return &Response{*res}, err
}

func (svc *labelResourcesService) IpBlockList(ipblockId string) (LabelResources, *Response, error) {
	req := svc.client.LabelApi.IpblocksLabelsGet(svc.context, ipblockId)
	ls, res, err := svc.client.LabelApi.IpblocksLabelsGetExecute(req)
	return LabelResources{ls}, &Response{*res}, err
}

func (svc *labelResourcesService) IpBlockGet(ipblockId, key string) (*LabelResource, *Response, error) {
	req := svc.client.LabelApi.IpblocksLabelsFindByKey(svc.context, ipblockId, key)
	ls, res, err := svc.client.LabelApi.IpblocksLabelsFindByKeyExecute(req)
	return &LabelResource{ls}, &Response{*res}, err
}

func (svc *labelResourcesService) IpBlockCreate(ipblockId, key, value string) (*LabelResource, *Response, error) {
	input := ionoscloud.LabelResource{
		Properties: &ionoscloud.LabelResourceProperties{
			Key:   &key,
			Value: &value,
		},
	}
	req := svc.client.LabelApi.IpblocksLabelsPost(svc.context, ipblockId).Label(input)
	ls, res, err := svc.client.LabelApi.IpblocksLabelsPostExecute(req)
	return &LabelResource{ls}, &Response{*res}, err
}

func (svc *labelResourcesService) IpBlockDelete(ipblockId, key string) (*Response, error) {
	req := svc.client.LabelApi.IpblocksLabelsDelete(svc.context, ipblockId, key)
	_, res, err := svc.client.LabelApi.IpblocksLabelsDeleteExecute(req)
	return &Response{*res}, err
}

func (svc *labelResourcesService) SnapshotList(snapshotId string) (LabelResources, *Response, error) {
	req := svc.client.LabelApi.SnapshotsLabelsGet(svc.context, snapshotId)
	ls, res, err := svc.client.LabelApi.SnapshotsLabelsGetExecute(req)
	return LabelResources{ls}, &Response{*res}, err
}

func (svc *labelResourcesService) SnapshotGet(snapshotId, key string) (*LabelResource, *Response, error) {
	req := svc.client.LabelApi.SnapshotsLabelsFindByKey(svc.context, snapshotId, key)
	ls, res, err := svc.client.LabelApi.SnapshotsLabelsFindByKeyExecute(req)
	return &LabelResource{ls}, &Response{*res}, err
}

func (svc *labelResourcesService) SnapshotCreate(snapshotId, key, value string) (*LabelResource, *Response, error) {
	input := ionoscloud.LabelResource{
		Properties: &ionoscloud.LabelResourceProperties{
			Key:   &key,
			Value: &value,
		},
	}
	req := svc.client.LabelApi.SnapshotsLabelsPost(svc.context, snapshotId).Label(input)
	ls, res, err := svc.client.LabelApi.SnapshotsLabelsPostExecute(req)
	return &LabelResource{ls}, &Response{*res}, err
}

func (svc *labelResourcesService) SnapshotDelete(snapshotId, key string) (*Response, error) {
	req := svc.client.LabelApi.SnapshotsLabelsDelete(svc.context, snapshotId, key)
	_, res, err := svc.client.LabelApi.SnapshotsLabelsDeleteExecute(req)
	return &Response{*res}, err
}
