package cluster

import (
	"context"
	"fmt"

	cloudapiv6completer "github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/sdk-go-bundle/products/kafka/v2"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/cobra"
)

func Create() *core.Command {
	cmd := core.NewCommand(
		context.Background(), nil, core.CommandBuilder{
			Namespace: "kafka",
			Resource:  "cluster",
			Verb:      "create",
			Aliases:   []string{"c", "post"},
			ShortDesc: "Create a kafka cluster. Wiki: https://docs.ionos.com/cloud/data-analytics/kafka/api-howtos/create-kafka",
			Example:   "ionosctl kafka cl create --name my-cluster --version 3.9.0 --size XS --location de/txl --datacenter-id DATACENTER_ID --lan-id LAN_ID --broker-addresses 127.0.0.1/24,127.0.0.2/24,127.0.0.3/24",
			PreCmdRun: func(c *core.PreCommandConfig) error {
				if err := core.CheckRequiredFlags(
					c.Command, c.NS,
					constants.FlagName, constants.FlagVersion, constants.FlagSize, constants.FlagLocation,
					constants.FlagDatacenterId, constants.FlagLanId, constants.FlagKafkaBrokerAddresses,
				); err != nil {
					return err
				}

				return nil
			},
			CmdRun: func(c *core.CommandConfig) error {
				input, err := setPropertiesFromFlags(c)
				if err != nil {
					return err
				}

				res, _, err := client.Must().Kafka.ClustersApi.ClustersPost(context.Background()).
					ClusterCreate(
						kafka.ClusterCreate{
							Properties: input,
						},
					).Execute()
				if err != nil {
					return err
				}

				return printCluster(c, res)
			},
			InitClient: true,
		},
	)

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false
	return addClusterCreateFlags(cmd)
}

func addClusterCreateFlags(cmd *core.Command) *core.Command {
	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "The name of the kafka cluster", core.RequiredFlagOption())
	cmd.AddStringFlag(constants.FlagVersion, "", "", "The version of the kafka cluster",
		core.RequiredFlagOption(), core.WithCompletion(func() []string {
			return []string{"3.9.0"}
		}, "", []string{}),
	)

	cmd.AddSetFlag(
		constants.FlagSize, "", "", []string{"XS", "S", "M", "L", "XL"}, "The size of the kafka cluster",
		core.RequiredFlagOption(),
	)

	cmd.AddStringFlag(
		constants.FlagDatacenterId, "", "", "The ID of the datacenter", core.RequiredFlagOption(),
	)
	_ = cmd.Command.RegisterFlagCompletionFunc(
		constants.FlagDatacenterId,
		func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return cloudapiv6completer.DataCentersIds(
				func(datacenter ionoscloud.Datacenter) bool {
					location, _ := cmd.Command.Flags().GetString(constants.FlagLocation)
					return *datacenter.Properties.Location == location
				},
			), cobra.ShellCompDirectiveNoFileComp
		},
	)

	cmd.AddStringFlag(constants.FlagLanId, "", "", "The ID of the LAN", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(
		constants.FlagLanId,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return cloudapiv6completer.LansIds(cmd.Flag(constants.FlagDatacenterId).Value.String()), cobra.ShellCompDirectiveNoFileComp
		},
	)

	cmd.AddStringSliceFlag(
		constants.FlagKafkaBrokerAddresses, "", []string{}, "The list of broker addresses", core.RequiredFlagOption(),
	)
	return cmd
}

func setPropertiesFromFlags(c *core.CommandConfig) (kafka.Cluster, error) {
	p := kafka.Cluster{}

	name, err := c.Command.Command.Flags().GetString(constants.FlagName)
	if err != nil {
		return p, err
	}
	p.Name = name

	version, err := c.Command.Command.Flags().GetString(constants.FlagVersion)
	if err != nil {
		return p, err
	}
	p.Version = version

	size, err := c.Command.Command.Flags().GetString(constants.FlagSize)
	if err != nil {
		return p, err
	}
	p.Size = size

	dcID, err := c.Command.Command.Flags().GetString(constants.FlagDatacenterId)
	if err != nil {
		return p, err
	}

	lanID, err := c.Command.Command.Flags().GetString(constants.FlagLanId)
	if err != nil {
		return p, err
	}

	brokerAddresses, err := c.Command.Command.Flags().GetStringSlice(constants.FlagKafkaBrokerAddresses)
	if err != nil {
		return p, err
	}

	p.Connections = []kafka.KafkaClusterConnection{
		{
			DatacenterId:    dcID,
			LanId:           lanID,
			BrokerAddresses: brokerAddresses,
		},
	}

	return p, nil
}

func printCluster(c *core.CommandConfig, d kafka.ClusterRead) error {
	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
	out, err := jsontabwriter.GenerateOutput(
		"", jsonpaths.KafkaCluster, d,
		tabheaders.GetHeadersAllDefault(allCols, cols),
	)
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)
	return nil
}
