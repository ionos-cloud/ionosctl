package waiter

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/viper"
)

func ServerStateInterrogator(c *core.CommandConfig, objId string) (*string, error) {
	obj, resp, err := c.CloudApiV6Services.Servers().Get(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)), objId)
	if err != nil {
		if resp != nil && resp.HttpNotFound() {
			return nil, nil // let it wait longer, maybe the resource is being created
		}
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
	obj, resp, err := c.CloudApiV6Services.K8s().GetCluster(objId)
	if err != nil {
		if resp != nil && resp.HttpNotFound() {
			return nil, nil // let it wait longer, maybe the resource is being created
		}
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
	obj, resp, err := c.CloudApiV6Services.K8s().GetNode(viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId)),
		viper.GetString(core.GetFlagName(c.NS, constants.FlagNodepoolId)), objId)
	if err != nil {
		if resp != nil && resp.HttpNotFound() {
			return nil, nil // let it wait longer, maybe the resource is being created
		}
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
	obj, resp, err := c.CloudApiV6Services.K8s().GetNodePool(viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId)), objId)
	if err != nil {
		if resp != nil && resp.HttpNotFound() {
			return nil, nil // let it wait longer, maybe the resource is being created
		}
		return nil, err
	}
	if metadata, ok := obj.GetMetadataOk(); ok && metadata != nil {
		if state, ok := metadata.GetStateOk(); ok && state != nil {
			return state, nil
		}
	}
	return nil, nil
}

func NatGatewayStateInterrogator(c *core.CommandConfig, objId string) (*string, error) {
	obj, resp, err := c.CloudApiV6Services.NatGateways().Get(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)), objId)
	if err != nil {
		if resp != nil && resp.HttpNotFound() {
			return nil, nil // let it wait longer, maybe the resource is being created
		}
		return nil, err
	}
	if metadata, ok := obj.GetMetadataOk(); ok && metadata != nil {
		if state, ok := metadata.GetStateOk(); ok && state != nil {
			return state, nil
		}
	}
	return nil, nil
}

func NetworkLoadBalancerStateInterrogator(c *core.CommandConfig, objId string) (*string, error) {
	obj, resp, err := c.CloudApiV6Services.NetworkLoadBalancers().Get(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)), objId)
	if err != nil {
		if resp != nil && resp.HttpNotFound() {
			return nil, nil // let it wait longer, maybe the resource is being created
		}
		return nil, err
	}
	if metadata, ok := obj.GetMetadataOk(); ok && metadata != nil {
		if state, ok := metadata.GetStateOk(); ok && state != nil {
			return state, nil
		}
	}
	return nil, nil
}

func ApplicationLoadBalancerStateInterrogator(c *core.CommandConfig, objId string) (*string, error) {
	obj, resp, err := c.CloudApiV6Services.ApplicationLoadBalancers().Get(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)), objId)
	if err != nil {
		if resp != nil && resp.HttpNotFound() {
			return nil, nil // let it wait longer, maybe the resource is being created
		}
		return nil, err
	}
	if metadata, ok := obj.GetMetadataOk(); ok && metadata != nil {
		if state, ok := metadata.GetStateOk(); ok && state != nil {
			return state, nil
		}
	}
	return nil, nil
}
