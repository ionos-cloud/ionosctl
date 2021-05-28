package resources

import (
	"context"

	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
)

type Group struct {
	ionoscloud.Group
}

type GroupProperties struct {
	ionoscloud.GroupProperties
}

type Groups struct {
	ionoscloud.Groups
}

type GroupMembers struct {
	ionoscloud.GroupMembers
}

type GroupShare struct {
	ionoscloud.GroupShare
}

type GroupShareProperties struct {
	ionoscloud.GroupShareProperties
}

type GroupShares struct {
	ionoscloud.GroupShares
}

type ResourceGroups struct {
	ionoscloud.ResourceGroups
}

// GroupsService is a wrapper around ionoscloud.Group
type GroupsService interface {
	List() (Groups, *Response, error)
	Get(groupId string) (*Group, *Response, error)
	Create(u Group) (*Group, *Response, error)
	Update(groupId string, input Group) (*Group, *Response, error)
	Delete(groupId string) (*Response, error)
	ListUsers(groupId string) (GroupMembers, *Response, error)
	AddUser(groupId string, input User) (*User, *Response, error)
	RemoveUser(groupId, userId string) (*Response, error)
	ListShares(groupId string) (GroupShares, *Response, error)
	GetShare(groupId, resourceId string) (*GroupShare, *Response, error)
	AddShare(groupId, resourceId string, input GroupShare) (*GroupShare, *Response, error)
	UpdateShare(groupId, resourceId string, input GroupShare) (*GroupShare, *Response, error)
	RemoveShare(groupId, resourceId string) (*Response, error)
	ListResources(groupId string) (ResourceGroups, *Response, error)
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
	gs, res, err := s.client.UserManagementApi.UmGroupsGetExecute(req)
	return Groups{gs}, &Response{*res}, err
}

func (s *groupsService) Get(groupId string) (*Group, *Response, error) {
	req := s.client.UserManagementApi.UmGroupsFindById(s.context, groupId)
	group, res, err := s.client.UserManagementApi.UmGroupsFindByIdExecute(req)
	return &Group{group}, &Response{*res}, err
}

func (s *groupsService) Create(g Group) (*Group, *Response, error) {
	req := s.client.UserManagementApi.UmGroupsPost(s.context).Group(g.Group)
	group, res, err := s.client.UserManagementApi.UmGroupsPostExecute(req)
	return &Group{group}, &Response{*res}, err
}

func (s *groupsService) Update(groupId string, input Group) (*Group, *Response, error) {
	req := s.client.UserManagementApi.UmGroupsPut(s.context, groupId).Group(input.Group)
	group, res, err := s.client.UserManagementApi.UmGroupsPutExecute(req)
	return &Group{group}, &Response{*res}, err
}

func (s *groupsService) Delete(groupId string) (*Response, error) {
	req := s.client.UserManagementApi.UmGroupsDelete(s.context, groupId)
	res, err := s.client.UserManagementApi.UmGroupsDeleteExecute(req)
	return &Response{*res}, err
}

// Users

func (s *groupsService) ListUsers(groupId string) (GroupMembers, *Response, error) {
	req := s.client.UserManagementApi.UmGroupsUsersGet(s.context, groupId)
	groupMembers, res, err := s.client.UserManagementApi.UmGroupsUsersGetExecute(req)
	return GroupMembers{groupMembers}, &Response{*res}, err
}

func (s *groupsService) AddUser(groupId string, input User) (*User, *Response, error) {
	req := s.client.UserManagementApi.UmGroupsUsersPost(s.context, groupId).User(input.User)
	u, res, err := s.client.UserManagementApi.UmGroupsUsersPostExecute(req)
	return &User{u}, &Response{*res}, err
}

func (s *groupsService) RemoveUser(groupId, userId string) (*Response, error) {
	req := s.client.UserManagementApi.UmGroupsUsersDelete(s.context, groupId, userId)
	res, err := s.client.UserManagementApi.UmGroupsUsersDeleteExecute(req)
	return &Response{*res}, err
}

// Shares

func (s *groupsService) ListShares(groupId string) (GroupShares, *Response, error) {
	req := s.client.UserManagementApi.UmGroupsSharesGet(s.context, groupId)
	groupShares, res, err := s.client.UserManagementApi.UmGroupsSharesGetExecute(req)
	return GroupShares{groupShares}, &Response{*res}, err
}

func (s *groupsService) GetShare(groupId, resourceId string) (*GroupShare, *Response, error) {
	req := s.client.UserManagementApi.UmGroupsSharesFindByResourceId(s.context, groupId, resourceId)
	groupShare, res, err := s.client.UserManagementApi.UmGroupsSharesFindByResourceIdExecute(req)
	return &GroupShare{groupShare}, &Response{*res}, err
}

func (s *groupsService) AddShare(groupId, resourceId string, input GroupShare) (*GroupShare, *Response, error) {
	req := s.client.UserManagementApi.UmGroupsSharesPost(s.context, groupId, resourceId).Resource(input.GroupShare)
	groupShare, res, err := s.client.UserManagementApi.UmGroupsSharesPostExecute(req)
	return &GroupShare{groupShare}, &Response{*res}, err
}

func (s *groupsService) UpdateShare(groupId, resourceId string, input GroupShare) (*GroupShare, *Response, error) {
	req := s.client.UserManagementApi.UmGroupsSharesPut(s.context, groupId, resourceId).Resource(input.GroupShare)
	groupShare, res, err := s.client.UserManagementApi.UmGroupsSharesPutExecute(req)
	return &GroupShare{groupShare}, &Response{*res}, err
}

func (s *groupsService) RemoveShare(groupId, resourceId string) (*Response, error) {
	req := s.client.UserManagementApi.UmGroupsSharesDelete(s.context, groupId, resourceId)
	res, err := s.client.UserManagementApi.UmGroupsSharesDeleteExecute(req)
	return &Response{*res}, err
}

// Resources

func (s *groupsService) ListResources(groupId string) (ResourceGroups, *Response, error) {
	req := s.client.UserManagementApi.UmGroupsResourcesGet(s.context, groupId)
	groupResources, res, err := s.client.UserManagementApi.UmGroupsResourcesGetExecute(req)
	return ResourceGroups{groupResources}, &Response{*res}, err
}
