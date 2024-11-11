package cluster

import (
	"context"
	"fmt"
	"github.com/ionos-cloud/ionosctl/v6/commands/kafka/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	kafka "github.com/ionos-cloud/sdk-go-kafka"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func FindByID() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "kafka",
		Resource:  "cluster",
		Verb:      "get",
		Aliases:   []string{"g"},
		ShortDesc: "Retrieve a cluster",
		Example:   "ionosctl kafka cl get --cluster-id ID",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			if err := core.CheckRequiredFlags(c.Command, c.NS, constants.FlagClusterId, constants.FlagLocation); err != nil {
				return err
			}

			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {
			changeLocation(client.Must().Kafka, viper.GetString(core.GetFlagName(c.NS, constants.FlagLocation)))
			clusterID := viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))
			r, _, err := client.Must().Kafka.ClustersApi.ClustersFindById(context.Background(),
				clusterID,
			).Execute()
			if err != nil {
				return err
			}

			cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
			out, err := jsontabwriter.GenerateOutput("", jsonpaths.KafkaCluster, r,
				tabheaders.GetHeadersAllDefault(defaultCols, cols))
			if err != nil {
				return err
			}

			fmt.Fprintf(c.Command.Command.OutOrStdout(), out)
			return nil
		},
		InitClient: true,
	})
	cmd.AddStringFlag(constants.FlagClusterId, constants.FlagIdShort, "", "The ID of the cluster you want to retrieve", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagClusterId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		changeLocation(client.Must().Kafka, viper.GetString(core.GetFlagName(cmd.NS, constants.FlagLocation)))
		return completer.ClustersProperty(func(r kafka.ClusterRead) string {
			return *r.Id
		}), cobra.ShellCompDirectiveNoFileComp
	})

	cmd.AddStringFlag(constants.FlagLocation, "", "", "The datacenter location", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagLocation, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		var locations []string
		for k := range locationToURL {
			locations = append(locations, k)
		}

		return locations, cobra.ShellCompDirectiveNoFileComp
	})

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false
	return cmd
}
