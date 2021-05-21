package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testDatacenterResourceVar = "test-datacenter-resource"
)

func TestNewDataCenterService(t *testing.T) {
	ctx := context.Background()
	t.Run("list_datacenters_error", func(t *testing.T) {
		svc := getTestClient(t)
		backupUnitSvc := NewDataCenterService(svc.Get(), ctx)
		_, _, err := backupUnitSvc.List()
		assert.Error(t, err)
	})
	t.Run("get_datacenter_error", func(t *testing.T) {
		svc := getTestClient(t)
		backupUnitSvc := NewDataCenterService(svc.Get(), ctx)
		_, _, err := backupUnitSvc.Get(testDatacenterResourceVar)
		assert.Error(t, err)
	})
	t.Run("create_datacenter_error", func(t *testing.T) {
		svc := getTestClient(t)
		backupUnitSvc := NewDataCenterService(svc.Get(), ctx)
		_, _, err := backupUnitSvc.Create(
			testDatacenterResourceVar,
			testDatacenterResourceVar,
			testDatacenterResourceVar,
		)
		assert.Error(t, err)
	})
	t.Run("update_datacenter_error", func(t *testing.T) {
		svc := getTestClient(t)
		backupUnitSvc := NewDataCenterService(svc.Get(), ctx)
		_, _, err := backupUnitSvc.Update(testDatacenterResourceVar, DatacenterProperties{})
		assert.Error(t, err)
	})
	t.Run("delete_datacenter_error", func(t *testing.T) {
		svc := getTestClient(t)
		backupUnitSvc := NewDataCenterService(svc.Get(), ctx)
		_, err := backupUnitSvc.Delete(testDatacenterResourceVar)
		assert.Error(t, err)
	})
}
