package location

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/mariadb-v2/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/spf13/viper"
)

func LocationGetCmd() *core.Command {
	ctx := context.TODO()
	cmd := core.NewCommand(ctx, nil, core.CommandBuilder{
		Namespace:  "dbaas-mariadb-v2",
		Resource:   "location",
		Verb:       "get",
		Aliases:    []string{"g"},
		ShortDesc:  "Get a MariaDB Backup Location",
		LongDesc:   "Use this command to retrieve details about a MariaDB Backup Location by using its ID.\n\nRequired values to run command:\n\n* Backup Location Id",
		Example:    "ionosctl dbaas mariadb-v2 backup location get --backup-location-id <backup-location-id>",
		PreCmdRun:  PreRunLocationId,
		CmdRun:     RunLocationGet,
		InitClient: true,
	})
	cmd.AddStringFlag(constants.FlagBackupLocationId, constants.FlagIdShort, "", "The ID of the MariaDB Backup Location", core.RequiredFlagOption(),
		core.WithCompletion(completer.BackupLocationIds, constants.MariaDBApiRegionalURL, constants.MariaDBLocations),
	)
	return cmd
}

func PreRunLocationId(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, constants.FlagBackupLocationId)
}

func RunLocationGet(c *core.CommandConfig) error {
	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
	locationId := viper.GetString(core.GetFlagName(c.NS, constants.FlagBackupLocationId))

	c.Verbose("Getting Backup Location...")

	location, _, err := client.Must().MariaClientV2.BackupLocationsApi.BackuplocationsFindById(context.Background(), locationId).Execute()
	if err != nil {
		return err
	}

	return c.Out(table.Sprint(locationCols, location, cols))
}
