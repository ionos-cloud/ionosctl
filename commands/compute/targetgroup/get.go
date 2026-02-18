package targetgroup

import (
	"context"

	cloudapiv6cmds "github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/cobra"
)

func TargetGroupGetCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace:  "targetgroup",
		Resource:   "targetgroup",
		Verb:       "get",
		Aliases:    []string{"g"},
		ShortDesc:  "Get a Target Group",
		LongDesc:   "Use this command to get information about a specified Target Group.\n\nRequired values to run command:\n\n* Target Group Id",
		Example:    `ionosctl targetgroup get -i TARGET_GROUP_ID`,
		PreCmdRun:  cloudapiv6cmds.PreRunTargetGroupId,
		CmdRun:     cloudapiv6cmds.RunTargetGroupGet,
		InitClient: true,
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgTargetGroupId, cloudapiv6.ArgIdShort, "", cloudapiv6.TargetGroupId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgTargetGroupId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.TargetGroupIds(), cobra.ShellCompDirectiveNoFileComp
	})

	return cmd
}
