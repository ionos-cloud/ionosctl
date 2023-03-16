package completer

import (
	"context"

	client2 "github.com/ionos-cloud/ionosctl/v6/internal/client"

	"github.com/ionos-cloud/ionosctl/v6/internal/functional"
	sdkgo "github.com/ionos-cloud/sdk-go-dataplatform"
)

func DataplatformClusterIds() []string {
	client, err := client2.Get()
	if err != nil {
		return nil
	}
	ls, _, err := client.DataplatformClient.DataPlatformClusterApi.GetClusters(context.Background()).Execute()
	if err != nil {
		return nil
	}
	return functional.Map(*ls.GetItems(), func(t sdkgo.ClusterResponseData) string {
		return *t.GetId()
	})
}

func DataplatformNodepoolsIds(clusterId string) []string {
	client, err := client2.Get()
	if err != nil {
		return nil
	}
	ls, _, err := client.DataplatformClient.DataPlatformNodePoolApi.GetClusterNodepools(context.Background(), clusterId).Execute()
	if err != nil {
		return nil
	}
	return functional.Map(*ls.GetItems(), func(t sdkgo.NodePoolResponseData) string {
		return *t.GetId()
	})
}
