package resources

import (
	"context"
	"github.com/ionos-cloud/ionosctl/v6/pkg/config"

	"github.com/fatih/structs"
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

func NewSnapshotService(client *config.Client, ctx context.Context) SnapshotsService {
	return &snapshotsService{
		client:  client.CloudClient,
		context: ctx,
	}
}

func (s *snapshotsService) List(params ListQueryParams) (Snapshots, *Response, error) {
	req := s.client.SnapshotsApi.SnapshotsGet(s.context)
	if !structs.IsZero(params) {
		if params.Filters != nil {
			for k, v := range *params.Filters {
				req = req.Filter(k, v)
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
	snapshots, resp, err := s.client.SnapshotsApi.SnapshotsGetExecute(req)
	return Snapshots{snapshots}, &Response{*resp}, err
}

func (s *snapshotsService) Get(snapshotId string, params QueryParams) (*Snapshot, *Response, error) {
	req := s.client.SnapshotsApi.SnapshotsFindById(s.context, snapshotId)
	if !structs.IsZero(params) {
		if params.Depth != nil {
			req = req.Depth(*params.Depth)
		}
		if params.Pretty != nil {
			// Currently not implemented
			req = req.Pretty(*params.Pretty)
		}
	}
	snapshot, resp, err := s.client.SnapshotsApi.SnapshotsFindByIdExecute(req)
	return &Snapshot{snapshot}, &Response{*resp}, err
}

func (s *snapshotsService) Create(datacenterId, volumeId, name, description, licenceType string, secAuthProtection bool, params QueryParams) (*Snapshot, *Response, error) {
	req := s.client.VolumesApi.DatacentersVolumesCreateSnapshotPost(s.context, datacenterId, volumeId).Name(name).Description(description).LicenceType(licenceType).SecAuthProtection(secAuthProtection)
	if !structs.IsZero(params) {
		if params.Depth != nil {
			req = req.Depth(*params.Depth)
		}
		if params.Pretty != nil {
			// Currently not implemented
			req = req.Pretty(*params.Pretty)
		}
	}
	snapshot, resp, err := s.client.VolumesApi.DatacentersVolumesCreateSnapshotPostExecute(req)
	return &Snapshot{snapshot}, &Response{*resp}, err
}

func (s *snapshotsService) Update(snapshotId string, snapshotProp SnapshotProperties, params QueryParams) (*Snapshot, *Response, error) {
	req := s.client.SnapshotsApi.SnapshotsPatch(s.context, snapshotId).Snapshot(snapshotProp.SnapshotProperties)
	if !structs.IsZero(params) {
		if params.Depth != nil {
			req = req.Depth(*params.Depth)
		}
		if params.Pretty != nil {
			// Currently not implemented
			req = req.Pretty(*params.Pretty)
		}
	}
	snapshot, resp, err := s.client.SnapshotsApi.SnapshotsPatchExecute(req)
	return &Snapshot{snapshot}, &Response{*resp}, err
}

func (s *snapshotsService) Restore(datacenterId, volumeId, snapshotId string, params QueryParams) (*Response, error) {
	req := s.client.VolumesApi.DatacentersVolumesRestoreSnapshotPost(s.context, datacenterId, volumeId).SnapshotId(snapshotId)
	resp, err := s.client.VolumesApi.DatacentersVolumesRestoreSnapshotPostExecute(req)
	return &Response{*resp}, err
}

func (s *snapshotsService) Delete(snapshotId string, params QueryParams) (*Response, error) {
	req := s.client.SnapshotsApi.SnapshotsDelete(s.context, snapshotId)
	resp, err := s.client.SnapshotsApi.SnapshotsDeleteExecute(req)
	return &Response{*resp}, err
}
