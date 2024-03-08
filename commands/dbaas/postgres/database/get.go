package database

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/postgres/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func GetCmd() *core.Command {
	c := core.NewCommand(
		context.Background(), nil, core.CommandBuilder{
			Verb:      "get",
			Namespace: "dbaas-postgres",
			Resource:  "database",
			ShortDesc: "Get database",
			LongDesc:  `Get the specified database from the given cluster`,
			Example:   `ionosctl dbaas postgres database get --cluster-id <cluster-id> --database <database>`,
			PreCmdRun: preRunGetCmd,
			CmdRun:    runGetCmd,
		},
	)
	c.Command.Flags().StringSlice(constants.ArgCols, []string{}, tabheaders.ColsMessage(allCols))
	_ = c.Command.RegisterFlagCompletionFunc(
		constants.ArgCols,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return allCols, cobra.ShellCompDirectiveNoFileComp
		},
	)

	c.AddStringFlag(constants.FlagClusterId, constants.FlagIdShort, "", "The ID of the Postgres cluster")
	_ = c.Command.RegisterFlagCompletionFunc(
		constants.FlagClusterId,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return completer.ClustersIds(), cobra.ShellCompDirectiveNoFileComp
		},
	)

	c.AddStringFlag(constants.FlagDatabase, "", "", "The name of the database")
	_ = c.Command.RegisterFlagCompletionFunc(
		constants.FlagDatabase,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return completer.DatabaseNames(c), cobra.ShellCompDirectiveNoFileComp
		},
	)

	return c
}

func preRunGetCmd(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, constants.FlagClusterId, constants.FlagDatabase)
}

func runGetCmd(c *core.CommandConfig) error {
	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
	clusterId := viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))
	databaseName := viper.GetString(core.GetFlagName(c.NS, constants.FlagDatabase))

	database, _, err := client.Must().PostgresClient.DatabasesApi.DatabasesGet(
		context.Background(), clusterId,
		databaseName,
	).Execute()
	if err != nil {
		return err
	}

	out, err := jsontabwriter.GenerateOutput(
		"", jsonpaths.DbaasPostgresDatabase, database,
		tabheaders.GetHeadersAllDefault(allCols, cols),
	)
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)
	return nil
}
