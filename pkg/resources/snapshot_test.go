package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testSnapshotResourceVar = "test-snapshot-resource"
)

func TestNewSnapshotService(t *testing.T) {
	ctx := context.Background()
	t.Run("list_snapshots_error", func(t *testing.T) {
		svc := getTestClient(t)
		backupUnitSvc := NewSnapshotService(svc.Get(), ctx)
		_, _, err := backupUnitSvc.List()
		assert.Error(t, err)
	})
	t.Run("get_snapshot_error", func(t *testing.T) {
		svc := getTestClient(t)
		backupUnitSvc := NewSnapshotService(svc.Get(), ctx)
		_, _, err := backupUnitSvc.Get(testSnapshotResourceVar)
		assert.Error(t, err)
	})
	t.Run("create_snapshot_error", func(t *testing.T) {
		svc := getTestClient(t)
		backupUnitSvc := NewSnapshotService(svc.Get(), ctx)
		_, _, err := backupUnitSvc.Create(
			testSnapshotResourceVar,
			testSnapshotResourceVar,
			testSnapshotResourceVar,
			testSnapshotResourceVar,
			testSnapshotResourceVar,
			false,
		)
		assert.Error(t, err)
	})
	t.Run("update_snapshot_error", func(t *testing.T) {
		svc := getTestClient(t)
		backupUnitSvc := NewSnapshotService(svc.Get(), ctx)
		_, _, err := backupUnitSvc.Update(testSnapshotResourceVar, SnapshotProperties{})
		assert.Error(t, err)
	})
	t.Run("restore_snapshot_error", func(t *testing.T) {
		svc := getTestClient(t)
		backupUnitSvc := NewSnapshotService(svc.Get(), ctx)
		_, err := backupUnitSvc.Restore(testSnapshotResourceVar, testSnapshotResourceVar, testSnapshotResourceVar)
		assert.Error(t, err)
	})
	t.Run("delete_snapshot_error", func(t *testing.T) {
		svc := getTestClient(t)
		backupUnitSvc := NewSnapshotService(svc.Get(), ctx)
		_, err := backupUnitSvc.Delete(testSnapshotResourceVar)
		assert.Error(t, err)
	})
}
