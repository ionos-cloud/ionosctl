package utils

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/config"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/functional"
	"github.com/ionos-cloud/sdk-go-bundle/products/dbaas/inmemorydb/v2"
	"github.com/spf13/viper"
)

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

func ReplicasetIDs() []string {
	return ReplicasetProperty(func(r inmemorydb.ReplicaSetRead) string {
		return fmt.Sprintf("%s\t%s (dns name '%s', '%d' replicas)", r.Id, r.Properties.DisplayName, r.Metadata.DnsName, r.Properties.Replicas)
	})
}
