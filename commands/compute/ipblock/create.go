package ipblock

import (
	"context"

	cloudapiv6cmds "github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/cobra"
)

func IpBlockCreateCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "ipblock",
		Resource:  "ipblock",
		Verb:      "create",
		Aliases:   []string{"c"},
		ShortDesc: "Create/Reserve an IpBlock",
		LongDesc: `Use this command to create/reserve an IpBlock in a specified location that can be used by resources within any Virtual Data Centers provisioned in that same location.
An IpBlock consists of one or more static IP addresses. The name, size of the IpBlock can be set.
You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option.`,
		Example:    "ionosctl ipblock create --name NAME --location LOCATION_ID --size IPBLOCK_SIZE",
		PreCmdRun:  core.NoPreRun,
		CmdRun:     cloudapiv6cmds.RunIpBlockCreate,
		InitClient: true,
	})
	cmd.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "", "Name of the IpBlock. If not set, it will automatically be set")
	cmd.AddStringFlag(cloudapiv6.ArgLocation, cloudapiv6.ArgLocationShort, "de/txl", "Location of the IpBlock. Location de/fra/2 is currently unavailable.")
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgLocation, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.LocationIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddIntFlag(cloudapiv6.ArgSize, "", 2, "Size of the IpBlock")
	cmd.AddBoolFlag(constants.ArgWaitForRequest, constants.ArgWaitForRequestShort, constants.DefaultWait, "Wait for the Request for IpBlock creation to be executed")
	cmd.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, constants.DefaultTimeoutSeconds, "Timeout option for Request for IpBlock creation [seconds]")

	return cmd
}
