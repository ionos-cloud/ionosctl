package version

import (
	"strings"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
)

func RunK8sVersionList(c *core.CommandConfig) error {
	u, resp, err := c.CloudApiV6Services.K8s().ListVersions()
	if resp != nil {
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}

	if err != nil {
		return err
	}

	c.Msg("%s", u)

	return nil
}

func RunK8sVersionGet(c *core.CommandConfig) error {
	u, err := getK8sVersion(c)
	if err != nil {
		return err
	}

	c.Msg("%s", u)

	return nil
}

func getK8sVersion(c *core.CommandConfig) (string, error) {
	k8sversion, resp, err := c.CloudApiV6Services.K8s().GetVersion()
	if err != nil {
		return "", err
	}

	k8sversion = strings.ReplaceAll(k8sversion, "\"", "")
	k8sversion = strings.ReplaceAll(k8sversion, "\n", "")
	if resp != nil {
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}

	return k8sversion, nil
}
