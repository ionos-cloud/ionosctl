package helpers

import (
	"strings"

	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
	"github.com/spf13/viper"
)

func GetFlowLogPropertiesSet(c *core.CommandConfig) resources.FlowLogProperties {
	properties := resources.FlowLogProperties{}

	name := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName))
	properties.SetName(name)
	c.Verbose("Property Name set: %v", name)

	action := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgAction))
	properties.SetAction(strings.ToUpper(action))
	c.Verbose("Property Action set: %v", action)

	direction := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDirection))
	properties.SetDirection(strings.ToUpper(direction))
	c.Verbose("Property Direction set: %v", direction)

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgS3Bucket)) {
		bucketName := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgS3Bucket))
		properties.SetBucket(bucketName)
		c.Verbose("Property Bucket set: %v", bucketName)
	}

	return properties
}

// GetFlowLogPropertiesUpdate returns FlowLog Properties set used for update commands.
func GetFlowLogPropertiesUpdate(c *core.CommandConfig) resources.FlowLogProperties {
	properties := resources.FlowLogProperties{}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgName)) {
		name := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName))
		properties.SetName(name)
		c.Verbose("Property Name set: %v", name)
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgAction)) {
		action := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgAction))
		properties.SetAction(strings.ToUpper(action))
		c.Verbose("Property Action set: %v", action)
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgDirection)) {
		direction := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDirection))
		properties.SetDirection(strings.ToUpper(direction))
		c.Verbose("Property Direction set: %v", direction)
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgS3Bucket)) {
		bucketName := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgS3Bucket))
		properties.SetBucket(bucketName)
		c.Verbose("Property Bucket set: %v", bucketName)
	}

	return properties
}
