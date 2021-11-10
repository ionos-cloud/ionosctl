package completer

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/ionos-cloud/ionosctl/internal/utils"
	"github.com/ionos-cloud/ionosctl/internal/utils/clierror"
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
		filters := NetworkLoadBalancersFilters()
		assert.True(t, len(filters) > 0)
	})
	t.Run("nlbrule_filters", func(t *testing.T) {
		filters := NetworkLoadBalancerForwardingRulesFilters()
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

// This struct is used for testing purposes.
type testStruct struct {
	a               *string
	test            *string
	TestSliceString *[]string
	testSliceInt    *[]int
	TestSliceStruct *[]testStruct
}

func TestGetPropertiesName(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() { return }
	w := bufio.NewWriter(&b)
	result := getPropertiesName(testStruct{})
	err := w.Flush()
	assert.NoError(t, err)
	expectedResult := []string{"a", "test", "testSliceString"}
	assert.True(t, utils.StringSlicesEqual(result, expectedResult))
}
