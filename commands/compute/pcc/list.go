package pcc

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
)

func PccListCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace:  "pcc",
		Resource:   "pcc",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List Cross-Connects",
		LongDesc:   "Use this command to get a list of the Cross-Connects (Private Cross-Connects) available on your contract, across all regions. To see which LANs are peered through a specific Cross-Connect, use `ionosctl compute pcc peers list --pcc-id <PCC_ID>`.\n\nYou can filter the results using `--filters` option. Use the following format to set filters: `--filters KEY1=VALUE1,KEY2=VALUE2`.\n" + completer.PccsFiltersUsage(),
		Example:    `ionosctl compute pcc list`,
		PreCmdRun:  core.NoPreRun,
		CmdRun:     RunPccList,
		InitClient: true,
	})

	return cmd
}
