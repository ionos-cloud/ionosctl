package version

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/mariadb-v2/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/spf13/viper"
)

func VersionGetCmd() *core.Command {
	ctx := context.TODO()
	cmd := core.NewCommand(ctx, nil, core.CommandBuilder{
		Namespace:  "dbaas-mariadb-v2",
		Resource:   "version",
		Verb:       "get",
		Aliases:    []string{"g"},
		ShortDesc:  "Get a MariaDB Version",
		LongDesc:   "Use this command to retrieve details about a MariaDB Version by using its ID.\n\nRequired values to run command:\n\n* Version Id",
		Example:    "ionosctl dbaas mariadb-v2 version get --version-id <version-id>",
		PreCmdRun:  PreRunVersionId,
		CmdRun:     RunVersionGet,
		InitClient: true,
	})
	cmd.AddStringFlag(constants.FlagVersionId, constants.FlagIdShort, "", "The ID of the MariaDB Version", core.RequiredFlagOption(),
		core.WithCompletion(completer.VersionIds, constants.MariaDBApiRegionalURL, constants.MariaDBLocations),
	)
	return cmd
}

func PreRunVersionId(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, constants.FlagVersionId)
}

func RunVersionGet(c *core.CommandConfig) error {
	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
	versionId := viper.GetString(core.GetFlagName(c.NS, constants.FlagVersionId))

	c.Verbose("Getting Version...")

	version, _, err := client.Must().MariaClientV2.VersionsApi.VersionsFindById(context.Background(), versionId).Execute()
	if err != nil {
		return err
	}

	return c.Out(table.Sprint(versionCols, version, cols))
}
