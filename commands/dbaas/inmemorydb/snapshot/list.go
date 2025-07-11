package snapshot

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
)

func List() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "inmemorydb",
		Resource:  "snapshot",
		Verb:      "list",
		Aliases:   []string{"l", "ls"},
		ShortDesc: "List In-Memory DB Snapshots",
		Example:   "ionosctl dbaas inmemorydb snapshot list",
		PreCmdRun: func(c *core.PreCommandConfig) error {

			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {
			ls, _, err := client.Must().InMemoryDBClient.SnapshotApi.
				SnapshotsGet(context.Background()).Execute()
			if err != nil {
				return err
			}

			cols, _ := c.Command.Command.Flags().GetStringSlice(constants.FlagCols)

			out, err := jsontabwriter.GenerateOutput("items", jsonpaths.DbaasInMemoryDBSnapshot, ls,
				tabheaders.GetHeadersAllDefault(allCols, cols))
			if err != nil {
				return err
			}

			fmt.Fprintf(c.Command.Command.OutOrStdout(), out)
			return nil
		},
		InitClient: true,
	})

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	return cmd
}
