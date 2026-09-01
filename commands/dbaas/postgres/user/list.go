package user

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/postgres/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
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
			LongDesc:  `List the users of a cluster. Provide --cluster-id to list the users of one cluster; omit it to list users across every cluster (the ClusterId column identifies each user's cluster). System users are hidden unless --system is given.`,
			Example:   `ionosctl dbaas postgres user list --cluster-id CLUSTER_ID --system`,
			PreCmdRun: core.NoPreRun,
			CmdRun:    runListCmd,
		},
	)
	cmd.AddStringFlag(constants.FlagClusterId, constants.FlagIdShort, "", "ID of the PostgreSQL cluster whose users to list. If omitted, users from all clusters are listed")
	_ = cmd.Command.RegisterFlagCompletionFunc(
		constants.FlagClusterId,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return completer.ClustersIds(), cobra.ShellCompDirectiveNoFileComp
		},
	)

	cmd.AddBoolFlag("system", "", false, "Also include service-managed system users (e.g. postgres), which are hidden by default")

	return cmd
}

func runListCmd(c *core.CommandConfig) error {
	if !viper.IsSet(core.GetFlagName(c.NS, constants.FlagClusterId)) {
		return listAll(c)
	}

	getSystemUsers := viper.GetBool(core.GetFlagName(c.NS, "system"))
	clusterId := viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))

	users, _, err := client.Must().PostgresClient.UsersApi.UsersList(
		context.Background(),
		clusterId,
	).System(getSystemUsers).Execute()
	if err != nil {
		return err
	}

	return c.Printer(allCols).Prefix("items").Print(users)
}

func listAll(c *core.CommandConfig) error {
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

	return c.Printer(allCols).Print(allUsers)
}

func userToRow(u psql.UserResource, clusterId string) map[string]any {
	row := map[string]any{
		"ClusterId":  clusterId,
		"id":         u.Id,
		"properties": map[string]any{},
	}
	if props, ok := u.GetPropertiesOk(); ok && props != nil {
		row["properties"] = map[string]any{
			"username": props.GetUsername(),
		}
		if sys, ok := props.GetSystemOk(); ok && sys != nil {
			row["properties"].(map[string]any)["system"] = *sys
		}
	}
	return row
}
