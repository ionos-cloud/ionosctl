package completer

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"

	"github.com/ionos-cloud/ionosctl/v6/internal/functional"
	sdkgo "github.com/ionos-cloud/sdk-go-dbaas-mongo"
)

func MongoClusterIds() []string {
	ls, _, err := client.Must().MongoClient.ClustersApi.ClustersGet(context.Background()).Execute()
	if err != nil {
		return nil
	}
	return functional.Map(*ls.GetItems(), func(t sdkgo.ClusterResponse) string {
		var completion string
		completion = *t.Id
		if props, ok := t.GetPropertiesOk(); ok {
			if name, ok := props.GetDisplayNameOk(); ok {
				// Here is where the completion descriptions start
				completion = fmt.Sprintf("%s\t%s", completion, *name)
			}
			if location, ok := props.GetLocationOk(); ok {
				completion = fmt.Sprintf("%s - %s", completion, *location)
			}
		}

		return completion
	})
}

func MongoSnapshots(clusterId string) []string {
	ls, _, err := client.Must().MongoClient.SnapshotsApi.ClustersSnapshotsGet(context.Background(), clusterId).Execute()
	if err != nil {
		return nil
	}
	return functional.Map(*ls.GetItems(), func(t sdkgo.SnapshotResponse) string {
		var completion string
		completion = *t.Id
		if props, ok := t.GetPropertiesOk(); ok {
			if time, ok := props.GetCreationTimeOk(); ok {
				// Here is where the completion descriptions start
				completion = fmt.Sprintf("%s\t%s", completion, time.String())
			}
			if v, ok := props.GetVersionOk(); ok {
				completion = fmt.Sprintf("%s - v%s", completion, *v)
			}
		}

		return completion
	})
}

func MongoTemplateIds() []string {
	ls, _, err := client.Must().MongoClient.TemplatesApi.TemplatesGet(context.Background()).Execute()
	if err != nil {
		return nil
	}
	return functional.Map(*ls.GetItems(), func(t sdkgo.TemplateResponse) string {
		return *t.GetId()
	})
}
