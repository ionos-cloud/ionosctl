package user

import (
	"context"

	cloudapiv6cmds "github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/cobra"
)

func UserGetCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace:  "user",
		Resource:   "user",
		Verb:       "get",
		Aliases:    []string{"g"},
		ShortDesc:  "Get a User",
		LongDesc:   "Use this command to retrieve details about a specific User.\n\nRequired values to run command:\n\n* User Id",
		Example:    "ionosctl user get --user-id USER_ID",
		PreCmdRun:  cloudapiv6cmds.PreRunUserId,
		CmdRun:     cloudapiv6cmds.RunUserGet,
		InitClient: true,
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgUserId, cloudapiv6.ArgIdShort, "", cloudapiv6.UserId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgUserId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.UsersIds(), cobra.ShellCompDirectiveNoFileComp
	})

	return cmd
}
