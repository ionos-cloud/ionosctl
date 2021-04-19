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

// UsersService is a wrapper around ionoscloud.User
type UsersService interface {
	List() (Users, *Response, error)
	Get(userId string) (*User, *Response, error)
	Create(u User) (*User, *Response, error)
	Update(userId string, input User) (*User, *Response, error)
	Delete(userId string) (*Response, error)
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
