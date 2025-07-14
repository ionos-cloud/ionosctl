package backup

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/mariadb/cluster"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/ionosctl/v6/pkg/functional"
	"github.com/ionos-cloud/sdk-go-bundle/products/dbaas/mariadb/v2"

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

	cmd.Command.PersistentFlags().StringSlice(constants.FlagCols, defaultCols, tabheaders.ColsMessage(allCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allCols, cobra.ShellCompDirectiveNoFileComp
	})

	cmd.AddCommand(List())
	cmd.AddCommand(Get())

	return cmd

}

var (
	allCols = []string{"BackupId", "ClusterId", "Size", "Items"}

	defaultCols = allCols
)

func Backups(fs ...Filter) (mariadb.BackupList, error) {
	cs, err := cluster.Clusters()

	if len(cs.Items) == 0 {
		return mariadb.BackupList{}, nil // no clusters -> empty response with no error
	}

	req := client.Must().MariaClient.BackupsApi.BackupsGet(context.Background())
	for _, f := range fs {
		req = f(req)
	}

	bs, _, err := req.Execute()

	return bs, err
}

type Filter func(request mariadb.ApiBackupsGetRequest) mariadb.ApiBackupsGetRequest

func FilterPaginationFlags(c *core.CommandConfig) Filter {
	return func(req mariadb.ApiBackupsGetRequest) mariadb.ApiBackupsGetRequest {
		if f := core.GetFlagName(c.NS, constants.FlagMaxResults); viper.IsSet(f) {
			req = req.Limit(viper.GetInt32(f))
		}
		if f := core.GetFlagName(c.NS, constants.FlagOffset); viper.IsSet(f) {
			req = req.Offset(viper.GetInt32(f))
		}
		return req
	}
}

func BackupsProperty[V any](f func(c mariadb.BackupResponse) V, fs ...Filter) []V {
	recs, err := Backups(fs...)
	if err != nil {
		return nil
	}
	return functional.Map(recs.Items, f)
}
