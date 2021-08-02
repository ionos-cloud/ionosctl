package v5

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testFirewallRuleResourceVar = "test-firewallrule-resource"
)

func TestNewFirewallRuleService(t *testing.T) {
	ctx := context.Background()
	t.Run("list_firewallrules_error", func(t *testing.T) {
		svc := getTestClient(t)
		firewallruleSvc := NewFirewallRuleService(svc.Get(), ctx)
		_, _, err := firewallruleSvc.List(
			testFirewallRuleResourceVar,
			testFirewallRuleResourceVar,
			testFirewallRuleResourceVar,
		)
		assert.Error(t, err)
	})
	t.Run("get_firewallrule_error", func(t *testing.T) {
		svc := getTestClient(t)
		firewallruleSvc := NewFirewallRuleService(svc.Get(), ctx)
		_, _, err := firewallruleSvc.Get(
			testFirewallRuleResourceVar,
			testFirewallRuleResourceVar,
			testFirewallRuleResourceVar,
			testFirewallRuleResourceVar,
		)
		assert.Error(t, err)
	})
	t.Run("create_firewallrule_error", func(t *testing.T) {
		svc := getTestClient(t)
		firewallruleSvc := NewFirewallRuleService(svc.Get(), ctx)
		_, _, err := firewallruleSvc.Create(
			testFirewallRuleResourceVar,
			testFirewallRuleResourceVar,
			testFirewallRuleResourceVar,
			FirewallRule{},
		)
		assert.Error(t, err)
	})
	t.Run("update_firewallrule_error", func(t *testing.T) {
		svc := getTestClient(t)
		firewallruleSvc := NewFirewallRuleService(svc.Get(), ctx)
		_, _, err := firewallruleSvc.Update(
			testFirewallRuleResourceVar,
			testFirewallRuleResourceVar,
			testFirewallRuleResourceVar,
			testFirewallRuleResourceVar,
			FirewallRuleProperties{},
		)
		assert.Error(t, err)
	})
	t.Run("delete_firewallrule_error", func(t *testing.T) {
		svc := getTestClient(t)
		firewallruleSvc := NewFirewallRuleService(svc.Get(), ctx)
		_, err := firewallruleSvc.Delete(
			testFirewallRuleResourceVar,
			testFirewallRuleResourceVar,
			testFirewallRuleResourceVar,
			testFirewallRuleResourceVar,
		)
		assert.Error(t, err)
	})
}
