package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testBackupUnitResourceVar = "test-backupunit-resource"
)

var (
	testListQueryParam = ListQueryParams{
		Filters: &map[string][]string{
			testQueryParamVar: {testQueryParamVar},
		},
		OrderBy:    &testQueryParamVar,
		MaxResults: &testMaxResultsVar,
	}
	testQueryParamVar = "test-filter"
	testMaxResultsVar = int32(2)
)

func TestNewBackupUnitService(t *testing.T) {
	ctx := context.Background()
	t.Run("list_backupunits_error", func(t *testing.T) {
		svc := getTestClient(t)
		backupUnitSvc := NewBackupUnitService(svc, ctx)
		_, _, err := backupUnitSvc.List(ListQueryParams{})
		assert.Error(t, err)
	})
	t.Run("list_backupunits_filters_error", func(t *testing.T) {
		svc := getTestClient(t)
		backupUnitSvc := NewBackupUnitService(svc, ctx)
		_, _, err := backupUnitSvc.List(testListQueryParam)
		assert.Error(t, err)
	})
	t.Run("get_backupunit_error", func(t *testing.T) {
		svc := getTestClient(t)
		backupUnitSvc := NewBackupUnitService(svc, ctx)
		_, _, err := backupUnitSvc.Get(testBackupUnitResourceVar, QueryParams{})
		assert.Error(t, err)
	})
	t.Run("getssourl_backupunit_error", func(t *testing.T) {
		svc := getTestClient(t)
		backupUnitSvc := NewBackupUnitService(svc, ctx)
		_, _, err := backupUnitSvc.GetSsoUrl(testBackupUnitResourceVar, QueryParams{})
		assert.Error(t, err)
	})
	t.Run("create_backupunit_error", func(t *testing.T) {
		svc := getTestClient(t)
		backupUnitSvc := NewBackupUnitService(svc, ctx)
		_, _, err := backupUnitSvc.Create(BackupUnit{}, QueryParams{})
		assert.Error(t, err)
	})
	t.Run("update_backupunit_error", func(t *testing.T) {
		svc := getTestClient(t)
		backupUnitSvc := NewBackupUnitService(svc, ctx)
		_, _, err := backupUnitSvc.Update(testBackupUnitResourceVar, BackupUnitProperties{}, QueryParams{})
		assert.Error(t, err)
	})
	t.Run("delete_backupunit_error", func(t *testing.T) {
		svc := getTestClient(t)
		backupUnitSvc := NewBackupUnitService(svc, ctx)
		_, err := backupUnitSvc.Delete(testBackupUnitResourceVar, QueryParams{})
		assert.Error(t, err)
	})
}
