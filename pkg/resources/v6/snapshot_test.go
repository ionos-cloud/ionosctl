package v6

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
		snapshotSvc := NewSnapshotService(svc.Get(), ctx)
		_, _, err := snapshotSvc.List()
		assert.Error(t, err)
	})
	t.Run("get_snapshot_error", func(t *testing.T) {
		svc := getTestClient(t)
		snapshotSvc := NewSnapshotService(svc.Get(), ctx)
		_, _, err := snapshotSvc.Get(testSnapshotResourceVar)
		assert.Error(t, err)
	})
	t.Run("create_snapshot_error", func(t *testing.T) {
		svc := getTestClient(t)
		snapshotSvc := NewSnapshotService(svc.Get(), ctx)
		_, _, err := snapshotSvc.Create(
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
		snapshotSvc := NewSnapshotService(svc.Get(), ctx)
		_, _, err := snapshotSvc.Update(testSnapshotResourceVar, SnapshotProperties{})
		assert.Error(t, err)
	})
	t.Run("restore_snapshot_error", func(t *testing.T) {
		svc := getTestClient(t)
		snapshotSvc := NewSnapshotService(svc.Get(), ctx)
		_, err := snapshotSvc.Restore(testSnapshotResourceVar, testSnapshotResourceVar, testSnapshotResourceVar)
		assert.Error(t, err)
	})
	t.Run("delete_snapshot_error", func(t *testing.T) {
		svc := getTestClient(t)
		snapshotSvc := NewSnapshotService(svc.Get(), ctx)
		_, err := snapshotSvc.Delete(testSnapshotResourceVar)
		assert.Error(t, err)
	})
}
