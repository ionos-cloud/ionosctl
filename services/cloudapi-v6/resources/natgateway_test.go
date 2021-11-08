package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testNatGatewayResourceVar = "test-natgateway-resource"
)

func TestNewNatGatewayService(t *testing.T) {
	ctx := context.Background()
	t.Run("list_natgateways_error", func(t *testing.T) {
		svc := getTestClient(t)
		natgatewaySvc := NewNatGatewayService(svc.Get(), ctx)
		_, _, err := natgatewaySvc.List(testNatGatewayResourceVar, ListQueryParams{})
		assert.Error(t, err)
	})
	t.Run("get_natgateway_error", func(t *testing.T) {
		svc := getTestClient(t)
		natgatewaySvc := NewNatGatewayService(svc.Get(), ctx)
		_, _, err := natgatewaySvc.Get(testNatGatewayResourceVar, testNatGatewayResourceVar)
		assert.Error(t, err)
	})
	t.Run("create_natgateway_error", func(t *testing.T) {
		svc := getTestClient(t)
		natgatewaySvc := NewNatGatewayService(svc.Get(), ctx)
		_, _, err := natgatewaySvc.Create(testNatGatewayResourceVar, NatGateway{})
		assert.Error(t, err)
	})
	t.Run("update_natgateway_error", func(t *testing.T) {
		svc := getTestClient(t)
		natgatewaySvc := NewNatGatewayService(svc.Get(), ctx)
		_, _, err := natgatewaySvc.Update(testNatGatewayResourceVar, testNatGatewayResourceVar, NatGatewayProperties{})
		assert.Error(t, err)
	})
	t.Run("delete_natgateway_error", func(t *testing.T) {
		svc := getTestClient(t)
		natgatewaySvc := NewNatGatewayService(svc.Get(), ctx)
		_, err := natgatewaySvc.Delete(testNatGatewayResourceVar, testNatGatewayResourceVar)
		assert.Error(t, err)
	})
	t.Run("listrules_natgateway_error", func(t *testing.T) {
		svc := getTestClient(t)
		natgatewaySvc := NewNatGatewayService(svc.Get(), ctx)
		_, _, err := natgatewaySvc.ListRules(testNatGatewayResourceVar, testNatGatewayResourceVar, ListQueryParams{})
		assert.Error(t, err)
	})
	t.Run("getrule_natgateway_error", func(t *testing.T) {
		svc := getTestClient(t)
		natgatewaySvc := NewNatGatewayService(svc.Get(), ctx)
		_, _, err := natgatewaySvc.GetRule(testNatGatewayResourceVar, testNatGatewayResourceVar, testNatGatewayResourceVar)
		assert.Error(t, err)
	})
	t.Run("createrule_natgateway_error", func(t *testing.T) {
		svc := getTestClient(t)
		natgatewaySvc := NewNatGatewayService(svc.Get(), ctx)
		_, _, err := natgatewaySvc.CreateRule(testNatGatewayResourceVar, testNatGatewayResourceVar, NatGatewayRule{})
		assert.Error(t, err)
	})
	t.Run("updaterule_natgateway_error", func(t *testing.T) {
		svc := getTestClient(t)
		natgatewaySvc := NewNatGatewayService(svc.Get(), ctx)
		_, _, err := natgatewaySvc.UpdateRule(testNatGatewayResourceVar, testNatGatewayResourceVar, testNatGatewayResourceVar, NatGatewayRuleProperties{})
		assert.Error(t, err)
	})
	t.Run("deleterule_natgateway_error", func(t *testing.T) {
		svc := getTestClient(t)
		natgatewaySvc := NewNatGatewayService(svc.Get(), ctx)
		_, err := natgatewaySvc.DeleteRule(testNatGatewayResourceVar, testNatGatewayResourceVar, testNatGatewayResourceVar)
		assert.Error(t, err)
	})
	t.Run("listflowlogs_natgateway_error", func(t *testing.T) {
		svc := getTestClient(t)
		natgatewaySvc := NewNatGatewayService(svc.Get(), ctx)
		_, _, err := natgatewaySvc.ListFlowLogs(testNatGatewayResourceVar, testNatGatewayResourceVar, ListQueryParams{})
		assert.Error(t, err)
	})
	t.Run("getflowlog_natgateway_error", func(t *testing.T) {
		svc := getTestClient(t)
		natgatewaySvc := NewNatGatewayService(svc.Get(), ctx)
		_, _, err := natgatewaySvc.GetFlowLog(testNatGatewayResourceVar, testNatGatewayResourceVar, testNatGatewayResourceVar)
		assert.Error(t, err)
	})
	t.Run("createflowlog_natgateway_error", func(t *testing.T) {
		svc := getTestClient(t)
		natgatewaySvc := NewNatGatewayService(svc.Get(), ctx)
		_, _, err := natgatewaySvc.CreateFlowLog(testNatGatewayResourceVar, testNatGatewayResourceVar, FlowLog{})
		assert.Error(t, err)
	})
	t.Run("updateflowlog_natgateway_error", func(t *testing.T) {
		svc := getTestClient(t)
		natgatewaySvc := NewNatGatewayService(svc.Get(), ctx)
		_, _, err := natgatewaySvc.UpdateFlowLog(testNatGatewayResourceVar, testNatGatewayResourceVar, testNatGatewayResourceVar, &FlowLogProperties{})
		assert.Error(t, err)
	})
	t.Run("deleteflowlog_natgateway_error", func(t *testing.T) {
		svc := getTestClient(t)
		natgatewaySvc := NewNatGatewayService(svc.Get(), ctx)
		_, err := natgatewaySvc.DeleteFlowLog(testNatGatewayResourceVar, testNatGatewayResourceVar, testNatGatewayResourceVar)
		assert.Error(t, err)
	})
}
