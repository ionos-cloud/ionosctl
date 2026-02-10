package cluster

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/postgres-v2/waiter"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/resource2table"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/ionosctl/v6/internal/waitfor"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ClusterGetCmd() *core.Command {
	ctx := context.TODO()
	get := core.NewCommand(ctx, nil, core.CommandBuilder{
		Namespace:  "dbaas-postgres-v2",
		Resource:   "cluster",
		Verb:       "get",
		Aliases:    []string{"g"},
		ShortDesc:  "Get a PostgreSQL Cluster",
		Example:    "ionosctl dbaas postgres cluster get --cluster-id <cluster-id>",
		LongDesc:   "Use this command to retrieve details about a PostgreSQL Cluster by using its ID.\n\nRequired values to run command:\n\n* Cluster Id",
		PreCmdRun:  PreRunClusterId,
		CmdRun:     RunClusterGet,
		InitClient: true,
	})
	get.AddUUIDFlag(constants.FlagClusterId, constants.FlagIdShort, "", constants.DescCluster, core.RequiredFlagOption())
	get.AddBoolFlag(constants.ArgWaitForState, constants.ArgWaitForStateShort, constants.DefaultWait, "Wait for Cluster to be in AVAILABLE state")
	get.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, constants.DefaultClusterTimeout, "Timeout option for Cluster to be in AVAILABLE state [seconds]")
	get.AddStringSliceFlag(constants.ArgCols, "", defaultClusterCols, tabheaders.ColsMessage(allClusterCols))
	_ = get.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allClusterCols, cobra.ShellCompDirectiveNoFileComp
	})
	return get
}

func PreRunClusterId(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, constants.FlagClusterId)
}

func RunClusterGet(c *core.CommandConfig) error {
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.ClusterId, viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Getting Cluster..."))

	if viper.GetBool(core.GetFlagName(c.NS, constants.ArgWaitForState)) {
		if err := waitfor.WaitForState(c, waiter.ClusterStateInterrogator, viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))); err != nil {
			return err
		}
	}

	cluster, _, err := client.Must().PostgresClientV2.ClustersApi.ClustersFindById(
		context.Background(), viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))).Execute()
	if err != nil {
		return fmt.Errorf("could not get cluster: %w", err)
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.NS, constants.ArgCols))

	clusterConverted, err := resource2table.ConvertDbaasPostgresClusterToTableV2(cluster)
	if err != nil {
		return err
	}

	out, err := jsontabwriter.GenerateOutputPreconverted(cluster, clusterConverted,
		tabheaders.GetHeaders(allClusterCols, defaultClusterCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)
	return nil
}
