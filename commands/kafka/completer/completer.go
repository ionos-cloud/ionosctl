package completer

import (
	"context"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/pkg/functional"
	kafka "github.com/ionos-cloud/sdk-go-kafka"
)

// ClustersProperty returns a list of properties of all clusters matching the given filters
func ClustersProperty[V any](f func(read kafka.ClusterRead) V, fs ...Filter) []V {
	recs, err := Clusters(fs...)
	if err != nil {
		return nil
	}
	return functional.Map(*recs.Items, f)
}

// Clusters returns all clusters matching the given filters
func Clusters(fs ...Filter) (kafka.ClusterReadList, error) {
	req := client.Must().Kafka.ClustersApi.ClustersGet(context.Background())
	for _, f := range fs {
		var err error
		req, err = f(req)
		if err != nil {
			return kafka.ClusterReadList{}, err
		}
	}

	ls, _, err := req.Execute()
	if err != nil {
		return kafka.ClusterReadList{}, err
	}

	return ls, nil
}

type Filter func(request kafka.ApiClustersGetRequest) (kafka.ApiClustersGetRequest, error)
