package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testLocationResourceVar = "test-location-resource"

func TestNewLocationService(t *testing.T) {
	ctx := context.Background()
	t.Run("list_locations_error", func(t *testing.T) {
		svc := getTestClient(t)
		locationSvc := NewLocationService(svc.Get(), ctx)
		_, _, err := locationSvc.List(ListQueryParams{})
		assert.Error(t, err)
	})
	t.Run("list_locations_filters_error", func(t *testing.T) {
		svc := getTestClient(t)
		locationSvc := NewLocationService(svc.Get(), ctx)
		_, _, err := locationSvc.List(testListQueryParam)
		assert.Error(t, err)
	})
	t.Run("get_region_location_error", func(t *testing.T) {
		svc := getTestClient(t)
		locationSvc := NewLocationService(svc.Get(), ctx)
		_, _, err := locationSvc.GetByRegionAndLocationId(testLocationResourceVar, testLocationResourceVar)
		assert.Error(t, err)
	})
	t.Run("get_location_error", func(t *testing.T) {
		svc := getTestClient(t)
		locationSvc := NewLocationService(svc.Get(), ctx)
		_, _, err := locationSvc.GetByRegionId(testLocationResourceVar)
		assert.Error(t, err)
	})
}
