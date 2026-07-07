package cluster

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/kafka/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/confirm"
	"github.com/ionos-cloud/ionosctl/v6/pkg/functional"
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
	// Fan out over every location (when --location is unset) so `delete --all`
	// spans all locations, matching `list`. Each location gets its own client.
	return c.RunForAllLocations(func(cfg *shared.Configuration, location string) error {
		kc := kafka.NewAPIClient(cfg)
		c.Verbose("Deleting all clusters in %s!", location)

		records, _, err := kc.ClustersApi.ClustersGet(context.Background()).Execute()
		if err != nil {
			return fmt.Errorf("failed listing kafka clusters: %w", err)
		}

		return functional.ApplyAndAggregateErrors(
			records.GetItems(), func(d kafka.ClusterRead) error {
				yes := confirm.FAsk(
					c.Command.Command.InOrStdin(),
					fmt.Sprintf(
						"Are you sure you want to delete cluster %s with name %s (location: %s) ",
						d.Id, d.Properties.Name, location,
					),
					viper.GetBool(constants.ArgForce),
				)
				if yes {
					_, delErr := kc.ClustersApi.ClustersDelete(context.Background(), d.Id).Execute()
					if delErr != nil {
						return fmt.Errorf("failed deleting %s (name: %s): %w", d.Id, d.Properties.Name, delErr)
					}
				}
				return nil
			},
		)
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
