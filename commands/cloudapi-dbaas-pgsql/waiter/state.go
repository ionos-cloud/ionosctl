package waiter

import (
	"github.com/ionos-cloud/ionosctl/internal/core"
)

func ClusterStateInterrogator(c *core.CommandConfig, objId string) (*string, error) {
	obj, _, err := c.CloudApiDbaasPgsqlServices.Clusters().Get(objId)
	if err != nil {
		return nil, err
	}
	if lifecycleStatusOk, ok := obj.GetLifecycleStatusOk(); ok && lifecycleStatusOk != nil {
		return lifecycleStatusOk, nil
	}
	return nil, nil
}
