package completer

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"

	sdkgo "github.com/ionos-cloud/sdk-go-bundle/products/dataplatform"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
)

func DataplatformClusterIds() []string {
	ls, _, err := client.Must().DataplatformClient.DataPlatformClusterApi.ClustersGet(context.Background()).Execute()
	if err != nil {
		return nil
	}
	return shared.Map(*ls.GetItems(), func(t sdkgo.ClusterResponseData) string {
		return *t.GetId()
	})
}

func DataplatformNodepoolsIds(clusterId string) []string {
	ls, _, err := client.Must().DataplatformClient.DataPlatformNodePoolApi.ClustersNodepoolsGet(context.Background(), clusterId).Execute()
	if err != nil {
		return nil
	}
	return shared.Map(*ls.GetItems(), func(t sdkgo.NodePoolResponseData) string {
		return *t.GetId()
	})
}
