package targetgroup

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/cobra"
)

func TargetGroupDeleteCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "targetgroup",
		Resource:  "targetgroup",
		Verb:      "delete",
		Aliases:   []string{"d"},
		ShortDesc: "Delete a Target Group",
		LongDesc:  "Delete a Target Group. This removes the backend pool definition itself. Deleting a group that is still referenced by an ALB forwarding rule may be rejected or leave that rule pointing at a missing group, so detach it from any FORWARD rules first. Use --all to delete every Target Group in the contract.\n\nRequired values to run command:\n\n* Target Group Id",
		Example: `# Delete one target group
ionosctl compute targetgroup delete --targetgroup-id TARGET_GROUP_ID --force

# Delete every target group in the contract
ionosctl compute targetgroup delete --all --force`,
		PreCmdRun:  PreRunTargetGroupDelete,
		CmdRun:     RunTargetGroupDelete,
		InitClient: true,
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgTargetGroupId, cloudapiv6.ArgIdShort, "", cloudapiv6.TargetGroupId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgTargetGroupId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.TargetGroupIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddBoolFlag(cloudapiv6.ArgAll, cloudapiv6.ArgAllShort, false, "Delete all Target Groups in the contract instead of a single one. Cannot be combined with --targetgroup-id.")

	return cmd
}
