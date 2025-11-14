package resources

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"

	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
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
	List(params ListQueryParams) (Snapshots, *Response, error)
	Get(snapshotId string, params QueryParams) (*Snapshot, *Response, error)
	Create(datacenterId, volumeId, name, description, licenceType string, secAuthProtection bool, params QueryParams) (*Snapshot, *Response, error)
	Update(snapshotId string, snapshotProp SnapshotProperties, params QueryParams) (*Snapshot, *Response, error)
	Restore(datacenterId, volumeId, snapshotId string, params QueryParams) (*Response, error)
	Delete(snapshotId string, params QueryParams) (*Response, error)
}

type snapshotsService struct {
	client  *ionoscloud.APIClient
	context context.Context
}

var _ SnapshotsService = &snapshotsService{}

func NewSnapshotService(client *client.Client, ctx context.Context) SnapshotsService {
	return &snapshotsService{
		client:  client.CloudClient,
		context: ctx,
	}
}

func (s *snapshotsService) List(params ListQueryParams) (Snapshots, *Response, error) {
	req := s.client.SnapshotsApi.SnapshotsGet(s.context)
	snapshots, resp, err := s.client.SnapshotsApi.SnapshotsGetExecute(req)
	return Snapshots{snapshots}, &Response{*resp}, err
}

func (s *snapshotsService) Get(snapshotId string, params QueryParams) (*Snapshot, *Response, error) {
	req := s.client.SnapshotsApi.SnapshotsFindById(s.context, snapshotId)
	snapshot, resp, err := s.client.SnapshotsApi.SnapshotsFindByIdExecute(req)
	return &Snapshot{snapshot}, &Response{*resp}, err
}

func (s *snapshotsService) Create(datacenterId, volumeId, name, description, licenceType string, secAuthProtection bool, params QueryParams) (*Snapshot, *Response, error) {
	req := s.client.VolumesApi.DatacentersVolumesCreateSnapshotPost(s.context, datacenterId, volumeId).Snapshot(
		ionoscloud.CreateSnapshot{
			Properties: &ionoscloud.CreateSnapshotProperties{
				Name:              &name,
				Description:       &description,
				LicenceType:       &licenceType,
				SecAuthProtection: &secAuthProtection,
			},
		},
	)
	snapshot, resp, err := s.client.VolumesApi.DatacentersVolumesCreateSnapshotPostExecute(req)
	return &Snapshot{snapshot}, &Response{*resp}, err
}

func (s *snapshotsService) Update(snapshotId string, snapshotProp SnapshotProperties, params QueryParams) (*Snapshot, *Response, error) {
	req := s.client.SnapshotsApi.SnapshotsPatch(s.context, snapshotId).Snapshot(snapshotProp.SnapshotProperties)
	snapshot, resp, err := s.client.SnapshotsApi.SnapshotsPatchExecute(req)
	return &Snapshot{snapshot}, &Response{*resp}, err
}

func (s *snapshotsService) Restore(datacenterId, volumeId, snapshotId string, params QueryParams) (*Response, error) {
	req := s.client.VolumesApi.DatacentersVolumesRestoreSnapshotPost(s.context, datacenterId, volumeId).RestoreSnapshot(
		ionoscloud.RestoreSnapshot{
			Properties: &ionoscloud.RestoreSnapshotProperties{
				SnapshotId: &snapshotId,
			},
		})
	resp, err := s.client.VolumesApi.DatacentersVolumesRestoreSnapshotPostExecute(req)
	return &Response{*resp}, err
}

func (s *snapshotsService) Delete(snapshotId string, params QueryParams) (*Response, error) {
	req := s.client.SnapshotsApi.SnapshotsDelete(s.context, snapshotId)
	resp, err := s.client.SnapshotsApi.SnapshotsDeleteExecute(req)
	return &Response{*resp}, err
}
