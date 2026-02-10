package waiter

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
)

func ClusterStateInterrogator(c *core.CommandConfig, objId string) (*string, error) {
	obj, _, err := client.Must().PostgresClientV2.ClustersApi.ClustersFindById(context.Background(), objId).Execute()
	if err != nil {
		return nil, err
	}
	if metadataOk, ok := obj.GetMetadataOk(); ok && metadataOk != nil {
		if stateOk, ok := metadataOk.GetStateOk(); ok && stateOk != nil {
			return (*string)(stateOk), nil
		}
	}
	return nil, nil
}
