package completer

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/ionos-cloud/ionosctl/v6/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestFiltersLength(t *testing.T) {
	t.Run("datacenters_filters", func(t *testing.T) {
		filters := DataCentersFilters()
		assert.True(t, len(filters) > 0)
	})
	t.Run("backupunits_filters", func(t *testing.T) {
		filters := BackupUnitsFilters()
		assert.True(t, len(filters) > 0)
	})
	t.Run("servers_filters", func(t *testing.T) {
		filters := ServersFilters()
		assert.True(t, len(filters) > 0)
	})
	t.Run("images_filters", func(t *testing.T) {
		filters := ImagesFilters()
		assert.True(t, len(filters) > 0)
	})
	t.Run("volumes_filters", func(t *testing.T) {
		filters := VolumesFilters()
		assert.True(t, len(filters) > 0)
	})
	t.Run("snapshots_filters", func(t *testing.T) {
		filters := SnapshotsFilters()
		assert.True(t, len(filters) > 0)
	})
	t.Run("ipblocks_filters", func(t *testing.T) {
		filters := IpBlocksFilters()
		assert.True(t, len(filters) > 0)
	})
	t.Run("lans_filters", func(t *testing.T) {
		filters := LANsFilters()
		assert.True(t, len(filters) > 0)
	})
	t.Run("nics_filters", func(t *testing.T) {
		filters := NICsFilters()
		assert.True(t, len(filters) > 0)
	})
	t.Run("locations_filters", func(t *testing.T) {
		filters := LocationsFilters()
		assert.True(t, len(filters) > 0)
	})
	t.Run("firewallrules_filters", func(t *testing.T) {
		filters := FirewallRulesFilters()
		assert.True(t, len(filters) > 0)
	})
	t.Run("loadbalancers_filters", func(t *testing.T) {
		filters := LoadBalancersFilters()
		assert.True(t, len(filters) > 0)
	})
	t.Run("requests_filters", func(t *testing.T) {
		filters := RequestsFilters()
		assert.True(t, len(filters) > 0)
	})
	t.Run("users_filters", func(t *testing.T) {
		filters := UsersFilters()
		assert.True(t, len(filters) > 0)
	})
	t.Run("k8sclusters_filters", func(t *testing.T) {
		filters := K8sClustersFilters()
		assert.True(t, len(filters) > 0)
	})
	t.Run("k8snodepools_filters", func(t *testing.T) {
		filters := K8sNodePoolsFilters()
		assert.True(t, len(filters) > 0)
	})
	t.Run("k8snodes_filters", func(t *testing.T) {
		filters := K8sNodesFilters()
		assert.True(t, len(filters) > 0)
	})
	t.Run("flowlogs_filters", func(t *testing.T) {
		filters := FlowLogsFilters()
		assert.True(t, len(filters) > 0)
	})
	t.Run("groups_filters", func(t *testing.T) {
		filters := GroupsFilters()
		assert.True(t, len(filters) > 0)
	})
	t.Run("natgateways_filters", func(t *testing.T) {
		filters := NATGatewaysFilters()
		assert.True(t, len(filters) > 0)
	})
	t.Run("natgatewayrules_filters", func(t *testing.T) {
		filters := NATGatewayRulesFilters()
		assert.True(t, len(filters) > 0)
	})
	t.Run("nlb_filters", func(t *testing.T) {
		filters := NlbsFilters()
		assert.True(t, len(filters) > 0)
	})
	t.Run("nlbrule_filters", func(t *testing.T) {
		filters := NlbRulesFilters()
		assert.True(t, len(filters) > 0)
	})
	t.Run("pcc_filters", func(t *testing.T) {
		filters := PccsFilters()
		assert.True(t, len(filters) > 0)
	})
	t.Run("templates_filters", func(t *testing.T) {
		filters := TemplatesFilters()
		assert.True(t, len(filters) > 0)
	})
}

func TestFiltersUsage(t *testing.T) {
	t.Run("datacenters_filters_usage", func(t *testing.T) {
		filtersUsage := DataCentersFiltersUsage()
		assert.True(t, filtersUsage != "")
	})
	t.Run("backupunits_filters_usage", func(t *testing.T) {
		filtersUsage := BackupUnitsFiltersUsage()
		assert.True(t, filtersUsage != "")
	})
	t.Run("servers_filters_usage", func(t *testing.T) {
		filtersUsage := ServersFiltersUsage()
		assert.True(t, filtersUsage != "")
	})
	t.Run("images_filters_usage", func(t *testing.T) {
		filtersUsage := ImagesFiltersUsage()
		assert.True(t, filtersUsage != "")
	})
	t.Run("volumes_filters_usage", func(t *testing.T) {
		filtersUsage := VolumesFiltersUsage()
		assert.True(t, filtersUsage != "")
	})
	t.Run("snapshots_filters_usage", func(t *testing.T) {
		filtersUsage := SnapshotsFiltersUsage()
		assert.True(t, filtersUsage != "")
	})
	t.Run("ipblocks_filters_usage", func(t *testing.T) {
		filtersUsage := IpBlocksFiltersUsage()
		assert.True(t, filtersUsage != "")
	})
	t.Run("lans_filters_usage", func(t *testing.T) {
		filtersUsage := LANsFiltersUsage()
		assert.True(t, filtersUsage != "")
	})
	t.Run("nics_filters_usage", func(t *testing.T) {
		filtersUsage := NICsFiltersUsage()
		assert.True(t, filtersUsage != "")
	})
	t.Run("locations_filters_usage", func(t *testing.T) {
		filtersUsage := LocationsFiltersUsage()
		assert.True(t, filtersUsage != "")
	})
	t.Run("firewallrules_filters_usage", func(t *testing.T) {
		filtersUsage := FirewallRulesFiltersUsage()
		assert.True(t, filtersUsage != "")
	})
	t.Run("loadbalancers_filters_usage", func(t *testing.T) {
		filtersUsage := LoadbalancersFiltersUsage()
		assert.True(t, filtersUsage != "")
	})
	t.Run("requests_filters_usage", func(t *testing.T) {
		filtersUsage := RequestsFiltersUsage()
		assert.True(t, filtersUsage != "")
	})
	t.Run("users_filters_usage", func(t *testing.T) {
		filtersUsage := UsersFiltersUsage()
		assert.True(t, filtersUsage != "")
	})
	t.Run("k8sclusters_filters_usage", func(t *testing.T) {
		filtersUsage := K8sClustersFiltersUsage()
		assert.True(t, filtersUsage != "")
	})
	t.Run("k8snodepools_filters_usage", func(t *testing.T) {
		filtersUsage := K8sNodePoolsFiltersUsage()
		assert.True(t, filtersUsage != "")
	})
	t.Run("k8snodes_filters_usage", func(t *testing.T) {
		filtersUsage := K8sNodesFiltersUsage()
		assert.True(t, filtersUsage != "")
	})
	t.Run("flowlogs_filters_usage", func(t *testing.T) {
		filtersUsage := FlowLogsFiltersUsage()
		assert.True(t, filtersUsage != "")
	})
	t.Run("groups_filters_usage", func(t *testing.T) {
		filtersUsage := GroupsFiltersUsage()
		assert.True(t, filtersUsage != "")
	})
	t.Run("natgateways_filters_usage", func(t *testing.T) {
		filtersUsage := NATGatewaysFiltersUsage()
		assert.True(t, filtersUsage != "")
	})
	t.Run("natgatewayrules_filters_usage", func(t *testing.T) {
		filtersUsage := NATGatewayRulesFiltersUsage()
		assert.True(t, filtersUsage != "")
	})
	t.Run("nlb_filters_usage", func(t *testing.T) {
		filtersUsage := NlbsFiltersUsage()
		assert.True(t, filtersUsage != "")
	})
	t.Run("nlbrule_filters_usage", func(t *testing.T) {
		filtersUsage := NlbRulesFiltersUsage()
		assert.True(t, filtersUsage != "")
	})
	t.Run("pcc_filters_usage", func(t *testing.T) {
		filtersUsage := PccsFiltersUsage()
		assert.True(t, filtersUsage != "")
	})
	t.Run("templates_filters_usage", func(t *testing.T) {
		filtersUsage := TemplatesFiltersUsage()
		assert.True(t, filtersUsage != "")
	})
}

// This struct is used for testing purposes.
type testStruct struct {
	a               *string
	test            *string
	TestSliceString *[]string
	testSliceInt    *[]int
	TestSliceStruct *[]testStruct
}

func TestGetPropertiesName(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	result := getPropertiesName(testStruct{})
	err := w.Flush()
	assert.NoError(t, err)
	expectedResult := []string{"a", "test", "testSliceString"}
	assert.True(t, utils.StringSlicesEqual(result, expectedResult))
}
