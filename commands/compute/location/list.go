package location

import (
	"context"

	cloudapiv6cmds "github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
)

func LocationListCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace:  "location",
		Resource:   "location",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List Locations",
		LongDesc:   "Use this command to get a list of available locations to create objects on.\n\nYou can filter the results using `--filters` option. Use the following format to set filters: `--filters KEY1=VALUE1,KEY2=VALUE2`.\n" + completer.LocationsFiltersUsage(),
		Example:    `ionosctl location list`,
		PreCmdRun:  core.NoPreRun,
		CmdRun:     cloudapiv6cmds.RunLocationList,
		InitClient: true,
	})

	return cmd
}
