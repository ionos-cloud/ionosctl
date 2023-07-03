package completer

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"

	"github.com/ionos-cloud/ionosctl/v6/internal/functional"
	sdkgo "github.com/ionos-cloud/sdk-go-dbaas-mongo"
)

func MongoClusterIds() []string {
	ls, _, err := client.Must().MongoClient.ClustersApi.ClustersGet(context.Background()).Execute()
	if err != nil {
		return nil
	}
	return functional.Map(
		*ls.GetItems(), func(t sdkgo.ClusterResponse) string {
			return *t.GetId()
		},
	)
}

func MongoSnapshots(clusterId string) []string {
	ls, _, err := client.Must().MongoClient.SnapshotsApi.ClustersSnapshotsGet(context.Background(), clusterId).Execute()
	if err != nil {
		return nil
	}
	return functional.Map(
		*ls.GetItems(), func(t sdkgo.SnapshotResponse) string {
			return *t.GetId()
		},
	)
}

func MongoTemplateIds() []string {
	ls, _, err := client.Must().MongoClient.TemplatesApi.TemplatesGet(context.Background()).Execute()
	if err != nil {
		return nil
	}
	return functional.Map(
		*ls.GetItems(), func(t sdkgo.TemplateResponse) string {
			return *t.GetId()
		},
	)
}
