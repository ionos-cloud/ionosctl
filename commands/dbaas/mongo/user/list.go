package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/mongo/cluster"
	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/mongo/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/printer"
	sdkgo "github.com/ionos-cloud/sdk-go-dbaas-mongo"
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
		Example:   "ionosctl dbaas mongo user list",
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

			req := client.Must().MongoClient.UsersApi.ClustersUsersGet(context.Background(), clusterId)
			c.Printer.Verbose("Getting Users from all cluster %s", clusterId)

			if f := core.GetFlagName(c.NS, constants.FlagMaxResults); viper.IsSet(f) {
				req = req.Limit(viper.GetInt32(f))
			}
			if f := core.GetFlagName(c.NS, constants.FlagOffset); viper.IsSet(f) {
				req = req.Offset(viper.GetInt32(f))
			}

			ls, _, err := req.Execute()
			if err != nil {
				return err
			}
			return c.Printer.Print(getUserPrint(c, ls.GetItems()))
		},
		InitClient: true,
	})

	cmd.AddStringFlag(flagFilterByClusterNameWhenListAll, "", "",
		"When listing all users, you can optionally filter by partial-match cluster name")

	cmd.AddStringFlag(constants.FlagClusterId, constants.FlagIdShort, "", "")
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.MongoClusterIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddBoolFlag(constants.ArgNoHeaders, "", false, "When using text output, don't print headers")
	cmd.AddStringSliceFlag(constants.ArgCols, "", nil, printer.ColsMessage(allCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allCols, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddInt32Flag(constants.FlagMaxResults, constants.FlagMaxResultsShort, 0, constants.DescMaxResults)
	cmd.AddInt32Flag(constants.FlagOffset, "", 0, "Skip a certain number of results")

	cmd.Command.SilenceUsage = true

	cmd.AddBoolFlag(constants.ArgAll, "", true, "This flag exists for backward-compatibility reasons. This is now the default behaviour")
	_ = cmd.Command.Flags().MarkHidden(constants.ArgAll)

	return cmd
}

func listAll(c *core.CommandConfig) error {
	c.Printer.Verbose("Getting Users from all clusters...")
	clusters, err := cluster.Clusters(cluster.FilterNameFlags(c))
	if err != nil {
		return fmt.Errorf("failed getting clusters: %w", err)
	}

	var ls []sdkgo.User
	var multiErr error
	for _, c := range *clusters.GetItems() {
		l, _, err := client.Must().MongoClient.UsersApi.ClustersUsersGet(context.Background(), *c.Id).Execute()
		if err != nil {
			multiErr = errors.Join(multiErr, fmt.Errorf("failed listing users of cluster %s: %w", *c.Properties.DisplayName, err))
		}
		ls = append(ls, *l.GetItems()...)
	}
	if multiErr != nil {
		return fmt.Errorf("failed getting users of at least one cluster: %w", err)
	}

	return c.Printer.Print(getUserPrint(c, &ls))
}
