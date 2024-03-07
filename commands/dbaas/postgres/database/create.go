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
	ionoscloud "github.com/ionos-cloud/sdk-go-dbaas-postgres"
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
			LongDesc:  `Create a new database in the specified cluster.`,
			Example:   `ionosctl dbaas postgres database create --cluster-id <cluster-id> --database <database> --owner <owner>`,
			PreCmdRun: preRunCreateCmd,
			CmdRun:    runCreateCmd,
		},
	)
	c.AddStringFlag(constants.FlagClusterId, "", "", "The ID of the cluster")
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

	c.AddStringFlag(constants.FlagDatabase, "", "", "The name of the database to create")

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

	databaseProps := ionoscloud.DatabaseProperties{Name: &databaseName, Owner: &owner}
	database, _, err := client.Must().PostgresClient.DatabasesApi.DatabasesPost(
		context.Background(), clusterId,
	).Database(ionoscloud.Database{Properties: &databaseProps}).Execute()
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
