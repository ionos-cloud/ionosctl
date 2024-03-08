package database

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/postgres/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
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
			LongDesc:  `Delete the specified database from the given cluster`,
			Example:   `ionosctl dbaas postgres database delete --cluster-id <cluster-id> --database <database>`,
			PreCmdRun: preRunDeleteCmd,
			CmdRun:    runDeleteCmd,
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

	out := jsontabwriter.GenerateLogOutput("DbaaS Postgres database %v successfully deleted", databaseName)

	fmt.Fprint(c.Command.Command.OutOrStdout(), out)
	return nil
}
