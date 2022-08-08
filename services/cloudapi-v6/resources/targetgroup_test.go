package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testTargetGroupResourceVar = "test-targetgroup-resource"

func TestNewTargetGroupService(t *testing.T) {
	ctx := context.Background()
	t.Run("list_targetgroups_error", func(t *testing.T) {
		svc := getTestClient(t)
		targetgroupSvc := NewTargetGroupService(svc.Get(), ctx)
		_, _, err := targetgroupSvc.List(ListQueryParams{})
		assert.Error(t, err)
	})
	t.Run("list_targetgroups_filters_error", func(t *testing.T) {
		svc := getTestClient(t)
		targetgroupSvc := NewTargetGroupService(svc.Get(), ctx)
		_, _, err := targetgroupSvc.List(testListQueryParam)
		assert.Error(t, err)
	})
	t.Run("get_targetgroup_error", func(t *testing.T) {
		svc := getTestClient(t)
		targetgroupSvc := NewTargetGroupService(svc.Get(), ctx)
		_, _, err := targetgroupSvc.Get(testTargetGroupResourceVar, QueryParams{})
		assert.Error(t, err)
	})
	t.Run("create_targetgroup_error", func(t *testing.T) {
		svc := getTestClient(t)
		targetgroupSvc := NewTargetGroupService(svc.Get(), ctx)
		_, _, err := targetgroupSvc.Create(TargetGroup{}, QueryParams{})
		assert.Error(t, err)
	})
	t.Run("update_targetgroup_error", func(t *testing.T) {
		svc := getTestClient(t)
		targetgroupSvc := NewTargetGroupService(svc.Get(), ctx)
		_, _, err := targetgroupSvc.Update(testTargetGroupResourceVar, &TargetGroupProperties{}, QueryParams{})
		assert.Error(t, err)
	})
	t.Run("delete_targetgroup_error", func(t *testing.T) {
		svc := getTestClient(t)
		targetgroupSvc := NewTargetGroupService(svc.Get(), ctx)
		_, err := targetgroupSvc.Delete(testTargetGroupResourceVar, QueryParams{})
		assert.Error(t, err)
	})
}
