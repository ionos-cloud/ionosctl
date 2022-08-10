package resources

import (
	"context"

	dp "github.com/ionos-cloud/sdk-go-autoscaling"
)

type NodePoolResponseData struct {
	dp.NodePoolResponseData
}

type NodePoolListResponseData struct {
	dp.NodePoolListResponseData
}

type CreateNodePoolRequest struct {
	dp.CreateNodePoolRequest
}

type CreateNodePoolProperties struct {
	dp.CreateNodePoolProperties
}

type PatchNodePoolRequest struct {
	dp.PatchNodePoolRequest
}

type PatchNodePoolProperties struct {
	dp.PatchNodePoolProperties
}

type AvailabilityZone struct {
	dp.AvailabilityZone
}

type StorageType struct {
	dp.StorageType
}

// NodePoolsService is a wrapper around dp.NodePool
type NodePoolsService interface {
	Get(clusterId, nodePoolId string) (NodePoolResponseData, *Response, error)
	List(clusterId string) (NodePoolListResponseData, *Response, error)
	Create(clusterId string, nodePool CreateNodePoolRequest) (NodePoolResponseData, *Response, error)
	Update(clusterId, nodePoolId string, cluster PatchNodePoolRequest) (NodePoolResponseData, *Response, error)
	Delete(clusterId, nodePoolId string) (NodePoolResponseData, *Response, error)
}

type nodePoolsService struct {
	client  *Client
	context context.Context
}

var _ NodePoolsService = &nodePoolsService{}

func NewNodePoolsService(client *Client, ctx context.Context) NodePoolsService {
	return &nodePoolsService{
		client:  client,
		context: ctx,
	}
}

func (svc *nodePoolsService) Get(clusterId, nodePoolId string) (NodePoolResponseData, *Response, error) {
	req := svc.client.DataPlatformNodePoolApi.GetClusterNodepool(svc.context, clusterId, nodePoolId)
	nodePoolResponse, res, err := svc.client.DataPlatformNodePoolApi.GetClusterNodepoolExecute(req)
	return NodePoolResponseData{nodePoolResponse}, &Response{*res}, err
}

func (svc *nodePoolsService) List(clusterId string) (NodePoolListResponseData, *Response, error) {
	request := svc.client.DataPlatformNodePoolApi.GetClusterNodepools(svc.context, clusterId)
	nodePoolListResponse, res, err := svc.client.DataPlatformNodePoolApi.GetClusterNodepoolsExecute(request)

	return NodePoolListResponseData{nodePoolListResponse}, &Response{*res}, err
}

func (svc *nodePoolsService) Create(clusterId string, nodePool CreateNodePoolRequest) (NodePoolResponseData, *Response, error) {
	req := svc.client.DataPlatformNodePoolApi.CreateClusterNodepool(svc.context, clusterId).CreateNodePoolRequest(nodePool.CreateNodePoolRequest)
	nodePoolResponse, res, err := svc.client.DataPlatformNodePoolApi.CreateClusterNodepoolExecute(req)
	return NodePoolResponseData{nodePoolResponse}, &Response{*res}, err
}

func (svc *nodePoolsService) Update(clusterId, nodePoolId string, nodePool PatchNodePoolRequest) (NodePoolResponseData, *Response, error) {
	req := svc.client.DataPlatformNodePoolApi.PatchClusterNodepool(svc.context, clusterId, nodePoolId).PatchNodePoolRequest(nodePool.PatchNodePoolRequest)
	nodePoolResponse, res, err := svc.client.DataPlatformNodePoolApi.PatchClusterNodepoolExecute(req)
	return NodePoolResponseData{nodePoolResponse}, &Response{*res}, err
}

func (svc *nodePoolsService) Delete(clusterId, nodePoolId string) (NodePoolResponseData, *Response, error) {
	req := svc.client.DataPlatformNodePoolApi.DeleteClusterNodepool(svc.context, clusterId, nodePoolId)
	nodePoolResponse, res, err := svc.client.DataPlatformNodePoolApi.DeleteClusterNodepoolExecute(req)
	return NodePoolResponseData{nodePoolResponse}, &Response{*res}, err
}
