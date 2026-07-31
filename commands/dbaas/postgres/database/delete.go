package database

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/postgres/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/confirm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func DeleteCmd() *core.Command {
	c := core.NewCommand(
		context.Background(), nil, core.CommandBuilder{
			Verb:      "delete",
			Namespace: "dbaas-postgres",
			Resource:  "database",
			ShortDesc: "Delete database",
			LongDesc:  `Delete a logical database and all of its data from the given cluster. This is irreversible. The owning user is not affected.` + "\n\n" + `Required values to run command:` + "\n\n" + `* Cluster Id` + "\n" + `* Database (name)`,
			Example:   `ionosctl dbaas postgres database delete --cluster-id CLUSTER_ID --database orders`,
			PreCmdRun: preRunDeleteCmd,
			CmdRun:    runDeleteCmd,
		},
	)

	c.AddStringFlag(constants.FlagClusterId, constants.FlagIdShort, "", "ID of the PostgreSQL cluster the database belongs to")
	_ = c.Command.RegisterFlagCompletionFunc(
		constants.FlagClusterId,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return completer.ClustersIds(), cobra.ShellCompDirectiveNoFileComp
		},
	)

	c.AddStringFlag(constants.FlagDatabase, "", "", "Name of the database to delete")
	_ = c.Command.RegisterFlagCompletionFunc(
		constants.FlagDatabase,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return completer.DatabaseNames(c), cobra.ShellCompDirectiveNoFileComp
		},
	)

	return c
}

func preRunDeleteCmd(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, constants.FlagClusterId, constants.FlagDatabase)
}

func runDeleteCmd(c *core.CommandConfig) error {
	clusterId := viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))
	databaseName := viper.GetString(core.GetFlagName(c.NS, constants.FlagDatabase))

	if !confirm.FAsk(
		c.Command.Command.InOrStdin(), fmt.Sprintf("delete database %s from cluster %s", databaseName, clusterId),
		viper.GetBool(constants.ArgForce),
	) {
		return fmt.Errorf(confirm.UserDenied)
	}

	_, err := client.Must().PostgresClient.DatabasesApi.DatabasesDelete(
		context.Background(), clusterId, databaseName,
	).Execute()
	if err != nil {
		return err
	}

	c.Msg("DbaaS Postgres database %v successfully deleted", databaseName)
	return nil
}
