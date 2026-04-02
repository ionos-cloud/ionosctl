package group

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
)

func GroupListCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace:  "group",
		Resource:   "group",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List Groups",
		LongDesc:   "Use this command to get a list of available Groups available on your account\n\nYou can filter the results using `--filters` option. Use the following format to set filters: `--filters KEY1=VALUE1,KEY2=VALUE2`.\n" + completer.GroupsFiltersUsage(),
		Example:    "ionosctl compute group list",
		PreCmdRun:  core.NoPreRun,
		CmdRun:     RunGroupList,
		InitClient: true,
	})

	return cmd
}
