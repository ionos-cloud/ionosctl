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
		Namespace:  "datacenter",
		Resource:   "datacenter",
		Verb:       "get",
		Aliases:    []string{"g"},
		ShortDesc:  "Get a Data Center",
		LongDesc:   "Use this command to retrieve details about a Virtual Data Center by using its ID. You can also retrieve relevant information about the Data Center resources.\n\nRequired values to run command:\n\n* Data Center Id",
		Example:    "ionosctl compute datacenter get --datacenter-id DATACENTER_ID",
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
