package resources

import (
	"context"

	sdkgo "github.com/ionos-cloud/sdk-go-dbaas-mongo"
)

type UsersService interface {
	List(clusterID string) (sdkgo.UsersList, *sdkgo.APIResponse, error)
	Get(clusterID, database, user string) (sdkgo.User, *sdkgo.APIResponse, error)
	Delete(clusterID, database, user string) (sdkgo.User, *sdkgo.APIResponse, error)
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

func (svc *usersService) Get(clusterID, database, username string) (sdkgo.User, *sdkgo.APIResponse, error) {
	req := svc.client.UsersApi.ClustersUsersFindById(svc.context, clusterID, database, username)
	u, res, err := svc.client.UsersApi.ClustersUsersFindByIdExecute(req)
	return u, res, err
}

func (svc *usersService) Delete(clusterID, database, username string) (sdkgo.User, *sdkgo.APIResponse, error) {
	req := svc.client.UsersApi.ClustersUsersDelete(svc.context, clusterID, database, username)
	u, res, err := svc.client.UsersApi.ClustersUsersDeleteExecute(req)
	return u, res, err
}
