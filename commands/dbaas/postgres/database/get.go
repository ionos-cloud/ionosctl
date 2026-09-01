package database

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/postgres/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
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
			LongDesc:  `Retrieve a single database of a cluster by name, showing its owner.` + "\n\n" + `Required values to run command:` + "\n\n" + `* Cluster Id` + "\n" + `* Database (name)`,
			Example:   `ionosctl dbaas postgres database get --cluster-id CLUSTER_ID --database orders`,
			PreCmdRun: preRunGetCmd,
			CmdRun:    runGetCmd,
		},
	)

	c.AddStringFlag(constants.FlagClusterId, constants.FlagIdShort, "", "ID of the PostgreSQL cluster the database belongs to")
	_ = c.Command.RegisterFlagCompletionFunc(
		constants.FlagClusterId,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return completer.ClustersIds(), cobra.ShellCompDirectiveNoFileComp
		},
	)

	c.AddStringFlag(constants.FlagDatabase, "", "", "Name of the database to retrieve")
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
	clusterId := viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))
	databaseName := viper.GetString(core.GetFlagName(c.NS, constants.FlagDatabase))

	database, _, err := client.Must().PostgresClient.DatabasesApi.DatabasesGet(
		context.Background(), clusterId,
		databaseName,
	).Execute()
	if err != nil {
		return err
	}

	return c.Printer(allCols).Print(database)
}
