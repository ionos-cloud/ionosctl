package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testBackupResourceVar = "test-backup-resource"

func TestNewBackupUnitService(t *testing.T) {
	ctx := context.Background()
	t.Run("list_backup_error", func(t *testing.T) {
		svc := getTestClient(t)
		backupSvc := NewBackupsService(svc.Get(), ctx)
		_, _, err := backupSvc.List()
		assert.Error(t, err)
	})
	t.Run("get_backups_error", func(t *testing.T) {
		svc := getTestClient(t)
		backupSvc := NewBackupsService(svc.Get(), ctx)
		_, _, err := backupSvc.Get(testBackupResourceVar)
		assert.Error(t, err)
	})
	t.Run("listbackups_cluster_backup_error", func(t *testing.T) {
		svc := getTestClient(t)
		backupSvc := NewBackupsService(svc.Get(), ctx)
		_, _, err := backupSvc.ListBackups(testBackupResourceVar)
		assert.Error(t, err)
	})
}
