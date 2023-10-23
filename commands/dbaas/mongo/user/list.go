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
	"github.com/ionos-cloud/ionosctl/v6/pkg/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/pkg/resource2table"
	"github.com/ionos-cloud/ionosctl/v6/pkg/tabheaders"
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

			req := client.Must().MongoClient.UsersApi.ClustersUsersGet(context.Background(), clusterId)
			fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Getting Users from all cluster %s", clusterId))

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

			cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

			lsConverted, err := resource2table.ConvertDbaasMongoUsersToTable(ls)
			if err != nil {
				return err
			}

			out, err := jsontabwriter.GenerateOutputPreconverted(ls, lsConverted, tabheaders.GetHeadersAllDefault(allCols, cols))
			if err != nil {
				return err
			}

			fmt.Fprintf(c.Command.Command.OutOrStdout(), out)
			return nil
		},
		InitClient: true,
	})

	cmd.AddStringFlag(flagFilterByClusterNameWhenListAll, "", "",
		"When listing all users, you can optionally filter by partial-match cluster name")

	cmd.AddStringFlag(constants.FlagClusterId, constants.FlagIdShort, "", "")
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.MongoClusterIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringSliceFlag(constants.ArgCols, "", nil, tabheaders.ColsMessage(allCols))
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
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Getting Users from all clusters..."))
	clusters, err := cluster.Clusters(func(r sdkgo.ApiClustersGetRequest) sdkgo.ApiClustersGetRequest {
		return r.FilterName(core.GetFlagName(c.NS, flagFilterByClusterNameWhenListAll))
	})
	if err != nil {
		return fmt.Errorf("failed getting clusters: %w", err)
	}

	var ls []sdkgo.UsersList
	var lsConverted []map[string]interface{}
	var multiErr error

	for _, c := range *clusters.GetItems() {
		l, _, err := client.Must().MongoClient.UsersApi.ClustersUsersGet(context.Background(), *c.Id).Execute()
		if err != nil {
			multiErr = errors.Join(multiErr, fmt.Errorf("failed listing users of cluster %s: %w", *c.Properties.DisplayName, err))
		}

		temp, err := resource2table.ConvertDbaasMongoUsersToTable(l)
		if err != nil {
			multiErr = errors.Join(multiErr, fmt.Errorf("failed converting users of cluster %s: %w", *c.Properties.DisplayName, err))
			continue
		}

		ls = append(ls, l)
		lsConverted = append(lsConverted, temp...)
	}
	if multiErr != nil {
		return fmt.Errorf("failed getting users of at least one cluster: %w", err)
	}

	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

	out, err := jsontabwriter.GenerateOutputPreconverted(ls, lsConverted, tabheaders.GetHeadersAllDefault(allCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)
	return nil
}
