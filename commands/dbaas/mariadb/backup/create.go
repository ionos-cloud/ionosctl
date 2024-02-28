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

func Create() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "dbaas-mariadb",
		Resource:  "backup",
		Verb:      "create",
		Aliases:   []string{"c"},
		ShortDesc: "Create an ad-hoc MariaDB Backup",
		Example:   "ionosctl dbaas mariadb backup create --cluster-id CLUSTER_ID",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			err := core.CheckRequiredFlags(c.Command, c.NS, constants.FlagClusterId)
			if err != nil {
				return err
			}

			return nil
		},
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

	return cmd
}
