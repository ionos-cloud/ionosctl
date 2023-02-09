package resources

import (
	"context"
	"github.com/ionos-cloud/ionosctl/pkg/config"

	sdkgo "github.com/ionos-cloud/sdk-go-dbaas-mongo"
)

type UsersService interface {
	List(clusterID string) (sdkgo.UsersList, *sdkgo.APIResponse, error)
	Create(clusterID string, user sdkgo.User) (sdkgo.User, *sdkgo.APIResponse, error)
	ListAll() ([]sdkgo.User, error)
	Get(clusterID, database, user string) (sdkgo.User, *sdkgo.APIResponse, error)
	Delete(clusterID, database, user string) (sdkgo.User, *sdkgo.APIResponse, error)
}

type usersService struct {
	client  *sdkgo.APIClient
	context context.Context
}

var _ UsersService = &usersService{}

func NewUsersService(client *config.Client, ctx context.Context) UsersService {
	return &usersService{
		client:  client.MongoClient,
		context: ctx,
	}
}

func (svc *usersService) ListAll() ([]sdkgo.User, error) {
	clusters, _, err := svc.client.ClustersApi.ClustersGet(svc.context).Execute()
	if err != nil {
		return nil, err
	}

	var users []sdkgo.User
	for _, c := range *clusters.GetItems() {
		ls, _, err := svc.client.UsersApi.ClustersUsersGet(svc.context, *c.GetId()).Execute()
		if err != nil {
			return nil, err
		}
		users = append(users, *ls.GetItems()...)
	}

	return users, err
}

func (svc *usersService) List(clusterID string) (sdkgo.UsersList, *sdkgo.APIResponse, error) {
	return svc.client.UsersApi.ClustersUsersGet(svc.context, clusterID).Execute()
}

func (svc *usersService) Create(clusterID string, user sdkgo.User) (sdkgo.User, *sdkgo.APIResponse, error) {
	return svc.client.UsersApi.ClustersUsersPost(svc.context, clusterID).User(user).Execute()
}

func (svc *usersService) Get(clusterID, database, username string) (sdkgo.User, *sdkgo.APIResponse, error) {
	return svc.client.UsersApi.ClustersUsersFindById(svc.context, clusterID, database, username).Execute()
}

func (svc *usersService) Delete(clusterID, database, username string) (sdkgo.User, *sdkgo.APIResponse, error) {
	return svc.client.UsersApi.ClustersUsersDelete(svc.context, clusterID, database, username).Execute()
}
