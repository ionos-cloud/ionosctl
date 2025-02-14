package resources

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	compute "github.com/ionos-cloud/sdk-go/v6"

	"github.com/fatih/structs"
)

type UserPost struct {
	compute.UserPost
}

type UserPut struct {
	compute.UserPut
}

type User struct {
	compute.User
}

type UserProperties struct {
	compute.UserProperties
}

type UserPropertiesPut struct {
	compute.UserPropertiesPut
}

type UserPropertiesPost struct {
	compute.UserPropertiesPost
}

type Users struct {
	compute.Users
}

type Resource struct {
	compute.Resource
}

type Resources struct {
	compute.Resources
}

// UsersService is a wrapper around compute.User
type UsersService interface {
	List(params ListQueryParams) (Users, *Response, error)
	Get(userId string, params QueryParams) (*User, *Response, error)
	Create(u UserPost, params QueryParams) (*User, *Response, error)
	Update(userId string, input UserPut, params QueryParams) (*User, *Response, error)
	Delete(userId string, params QueryParams) (*Response, error)
	ListResources() (Resources, *Response, error)
	GetResourcesByType(resourceType string) (Resources, *Response, error)
	GetResourceByTypeAndId(resourceType, resourceId string) (*Resource, *Response, error)
}

type usersService struct {
	client  *compute.APIClient
	context context.Context
}

var _ UsersService = &usersService{}

func NewUserService(client *client.Client, ctx context.Context) UsersService {
	return &usersService{
		client:  client.CloudClient,
		context: ctx,
	}
}

func (s *usersService) List(params ListQueryParams) (Users, *Response, error) {
	req := s.client.UserManagementApi.UmUsersGet(s.context)
	if !structs.IsZero(params) {
		if params.Filters != nil {
			for k, v := range *params.Filters {
				for _, val := range v {
					req = req.Filter(k, val)
				}
			}
		}
		if params.OrderBy != nil {
			req = req.OrderBy(*params.OrderBy)
		}
		if params.MaxResults != nil {
			req = req.MaxResults(*params.MaxResults)
		}
		if !structs.IsZero(params.QueryParams) {
			if params.QueryParams.Depth != nil {
				req = req.Depth(*params.QueryParams.Depth)
			}
			if params.QueryParams.Pretty != nil {
				// Currently not implemented
				req = req.Pretty(*params.QueryParams.Pretty)
			}
		}
	}
	dcs, res, err := s.client.UserManagementApi.UmUsersGetExecute(req)
	return Users{dcs}, &Response{*res}, err
}

func (s *usersService) Get(userId string, params QueryParams) (*User, *Response, error) {
	req := s.client.UserManagementApi.UmUsersFindById(s.context, userId)
	if !structs.IsZero(params) {
		if params.Depth != nil {
			req = req.Depth(*params.Depth)
		}
		if params.Pretty != nil {
			// Currently not implemented
			req = req.Pretty(*params.Pretty)
		}
	}
	user, res, err := s.client.UserManagementApi.UmUsersFindByIdExecute(req)
	return &User{user}, &Response{*res}, err
}

func (s *usersService) Create(u UserPost, params QueryParams) (*User, *Response, error) {
	req := s.client.UserManagementApi.UmUsersPost(s.context).User(u.UserPost)
	if !structs.IsZero(params) {
		if params.Depth != nil {
			req = req.Depth(*params.Depth)
		}
		if params.Pretty != nil {
			// Currently not implemented
			req = req.Pretty(*params.Pretty)
		}
	}
	user, res, err := s.client.UserManagementApi.UmUsersPostExecute(req)
	return &User{user}, &Response{*res}, err
}

func (s *usersService) Update(userId string, input UserPut, params QueryParams) (*User, *Response, error) {
	req := s.client.UserManagementApi.UmUsersPut(s.context, userId).User(input.UserPut)
	if !structs.IsZero(params) {
		if params.Depth != nil {
			req = req.Depth(*params.Depth)
		}
		if params.Pretty != nil {
			// Currently not implemented
			req = req.Pretty(*params.Pretty)
		}
	}
	user, res, err := s.client.UserManagementApi.UmUsersPutExecute(req)
	return &User{user}, &Response{*res}, err
}

func (s *usersService) Delete(userId string, params QueryParams) (*Response, error) {
	req := s.client.UserManagementApi.UmUsersDelete(s.context, userId)
	if !structs.IsZero(params) {
		if params.Depth != nil {
			req = req.Depth(*params.Depth)
		}
		if params.Pretty != nil {
			// Currently not implemented
			req = req.Pretty(*params.Pretty)
		}
	}
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
