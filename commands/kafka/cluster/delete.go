package cluster

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/kafka/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/confirm"
	"github.com/ionos-cloud/sdk-go-bundle/products/kafka/v2"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
	"github.com/spf13/viper"

	"github.com/ionos-cloud/ionosctl/v6/internal/core"
)

func Delete() *core.Command {
	cmd := core.NewCommand(
		context.Background(), nil, core.CommandBuilder{
			Namespace: "kafka",
			Resource:  "cluster",
			Verb:      "delete",
			Aliases:   []string{"del", "d"},
			ShortDesc: "Delete a cluster",
			Example:   `ionosctl kafka cl delete --cluster-id ID`,
			PreCmdRun: func(c *core.PreCommandConfig) error {
				if err := core.CheckRequiredFlagsSets(
					c.Command, c.NS,
					[]string{constants.FlagClusterId, constants.FlagLocation},
					[]string{constants.ArgAll},
				); err != nil {
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

				return deleteSingle(c, viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId)))
			},
			InitClient: true,
		},
	)

	cmd.AddStringFlag(
		constants.FlagClusterId, constants.FlagIdShort, "", "The ID of the cluster you want to retrieve",
		core.RequiredFlagOption(),
		core.WithCompletion(
			func() []string {
				return completer.ClustersProperty(
					func(k kafka.ClusterRead) string {
						return k.Id
					},
				)
			}, constants.KafkaApiRegionalURL, constants.KafkaLocations,
		),
	)

	cmd.AddBoolFlag(
		constants.ArgAll, constants.ArgAllShort, false, "Delete all records if set", core.RequiredFlagOption(),
	)

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	return cmd
}

func deleteAll(c *core.CommandConfig) error {
	// Gather clusters from every location (unless --location pins one), tagging each with its
	// location and location-scoped client, then hand the flat list to core.DeleteAll for a
	// consistent preview / per-item confirm-skip / summary flow.
	type located struct {
		cluster kafka.ClusterRead
		loc     string
		api     *kafka.APIClient
	}
	var items []located
	if err := c.RunForAllLocations(func(cfg *shared.Configuration, location string) error {
		kc := kafka.NewAPIClient(cfg)
		records, _, err := kc.ClustersApi.ClustersGet(context.Background()).Execute()
		if err != nil {
			return fmt.Errorf("failed listing kafka clusters: %w", err)
		}
		for _, d := range records.GetItems() {
			items = append(items, located{cluster: d, loc: location, api: kc})
		}
		return nil
	}); err != nil {
		return err
	}

	return core.DeleteAll(c, core.DeleteAllOptions[located]{
		Resource: "cluster",
		List:     func() ([]located, error) { return items, nil },
		Summary: func(l located) string {
			return fmt.Sprintf("name: %s, id: %s (location: %s)", l.cluster.Properties.Name, l.cluster.Id, l.loc)
		},
		ID: func(l located) string { return l.cluster.Id },
		Delete: func(l located) error {
			_, err := l.api.ClustersApi.ClustersDelete(context.Background(), l.cluster.Id).Execute()
			return err
		},
	})
}

func deleteSingle(c *core.CommandConfig, id string) error {
	d, _, err := client.Must().Kafka.ClustersApi.ClustersFindById(context.Background(), id).Execute()
	if err != nil {
		return fmt.Errorf("cluster not found: %w", err)
	}

	yes := confirm.FAsk(
		c.Command.Command.InOrStdin(),
		fmt.Sprintf("Are you sure you want to delete cluster %s with name %s", d.Id, d.Properties.Name),
		viper.GetBool(constants.ArgForce),
	)
	if !yes {
		return fmt.Errorf("user cancelled deletion")
	}

	_, err = client.Must().Kafka.ClustersApi.ClustersDelete(context.Background(), d.Id).Execute()
	return err
}
