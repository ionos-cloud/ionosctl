package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testNicResourceVar = "test-nic-resource"
)

func TestNewNicService(t *testing.T) {
	ctx := context.Background()
	t.Run("list_ nics_error", func(t *testing.T) {
		svc := getTestClient(t)
		backupUnitSvc := NewNicService(svc.Get(), ctx)
		_, _, err := backupUnitSvc.List(testNicResourceVar, testNicResourceVar)
		assert.Error(t, err)
	})
	t.Run("get_ nic_error", func(t *testing.T) {
		svc := getTestClient(t)
		backupUnitSvc := NewNicService(svc.Get(), ctx)
		_, _, err := backupUnitSvc.Get(testNicResourceVar, testNicResourceVar, testNicResourceVar)
		assert.Error(t, err)
	})
	t.Run("create_ nic_error", func(t *testing.T) {
		svc := getTestClient(t)
		backupUnitSvc := NewNicService(svc.Get(), ctx)
		_, _, err := backupUnitSvc.Create(testNicResourceVar, testNicResourceVar, testNicResourceVar,
			[]string{testNicResourceVar}, false, int32(1))
		assert.Error(t, err)
	})
	t.Run("update_ nic_error", func(t *testing.T) {
		svc := getTestClient(t)
		backupUnitSvc := NewNicService(svc.Get(), ctx)
		_, _, err := backupUnitSvc.Update(testNicResourceVar, testNicResourceVar, testNicResourceVar, NicProperties{})
		assert.Error(t, err)
	})
	t.Run("delete_ nic_error", func(t *testing.T) {
		svc := getTestClient(t)
		backupUnitSvc := NewNicService(svc.Get(), ctx)
		_, err := backupUnitSvc.Delete(testNicResourceVar, testNicResourceVar, testNicResourceVar)
		assert.Error(t, err)
	})
}
