package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testVolumeResourceVar = "test-volume-resource"
)

func TestNewVolumeService(t *testing.T) {
	ctx := context.Background()
	t.Run("list_volumes_error", func(t *testing.T) {
		svc := getTestClient(t)
		volumeSvc := NewVolumeService(svc.Get(), ctx)
		_, _, err := volumeSvc.List(testVolumeResourceVar, ListQueryParams{})
		assert.Error(t, err)
	})
	t.Run("get_volume_error", func(t *testing.T) {
		svc := getTestClient(t)
		volumeSvc := NewVolumeService(svc.Get(), ctx)
		_, _, err := volumeSvc.Get(testVolumeResourceVar, testVolumeResourceVar)
		assert.Error(t, err)
	})
	t.Run("create_volume_error", func(t *testing.T) {
		svc := getTestClient(t)
		volumeSvc := NewVolumeService(svc.Get(), ctx)
		_, _, err := volumeSvc.Create(testVolumeResourceVar, Volume{})
		assert.Error(t, err)
	})
	t.Run("update_volume_error", func(t *testing.T) {
		svc := getTestClient(t)
		volumeSvc := NewVolumeService(svc.Get(), ctx)
		_, _, err := volumeSvc.Update(testVolumeResourceVar, testVolumeResourceVar, VolumeProperties{})
		assert.Error(t, err)
	})
	t.Run("delete_volume_error", func(t *testing.T) {
		svc := getTestClient(t)
		volumeSvc := NewVolumeService(svc.Get(), ctx)
		_, err := volumeSvc.Delete(testVolumeResourceVar, testVolumeResourceVar)
		assert.Error(t, err)
	})
}
