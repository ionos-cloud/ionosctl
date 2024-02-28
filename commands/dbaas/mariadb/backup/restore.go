package backup

import (
	"context"

	ionoscloud "github.com/avirtopeanu-ionos/alpha-sdk-go-dbaas-mariadb"
	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/mariadb/cluster"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Restore() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "dbaas-mariadb",
		Resource:  "backup",
		Verb:      "restore",
		Aliases:   []string{"r", "rs"},
		ShortDesc: "Use a MariaDB Backup to restore the cluster to its previous state",
		Example:   "ionosctl dbaas mariadb backup restore --cluster-id CLUSTER_ID --backup-id BACKUP_ID",
		PreCmdRun: core.NoPreRun,
		CmdRun: func(c *core.CommandConfig) error {
			clusterId := viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))

			_, err := client.Must().MariaClient.BackupsApi.ClusterBackupsPost(context.Background(), clusterId).Execute()
			if err != nil {
				return err
			}

			return nil
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagClusterId, constants.FlagIdShort, "", "The unique ID of the cluster", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return cluster.ClustersProperty(func(c ionoscloud.ClusterResponse) string {
			if c.Id == nil {
				return ""
			}
			return *c.Id
		}), cobra.ShellCompDirectiveNoFileComp
	})

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "Response filter to list only the MariaDB Clusters that contain the specified name in the DisplayName field. The value is case insensitive")
	cmd.AddInt32Flag(constants.FlagMaxResults, constants.FlagMaxResultsShort, 0, constants.DescMaxResults)
	cmd.AddInt32Flag(constants.FlagOffset, "", 0, "Skip a certain number of results")

	return cmd
}
