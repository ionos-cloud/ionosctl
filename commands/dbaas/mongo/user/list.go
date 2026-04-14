package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/mongo/cluster"
	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/mongo/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	sdkgo "github.com/ionos-cloud/sdk-go-bundle/products/dbaas/mongo/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const flagFilterByClusterNameWhenListAll = "cluster-name"

func UserListCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "dbaas-mongo",
		Resource:  "user",
		Verb:      "list",
		Aliases:   []string{"l", "ls"},
		ShortDesc: fmt.Sprintf("Retrieves a list of MongoDB users. "+
			"You can either list users of a certain cluster (--%s), "+
			"or all clusters with an optional partial-match name filter (--%s)",
			constants.FlagClusterId, flagFilterByClusterNameWhenListAll),
		Example: `ionosctl dbaas mongo user list
ionosctl dbaas mongo user list --cluster-name <cluster-name>,
ionosctl dbaas mongo user list --cluster-id <cluster-id>`,
		PreCmdRun: core.NoPreRun,
		CmdRun: func(c *core.CommandConfig) error {
			fnClusterId := core.GetFlagName(c.NS, constants.FlagClusterId)
			if !viper.IsSet(fnClusterId) {
				err := listAll(c)
				if err != nil {
					return fmt.Errorf("failed listing users across all clusters: %w", err)
				}
				return nil
			}
			clusterId := viper.GetString(fnClusterId)

			c.Verbose("Getting Users from cluster %s", clusterId)

			ls, _, err := client.Must().MongoClient.UsersApi.ClustersUsersGet(context.Background(), clusterId).Execute()
			if err != nil {
				return err
			}

			return c.Printer(allCols).Prefix("items").Print(ls)
		},
		InitClient: true,
	})

	cmd.AddStringFlag(flagFilterByClusterNameWhenListAll, "", "",
		"When listing all users, you can optionally filter by partial-match cluster name")

	cmd.AddStringFlag(constants.FlagClusterId, constants.FlagIdShort, "", "")
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.MongoClusterIds(), cobra.ShellCompDirectiveNoFileComp
	})

	cmd.Command.SilenceUsage = true

	cmd.AddBoolFlag(constants.ArgAll, "", true, "This flag exists for backward-compatibility reasons. This is now the default behaviour")
	_ = cmd.Command.Flags().MarkHidden(constants.ArgAll)

	return cmd
}

func listAll(c *core.CommandConfig) error {
	c.Verbose("Getting Users from all clusters...")
	clusters, err := cluster.Clusters(func(r sdkgo.ApiClustersGetRequest) sdkgo.ApiClustersGetRequest {
		return r.FilterName(core.GetFlagName(c.NS, flagFilterByClusterNameWhenListAll))
	})
	if err != nil {
		return fmt.Errorf("failed getting clusters: %w", err)
	}

	var allUsers []sdkgo.User
	var multiErr error

	for _, cl := range clusters.GetItems() {
		l, _, err := client.Must().MongoClient.UsersApi.ClustersUsersGet(context.Background(), *cl.Id).Execute()
		if err != nil {
			multiErr = errors.Join(multiErr, fmt.Errorf("failed listing users of cluster %s: %w", *cl.Properties.DisplayName, err))
			continue
		}

		allUsers = append(allUsers, l.GetItems()...)
	}
	if multiErr != nil {
		return fmt.Errorf("failed getting users of at least one cluster: %w", multiErr)
	}

	return c.Printer(allCols).Print(allUsers)
}
