package resources

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"

	"github.com/ionos-cloud/sdk-go-bundle/products/dbaas/psql/v2"
)

type ClusterBackup struct {
	psql.ClusterBackup
}

type BackupResponse struct {
	psql.BackupResponse
}

type ClusterBackupList struct {
	psql.ClusterBackupList
}

// BackupsService is a wrapper around ionoscloud.ClusterBackup
type BackupsService interface {
	List() (ClusterBackupList, *Response, error)
	Get(backupId string) (*BackupResponse, *Response, error)
	ListBackups(clusterId string) (ClusterBackupList, *Response, error)
}

type backupsService struct {
	client  *psql.APIClient
	context context.Context
}

var _ BackupsService = &backupsService{}

func NewBackupsService(client *client.Client, ctx context.Context) BackupsService {
	return &backupsService{
		client:  client.PostgresClient,
		context: ctx,
	}
}

func (svc *backupsService) List() (ClusterBackupList, *Response, error) {
	req := svc.client.BackupsApi.ClustersBackupsGet(svc.context)
	backupList, res, err := svc.client.BackupsApi.ClustersBackupsGetExecute(req)
	return ClusterBackupList{backupList}, &Response{*res}, err
}

func (svc *backupsService) Get(backupId string) (*BackupResponse, *Response, error) {
	req := svc.client.BackupsApi.ClustersBackupsFindById(svc.context, backupId)
	backup, res, err := svc.client.BackupsApi.ClustersBackupsFindByIdExecute(req)
	return &BackupResponse{backup}, &Response{*res}, err
}

func (svc *backupsService) ListBackups(clusterId string) (ClusterBackupList, *Response, error) {
	req := svc.client.BackupsApi.ClusterBackupsGet(svc.context, clusterId)
	clusterBackups, res, err := svc.client.BackupsApi.ClusterBackupsGetExecute(req)
	return ClusterBackupList{clusterBackups}, &Response{*res}, err
}
