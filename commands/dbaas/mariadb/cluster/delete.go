package cluster

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/confirm"
	"github.com/ionos-cloud/sdk-go-bundle/products/dbaas/mariadb/v2"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
	"github.com/spf13/viper"
)

func confirmStringForCluster(c mariadb.ClusterResponse) string {
	askString := ""
	if p := c.Properties; p != nil {
		if c.Id != nil {
			askString = fmt.Sprintf("%s cluster %s", askString, *c.Id)
		}
		if n := p.DisplayName; n != nil {
			askString = fmt.Sprintf("%s (%s)", askString, *n)
		}
		if v := p.MariadbVersion; v != nil {
			askString = fmt.Sprintf("%s version v%s", askString, *v)
		}
	}
	return fmt.Sprintf("delete%s and its snapshots", askString)
}

func Delete() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "dbaas-mariadb",
		Resource:  "cluster",
		Verb:      "delete",
		Aliases:   []string{"del", "d"},
		ShortDesc: "Delete a MariaDB Cluster by ID",
		Example: `ionosctl dbaas mariadb cluster delete --cluster-id <cluster-id>
ionosctl db mar c d --all
ionosctl db mar c d --all --name <name>`,
		PreCmdRun: func(c *core.PreCommandConfig) error {
			if err := core.CheckRequiredFlagsSets(c.Command, c.NS, []string{constants.ArgAll}, []string{constants.FlagClusterId}); err != nil {
				return err
			}
			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {
			if all := viper.GetBool(core.GetFlagName(c.NS, constants.ArgAll)); all {
				return deleteAll(c)
			}

			if err := c.RequireExplicitLocation(); err != nil {
				return err
			}

			clusterId := viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))
			chosenCluster, _, err := client.Must().MariaClient.ClustersApi.ClustersFindById(context.Background(), clusterId).Execute()
			if err != nil {
				wrapped := fmt.Errorf("failed trying to find cluster by id: %w", err)
				keepGoing := confirm.FAsk(c.Command.Command.InOrStdin(), fmt.Sprintf("%s, try deleting %s anyways", wrapped.Error(), clusterId))
				if !keepGoing {
					return wrapped
				}
			}

			ok := confirm.FAsk(c.Command.Command.InOrStdin(), confirmStringForCluster(chosenCluster), viper.GetBool(constants.ArgForce))
			if !ok {
				return fmt.Errorf(confirm.UserDenied)
			}
			c.Verbose("Deleting cluster: %s", clusterId)

			_, _, err = client.Must().MariaClient.ClustersApi.ClustersDelete(context.Background(), clusterId).Execute()
			return err
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagClusterId, constants.FlagIdShort, "", "The unique ID of the cluster",
		core.RequiredFlagOption(),
		core.WithCompletion(
			func() []string {
				return ClustersProperty(func(c mariadb.ClusterResponse) string {
					if c.Id == nil {
						return ""
					}
					return *c.Id
				})
			}, constants.MariaDBApiRegionalURL, constants.MariaDBLocations),
	)
	cmd.AddBoolFlag(constants.ArgAll, constants.ArgAllShort, false, "Delete all mariadb clusters")
	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "When deleting all clusters, filter the clusters by a name")

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	return cmd
}

func clusterSummary(c mariadb.ClusterResponse) string {
	s := ""
	if c.Id != nil {
		s = *c.Id
	}
	if p := c.Properties; p != nil {
		if n := p.DisplayName; n != nil {
			s = fmt.Sprintf("%s (%s)", s, *n)
		}
		if v := p.MariadbVersion; v != nil {
			s = fmt.Sprintf("%s version v%s", s, *v)
		}
	}
	return s
}

func deleteAll(c *core.CommandConfig) error {
	// Gather clusters from every location (unless --location pins one), tagging each with
	// its location and the location-scoped client, then hand the flat list to core.DeleteAll
	// so the preview / per-item confirm-skip / summary flow is consistent across all resources.
	type located struct {
		cluster mariadb.ClusterResponse
		loc     string
		api     *mariadb.APIClient
	}
	var items []located
	if err := c.RunForAllLocations(func(cfg *shared.Configuration, location string) error {
		api := mariadb.NewAPIClient(cfg)
		req := FilterNameFlags(c)(api.ClustersApi.ClustersGet(context.Background()))
		xs, _, err := req.Execute()
		if err != nil {
			return fmt.Errorf("failed getting clusters: %w", err)
		}
		for _, x := range xs.GetItems() {
			items = append(items, located{cluster: x, loc: location, api: api})
		}
		return nil
	}); err != nil {
		return err
	}

	return core.DeleteAll(c, core.DeleteAllOptions[located]{
		Resource: "cluster",
		List:     func() ([]located, error) { return items, nil },
		Summary: func(l located) string {
			return fmt.Sprintf("%s (location: %s)", clusterSummary(l.cluster), l.loc)
		},
		ID: func(l located) string { return *l.cluster.Id },
		Delete: func(l located) error {
			_, _, delErr := l.api.ClustersApi.ClustersDelete(context.Background(), *l.cluster.Id).Execute()
			return delErr
		},
	})
}
