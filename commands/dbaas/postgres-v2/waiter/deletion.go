package waiter

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
)

func ClusterDeleteInterrogator(c *core.CommandConfig, objId string) (*int, error) {
	_, resp, err := client.Must().PostgresClientV2.ClustersApi.ClustersFindById(context.Background(), objId).Execute()
	if resp != nil && resp.Response != nil {
		return &resp.Response.StatusCode, nil
	}
	if err != nil {
		return nil, err
	}
	return nil, nil
}
