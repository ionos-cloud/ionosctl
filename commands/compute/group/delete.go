package group

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/cobra"
)

func GroupDeleteCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "group",
		Resource:  "group",
		Verb:      "delete",
		Aliases:   []string{"d"},
		ShortDesc: "Delete a Group",
		LongDesc: `Delete a single Group. This removes the Group and the privileges it granted; its members lose those privileges (unless another Group still grants them) and lose access to the Group's shared resources.

Deleting a Group does NOT delete the Users in it, nor the resources shared with it - the resources simply become inaccessible to former members, except to a Contract Owner, Admin (administrator user), or Resource Owner.

Required values to run command:

* Group Id`,
		Example:    "ionosctl compute group delete --group-id GROUP_ID",
		PreCmdRun:  PreRunGroupDelete,
		CmdRun:     RunGroupDelete,
		InitClient: true,
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgGroupId, cloudapiv6.ArgIdShort, "", cloudapiv6.GroupId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgGroupId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.GroupsIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddBoolFlag(cloudapiv6.ArgAll, cloudapiv6.ArgAllShort, false, "Delete every Group on the contract. Use with caution: this strips group-granted privileges and shared-resource access from all members")

	return cmd
}
