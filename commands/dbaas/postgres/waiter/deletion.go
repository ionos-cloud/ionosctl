package waiter

import "github.com/ionos-cloud/ionosctl/v6/pkg/core"

func ClusterDeleteInterrogator(c *core.CommandConfig, objId string) (*int, error) {
	_, resp, err := c.CloudApiDbaasPgsqlServices.Clusters().Get(objId)
	// Return HTTP Response Status Code
	if resp != nil && resp.Response != nil {
		return &resp.Response.StatusCode, nil
	}
	if err != nil {
		return nil, err
	}
	return nil, nil
}
