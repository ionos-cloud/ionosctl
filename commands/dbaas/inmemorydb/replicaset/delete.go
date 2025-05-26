package replicaset

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/inmemorydb/utils"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/confirm"
	"github.com/ionos-cloud/ionosctl/v6/pkg/functional"
	"github.com/ionos-cloud/sdk-go-bundle/products/dbaas/inmemorydb/v2"
	"github.com/spf13/viper"
)

func Delete() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "dbaas inmemorydb",
		Resource:  "replicaset",
		Verb:      "delete",
		Aliases:   []string{"d", "del", "remove", "rm"},
		ShortDesc: "Delete In-Memory DB Replica Sets",
		Example: fmt.Sprintf(`ionosctl dbaas inmemorydb replicaset delete %s
ionosctl dbaas inmemorydb replicaset delete %s`,
			core.FlagsUsage(constants.FlagReplicasetID, constants.ArgForce),
			core.FlagsUsage(constants.ArgAll, constants.ArgForce)),
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return core.CheckRequiredFlagsSets(
				c.Command, c.NS,
				[]string{constants.ArgAll},
				[]string{constants.FlagReplicasetID},
			)
		},
		CmdRun: func(c *core.CommandConfig) error {
			if viper.GetBool(core.GetFlagName(c.NS, constants.ArgAll)) {
				return deleteAll(c)
			}
			id := viper.GetString(core.GetFlagName(c.NS, constants.FlagReplicasetID))
			return deleteSingle(c, id)
		},
		InitClient: true,
	})

	cmd.AddBoolFlag(constants.ArgAll, constants.ArgAllShort, false,
		fmt.Sprintf("Delete all replica-sets. Required or -%s", constants.FlagReplicasetID),
	)
	cmd.AddStringFlag(constants.FlagReplicasetID, constants.FlagIdShort, "",
		"The ID of the Replica Set you want to delete",
		core.WithCompletion(utils.ReplicasetIDs, constants.InMemoryDBApiRegionalURL, constants.InMemoryDBLocations),
	)

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	return cmd
}

// deleteSingle encapsulates: fetch the replica-set, ask to confirm deletion, delete.
func deleteSingle(c *core.CommandConfig, id string) error {
	rs, _, err := client.Must().InMemoryDBClient.
		ReplicaSetApi.
		ReplicasetsFindById(context.Background(), id).
		Execute()
	if err != nil {
		return fmt.Errorf("failed getting replicaset with id %q: %w", id, err)
	}

	prompt := fmt.Sprintf(
		"Are you sure you want to delete replica-set '%s' (dns name: '%s', replicas: %d)",
		rs.Properties.DisplayName,
		rs.Metadata.DnsName,
		rs.Properties.Replicas,
	)
	yes := confirm.FAsk(
		c.Command.Command.InOrStdin(),
		prompt,
		viper.GetBool(constants.ArgForce),
	)
	if !yes {
		return fmt.Errorf(confirm.UserDenied)
	}

	_, err = client.Must().InMemoryDBClient.
		ReplicaSetApi.
		ReplicasetsDelete(context.Background(), id).
		Execute()
	if err != nil {
		return fmt.Errorf("failed deleting replicaset %q: %w", id, err)
	}

	return nil
}

func deleteAll(c *core.CommandConfig) error {
	list, _, err := client.Must().InMemoryDBClient.
		ReplicaSetApi.
		ReplicasetsGet(c.Context).
		Execute()
	if err != nil {
		return fmt.Errorf("failed listing replica-sets: %w", err)
	}

	return functional.ApplyAndAggregateErrors(list.GetItems(), func(rs inmemorydb.ReplicaSetRead) error {
		return deleteSingle(c, rs.Id)
	})
}
