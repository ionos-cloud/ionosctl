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
	t.Run("list_nics_error", func(t *testing.T) {
		svc := getTestClient(t)
		nicSvc := NewNicService(svc.Get(), ctx)
		_, _, err := nicSvc.List(testNicResourceVar, testNicResourceVar, ListQueryParams{})
		assert.Error(t, err)
	})
	t.Run("list_nics_filters_error", func(t *testing.T) {
		svc := getTestClient(t)
		nicSvc := NewNicService(svc.Get(), ctx)
		_, _, err := nicSvc.List(testNicResourceVar, testNicResourceVar, testListQueryParam)
		assert.Error(t, err)
	})
	t.Run("get_nic_error", func(t *testing.T) {
		svc := getTestClient(t)
		nicSvc := NewNicService(svc.Get(), ctx)
		_, _, err := nicSvc.Get(testNicResourceVar, testNicResourceVar, testNicResourceVar)
		assert.Error(t, err)
	})
	t.Run("create_nic_error", func(t *testing.T) {
		svc := getTestClient(t)
		nicSvc := NewNicService(svc.Get(), ctx)
		_, _, err := nicSvc.Create(testNicResourceVar, testNicResourceVar, Nic{})
		assert.Error(t, err)
	})
	t.Run("update_nic_error", func(t *testing.T) {
		svc := getTestClient(t)
		nicSvc := NewNicService(svc.Get(), ctx)
		_, _, err := nicSvc.Update(testNicResourceVar, testNicResourceVar, testNicResourceVar, NicProperties{})
		assert.Error(t, err)
	})
	t.Run("delete_nic_error", func(t *testing.T) {
		svc := getTestClient(t)
		nicSvc := NewNicService(svc.Get(), ctx)
		_, err := nicSvc.Delete(testNicResourceVar, testNicResourceVar, testNicResourceVar)
		assert.Error(t, err)
	})
}
