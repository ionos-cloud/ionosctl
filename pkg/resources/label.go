package resources

import (
	"context"

	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
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
	req := svc.client.LabelsApi.LabelsFindByUrn(svc.context, labelurn)
	ls, res, err := svc.client.LabelsApi.LabelsFindByUrnExecute(req)
	return &Label{ls}, &Response{*res}, err
}

func (svc *labelResourcesService) List() (Labels, *Response, error) {
	req := svc.client.LabelsApi.LabelsGet(svc.context)
	ls, res, err := svc.client.LabelsApi.LabelsGetExecute(req)
	return Labels{ls}, &Response{*res}, err
}

func (svc *labelResourcesService) DatacenterList(datacenterId string) (LabelResources, *Response, error) {
	req := svc.client.LabelsApi.DatacentersLabelsGet(svc.context, datacenterId)
	ls, res, err := svc.client.LabelsApi.DatacentersLabelsGetExecute(req)
	return LabelResources{ls}, &Response{*res}, err
}

func (svc *labelResourcesService) DatacenterGet(datacenterId, key string) (*LabelResource, *Response, error) {
	req := svc.client.LabelsApi.DatacentersLabelsFindByKey(svc.context, datacenterId, key)
	ls, res, err := svc.client.LabelsApi.DatacentersLabelsFindByKeyExecute(req)
	return &LabelResource{ls}, &Response{*res}, err
}

func (svc *labelResourcesService) DatacenterCreate(datacenterId, key, value string) (*LabelResource, *Response, error) {
	input := ionoscloud.LabelResource{
		Properties: &ionoscloud.LabelResourceProperties{
			Key:   &key,
			Value: &value,
		},
	}
	req := svc.client.LabelsApi.DatacentersLabelsPost(svc.context, datacenterId).Label(input)
	ls, res, err := svc.client.LabelsApi.DatacentersLabelsPostExecute(req)
	return &LabelResource{ls}, &Response{*res}, err
}

func (svc *labelResourcesService) DatacenterDelete(datacenterId, key string) (*Response, error) {
	req := svc.client.LabelsApi.DatacentersLabelsDelete(svc.context, datacenterId, key)
	_, res, err := svc.client.LabelsApi.DatacentersLabelsDeleteExecute(req)
	return &Response{*res}, err
}

func (svc *labelResourcesService) ServerList(datacenterId, serverId string) (LabelResources, *Response, error) {
	req := svc.client.LabelsApi.DatacentersServersLabelsGet(svc.context, datacenterId, serverId)
	ls, res, err := svc.client.LabelsApi.DatacentersServersLabelsGetExecute(req)
	return LabelResources{ls}, &Response{*res}, err
}

func (svc *labelResourcesService) ServerGet(datacenterId, serverId, key string) (*LabelResource, *Response, error) {
	req := svc.client.LabelsApi.DatacentersServersLabelsFindByKey(svc.context, datacenterId, serverId, key)
	ls, res, err := svc.client.LabelsApi.DatacentersServersLabelsFindByKeyExecute(req)
	return &LabelResource{ls}, &Response{*res}, err
}

func (svc *labelResourcesService) ServerCreate(datacenterId, serverId, key, value string) (*LabelResource, *Response, error) {
	input := ionoscloud.LabelResource{
		Properties: &ionoscloud.LabelResourceProperties{
			Key:   &key,
			Value: &value,
		},
	}
	req := svc.client.LabelsApi.DatacentersServersLabelsPost(svc.context, datacenterId, serverId).Label(input)
	ls, res, err := svc.client.LabelsApi.DatacentersServersLabelsPostExecute(req)
	return &LabelResource{ls}, &Response{*res}, err
}

func (svc *labelResourcesService) ServerDelete(datacenterId, serverId, key string) (*Response, error) {
	req := svc.client.LabelsApi.DatacentersServersLabelsDelete(svc.context, datacenterId, serverId, key)
	_, res, err := svc.client.LabelsApi.DatacentersServersLabelsDeleteExecute(req)
	return &Response{*res}, err
}

func (svc *labelResourcesService) VolumeList(datacenterId, volumeId string) (LabelResources, *Response, error) {
	req := svc.client.LabelsApi.DatacentersVolumesLabelsGet(svc.context, datacenterId, volumeId)
	ls, res, err := svc.client.LabelsApi.DatacentersVolumesLabelsGetExecute(req)
	return LabelResources{ls}, &Response{*res}, err
}

func (svc *labelResourcesService) VolumeGet(datacenterId, volumeId, key string) (*LabelResource, *Response, error) {
	req := svc.client.LabelsApi.DatacentersVolumesLabelsFindByKey(svc.context, datacenterId, volumeId, key)
	ls, res, err := svc.client.LabelsApi.DatacentersVolumesLabelsFindByKeyExecute(req)
	return &LabelResource{ls}, &Response{*res}, err
}

func (svc *labelResourcesService) VolumeCreate(datacenterId, volumeId, key, value string) (*LabelResource, *Response, error) {
	input := ionoscloud.LabelResource{
		Properties: &ionoscloud.LabelResourceProperties{
			Key:   &key,
			Value: &value,
		},
	}
	req := svc.client.LabelsApi.DatacentersVolumesLabelsPost(svc.context, datacenterId, volumeId).Label(input)
	ls, res, err := svc.client.LabelsApi.DatacentersVolumesLabelsPostExecute(req)
	return &LabelResource{ls}, &Response{*res}, err
}

func (svc *labelResourcesService) VolumeDelete(datacenterId, volumeId, key string) (*Response, error) {
	req := svc.client.LabelsApi.DatacentersVolumesLabelsDelete(svc.context, datacenterId, volumeId, key)
	_, res, err := svc.client.LabelsApi.DatacentersVolumesLabelsDeleteExecute(req)
	return &Response{*res}, err
}

func (svc *labelResourcesService) IpBlockList(ipblockId string) (LabelResources, *Response, error) {
	req := svc.client.LabelsApi.IpblocksLabelsGet(svc.context, ipblockId)
	ls, res, err := svc.client.LabelsApi.IpblocksLabelsGetExecute(req)
	return LabelResources{ls}, &Response{*res}, err
}

func (svc *labelResourcesService) IpBlockGet(ipblockId, key string) (*LabelResource, *Response, error) {
	req := svc.client.LabelsApi.IpblocksLabelsFindByKey(svc.context, ipblockId, key)
	ls, res, err := svc.client.LabelsApi.IpblocksLabelsFindByKeyExecute(req)
	return &LabelResource{ls}, &Response{*res}, err
}

func (svc *labelResourcesService) IpBlockCreate(ipblockId, key, value string) (*LabelResource, *Response, error) {
	input := ionoscloud.LabelResource{
		Properties: &ionoscloud.LabelResourceProperties{
			Key:   &key,
			Value: &value,
		},
	}
	req := svc.client.LabelsApi.IpblocksLabelsPost(svc.context, ipblockId).Label(input)
	ls, res, err := svc.client.LabelsApi.IpblocksLabelsPostExecute(req)
	return &LabelResource{ls}, &Response{*res}, err
}

func (svc *labelResourcesService) IpBlockDelete(ipblockId, key string) (*Response, error) {
	req := svc.client.LabelsApi.IpblocksLabelsDelete(svc.context, ipblockId, key)
	_, res, err := svc.client.LabelsApi.IpblocksLabelsDeleteExecute(req)
	return &Response{*res}, err
}

func (svc *labelResourcesService) SnapshotList(snapshotId string) (LabelResources, *Response, error) {
	req := svc.client.LabelsApi.SnapshotsLabelsGet(svc.context, snapshotId)
	ls, res, err := svc.client.LabelsApi.SnapshotsLabelsGetExecute(req)
	return LabelResources{ls}, &Response{*res}, err
}

func (svc *labelResourcesService) SnapshotGet(snapshotId, key string) (*LabelResource, *Response, error) {
	req := svc.client.LabelsApi.SnapshotsLabelsFindByKey(svc.context, snapshotId, key)
	ls, res, err := svc.client.LabelsApi.SnapshotsLabelsFindByKeyExecute(req)
	return &LabelResource{ls}, &Response{*res}, err
}

func (svc *labelResourcesService) SnapshotCreate(snapshotId, key, value string) (*LabelResource, *Response, error) {
	input := ionoscloud.LabelResource{
		Properties: &ionoscloud.LabelResourceProperties{
			Key:   &key,
			Value: &value,
		},
	}
	req := svc.client.LabelsApi.SnapshotsLabelsPost(svc.context, snapshotId).Label(input)
	ls, res, err := svc.client.LabelsApi.SnapshotsLabelsPostExecute(req)
	return &LabelResource{ls}, &Response{*res}, err
}

func (svc *labelResourcesService) SnapshotDelete(snapshotId, key string) (*Response, error) {
	req := svc.client.LabelsApi.SnapshotsLabelsDelete(svc.context, snapshotId, key)
	_, res, err := svc.client.LabelsApi.SnapshotsLabelsDeleteExecute(req)
	return &Response{*res}, err
}
