package resources

import (
	"context"

	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
)

type Group struct {
	ionoscloud.Group
}

type Groups struct {
	ionoscloud.Groups
}

// GroupsService is a wrapper around ionoscloud.Group
type GroupsService interface {
	List() (Groups, *Response, error)
	Get(groupId string) (*Group, *Response, error)
	Create(u *Group) (*Group, *Response, error)
	Update(groupId string, input *Group) (*Group, *Response, error)
	Delete(groupId string) (*Response, error)
	// Users
	ListUsers()
	AddUser()
	RemoveUser()
	// Shares
	ListShares()
	GetShare()
	AddShare()
	UpdateShare()
	DeleteShare()
	// Resources
	ListResources()
}

type groupsService struct {
	client  *Client
	context context.Context
}

var _ GroupsService = &groupsService{}

func NewGroupService(client *Client, ctx context.Context) GroupsService {
	return &groupsService{
		client:  client,
		context: ctx,
	}
}

func (s *groupsService) List() (Groups, *Response, error) {
	req := s.client.UserManagementApi.UmGroupsGet(s.context)
	dcs, res, err := s.client.UserManagementApi.UmGroupsGetExecute(req)
	return Groups{dcs}, &Response{*res}, err
}

func (s *groupsService) Get(groupId string) (*Group, *Response, error) {
	req := s.client.UserManagementApi.UmGroupsFindById(s.context, groupId)
	group, res, err := s.client.UserManagementApi.UmGroupsFindByIdExecute(req)
	return &Group{group}, &Response{*res}, err
}

func (s *groupsService) Create(g *Group) (*Group, *Response, error) {
	req := s.client.UserManagementApi.UmGroupsPost(s.context).Group(g.Group)
	group, res, err := s.client.UserManagementApi.UmGroupsPostExecute(req)
	return &Group{group}, &Response{*res}, err
}

func (s *groupsService) Update(groupId string, input *Group) (*Group, *Response, error) {
	req := s.client.UserManagementApi.UmGroupsPut(s.context, groupId).Group(input.Group)
	group, res, err := s.client.UserManagementApi.UmGroupsPutExecute(req)
	return &Group{group}, &Response{*res}, err
}

func (s *groupsService) Delete(groupId string) (*Response, error) {
	req := s.client.UserManagementApi.UmGroupsDelete(s.context, groupId)
	_, res, err := s.client.UserManagementApi.UmGroupsDeleteExecute(req)
	return &Response{*res}, err
}
