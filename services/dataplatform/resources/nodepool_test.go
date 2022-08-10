package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testNodePoolResourceVar = "test-node-pool-resource"
)

func TestNewNodePoolsService(t *testing.T) {
	ctx := context.Background()
	t.Run("list_node_pools_error", func(t *testing.T) {
		svc := getTestClient(t)
		nodePoolSvc := NewNodePoolsService(svc.Get(), ctx)
		_, _, err := nodePoolSvc.List(testClusterResourceVar)
		assert.Error(t, err)
	})
	t.Run("get_node_pools_error", func(t *testing.T) {
		svc := getTestClient(t)
		nodePoolSvc := NewNodePoolsService(svc.Get(), ctx)
		_, _, err := nodePoolSvc.Get(testClusterResourceVar, testNodePoolResourceVar)
		assert.Error(t, err)
	})
	t.Run("create_node_pool_error", func(t *testing.T) {
		svc := getTestClient(t)
		nodePoolSvc := NewNodePoolsService(svc.Get(), ctx)
		_, _, err := nodePoolSvc.Create(testClusterResourceVar, CreateNodePoolRequest{})
		assert.Error(t, err)
	})
	t.Run("update_node_pool_error", func(t *testing.T) {
		svc := getTestClient(t)
		nodePoolSvc := NewNodePoolsService(svc.Get(), ctx)
		_, _, err := nodePoolSvc.Update(testClusterResourceVar, testNodePoolResourceVar, PatchNodePoolRequest{})
		assert.Error(t, err)
	})
	t.Run("delete_node_pool_error", func(t *testing.T) {
		svc := getTestClient(t)
		nodePoolSvc := NewNodePoolsService(svc.Get(), ctx)
		_, _, err := nodePoolSvc.Delete(testClusterResourceVar, testNodePoolResourceVar)
		assert.Error(t, err)
	})
}
