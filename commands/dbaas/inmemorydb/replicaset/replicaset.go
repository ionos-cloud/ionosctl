package replicaset

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/config"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/ionosctl/v6/pkg/functional"
	"github.com/ionos-cloud/sdk-go-bundle/products/dbaas/inmemorydb/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

	cmd.Command.PersistentFlags().StringSlice(constants.ArgCols, nil, tabheaders.ColsMessage(allCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(
		constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return allCols, cobra.ShellCompDirectiveNoFileComp
		},
	)

	cmd.AddCommand(Create())
	cmd.AddCommand(List())

	return cmd
}

func ReplicasetProperty[V any](f func(inmemorydb.ReplicaSetRead) V, fs ...Filter) []V {
	recs, err := Replicasets(fs...)
	if err != nil {
		return nil
	}
	return functional.Map(recs.Items, f)
}

func Replicasets(fs ...Filter) (inmemorydb.ReplicaSetReadList, error) {
	if url := config.GetServerUrl(); url == constants.DefaultApiURL {
		viper.Set(constants.ArgServerUrl, "")
	}

	req := client.Must().InMemoryDBClient.ReplicaSetApi.ReplicasetsGet(context.Background())

	for _, f := range fs {
		var err error
		req, err = f(req)
		if err != nil {
			return inmemorydb.ReplicaSetReadList{}, err
		}
	}

	ls, _, err := req.Execute()
	if err != nil {
		return inmemorydb.ReplicaSetReadList{}, err
	}
	return ls, nil
}

type Filter func(inmemorydb.ApiReplicasetsGetRequest) (inmemorydb.ApiReplicasetsGetRequest, error)
