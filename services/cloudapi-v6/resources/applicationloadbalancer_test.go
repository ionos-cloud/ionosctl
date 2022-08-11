package resources

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testApplicationLoadBalancerResourceVar = "test-applicationloadbalancer-resource"

func TestNewApplicationLoadBalancerService(t *testing.T) {
	ctx := context.Background()
	t.Run("list_applicationloadbalancers_error", func(t *testing.T) {
		svc := getTestClient(t)
		applicationloadbalancerSvc := NewApplicationLoadBalancerService(svc.Get(), ctx)
		_, _, err := applicationloadbalancerSvc.List(testApplicationLoadBalancerResourceVar, ListQueryParams{})
		assert.Error(t, err)
	})
	t.Run("list_applicationloadbalancers_filters_error", func(t *testing.T) {
		svc := getTestClient(t)
		applicationloadbalancerSvc := NewApplicationLoadBalancerService(svc.Get(), ctx)
		_, _, err := applicationloadbalancerSvc.List(testApplicationLoadBalancerResourceVar, testListQueryParam)
		assert.Error(t, err)
	})
	t.Run("get_applicationloadbalancer_error", func(t *testing.T) {
		svc := getTestClient(t)
		applicationloadbalancerSvc := NewApplicationLoadBalancerService(svc.Get(), ctx)
		_, _, err := applicationloadbalancerSvc.Get(testApplicationLoadBalancerResourceVar, testApplicationLoadBalancerResourceVar, QueryParams{})
		assert.Error(t, err)
	})
	t.Run("create_applicationloadbalancer_error", func(t *testing.T) {
		svc := getTestClient(t)
		applicationloadbalancerSvc := NewApplicationLoadBalancerService(svc.Get(), ctx)
		_, _, err := applicationloadbalancerSvc.Create(testApplicationLoadBalancerResourceVar, ApplicationLoadBalancer{}, QueryParams{})
		assert.Error(t, err)
	})
	t.Run("update_applicationloadbalancer_error", func(t *testing.T) {
		svc := getTestClient(t)
		applicationloadbalancerSvc := NewApplicationLoadBalancerService(svc.Get(), ctx)
		_, _, err := applicationloadbalancerSvc.Update(testApplicationLoadBalancerResourceVar, testApplicationLoadBalancerResourceVar, ApplicationLoadBalancerProperties{}, QueryParams{})
		assert.Error(t, err)
	})
	t.Run("delete_applicationloadbalancer_error", func(t *testing.T) {
		svc := getTestClient(t)
		applicationloadbalancerSvc := NewApplicationLoadBalancerService(svc.Get(), ctx)
		_, err := applicationloadbalancerSvc.Delete(testApplicationLoadBalancerResourceVar, testApplicationLoadBalancerResourceVar, QueryParams{})
		assert.Error(t, err)
	})
	t.Run("listrules_applicationloadbalancer_error", func(t *testing.T) {
		svc := getTestClient(t)
		applicationloadbalancerSvc := NewApplicationLoadBalancerService(svc.Get(), ctx)
		_, _, err := applicationloadbalancerSvc.ListForwardingRules(testApplicationLoadBalancerResourceVar, testApplicationLoadBalancerResourceVar, ListQueryParams{})
		assert.Error(t, err)
	})
	t.Run("listrules_applicationloadbalancer_filters_error", func(t *testing.T) {
		svc := getTestClient(t)
		applicationloadbalancerSvc := NewApplicationLoadBalancerService(svc.Get(), ctx)
		_, _, err := applicationloadbalancerSvc.ListForwardingRules(testApplicationLoadBalancerResourceVar, testApplicationLoadBalancerResourceVar, testListQueryParam)
		assert.Error(t, err)
	})
	t.Run("getrule_applicationloadbalancer_error", func(t *testing.T) {
		svc := getTestClient(t)
		applicationloadbalancerSvc := NewApplicationLoadBalancerService(svc.Get(), ctx)
		_, _, err := applicationloadbalancerSvc.GetForwardingRule(testApplicationLoadBalancerResourceVar, testApplicationLoadBalancerResourceVar, testApplicationLoadBalancerResourceVar, QueryParams{})
		assert.Error(t, err)
	})
	t.Run("createrule_applicationloadbalancer_error", func(t *testing.T) {
		svc := getTestClient(t)
		applicationloadbalancerSvc := NewApplicationLoadBalancerService(svc.Get(), ctx)
		_, _, err := applicationloadbalancerSvc.CreateForwardingRule(testApplicationLoadBalancerResourceVar, testApplicationLoadBalancerResourceVar,
			ApplicationLoadBalancerForwardingRule{}, QueryParams{})
		assert.Error(t, err)
	})
	t.Run("updaterule_applicationloadbalancer_error", func(t *testing.T) {
		svc := getTestClient(t)
		applicationloadbalancerSvc := NewApplicationLoadBalancerService(svc.Get(), ctx)
		_, _, err := applicationloadbalancerSvc.UpdateForwardingRule(testApplicationLoadBalancerResourceVar, testApplicationLoadBalancerResourceVar,
			testApplicationLoadBalancerResourceVar, &ApplicationLoadBalancerForwardingRuleProperties{}, QueryParams{})
		assert.Error(t, err)
	})
	t.Run("deleterule_applicationloadbalancer_error", func(t *testing.T) {
		svc := getTestClient(t)
		applicationloadbalancerSvc := NewApplicationLoadBalancerService(svc.Get(), ctx)
		_, err := applicationloadbalancerSvc.DeleteForwardingRule(testApplicationLoadBalancerResourceVar, testApplicationLoadBalancerResourceVar, testApplicationLoadBalancerResourceVar, QueryParams{})
		assert.Error(t, err)
	})
	t.Run("listflowlogs_applicationloadbalancer_error", func(t *testing.T) {
		svc := getTestClient(t)
		applicationloadbalancerSvc := NewApplicationLoadBalancerService(svc.Get(), ctx)
		_, _, err := applicationloadbalancerSvc.ListFlowLogs(testApplicationLoadBalancerResourceVar, testApplicationLoadBalancerResourceVar, ListQueryParams{})
		assert.Error(t, err)
	})
	t.Run("listflowlogs_applicationloadbalancer_filters_error", func(t *testing.T) {
		svc := getTestClient(t)
		applicationloadbalancerSvc := NewApplicationLoadBalancerService(svc.Get(), ctx)
		_, _, err := applicationloadbalancerSvc.ListFlowLogs(testApplicationLoadBalancerResourceVar, testApplicationLoadBalancerResourceVar, testListQueryParam)
		assert.Error(t, err)
	})
	t.Run("getflowlog_applicationloadbalancer_error", func(t *testing.T) {
		svc := getTestClient(t)
		applicationloadbalancerSvc := NewApplicationLoadBalancerService(svc.Get(), ctx)
		_, _, err := applicationloadbalancerSvc.GetFlowLog(testApplicationLoadBalancerResourceVar, testApplicationLoadBalancerResourceVar, testApplicationLoadBalancerResourceVar, QueryParams{})
		assert.Error(t, err)
	})
	t.Run("createflowlog_applicationloadbalancer_error", func(t *testing.T) {
		svc := getTestClient(t)
		applicationloadbalancerSvc := NewApplicationLoadBalancerService(svc.Get(), ctx)
		_, _, err := applicationloadbalancerSvc.CreateFlowLog(testApplicationLoadBalancerResourceVar, testApplicationLoadBalancerResourceVar, FlowLog{}, QueryParams{})
		assert.Error(t, err)
	})
	t.Run("updateflowlog_applicationloadbalancer_error", func(t *testing.T) {
		svc := getTestClient(t)
		applicationloadbalancerSvc := NewApplicationLoadBalancerService(svc.Get(), ctx)
		_, _, err := applicationloadbalancerSvc.UpdateFlowLog(testApplicationLoadBalancerResourceVar, testApplicationLoadBalancerResourceVar,
			testApplicationLoadBalancerResourceVar, &FlowLogProperties{}, QueryParams{})
		assert.Error(t, err)
	})
	t.Run("deleteflowlog_applicationloadbalancer_error", func(t *testing.T) {
		svc := getTestClient(t)
		applicationloadbalancerSvc := NewApplicationLoadBalancerService(svc.Get(), ctx)
		_, err := applicationloadbalancerSvc.DeleteFlowLog(testApplicationLoadBalancerResourceVar, testApplicationLoadBalancerResourceVar, testApplicationLoadBalancerResourceVar, QueryParams{})
		assert.Error(t, err)
	})
}
