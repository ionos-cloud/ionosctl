package group

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
)

var allResourceCols = []table.Column{
	{Name: "ResourceId", JSONPath: "id", Default: true},
	{Name: "Name", JSONPath: "properties.name", Default: true},
	{Name: "SecAuthProtection", JSONPath: "properties.secAuthProtection", Default: true},
	{Name: "Type", JSONPath: "type", Default: true},
	{Name: "State", JSONPath: "metadata.state", Default: true},
}

func GroupResourceCmd() *core.Command {
	resourceCmd := &core.Command{
		Command: &cobra.Command{
			Use:     "resource",
			Aliases: []string{"res"},
			Short:   "List the resources a Group has been granted access to",
			Long: `List the specific resources (datacenters, snapshots, images, IP blocks, etc.) that have been SHARED with a Group. These shares are what give a Group access to individual existing resources, as opposed to the contract-wide privileges set on the Group itself.

This is the read-only companion to ` + "`ionosctl compute share`" + `, which is where you create, update and delete the shares (with edit / re-share permissions).`,
			TraverseChildren: true,
		},
	}

	resourceCmd.AddCommand(groupResourceListCmd())

	return core.WithConfigOverride(resourceCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}

func groupResourceListCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace:  "group",
		Resource:   "resource",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List the resources shared with a Group",
		LongDesc:   "List every resource that has been shared with the given Group, i.e. the resources its members can access by virtue of membership. To inspect any resource across the whole contract (not just those shared with a Group), use the `ionosctl compute resource` commands; to grant or change a share, use `ionosctl compute share`.\n\nRequired values to run command:\n\n* Group Id",
		Example:    "ionosctl compute group resource list --group-id GROUP_ID",
		PreCmdRun:  PreRunGroupId,
		CmdRun:     RunGroupResourceList,
		InitClient: true,
	})
	cmd.AddColsFlag(allResourceCols)
	cmd.AddUUIDFlag(cloudapiv6.ArgGroupId, "", "", cloudapiv6.GroupId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgGroupId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.GroupsIds(), cobra.ShellCompDirectiveNoFileComp
	})

	return cmd
}
