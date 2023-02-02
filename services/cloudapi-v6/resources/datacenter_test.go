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
		datacenterSvc := NewDataCenterService(svc, ctx)
		_, _, err := datacenterSvc.List(ListQueryParams{})
		assert.Error(t, err)
	})
	t.Run("list_datacenters_filters_error", func(t *testing.T) {
		svc := getTestClient(t)
		datacenterSvc := NewDataCenterService(svc, ctx)
		_, _, err := datacenterSvc.List(testListQueryParam)
		assert.Error(t, err)
	})
	t.Run("get_datacenter_error", func(t *testing.T) {
		svc := getTestClient(t)
		datacenterSvc := NewDataCenterService(svc, ctx)
		_, _, err := datacenterSvc.Get(testDatacenterResourceVar, QueryParams{})
		assert.Error(t, err)
	})
	t.Run("create_datacenter_error", func(t *testing.T) {
		svc := getTestClient(t)
		datacenterSvc := NewDataCenterService(svc, ctx)
		_, _, err := datacenterSvc.Create(
			testDatacenterResourceVar,
			testDatacenterResourceVar,
			testDatacenterResourceVar,
			QueryParams{},
		)
		assert.Error(t, err)
	})
	t.Run("update_datacenter_error", func(t *testing.T) {
		svc := getTestClient(t)
		datacenterSvc := NewDataCenterService(svc, ctx)
		_, _, err := datacenterSvc.Update(testDatacenterResourceVar, DatacenterProperties{}, QueryParams{})
		assert.Error(t, err)
	})
	t.Run("delete_datacenter_error", func(t *testing.T) {
		svc := getTestClient(t)
		datacenterSvc := NewDataCenterService(svc, ctx)
		_, err := datacenterSvc.Delete(testDatacenterResourceVar, QueryParams{})
		assert.Error(t, err)
	})
}
