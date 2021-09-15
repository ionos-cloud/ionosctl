package resources

import (
	"context"

	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
)

type TargetGroup struct {
	ionoscloud.TargetGroup
}

type TargetGroupTarget struct {
	ionoscloud.TargetGroupTarget
}

type TargetGroupTargetHealthCheck struct {
	ionoscloud.TargetGroupTargetHealthCheck
}

type TargetGroupProperties struct {
	ionoscloud.TargetGroupProperties
}

type TargetGroupHealthCheck struct {
	ionoscloud.TargetGroupHealthCheck
}

type TargetGroupHttpHealthCheck struct {
	ionoscloud.TargetGroupHttpHealthCheck
}

type TargetGroups struct {
	ionoscloud.TargetGroups
}

// TargetGroupsService is a wrapper around ionoscloud.TargetGroup
type TargetGroupsService interface {
	List() (TargetGroups, *Response, error)
	Get(targetGroupId string) (*TargetGroup, *Response, error)
	Create(tg TargetGroup) (*TargetGroup, *Response, error)
	Update(targetGroupId string, input *TargetGroupProperties) (*TargetGroup, *Response, error)
	Delete(targetGroupId string) (*Response, error)
}

type targetGroupsService struct {
	client  *Client
	context context.Context
}

var _ TargetGroupsService = &targetGroupsService{}

func NewTargetGroupService(client *Client, ctx context.Context) TargetGroupsService {
	return &targetGroupsService{
		client:  client,
		context: ctx,
	}
}

func (svc *targetGroupsService) List() (TargetGroups, *Response, error) {
	req := svc.client.TargetGroupsApi.TargetgroupsGet(svc.context)
	dcs, res, err := svc.client.TargetGroupsApi.TargetgroupsGetExecute(req)
	return TargetGroups{dcs}, &Response{*res}, err
}

func (svc *targetGroupsService) Get(targetGroupId string) (*TargetGroup, *Response, error) {
	req := svc.client.TargetGroupsApi.TargetgroupsFindByTargetGroupId(svc.context, targetGroupId)
	targetGroup, res, err := svc.client.TargetGroupsApi.TargetgroupsFindByTargetGroupIdExecute(req)
	return &TargetGroup{targetGroup}, &Response{*res}, err
}

func (svc *targetGroupsService) Create(tg TargetGroup) (*TargetGroup, *Response, error) {
	req := svc.client.TargetGroupsApi.TargetgroupsPost(svc.context).TargetGroup(tg.TargetGroup)
	targetGroup, res, err := svc.client.TargetGroupsApi.TargetgroupsPostExecute(req)
	return &TargetGroup{targetGroup}, &Response{*res}, err
}

func (svc *targetGroupsService) Update(targetGroupId string, input *TargetGroupProperties) (*TargetGroup, *Response, error) {
	req := svc.client.TargetGroupsApi.TargetgroupsPatch(svc.context, targetGroupId).TargetGroupProperties(input.TargetGroupProperties)
	targetGroup, res, err := svc.client.TargetGroupsApi.TargetgroupsPatchExecute(req)
	return &TargetGroup{targetGroup}, &Response{*res}, err
}

func (svc *targetGroupsService) Delete(targetGroupId string) (*Response, error) {
	req := svc.client.TargetGroupsApi.TargetGroupsDelete(context.Background(), targetGroupId)
	res, err := svc.client.TargetGroupsApi.TargetGroupsDeleteExecute(req)
	return &Response{*res}, err
}
