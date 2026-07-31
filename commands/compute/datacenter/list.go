package datacenter

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
)

func DatacenterListCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "datacenter",
		Resource:  "datacenter",
		Verb:      "list",
		Aliases:   []string{"l", "ls"},
		ShortDesc: "List all Virtual Data Centers on your account",
		LongDesc: "List every Virtual Data Center available on your account, across all regions, with their ID, name, region, CPU family and current state. This is the usual way to discover the `--datacenter-id` values that the other compute commands need.\n\nNarrow the results with `--filters` using the format `--filters KEY1=VALUE1,KEY2=VALUE2` (for example `--filters location=de/fra`).\n" +
			completer.DataCentersFiltersUsage(),
		Example:    "# List all VDCs\nionosctl compute datacenter list\n\n# List only the columns you care about\nionosctl compute datacenter list --cols \"DatacenterId,Name,Location,Version\"\n\n# List only VDCs in Frankfurt\nionosctl compute datacenter list --filters location=de/fra",
		PreCmdRun:  core.NoPreRun,
		CmdRun:     RunDataCenterList,
		InitClient: true,
	})

	return cmd
}
