package version

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

func VersionListCmd() *core.Command {
	ctx := context.TODO()
	cmd := core.NewCommand(ctx, nil, core.CommandBuilder{
		Namespace:  "dbaas-postgres-v2",
		Resource:   "version",
		Verb:       "list",
		Aliases:    []string{"ls"},
		ShortDesc:  "List PostgreSQL Versions",
		LongDesc:   "Use this command to retrieve a list of available PostgreSQL Versions.",
		Example:    "ionosctl dbaas postgres version list",
		PreCmdRun:  core.NoPreRun,
		CmdRun:     RunVersionList,
		InitClient: true,
	})
	cmd.AddStringSliceFlag(constants.ArgCols, "", defaultVersionCols, tabheaders.ColsMessage(allVersionCols))
	return cmd
}

func RunVersionList(c *core.CommandConfig) error {

	versions, _, err := client.Must().PostgresClientV2.VersionsApi.VersionsGet(context.Background()).Execute()
	if err != nil {
		return err
	}

	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

	out, err := jsontabwriter.GenerateOutput("items", jsonpaths.DbaasPostgresV2Cluster, versions,
		tabheaders.GetHeaders(allVersionCols, defaultVersionCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)
	return nil
}
