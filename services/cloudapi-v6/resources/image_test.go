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
		imageSvc := NewImageService(svc.Get(), ctx)
		_, _, err := imageSvc.List(ListQueryParams{})
		assert.Error(t, err)
	})
	t.Run("list_images_filters_error", func(t *testing.T) {
		svc := getTestClient(t)
		imageSvc := NewImageService(svc.Get(), ctx)
		_, _, err := imageSvc.List(testListQueryParam)
		assert.Error(t, err)
	})
	t.Run("get_image_error", func(t *testing.T) {
		svc := getTestClient(t)
		imageSvc := NewImageService(svc.Get(), ctx)
		_, _, err := imageSvc.Get(testImageResourceVar, QueryParams{})
		assert.Error(t, err)
	})
	t.Run("update_image_error", func(t *testing.T) {
		svc := getTestClient(t)
		imageSvc := NewImageService(svc.Get(), ctx)
		_, _, err := imageSvc.Update(
			testImageResourceVar,
			ImageProperties{},
			QueryParams{},
		)
		assert.Error(t, err)
	})
	t.Run("delete_image_error", func(t *testing.T) {
		svc := getTestClient(t)
		imageSvc := NewImageService(svc.Get(), ctx)
		_, err := imageSvc.Delete(testImageResourceVar, QueryParams{})
		assert.Error(t, err)
	})
}
