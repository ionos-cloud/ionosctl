package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testFlowLogResourceVar = "test-flowlog-resource"
)

func TestNewFlowLogService(t *testing.T) {
	ctx := context.Background()
	t.Run("list_flowlogs_error", func(t *testing.T) {
		svc := getTestClient(t)
		flowlogSvc := NewFlowLogService(svc.Get(), ctx)
		_, _, err := flowlogSvc.List(
			testFlowLogResourceVar,
			testFlowLogResourceVar,
			testFlowLogResourceVar,
			ListQueryParams{},
		)
		assert.Error(t, err)
	})
	t.Run("list_flowlogs_filters_error", func(t *testing.T) {
		svc := getTestClient(t)
		flowlogSvc := NewFlowLogService(svc.Get(), ctx)
		_, _, err := flowlogSvc.List(
			testFlowLogResourceVar,
			testFlowLogResourceVar,
			testFlowLogResourceVar,
			testListQueryParam,
		)
		assert.Error(t, err)
	})
	t.Run("get_flowlog_error", func(t *testing.T) {
		svc := getTestClient(t)
		flowlogSvc := NewFlowLogService(svc.Get(), ctx)
		_, _, err := flowlogSvc.Get(
			testFlowLogResourceVar,
			testFlowLogResourceVar,
			testFlowLogResourceVar,
			testFlowLogResourceVar,
		)
		assert.Error(t, err)
	})
	t.Run("create_flowlog_error", func(t *testing.T) {
		svc := getTestClient(t)
		flowlogSvc := NewFlowLogService(svc.Get(), ctx)
		_, _, err := flowlogSvc.Create(
			testFlowLogResourceVar,
			testFlowLogResourceVar,
			testFlowLogResourceVar,
			FlowLog{},
		)
		assert.Error(t, err)
	})
	t.Run("update_flowlog_error", func(t *testing.T) {
		svc := getTestClient(t)
		flowlogSvc := NewFlowLogService(svc.Get(), ctx)
		_, _, err := flowlogSvc.Update(
			testFlowLogResourceVar,
			testFlowLogResourceVar,
			testFlowLogResourceVar,
			testFlowLogResourceVar,
			FlowLogPut{},
		)
		assert.Error(t, err)
	})
	t.Run("delete_flowlog_error", func(t *testing.T) {
		svc := getTestClient(t)
		flowlogSvc := NewFlowLogService(svc.Get(), ctx)
		_, err := flowlogSvc.Delete(
			testFlowLogResourceVar,
			testFlowLogResourceVar,
			testFlowLogResourceVar,
			testFlowLogResourceVar,
		)
		assert.Error(t, err)
	})
}
