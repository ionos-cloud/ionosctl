package completer

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/pkg/convbytes"
	"github.com/ionos-cloud/ionosctl/v6/pkg/functional"

	sdkgo "github.com/ionos-cloud/sdk-go-bundle/products/dbaas/mongo/v2"
)

func MongoClusterIds() []string {
	ls, _, err := client.Must().MongoClient.ClustersApi.ClustersGet(context.Background()).Execute()
	if err != nil {
		return nil
	}
	return functional.Map(ls.GetItems(), func(t sdkgo.ClusterResponse) string {
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
	return functional.Map(ls.GetItems(), func(t sdkgo.SnapshotResponse) string {
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
			if size, ok := props.GetSizeOk(); ok {
				completion = fmt.Sprintf("%s (%d GB)", completion, convbytes.Convert(int64(*size), convbytes.MB, convbytes.GB))
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
	return functional.Map(ls.GetItems(), func(t sdkgo.TemplateResponse) string {
		return t.GetId()
	})
}

func MongoDBVersions() []string {
	ls, _, err := client.Must().MongoClient.MetadataApi.VersionsGet(context.Background()).Execute()
	if err != nil {
		return nil
	}
	return functional.Map(ls.GetData(), func(t sdkgo.MongoDBVersionListData) string {
		return *t.Name
	})
}

func MongoClusterVersions(clusterid string) []string {
	ls, _, err := client.Must().MongoClient.ClustersApi.ClustersVersionsGet(context.Background(), clusterid).Execute()
	if err != nil {
		return nil
	}
	return functional.Map(ls.GetData(), func(t sdkgo.MongoDBVersionListData) string {
		return *t.Name
	})
}
