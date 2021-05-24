package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testLanResourceVar = "test-lan-resource"
)

func TestNewLanService(t *testing.T) {
	ctx := context.Background()
	t.Run("list_lans_error", func(t *testing.T) {
		svc := getTestClient(t)
		backupUnitSvc := NewLanService(svc.Get(), ctx)
		_, _, err := backupUnitSvc.List(testLanResourceVar)
		assert.Error(t, err)
	})
	t.Run("get_lan_error", func(t *testing.T) {
		svc := getTestClient(t)
		backupUnitSvc := NewLanService(svc.Get(), ctx)
		_, _, err := backupUnitSvc.Get(testLanResourceVar, testLanResourceVar)
		assert.Error(t, err)
	})
	t.Run("create_lan_error", func(t *testing.T) {
		svc := getTestClient(t)
		backupUnitSvc := NewLanService(svc.Get(), ctx)
		_, _, err := backupUnitSvc.Create(
			testLanResourceVar,
			LanPost{},
		)
		assert.Error(t, err)
	})
	t.Run("update_lan_error", func(t *testing.T) {
		svc := getTestClient(t)
		backupUnitSvc := NewLanService(svc.Get(), ctx)
		_, _, err := backupUnitSvc.Update(testLanResourceVar, testLanResourceVar, LanProperties{})
		assert.Error(t, err)
	})
	t.Run("delete_lan_error", func(t *testing.T) {
		svc := getTestClient(t)
		backupUnitSvc := NewLanService(svc.Get(), ctx)
		_, err := backupUnitSvc.Delete(testLanResourceVar, testLanResourceVar)
		assert.Error(t, err)
	})
}
