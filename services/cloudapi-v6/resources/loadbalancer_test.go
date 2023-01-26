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
		loadbalancerSvc := NewLoadbalancerService(svc, ctx)
		_, _, err := loadbalancerSvc.List(testLoadbalancerResourceVar, ListQueryParams{})
		assert.Error(t, err)
	})
	t.Run("list_loadbalancers_filters_error", func(t *testing.T) {
		svc := getTestClient(t)
		loadbalancerSvc := NewLoadbalancerService(svc, ctx)
		_, _, err := loadbalancerSvc.List(testLoadbalancerResourceVar, ListQueryParams{})
		assert.Error(t, err)
	})
	t.Run("get_loadbalancer_error", func(t *testing.T) {
		svc := getTestClient(t)
		loadbalancerSvc := NewLoadbalancerService(svc, ctx)
		_, _, err := loadbalancerSvc.Get(testLoadbalancerResourceVar, testLoadbalancerResourceVar, QueryParams{})
		assert.Error(t, err)
	})
	t.Run("create_loadbalancer_error", func(t *testing.T) {
		svc := getTestClient(t)
		loadbalancerSvc := NewLoadbalancerService(svc, ctx)
		_, _, err := loadbalancerSvc.Create(testLoadbalancerResourceVar, testLoadbalancerResourceVar, false, QueryParams{})
		assert.Error(t, err)
	})
	t.Run("update_loadbalancer_error", func(t *testing.T) {
		svc := getTestClient(t)
		loadbalancerSvc := NewLoadbalancerService(svc, ctx)
		_, _, err := loadbalancerSvc.Update(testLoadbalancerResourceVar, testLoadbalancerResourceVar, LoadbalancerProperties{}, QueryParams{})
		assert.Error(t, err)
	})
	t.Run("delete_loadbalancer_error", func(t *testing.T) {
		svc := getTestClient(t)
		loadbalancerSvc := NewLoadbalancerService(svc, ctx)
		_, err := loadbalancerSvc.Delete(testLoadbalancerResourceVar, testLoadbalancerResourceVar, QueryParams{})
		assert.Error(t, err)
	})
	t.Run("attach_nic_loadbalancer_error", func(t *testing.T) {
		svc := getTestClient(t)
		loadbalancerSvc := NewLoadbalancerService(svc, ctx)
		_, _, err := loadbalancerSvc.AttachNic(testLoadbalancerResourceVar, testLoadbalancerResourceVar, testLoadbalancerResourceVar, QueryParams{})
		assert.Error(t, err)
	})
	t.Run("list_nics_loadbalancer_error", func(t *testing.T) {
		svc := getTestClient(t)
		loadbalancerSvc := NewLoadbalancerService(svc, ctx)
		_, _, err := loadbalancerSvc.ListNics(testLoadbalancerResourceVar, testLoadbalancerResourceVar, ListQueryParams{})
		assert.Error(t, err)
	})
	t.Run("list_nics_filters_loadbalancer_error", func(t *testing.T) {
		svc := getTestClient(t)
		loadbalancerSvc := NewLoadbalancerService(svc, ctx)
		_, _, err := loadbalancerSvc.ListNics(testLoadbalancerResourceVar, testLoadbalancerResourceVar, ListQueryParams{})
		assert.Error(t, err)
	})
	t.Run("get_nic_loadbalancer_error", func(t *testing.T) {
		svc := getTestClient(t)
		loadbalancerSvc := NewLoadbalancerService(svc, ctx)
		_, _, err := loadbalancerSvc.GetNic(testLoadbalancerResourceVar, testLoadbalancerResourceVar, testLoadbalancerResourceVar, QueryParams{})
		assert.Error(t, err)
	})
	t.Run("detach_nic_loadbalancer_error", func(t *testing.T) {
		svc := getTestClient(t)
		loadbalancerSvc := NewLoadbalancerService(svc, ctx)
		_, err := loadbalancerSvc.DetachNic(testLoadbalancerResourceVar, testLoadbalancerResourceVar, testLoadbalancerResourceVar, QueryParams{})
		assert.Error(t, err)
	})
}
