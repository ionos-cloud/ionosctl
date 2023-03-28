package cluster

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"

	"github.com/ionos-cloud/ionosctl/v6/internal/confirm"
	"github.com/ionos-cloud/ionosctl/v6/internal/functional"
	ionoscloud "github.com/ionos-cloud/sdk-go-dataplatform"

	"github.com/ionos-cloud/ionosctl/v6/commands/dataplatform/completer"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ClusterDeleteCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "dataplatform",
		Resource:  "cluster",
		Verb:      "delete",
		Aliases:   []string{"del", "d"},
		ShortDesc: "Delete a Dataplatform Cluster by ID",
		Example:   "ionosctl dataplatform cluster delete --cluster-id <cluster-id>",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return core.CheckRequiredFlagsSets(c.Command, c.NS, []string{constants.ArgAll}, []string{constants.FlagClusterId})
		},
		CmdRun: func(c *core.CommandConfig) error {
			if all := viper.GetBool(core.GetFlagName(c.NS, constants.ArgAll)); all {
				return deleteAll(c)
			}

			clusterId := viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))
			err := utils.AskForConfirm(c.Stdin, c.Printer, fmt.Sprintf("delete cluster %s", clusterId))
			if err != nil {
				return err
			}
			c.Printer.Verbose("Deleting cluster: %s", clusterId)
			_, _, err = client.Must().DataplatformClient.DataPlatformClusterApi.ClustersDelete(c.Context, clusterId).Execute()
			if err != nil {
				return err
			}
			return err
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagClusterId, constants.FlagIdShort, "", "The unique ID of the cluster", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataplatformClusterIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddBoolFlag(constants.ArgAll, constants.ArgAllShort, false, "Delete all clusters")
	cmd.AddBoolFlag(constants.ArgForce, constants.ArgForceShort, false, "Skip yes/no verification")

	cmd.Command.SilenceUsage = true

	return cmd
}

func deleteAll(c *core.CommandConfig) error {
	c.Printer.Verbose("Deleting All Clusters!")
	xs, _, err := client.Must().DataplatformClient.DataPlatformClusterApi.ClustersGet(c.Context).Execute()
	if err != nil {
		return err
	}

	err = functional.ApplyAndAggregateErrors(*xs.GetItems(), func(x ionoscloud.ClusterResponseData) error {
		yes := confirm.Ask(fmt.Sprintf("delete cluster %s (%s)", *x.Id, *x.Properties.Name), viper.GetBool(core.GetFlagName(c.NS, constants.ArgForce)))
		if yes {
			_, _, delErr := client.Must().DataplatformClient.DataPlatformClusterApi.ClustersDelete(c.Context, *x.Id).Execute()
			if delErr != nil {
				return fmt.Errorf("failed deleting %s: %w", *x.Id, delErr)
			}
		}
		return nil
	})

	_, _, err = client.Must().DataplatformClient.DataPlatformClusterApi.ClustersDelete(c.Context, viper.GetString(core.GetFlagName(c.NS, constants.FlagName))).Execute()
	return err
}
