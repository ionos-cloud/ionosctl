package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testRestoreResourceVar = "test-restore-resource"

func TestNewRestoresService(t *testing.T) {
	ctx := context.Background()
	t.Run("restore_cluster_error", func(t *testing.T) {
		svc := getTestClient(t)
		infoUnitSvc := NewRestoresService(svc, ctx)
		_, err := infoUnitSvc.Restore(testRestoreResourceVar, CreateRestoreRequest{})
		assert.Error(t, err)
	})
}
