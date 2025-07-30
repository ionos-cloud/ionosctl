package user

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
	cmd.Command.Flags().StringSlice(constants.ArgCols, []string{}, tabheaders.ColsMessage(allCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(
		constants.ArgCols,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return allCols, cobra.ShellCompDirectiveNoFileComp
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

	out, err := jsontabwriter.GenerateOutput(
		"items", jsonpaths.DbaasPostgresUser, users,
		tabheaders.GetHeadersAllDefault(defaultCols, cols),
	)
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)
	return nil
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

	var usersRaw []psql.UserList
	var usersConverted []map[string]interface{}
	for _, cluster := range clusters {
		tempUsers, tempConverted, err := getUsersFromCluster(cluster, getSystemUsers)
		if err != nil {
			return err
		}

		usersRaw = append(usersRaw, tempUsers)
		usersConverted = append(usersConverted, tempConverted...)
	}

	out, err := jsontabwriter.GenerateOutputPreconverted(
		usersRaw, usersConverted,
		tabheaders.GetHeaders(allCols, defaultCols, cols),
	)
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)
	return nil
}

func getUsersFromCluster(cluster psql.ClusterResponse, getSystemUsers bool) (
	psql.UserList, []map[string]interface{}, error,
) {
	clusterId, ok := cluster.GetIdOk()
	if !ok || clusterId == nil {
		return psql.UserList{}, nil, fmt.Errorf("failed to retrieve Postgres Cluster ID")
	}

	userList, _, err := client.Must().PostgresClient.UsersApi.UsersList(
		context.Background(), *clusterId,
	).System(getSystemUsers).Execute()
	if err != nil {
		return psql.UserList{}, nil, err
	}

	users, ok := userList.GetItemsOk()
	if !ok || users == nil {
		return psql.UserList{}, nil, fmt.Errorf("failed to retrieve Postgres Users")
	}

	convertedUserList := functional.Map(
		users, func(u psql.UserResource) map[string]interface{} {
			uConv, err := json2table.ConvertJSONToTable("", jsonpaths.DbaasPostgresUser, u)
			if err != nil {
				return nil
			}

			uConv[0]["ClusterId"] = *clusterId
			return uConv[0]
		},
	)

	return userList, convertedUserList, nil
}
