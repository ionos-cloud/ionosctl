package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testNetworkLoadBalancerResourceVar = "test-networkloadbalancer-resource"

func TestNewNetworkLoadBalancerService(t *testing.T) {
	ctx := context.Background()
	t.Run("list_networkloadbalancers_error", func(t *testing.T) {
		svc := getTestClient(t)
		networkloadbalancerSvc := NewNetworkLoadBalancerService(svc.Get(), ctx)
		_, _, err := networkloadbalancerSvc.List(testNetworkLoadBalancerResourceVar)
		assert.Error(t, err)
	})
	t.Run("get_networkloadbalancer_error", func(t *testing.T) {
		svc := getTestClient(t)
		networkloadbalancerSvc := NewNetworkLoadBalancerService(svc.Get(), ctx)
		_, _, err := networkloadbalancerSvc.Get(testNetworkLoadBalancerResourceVar, testNetworkLoadBalancerResourceVar)
		assert.Error(t, err)
	})
	t.Run("create_networkloadbalancer_error", func(t *testing.T) {
		svc := getTestClient(t)
		networkloadbalancerSvc := NewNetworkLoadBalancerService(svc.Get(), ctx)
		_, _, err := networkloadbalancerSvc.Create(testNetworkLoadBalancerResourceVar, NetworkLoadBalancer{})
		assert.Error(t, err)
	})
	t.Run("update_networkloadbalancer_error", func(t *testing.T) {
		svc := getTestClient(t)
		networkloadbalancerSvc := NewNetworkLoadBalancerService(svc.Get(), ctx)
		_, _, err := networkloadbalancerSvc.Update(testNetworkLoadBalancerResourceVar, testNetworkLoadBalancerResourceVar, NetworkLoadBalancerProperties{})
		assert.Error(t, err)
	})
	t.Run("delete_networkloadbalancer_error", func(t *testing.T) {
		svc := getTestClient(t)
		networkloadbalancerSvc := NewNetworkLoadBalancerService(svc.Get(), ctx)
		_, err := networkloadbalancerSvc.Delete(testNetworkLoadBalancerResourceVar, testNetworkLoadBalancerResourceVar)
		assert.Error(t, err)
	})
	t.Run("listrules_networkloadbalancer_error", func(t *testing.T) {
		svc := getTestClient(t)
		networkloadbalancerSvc := NewNetworkLoadBalancerService(svc.Get(), ctx)
		_, _, err := networkloadbalancerSvc.ListForwardingRules(testNetworkLoadBalancerResourceVar, testNetworkLoadBalancerResourceVar)
		assert.Error(t, err)
	})
	t.Run("getrule_networkloadbalancer_error", func(t *testing.T) {
		svc := getTestClient(t)
		networkloadbalancerSvc := NewNetworkLoadBalancerService(svc.Get(), ctx)
		_, _, err := networkloadbalancerSvc.GetForwardingRule(testNetworkLoadBalancerResourceVar, testNetworkLoadBalancerResourceVar, testNetworkLoadBalancerResourceVar)
		assert.Error(t, err)
	})
	t.Run("createrule_networkloadbalancer_error", func(t *testing.T) {
		svc := getTestClient(t)
		networkloadbalancerSvc := NewNetworkLoadBalancerService(svc.Get(), ctx)
		_, _, err := networkloadbalancerSvc.CreateForwardingRule(testNetworkLoadBalancerResourceVar, testNetworkLoadBalancerResourceVar, NetworkLoadBalancerForwardingRule{})
		assert.Error(t, err)
	})
	t.Run("updaterule_networkloadbalancer_error", func(t *testing.T) {
		svc := getTestClient(t)
		networkloadbalancerSvc := NewNetworkLoadBalancerService(svc.Get(), ctx)
		_, _, err := networkloadbalancerSvc.UpdateForwardingRule(testNetworkLoadBalancerResourceVar, testNetworkLoadBalancerResourceVar,
			testNetworkLoadBalancerResourceVar, &NetworkLoadBalancerForwardingRuleProperties{})
		assert.Error(t, err)
	})
	t.Run("deleterule_networkloadbalancer_error", func(t *testing.T) {
		svc := getTestClient(t)
		networkloadbalancerSvc := NewNetworkLoadBalancerService(svc.Get(), ctx)
		_, err := networkloadbalancerSvc.DeleteForwardingRule(testNetworkLoadBalancerResourceVar, testNetworkLoadBalancerResourceVar, testNetworkLoadBalancerResourceVar)
		assert.Error(t, err)
	})
	t.Run("listflowlogs_networkloadbalancer_error", func(t *testing.T) {
		svc := getTestClient(t)
		networkloadbalancerSvc := NewNetworkLoadBalancerService(svc.Get(), ctx)
		_, _, err := networkloadbalancerSvc.ListFlowLogs(testNetworkLoadBalancerResourceVar, testNetworkLoadBalancerResourceVar)
		assert.Error(t, err)
	})
	t.Run("getflowlog_networkloadbalancer_error", func(t *testing.T) {
		svc := getTestClient(t)
		networkloadbalancerSvc := NewNetworkLoadBalancerService(svc.Get(), ctx)
		_, _, err := networkloadbalancerSvc.GetFlowLog(testNetworkLoadBalancerResourceVar, testNetworkLoadBalancerResourceVar, testNetworkLoadBalancerResourceVar)
		assert.Error(t, err)
	})
	t.Run("createflowlog_networkloadbalancer_error", func(t *testing.T) {
		svc := getTestClient(t)
		networkloadbalancerSvc := NewNetworkLoadBalancerService(svc.Get(), ctx)
		_, _, err := networkloadbalancerSvc.CreateFlowLog(testNetworkLoadBalancerResourceVar, testNetworkLoadBalancerResourceVar, FlowLog{})
		assert.Error(t, err)
	})
	t.Run("updateflowlog_networkloadbalancer_error", func(t *testing.T) {
		svc := getTestClient(t)
		networkloadbalancerSvc := NewNetworkLoadBalancerService(svc.Get(), ctx)
		_, _, err := networkloadbalancerSvc.UpdateFlowLog(testNetworkLoadBalancerResourceVar, testNetworkLoadBalancerResourceVar,
			testNetworkLoadBalancerResourceVar, &FlowLogProperties{})
		assert.Error(t, err)
	})
	t.Run("deleteflowlog_networkloadbalancer_error", func(t *testing.T) {
		svc := getTestClient(t)
		networkloadbalancerSvc := NewNetworkLoadBalancerService(svc.Get(), ctx)
		_, err := networkloadbalancerSvc.DeleteFlowLog(testNetworkLoadBalancerResourceVar, testNetworkLoadBalancerResourceVar, testNetworkLoadBalancerResourceVar)
		assert.Error(t, err)
	})
}
