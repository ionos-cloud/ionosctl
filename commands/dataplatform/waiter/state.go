package waiter

import (
	"github.com/ionos-cloud/ionosctl/pkg/core"
	dp "github.com/ionos-cloud/ionosctl/services/dataplatform"
	"github.com/spf13/viper"
)

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

func NodePoolStateInterrogator(c *core.CommandConfig, objId string) (*string, error) {
	obj, _, err := c.DataPlatformServices.NodePools().Get(viper.GetString(core.GetFlagName(c.NS, dp.ArgClusterId)), objId)
	if err != nil {
		return nil, err
	}
	if metadata, ok := obj.GetMetadataOk(); ok && metadata != nil {
		if state, ok := metadata.GetStateOk(); ok && state != nil {
			return state, nil
		}
	}
	return nil, nil
}
