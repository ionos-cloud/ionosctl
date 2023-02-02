package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testPrivateCrossConnectResourceVar = "test-pcc-resource"
)

func TestNewPrivateCrossConnectService(t *testing.T) {
	ctx := context.Background()
	t.Run("list_pccs_error", func(t *testing.T) {
		svc := getTestClient(t)
		pccSvc := NewPrivateCrossConnectService(svc, ctx)
		_, _, err := pccSvc.List(ListQueryParams{})
		assert.Error(t, err)
	})
	t.Run("list_pccs_filters_error", func(t *testing.T) {
		svc := getTestClient(t)
		pccSvc := NewPrivateCrossConnectService(svc, ctx)
		_, _, err := pccSvc.List(testListQueryParam)
		assert.Error(t, err)
	})
	t.Run("get_pcc_error", func(t *testing.T) {
		svc := getTestClient(t)
		pccSvc := NewPrivateCrossConnectService(svc, ctx)
		_, _, err := pccSvc.Get(testPrivateCrossConnectResourceVar, QueryParams{})
		assert.Error(t, err)
	})
	t.Run("getpeers_pcc_error", func(t *testing.T) {
		svc := getTestClient(t)
		pccSvc := NewPrivateCrossConnectService(svc, ctx)
		_, _, err := pccSvc.GetPeers(testPrivateCrossConnectResourceVar)
		assert.Error(t, err)
	})
	t.Run("create_pcc_error", func(t *testing.T) {
		svc := getTestClient(t)
		pccSvc := NewPrivateCrossConnectService(svc, ctx)
		_, _, err := pccSvc.Create(PrivateCrossConnect{}, QueryParams{})
		assert.Error(t, err)
	})
	t.Run("update_pcc_error", func(t *testing.T) {
		svc := getTestClient(t)
		pccSvc := NewPrivateCrossConnectService(svc, ctx)
		_, _, err := pccSvc.Update(testPrivateCrossConnectResourceVar, PrivateCrossConnectProperties{}, QueryParams{})
		assert.Error(t, err)
	})
	t.Run("delete_pcc_error", func(t *testing.T) {
		svc := getTestClient(t)
		pccSvc := NewPrivateCrossConnectService(svc, ctx)
		_, err := pccSvc.Delete(testPrivateCrossConnectResourceVar, QueryParams{})
		assert.Error(t, err)
	})
}
