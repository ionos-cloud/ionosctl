package resources

import (
	"context"
	sdkgo "github.com/ionos-cloud/sdk-go-dbaas-mongo"
)

type UsersService interface {
	List(clusterID string) (sdkgo.UsersList, *sdkgo.APIResponse, error)
}

type usersService struct {
	client  *Client
	context context.Context
}

var _ UsersService = &usersService{}

func NewUsersService(client *Client, ctx context.Context) UsersService {
	return &usersService{
		client:  client,
		context: ctx,
	}
}

func (svc *usersService) List(clusterID string) (sdkgo.UsersList, *sdkgo.APIResponse, error) {
	req := svc.client.UsersApi.ClustersUsersGet(svc.context, clusterID)
	ls, res, err := svc.client.UsersApi.ClustersUsersGetExecute(req)
	return ls, res, err
}
