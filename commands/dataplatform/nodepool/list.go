package nodepool

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/functional"
	"github.com/ionos-cloud/ionosctl/v6/pkg/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/pkg/printer"
	ionoscloud "github.com/ionos-cloud/sdk-go-dataplatform"

	"github.com/ionos-cloud/ionosctl/v6/commands/dataplatform/completer"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NodepoolListCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "dataplatform",
		Resource:  "nodepool",
		Verb:      "list",
		Aliases:   []string{"l", "ls"},
		ShortDesc: "List Dataplatform Nodepools of a certain cluster",
		Example:   "ionosctl dataplatform nodepool list",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return core.CheckRequiredFlagsSets(c.Command, c.NS, []string{constants.ArgAll}, []string{constants.FlagClusterId})
		},
		CmdRun: func(c *core.CommandConfig) error {
			if viper.GetBool(core.GetFlagName(c.NS, constants.ArgAll)) {
				return listAll(c)
			}

			fmt.Fprintf(c.Stderr, jsontabwriter.GenerateVerboseOutput("Getting Nodepools..."))

			clusterId := viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))

			np, _, err := client.Must().DataplatformClient.DataPlatformNodePoolApi.ClustersNodepoolsGet(c.Context, clusterId).Execute()
			if err != nil {
				return err
			}

			cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

			npConverted, err := convertNodePoolsToTable(np)
			if err != nil {
				return err
			}

			out, err := jsontabwriter.GenerateOutputPreconverted(np, npConverted, printer.GetHeaders(allCols, defaultCols, cols))
			if err != nil {
				return err
			}

			fmt.Fprintf(c.Stdout, out)
			return nil
		},
		InitClient: true,
	})

	cmd.AddBoolFlag(constants.ArgAll, constants.ArgAllShort, false, "List all account nodepools, by iterating through all clusters first. May invoke a lot of GET calls")
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataplatformClusterIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(constants.FlagClusterId, constants.FlagIdShort, "", "The unique ID of the cluster. Must conform to the UUID format")
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataplatformClusterIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.Command.SilenceUsage = true

	return cmd
}

func listAll(c *core.CommandConfig) error {
	fmt.Fprintf(c.Stderr, jsontabwriter.GenerateVerboseOutput("Getting all nodepools..."))

	ls, _, err := client.Must().DataplatformClient.DataPlatformClusterApi.ClustersGet(context.Background()).Execute()
	if err != nil {
		return err
	}
	clusterIds := functional.Map(*ls.GetItems(), func(t ionoscloud.ClusterResponseData) string {
		return *t.GetId()
	})

	nps := make([]ionoscloud.NodePoolResponseData, 0)
	npsConverted := make([]map[string]interface{}, 0)
	for _, cID := range clusterIds {
		np, _, err := client.Must().DataplatformClient.DataPlatformNodePoolApi.ClustersNodepoolsGet(c.Context, cID).Execute()
		if err != nil {
			return err
		}

		temp, err := convertNodePoolsToTable(np)
		if err != nil {
			return err
		}

		npsConverted = append(npsConverted, temp...)
		nps = append(nps, *np.GetItems()...)
	}

	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

	out, err := jsontabwriter.GenerateOutputPreconverted(nps, npsConverted, printer.GetHeaders(allCols, defaultCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Stdout, out)
	return nil
}
