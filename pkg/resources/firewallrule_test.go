package resources

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
		backupUnitSvc := NewFirewallRuleService(svc.Get(), ctx)
		_, _, err := backupUnitSvc.List(
			testFirewallRuleResourceVar,
			testFirewallRuleResourceVar,
			testFirewallRuleResourceVar,
		)
		assert.Error(t, err)
	})
	t.Run("get_firewallrule_error", func(t *testing.T) {
		svc := getTestClient(t)
		backupUnitSvc := NewFirewallRuleService(svc.Get(), ctx)
		_, _, err := backupUnitSvc.Get(
			testFirewallRuleResourceVar,
			testFirewallRuleResourceVar,
			testFirewallRuleResourceVar,
			testFirewallRuleResourceVar,
		)
		assert.Error(t, err)
	})
	t.Run("create_firewallrule_error", func(t *testing.T) {
		svc := getTestClient(t)
		backupUnitSvc := NewFirewallRuleService(svc.Get(), ctx)
		_, _, err := backupUnitSvc.Create(
			testFirewallRuleResourceVar,
			testFirewallRuleResourceVar,
			testFirewallRuleResourceVar,
			FirewallRule{},
		)
		assert.Error(t, err)
	})
	t.Run("update_firewallrule_error", func(t *testing.T) {
		svc := getTestClient(t)
		backupUnitSvc := NewFirewallRuleService(svc.Get(), ctx)
		_, _, err := backupUnitSvc.Update(
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
		backupUnitSvc := NewFirewallRuleService(svc.Get(), ctx)
		_, err := backupUnitSvc.Delete(
			testFirewallRuleResourceVar,
			testFirewallRuleResourceVar,
			testFirewallRuleResourceVar,
			testFirewallRuleResourceVar,
		)
		assert.Error(t, err)
	})
}
