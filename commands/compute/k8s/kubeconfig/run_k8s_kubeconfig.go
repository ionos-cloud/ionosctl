package kubeconfig

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
)

func RunK8sKubeconfigGet(c *core.CommandConfig) error {
	c.Verbose("K8s kube config with id: %v is getting...", c.Flags().String(constants.FlagClusterId))

	u, resp, err := c.CloudApiV6Services.K8s().ReadKubeConfig(c.Flags().String(constants.FlagClusterId))
	if resp != nil {
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	c.Msg("%s", u)

	return nil
}

func PreRunK8sClusterId(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, constants.FlagClusterId)
}
