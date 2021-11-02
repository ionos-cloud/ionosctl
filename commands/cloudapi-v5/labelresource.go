package cloudapi_v5

import (
	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/commands/cloudapi-v5/waiter"
	"github.com/ionos-cloud/ionosctl/internal/config"
	"github.com/ionos-cloud/ionosctl/internal/core"
	"github.com/ionos-cloud/ionosctl/internal/printer"
	"github.com/ionos-cloud/ionosctl/internal/utils"
	cloudapiv5 "github.com/ionos-cloud/ionosctl/services/cloudapi-v5"
	"github.com/ionos-cloud/ionosctl/services/cloudapi-v5/resources"
	"github.com/spf13/viper"
)

func RunDataCenterLabelsList(c *core.CommandConfig) error {
	labelDcs, resp, err := c.CloudApiV5Services.Labels().DatacenterList(viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgDataCenterId)))
	if resp != nil {
		c.Printer.Verbose(config.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, getLabelResources(labelDcs)))
}

func RunDataCenterLabelGet(c *core.CommandConfig) error {
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgDataCenterId))
	labelKey := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgLabelKey))
	c.Printer.Verbose("Getting label with key: %v for Datacenter with id: %v...", labelKey, dcId)
	labelDc, resp, err := c.CloudApiV5Services.Labels().DatacenterGet(dcId, labelKey)
	if resp != nil {
		c.Printer.Verbose(config.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, getLabelResource(labelDc)))
}

func RunDataCenterLabelAdd(c *core.CommandConfig) error {
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgDataCenterId))
	labelKey := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgLabelKey))
	labelValue := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgLabelValue))
	c.Printer.Verbose("Adding label with key: %v and value: %v to Datacenter with id: %v...", labelKey, labelValue, dcId)
	labelDc, resp, err := c.CloudApiV5Services.Labels().DatacenterCreate(dcId, labelKey, labelValue)
	if resp != nil {
		c.Printer.Verbose(config.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, getLabelResource(labelDc)))
}

func RunDataCenterLabelRemove(c *core.CommandConfig) error {
	var resp *resources.Response
	var err error
	allFlag := viper.GetBool(core.GetFlagName(c.NS, cloudapiv5.ArgAll))
	if allFlag {
		resp, err = RemoveAllDatacenterLabels(c)
		if err != nil {
			return err
		}
	} else {
		dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgDataCenterId))
		labelKey := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgLabelKey))
		c.Printer.Verbose("Removing label with key: %v for Datacenter with id: %v...", labelKey, dcId)
		resp, err = c.CloudApiV5Services.Labels().DatacenterDelete(dcId, labelKey)
		if resp != nil {
			c.Printer.Verbose(config.RequestTimeMessage, resp.RequestTime)
		}
		if err != nil {
			return err
		}
	}

	return c.Printer.Print(getLabelResourcePrint(c, nil))
}

func RemoveAllDatacenterLabels(c *core.CommandConfig) (*resources.Response, error) {
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgDataCenterId))
	_ = c.Printer.Print("Labels to be removed from Datacenter with Id: " + dcId)
	labels, resp, err := c.CloudApiV5Services.Labels().DatacenterList(dcId)
	if err != nil {
		return nil, err
	}
	if labelsItems, ok := labels.GetItemsOk(); ok && labelsItems != nil {
		for _, label := range *labelsItems {
			if properties, ok := label.GetPropertiesOk(); ok && properties != nil {
				if key, ok := properties.GetKeyOk(); ok && key != nil {
					_ = c.Printer.Print(" Label Key: " + *key)
				}
				if value, ok := properties.GetValueOk(); ok && value != nil {
					_ = c.Printer.Print(" Label Value: " + *value)
				}
			}
		}

		if err := utils.AskForConfirm(c.Stdin, c.Printer, "delete all the Labels from Datacenter with Id: "+dcId); err != nil {
			return nil, err
		}
		c.Printer.Verbose("Deleting all the Labels from Datacenter with Id: %v...", dcId)
		for _, label := range *labelsItems {
			if properties, ok := label.GetPropertiesOk(); ok && properties != nil {
				if key, ok := properties.GetKeyOk(); ok && key != nil {
					c.Printer.Verbose("Starting deleting Label with id: %v...", *key)
					resp, err = c.CloudApiV5Services.Labels().DatacenterDelete(dcId, *key)
					if resp != nil {
						c.Printer.Verbose("Request Id: %v", printer.GetId(resp))
						c.Printer.Verbose(config.RequestTimeMessage, resp.RequestTime)
					}
					if err != nil {
						return nil, err
					}
					if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
						return nil, err
					}
				}
			}
		}
	}
	return resp, nil
}

func RunServerLabelsList(c *core.CommandConfig) error {
	labelDcs, resp, err := c.CloudApiV5Services.Labels().ServerList(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgServerId)),
	)
	if resp != nil {
		c.Printer.Verbose(config.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, getLabelResources(labelDcs)))
}

func RunServerLabelGet(c *core.CommandConfig) error {
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgDataCenterId))
	serverId := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgServerId))
	labelkey := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgLabelKey))
	c.Printer.Verbose("Getting label with key: %v for Server with id: %v...", labelkey, serverId)
	labelDc, resp, err := c.CloudApiV5Services.Labels().ServerGet(dcId, serverId, labelkey)
	if resp != nil {
		c.Printer.Verbose(config.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, getLabelResource(labelDc)))
}

func RunServerLabelAdd(c *core.CommandConfig) error {
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgDataCenterId))
	serverId := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgServerId))
	labelKey := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgLabelKey))
	labelValue := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgLabelValue))
	c.Printer.Verbose("Adding label with key: %v and value: %v to Server with id: %v...", labelKey, labelValue, serverId)
	labelDc, resp, err := c.CloudApiV5Services.Labels().ServerCreate(dcId, serverId, labelKey, labelValue)
	if resp != nil {
		c.Printer.Verbose(config.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, getLabelResource(labelDc)))
}

func RunServerLabelRemove(c *core.CommandConfig) error {
	var resp *resources.Response
	var err error
	allFlag := viper.GetBool(core.GetFlagName(c.NS, cloudapiv5.ArgAll))
	if allFlag {
		resp, err = RemoveAllServerLabels(c)
		if err != nil {
			return err
		}
	} else {
		dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgDataCenterId))
		serverId := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgServerId))
		labelKey := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgLabelKey))
		c.Printer.Verbose("Removing label with key: %v for Server with id: %v...", labelKey, serverId)
		resp, err = c.CloudApiV5Services.Labels().ServerDelete(dcId, serverId, labelKey)
		if resp != nil {
			c.Printer.Verbose(config.RequestTimeMessage, resp.RequestTime)
		}
		if err != nil {
			return err
		}
	}

	return c.Printer.Print(getLabelResourcePrint(c, nil))
}

func RemoveAllServerLabels(c *core.CommandConfig) (*resources.Response, error) {
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgDataCenterId))
	serverId := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgServerId))
	_ = c.Printer.Print("Labels to be removed from Server with Id: " + serverId)
	labels, resp, err := c.CloudApiV5Services.Labels().ServerList(dcId, serverId)
	if err != nil {
		return nil, err
	}
	if labelsItems, ok := labels.GetItemsOk(); ok && labelsItems != nil {
		for _, label := range *labelsItems {
			if properties, ok := label.GetPropertiesOk(); ok && properties != nil {
				if key, ok := properties.GetKeyOk(); ok && key != nil {
					_ = c.Printer.Print(" Label Key: " + *key)
				}
				if value, ok := properties.GetValueOk(); ok && value != nil {
					_ = c.Printer.Print(" Label Value: " + *value)
				}
			}
		}

		if err := utils.AskForConfirm(c.Stdin, c.Printer, "delete all the Labels from Server with Id: "+serverId); err != nil {
			return nil, err
		}
		c.Printer.Verbose("Deleting all the Labels from Server with Id: %v...", serverId)
		for _, label := range *labelsItems {
			if properties, ok := label.GetPropertiesOk(); ok && properties != nil {
				if key, ok := properties.GetKeyOk(); ok && key != nil {
					c.Printer.Verbose("Starting deleting Label with id: %v...", *key)
					resp, err = c.CloudApiV5Services.Labels().ServerDelete(dcId, serverId, *key)
					if resp != nil {
						c.Printer.Verbose("Request Id: %v", printer.GetId(resp))
						c.Printer.Verbose(config.RequestTimeMessage, resp.RequestTime)
					}
					if err != nil {
						return nil, err
					}
					if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
						return nil, err
					}
				}

			}
		}
	}
	return resp, nil
}

func RunVolumeLabelsList(c *core.CommandConfig) error {
	labelDcs, resp, err := c.CloudApiV5Services.Labels().VolumeList(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgVolumeId)),
	)
	if resp != nil {
		c.Printer.Verbose(config.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, getLabelResources(labelDcs)))
}

func RunVolumeLabelGet(c *core.CommandConfig) error {
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgDataCenterId))
	volumeId := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgVolumeId))
	labelKey := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgLabelKey))
	c.Printer.Verbose("Getting label with key: %v for Volume with id: %v...", labelKey, volumeId)
	labelDc, resp, err := c.CloudApiV5Services.Labels().VolumeGet(dcId, volumeId, labelKey)
	if resp != nil {
		c.Printer.Verbose(config.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, getLabelResource(labelDc)))
}

func RunVolumeLabelAdd(c *core.CommandConfig) error {
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgDataCenterId))
	volumeId := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgVolumeId))
	labelKey := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgLabelKey))
	labelValue := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgLabelValue))
	c.Printer.Verbose("Adding label with key: %v and value: %v to Volume with id: %v...", labelKey, labelValue, volumeId)
	labelDc, resp, err := c.CloudApiV5Services.Labels().VolumeCreate(dcId, volumeId, labelKey, labelValue)
	if resp != nil {
		c.Printer.Verbose(config.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, getLabelResource(labelDc)))
}

func RunVolumeLabelRemove(c *core.CommandConfig) error {
	var resp *resources.Response
	var err error
	allFlag := viper.GetBool(core.GetFlagName(c.NS, cloudapiv5.ArgAll))
	if allFlag {
		resp, err = RemoveAllVolumeLabels(c)
		if err != nil {
			return err
		}
	} else {
		dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgDataCenterId))
		volumeId := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgVolumeId))
		labelKey := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgLabelKey))
		c.Printer.Verbose("Removing label with key: %v for Volume with id: %v...", labelKey, volumeId)
		resp, err = c.CloudApiV5Services.Labels().VolumeDelete(dcId, volumeId, labelKey)
		if resp != nil {
			c.Printer.Verbose(config.RequestTimeMessage, resp.RequestTime)
		}
		if err != nil {
			return err
		}
	}
	return c.Printer.Print(getLabelResourcePrint(c, nil))
}

func RemoveAllVolumeLabels(c *core.CommandConfig) (*resources.Response, error) {
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgDataCenterId))
	volumeId := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgVolumeId))
	_ = c.Printer.Print("Labels to be removed from Volume with Id: " + volumeId)
	labels, resp, err := c.CloudApiV5Services.Labels().VolumeList(dcId, volumeId)
	if err != nil {
		return nil, err
	}
	if labelsItems, ok := labels.GetItemsOk(); ok && labelsItems != nil {
		for _, label := range *labelsItems {
			if properties, ok := label.GetPropertiesOk(); ok && properties != nil {
				if key, ok := properties.GetKeyOk(); ok && key != nil {
					_ = c.Printer.Print(" Label Key: " + *key)
				}
				if value, ok := properties.GetValueOk(); ok && value != nil {
					_ = c.Printer.Print(" Label Value: " + *value)
				}
			}
		}

		if err := utils.AskForConfirm(c.Stdin, c.Printer, "delete all the Labels from Volume with Id: "+volumeId); err != nil {
			return nil, err
		}
		c.Printer.Verbose("Deleting all the Labels from Volume with Id: %v...", volumeId)
		for _, label := range *labelsItems {
			if properties, ok := label.GetPropertiesOk(); ok && properties != nil {
				if key, ok := properties.GetKeyOk(); ok && key != nil {
					c.Printer.Verbose("Starting deleting Label with id: %v...", *key)
					resp, err = c.CloudApiV5Services.Labels().VolumeDelete(dcId, volumeId, *key)
					if resp != nil {
						c.Printer.Verbose("Request Id: %v", printer.GetId(resp))
						c.Printer.Verbose(config.RequestTimeMessage, resp.RequestTime)
					}
					if err != nil {
						return nil, err
					}
					if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
						return nil, err
					}
				}

			}
		}
	}
	return resp, nil
}

func RunIpBlockLabelsList(c *core.CommandConfig) error {
	labelDcs, resp, err := c.CloudApiV5Services.Labels().IpBlockList(viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgIpBlockId)))
	if resp != nil {
		c.Printer.Verbose(config.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, getLabelResources(labelDcs)))
}

func RunIpBlockLabelGet(c *core.CommandConfig) error {
	ipBlockId := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgIpBlockId))
	labelKey := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgLabelKey))
	c.Printer.Verbose("Getting label with key: %v for IpBlock with id: %v...", labelKey, ipBlockId)
	labelDc, resp, err := c.CloudApiV5Services.Labels().IpBlockGet(ipBlockId, labelKey)
	if resp != nil {
		c.Printer.Verbose(config.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, getLabelResource(labelDc)))
}

func RunIpBlockLabelAdd(c *core.CommandConfig) error {
	ipBlockId := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgIpBlockId))
	labelKey := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgLabelKey))
	labelValue := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgLabelValue))
	c.Printer.Verbose("Adding label with key: %v and value: %v to IpBlock with id: %v...", labelKey, labelValue, ipBlockId)
	labelDc, resp, err := c.CloudApiV5Services.Labels().IpBlockCreate(
		ipBlockId,
		labelKey,
		labelValue,
	)
	if resp != nil {
		c.Printer.Verbose(config.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, getLabelResource(labelDc)))
}

func RunIpBlockLabelRemove(c *core.CommandConfig) error {
	var resp *resources.Response
	var err error
	allFlag := viper.GetBool(core.GetFlagName(c.NS, cloudapiv5.ArgAll))
	if allFlag {
		resp, err = RemoveAllIpBlockLabels(c)
		if err != nil {
			return err
		}
	} else {
		ipBlockId := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgIpBlockId))
		labelKey := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgLabelKey))
		c.Printer.Verbose("Removing label with key: %v for IpBlock with id: %v...", labelKey, ipBlockId)
		resp, err = c.CloudApiV5Services.Labels().IpBlockDelete(ipBlockId, labelKey)
		if resp != nil {
			c.Printer.Verbose(config.RequestTimeMessage, resp.RequestTime)
		}
		if err != nil {
			return err
		}
	}
	return c.Printer.Print(getLabelResourcePrint(c, nil))
}

func RemoveAllIpBlockLabels(c *core.CommandConfig) (*resources.Response, error) {
	ipBlockId := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgIpBlockId))
	_ = c.Printer.Print("Labels to be removed from IpBlock with Id: " + ipBlockId)
	labels, resp, err := c.CloudApiV5Services.Labels().IpBlockList(ipBlockId)
	if err != nil {
		return nil, err
	}
	if labelsItems, ok := labels.GetItemsOk(); ok && labelsItems != nil {
		for _, label := range *labelsItems {
			if properties, ok := label.GetPropertiesOk(); ok && properties != nil {
				if key, ok := properties.GetKeyOk(); ok && key != nil {
					_ = c.Printer.Print(" Label Key: " + *key)
				}
				if value, ok := properties.GetValueOk(); ok && value != nil {
					_ = c.Printer.Print(" Label Value: " + *value)
				}
			}
		}

		if err := utils.AskForConfirm(c.Stdin, c.Printer, "delete all the Labels from IpBlock with Id: "+ipBlockId); err != nil {
			return nil, err
		}
		c.Printer.Verbose("Deleting all the Labels from IpBlock with Id: %v...", ipBlockId)
		for _, label := range *labelsItems {
			if properties, ok := label.GetPropertiesOk(); ok && properties != nil {
				if key, ok := properties.GetKeyOk(); ok && key != nil {
					c.Printer.Verbose("Starting deleting Label with id: %v...", *key)
					resp, err = c.CloudApiV5Services.Labels().IpBlockDelete(ipBlockId, *key)
					if resp != nil {
						c.Printer.Verbose("Request Id: %v", printer.GetId(resp))
						c.Printer.Verbose(config.RequestTimeMessage, resp.RequestTime)
					}
					if err != nil {
						return nil, err
					}
					if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
						return nil, err
					}
				}

			}
		}
	}
	return resp, nil
}

func RunSnapshotLabelsList(c *core.CommandConfig) error {
	labelDcs, resp, err := c.CloudApiV5Services.Labels().SnapshotList(viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgSnapshotId)))
	if resp != nil {
		c.Printer.Verbose(config.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, getLabelResources(labelDcs)))
}

func RunSnapshotLabelGet(c *core.CommandConfig) error {
	snapshotId := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgSnapshotId))
	labelKey := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgLabelKey))
	c.Printer.Verbose("Getting label with key: %v for Snapshot with id: %v...", labelKey, snapshotId)
	labelDc, resp, err := c.CloudApiV5Services.Labels().SnapshotGet(snapshotId, labelKey)
	if resp != nil {
		c.Printer.Verbose(config.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, getLabelResource(labelDc)))
}

func RunSnapshotLabelAdd(c *core.CommandConfig) error {
	snapshotId := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgSnapshotId))
	labelKey := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgLabelKey))
	labelValue := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgLabelValue))
	c.Printer.Verbose("Adding label with key: %v and value: %v to Snapshot with id: %v...", labelKey, labelValue, snapshotId)
	labelDc, resp, err := c.CloudApiV5Services.Labels().SnapshotCreate(snapshotId, labelKey, labelValue)
	if resp != nil {
		c.Printer.Verbose(config.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, getLabelResource(labelDc)))
}

func RunSnapshotLabelRemove(c *core.CommandConfig) error {
	var resp *resources.Response
	var err error
	allFlag := viper.GetBool(core.GetFlagName(c.NS, cloudapiv5.ArgAll))
	if allFlag {
		resp, err = RemoveAllSnapshotLabels(c)
		if err != nil {
			return err
		}
	} else {
		snapshotId := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgSnapshotId))
		labelKey := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgLabelKey))
		c.Printer.Verbose("Removing label with key: %v for Snapshot with id: %v...", labelKey, snapshotId)
		resp, err = c.CloudApiV5Services.Labels().SnapshotDelete(snapshotId, labelKey)
		if resp != nil {
			c.Printer.Verbose(config.RequestTimeMessage, resp.RequestTime)
		}
		if err != nil {
			return err
		}
	}
	return c.Printer.Print(getLabelResourcePrint(c, nil))
}

func RemoveAllSnapshotLabels(c *core.CommandConfig) (*resources.Response, error) {
	snapshotId := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgSnapshotId))
	_ = c.Printer.Print("Labels to be removed from Snapshot with Id: " + snapshotId)
	labels, resp, err := c.CloudApiV5Services.Labels().SnapshotList(snapshotId)
	if err != nil {
		return nil, err
	}
	if labelsItems, ok := labels.GetItemsOk(); ok && labelsItems != nil {
		for _, label := range *labelsItems {
			if properties, ok := label.GetPropertiesOk(); ok && properties != nil {
				if key, ok := properties.GetKeyOk(); ok && key != nil {
					_ = c.Printer.Print(" Label Key: " + *key)
				}
				if value, ok := properties.GetValueOk(); ok && value != nil {
					_ = c.Printer.Print(" Label Value: " + *value)
				}
			}
		}

		if err := utils.AskForConfirm(c.Stdin, c.Printer, "delete all the Labels from Snapshot with Id: "+snapshotId); err != nil {
			return nil, err
		}
		c.Printer.Verbose("Deleting all the Labels from Snapshot with Id: %v...", snapshotId)
		for _, label := range *labelsItems {
			if properties, ok := label.GetPropertiesOk(); ok && properties != nil {
				if key, ok := properties.GetKeyOk(); ok && key != nil {
					c.Printer.Verbose("Starting deleting Label with id: %v...", *key)
					resp, err = c.CloudApiV5Services.Labels().SnapshotDelete(snapshotId, *key)
					if resp != nil {
						c.Printer.Verbose("Request Id: %v", printer.GetId(resp))
						c.Printer.Verbose(config.RequestTimeMessage, resp.RequestTime)
					}
					if err != nil {
						return nil, err
					}
					if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
						return nil, err
					}
				}

			}
		}
	}
	return resp, nil
}

// Output Printing

var defaultLabelResourceCols = []string{"Key", "Value"}

type LabelResourcePrint struct {
	Key   string `json:"Key,omitempty"`
	Value string `json:"Value,omitempty"`
}

func getLabelResourcePrint(c *core.CommandConfig, s []resources.LabelResource) printer.Result {
	r := printer.Result{}
	if c != nil {
		if s != nil {
			r.OutputJSON = s
			r.KeyValue = getLabelResourcesKVMaps(s)
			r.Columns = defaultLabelResourceCols
		}
	}
	return r
}

func getLabelResources(labelResources resources.LabelResources) []resources.LabelResource {
	ss := make([]resources.LabelResource, 0)
	if items, ok := labelResources.GetItemsOk(); ok && items != nil {
		for _, s := range *items {
			ss = append(ss, resources.LabelResource{LabelResource: s})
		}
	}
	return ss
}

func getLabelResource(s *resources.LabelResource) []resources.LabelResource {
	ss := make([]resources.LabelResource, 0)
	if s != nil {
		ss = append(ss, resources.LabelResource{LabelResource: s.LabelResource})
	}
	return ss
}

func getLabelResourcesKVMaps(ss []resources.LabelResource) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(ss))
	for _, s := range ss {
		o := getLabelResourceKVMap(s)
		out = append(out, o)
	}
	return out
}

func getLabelResourceKVMap(s resources.LabelResource) map[string]interface{} {
	var ssPrint LabelResourcePrint
	if properties, ok := s.GetPropertiesOk(); ok && properties != nil {
		if key, ok := properties.GetKeyOk(); ok && key != nil {
			ssPrint.Key = *key
		}
		if value, ok := properties.GetValueOk(); ok && value != nil {
			ssPrint.Value = *value
		}
	}
	return structs.Map(ssPrint)
}
