package group

import (
	"context"

	cloudapiv6cmds "github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
)

var (
	defaultResourceCols = []string{"ResourceId", "Name", "SecAuthProtection", "Type", "State"}
)

func GroupResourceCmd() *core.Command {
	resourceCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "resource",
			Aliases:          []string{"res"},
			Short:            "Group Resource Operations",
			Long:             "The sub-command of `ionosctl group resource` allows you to list Resources from a Group.",
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
		ShortDesc:  "List Resources from a Group",
		LongDesc:   "Use this command to get a list of Resources assigned to a Group. To see more details about existing Resources, use `ionosctl resource` commands.\n\nRequired values to run command:\n\n* Group Id",
		Example:    "ionosctl group resource list --group-id GROUP_ID",
		PreCmdRun:  cloudapiv6cmds.PreRunGroupId,
		CmdRun:     cloudapiv6cmds.RunGroupResourceList,
		InitClient: true,
	})
	cmd.AddStringSliceFlag(constants.ArgCols, "", defaultResourceCols, tabheaders.ColsMessage(defaultResourceCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultResourceCols, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddUUIDFlag(cloudapiv6.ArgGroupId, "", "", cloudapiv6.GroupId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgGroupId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.GroupsIds(), cobra.ShellCompDirectiveNoFileComp
	})

	return cmd
}
