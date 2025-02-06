package cluster

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
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
		Example:   "ionosctl dbaas mariadb cluster get --cluster-id <cluster-id>",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return c.Command.Command.MarkFlagRequired(constants.FlagClusterId)
		},
		CmdRun: func(c *core.CommandConfig) error {
			clusterId := viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))

			fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Getting Cluster by id: %s", clusterId))

			cluster, _, err := client.Must().MariaClient.ClustersApi.ClustersFindById(context.Background(), clusterId).Execute()
			if err != nil {
				return err
			}

			cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

			out, err := jsontabwriter.GenerateOutput("", jsonpaths.DbaasMariadbCluster, cluster,
				tabheaders.GetHeaders(allCols, defaultCols, cols))
			if err != nil {
				return err
			}

			fmt.Fprintf(c.Command.Command.OutOrStdout(), out)
			return nil
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagClusterId, constants.FlagIdShort, "", "The unique ID of the cluster",
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
