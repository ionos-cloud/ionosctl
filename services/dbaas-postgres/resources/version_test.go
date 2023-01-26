package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testVersionResourceVar = "test-version-resource"

func TestNewVersionsService(t *testing.T) {
	ctx := context.Background()
	t.Run("list_version_error", func(t *testing.T) {
		svc := getTestClient(t)
		versionSvc := NewVersionsService(svc, ctx)
		_, _, err := versionSvc.List()
		assert.Error(t, err)
	})
	t.Run("get_versions_error", func(t *testing.T) {
		svc := getTestClient(t)
		versionSvc := NewVersionsService(svc, ctx)
		_, _, err := versionSvc.Get(testVersionResourceVar)
		assert.Error(t, err)
	})
}
