package datacenter

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
)

var allDatacenterCols = []table.Column{
	{Name: "DatacenterId", JSONPath: "id", Default: true},
	{Name: "Name", JSONPath: "properties.name", Default: true},
	{Name: "Location", JSONPath: "properties.location", Default: true},
	{Name: "CpuFamily", JSONPath: "properties.cpuArchitecture.*.cpuFamily", Default: true},
	{Name: "IPv6CidrBlock", JSONPath: "properties.ipv6CidrBlock", Default: true},
	{Name: "State", JSONPath: "metadata.state", Default: true},
	{Name: "Description", JSONPath: "properties.description"},
	{Name: "Version", JSONPath: "properties.version"},
	{Name: "Features", JSONPath: "properties.features"},
	{Name: "SecAuthProtection", JSONPath: "properties.secAuthProtection"},
}

func DatacenterCmd() *core.Command {
	datacenterCmd := &core.Command{
		Command: &cobra.Command{
			Use:     "datacenter",
			Aliases: []string{"d", "dc", "vdc"},
			Args:    cobra.ExactValidArgs(1),
			Short:   "Create, list, and manage Virtual Data Centers (the top-level container for your compute resources)",
			Long: `Manage Virtual Data Centers (VDCs), the top-level container of IONOS Compute Engine.

A Virtual Data Center is a logically isolated, private environment pinned to a single region (its ` + "`location`" + `, e.g. ` + "`de/txl`" + `). Every compute resource you create - servers, volumes, LANs, NICs, firewalls and virtual routers - lives inside a VDC and inherits its region. Resources in the same VDC can be spread across the availability zones of that region for redundancy, but they cannot span two regions: to run workloads in another region you create a second VDC there. VDCs are network-isolated from one another, so a private LAN in one VDC is not reachable from another.

The region is chosen at creation and is IMMUTABLE - it cannot be changed later (only ` + "`name`" + ` and ` + "`description`" + ` are editable via update). Deleting a VDC permanently destroys everything inside it, so treat delete as a highly destructive operation.

The sub-commands below let you create, list, get, update and delete Virtual Data Centers.`,
			TraverseChildren: true,
		},
	}
	datacenterCmd.AddColsFlag(allDatacenterCols)

	datacenterCmd.AddCommand(DatacenterListCmd())
	datacenterCmd.AddCommand(DatacenterGetCmd())
	datacenterCmd.AddCommand(DatacenterCreateCmd())
	datacenterCmd.AddCommand(DatacenterUpdateCmd())
	datacenterCmd.AddCommand(DatacenterDeleteCmd())

	return core.WithConfigOverride(datacenterCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}
