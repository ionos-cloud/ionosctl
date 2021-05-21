package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testS3KeyResourceVar = "test-s3key-resource"
)

func TestNewS3KeyService(t *testing.T) {
	ctx := context.Background()
	t.Run("list_s3keys_error", func(t *testing.T) {
		svc := getTestClient(t)
		backupUnitSvc := NewS3KeyService(svc.Get(), ctx)
		_, _, err := backupUnitSvc.List(testS3KeyResourceVar)
		assert.Error(t, err)
	})
	t.Run("get_s3key_error", func(t *testing.T) {
		svc := getTestClient(t)
		backupUnitSvc := NewS3KeyService(svc.Get(), ctx)
		_, _, err := backupUnitSvc.Get(testS3KeyResourceVar, testS3KeyResourceVar)
		assert.Error(t, err)
	})
	t.Run("create_s3key_error", func(t *testing.T) {
		svc := getTestClient(t)
		backupUnitSvc := NewS3KeyService(svc.Get(), ctx)
		_, _, err := backupUnitSvc.Create(testS3KeyResourceVar)
		assert.Error(t, err)
	})
	t.Run("update_s3key_error", func(t *testing.T) {
		svc := getTestClient(t)
		backupUnitSvc := NewS3KeyService(svc.Get(), ctx)
		_, _, err := backupUnitSvc.Update(testS3KeyResourceVar, testS3KeyResourceVar, S3Key{})
		assert.Error(t, err)
	})
	t.Run("delete_s3key_error", func(t *testing.T) {
		svc := getTestClient(t)
		backupUnitSvc := NewS3KeyService(svc.Get(), ctx)
		_, err := backupUnitSvc.Delete(testS3KeyResourceVar, testS3KeyResourceVar)
		assert.Error(t, err)
	})
}
