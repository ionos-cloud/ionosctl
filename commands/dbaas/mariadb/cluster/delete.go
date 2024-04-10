package cluster

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/ionosctl/v6/pkg/confirm"
	"github.com/ionos-cloud/ionosctl/v6/pkg/functional"
	sdkgo "github.com/ionos-cloud/sdk-go-dbaas-mariadb"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func confirmStringForCluster(c sdkgo.ClusterResponse) string {
	askString := ""
	if p := c.Properties; p != nil {
		if c.Id != nil {
			askString = fmt.Sprintf("%s cluster %s", askString, *c.Id)
		}
		if n := p.DisplayName; n != nil {
			askString = fmt.Sprintf("%s (%s)", askString, *n)
		}
		if v := p.MariadbVersion; v != nil {
			askString = fmt.Sprintf("%s version v%s", askString, *v)
		}
	}
	return fmt.Sprintf("delete%s and its snapshots", askString)
}

func Delete() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "dbaas-mariadb",
		Resource:  "cluster",
		Verb:      "delete",
		Aliases:   []string{"del", "d"},
		ShortDesc: "Delete a MariaDB Cluster by ID",
		Example: `ionosctl dbaas mariadb cluster delete --cluster-id <cluster-id>
ionosctl db mar c d --all
ionosctl db mar c d --all --name <name>`,
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return core.CheckRequiredFlagsSets(c.Command, c.NS, []string{constants.ArgAll}, []string{constants.FlagClusterId})
		},
		CmdRun: func(c *core.CommandConfig) error {
			if all := viper.GetBool(core.GetFlagName(c.NS, constants.ArgAll)); all {
				return deleteAll(c)
			}

			clusterId := viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))
			chosenCluster, _, err := client.Must().MariaClient.ClustersApi.ClustersFindById(context.Background(), clusterId).Execute()
			if err != nil {
				wrapped := fmt.Errorf("failed trying to find cluster by id: %w", err)
				keepGoing := confirm.FAsk(c.Command.Command.InOrStdin(), fmt.Sprintf("%s, try deleting %s anyways", wrapped.Error(), clusterId))
				if !keepGoing {
					return wrapped
				}
			}

			ok := confirm.FAsk(c.Command.Command.InOrStdin(), confirmStringForCluster(chosenCluster), viper.GetBool(constants.ArgForce))
			if !ok {
				return fmt.Errorf(confirm.UserDenied)
			}
			fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Deleting cluster: %s", clusterId))

			_, _, err = client.Must().MariaClient.ClustersApi.ClustersDelete(context.Background(), clusterId).Execute()
			return err
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagClusterId, constants.FlagIdShort, "", "The unique ID of the cluster", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return ClustersProperty(func(c sdkgo.ClusterResponse) string {
			if c.Id == nil {
				return ""
			}
			return *c.Id
		}), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddBoolFlag(constants.ArgAll, constants.ArgAllShort, false, "Delete all mariadb clusters")
	cmd.AddBoolFlag(constants.FlagName, "", false, "When deleting all clusters, filter the clusters by a name")
	cmd.AddStringSliceFlag(constants.ArgCols, "", nil, tabheaders.ColsMessage(allCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allCols, cobra.ShellCompDirectiveNoFileComp
	})

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	return cmd
}

func deleteAll(c *core.CommandConfig) error {
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Deleting All Clusters!"))
	xs, err := Clusters(FilterNameFlags(c))
	if err != nil {
		return err
	}

	return functional.ApplyAndAggregateErrors(*xs.GetItems(), func(x sdkgo.ClusterResponse) error {
		yes := confirm.FAsk(c.Command.Command.InOrStdin(), confirmStringForCluster(x), viper.GetBool(constants.ArgForce))
		if yes {
			_, _, delErr := client.Must().MariaClient.ClustersApi.ClustersDelete(c.Context, *x.Id).Execute()
			if delErr != nil {
				return fmt.Errorf("failed deleting cluster %s: %w", *x.Properties.DisplayName, delErr)
			}
		}
		return nil
	})
}
