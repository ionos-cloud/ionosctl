package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewVersionsService(t *testing.T) {
	ctx := context.Background()
	t.Run("list_versions_error", func(t *testing.T) {
		svc := getTestClient(t)
		versionSvc := NewVersionsService(svc.Get(), ctx)
		_, _, err := versionSvc.List()
		assert.Error(t, err)
	})
}
