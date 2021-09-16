package resources

import (
	"context"

	sdkgo "github.com/ionos-cloud/sdk-go-autoscaling"
)

type PostgresVersionList struct {
	sdkgo.PostgresVersionList
}

// VersionsService is a wrapper around ionoscloud.Cluster
type VersionsService interface {
	List() ([]PostgresVersionList, *Response, error)
	Get(clusterId string) ([]PostgresVersionList, *Response, error)
}

type versionsService struct {
	client  *Client
	context context.Context
}

var _ VersionsService = &versionsService{}

func NewVersionsService(client *Client, ctx context.Context) VersionsService {
	return &versionsService{
		client:  client,
		context: ctx,
	}
}

func (vs *versionsService) List() ([]PostgresVersionList, *Response, error) {
	req := vs.client.ClustersApi.PostgresVersionsGet(vs.context)
	versions, res, err := vs.client.ClustersApi.PostgresVersionsGetExecute(req)
	pgsqlVersions := make([]PostgresVersionList, 0)
	for _, version := range versions {
		pgsqlVersions = append(pgsqlVersions, PostgresVersionList{version})
	}
	return pgsqlVersions, &Response{*res}, err
}

func (vs *versionsService) Get(clusterId string) ([]PostgresVersionList, *Response, error) {
	req := vs.client.ClustersApi.ClusterPostgresVersionsGet(vs.context, clusterId)
	versions, res, err := vs.client.ClustersApi.ClusterPostgresVersionsGetExecute(req)
	pgsqlVersions := make([]PostgresVersionList, 0)
	for _, version := range versions {
		pgsqlVersions = append(pgsqlVersions, PostgresVersionList{version})
	}
	return pgsqlVersions, &Response{*res}, err
}
