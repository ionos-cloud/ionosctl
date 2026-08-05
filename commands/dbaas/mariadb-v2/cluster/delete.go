package cluster

import (
	"context"
	"errors"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/mariadb-v2/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/confirm"
	mariadb "github.com/ionos-cloud/sdk-go-bundle/products/dbaas/mariadb/v3"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
	"github.com/spf13/viper"
)

func ClusterDeleteCmd() *core.Command {
	ctx := context.TODO()
	deleteCmd := core.NewCommand(ctx, nil, core.CommandBuilder{
		Namespace: "dbaas-mariadb-v2",
		Resource:  "cluster",
		Verb:      "delete",
		Aliases:   []string{"del", "d"},
		ShortDesc: "Delete a MariaDB Cluster",
		LongDesc: `Use this command to delete a specified MariaDB Cluster from your account.

Required values to run command:

* Cluster Id`,
		Example: `ionosctl dbaas mariadb-v2 cluster delete --cluster-id <cluster-id>
ionosctl dbaas mariadb-v2 cluster delete --all
ionosctl dbaas mariadb-v2 cluster delete --all --name <name>`,
		PreCmdRun:  PreRunClusterDelete,
		CmdRun:     RunClusterDelete,
		InitClient: true,
	})
	deleteCmd.AddUUIDFlag(constants.FlagClusterId, constants.FlagIdShort, "", constants.DescCluster, core.RequiredFlagOption(),
		core.WithCompletion(completer.ClusterIds, constants.MariaDBApiRegionalURL, constants.MariaDBLocations),
	)
	deleteCmd.AddBoolFlag(constants.ArgAll, constants.ArgAllShort, false, "Delete all Clusters")
	deleteCmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "Delete all Clusters after filtering based on name. It does not require an exact match. Can be used with --all flag")
	deleteCmd.AddSetFlag(constants.FlagState, "", "", []string{"PROVISIONING", "AVAILABLE", "UPDATING", "DESTROYING", "FAILED"},
		"When used with --all, only delete clusters in this state")

	deleteCmd.Command.Flags().SortFlags = false
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
		return ClusterDeleteAll(c)
	}

	if err := c.RequireExplicitLocation(); err != nil {
		return err
	}

	clusterId := viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))

	c.Verbose(constants.ClusterId, clusterId)

	if !confirm.FAsk(c.Command.Command.InOrStdin(), fmt.Sprintf("delete cluster with id: %v and its backups", clusterId), viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	_, err := client.Must().MariaClientV2.ClustersApi.ClustersDelete(context.Background(), clusterId).Execute()
	if err != nil {
		return err
	}
	return nil
}

func ClusterDeleteAll(c *core.CommandConfig) error {
	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagName)) {
		c.Verbose("Filtering based on Cluster Name: %v", viper.GetString(core.GetFlagName(c.NS, constants.FlagName)))
	}

	// Gather clusters from every location (unless --location pins one), tagging each with its
	// location and location-scoped client, then hand the flat list to core.DeleteAll for a
	// consistent preview / per-item confirm-skip / summary flow.
	type located struct {
		cluster mariadb.ClusterRead
		loc     string
		api     *mariadb.APIClient
	}
	var items []located
	if err := c.RunForAllLocations(func(cfg *shared.Configuration, location string) error {
		api := mariadb.NewAPIClient(cfg)
		req := api.ClustersApi.ClustersGet(context.Background())
		if fn := core.GetFlagName(c.NS, constants.FlagName); viper.IsSet(fn) && viper.GetString(fn) != "" {
			req = req.FilterName(viper.GetString(fn))
		}
		if fn := core.GetFlagName(c.NS, constants.FlagState); viper.IsSet(fn) {
			req = req.FilterState(mariadb.MariadbClusterStates(viper.GetString(fn)))
		}
		clusters, _, err := req.Execute()
		if err != nil {
			return fmt.Errorf("failed listing clusters in location %s: %w", location, err)
		}
		for _, cluster := range clusters.GetItems() {
			items = append(items, located{cluster: cluster, loc: location, api: api})
		}
		return nil
	}); err != nil {
		return err
	}

	return core.DeleteAll(c, core.DeleteAllOptions[located]{
		Resource: "Cluster",
		List:     func() ([]located, error) { return items, nil },
		Summary: func(l located) string {
			return fmt.Sprintf("%s (%s) (location: %s)", l.cluster.Id, l.cluster.Properties.Name, l.loc)
		},
		ID: func(l located) string { return l.cluster.Id },
		Delete: func(l located) error {
			_, delErr := l.api.ClustersApi.ClustersDelete(context.Background(), l.cluster.Id).Execute()
			return delErr
		},
	})
}
