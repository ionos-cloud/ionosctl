package cluster

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/ionos-cloud/ionosctl/v6/pkg/functional"
	"github.com/ionos-cloud/sdk-go-bundle/products/dbaas/mariadb/v2"
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

	cmd.Command.PersistentFlags().StringSlice(constants.ArgCols, nil, table.ColsMessage(allCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return table.AllCols(allCols), cobra.ShellCompDirectiveNoFileComp
	})

	cmd.AddCommand(List())
	cmd.AddCommand(Create())
	cmd.AddCommand(Update())
	cmd.AddCommand(Get())
	cmd.AddCommand(Delete())

	return cmd
}

var allCols = []table.Column{
	{Name: "ClusterId", JSONPath: "id", Default: true},
	{Name: "Name", JSONPath: "properties.displayName", Default: true},
	{Name: "DNS", JSONPath: "properties.dnsName", Default: true},
	{Name: "Instances", JSONPath: "properties.instances", Default: true},
	{Name: "Version", JSONPath: "properties.mariadbVersion", Default: true},
	{Name: "State", JSONPath: "metadata.state", Default: true},
	{Name: "Cores", JSONPath: "properties.cores"},
	{Name: "RAM", Format: func(item map[string]any) any {
		v := table.Navigate(item, "properties.ram")
		if v == nil {
			return nil
		}
		f, ok := v.(float64)
		if !ok {
			return v
		}
		return fmt.Sprintf("%d GB", int(f))
	}},
	{Name: "StorageSize", JSONPath: "properties.storageSize"},
	{Name: "MaintenanceDay", JSONPath: "properties.maintenanceWindow.dayOfTheWeek"},
	{Name: "MaintenanceTime", JSONPath: "properties.maintenanceWindow.time"},
}

func Clusters(fs ...Filter) (mariadb.ClusterList, error) {
	req := client.Must().MariaClient.ClustersApi.ClustersGet(context.Background())

	for _, f := range fs {
		req = f(req)
	}

	clusters, _, err := req.Execute()
	if err != nil {
		return mariadb.ClusterList{}, fmt.Errorf("failed getting clusters: %w", err)
	}
	return clusters, err
}

func ClustersProperty[V any](f func(c mariadb.ClusterResponse) V, fs ...Filter) []V {
	recs, err := Clusters(fs...)
	if err != nil {
		return nil
	}
	return functional.Map(recs.Items, f)
}

type Filter func(mariadb.ApiClustersGetRequest) mariadb.ApiClustersGetRequest

func FilterNameFlags(c *core.CommandConfig) Filter {
	return func(req mariadb.ApiClustersGetRequest) mariadb.ApiClustersGetRequest {
		if f := core.GetFlagName(c.NS, constants.FlagName); viper.IsSet(f) {
			req = req.FilterName(viper.GetString(f))
		}
		return req
	}
}
