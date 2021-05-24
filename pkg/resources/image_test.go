package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testImageResourceVar = "test-image-resource"
)

func TestNewImageService(t *testing.T) {
	ctx := context.Background()
	t.Run("list_images_error", func(t *testing.T) {
		svc := getTestClient(t)
		backupUnitSvc := NewImageService(svc.Get(), ctx)
		_, _, err := backupUnitSvc.List()
		assert.Error(t, err)
	})
	t.Run("get_image_error", func(t *testing.T) {
		svc := getTestClient(t)
		backupUnitSvc := NewImageService(svc.Get(), ctx)
		_, _, err := backupUnitSvc.Get(testImageResourceVar)
		assert.Error(t, err)
	})
	t.Run("update_image_error", func(t *testing.T) {
		svc := getTestClient(t)
		backupUnitSvc := NewImageService(svc.Get(), ctx)
		_, _, err := backupUnitSvc.Update(
			testImageResourceVar,
			ImageProperties{},
		)
		assert.Error(t, err)
	})
	t.Run("delete_image_error", func(t *testing.T) {
		svc := getTestClient(t)
		backupUnitSvc := NewImageService(svc.Get(), ctx)
		_, err := backupUnitSvc.Delete(testImageResourceVar)
		assert.Error(t, err)
	})
}
