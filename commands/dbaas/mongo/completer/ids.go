package completer

import (
	"context"
	client2 "github.com/ionos-cloud/ionosctl/v6/internal/client"

	"github.com/ionos-cloud/ionosctl/v6/internal/functional"
	sdkgo "github.com/ionos-cloud/sdk-go-dbaas-mongo"
)

func MongoClusterIds() []string {
	client, err := client2.Get()
	if err != nil {
		return nil
	}
	ls, _, err := client.MongoClient.ClustersApi.ClustersGet(context.Background()).Execute()
	if err != nil {
		return nil
	}
	return functional.Map(*ls.GetItems(), func(t sdkgo.ClusterResponse) string {
		return *t.GetId()
	})
}

func MongoSnapshots(clusterId string) []string {
	client, err := client2.Get()
	if err != nil {
		return nil
	}
	ls, _, err := client.MongoClient.SnapshotsApi.ClustersSnapshotsGet(context.Background(), clusterId).Execute()
	if err != nil {
		return nil
	}
	return functional.Map(*ls.GetItems(), func(t sdkgo.SnapshotResponse) string {
		return *t.GetId()
	})
}

func MongoTemplateIds() []string {
	client, err := client2.Get()
	if err != nil {
		return nil
	}
	ls, _, err := client.MongoClient.TemplatesApi.TemplatesGet(context.Background()).Execute()
	if err != nil {
		return nil
	}
	return functional.Map(*ls.GetItems(), func(t sdkgo.TemplateResponse) string {
		return *t.GetId()
	})
}
