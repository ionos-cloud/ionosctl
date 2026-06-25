package helpers

import (
	"strings"

	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
)

func GetFlowLogPropertiesSet(c *core.CommandConfig) resources.FlowLogProperties {
	properties := resources.FlowLogProperties{}

	name := c.Flags().String(cloudapiv6.ArgName)
	properties.SetName(name)
	c.Verbose("Property Name set: %v", name)

	action := c.Flags().String(cloudapiv6.ArgAction)
	properties.SetAction(strings.ToUpper(action))
	c.Verbose("Property Action set: %v", action)

	direction := c.Flags().String(cloudapiv6.ArgDirection)
	properties.SetDirection(strings.ToUpper(direction))
	c.Verbose("Property Direction set: %v", direction)

	if c.Flags().Changed(cloudapiv6.ArgS3Bucket) {
		bucketName := c.Flags().String(cloudapiv6.ArgS3Bucket)
		properties.SetBucket(bucketName)
		c.Verbose("Property Bucket set: %v", bucketName)
	}

	return properties
}

// GetFlowLogPropertiesUpdate returns FlowLog Properties set used for update commands.
func GetFlowLogPropertiesUpdate(c *core.CommandConfig) resources.FlowLogProperties {
	properties := resources.FlowLogProperties{}

	if c.Flags().Changed(cloudapiv6.ArgName) {
		name := c.Flags().String(cloudapiv6.ArgName)
		properties.SetName(name)
		c.Verbose("Property Name set: %v", name)
	}

	if c.Flags().Changed(cloudapiv6.ArgAction) {
		action := c.Flags().String(cloudapiv6.ArgAction)
		properties.SetAction(strings.ToUpper(action))
		c.Verbose("Property Action set: %v", action)
	}

	if c.Flags().Changed(cloudapiv6.ArgDirection) {
		direction := c.Flags().String(cloudapiv6.ArgDirection)
		properties.SetDirection(strings.ToUpper(direction))
		c.Verbose("Property Direction set: %v", direction)
	}

	if c.Flags().Changed(cloudapiv6.ArgS3Bucket) {
		bucketName := c.Flags().String(cloudapiv6.ArgS3Bucket)
		properties.SetBucket(bucketName)
		c.Verbose("Property Bucket set: %v", bucketName)
	}

	return properties
}
