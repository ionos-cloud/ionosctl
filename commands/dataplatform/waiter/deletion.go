package waiter

import (
	"github.com/ionos-cloud/ionosctl/pkg/core"
	dp "github.com/ionos-cloud/ionosctl/services/dataplatform"
	"github.com/spf13/viper"
)

func ClusterDeleteInterrogator(c *core.CommandConfig, objId string) (*int, error) {
	_, resp, err := c.DataPlatformServices.Clusters().Get(objId)
	// Return HTTP Response Status Code
	if resp != nil && resp.Response != nil {
		return &resp.Response.StatusCode, nil
	}
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func NodePoolDeleteInterrogator(c *core.CommandConfig, objId string) (*int, error) {
	_, resp, err := c.DataPlatformServices.NodePools().Get(viper.GetString(core.GetFlagName(c.NS, dp.ArgNodePoolId)), objId)
	// Return HTTP Response Status Code
	if resp != nil && resp.Response != nil {
		return &resp.Response.StatusCode, nil
	}
	if err != nil {
		return nil, err
	}
	return nil, nil
}
