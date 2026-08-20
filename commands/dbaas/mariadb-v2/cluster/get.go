package cluster

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/mariadb-v2/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/spf13/viper"
)

func ClusterGetCmd() *core.Command {
	ctx := context.TODO()
	get := core.NewCommand(ctx, nil, core.CommandBuilder{
		Namespace:  "dbaas-mariadb-v2",
		Resource:   "cluster",
		Verb:       "get",
		Aliases:    []string{"g"},
		ShortDesc:  "Get a MariaDB Cluster",
		Example:    "ionosctl dbaas mariadb-v2 cluster get --location <location> --cluster-id <cluster-id>",
		LongDesc:   "Use this command to retrieve details about a MariaDB Cluster by using its ID.\n\nRequired values to run command:\n\n* Cluster Id",
		PreCmdRun:  PreRunClusterId,
		CmdRun:     RunClusterGet,
		InitClient: true,
	})

	get.AddUUIDFlag(constants.FlagClusterId, constants.FlagIdShort, "", constants.DescCluster, core.RequiredFlagOption(),
		core.WithCompletion(completer.ClusterIds, constants.MariaDBApiRegionalURL, constants.MariaDBLocations),
	)
	return get
}

func PreRunClusterId(c *core.PreCommandConfig) error {
	return c.CheckRequiredFlagsAndLocation(constants.FlagClusterId)
}

func RunClusterGet(c *core.CommandConfig) error {
	clusterId := viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))
	c.Verbose(constants.ClusterId, clusterId)
	c.Verbose("Getting Cluster...")

	cluster, _, err := client.Must().MariaClientV2.ClustersApi.ClustersFindById(
		context.Background(), clusterId).Execute()
	if err != nil {
		return fmt.Errorf("could not get cluster: %w", err)
	}

	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

	return c.Out(table.Sprint(clusterCols, cluster, cols))
}
