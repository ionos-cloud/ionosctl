package resources

import (
	"context"

	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
)

type User struct {
	ionoscloud.User
}

type UserProperties struct {
	ionoscloud.UserProperties
}

type Users struct {
	ionoscloud.Users
}

type Resource struct {
	ionoscloud.Resource
}

type Resources struct {
	ionoscloud.Resources
}

// UsersService is a wrapper around ionoscloud.User
type UsersService interface {
	List() (Users, *Response, error)
	Get(userId string) (*User, *Response, error)
	Create(u User) (*User, *Response, error)
	Update(userId string, input User) (*User, *Response, error)
	Delete(userId string) (*Response, error)
	// Resources
	ListResources() (Resources, *Response, error)
	GetResourcesByType(resourceType string) (Resources, *Response, error)
	GetResourceByTypeAndId(resourceType, resourceId string) (*Resource, *Response, error)
}

type usersService struct {
	client  *Client
	context context.Context
}

var _ UsersService = &usersService{}

func NewUserService(client *Client, ctx context.Context) UsersService {
	return &usersService{
		client:  client,
		context: ctx,
	}
}

func (s *usersService) List() (Users, *Response, error) {
	req := s.client.UserManagementApi.UmUsersGet(s.context)
	dcs, res, err := s.client.UserManagementApi.UmUsersGetExecute(req)
	return Users{dcs}, &Response{*res}, err
}

func (s *usersService) Get(userId string) (*User, *Response, error) {
	req := s.client.UserManagementApi.UmUsersFindById(s.context, userId)
	user, res, err := s.client.UserManagementApi.UmUsersFindByIdExecute(req)
	return &User{user}, &Response{*res}, err
}

func (s *usersService) Create(u User) (*User, *Response, error) {
	req := s.client.UserManagementApi.UmUsersPost(s.context).User(u.User)
	user, res, err := s.client.UserManagementApi.UmUsersPostExecute(req)
	return &User{user}, &Response{*res}, err
}

func (s *usersService) Update(userId string, input User) (*User, *Response, error) {
	req := s.client.UserManagementApi.UmUsersPut(s.context, userId).User(input.User)
	user, res, err := s.client.UserManagementApi.UmUsersPutExecute(req)
	return &User{user}, &Response{*res}, err
}

func (s *usersService) Delete(userId string) (*Response, error) {
	req := s.client.UserManagementApi.UmUsersDelete(s.context, userId)
	_, res, err := s.client.UserManagementApi.UmUsersDeleteExecute(req)
	return &Response{*res}, err
}

func (s *usersService) ListResources() (Resources, *Response, error) {
	req := s.client.UserManagementApi.UmResourcesGet(s.context)
	gs, res, err := s.client.UserManagementApi.UmResourcesGetExecute(req)
	return Resources{gs}, &Response{*res}, err
}

func (s *usersService) GetResourcesByType(resourceType string) (Resources, *Response, error) {
	req := s.client.UserManagementApi.UmResourcesFindByType(s.context, resourceType)
	ss, res, err := s.client.UserManagementApi.UmResourcesFindByTypeExecute(req)
	return Resources{ss}, &Response{*res}, err
}

func (s *usersService) GetResourceByTypeAndId(resourceType, resourceId string) (*Resource, *Response, error) {
	req := s.client.UserManagementApi.UmResourcesFindByTypeAndId(s.context, resourceType, resourceId)
	ss, res, err := s.client.UserManagementApi.UmResourcesFindByTypeAndIdExecute(req)
	return &Resource{ss}, &Response{*res}, err
}
