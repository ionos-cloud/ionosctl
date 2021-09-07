package waiter

import (
	"github.com/ionos-cloud/ionosctl/internal/core"
	cloudapiv5 "github.com/ionos-cloud/ionosctl/services/cloudapi-v5"
	"github.com/spf13/viper"
)

func ServerStateInterrogator(c *core.CommandConfig, objId string) (*string, error) {
	obj, _, err := c.CloudApiV5Services.Servers().Get(viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgDataCenterId)), objId)
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

func K8sClusterStateInterrogator(c *core.CommandConfig, objId string) (*string, error) {
	obj, _, err := c.CloudApiV5Services.K8s().GetCluster(objId)
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

func K8sNodeStateInterrogator(c *core.CommandConfig, objId string) (*string, error) {
	obj, _, err := c.CloudApiV5Services.K8s().GetNode(viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgK8sClusterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgK8sNodePoolId)), objId)
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

func K8sNodePoolStateInterrogator(c *core.CommandConfig, objId string) (*string, error) {
	obj, _, err := c.CloudApiV5Services.K8s().GetNodePool(viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgK8sClusterId)), objId)
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
