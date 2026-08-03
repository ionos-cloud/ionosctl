package resource

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/cobra"
)

func ResourceGetCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "resource",
		Resource:  "resource",
		Verb:      "get",
		Aliases:   []string{"g"},
		ShortDesc: "List resources of one type, or get a single resource",
		LongDesc: `List all resources of a given type, or - if you also pass --resource-id - fetch one specific resource of that type. Valid types are: datacenter, snapshot, image, ipblock, pcc, backupunit, k8s.

Required values to run command:

* Type`,
		Example: `# List every IP block on the contract
ionosctl compute resource get --resource-type ipblock

# Fetch one specific datacenter (to grab its ID/type before sharing it)
ionosctl compute resource get --resource-type datacenter --resource-id DATACENTER_ID`,
		PreCmdRun:  PreRunResourceType,
		CmdRun:     RunResourceGet,
		InitClient: true,
	})
	cmd.AddStringFlag(constants.FlagType, "", "", "The type of resources to retrieve. One of: datacenter, snapshot, image, ipblock, pcc, backupunit, k8s", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagType, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"datacenter", "snapshot", "image", "ipblock", "pcc", "backupunit", "k8s"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgResourceId, cloudapiv6.ArgIdShort, "", "Optional: the ID of a single resource (of the given --resource-type) to fetch. Omit to list all resources of that type")
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgResourceId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ResourcesIds(), cobra.ShellCompDirectiveNoFileComp
	})

	return cmd
}
