package cluster

import (
	"context"
	"fmt"
	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/completer"
	cloudapiv6completer "github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/ionosctl/v6/pkg/pointer"
	kafka "github.com/ionos-cloud/sdk-go-kafka"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Create() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "kafka",
		Resource:  "cluster",
		Verb:      "create",
		Aliases:   []string{"c", "post"},
		ShortDesc: "Create a kafka cluster. Wiki: https://docs.ionos.com/cloud/data-analytics/kafka/api-howtos/create-kafka",
		Example:   "ionosctl kafka cl create --name my-cluster --version 1.0 --size S --location de/txl --datacenter-id DATACENTER_ID --lan-id LAN_ID --broker-addresses 127.0.0.1,127.0.0.2",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			if err := core.CheckRequiredFlags(c.Command, c.NS,
				constants.FlagName, constants.FlagVersion, constants.FlagSize, constants.FlagLocation,
				constants.FlagDatacenterId, constants.FlagLanId, constants.FlagKafkaBrokerAddresses); err != nil {
				return err
			}

			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {
			input := &kafka.Cluster{}
			if err := setPropertiesFromFlags(c, input); err != nil {
				return err
			}

			changeLocation(client.Must().Kafka, viper.GetString(core.GetFlagName(c.NS, constants.FlagLocation)))
			res, _, err := client.Must().Kafka.ClustersApi.ClustersPost(context.Background()).
				ClusterCreate(kafka.ClusterCreate{
					Properties: input,
				}).Execute()
			if err != nil {
				return err
			}

			return printCluster(c, res)
		},
		InitClient: true,
	})

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false
	return addClusterCreateFlags(cmd)
}

func changeLocation(client *kafka.APIClient, location string) {
	cfg := client.GetConfig()
	cfg.Servers = kafka.ServerConfigurations{
		{
			URL: locationToURL[location],
		},
	}
	return
}

func addClusterCreateFlags(cmd *core.Command) *core.Command {
	cmd.AddStringFlag(constants.FlagName, "", "", "The name of the kafka cluster", core.RequiredFlagOption())
	cmd.AddStringFlag(constants.FlagVersion, "", "", "The version of the kafka cluster", core.RequiredFlagOption())
	cmd.AddSetFlag(constants.FlagSize, "", "", []string{"XS", "S", "M", "L", "XL"}, "The size of the kafka cluster", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagSize, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"XS", "S", "M", "L", "XL"}, cobra.ShellCompDirectiveNoFileComp
	})

	cmd.AddStringFlag(constants.FlagLocation, "", "", "The location of the kafka cluster", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagLocation, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		var locations []string
		for k := range locationToURL {
			locations = append(locations, k)
		}

		return locations, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(constants.FlagDatacenterId, "", "", "The ID of the datacenter", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagDatacenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})

	cmd.AddStringFlag(constants.FlagLanId, "", "", "The ID of the LAN", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagLanId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return cloudapiv6completer.LansIds(cmd.Flag(constants.FlagDatacenterId).Value.String()), cobra.ShellCompDirectiveNoFileComp
	})

	cmd.AddStringSliceFlag(constants.FlagKafkaBrokerAddresses, "", []string{}, "The list of broker addresses", core.RequiredFlagOption())
	return cmd
}

func setPropertiesFromFlags(c *core.CommandConfig, p *kafka.Cluster) error {
	p.Name = pointer.From(viper.GetString(core.GetFlagName(c.NS, constants.FlagName)))
	p.Version = pointer.From(viper.GetString(core.GetFlagName(c.NS, constants.FlagVersion)))
	p.Size = pointer.From(viper.GetString(core.GetFlagName(c.NS, constants.FlagSize)))

	p.Connections = &[]kafka.KafkaClusterConnection{{
		DatacenterId:    pointer.From(viper.GetString(core.GetFlagName(c.NS, constants.FlagDatacenterId))),
		LanId:           pointer.From(viper.GetString(core.GetFlagName(c.NS, constants.FlagLanId))),
		BrokerAddresses: pointer.From(viper.GetStringSlice(core.GetFlagName(c.NS, constants.FlagKafkaBrokerAddresses))),
	}}

	return nil
}

func printCluster(c *core.CommandConfig, d kafka.ClusterRead) error {
	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
	out, err := jsontabwriter.GenerateOutput("", jsonpaths.KafkaCluster, d,
		tabheaders.GetHeadersAllDefault(defaultCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)
	return nil
}

var locationToURL = map[string]string{
	"de/fra": "https://kafka.de-fra.ionos.com",
	"de/txl": "https://kafka.de-txl.ionos.com",
	// other locations not yet available. will be added in the future.
	// "es/vit": "https://kafka.es-vit.ionos.com",
	// "gb/lhr": "https://kafka.gb-lhr.ionos.com",
	// "us/ewr": "https://kafka.us-ewr.ionos.com",
	// "us/las": "https://kafka.us-las.ionos.com",
	// "us/mci": "https://kafka.us-mci.ionos.com",
	// "fr/par": "https://kafka.fr-par.ionos.com",
}
