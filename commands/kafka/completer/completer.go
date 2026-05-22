package completer

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/completions"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/ionos-cloud/ionosctl/v6/pkg/functional"
	"github.com/ionos-cloud/sdk-go-bundle/products/kafka/v2"
	"github.com/spf13/viper"
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
	kafkaClient := kafka.NewAPIClient(client.NewRegionalConfig(viper.GetString(constants.ArgServerUrl)))

	req := kafkaClient.ClustersApi.ClustersGet(context.Background())
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

var topicCompleterCols = []table.Column{
	{Name: "Id", JSONPath: "id"},
	{Name: "Name", JSONPath: "properties.name"},
}

// Topics returns all topics in the given cluster
func Topics(clusterID string) []string {
	kafkaClient := kafka.NewAPIClient(client.NewRegionalConfig(viper.GetString(constants.ArgServerUrl)))

	topicsList, _, err := kafkaClient.TopicsApi.ClustersTopicsGet(
		context.Background(), clusterID,
	).Execute()
	if err != nil {
		return nil
	}

	t := table.New(topicCompleterCols, table.WithPrefix("items"))
	if err := t.Extract(topicsList); err != nil {
		return nil
	}

	return completions.NewCompleter(t.Rows(), "Id").AddInfo("Name").ToString()
}

func Users(clusterID string) []string {
	kafkaClient := kafka.NewAPIClient(client.NewRegionalConfig(viper.GetString(constants.ArgServerUrl)))

	users, _, err := kafkaClient.UsersApi.ClustersUsersGet(
		context.Background(), clusterID,
	).Execute()
	if err != nil {
		return nil
	}

	ids := []string{}
	for _, u := range users.Items {
		ids = append(ids, u.Id+"\t"+u.Properties.Name)
	}
	return ids
}
