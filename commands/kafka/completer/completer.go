package completer

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/completions"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/pkg/functional"
	"github.com/ionos-cloud/sdk-go-bundle/products/kafka/v2"
)

// ClustersProperty returns a list of properties of all clusters matching the given filters
func ClustersProperty[V any](f func(read kafka.ClusterRead) V, fs ...Filter) []V {
	recs, err := Clusters(fs...)
	if err != nil {
		return nil
	}
	return functional.Map(recs.Items, f)
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

// Topics returns all topics in the given cluster
func Topics(clusterID string) []string {
	topicsList, _, err := client.Must().Kafka.TopicsApi.ClustersTopicsGet(
		context.Background(), clusterID,
	).Execute()
	if err != nil {
		return nil
	}

	topicsConverted, err := json2table.ConvertJSONToTable("items", jsonpaths.KafkaTopic, topicsList)
	if err != nil {
		return nil
	}

	return completions.NewCompleter(topicsConverted, "Id").AddInfo("Name").ToString()
}
