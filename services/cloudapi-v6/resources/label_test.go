package resources

import (
	"context"
	"testing"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"

	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/stretchr/testify/assert"
)

const (
	testLabelResourceVar = "test-label-resource"
)

func TestNewLabelResourceService(t *testing.T) {
	ctx := context.Background()
	t.Run("list_labels_error", func(t *testing.T) {
		svc := getTestClient(t)
		labelSvc := NewLabelResourceService(svc, ctx)
		_, _, err := labelSvc.List(testListQueryParam)
		assert.Error(t, err)
	})
	t.Run("get_label_error", func(t *testing.T) {
		svc := getTestClient(t)
		labelSvc := NewLabelResourceService(svc, ctx)
		_, _, err := labelSvc.GetByUrn(testLabelResourceVar)
		assert.Error(t, err)
	})
	t.Run("datacenter_list_labels_error", func(t *testing.T) {
		svc := getTestClient(t)
		labelSvc := NewLabelResourceService(svc, ctx)
		_, _, err := labelSvc.DatacenterList(testListQueryParam, testLabelResourceVar)
		assert.Error(t, err)
	})
	t.Run("datacenter_get_label_error", func(t *testing.T) {
		svc := getTestClient(t)
		labelSvc := NewLabelResourceService(svc, ctx)
		_, _, err := labelSvc.DatacenterGet(testLabelResourceVar, testLabelResourceVar)
		assert.Error(t, err)
	})
	t.Run("datacenter_create_label_error", func(t *testing.T) {
		svc := getTestClient(t)
		labelSvc := NewLabelResourceService(svc, ctx)
		_, _, err := labelSvc.DatacenterCreate(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar)
		assert.Error(t, err)
	})
	t.Run("datacenter_delete_label_error", func(t *testing.T) {
		svc := getTestClient(t)
		labelSvc := NewLabelResourceService(svc, ctx)
		_, err := labelSvc.DatacenterDelete(testLabelResourceVar, testLabelResourceVar)
		assert.Error(t, err)
	})
	t.Run("server_list_labels_error", func(t *testing.T) {
		svc := getTestClient(t)
		labelSvc := NewLabelResourceService(svc, ctx)
		_, _, err := labelSvc.ServerList(testListQueryParam, testLabelResourceVar, testLabelResourceVar)
		assert.Error(t, err)
	})
	t.Run("server_get_label_error", func(t *testing.T) {
		svc := getTestClient(t)
		labelSvc := NewLabelResourceService(svc, ctx)
		_, _, err := labelSvc.ServerGet(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar)
		assert.Error(t, err)
	})
	t.Run("server_create_label_error", func(t *testing.T) {
		svc := getTestClient(t)
		labelSvc := NewLabelResourceService(svc, ctx)
		_, _, err := labelSvc.ServerCreate(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar, testLabelResourceVar)
		assert.Error(t, err)
	})
	t.Run("server_delete_label_error", func(t *testing.T) {
		svc := getTestClient(t)
		labelSvc := NewLabelResourceService(svc, ctx)
		_, err := labelSvc.ServerDelete(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar)
		assert.Error(t, err)
	})
	t.Run("volume_list_labels_error", func(t *testing.T) {
		svc := getTestClient(t)
		labelSvc := NewLabelResourceService(svc, ctx)
		_, _, err := labelSvc.VolumeList(testListQueryParam, testLabelResourceVar, testLabelResourceVar)
		assert.Error(t, err)
	})
	t.Run("volume_get_label_error", func(t *testing.T) {
		svc := getTestClient(t)
		labelSvc := NewLabelResourceService(svc, ctx)
		_, _, err := labelSvc.VolumeGet(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar)
		assert.Error(t, err)
	})
	t.Run("volume_create_label_error", func(t *testing.T) {
		svc := getTestClient(t)
		labelSvc := NewLabelResourceService(svc, ctx)
		_, _, err := labelSvc.VolumeCreate(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar, testLabelResourceVar)
		assert.Error(t, err)
	})
	t.Run("volume_delete_label_error", func(t *testing.T) {
		svc := getTestClient(t)
		labelSvc := NewLabelResourceService(svc, ctx)
		_, err := labelSvc.VolumeDelete(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar)
		assert.Error(t, err)
	})
	t.Run("snapshot_list_labels_error", func(t *testing.T) {
		svc := getTestClient(t)
		labelSvc := NewLabelResourceService(svc, ctx)
		_, _, err := labelSvc.SnapshotList(testListQueryParam, testLabelResourceVar)
		assert.Error(t, err)
	})
	t.Run("snapshot_get_label_error", func(t *testing.T) {
		svc := getTestClient(t)
		labelSvc := NewLabelResourceService(svc, ctx)
		_, _, err := labelSvc.SnapshotGet(testLabelResourceVar, testLabelResourceVar)
		assert.Error(t, err)
	})
	t.Run("snapshot_create_label_error", func(t *testing.T) {
		svc := getTestClient(t)
		labelSvc := NewLabelResourceService(svc, ctx)
		_, _, err := labelSvc.SnapshotCreate(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar)
		assert.Error(t, err)
	})
	t.Run("snapshot_delete_label_error", func(t *testing.T) {
		svc := getTestClient(t)
		labelSvc := NewLabelResourceService(svc, ctx)
		_, err := labelSvc.SnapshotDelete(testLabelResourceVar, testLabelResourceVar)
		assert.Error(t, err)
	})
	t.Run("ipblock_list_labels_error", func(t *testing.T) {
		svc := getTestClient(t)
		labelSvc := NewLabelResourceService(svc, ctx)
		_, _, err := labelSvc.IpBlockList(testListQueryParam, testLabelResourceVar)
		assert.Error(t, err)
	})
	t.Run("ipblock_get_label_error", func(t *testing.T) {
		svc := getTestClient(t)
		labelSvc := NewLabelResourceService(svc, ctx)
		_, _, err := labelSvc.IpBlockGet(testLabelResourceVar, testLabelResourceVar)
		assert.Error(t, err)
	})
	t.Run("ipblock_create_label_error", func(t *testing.T) {
		svc := getTestClient(t)
		labelSvc := NewLabelResourceService(svc, ctx)
		_, _, err := labelSvc.IpBlockCreate(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar)
		assert.Error(t, err)
	})
	t.Run("ipblock_delete_label_error", func(t *testing.T) {
		svc := getTestClient(t)
		labelSvc := NewLabelResourceService(svc, ctx)
		_, err := labelSvc.IpBlockDelete(testLabelResourceVar, testLabelResourceVar)
		assert.Error(t, err)
	})
}

func getTestClient(t *testing.T) *client.Client {
	svc, err := client.NewTestClient("user", "pass", "", constants.DefaultApiURL)
	assert.NotNil(t, svc)
	assert.NoError(t, err)
	assert.Equal(t, "user", svc.CloudClient.GetConfig().Username)
	assert.Equal(t, "pass", svc.CloudClient.GetConfig().Password)
	return svc
}
