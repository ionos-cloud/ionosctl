package version

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/resource2table"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/spf13/viper"
)

func VersionGetCmd() *core.Command {
	ctx := context.TODO()
	cmd := core.NewCommand(ctx, nil, core.CommandBuilder{
		Namespace:  "dbaas-postgres-v2",
		Resource:   "version",
		Verb:       "get",
		Aliases:    []string{"g"},
		ShortDesc:  "Get a PostgreSQL Version",
		LongDesc:   "Use this command to retrieve details about a PostgreSQL Version by using its ID.\n\nRequired values to run command:\n\n* Version Id",
		Example:    "ionosctl dbaas postgres version get --version-id <version-id>",
		PreCmdRun:  PreRunVersionId,
		CmdRun:     RunVersionGet,
		InitClient: true,
	})
	cmd.AddStringFlag(constants.FlagVersionId, constants.FlagIdShort, "", "The ID of the PostgreSQL Version", core.RequiredFlagOption())
	cmd.AddStringSliceFlag(constants.ArgCols, "", defaultVersionCols, tabheaders.ColsMessage(allVersionCols))
	return cmd
}

func PreRunVersionId(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, constants.FlagVersionId)
}

func RunVersionGet(c *core.CommandConfig) error {
	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
	versionId := viper.GetString(core.GetFlagName(c.NS, constants.FlagVersionId))

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Getting Version..."))

	version, _, err := client.Must().PostgresClientV2.VersionsApi.VersionsFindById(context.Background(), versionId).Execute()
	if err != nil {
		return err
	}

	convertedVersion, err := resource2table.ConvertDbaasPostgresVersionToTable(version)
	if err != nil {
		return err
	}

	out, err := jsontabwriter.GenerateOutputPreconverted(version, convertedVersion,
		tabheaders.GetHeaders(allVersionCols, defaultVersionCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)
	return nil
}
