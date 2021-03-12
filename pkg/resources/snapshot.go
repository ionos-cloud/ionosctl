package resources

import (
	"context"

	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
)

type Snapshot struct {
	ionoscloud.Snapshot
}

type Snapshots struct {
	ionoscloud.Snapshots
}

type SnapshotProperties struct {
	ionoscloud.SnapshotProperties
}

// SnapshotsService is a wrapper around ionoscloud.Snapshot
type SnapshotsService interface {
	List() (Snapshots, *Response, error)
	Get(snapshotId string) (*Snapshot, *Response, error)
	Create(datacenterId, volumeId, name, description, licenceType string) (*Snapshot, *Response, error)
	Update(snapshotId string, snapshotProp SnapshotProperties) (*Snapshot, *Response, error)
	Restore(datacenterId, volumeId, snapshotId string) (*Response, error)
	Delete(snapshotId string) (*Response, error)
}

type snapshotsService struct {
	client  *Client
	context context.Context
}

var _ SnapshotsService = &snapshotsService{}

func NewSnapshotService(client *Client, ctx context.Context) SnapshotsService {
	return &snapshotsService{
		client:  client,
		context: ctx,
	}
}

func (s *snapshotsService) List() (Snapshots, *Response, error) {
	req := s.client.SnapshotApi.SnapshotsGet(s.context)
	snapshots, resp, err := s.client.SnapshotApi.SnapshotsGetExecute(req)
	return Snapshots{snapshots}, &Response{*resp}, err
}

func (s *snapshotsService) Get(snapshotId string) (*Snapshot, *Response, error) {
	req := s.client.SnapshotApi.SnapshotsFindById(s.context, snapshotId)
	snapshot, resp, err := s.client.SnapshotApi.SnapshotsFindByIdExecute(req)
	return &Snapshot{snapshot}, &Response{*resp}, err
}

func (s *snapshotsService) Create(datacenterId, volumeId, name, description, licenceType string) (*Snapshot, *Response, error) {
	req := s.client.VolumeApi.DatacentersVolumesCreateSnapshotPost(s.context, datacenterId, volumeId)
	if name != "" {
		req = req.Name(name)
	}
	if description != "" {
		req = req.Description(description)
	}
	if licenceType != "" {
		req = req.LicenceType(licenceType)
	}
	snapshot, resp, err := s.client.VolumeApi.DatacentersVolumesCreateSnapshotPostExecute(req)
	return &Snapshot{snapshot}, &Response{*resp}, err
}

func (s *snapshotsService) Update(snapshotId string, snapshotProp SnapshotProperties) (*Snapshot, *Response, error) {
	req := s.client.SnapshotApi.SnapshotsPatch(s.context, snapshotId).Snapshot(snapshotProp.SnapshotProperties)
	snapshot, resp, err := s.client.SnapshotApi.SnapshotsPatchExecute(req)
	return &Snapshot{snapshot}, &Response{*resp}, err
}

func (s *snapshotsService) Restore(datacenterId, volumeId, snapshotId string) (*Response, error) {
	req := s.client.VolumeApi.DatacentersVolumesRestoreSnapshotPost(s.context, datacenterId, volumeId).SnapshotId(snapshotId)
	_, resp, err := s.client.VolumeApi.DatacentersVolumesRestoreSnapshotPostExecute(req)
	return &Response{*resp}, err
}

func (s *snapshotsService) Delete(snapshotId string) (*Response, error) {
	req := s.client.SnapshotApi.SnapshotsDelete(s.context, snapshotId)
	_, resp, err := s.client.SnapshotApi.SnapshotsDeleteExecute(req)
	return &Response{*resp}, err
}
