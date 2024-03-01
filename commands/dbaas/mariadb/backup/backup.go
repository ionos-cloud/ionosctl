package backup

import (
	"context"

	ionoscloud "github.com/avirtopeanu-ionos/alpha-sdk-go-dbaas-mariadb"
	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/mongo/cluster"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Root() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "backup",
			Aliases:          []string{"b"},
			Short:            "MariaDB Backup Operations",
			Long:             "The sub-commands of `ionosctl dbaas mariadb backup` allow you to manage the MariaDB Cluster Backups under your account.",
			TraverseChildren: true,
		},
	}

	cmd.AddCommand(List())
	cmd.AddCommand(Get())

	return cmd

}

var (
	allCols = []string{"BackupId", "ClusterId", "Items"}

	defaultCols = allCols
)

func Backups(fs ...Filter) (ionoscloud.BackupList, error) {
	cs, err := cluster.Clusters()

	if cs.Items == nil || len(*cs.Items) == 0 {
		return ionoscloud.BackupList{}, nil // no clusters -> empty response with no error
	}

	req := client.Must().MariaClient.BackupsApi.BackupsGet(context.Background())
	for _, f := range fs {
		req = f(req)
	}

	bs, _, err := req.Execute()

	return bs, err
}

type Filter func(request ionoscloud.ApiBackupsGetRequest) ionoscloud.ApiBackupsGetRequest

func FilterPaginationFlags(c *core.CommandConfig) Filter {
	return func(req ionoscloud.ApiBackupsGetRequest) ionoscloud.ApiBackupsGetRequest {
		if f := core.GetFlagName(c.NS, constants.FlagMaxResults); viper.IsSet(f) {
			req = req.Limit(viper.GetInt32(f))
		}
		if f := core.GetFlagName(c.NS, constants.FlagOffset); viper.IsSet(f) {
			req = req.Offset(viper.GetInt32(f))
		}
		return req
	}
}
