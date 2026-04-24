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
			Use:              "datacenter",
			Aliases:          []string{"d", "dc", "vdc"},
			Args:             cobra.ExactValidArgs(1),
			Short:            "Data Center Operations",
			Long:             "The sub-commands of `ionosctl compute datacenter` allow you to create, list, get, update and delete Data Centers.",
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
