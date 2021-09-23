package commands

import (
	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/internal/core"
	"github.com/ionos-cloud/ionosctl/internal/printer"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/services/cloudapi-v6/resources"
	"github.com/spf13/viper"
)

func RunDataCenterLabelsList(c *core.CommandConfig) error {
	labelDcs, resp, err := c.CloudApiV6Services.Labels().DatacenterList(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)))
	if resp != nil {
		c.Printer.Verbose(cloudapiv6.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, getLabelResources(labelDcs)))
}

func RunDataCenterLabelGet(c *core.CommandConfig) error {
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))
	labelKey := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLabelKey))
	c.Printer.Verbose("Getting label with key: %v for Datacenter with id: %v...", labelKey, dcId)
	labelDc, resp, err := c.CloudApiV6Services.Labels().DatacenterGet(dcId, labelKey)
	if resp != nil {
		c.Printer.Verbose(cloudapiv6.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, getLabelResource(labelDc)))
}

func RunDataCenterLabelAdd(c *core.CommandConfig) error {
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))
	labelKey := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLabelKey))
	labelValue := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLabelValue))
	c.Printer.Verbose("Adding label with key: %v and value: %v to Datacenter with id: %v...", labelKey, labelValue, dcId)
	labelDc, resp, err := c.CloudApiV6Services.Labels().DatacenterCreate(dcId, labelKey, labelValue)
	if resp != nil {
		c.Printer.Verbose(cloudapiv6.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, getLabelResource(labelDc)))
}

func RunDataCenterLabelRemove(c *core.CommandConfig) error {
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))
	labelKey := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLabelKey))
	c.Printer.Verbose("Removing label with key: %v for Datacenter with id: %v...", labelKey, dcId)
	resp, err := c.CloudApiV6Services.Labels().DatacenterDelete(dcId, labelKey)
	if resp != nil {
		c.Printer.Verbose(cloudapiv6.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, nil))
}

func RunServerLabelsList(c *core.CommandConfig) error {
	labelDcs, resp, err := c.CloudApiV6Services.Labels().ServerList(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgServerId)),
	)
	if resp != nil {
		c.Printer.Verbose(cloudapiv6.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, getLabelResources(labelDcs)))
}

func RunServerLabelGet(c *core.CommandConfig) error {
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))
	serverId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgServerId))
	labelkey := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLabelKey))
	c.Printer.Verbose("Getting label with key: %v for Server with id: %v...", labelkey, serverId)
	labelDc, resp, err := c.CloudApiV6Services.Labels().ServerGet(dcId, serverId, labelkey)
	if resp != nil {
		c.Printer.Verbose(cloudapiv6.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, getLabelResource(labelDc)))
}

func RunServerLabelAdd(c *core.CommandConfig) error {
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))
	serverId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgServerId))
	labelKey := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLabelKey))
	labelValue := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLabelValue))
	c.Printer.Verbose("Adding label with key: %v and value: %v to Server with id: %v...", labelKey, labelValue, serverId)
	labelDc, resp, err := c.CloudApiV6Services.Labels().ServerCreate(dcId, serverId, labelKey, labelValue)
	if resp != nil {
		c.Printer.Verbose(cloudapiv6.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, getLabelResource(labelDc)))
}

func RunServerLabelRemove(c *core.CommandConfig) error {
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))
	serverId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgServerId))
	labelKey := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLabelKey))
	c.Printer.Verbose("Removing label with key: %v for Server with id: %v...", labelKey, serverId)
	resp, err := c.CloudApiV6Services.Labels().ServerDelete(dcId, serverId, labelKey)
	if resp != nil {
		c.Printer.Verbose(cloudapiv6.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, nil))
}

func RunVolumeLabelsList(c *core.CommandConfig) error {
	labelDcs, resp, err := c.CloudApiV6Services.Labels().VolumeList(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgVolumeId)),
	)
	if resp != nil {
		c.Printer.Verbose(cloudapiv6.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, getLabelResources(labelDcs)))
}

func RunVolumeLabelGet(c *core.CommandConfig) error {
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))
	volumeId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgVolumeId))
	labelKey := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLabelKey))
	c.Printer.Verbose("Getting label with key: %v for Volume with id: %v...", labelKey, volumeId)
	labelDc, resp, err := c.CloudApiV6Services.Labels().VolumeGet(dcId, volumeId, labelKey)
	if resp != nil {
		c.Printer.Verbose(cloudapiv6.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, getLabelResource(labelDc)))
}

func RunVolumeLabelAdd(c *core.CommandConfig) error {
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))
	volumeId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgVolumeId))
	labelKey := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLabelKey))
	labelValue := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLabelValue))
	c.Printer.Verbose("Adding label with key: %v and value: %v to Volume with id: %v...", labelKey, labelValue, volumeId)
	labelDc, resp, err := c.CloudApiV6Services.Labels().VolumeCreate(dcId, volumeId, labelKey, labelValue)
	if resp != nil {
		c.Printer.Verbose(cloudapiv6.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, getLabelResource(labelDc)))
}

func RunVolumeLabelRemove(c *core.CommandConfig) error {
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))
	volumeId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgVolumeId))
	labelKey := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLabelKey))
	c.Printer.Verbose("Removing label with key: %v for Volume with id: %v...", labelKey, volumeId)
	resp, err := c.CloudApiV6Services.Labels().VolumeDelete(dcId, volumeId, labelKey)
	if resp != nil {
		c.Printer.Verbose(cloudapiv6.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, nil))
}

func RunIpBlockLabelsList(c *core.CommandConfig) error {
	labelDcs, resp, err := c.CloudApiV6Services.Labels().IpBlockList(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgIpBlockId)))
	if resp != nil {
		c.Printer.Verbose(cloudapiv6.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, getLabelResources(labelDcs)))
}

func RunIpBlockLabelGet(c *core.CommandConfig) error {
	ipBlockId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgIpBlockId))
	labelKey := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLabelKey))
	c.Printer.Verbose("Getting label with key: %v for IpBlock with id: %v...", labelKey, ipBlockId)
	labelDc, resp, err := c.CloudApiV6Services.Labels().IpBlockGet(ipBlockId, labelKey)
	if resp != nil {
		c.Printer.Verbose(cloudapiv6.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, getLabelResource(labelDc)))
}

func RunIpBlockLabelAdd(c *core.CommandConfig) error {
	ipBlockId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgIpBlockId))
	labelKey := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLabelKey))
	labelValue := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLabelValue))
	c.Printer.Verbose("Adding label with key: %v and value: %v to IpBlock with id: %v...", labelKey, labelValue, ipBlockId)
	labelDc, resp, err := c.CloudApiV6Services.Labels().IpBlockCreate(ipBlockId, labelKey, labelValue)
	if resp != nil {
		c.Printer.Verbose(cloudapiv6.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, getLabelResource(labelDc)))
}

func RunIpBlockLabelRemove(c *core.CommandConfig) error {
	ipBlockId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgIpBlockId))
	labelKey := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLabelKey))
	c.Printer.Verbose("Removing label with key: %v for IpBlock with id: %v...", labelKey, ipBlockId)
	resp, err := c.CloudApiV6Services.Labels().IpBlockDelete(ipBlockId, labelKey)
	if resp != nil {
		c.Printer.Verbose(cloudapiv6.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, nil))
}

func RunSnapshotLabelsList(c *core.CommandConfig) error {
	labelDcs, resp, err := c.CloudApiV6Services.Labels().SnapshotList(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgSnapshotId)))
	if resp != nil {
		c.Printer.Verbose(cloudapiv6.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, getLabelResources(labelDcs)))
}

func RunSnapshotLabelGet(c *core.CommandConfig) error {
	snapshotId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgSnapshotId))
	labelKey := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLabelKey))
	c.Printer.Verbose("Getting label with key: %v for Snapshot with id: %v...", labelKey, snapshotId)
	labelDc, resp, err := c.CloudApiV6Services.Labels().SnapshotGet(snapshotId, labelKey)
	if resp != nil {
		c.Printer.Verbose(cloudapiv6.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, getLabelResource(labelDc)))
}

func RunSnapshotLabelAdd(c *core.CommandConfig) error {
	snapshotId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgSnapshotId))
	labelKey := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLabelKey))
	labelValue := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLabelValue))
	c.Printer.Verbose("Adding label with key: %v and value: %v to Snapshot with id: %v...", labelKey, labelValue, snapshotId)
	labelDc, resp, err := c.CloudApiV6Services.Labels().SnapshotCreate(snapshotId, labelKey, labelValue)
	if resp != nil {
		c.Printer.Verbose(cloudapiv6.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, getLabelResource(labelDc)))
}

func RunSnapshotLabelRemove(c *core.CommandConfig) error {
	snapshotId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgSnapshotId))
	labelKey := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLabelKey))
	c.Printer.Verbose("Removing label with key: %v for Snapshot with id: %v...", labelKey, snapshotId)
	resp, err := c.CloudApiV6Services.Labels().SnapshotDelete(snapshotId, labelKey)
	if resp != nil {
		c.Printer.Verbose(cloudapiv6.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, nil))
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
