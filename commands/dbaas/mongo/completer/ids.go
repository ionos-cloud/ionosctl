package completer

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"

	sdkgo "github.com/ionos-cloud/sdk-go-bundle/products/dbaas/mongo"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
)

func MongoClusterIds() []string {
	ls, _, err := client.Must().MongoClient.ClustersApi.ClustersGet(context.Background()).Execute()
	if err != nil {
		return nil
	}
	return shared.Map(*ls.GetItems(), func(t sdkgo.ClusterResponse) string {
		return *t.GetId()
	})
}

func MongoSnapshots(clusterId string) []string {
	ls, _, err := client.Must().MongoClient.SnapshotsApi.ClustersSnapshotsGet(context.Background(), clusterId).Execute()
	if err != nil {
		return nil
	}
	return shared.Map(*ls.GetItems(), func(t sdkgo.SnapshotResponse) string {
		return *t.GetId()
	})
}

func MongoTemplateIds() []string {
	ls, _, err := client.Must().MongoClient.TemplatesApi.TemplatesGet(context.Background()).Execute()
	if err != nil {
		return nil
	}
	return shared.Map(*ls.GetItems(), func(t sdkgo.TemplateResponse) string {
		return *t.GetId()
	})
}
