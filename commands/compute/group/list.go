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
		ShortDesc:  "List all IAM Groups on the contract",
		LongDesc:   "List every IAM Group on your contract. The default columns show the Group's ID, name and its most common privilege toggles; add more privilege columns with `--cols`.\n\nYou can filter the results using `--filters` option. Use the following format to set filters: `--filters KEY1=VALUE1,KEY2=VALUE2`.\n" + completer.GroupsFiltersUsage(),
		Example:    "ionosctl compute group list",
		PreCmdRun:  core.NoPreRun,
		CmdRun:     RunGroupList,
		InitClient: true,
	})

	return cmd
}
