package resources

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const testLogsResourceVar = "test-cluster-logs-resource"

func TestNewLogsService(t *testing.T) {
	ctx := context.Background()
	t.Run("get_logs_error", func(t *testing.T) {
		svc := getTestClient(t)
		infoUnitSvc := NewLogsService(svc.Get(), ctx)
		_, _, err := infoUnitSvc.Get(testLogsResourceVar, LogsQueryParams{})
		assert.Error(t, err)
	})
	t.Run("get_logs_query_error", func(t *testing.T) {
		svc := getTestClient(t)
		infoUnitSvc := NewLogsService(svc.Get(), ctx)
		_, _, err := infoUnitSvc.Get(testLogsResourceVar, LogsQueryParams{
			Direction: "",
			Limit:     10,
			StartTime: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			EndTime:   time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
		})
		assert.Error(t, err)
	})
}
