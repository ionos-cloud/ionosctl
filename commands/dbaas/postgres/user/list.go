package user

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/postgres/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/ionos-cloud/sdk-go-bundle/products/dbaas/psql/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ListCmd() *core.Command {
	cmd := core.NewCommand(
		context.Background(), nil, core.CommandBuilder{
			Verb:      "list",
			Aliases:   []string{"ls"},
			Namespace: "dbaas-postgres",
			Resource:  "user",
			ShortDesc: "List users",
			LongDesc:  `List all users in the given cluster`,
			Example:   `ionosctl dbaas postgres user list --cluster-id <cluster-id>`,
			PreCmdRun: core.NoPreRun,
			CmdRun:    runListCmd,
		},
	)
	cmd.Command.Flags().StringSlice(constants.ArgCols, []string{}, table.ColsMessage(allCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(
		constants.ArgCols,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return table.AllCols(allCols), cobra.ShellCompDirectiveNoFileComp
		},
	)

	cmd.AddStringFlag(constants.FlagClusterId, constants.FlagIdShort, "", "The ID of the Postgres cluster")
	_ = cmd.Command.RegisterFlagCompletionFunc(
		constants.FlagClusterId,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return completer.ClustersIds(), cobra.ShellCompDirectiveNoFileComp
		},
	)

	cmd.AddBoolFlag("system", "", false, "List system users along with regular users")

	return cmd
}

func runListCmd(c *core.CommandConfig) error {
	if !viper.IsSet(core.GetFlagName(c.NS, constants.FlagClusterId)) {
		return listAll(c)
	}

	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
	getSystemUsers := viper.GetBool(core.GetFlagName(c.NS, "system"))
	clusterId := viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))

	users, _, err := client.Must().PostgresClient.UsersApi.UsersList(
		context.Background(),
		clusterId,
	).System(getSystemUsers).Execute()
	if err != nil {
		return err
	}

	return c.Out(table.Sprint(allCols, users, cols, table.WithPrefix("items")))
}

func listAll(c *core.CommandConfig) error {
	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
	getSystemUsers := viper.GetBool(core.GetFlagName(c.NS, "system"))

	clusterList, _, err := client.Must().PostgresClient.ClustersApi.ClustersGet(context.Background()).Execute()
	if err != nil {
		return err
	}

	clusters, ok := clusterList.GetItemsOk()
	if !ok || clusters == nil {
		return fmt.Errorf("failed to retrieve Postgres Clusters")
	}

	var allUsers []map[string]any
	for _, cluster := range clusters {
		clusterId, ok := cluster.GetIdOk()
		if !ok || clusterId == nil {
			continue
		}

		userList, _, err := client.Must().PostgresClient.UsersApi.UsersList(
			context.Background(), *clusterId,
		).System(getSystemUsers).Execute()
		if err != nil {
			return err
		}

		users, ok := userList.GetItemsOk()
		if !ok || users == nil {
			continue
		}

		for _, u := range users {
			allUsers = append(allUsers, userToRow(u, *clusterId))
		}
	}

	return c.Out(table.Sprint(allCols, allUsers, cols))
}

func userToRow(u psql.UserResource, clusterId string) map[string]any {
	row := map[string]any{
		"ClusterId": clusterId,
		"Id":        u.Id,
	}
	if props, ok := u.GetPropertiesOk(); ok && props != nil {
		row["Username"] = props.GetUsername()
		if sys, ok := props.GetSystemOk(); ok && sys != nil {
			row["System"] = *sys
		}
	}
	return row
}
