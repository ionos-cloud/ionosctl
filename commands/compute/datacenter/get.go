package datacenter

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/cobra"
)

func DatacenterGetCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "datacenter",
		Resource:  "datacenter",
		Verb:      "get",
		Aliases:   []string{"g"},
		ShortDesc: "Get details of a single Virtual Data Center",
		LongDesc: `Retrieve the full properties and current state of one Virtual Data Center by its ID.

Beyond the name, description and region, the output surfaces read-only, server-computed fields useful for troubleshooting: the CPU family available in the region (CpuFamily), whether two-step verification is required to touch this VDC (SecAuthProtection), the allocated IPv6 CIDR block, the current provisioning State (e.g. AVAILABLE, BUSY), and the resource Version. Use ` + "`--cols`" + ` to widen the table to any of these.

Required values to run command:

* Data Center Id`,
		Example:    "ionosctl compute datacenter get --datacenter-id DATACENTER_ID\nionosctl compute datacenter get --datacenter-id DATACENTER_ID --cols \"DatacenterId,Name,Location,State,CpuFamily,SecAuthProtection\"",
		PreCmdRun:  PreRunDataCenterId,
		CmdRun:     RunDataCenterGet,
		InitClient: true,
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgDataCenterId, cloudapiv6.ArgIdShort, "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})

	return cmd
}
