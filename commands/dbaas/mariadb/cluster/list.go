package cluster

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/sdk-go-bundle/products/dbaas/mariadb/v2"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
	"github.com/spf13/viper"
)

func List() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "dbaas-mariadb",
		Resource:  "cluster",
		Verb:      "list",
		Aliases:   []string{"l", "ls"},
		ShortDesc: "List MariaDB Clusters",
		LongDesc:  "Use this command to retrieve a list of MariaDB Clusters provisioned under your account. You can filter the result based on Cluster Name using `--name` option.",
		Example:   "ionosctl dbaas mariadb cluster list",
		PreCmdRun: core.NoPreRun,
		CmdRun: func(c *core.CommandConfig) error {
			c.Verbose("Getting Clusters...")

			return c.ListAllLocations(allCols, func(cfg *shared.Configuration) (any, error) {
				apiClient := mariadb.NewAPIClient(cfg)
				req := apiClient.ClustersApi.ClustersGet(context.Background())

				if fn := core.GetFlagName(c.NS, constants.FlagName); viper.IsSet(fn) {
					req = req.FilterName(viper.GetString(fn))
				}

				clusters, _, err := req.Execute()
				return clusters, err
			})
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "Response filter to list only the MariaDB Clusters that contain the specified name in the DisplayName field. The value is case insensitive")

	return cmd
}
