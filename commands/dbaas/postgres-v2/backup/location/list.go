package location

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/spf13/viper"
)

func BackupLocationListCmd() *core.Command {
	ctx := context.TODO()
	list := core.NewCommand(ctx, nil, core.CommandBuilder{
		Namespace:  "dbaas-postgres-v2",
		Resource:   "backup-location",
		Verb:       "list",
		Aliases:    []string{"ls"},
		ShortDesc:  "List PostgreSQL Backup Locations",
		LongDesc:   "Use this command to retrieve a list of PostgreSQL Backup Locations.",
		Example:    "ionosctl dbaas postgres backup location list",
		PreCmdRun:  core.NoPreRun,
		CmdRun:     RunBackupLocationList,
		InitClient: true,
	})
	list.AddInt32Flag(constants.FlagLimit, "", 100, "The limit of the number of items to return")
	list.AddInt32Flag(constants.FlagOffset, "", 0, "The offset of the listing")

	list.AddStringSliceFlag(constants.ArgCols, "", defaultBackupLocationCols, tabheaders.ColsMessage(allBackupLocationCols))

	return list
}

func RunBackupLocationList(c *core.CommandConfig) error {
	req := client.Must().PostgresClientV2.BackupLocationsApi.BackuplocationsGet(context.Background())

	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagLimit)) {
		req = req.Limit(viper.GetInt32(core.GetFlagName(c.NS, constants.FlagLimit)))
	}
	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagOffset)) {
		req = req.Offset(viper.GetInt32(core.GetFlagName(c.NS, constants.FlagOffset)))
	}

	locations, _, err := req.Execute()
	if err != nil {
		return err
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.NS, constants.ArgCols))

	out, err := jsontabwriter.GenerateOutput("items", jsonpaths.DbaasPostgresV2BackupLocation, locations,
		tabheaders.GetHeaders(allBackupLocationCols, defaultBackupLocationCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)
	return nil
}
