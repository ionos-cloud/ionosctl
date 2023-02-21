package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testClusterResourceVar = "test-cluster-resource"

func TestNewClustersService(t *testing.T) {
	ctx := context.Background()
	t.Run("list_cluster_error", func(t *testing.T) {
		svc := getTestClient(t)
		clusterUnitSvc := NewClustersService(svc, ctx)
		_, _, err := clusterUnitSvc.List("")
		assert.Error(t, err)
	})
	t.Run("list_cluster_filter_error", func(t *testing.T) {
		svc := getTestClient(t)
		clusterUnitSvc := NewClustersService(svc, ctx)
		_, _, err := clusterUnitSvc.List(testClusterResourceVar)
		assert.Error(t, err)
	})
	t.Run("get_clusters_error", func(t *testing.T) {
		svc := getTestClient(t)
		clusterUnitSvc := NewClustersService(svc, ctx)
		_, _, err := clusterUnitSvc.Get(testClusterResourceVar)
		assert.Error(t, err)
	})
	t.Run("create_cluster_error", func(t *testing.T) {
		svc := getTestClient(t)
		clusterUnitSvc := NewClustersService(svc, ctx)
		_, _, err := clusterUnitSvc.Create(CreateClusterRequest{})
		assert.Error(t, err)
	})
	t.Run("update_cluster_error", func(t *testing.T) {
		svc := getTestClient(t)
		clusterUnitSvc := NewClustersService(svc, ctx)
		_, _, err := clusterUnitSvc.Update(testClusterResourceVar, PatchClusterRequest{})
		assert.Error(t, err)
	})
	t.Run("delete_cluster_error", func(t *testing.T) {
		svc := getTestClient(t)
		clusterUnitSvc := NewClustersService(svc, ctx)
		_, err := clusterUnitSvc.Delete(testClusterResourceVar)
		assert.Error(t, err)
	})
}
