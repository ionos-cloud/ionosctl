package database

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/postgres/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
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

	return c.Out(table.Sprint(allCols, databases, cols, table.WithPrefix("items")))
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

	// Collect databases from all clusters, injecting ClusterId into each item's raw structure.
	var rows []map[string]any
	for _, cluster := range clusters {
		clusterId, ok := cluster.GetIdOk()
		if !ok || clusterId == nil {
			return fmt.Errorf("failed to retrieve Postgres Cluster ID")
		}

		databaseList, _, err := client.Must().PostgresClient.DatabasesApi.DatabasesList(
			context.Background(), *clusterId,
		).Execute()
		if err != nil {
			return err
		}

		databases, ok := databaseList.GetItemsOk()
		if !ok || databases == nil {
			return fmt.Errorf("failed to retrieve Postgres Databases")
		}

		for _, db := range databases {
			b, err := json.Marshal(db)
			if err != nil {
				return err
			}
			var m map[string]any
			if err := json.Unmarshal(b, &m); err != nil {
				return err
			}
			m["ClusterId"] = *clusterId
			rows = append(rows, m)
		}
	}

	return c.Out(table.Sprint(allCols, rows, cols))
}
