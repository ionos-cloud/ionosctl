package completer

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/functional"
	"github.com/ionos-cloud/ionosctl/v6/pkg/config"
	sdkgo "github.com/ionos-cloud/sdk-go-dataplatform"
)

func DataplatformClusterIds() []string {
	client, err := config.GetClient()
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
	client, err := config.GetClient()
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
