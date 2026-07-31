package cluster

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/sdk-go-bundle/products/dbaas/mariadb/v2"
	"github.com/spf13/viper"
)

func Get() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "dbaas-mariadb",
		Resource:  "cluster",
		Verb:      "get",
		Aliases:   []string{"g"},
		ShortDesc: "Get a MariaDB Cluster by ID",
		LongDesc:  "Retrieve the full details of a single MariaDB cluster by its ID, including its current state, MariaDB version, instance count, compute/storage sizing, DNS name, connection and maintenance window. Use this to check whether a cluster has reached the AVAILABLE state before connecting to it or running further operations.",
		Example:   "ionosctl dbaas mariadb cluster get --cluster-id <cluster-id>",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return c.CheckRequiredFlagsAndLocation(constants.FlagClusterId)
		},
		CmdRun: func(c *core.CommandConfig) error {
			clusterId := viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))

			c.Verbose("Getting Cluster by id: %s", clusterId)

			cluster, _, err := client.Must().MariaClient.ClustersApi.ClustersFindById(context.Background(), clusterId).Execute()
			if err != nil {
				return err
			}

			return c.Printer(allCols).Print(cluster)
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagClusterId, constants.FlagIdShort, "", "The unique ID of the cluster to retrieve",
		core.RequiredFlagOption(),
		core.WithCompletion(
			func() []string {
				return ClustersProperty(func(c mariadb.ClusterResponse) string {
					if c.Id == nil {
						return ""
					}
					return *c.Id
				})
			}, constants.MariaDBApiRegionalURL, constants.MariaDBLocations),
	)

	cmd.Command.SilenceUsage = true

	return cmd
}
