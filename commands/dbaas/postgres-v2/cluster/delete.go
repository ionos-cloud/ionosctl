package cluster

import (
	"context"
	"errors"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/postgres-v2/waiter"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/ionosctl/v6/internal/waitfor"
	"github.com/ionos-cloud/ionosctl/v6/pkg/confirm"
	"github.com/ionos-cloud/ionosctl/v6/pkg/functional"
	psqlv2 "github.com/ionos-cloud/sdk-go-dbaas-psql"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ClusterDeleteCmd() *core.Command {
	ctx := context.TODO()
	deleteCmd := core.NewCommand(ctx, nil, core.CommandBuilder{
		Namespace: "dbaas-postgres-v2",
		Resource:  "cluster",
		Verb:      "delete",
		Aliases:   []string{"d"},
		ShortDesc: "Delete a PostgreSQL Cluster",
		LongDesc: `Use this command to delete a specified PostgreSQL Cluster from your account. You can wait for the cluster to be deleted with the wait-for-deletion option.

Required values to run command:

* Cluster Id`,
		Example:    "ionosctl dbaas postgres cluster delete --cluster-id <cluster-id>",
		PreCmdRun:  PreRunClusterDelete,
		CmdRun:     RunClusterDelete,
		InitClient: true,
	})
	deleteCmd.AddUUIDFlag(constants.FlagClusterId, constants.FlagIdShort, "", constants.DescCluster, core.RequiredFlagOption(),
		core.WithCompletion(func() []string {
			clusters, err := Clusters()
			if err != nil {
				return []string{}
			}
			return functional.Map(clusters.Items, func(c psqlv2.ClusterRead) string {
				return fmt.Sprintf("%s\t%s: %d instances, datacenter: %s",
					c.Id, c.Properties.Name, c.Properties.Instances.Count, c.Properties.Connection.DatacenterId)

			})
		}, constants.PostgresApiRegionalURL, constants.PostgresLocations),
	)
	deleteCmd.AddBoolFlag(constants.ArgAll, constants.ArgAllShort, false, "Delete all Clusters")
	deleteCmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "Delete all Clusters after filtering based on name. It does not require an exact match. Can be used with --all flag")
	deleteCmd.AddBoolFlag(constants.ArgWaitForDelete, constants.ArgWaitForStateShort, constants.DefaultWait, "Wait for Cluster to be completely removed")
	deleteCmd.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, constants.DefaultClusterTimeout, "Timeout option for Cluster to be completely removed[seconds]")
	deleteCmd.AddStringSliceFlag(constants.ArgCols, "", defaultClusterCols, tabheaders.ColsMessage(allClusterCols))
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allClusterCols, cobra.ShellCompDirectiveNoFileComp
	})

	return deleteCmd
}

func PreRunClusterDelete(c *core.PreCommandConfig) error {
	err := core.CheckRequiredFlagsSets(c.Command, c.NS, []string{constants.FlagClusterId}, []string{constants.ArgAll})
	if err != nil {
		return err
	}
	// Validate Flags
	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagName)) {
		if !viper.IsSet(core.GetFlagName(c.NS, constants.ArgAll)) {
			return errors.New("error: name flag can to be used with the --all flag")
		}
	}
	return nil
}

func RunClusterDelete(c *core.CommandConfig) error {
	if viper.GetBool(core.GetFlagName(c.NS, constants.ArgAll)) {
		if err := ClusterDeleteAll(c); err != nil {
			return err
		}
		return nil
	}

	clusterId := viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.ClusterId, clusterId))

	if !confirm.FAsk(c.Command.Command.InOrStdin(), fmt.Sprintf("delete cluster with id: %v", clusterId), viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	_, err := client.Must().PostgresClientV2.ClustersApi.ClustersDelete(context.Background(), clusterId).Execute()
	if err != nil {
		return err
	}
	if viper.GetBool(core.GetFlagName(c.NS, constants.ArgWaitForDelete)) {
		if err = waitfor.WaitForDelete(c, waiter.ClusterDeleteInterrogator, viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))); err != nil {
			return err
		}
	}
	return nil
}

func ClusterDeleteAll(c *core.CommandConfig) error {
	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagName)) {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Filtering based on Cluster Name: %v", viper.GetString(core.GetFlagName(c.NS, constants.FlagName))))
	}

	req := client.Must().PostgresClientV2.ClustersApi.ClustersGet(context.Background())
	if fn := core.GetFlagName(c.NS, constants.FlagName); viper.IsSet(fn) && viper.GetString(fn) != "" {
		req = req.FilterName(viper.GetString(fn))
	}
	clusters, _, err := req.Execute()
	if err != nil {
		return err
	}

	dataOk, ok := clusters.GetItemsOk()
	if !ok || dataOk == nil {
		return fmt.Errorf("could not get items of Clusters")
	}

	if len(dataOk) <= 0 {
		return fmt.Errorf("no Clusters found")
	}

	var multiErr error
	for _, cluster := range dataOk {
		idOk, ok := cluster.GetIdOk()
		if !ok || idOk == nil {
			continue
		}

		clusterId := *idOk
		propertiesOk, ok := cluster.GetPropertiesOk()
		clusterName := "Unnamed Cluster"
		if ok && propertiesOk != nil {
			if n, ok := propertiesOk.GetNameOk(); ok && n != nil {
				clusterName = *n
			}
		}

		if !confirm.FAsk(c.Command.Command.InOrStdin(), fmt.Sprintf("delete cluster with id: %v, name: %v", clusterId, clusterName), viper.GetBool(constants.ArgForce)) {
			multiErr = errors.Join(multiErr, fmt.Errorf(confirm.UserDenied))
			continue
		}

		_, err := client.Must().PostgresClientV2.ClustersApi.ClustersDelete(context.Background(), *idOk).Execute()
		if err != nil {
			multiErr = errors.Join(multiErr, fmt.Errorf(constants.ErrDeleteAll, c.Resource, *idOk, err))
			continue
		}

		if viper.GetBool(core.GetFlagName(c.NS, constants.ArgWaitForDelete)) {
			if err = waitfor.WaitForDelete(c, waiter.ClusterDeleteInterrogator, *idOk); err != nil {
				multiErr = errors.Join(multiErr, fmt.Errorf(constants.ErrWaitDeleteAll, c.Resource, *idOk, err))
			}
		}
	}

	if multiErr != nil {
		return multiErr
	}

	return nil
}
