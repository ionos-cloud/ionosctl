package waiter

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/viper"
)

// Deprecated: ServerStateInterrogator is a legacy interrogator for WaitForState.
// New code should rely on globalwait.WaitAndRerender.
// Only remaining caller: commands/compute/server/run_server.go (--promote-volume).
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
