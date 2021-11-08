package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testK8sResourceVar = "test-k8s-resource"
)

func TestNewK8sService(t *testing.T) {
	ctx := context.Background()
	t.Run("listclusters_k8s_error", func(t *testing.T) {
		svc := getTestClient(t)
		k8sSvc := NewK8sService(svc.Get(), ctx)
		_, _, err := k8sSvc.ListClusters(ListQueryParams{})
		assert.Error(t, err)
	})
	t.Run("getcluster_k8s_error", func(t *testing.T) {
		svc := getTestClient(t)
		k8sSvc := NewK8sService(svc.Get(), ctx)
		_, _, err := k8sSvc.GetCluster(testK8sResourceVar)
		assert.Error(t, err)
	})
	t.Run("createcluster_k8s_error", func(t *testing.T) {
		svc := getTestClient(t)
		k8sSvc := NewK8sService(svc.Get(), ctx)
		_, _, err := k8sSvc.CreateCluster(K8sClusterForPost{})
		assert.Error(t, err)
	})
	t.Run("updatecluster_k8s_error", func(t *testing.T) {
		svc := getTestClient(t)
		k8sSvc := NewK8sService(svc.Get(), ctx)
		_, _, err := k8sSvc.UpdateCluster(testK8sResourceVar, K8sClusterForPut{})
		assert.Error(t, err)
	})
	t.Run("deletecluster_k8s_error", func(t *testing.T) {
		svc := getTestClient(t)
		k8sSvc := NewK8sService(svc.Get(), ctx)
		_, err := k8sSvc.DeleteCluster(testK8sResourceVar)
		assert.Error(t, err)
	})
	t.Run("readkubeconfig_k8s_error", func(t *testing.T) {
		svc := getTestClient(t)
		k8sSvc := NewK8sService(svc.Get(), ctx)
		_, _, err := k8sSvc.ReadKubeConfig(testK8sResourceVar)
		assert.Error(t, err)
	})
	t.Run("listnodepools_k8s_error", func(t *testing.T) {
		svc := getTestClient(t)
		k8sSvc := NewK8sService(svc.Get(), ctx)
		_, _, err := k8sSvc.ListNodePools(testK8sResourceVar)
		assert.Error(t, err)
	})
	t.Run("getnodepool_k8s_error", func(t *testing.T) {
		svc := getTestClient(t)
		k8sSvc := NewK8sService(svc.Get(), ctx)
		_, _, err := k8sSvc.GetNodePool(testK8sResourceVar, testK8sResourceVar)
		assert.Error(t, err)
	})
	t.Run("createnodepool_k8s_error", func(t *testing.T) {
		svc := getTestClient(t)
		k8sSvc := NewK8sService(svc.Get(), ctx)
		_, _, err := k8sSvc.CreateNodePool(testK8sResourceVar, K8sNodePool{})
		assert.Error(t, err)
	})
	t.Run("updatenodepool_k8s_error", func(t *testing.T) {
		svc := getTestClient(t)
		k8sSvc := NewK8sService(svc.Get(), ctx)
		_, _, err := k8sSvc.UpdateNodePool(testK8sResourceVar, testK8sResourceVar, K8sNodePoolForPut{})
		assert.Error(t, err)
	})
	t.Run("deletenodepool_k8s_error", func(t *testing.T) {
		svc := getTestClient(t)
		k8sSvc := NewK8sService(svc.Get(), ctx)
		_, err := k8sSvc.DeleteNodePool(testK8sResourceVar, testK8sResourceVar)
		assert.Error(t, err)
	})
	t.Run("recreatenode_k8s_error", func(t *testing.T) {
		svc := getTestClient(t)
		k8sSvc := NewK8sService(svc.Get(), ctx)
		_, err := k8sSvc.RecreateNode(testK8sResourceVar, testK8sResourceVar, testK8sResourceVar)
		assert.Error(t, err)
	})
	t.Run("getnode_k8s_error", func(t *testing.T) {
		svc := getTestClient(t)
		k8sSvc := NewK8sService(svc.Get(), ctx)
		_, _, err := k8sSvc.GetNode(testK8sResourceVar, testK8sResourceVar, testK8sResourceVar)
		assert.Error(t, err)
	})
	t.Run("listnodes_k8s_error", func(t *testing.T) {
		svc := getTestClient(t)
		k8sSvc := NewK8sService(svc.Get(), ctx)
		_, _, err := k8sSvc.ListNodes(testK8sResourceVar, testK8sResourceVar)
		assert.Error(t, err)
	})
	t.Run("deletenode_k8s_error", func(t *testing.T) {
		svc := getTestClient(t)
		k8sSvc := NewK8sService(svc.Get(), ctx)
		_, err := k8sSvc.DeleteNode(testK8sResourceVar, testK8sResourceVar, testK8sResourceVar)
		assert.Error(t, err)
	})
	t.Run("listversions_k8s_error", func(t *testing.T) {
		svc := getTestClient(t)
		k8sSvc := NewK8sService(svc.Get(), ctx)
		_, _, err := k8sSvc.ListVersions()
		assert.Error(t, err)
	})
	t.Run("getversion_k8s_error", func(t *testing.T) {
		svc := getTestClient(t)
		k8sSvc := NewK8sService(svc.Get(), ctx)
		_, _, err := k8sSvc.GetVersion()
		assert.Error(t, err)
	})
}
