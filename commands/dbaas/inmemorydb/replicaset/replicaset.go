package replicaset

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/spf13/cobra"
)

var (
	allCols     = []string{"Id", "Name", "Version", "DNSName", "Replicas", "Cores", "RAM", "StorageSize", "State", "BackupLocation", "PersistenceMode", "EvictionPolicy", "MaintenanceDay", "MaintenanceTime", "DatacenterId", "LanId", "Username"}
	defaultCols = []string{"Id", "Name", "Version", "DNSName", "Replicas", "Cores", "RAM", "StorageSize", "State"}
)

func Root() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "replicaset",
			Aliases:          []string{"rs", "replica-set", "replicasets", "cluster"},
			Short:            "The sub-commands of 'ionosctl dbaas inmemorydb replicaset' allow you to manage In-Memory DB Replica Sets.",
			Long:             "In-Memory DB replica set with support for a single instance or a In-Memory DB replication in leader follower mode. The mode is determined by the number of replicas. One replica is standalone, everything else an In-Memory DB replication as leader follower mode with one active and n-1 passive replicas.",
			TraverseChildren: true,
		},
	}

	cmd.Command.PersistentFlags().StringSlice(constants.FlagCols, nil, tabheaders.ColsMessage(allCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(
		constants.FlagCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return allCols, cobra.ShellCompDirectiveNoFileComp
		},
	)

	cmd.AddCommand(Create())
	cmd.AddCommand(Get())
	cmd.AddCommand(List())
	cmd.AddCommand(Delete())
	// cmd.AddCommand(Update()) // Update is disabled until an API fix is rolled out

	return cmd
}
