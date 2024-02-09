package backup

import (
	"fmt"

	ionoscloud "github.com/avirtopeanu-ionos/alpha-sdk-go-dbaas-mariadb"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Root() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "cluster",
			Aliases:          []string{"c"},
			Short:            "MariaDB Cluster Operations",
			Long:             "The sub-commands of `ionosctl dbaas mariadb cluster` allow you to manage the MariaDB Clusters under your account.",
			TraverseChildren: true,
		},
	}

	cmd.AddCommand(List())

}

var (
	allCols = []string{"BackupId"}

	defaultCols = allCols[0:0]
)

func Backups(fs ...Filter) (ionoscloud.ClusterList, error) {
	req := client.Must().MongoClient.BackupsApi.ClustersRestorePost()

	for _, f := range fs {
		req = f(req)
	}

	clusters, _, err := req.Execute()
	if err != nil {
		return ionoscloud.ClusterList{}, fmt.Errorf("failed getting clusters: %w", err)
	}
	return clusters, err
}

type Filter func(ionoscloud.ApiClustersGetRequest) ionoscloud.ApiClustersGetRequest

func FilterPaginationFlags(c *core.CommandConfig) Filter {
	return func(req ionoscloud.ApiClustersGetRequest) ionoscloud.ApiClustersGetRequest {
		if f := core.GetFlagName(c.NS, constants.FlagMaxResults); viper.IsSet(f) {
			req = req.Limit(viper.GetInt32(f))
		}
		if f := core.GetFlagName(c.NS, constants.FlagOffset); viper.IsSet(f) {
			req = req.Offset(viper.GetInt32(f))
		}
		return req
	}
}

func FilterNameFlags(c *core.CommandConfig) Filter {
	return func(req ionoscloud.ApiClustersGetRequest) ionoscloud.ApiClustersGetRequest {
		if f := core.GetFlagName(c.NS, constants.FlagName); viper.IsSet(f) {
			req = req.FilterName(viper.GetString(f))
		}
		return req
	}
}
