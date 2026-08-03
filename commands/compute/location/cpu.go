package location

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
)

var allCpuCols = []table.Column{
	{Name: "CpuFamily", JSONPath: "cpuFamily", Default: true},
	{Name: "MaxCores", JSONPath: "maxCores", Default: true},
	{Name: "MaxRam", JSONPath: "maxRam", Default: true},
	{Name: "Vendor", JSONPath: "vendor", Default: true},
	{Name: "EnabledFeatures", JSONPath: "enabledFeatures", Default: true},
}

func CpuCmd() *core.Command {
	cpuCmd := &core.Command{
		Command: &cobra.Command{
			Use:   "cpu",
			Short: "List the CPU families available in a location",
			Long: `The ` + "`ionosctl compute location cpu`" + ` command lists the CPU architectures offered at a given location.

Each entry describes one CPU family: its name (the value you pass as ` + "`--cpu-family`" + ` when creating a server, e.g. INTEL_SKYLAKE, INTEL_ICELAKE, AMD_OPTERON), the vendor, the maximum cores and maximum RAM you may request per server on that family, and the enabled hardware features. Because families differ between sites, always check the target location before choosing a ` + "`--cpu-family`" + `. This command is read-only.`,
			TraverseChildren: true,
		},
	}
	cpuCmd.AddColsFlag(allCpuCols)

	cpuCmd.AddCommand(CpuListCmd())

	return core.WithConfigOverride(cpuCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}

func CpuListCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "location",
		Resource:  "cpu",
		Verb:      "list",
		Aliases:   []string{"l", "ls"},
		ShortDesc: "List available CPU families for a location",
		LongDesc: `Use this command to list the CPU families available at a location, along with the max cores, max RAM, vendor and enabled features of each. Use the CpuFamily value as ` + "`--cpu-family`" + ` when creating a server.

Required values to run command:

* Location Id, in ` + "`region/location`" + ` form, e.g. ` + "`de/txl`" + ``,
		Example:    `ionosctl compute location cpu list --location-id de/txl`,
		PreCmdRun:  PreRunLocationId,
		CmdRun:     RunLocationCpuList,
		InitClient: true,
	})
	cmd.AddStringFlag(cloudapiv6.ArgLocationId, "", "", "The location ID in `region/location` form, e.g. de/txl", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgLocationId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.LocationIds(), cobra.ShellCompDirectiveNoFileComp
	})

	return cmd
}
