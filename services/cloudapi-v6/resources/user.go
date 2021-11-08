package resources

import (
	"context"

	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
)

type UserPost struct {
	ionoscloud.UserPost
}

type UserPut struct {
	ionoscloud.UserPut
}

type User struct {
	ionoscloud.User
}

type UserProperties struct {
	ionoscloud.UserProperties
}

type UserPropertiesPut struct {
	ionoscloud.UserPropertiesPut
}

type UserPropertiesPost struct {
	ionoscloud.UserPropertiesPost
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
	List(params ListQueryParams) (Users, *Response, error)
	Get(userId string) (*User, *Response, error)
	Create(u UserPost) (*User, *Response, error)
	Update(userId string, input UserPut) (*User, *Response, error)
	Delete(userId string) (*Response, error)
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

func (s *usersService) List(params ListQueryParams) (Users, *Response, error) {
	req := s.client.UserManagementApi.UmUsersGet(s.context)
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
	dcs, res, err := s.client.UserManagementApi.UmUsersGetExecute(req)
	return Users{dcs}, &Response{*res}, err
}

func (s *usersService) Get(userId string) (*User, *Response, error) {
	req := s.client.UserManagementApi.UmUsersFindById(s.context, userId)
	user, res, err := s.client.UserManagementApi.UmUsersFindByIdExecute(req)
	return &User{user}, &Response{*res}, err
}

func (s *usersService) Create(u UserPost) (*User, *Response, error) {
	req := s.client.UserManagementApi.UmUsersPost(s.context).User(u.UserPost)
	user, res, err := s.client.UserManagementApi.UmUsersPostExecute(req)
	return &User{user}, &Response{*res}, err
}

func (s *usersService) Update(userId string, input UserPut) (*User, *Response, error) {
	req := s.client.UserManagementApi.UmUsersPut(s.context, userId).User(input.UserPut)
	user, res, err := s.client.UserManagementApi.UmUsersPutExecute(req)
	return &User{user}, &Response{*res}, err
}

func (s *usersService) Delete(userId string) (*Response, error) {
	req := s.client.UserManagementApi.UmUsersDelete(s.context, userId)
	res, err := s.client.UserManagementApi.UmUsersDeleteExecute(req)
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
