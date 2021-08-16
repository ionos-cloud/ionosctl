package commands

import (
	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/resources/v5"
	"github.com/ionos-cloud/ionosctl/pkg/utils/printer"
	"github.com/spf13/viper"
)

func RunDataCenterLabelsList(c *core.CommandConfig) error {
	labelDcs, resp, err := c.Labels().DatacenterList(viper.GetString(core.GetGlobalFlagName(c.Resource, config.ArgDataCenterId)))
	if resp != nil {
		c.Printer.Verbose("the execution time of the command is: %v", resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, getLabelResources(labelDcs)))
}

func RunDataCenterLabelGet(c *core.CommandConfig) error {
	dcId := viper.GetString(core.GetGlobalFlagName(c.Resource, config.ArgDataCenterId))
	labelKey := viper.GetString(core.GetFlagName(c.NS, config.ArgLabelKey))
	c.Printer.Verbose("Getting label with key: %v for Datacenter with id: %v...", labelKey, dcId)
	labelDc, resp, err := c.Labels().DatacenterGet(dcId, labelKey)
	if resp != nil {
		c.Printer.Verbose("the execution time of the command is: %v", resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, getLabelResource(labelDc)))
}

func RunDataCenterLabelAdd(c *core.CommandConfig) error {
	dcId := viper.GetString(core.GetGlobalFlagName(c.Resource, config.ArgDataCenterId))
	labelKey := viper.GetString(core.GetFlagName(c.NS, config.ArgLabelKey))
	labelValue := viper.GetString(core.GetFlagName(c.NS, config.ArgLabelValue))
	c.Printer.Verbose("Adding label with key: %v and value: %v to Datacenter with id: %v...", labelKey, labelValue, dcId)
	labelDc, resp, err := c.Labels().DatacenterCreate(dcId, labelKey, labelValue)
	if resp != nil {
		c.Printer.Verbose("the execution time of the command is: %v", resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, getLabelResource(labelDc)))
}

func RunDataCenterLabelRemove(c *core.CommandConfig) error {
	dcId := viper.GetString(core.GetGlobalFlagName(c.Resource, config.ArgDataCenterId))
	labelKey := viper.GetString(core.GetFlagName(c.NS, config.ArgLabelKey))
	c.Printer.Verbose("Removing label with key: %v for Datacenter with id: %v...", labelKey, dcId)
	resp, err := c.Labels().DatacenterDelete(dcId, labelKey)
	if resp != nil {
		c.Printer.Verbose("the execution time of the command is: %v", resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, nil))
}

func RunServerLabelsList(c *core.CommandConfig) error {
	labelDcs, resp, err := c.Labels().ServerList(
		viper.GetString(core.GetGlobalFlagName(c.Resource, config.ArgDataCenterId)),
		viper.GetString(core.GetGlobalFlagName(c.Resource, config.ArgServerId)),
	)
	if resp != nil {
		c.Printer.Verbose("the execution time of the command is: %v", resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, getLabelResources(labelDcs)))
}

func RunServerLabelGet(c *core.CommandConfig) error {
	dcId := viper.GetString(core.GetGlobalFlagName(c.Resource, config.ArgDataCenterId))
	serverId := viper.GetString(core.GetGlobalFlagName(c.Resource, config.ArgServerId))
	labelkey := viper.GetString(core.GetFlagName(c.NS, config.ArgLabelKey))
	c.Printer.Verbose("Getting label with key: %v for Server with id: %v...", labelkey, serverId)
	labelDc, resp, err := c.Labels().ServerGet(dcId, serverId, labelkey)
	if resp != nil {
		c.Printer.Verbose("the execution time of the command is: %v", resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, getLabelResource(labelDc)))
}

func RunServerLabelAdd(c *core.CommandConfig) error {
	dcId := viper.GetString(core.GetGlobalFlagName(c.Resource, config.ArgDataCenterId))
	serverId := viper.GetString(core.GetGlobalFlagName(c.Resource, config.ArgServerId))
	labelKey := viper.GetString(core.GetFlagName(c.NS, config.ArgLabelKey))
	labelValue := viper.GetString(core.GetFlagName(c.NS, config.ArgLabelValue))
	c.Printer.Verbose("Adding label with key: %v and value: %v to Server with id: %v...", labelKey, labelValue, serverId)
	labelDc, resp, err := c.Labels().ServerCreate(dcId, serverId, labelKey, labelValue)
	if resp != nil {
		c.Printer.Verbose("the execution time of the command is: %v", resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, getLabelResource(labelDc)))
}

func RunServerLabelRemove(c *core.CommandConfig) error {
	dcId := viper.GetString(core.GetGlobalFlagName(c.Resource, config.ArgDataCenterId))
	serverId := viper.GetString(core.GetGlobalFlagName(c.Resource, config.ArgServerId))
	labelKey := viper.GetString(core.GetFlagName(c.NS, config.ArgLabelKey))
	c.Printer.Verbose("Removing label with key: %v for Server with id: %v...", labelKey, serverId)
	resp, err := c.Labels().ServerDelete(dcId, serverId, labelKey)
	if resp != nil {
		c.Printer.Verbose("the execution time of the command is: %v", resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, nil))
}

func RunVolumeLabelsList(c *core.CommandConfig) error {
	labelDcs, resp, err := c.Labels().VolumeList(
		viper.GetString(core.GetGlobalFlagName(c.Resource, config.ArgDataCenterId)),
		viper.GetString(core.GetGlobalFlagName(c.Resource, config.ArgVolumeId)),
	)
	if resp != nil {
		c.Printer.Verbose("the execution time of the command is: %v", resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, getLabelResources(labelDcs)))
}

func RunVolumeLabelGet(c *core.CommandConfig) error {
	dcId := viper.GetString(core.GetGlobalFlagName(c.Resource, config.ArgDataCenterId))
	volumeId := viper.GetString(core.GetGlobalFlagName(c.Resource, config.ArgVolumeId))
	labelKey := viper.GetString(core.GetFlagName(c.NS, config.ArgLabelKey))
	c.Printer.Verbose("Getting label with key: %v for Volume with id: %v...", labelKey, volumeId)
	labelDc, resp, err := c.Labels().VolumeGet(dcId, volumeId, labelKey)
	if resp != nil {
		c.Printer.Verbose("the execution time of the command is: %v", resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, getLabelResource(labelDc)))
}

func RunVolumeLabelAdd(c *core.CommandConfig) error {
	dcId := viper.GetString(core.GetGlobalFlagName(c.Resource, config.ArgDataCenterId))
	volumeId := viper.GetString(core.GetGlobalFlagName(c.Resource, config.ArgVolumeId))
	labelKey := viper.GetString(core.GetFlagName(c.NS, config.ArgLabelKey))
	labelValue := viper.GetString(core.GetFlagName(c.NS, config.ArgLabelValue))
	c.Printer.Verbose("Adding label with key: %v and value: %v to Volume with id: %v...", labelKey, labelValue, volumeId)
	labelDc, resp, err := c.Labels().VolumeCreate(dcId, volumeId, labelKey, labelValue)
	if resp != nil {
		c.Printer.Verbose("the execution time of the command is: %v", resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, getLabelResource(labelDc)))
}

func RunVolumeLabelRemove(c *core.CommandConfig) error {
	dcId := viper.GetString(core.GetGlobalFlagName(c.Resource, config.ArgDataCenterId))
	volumeId := viper.GetString(core.GetGlobalFlagName(c.Resource, config.ArgVolumeId))
	labelKey := viper.GetString(core.GetFlagName(c.NS, config.ArgLabelKey))
	c.Printer.Verbose("Removing label with key: %v for Volume with id: %v...", labelKey, volumeId)
	resp, err := c.Labels().VolumeDelete(dcId, volumeId, labelKey)
	if resp != nil {
		c.Printer.Verbose("the execution time of the command is: %v", resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, nil))
}

func RunIpBlockLabelsList(c *core.CommandConfig) error {
	labelDcs, resp, err := c.Labels().IpBlockList(viper.GetString(core.GetGlobalFlagName(c.Resource, config.ArgIpBlockId)))
	if resp != nil {
		c.Printer.Verbose("the execution time of the command is: %v", resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, getLabelResources(labelDcs)))
}

func RunIpBlockLabelGet(c *core.CommandConfig) error {
	ipBlockId := viper.GetString(core.GetGlobalFlagName(c.Resource, config.ArgIpBlockId))
	labelKey := viper.GetString(core.GetFlagName(c.NS, config.ArgLabelKey))
	c.Printer.Verbose("Getting label with key: %v for IpBlock with id: %v...", labelKey, ipBlockId)
	labelDc, resp, err := c.Labels().IpBlockGet(ipBlockId, labelKey)
	if resp != nil {
		c.Printer.Verbose("the execution time of the command is: %v", resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, getLabelResource(labelDc)))
}

func RunIpBlockLabelAdd(c *core.CommandConfig) error {
	ipBlockId := viper.GetString(core.GetGlobalFlagName(c.Resource, config.ArgIpBlockId))
	labelKey := viper.GetString(core.GetFlagName(c.NS, config.ArgLabelKey))
	labelValue := viper.GetString(core.GetFlagName(c.NS, config.ArgLabelValue))
	c.Printer.Verbose("Adding label with key: %v and value: %v to IpBlock with id: %v...", labelKey, labelValue, ipBlockId)
	labelDc, resp, err := c.Labels().IpBlockCreate(
		ipBlockId,
		labelKey,
		labelValue,
	)
	if resp != nil {
		c.Printer.Verbose("the execution time of the command is: %v", resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, getLabelResource(labelDc)))
}

func RunIpBlockLabelRemove(c *core.CommandConfig) error {
	ipBlockId := viper.GetString(core.GetGlobalFlagName(c.Resource, config.ArgIpBlockId))
	labelKey := viper.GetString(core.GetFlagName(c.NS, config.ArgLabelKey))
	c.Printer.Verbose("Removing label with key: %v for IpBlock with id: %v...", labelKey, ipBlockId)
	resp, err := c.Labels().IpBlockDelete(ipBlockId, labelKey)
	if resp != nil {
		c.Printer.Verbose("the execution time of the command is: %v", resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, nil))
}

func RunSnapshotLabelsList(c *core.CommandConfig) error {
	labelDcs, resp, err := c.Labels().SnapshotList(viper.GetString(core.GetGlobalFlagName(c.Resource, config.ArgSnapshotId)))
	if resp != nil {
		c.Printer.Verbose("the execution time of the command is: %v", resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, getLabelResources(labelDcs)))
}

func RunSnapshotLabelGet(c *core.CommandConfig) error {
	snapshotId := viper.GetString(core.GetGlobalFlagName(c.Resource, config.ArgSnapshotId))
	labelKey := viper.GetString(core.GetFlagName(c.NS, config.ArgLabelKey))
	c.Printer.Verbose("Getting label with key: %v for Snapshot with id: %v...", labelKey, snapshotId)
	labelDc, resp, err := c.Labels().SnapshotGet(snapshotId, labelKey)
	if resp != nil {
		c.Printer.Verbose("the execution time of the command is: %v", resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, getLabelResource(labelDc)))
}

func RunSnapshotLabelAdd(c *core.CommandConfig) error {
	snapshotId := viper.GetString(core.GetGlobalFlagName(c.Resource, config.ArgSnapshotId))
	labelKey := viper.GetString(core.GetFlagName(c.NS, config.ArgLabelKey))
	labelValue := viper.GetString(core.GetFlagName(c.NS, config.ArgLabelValue))
	c.Printer.Verbose("Adding label with key: %v and value: %v to Snapshot with id: %v...", labelKey, labelValue, snapshotId)
	labelDc, resp, err := c.Labels().SnapshotCreate(snapshotId, labelKey, labelValue)
	if resp != nil {
		c.Printer.Verbose("the execution time of the command is: %v", resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, getLabelResource(labelDc)))
}

func RunSnapshotLabelRemove(c *core.CommandConfig) error {
	snapshotId := viper.GetString(core.GetGlobalFlagName(c.Resource, config.ArgSnapshotId))
	labelKey := viper.GetString(core.GetFlagName(c.NS, config.ArgLabelKey))
	c.Printer.Verbose("Removing label with key: %v for Snapshot with id: %v...", labelKey, snapshotId)
	resp, err := c.Labels().SnapshotDelete(snapshotId, labelKey)
	if resp != nil {
		c.Printer.Verbose("the execution time of the command is: %v", resp.RequestTime)
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

func getLabelResourcePrint(c *core.CommandConfig, s []v5.LabelResource) printer.Result {
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

func getLabelResources(labelResources v5.LabelResources) []v5.LabelResource {
	ss := make([]v5.LabelResource, 0)
	if items, ok := labelResources.GetItemsOk(); ok && items != nil {
		for _, s := range *items {
			ss = append(ss, v5.LabelResource{LabelResource: s})
		}
	}
	return ss
}

func getLabelResource(s *v5.LabelResource) []v5.LabelResource {
	ss := make([]v5.LabelResource, 0)
	if s != nil {
		ss = append(ss, v5.LabelResource{LabelResource: s.LabelResource})
	}
	return ss
}

func getLabelResourcesKVMaps(ss []v5.LabelResource) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(ss))
	for _, s := range ss {
		o := getLabelResourceKVMap(s)
		out = append(out, o)
	}
	return out
}

func getLabelResourceKVMap(s v5.LabelResource) map[string]interface{} {
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
