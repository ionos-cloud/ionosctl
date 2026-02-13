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
	"github.com/ionos-cloud/sdk-go-bundle/products/dbaas/psql/v2"
	"github.com/spf13/cobra"
)

func CreateCmd() *core.Command {
	c := core.NewCommand(
		context.Background(), nil, core.CommandBuilder{
			Verb:      "create",
			Namespace: "dbaas-postgres",
			Resource:  "database",
			ShortDesc: "Create database",
			LongDesc:  `Create a new database in the specified cluster`,
			Example:   `ionosctl dbaas postgres database create --cluster-id <cluster-id> --database <database> --owner <owner>`,
			PreCmdRun: preRunCreateCmd,
			CmdRun:    runCreateCmd,
		},
	)
	c.AddStringFlag(constants.FlagClusterId, "", "", "The ID of the Postgres cluster")
	_ = c.Command.RegisterFlagCompletionFunc(
		constants.FlagClusterId,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return completer.ClustersIds(), cobra.ShellCompDirectiveNoFileComp
		},
	)

	c.AddStringFlag(constants.FlagOwner, "", "", "The owner of the database")
	_ = c.Command.RegisterFlagCompletionFunc(
		constants.FlagOwner,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return completer.UserNames(c), cobra.ShellCompDirectiveNoFileComp
		},
	)

	c.AddStringFlag(constants.FlagDatabase, "", "", "The name of the database")

	return c
}

func preRunCreateCmd(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(
		c.Command, c.NS, constants.FlagClusterId, constants.FlagDatabase, constants.FlagOwner,
	)
}

func runCreateCmd(c *core.CommandConfig) error {
	clusterId, err := c.Command.Command.Flags().GetString(constants.FlagClusterId)
	if err != nil {
		return err
	}
	databaseName, err := c.Command.Command.Flags().GetString(constants.FlagDatabase)
	if err != nil {
		return err
	}
	owner, err := c.Command.Command.Flags().GetString(constants.FlagOwner)
	if err != nil {
		return err
	}

	databaseProps := psql.DatabaseProperties{Name: databaseName, Owner: owner}
	database, _, err := client.Must().PostgresClient.DatabasesApi.DatabasesPost(
		context.Background(), clusterId,
	).Database(psql.Database{Properties: databaseProps}).Execute()
	if err != nil {
		return err
	}

	out, err := jsontabwriter.GenerateOutput("", jsonpaths.DbaasPostgresDatabase, database, defaultCols)
	if err != nil {
		return err
	}

	fmt.Fprint(c.Command.Command.OutOrStdout(), out)
	return nil
}
