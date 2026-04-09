package location

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/ionos-cloud/ionosctl/v6/pkg/functional"
	psqlv2 "github.com/ionos-cloud/sdk-go-bundle/products/dbaas/psql/v3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func BackupLocationGetCmd() *core.Command {
	ctx := context.TODO()
	get := core.NewCommand(ctx, nil, core.CommandBuilder{
		Namespace:  "dbaas-postgres-v2",
		Resource:   "backup-location",
		Verb:       "get",
		Aliases:    []string{"g"},
		ShortDesc:  "Get a PostgreSQL Backup Location",
		LongDesc:   "Use this command to retrieve details about a PostgreSQL Backup Location by using its ID.\n\nRequired values to run command:\n\n* Backup Location Id",
		Example:    "ionosctl dbaas postgres-v2 backup location get --backup-location-id <backup-location-id>",
		PreCmdRun:  PreRunBackupLocationId,
		CmdRun:     RunBackupLocationGet,
		InitClient: true,
	})
	get.AddStringFlag(constants.FlagBackupLocationId, constants.FlagIdShort, "", "The unique ID of the Backup Location", core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(constants.FlagBackupLocationId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		locations, _, err := client.Must().PostgresClientV2.BackupLocationsApi.BackuplocationsGet(context.Background()).Execute()
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}
		return functional.Map(locations.Items, func(l psqlv2.BackupLocationRead) string {
			loc := ""
			if l.Properties.Location != nil {
				loc = *l.Properties.Location
			}
			return fmt.Sprintf("%s\t%s", l.Id, loc)
		}), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddStringSliceFlag(constants.ArgCols, "", defaultBackupLocationCols, table.ColsMessage(backupLocationCols))

	return get
}

func PreRunBackupLocationId(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, constants.FlagBackupLocationId)
}

func RunBackupLocationGet(c *core.CommandConfig) error {
	locationId := viper.GetString(core.GetFlagName(c.NS, constants.FlagBackupLocationId))

	c.Verbose(constants.FlagBackupLocationId, locationId)

	location, _, err := client.Must().PostgresClientV2.BackupLocationsApi.BackuplocationsFindById(context.Background(), locationId).Execute()
	if err != nil {
		return err
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.NS, constants.ArgCols))
	return c.Out(table.Sprint(backupLocationCols, location, cols))
}
