package targetgroup

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
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
		LongDesc:   "Get the full configuration of a Target Group: its --algorithm, --protocol, connection and HTTP health-check settings, and its list of targets (backend servers). Use `--cols` to reveal health-check columns not shown by default.\n\nRequired values to run command:\n\n* Target Group Id",
		Example:    `ionosctl compute targetgroup get -i TARGET_GROUP_ID`,
		PreCmdRun:  PreRunTargetGroupId,
		CmdRun:     RunTargetGroupGet,
		InitClient: true,
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgTargetGroupId, cloudapiv6.ArgIdShort, "", cloudapiv6.TargetGroupId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgTargetGroupId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.TargetGroupIds(), cobra.ShellCompDirectiveNoFileComp
	})

	return cmd
}
