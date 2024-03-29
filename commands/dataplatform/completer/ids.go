package completer

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/pkg/functional"

	sdkgo "github.com/ionos-cloud/sdk-go-dataplatform"
)

func DataplatformClusterIds() []string {
	ls, _, err := client.Must().DataplatformClient.DataPlatformClusterApi.ClustersGet(context.Background()).Execute()
	if err != nil {
		return nil
	}
	return functional.Map(*ls.GetItems(), func(t sdkgo.ClusterResponseData) string {
		return *t.GetId()
	})
}

func DataplatformNodepoolsIds(clusterId string) []string {
	ls, _, err := client.Must().DataplatformClient.DataPlatformNodePoolApi.ClustersNodepoolsGet(context.Background(), clusterId).Execute()
	if err != nil {
		return nil
	}
	return functional.Map(*ls.GetItems(), func(t sdkgo.NodePoolResponseData) string {
		return *t.GetId()
	})
}
