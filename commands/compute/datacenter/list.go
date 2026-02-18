package datacenter

import (
	"context"

	cloudapiv6cmds "github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
)

func DatacenterListCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "datacenter",
		Resource:  "datacenter",
		Verb:      "list",
		Aliases:   []string{"l", "ls"},
		ShortDesc: "List Data Centers",
		LongDesc: "Use this command to get a list of existing Data Centers available on your account.\n\nYou can filter the results using `--filters` option. Use the following format to set filters: `--filters KEY1=VALUE1,KEY2=VALUE2`.\n" +
			completer.DataCentersFiltersUsage(),
		Example:    "ionosctl datacenter list\nionosctl datacenter list --cols \"DatacenterId,Name,Location,Version\"",
		PreCmdRun:  core.NoPreRun,
		CmdRun:     cloudapiv6cmds.RunDataCenterList,
		InitClient: true,
	})

	return cmd
}
