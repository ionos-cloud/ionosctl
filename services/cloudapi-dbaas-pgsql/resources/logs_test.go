package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testLogsResourceVar = "test-cluster-logs-resource"

func TestNewLogsService(t *testing.T) {
	ctx := context.Background()
	t.Run("get_logs_error", func(t *testing.T) {
		svc := getTestClient(t)
		infoUnitSvc := NewLogsService(svc.Get(), ctx)
		_, _, err := infoUnitSvc.Get(testLogsResourceVar, "", "", 0)
		assert.Error(t, err)
	})
}
