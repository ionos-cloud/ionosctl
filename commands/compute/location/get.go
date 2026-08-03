package location

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/cobra"
)

func LocationGetCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "location",
		Resource:  "location",
		Verb:      "get",
		Aliases:   []string{"g"},
		ShortDesc: "Get details (features, images, CPU families) of one location",
		LongDesc: `Use this command to get the full details of a single location: the features it supports, the image aliases available there, and its CPU architectures.

Required values to run command:

* Location Id, in ` + "`region/location`" + ` form, e.g. ` + "`de/txl`" + ` (list them with ` + "`ionosctl compute location list`" + `)`,
		Example:    `ionosctl compute location get --location-id de/txl`,
		PreCmdRun:  PreRunLocationId,
		CmdRun:     RunLocationGet,
		InitClient: true,
	})
	cmd.AddStringFlag(cloudapiv6.ArgLocationId, cloudapiv6.ArgIdShort, "", "The location ID in `region/location` form, e.g. de/txl", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgLocationId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.LocationIds(), cobra.ShellCompDirectiveNoFileComp
	})

	return cmd
}
