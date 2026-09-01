package database

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/postgres/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/sdk-go-bundle/products/dbaas/psql/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func CreateCmd() *core.Command {
	c := core.NewCommand(
		context.Background(), nil, core.CommandBuilder{
			Verb:      "create",
			Namespace: "dbaas-postgres",
			Resource:  "database",
			ShortDesc: "Create database",
			LongDesc: `Create a new logical database in the specified cluster.

The --owner must be an existing user (role) in the same cluster and is granted ownership (full privileges) of the new database. The cluster must be AVAILABLE.

Required values to run command:

* Cluster Id
* Database (name)
* Owner`,
			Example:   `ionosctl dbaas postgres database create --cluster-id CLUSTER_ID --database orders --owner appuser`,
			PreCmdRun: preRunCreateCmd,
			CmdRun:    runCreateCmd,
		},
	)
	c.AddStringFlag(constants.FlagClusterId, "", "", "ID of the PostgreSQL cluster to create the database in (must be AVAILABLE)")
	_ = c.Command.RegisterFlagCompletionFunc(
		constants.FlagClusterId,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return completer.ClustersIds(), cobra.ShellCompDirectiveNoFileComp
		},
	)

	c.AddStringFlag(constants.FlagOwner, "", "", "Name of an existing user (role) in the same cluster that will own the database and hold full privileges on it")
	_ = c.Command.RegisterFlagCompletionFunc(
		constants.FlagOwner,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return completer.UserNames(c), cobra.ShellCompDirectiveNoFileComp
		},
	)

	c.AddStringFlag(constants.FlagDatabase, "", "", "Name of the database to create (1-63 characters)")

	return c
}

func preRunCreateCmd(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(
		c.Command, c.NS, constants.FlagClusterId, constants.FlagDatabase, constants.FlagOwner,
	)
}

func runCreateCmd(c *core.CommandConfig) error {
	clusterId := viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))
	databaseName := viper.GetString(core.GetFlagName(c.NS, constants.FlagDatabase))
	owner := viper.GetString(core.GetFlagName(c.NS, constants.FlagOwner))

	databaseProps := psql.DatabaseProperties{Name: databaseName, Owner: owner}
	database, _, err := client.Must().PostgresClient.DatabasesApi.DatabasesPost(
		context.Background(), clusterId,
	).Database(psql.Database{Properties: databaseProps}).Execute()
	if err != nil {
		return err
	}

	return c.Printer(allCols).Print(database)
}
