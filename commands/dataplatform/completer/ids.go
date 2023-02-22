package completer

import (
	"context"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/utils"
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
	return utils.MapNoIdx(*ls.GetItems(), func(t sdkgo.ClusterResponseData) string {
		return *t.GetId()
	})
}
