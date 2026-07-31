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
			Use:     "cluster",
			Aliases: []string{"c"},
			Short:   "MariaDB Cluster Operations",
			Long: `Manage DBaaS MariaDB clusters.

A MariaDB cluster is a managed, regional group of instances running one MariaDB version: a single standalone instance, or a high-availability replica set of one primary plus secondaries (an odd instance count of 3 or 5). Every instance shares the same compute and storage sizing (cores, RAM, storage). A cluster lives in one location (region) and is reachable only over a private LAN in one of your Virtual Data Centers - never the public internet.

IONOS continuously backs each cluster up (see the ` + "`backup`" + ` sub-commands), enabling point-in-time restore between a backup's earliest recovery time and now. Clusters progress through states such as CREATING, BUSY, AVAILABLE and UPDATING; most operations require the cluster to be AVAILABLE first.

Use these sub-commands to list, create, inspect, update and delete clusters.`,
			TraverseChildren: true,
		},
	}

	cmd.AddColsFlag(allCols)

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
