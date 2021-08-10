package v6

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
		pccSvc := NewPrivateCrossConnectService(svc.Get(), ctx)
		_, _, err := pccSvc.List()
		assert.Error(t, err)
	})
	t.Run("get_pcc_error", func(t *testing.T) {
		svc := getTestClient(t)
		pccSvc := NewPrivateCrossConnectService(svc.Get(), ctx)
		_, _, err := pccSvc.Get(testPrivateCrossConnectResourceVar)
		assert.Error(t, err)
	})
	t.Run("getpeers_pcc_error", func(t *testing.T) {
		svc := getTestClient(t)
		pccSvc := NewPrivateCrossConnectService(svc.Get(), ctx)
		_, _, err := pccSvc.GetPeers(testPrivateCrossConnectResourceVar)
		assert.Error(t, err)
	})
	t.Run("create_pcc_error", func(t *testing.T) {
		svc := getTestClient(t)
		pccSvc := NewPrivateCrossConnectService(svc.Get(), ctx)
		_, _, err := pccSvc.Create(PrivateCrossConnect{})
		assert.Error(t, err)
	})
	t.Run("update_pcc_error", func(t *testing.T) {
		svc := getTestClient(t)
		pccSvc := NewPrivateCrossConnectService(svc.Get(), ctx)
		_, _, err := pccSvc.Update(testPrivateCrossConnectResourceVar, PrivateCrossConnectProperties{})
		assert.Error(t, err)
	})
	t.Run("delete_pcc_error", func(t *testing.T) {
		svc := getTestClient(t)
		pccSvc := NewPrivateCrossConnectService(svc.Get(), ctx)
		_, err := pccSvc.Delete(testPrivateCrossConnectResourceVar)
		assert.Error(t, err)
	})
}
