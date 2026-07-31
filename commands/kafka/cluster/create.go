package cluster

import (
	"context"

	cloudapiv6completer "github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/sdk-go-bundle/products/kafka/v2"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Create() *core.Command {
	cmd := core.NewCommand(
		context.Background(), nil, core.CommandBuilder{
			Namespace: "kafka",
			Resource:  "cluster",
			Verb:      "create",
			Aliases:   []string{"c", "post"},
			ShortDesc: "Create a Kafka cluster",
			LongDesc: `Create a Kafka cluster: 3 broker nodes attached to one of your private LANs.

Sizing (--size). Every size has 3 brokers; sizes differ in the cores/RAM/storage each broker gets:

  XS   1 core   2 GB    195 GB/broker   (585 GB total)   dev / test
  S    2 cores  4 GB    250 GB/broker   (750 GB total)
  M    2 cores  8 GB    400 GB/broker   (1200 GB total)
  L    4 cores  16 GB   800 GB/broker   (2400 GB total)
  XL   8 cores  32 GB   1500 GB/broker  (4500 GB total)  production

Storage is shared by all topics on the cluster; size for (data-rate x retention x replication-factor) across your topics.

Networking. The cluster lives in the LAN --lan-id inside the datacenter --datacenter-id, both in --location. --broker-addresses assigns a private IP to each broker in CIDR notation (e.g. 10.0.0.1/24) taken from that LAN's subnet, so pass exactly 3 addresses. Clients connect over TLS with mutual authentication; port 9093 is appended automatically (e.g. 10.0.0.1:9093).

The cluster is BUSY while it deploys and becomes AVAILABLE when ready; topics and users can only be created once it is AVAILABLE.

Wiki: https://docs.ionos.com/cloud/data-analytics/kafka/api-howtos/create-kafka`,
			Example: "ionosctl kafka cl create --name my-cluster --version 3.9.0 --size XS --location de/txl --datacenter-id DATACENTER_ID --lan-id LAN_ID --broker-addresses 10.0.0.1/24,10.0.0.2/24,10.0.0.3/24",
			PreCmdRun: func(c *core.PreCommandConfig) error {
				return c.CheckRequiredFlagsAndLocation(
					constants.FlagName, constants.FlagVersion, constants.FlagSize,
					constants.FlagDatacenterId, constants.FlagLanId, constants.FlagKafkaBrokerAddresses,
				)
			},
			CmdRun: func(c *core.CommandConfig) error {
				input := setPropertiesFromFlags(c)

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
	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "Human-readable name for the cluster", core.RequiredFlagOption())
	cmd.AddStringFlag(constants.FlagVersion, "", "", "Kafka version to deploy, e.g. 3.9.0",
		core.RequiredFlagOption(), core.WithCompletion(func() []string {
			return []string{"3.9.0"}
		}, "", []string{}),
	)

	cmd.AddSetFlag(
		constants.FlagSize, "", "", []string{"XS", "S", "M", "L", "XL"},
		"Cluster size: sets cores/RAM/storage per broker (all sizes have 3 brokers). XS=1c/2GB, S=2c/4GB, M=2c/8GB, L=4c/16GB, XL=8c/32GB",
		core.RequiredFlagOption(),
	)

	cmd.AddStringFlag(
		constants.FlagDatacenterId, "", "", "ID of the Virtual Data Center holding the LAN the brokers attach to; must be in --location", core.RequiredFlagOption(),
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

	cmd.AddStringFlag(constants.FlagLanId, "", "", "ID of the private LAN (inside --datacenter-id) the brokers attach to; clients reach the cluster over this LAN", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(
		constants.FlagLanId,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return cloudapiv6completer.LansIds(cmd.Flag(constants.FlagDatacenterId).Value.String()), cobra.ShellCompDirectiveNoFileComp
		},
	)

	cmd.AddStringSliceFlag(
		constants.FlagKafkaBrokerAddresses, "", []string{},
		"Private IP per broker in CIDR notation (e.g. 10.0.0.1/24), one per broker — pass exactly 3, all from the --lan-id subnet. Port 9093 (TLS) is appended for clients",
		core.RequiredFlagOption(),
	)
	return cmd
}

func setPropertiesFromFlags(c *core.CommandConfig) kafka.Cluster {
	p := kafka.Cluster{}

	p.Name = viper.GetString(core.GetFlagName(c.NS, constants.FlagName))
	p.Version = viper.GetString(core.GetFlagName(c.NS, constants.FlagVersion))
	p.Size = viper.GetString(core.GetFlagName(c.NS, constants.FlagSize))

	p.Connections = []kafka.KafkaClusterConnection{
		{
			DatacenterId: viper.GetString(core.GetFlagName(c.NS, constants.FlagDatacenterId)),
			LanId:        viper.GetString(core.GetFlagName(c.NS, constants.FlagLanId)),
			BrokerAddresses: viper.GetStringSlice(
				core.GetFlagName(c.NS, constants.FlagKafkaBrokerAddresses),
			),
		},
	}

	return p
}

func printCluster(c *core.CommandConfig, d kafka.ClusterRead) error {
	return c.Printer(allCols).Print(d)
}
