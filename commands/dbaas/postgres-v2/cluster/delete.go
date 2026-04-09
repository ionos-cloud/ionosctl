package cluster

import (
	"context"
	"errors"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/postgres-v2/completer"
	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/postgres-v2/waiter"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/waitfor"
	"github.com/ionos-cloud/ionosctl/v6/pkg/confirm"
	"github.com/ionos-cloud/ionosctl/v6/pkg/functional"
	psqlv2 "github.com/ionos-cloud/sdk-go-bundle/products/dbaas/psql/v3"
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
		Example:    "ionosctl dbaas postgres-v2 cluster delete --cluster-id <cluster-id>",
		PreCmdRun:  PreRunClusterDelete,
		CmdRun:     RunClusterDelete,
		InitClient: true,
	})
	deleteCmd.AddUUIDFlag(constants.FlagClusterId, constants.FlagIdShort, "", constants.DescCluster, core.RequiredFlagOption(),
		core.WithCompletion(completer.ClusterIds, constants.PostgresApiRegionalURL, constants.PostgresLocations),
	)
	deleteCmd.AddBoolFlag(constants.ArgAll, constants.ArgAllShort, false, "Delete all Clusters")
	deleteCmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "Delete all Clusters after filtering based on name. It does not require an exact match. Can be used with --all flag")
	deleteCmd.AddSetFlag(constants.FlagState, "", "", []string{"PROVISIONING", "AVAILABLE", "UPDATING", "DESTROYING", "FAILED"},
		"When used with --all, only delete clusters in this state")
	deleteCmd.AddBoolFlag(constants.ArgWaitForDelete, constants.ArgWaitForStateShort, constants.DefaultWait, "Wait for Cluster to be completely removed")
	deleteCmd.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, constants.DefaultClusterTimeout, "Timeout option for Cluster to be completely removed[seconds]")
	return deleteCmd
}

func PreRunClusterDelete(c *core.PreCommandConfig) error {
	err := core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{constants.FlagClusterId},
		[]string{constants.ArgAll},
		[]string{constants.ArgAll, constants.FlagName},
		[]string{constants.ArgAll, constants.FlagState},
	)
	if err != nil {
		return err
	}
	// Validate Flags
	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagName)) && !viper.IsSet(core.GetFlagName(c.NS, constants.ArgAll)) {
		return errors.New("error: --name flag can only be used with the --all flag")
	}
	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagState)) && !viper.IsSet(core.GetFlagName(c.NS, constants.ArgAll)) {
		return errors.New("error: --state flag can only be used with the --all flag")
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

	c.Verbose(constants.ClusterId, clusterId)

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
		c.Verbose("Filtering based on Cluster Name: %v", viper.GetString(core.GetFlagName(c.NS, constants.FlagName)))
	}

	req := client.Must().PostgresClientV2.ClustersApi.ClustersGet(context.Background())
	if fn := core.GetFlagName(c.NS, constants.FlagName); viper.IsSet(fn) && viper.GetString(fn) != "" {
		req = req.FilterName(viper.GetString(fn))
	}
	if fn := core.GetFlagName(c.NS, constants.FlagState); viper.IsSet(fn) {
		req = req.FilterState(psqlv2.PostgresClusterStates(viper.GetString(fn)))
	}
	clusters, _, err := req.Execute()
	if err != nil {
		return err
	}

	items := clusters.GetItems()
	if len(items) == 0 {
		return fmt.Errorf("no Clusters found")
	}

	return functional.ApplyAndAggregateErrors(items, func(cluster psqlv2.ClusterRead) error {
		if !confirm.FAsk(c.Command.Command.InOrStdin(),
			fmt.Sprintf("delete cluster %s (%s)", cluster.Id, cluster.Properties.Name),
			viper.GetBool(constants.ArgForce)) {
			return fmt.Errorf(confirm.UserDenied)
		}

		c.Verbose("Deleting cluster: %s (%s)", cluster.Id, cluster.Properties.Name)
		_, delErr := client.Must().PostgresClientV2.ClustersApi.ClustersDelete(context.Background(), cluster.Id).Execute()
		if delErr != nil {
			return fmt.Errorf("failed deleting cluster %s (%s): %w", cluster.Id, cluster.Properties.Name, delErr)
		}

		if viper.GetBool(core.GetFlagName(c.NS, constants.ArgWaitForDelete)) {
			if waitErr := waitfor.WaitForDelete(c, waiter.ClusterDeleteInterrogator, cluster.Id); waitErr != nil {
				return fmt.Errorf("failed waiting for deletion of cluster %s (%s): %w", cluster.Id, cluster.Properties.Name, waitErr)
			}
		}

		return nil
	})
}
