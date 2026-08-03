package location

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
)

func LocationListCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "location",
		Resource:  "location",
		Verb:      "list",
		Aliases:   []string{"l", "ls"},
		ShortDesc: "List all IONOS locations available for provisioning",
		LongDesc:  "Use this command to list every location (data center site) where you can create resources. The LocationId column (e.g. `de/txl`) is the value you pass as `--location` elsewhere.\n\nYou can filter the results using `--filters` option. Use the following format to set filters: `--filters KEY1=VALUE1,KEY2=VALUE2`.\n" + completer.LocationsFiltersUsage(),
		Example: `# List every available location
ionosctl compute location list

# Show only the ID and name columns
ionosctl compute location list --cols LocationId,Name`,
		PreCmdRun:  core.NoPreRun,
		CmdRun:     RunLocationList,
		InitClient: true,
	})
	return cmd
}
