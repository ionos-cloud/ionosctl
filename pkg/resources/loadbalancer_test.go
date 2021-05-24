package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testLoadbalancerResourceVar = "test-loadbalancer-resource"
)

func TestNewLoadbalancerService(t *testing.T) {
	ctx := context.Background()
	t.Run("list_loadbalancers_error", func(t *testing.T) {
		svc := getTestClient(t)
		backupUnitSvc := NewLoadbalancerService(svc.Get(), ctx)
		_, _, err := backupUnitSvc.List(testLoadbalancerResourceVar)
		assert.Error(t, err)
	})
	t.Run("get_loadbalancer_error", func(t *testing.T) {
		svc := getTestClient(t)
		backupUnitSvc := NewLoadbalancerService(svc.Get(), ctx)
		_, _, err := backupUnitSvc.Get(testLoadbalancerResourceVar, testLoadbalancerResourceVar)
		assert.Error(t, err)
	})
	t.Run("create_loadbalancer_error", func(t *testing.T) {
		svc := getTestClient(t)
		backupUnitSvc := NewLoadbalancerService(svc.Get(), ctx)
		_, _, err := backupUnitSvc.Create(testLoadbalancerResourceVar, testLoadbalancerResourceVar, false)
		assert.Error(t, err)
	})
	t.Run("update_loadbalancer_error", func(t *testing.T) {
		svc := getTestClient(t)
		backupUnitSvc := NewLoadbalancerService(svc.Get(), ctx)
		_, _, err := backupUnitSvc.Update(testLoadbalancerResourceVar, testLoadbalancerResourceVar, LoadbalancerProperties{})
		assert.Error(t, err)
	})
	t.Run("delete_loadbalancer_error", func(t *testing.T) {
		svc := getTestClient(t)
		backupUnitSvc := NewLoadbalancerService(svc.Get(), ctx)
		_, err := backupUnitSvc.Delete(testLoadbalancerResourceVar, testLoadbalancerResourceVar)
		assert.Error(t, err)
	})
	t.Run("attach_nic_loadbalancer_error", func(t *testing.T) {
		svc := getTestClient(t)
		backupUnitSvc := NewLoadbalancerService(svc.Get(), ctx)
		_, _, err := backupUnitSvc.AttachNic(testLoadbalancerResourceVar, testLoadbalancerResourceVar, testLoadbalancerResourceVar)
		assert.Error(t, err)
	})
	t.Run("list_nics_loadbalancer_error", func(t *testing.T) {
		svc := getTestClient(t)
		backupUnitSvc := NewLoadbalancerService(svc.Get(), ctx)
		_, _, err := backupUnitSvc.ListNics(testLoadbalancerResourceVar, testLoadbalancerResourceVar)
		assert.Error(t, err)
	})
	t.Run("get_nic_loadbalancer_error", func(t *testing.T) {
		svc := getTestClient(t)
		backupUnitSvc := NewLoadbalancerService(svc.Get(), ctx)
		_, _, err := backupUnitSvc.GetNic(testLoadbalancerResourceVar, testLoadbalancerResourceVar, testLoadbalancerResourceVar)
		assert.Error(t, err)
	})
	t.Run("detach_nic_loadbalancer_error", func(t *testing.T) {
		svc := getTestClient(t)
		backupUnitSvc := NewLoadbalancerService(svc.Get(), ctx)
		_, err := backupUnitSvc.DetachNic(testLoadbalancerResourceVar, testLoadbalancerResourceVar, testLoadbalancerResourceVar)
		assert.Error(t, err)
	})
}
