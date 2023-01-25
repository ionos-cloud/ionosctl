package resources

import (
	"context"
	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/pkg/config"
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
	List(params ListQueryParams) (Groups, *Response, error)
	Get(groupId string, params QueryParams) (*Group, *Response, error)
	Create(u Group, params QueryParams) (*Group, *Response, error)
	Update(groupId string, input Group, params QueryParams) (*Group, *Response, error)
	Delete(groupId string, params QueryParams) (*Response, error)
	ListUsers(groupId string, params ListQueryParams) (GroupMembers, *Response, error)
	AddUser(groupId string, input User, params QueryParams) (*User, *Response, error)
	RemoveUser(groupId, userId string, params QueryParams) (*Response, error)
	ListShares(groupId string, params ListQueryParams) (GroupShares, *Response, error)
	GetShare(groupId, resourceId string, params QueryParams) (*GroupShare, *Response, error)
	AddShare(groupId, resourceId string, input GroupShare, params QueryParams) (*GroupShare, *Response, error)
	UpdateShare(groupId, resourceId string, input GroupShare, params QueryParams) (*GroupShare, *Response, error)
	RemoveShare(groupId, resourceId string, params QueryParams) (*Response, error)
	ListResources(groupId string, params ListQueryParams) (ResourceGroups, *Response, error)
}

type groupsService struct {
	client  *config.Client
	context context.Context
}

var _ GroupsService = &groupsService{}

func NewGroupService(client *config.Client, ctx context.Context) GroupsService {
	return &groupsService{
		client:  client,
		context: ctx,
	}
}

func (s *groupsService) List(params ListQueryParams) (Groups, *Response, error) {
	req := s.client.UserManagementApi.UmGroupsGet(s.context)
	if !structs.IsZero(params) {
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
	gs, res, err := s.client.UserManagementApi.UmGroupsGetExecute(req)
	return Groups{gs}, &Response{*res}, err
}

func (s *groupsService) Get(groupId string, params QueryParams) (*Group, *Response, error) {
	req := s.client.UserManagementApi.UmGroupsFindById(s.context, groupId)
	if !structs.IsZero(params) {
		if params.Depth != nil {
			req = req.Depth(*params.Depth)
		}
		if params.Pretty != nil {
			// Currently not implemented
			req = req.Pretty(*params.Pretty)
		}
	}
	group, res, err := s.client.UserManagementApi.UmGroupsFindByIdExecute(req)
	return &Group{group}, &Response{*res}, err
}

func (s *groupsService) Create(g Group, params QueryParams) (*Group, *Response, error) {
	req := s.client.UserManagementApi.UmGroupsPost(s.context).Group(g.Group)
	if !structs.IsZero(params) {
		if params.Depth != nil {
			req = req.Depth(*params.Depth)
		}
		if params.Pretty != nil {
			// Currently not implemented
			req = req.Pretty(*params.Pretty)
		}
	}
	group, res, err := s.client.UserManagementApi.UmGroupsPostExecute(req)
	return &Group{group}, &Response{*res}, err
}

func (s *groupsService) Update(groupId string, input Group, params QueryParams) (*Group, *Response, error) {
	req := s.client.UserManagementApi.UmGroupsPut(s.context, groupId).Group(input.Group)
	if !structs.IsZero(params) {
		if params.Depth != nil {
			req = req.Depth(*params.Depth)
		}
		if params.Pretty != nil {
			// Currently not implemented
			req = req.Pretty(*params.Pretty)
		}
	}
	group, res, err := s.client.UserManagementApi.UmGroupsPutExecute(req)
	return &Group{group}, &Response{*res}, err
}

func (s *groupsService) Delete(groupId string, params QueryParams) (*Response, error) {
	req := s.client.UserManagementApi.UmGroupsDelete(s.context, groupId)
	if !structs.IsZero(params) {
		if params.Depth != nil {
			req = req.Depth(*params.Depth)
		}
		if params.Pretty != nil {
			// Currently not implemented
			req = req.Pretty(*params.Pretty)
		}
	}
	res, err := s.client.UserManagementApi.UmGroupsDeleteExecute(req)
	return &Response{*res}, err
}

// Users

func (s *groupsService) ListUsers(groupId string, params ListQueryParams) (GroupMembers, *Response, error) {
	req := s.client.UserManagementApi.UmGroupsUsersGet(s.context, groupId)
	if !structs.IsZero(params) {
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
	groupMembers, res, err := s.client.UserManagementApi.UmGroupsUsersGetExecute(req)
	return GroupMembers{groupMembers}, &Response{*res}, err
}

func (s *groupsService) AddUser(groupId string, input User, params QueryParams) (*User, *Response, error) {
	req := s.client.UserManagementApi.UmGroupsUsersPost(s.context, groupId).User(input.User)
	if !structs.IsZero(params) {
		if params.Depth != nil {
			req = req.Depth(*params.Depth)
		}
		if params.Pretty != nil {
			// Currently not implemented
			req = req.Pretty(*params.Pretty)
		}
	}
	u, res, err := s.client.UserManagementApi.UmGroupsUsersPostExecute(req)
	return &User{u}, &Response{*res}, err
}

func (s *groupsService) RemoveUser(groupId, userId string, params QueryParams) (*Response, error) {
	req := s.client.UserManagementApi.UmGroupsUsersDelete(s.context, groupId, userId)
	if !structs.IsZero(params) {
		if params.Depth != nil {
			req = req.Depth(*params.Depth)
		}
		if params.Pretty != nil {
			// Currently not implemented
			req = req.Pretty(*params.Pretty)
		}
	}
	res, err := s.client.UserManagementApi.UmGroupsUsersDeleteExecute(req)
	return &Response{*res}, err
}

// Shares

func (s *groupsService) ListShares(groupId string, params ListQueryParams) (GroupShares, *Response, error) {
	req := s.client.UserManagementApi.UmGroupsSharesGet(s.context, groupId)
	if !structs.IsZero(params) {
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
	groupShares, res, err := s.client.UserManagementApi.UmGroupsSharesGetExecute(req)
	return GroupShares{groupShares}, &Response{*res}, err
}

func (s *groupsService) GetShare(groupId, resourceId string, params QueryParams) (*GroupShare, *Response, error) {
	req := s.client.UserManagementApi.UmGroupsSharesFindByResourceId(s.context, groupId, resourceId)
	if !structs.IsZero(params) {
		if params.Depth != nil {
			req = req.Depth(*params.Depth)
		}
		if params.Pretty != nil {
			// Currently not implemented
			req = req.Pretty(*params.Pretty)
		}
	}
	groupShare, res, err := s.client.UserManagementApi.UmGroupsSharesFindByResourceIdExecute(req)
	return &GroupShare{groupShare}, &Response{*res}, err
}

func (s *groupsService) AddShare(groupId, resourceId string, input GroupShare, params QueryParams) (*GroupShare, *Response, error) {
	req := s.client.UserManagementApi.UmGroupsSharesPost(s.context, groupId, resourceId).Resource(input.GroupShare)
	if !structs.IsZero(params) {
		if params.Depth != nil {
			req = req.Depth(*params.Depth)
		}
		if params.Pretty != nil {
			// Currently not implemented
			req = req.Pretty(*params.Pretty)
		}
	}
	groupShare, res, err := s.client.UserManagementApi.UmGroupsSharesPostExecute(req)
	return &GroupShare{groupShare}, &Response{*res}, err
}

func (s *groupsService) UpdateShare(groupId, resourceId string, input GroupShare, params QueryParams) (*GroupShare, *Response, error) {
	req := s.client.UserManagementApi.UmGroupsSharesPut(s.context, groupId, resourceId).Resource(input.GroupShare)
	if !structs.IsZero(params) {
		if params.Depth != nil {
			req = req.Depth(*params.Depth)
		}
		if params.Pretty != nil {
			// Currently not implemented
			req = req.Pretty(*params.Pretty)
		}
	}
	groupShare, res, err := s.client.UserManagementApi.UmGroupsSharesPutExecute(req)
	return &GroupShare{groupShare}, &Response{*res}, err
}

func (s *groupsService) RemoveShare(groupId, resourceId string, params QueryParams) (*Response, error) {
	req := s.client.UserManagementApi.UmGroupsSharesDelete(s.context, groupId, resourceId)
	if !structs.IsZero(params) {
		if params.Depth != nil {
			req = req.Depth(*params.Depth)
		}
		if params.Pretty != nil {
			// Currently not implemented
			req = req.Pretty(*params.Pretty)
		}
	}
	res, err := s.client.UserManagementApi.UmGroupsSharesDeleteExecute(req)
	return &Response{*res}, err
}

// Resources

func (s *groupsService) ListResources(groupId string, params ListQueryParams) (ResourceGroups, *Response, error) {
	req := s.client.UserManagementApi.UmGroupsResourcesGet(s.context, groupId)
	if !structs.IsZero(params) {
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
	groupResources, res, err := s.client.UserManagementApi.UmGroupsResourcesGetExecute(req)
	return ResourceGroups{groupResources}, &Response{*res}, err
}
