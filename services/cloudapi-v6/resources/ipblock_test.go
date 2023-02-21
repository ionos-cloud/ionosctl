package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testIpBlockResourceVar = "test-ipblock-resource"
)

func TestNewIpBlockService(t *testing.T) {
	ctx := context.Background()
	t.Run("list_ipblocks_error", func(t *testing.T) {
		svc := getTestClient(t)
		ipblockSvc := NewIpBlockService(svc, ctx)
		_, _, err := ipblockSvc.List(ListQueryParams{})
		assert.Error(t, err)
	})
	t.Run("list_ipblocks_filters_error", func(t *testing.T) {
		svc := getTestClient(t)
		ipblockSvc := NewIpBlockService(svc, ctx)
		_, _, err := ipblockSvc.List(testListQueryParam)
		assert.Error(t, err)
	})
	t.Run("get_ipblock_error", func(t *testing.T) {
		svc := getTestClient(t)
		ipblockSvc := NewIpBlockService(svc, ctx)
		_, _, err := ipblockSvc.Get(testIpBlockResourceVar, QueryParams{})
		assert.Error(t, err)
	})
	t.Run("create_ipblock_error", func(t *testing.T) {
		svc := getTestClient(t)
		ipblockSvc := NewIpBlockService(svc, ctx)
		_, _, err := ipblockSvc.Create(
			testIpBlockResourceVar,
			testIpBlockResourceVar,
			int32(1),
			QueryParams{},
		)
		assert.Error(t, err)
	})
	t.Run("update_ipblock_error", func(t *testing.T) {
		svc := getTestClient(t)
		ipblockSvc := NewIpBlockService(svc, ctx)
		_, _, err := ipblockSvc.Update(testIpBlockResourceVar, IpBlockProperties{}, QueryParams{})
		assert.Error(t, err)
	})
	t.Run("delete_ipblock_error", func(t *testing.T) {
		svc := getTestClient(t)
		ipblockSvc := NewIpBlockService(svc, ctx)
		_, err := ipblockSvc.Delete(testIpBlockResourceVar, QueryParams{})
		assert.Error(t, err)
	})
}
