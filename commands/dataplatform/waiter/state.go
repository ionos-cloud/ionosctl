package waiter

import "github.com/ionos-cloud/ionosctl/pkg/core"

func ClusterStateInterrogator(c *core.CommandConfig, objId string) (*string, error) {
	obj, _, err := c.DataPlatformServices.Clusters().Get(objId)
	if err != nil {
		return nil, err
	}
	if metadataOk, ok := obj.GetMetadataOk(); ok && metadataOk != nil {
		if stateOk, ok := metadataOk.GetStateOk(); ok && stateOk != nil {
			return stateOk, nil
		}
	}
	return nil, nil
}
