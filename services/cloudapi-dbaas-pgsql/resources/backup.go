package resources

import (
	"context"

	sdkgo "github.com/ionos-cloud/sdk-go-autoscaling"
)

type ClusterBackup struct {
	sdkgo.ClusterBackup
}

type ClusterBackupList struct {
	sdkgo.ClusterBackupList
}

// BackupsService is a wrapper around ionoscloud.Cluster
type BackupsService interface {
	List() (ClusterBackupList, *Response, error)
	Get(backupId string) (*ClusterBackup, *Response, error)
	ListBackups(clusterId string) (*ClusterBackupList, *Response, error)
}

type backupsService struct {
	client  *Client
	context context.Context
}

var _ BackupsService = &backupsService{}

func NewBackupsService(client *Client, ctx context.Context) BackupsService {
	return &backupsService{
		client:  client,
		context: ctx,
	}
}

func (cs *backupsService) List() (ClusterBackupList, *Response, error) {
	req := cs.client.BackupsApi.ClustersBackupsGet(cs.context)
	backupList, res, err := cs.client.BackupsApi.ClustersBackupsGetExecute(req)
	return ClusterBackupList{backupList}, &Response{*res}, err
}

func (cs *backupsService) Get(backupId string) (*ClusterBackup, *Response, error) {
	req := cs.client.BackupsApi.ClustersBackupsFindById(cs.context, backupId)
	backup, res, err := cs.client.BackupsApi.ClustersBackupsFindByIdExecute(req)
	return &ClusterBackup{backup}, &Response{*res}, err
}

func (cs *backupsService) ListBackups(clusterId string) (*ClusterBackupList, *Response, error) {
	req := cs.client.BackupsApi.ClusterBackupsGet(cs.context, clusterId)
	clusterBackups, res, err := cs.client.BackupsApi.ClusterBackupsGetExecute(req)
	return &ClusterBackupList{clusterBackups}, &Response{*res}, err
}
