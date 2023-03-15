package resources

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/pkg/config"

	sdkgo "github.com/ionos-cloud/sdk-go-dbaas-postgres"
)

type CreateRestoreRequest struct {
	sdkgo.CreateRestoreRequest
}

// RestoresService is a wrapper around ionoscloud.CreateRestoreRequest
type RestoresService interface {
	Restore(clusterId string, input CreateRestoreRequest) (*Response, error)
}

type restoresService struct {
	client  *sdkgo.APIClient
	context context.Context
}

var _ RestoresService = &restoresService{}

func NewRestoresService(client *config.Client, ctx context.Context) RestoresService {
	return &restoresService{
		client:  client.PostgresClient,
		context: ctx,
	}
}

func (svc *restoresService) Restore(clusterId string, input CreateRestoreRequest) (*Response, error) {
	req := svc.client.RestoresApi.ClusterRestorePost(svc.context, clusterId).CreateRestoreRequest(input.CreateRestoreRequest)
	res, err := svc.client.RestoresApi.ClusterRestorePostExecute(req)
	return &Response{*res}, err
}
