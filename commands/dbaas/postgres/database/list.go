package database

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/postgres/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/ionosctl/v6/pkg/functional"
	psql "github.com/ionos-cloud/sdk-go-dbaas-postgres"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ListCmd() *core.Command {
	c := core.NewCommand(
		context.Background(), nil, core.CommandBuilder{
			Verb:      "list",
			Aliases:   []string{"ls"},
			Resource:  "database",
			Namespace: "dbaas-postgres",
			ShortDesc: "List databases",
			LongDesc:  `List databases in the given cluster`,
			Example:   `ionosctl dbaas postgres database list`,
			PreCmdRun: core.NoPreRun,
			CmdRun:    runListCmd,
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

	return c
}

func runListCmd(c *core.CommandConfig) error {
	if !viper.IsSet(core.GetFlagName(c.NS, constants.FlagClusterId)) {
		return listAll(c)
	}

	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
	clusterId := viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))

	databases, _, err := client.Must().PostgresClient.DatabasesApi.DatabasesList(
		context.Background(),
		clusterId,
	).Execute()
	if err != nil {
		return err
	}

	out, err := jsontabwriter.GenerateOutput(
		"items", jsonpaths.DbaasPostgresDatabase, databases,
		tabheaders.GetHeadersAllDefault(defaultCols, cols),
	)
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)
	return nil
}

func listAll(c *core.CommandConfig) error {
	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

	clusterList, _, err := client.Must().PostgresClient.ClustersApi.ClustersGet(context.Background()).Execute()
	if err != nil {
		return err
	}

	clusters, ok := clusterList.GetItemsOk()
	if !ok || clusters == nil {
		return fmt.Errorf("failed to retrieve Postgres Clusters")
	}

	var databasesRaw []psql.DatabaseList
	var usersConverted []map[string]interface{}
	for _, cluster := range *clusters {
		tempDatabases, tempConverted, err := getDatabasesFromCluster(cluster)
		if err != nil {
			return err
		}

		databasesRaw = append(databasesRaw, tempDatabases)
		usersConverted = append(usersConverted, tempConverted...)
	}

	out, err := jsontabwriter.GenerateOutputPreconverted(
		databasesRaw, usersConverted,
		tabheaders.GetHeaders(allCols, defaultCols, cols),
	)
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)
	return nil
}

func getDatabasesFromCluster(cluster psql.ClusterResponse) (
	psql.DatabaseList, []map[string]interface{}, error,
) {
	clusterId, ok := cluster.GetIdOk()
	if !ok || clusterId == nil {
		return psql.DatabaseList{}, nil, fmt.Errorf("failed to retrieve Postgres Cluster ID")
	}

	databaseList, _, err := client.Must().PostgresClient.DatabasesApi.DatabasesList(
		context.Background(), *clusterId,
	).Execute()
	if err != nil {
		return psql.DatabaseList{}, nil, err
	}

	databases, ok := databaseList.GetItemsOk()
	if !ok || databases == nil {
		return psql.DatabaseList{}, nil, fmt.Errorf("failed to retrieve Postgres Databases")
	}

	convertedDatabaseList := functional.Map(
		*databases, func(db psql.DatabaseResource) map[string]interface{} {
			dbConv, err := json2table.ConvertJSONToTable("", jsonpaths.DbaasPostgresDatabase, db)
			if err != nil {
				return nil
			}

			dbConv[0]["ClusterId"] = *clusterId
			return dbConv[0]
		},
	)

	return databaseList, convertedDatabaseList, nil
}
