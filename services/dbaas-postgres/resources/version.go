package resources

import (
	"context"
	"github.com/ionos-cloud/ionosctl/v6/pkg/config"

	sdkgo "github.com/ionos-cloud/sdk-go-dbaas-postgres"
)

type PostgresVersionList struct {
	sdkgo.PostgresVersionList
}

// VersionsService is a wrapper around ionoscloud.PostgresVersionList
type VersionsService interface {
	List() (PostgresVersionList, *Response, error)
	Get(clusterId string) (PostgresVersionList, *Response, error)
}

type versionsService struct {
	client  *sdkgo.APIClient
	context context.Context
}

var _ VersionsService = &versionsService{}

func NewVersionsService(client *config.Client, ctx context.Context) VersionsService {
	return &versionsService{
		client:  client.DbaasClient,
		context: ctx,
	}
}

func (svc *versionsService) List() (PostgresVersionList, *Response, error) {
	req := svc.client.ClustersApi.PostgresVersionsGet(svc.context)
	versions, res, err := svc.client.ClustersApi.PostgresVersionsGetExecute(req)
	return PostgresVersionList{versions}, &Response{*res}, err
}

func (svc *versionsService) Get(clusterId string) (PostgresVersionList, *Response, error) {
	req := svc.client.ClustersApi.ClusterPostgresVersionsGet(svc.context, clusterId)
	versions, res, err := svc.client.ClustersApi.ClusterPostgresVersionsGetExecute(req)
	return PostgresVersionList{versions}, &Response{*res}, err
}
