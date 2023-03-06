package waiter

import "github.com/ionos-cloud/ionosctl/v6/pkg/core"

func ClusterStateInterrogator(c *core.CommandConfig, objId string) (*string, error) {
	obj, _, err := c.CloudApiDbaasPgsqlServices.Clusters().Get(objId)
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
