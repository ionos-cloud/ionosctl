package resource

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
)

func PreRunResourceType(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, constants.FlagType)
}

func RunResourceList(c *core.CommandConfig) error {
	resourcesListed, resp, err := c.CloudApiV6Services.Users().ListResources()
	if resp != nil {
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allResourceCols).Prefix("items").Print(resourcesListed.Resources)
}

func RunResourceGet(c *core.CommandConfig) error {
	c.Verbose("Resource with id: %v is getting...", c.Flags().String(cloudapiv6.ArgResourceId))

	if c.Flags().Changed(cloudapiv6.ArgResourceId) {
		resourceListed, resp, err := c.CloudApiV6Services.Users().GetResourceByTypeAndId(
			c.Flags().String(constants.FlagType),
			c.Flags().String(cloudapiv6.ArgResourceId),
		)
		if resp != nil {
			c.Verbose(constants.MessageRequestTime, resp.RequestTime)
		}
		if err != nil {
			return err
		}

		return c.Printer(allResourceCols).Print(resourceListed.Resource)
	}

	resourcesListed, resp, err := c.CloudApiV6Services.Users().GetResourcesByType(c.Flags().String(constants.FlagType))
	if resp != nil {
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allResourceCols).Prefix("items").Print(resourcesListed.Resources)
}

func RunGroupResourceList(c *core.CommandConfig) error {
	c.Verbose("Listing Resources from Group with ID: %v...", c.Flags().String(cloudapiv6.ArgGroupId))

	resourcesListed, resp, err := c.CloudApiV6Services.Groups().ListResources(c.Flags().String(cloudapiv6.ArgGroupId))
	if resp != nil {
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allResourceCols).Prefix("items").Print(resourcesListed.ResourceGroups)
}
