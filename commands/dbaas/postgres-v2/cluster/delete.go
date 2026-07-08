package cluster

import (
	"context"
	"errors"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/postgres-v2/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/confirm"
	"github.com/ionos-cloud/ionosctl/v6/pkg/functional"
	psqlv2 "github.com/ionos-cloud/sdk-go-bundle/products/dbaas/psql/v3"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
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
		LongDesc: `Use this command to delete a specified PostgreSQL Cluster from your account. Use ` + "`--wait` (`-w`)" + ` to wait for the deletion to complete.

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

	if err := c.RequireExplicitLocation(); err != nil {
		return err
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
	return nil
}

func ClusterDeleteAll(c *core.CommandConfig) error {
	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagName)) {
		c.Verbose("Filtering based on Cluster Name: %v", viper.GetString(core.GetFlagName(c.NS, constants.FlagName)))
	}

	return c.RunForAllLocations(func(cfg *shared.Configuration, location string) error {
		apiClient := psqlv2.NewAPIClient(cfg)

		req := apiClient.ClustersApi.ClustersGet(context.Background())
		if fn := core.GetFlagName(c.NS, constants.FlagName); viper.IsSet(fn) && viper.GetString(fn) != "" {
			req = req.FilterName(viper.GetString(fn))
		}
		if fn := core.GetFlagName(c.NS, constants.FlagState); viper.IsSet(fn) {
			req = req.FilterState(psqlv2.PostgresClusterStates(viper.GetString(fn)))
		}
		clusters, _, err := req.Execute()
		if err != nil {
			return fmt.Errorf("failed listing clusters in location %s: %w", location, err)
		}

		return functional.ApplyAndAggregateErrors(clusters.GetItems(), func(cluster psqlv2.ClusterRead) error {
			// Skip (not fail) on decline, matching the other delete --all
			// commands: a "no" for one cluster must not abort the whole run.
			if !confirm.FAsk(c.Command.Command.InOrStdin(),
				fmt.Sprintf("delete cluster %s (%s) (location: %s)", cluster.Id, cluster.Properties.Name, location),
				viper.GetBool(constants.ArgForce)) {
				return nil
			}

			c.Verbose("Deleting cluster: %s (%s)", cluster.Id, cluster.Properties.Name)
			_, delErr := apiClient.ClustersApi.ClustersDelete(context.Background(), cluster.Id).Execute()
			if delErr != nil {
				return fmt.Errorf("failed deleting cluster %s (%s): %w", cluster.Id, cluster.Properties.Name, delErr)
			}

			return nil
		})
	})
}
