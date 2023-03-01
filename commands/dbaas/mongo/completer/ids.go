package completer

import (
	"context"
	"github.com/ionos-cloud/ionosctl/internal/functional"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	sdkgo "github.com/ionos-cloud/sdk-go-dbaas-mongo"
)

func MongoClusterIds() []string {
	client, err := config.GetClient()
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

func MongoTemplateIds() []string {
	client, err := config.GetClient()
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
