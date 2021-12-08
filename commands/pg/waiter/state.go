package waiter

import (
	"github.com/ionos-cloud/ionosctl/internal/core"
)

var deletionDone = "DONE"

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

func ClusterDeleteInterrogator(c *core.CommandConfig, objId string) (*string, error) {
	obj, resp, err := c.CloudApiDbaasPgsqlServices.Clusters().Get(objId)
	// If the cluster is no longer existing, return DONE
	// state to signal that the deletion is done.
	if resp != nil && resp.StatusCode == 404 {
		return &deletionDone, nil
	}
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
