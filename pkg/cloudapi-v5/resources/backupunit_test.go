package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testBackupUnitResourceVar = "test-backupunit-resource"
)

func TestNewBackupUnitService(t *testing.T) {
	ctx := context.Background()
	t.Run("list_backupunits_error", func(t *testing.T) {
		svc := getTestClient(t)
		backupUnitSvc := NewBackupUnitService(svc.Get(), ctx)
		_, _, err := backupUnitSvc.List()
		assert.Error(t, err)
	})
	t.Run("get_backupunit_error", func(t *testing.T) {
		svc := getTestClient(t)
		backupUnitSvc := NewBackupUnitService(svc.Get(), ctx)
		_, _, err := backupUnitSvc.Get(testBackupUnitResourceVar)
		assert.Error(t, err)
	})
	t.Run("getssourl_backupunit_error", func(t *testing.T) {
		svc := getTestClient(t)
		backupUnitSvc := NewBackupUnitService(svc.Get(), ctx)
		_, _, err := backupUnitSvc.GetSsoUrl(testBackupUnitResourceVar)
		assert.Error(t, err)
	})
	t.Run("create_backupunit_error", func(t *testing.T) {
		svc := getTestClient(t)
		backupUnitSvc := NewBackupUnitService(svc.Get(), ctx)
		_, _, err := backupUnitSvc.Create(BackupUnit{})
		assert.Error(t, err)
	})
	t.Run("update_backupunit_error", func(t *testing.T) {
		svc := getTestClient(t)
		backupUnitSvc := NewBackupUnitService(svc.Get(), ctx)
		_, _, err := backupUnitSvc.Update(testBackupUnitResourceVar, BackupUnitProperties{})
		assert.Error(t, err)
	})
	t.Run("delete_backupunit_error", func(t *testing.T) {
		svc := getTestClient(t)
		backupUnitSvc := NewBackupUnitService(svc.Get(), ctx)
		_, err := backupUnitSvc.Delete(testBackupUnitResourceVar)
		assert.Error(t, err)
	})
}
