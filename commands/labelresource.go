package commands

import (
	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/resources/v6"
	"github.com/ionos-cloud/ionosctl/pkg/utils/printer"
	"github.com/spf13/viper"
)

func RunDataCenterLabelsList(c *core.CommandConfig) error {
	labelDcs, _, err := c.Labels().DatacenterList(viper.GetString(core.GetFlagName(c.NS, config.ArgDataCenterId)))
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, getLabelResources(labelDcs)))
}

func RunDataCenterLabelGet(c *core.CommandConfig) error {
	dcId := viper.GetString(core.GetFlagName(c.NS, config.ArgDataCenterId))
	labelKey := viper.GetString(core.GetFlagName(c.NS, config.ArgLabelKey))
	c.Printer.Verbose("Getting label with key: %v for Datacenter with id: %v...", labelKey, dcId)
	labelDc, _, err := c.Labels().DatacenterGet(dcId, labelKey)
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, getLabelResource(labelDc)))
}

func RunDataCenterLabelAdd(c *core.CommandConfig) error {
	dcId := viper.GetString(core.GetFlagName(c.NS, config.ArgDataCenterId))
	labelKey := viper.GetString(core.GetFlagName(c.NS, config.ArgLabelKey))
	labelValue := viper.GetString(core.GetFlagName(c.NS, config.ArgLabelValue))
	c.Printer.Verbose("Adding label with key: %v and value: %v to Datacenter with id: %v...", labelKey, labelValue, dcId)
	labelDc, _, err := c.Labels().DatacenterCreate(dcId, labelKey, labelValue)
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, getLabelResource(labelDc)))
}

func RunDataCenterLabelRemove(c *core.CommandConfig) error {
	dcId := viper.GetString(core.GetFlagName(c.NS, config.ArgDataCenterId))
	labelKey := viper.GetString(core.GetFlagName(c.NS, config.ArgLabelKey))
	c.Printer.Verbose("Removing label with key: %v for Datacenter with id: %v...", labelKey, dcId)
	_, err := c.Labels().DatacenterDelete(dcId, labelKey)
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, nil))
}

func RunServerLabelsList(c *core.CommandConfig) error {
	labelDcs, _, err := c.Labels().ServerList(
		viper.GetString(core.GetFlagName(c.NS, config.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgServerId)),
	)
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, getLabelResources(labelDcs)))
}

func RunServerLabelGet(c *core.CommandConfig) error {
	dcId := viper.GetString(core.GetFlagName(c.NS, config.ArgDataCenterId))
	serverId := viper.GetString(core.GetFlagName(c.NS, config.ArgServerId))
	labelkey := viper.GetString(core.GetFlagName(c.NS, config.ArgLabelKey))
	c.Printer.Verbose("Getting label with key: %v for Server with id: %v...", labelkey, serverId)
	labelDc, _, err := c.Labels().ServerGet(dcId, serverId, labelkey)
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, getLabelResource(labelDc)))
}

func RunServerLabelAdd(c *core.CommandConfig) error {
	dcId := viper.GetString(core.GetFlagName(c.NS, config.ArgDataCenterId))
	serverId := viper.GetString(core.GetFlagName(c.NS, config.ArgServerId))
	labelKey := viper.GetString(core.GetFlagName(c.NS, config.ArgLabelKey))
	labelValue := viper.GetString(core.GetFlagName(c.NS, config.ArgLabelValue))
	c.Printer.Verbose("Adding label with key: %v and value: %v to Server with id: %v...", labelKey, labelValue, serverId)
	labelDc, _, err := c.Labels().ServerCreate(dcId, serverId, labelKey, labelValue)
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, getLabelResource(labelDc)))
}

func RunServerLabelRemove(c *core.CommandConfig) error {
	dcId := viper.GetString(core.GetFlagName(c.NS, config.ArgDataCenterId))
	serverId := viper.GetString(core.GetFlagName(c.NS, config.ArgServerId))
	labelKey := viper.GetString(core.GetFlagName(c.NS, config.ArgLabelKey))
	c.Printer.Verbose("Removing label with key: %v for Server with id: %v...", labelKey, serverId)
	_, err := c.Labels().ServerDelete(dcId, serverId, labelKey)
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, nil))
}

func RunVolumeLabelsList(c *core.CommandConfig) error {
	labelDcs, _, err := c.Labels().VolumeList(
		viper.GetString(core.GetFlagName(c.NS, config.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgVolumeId)),
	)
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, getLabelResources(labelDcs)))
}

func RunVolumeLabelGet(c *core.CommandConfig) error {
	dcId := viper.GetString(core.GetFlagName(c.NS, config.ArgDataCenterId))
	volumeId := viper.GetString(core.GetFlagName(c.NS, config.ArgVolumeId))
	labelKey := viper.GetString(core.GetFlagName(c.NS, config.ArgLabelKey))
	c.Printer.Verbose("Getting label with key: %v for Volume with id: %v...", labelKey, volumeId)
	labelDc, _, err := c.Labels().VolumeGet(dcId, volumeId, labelKey)
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, getLabelResource(labelDc)))
}

func RunVolumeLabelAdd(c *core.CommandConfig) error {
	dcId := viper.GetString(core.GetFlagName(c.NS, config.ArgDataCenterId))
	volumeId := viper.GetString(core.GetFlagName(c.NS, config.ArgVolumeId))
	labelKey := viper.GetString(core.GetFlagName(c.NS, config.ArgLabelKey))
	labelValue := viper.GetString(core.GetFlagName(c.NS, config.ArgLabelValue))
	c.Printer.Verbose("Adding label with key: %v and value: %v to Volume with id: %v...", labelKey, labelValue, volumeId)
	labelDc, _, err := c.Labels().VolumeCreate(dcId, volumeId, labelKey, labelValue)
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, getLabelResource(labelDc)))
}

func RunVolumeLabelRemove(c *core.CommandConfig) error {
	dcId := viper.GetString(core.GetFlagName(c.NS, config.ArgDataCenterId))
	volumeId := viper.GetString(core.GetFlagName(c.NS, config.ArgVolumeId))
	labelKey := viper.GetString(core.GetFlagName(c.NS, config.ArgLabelKey))
	c.Printer.Verbose("Removing label with key: %v for Volume with id: %v...", labelKey, volumeId)
	_, err := c.Labels().VolumeDelete(dcId, volumeId, labelKey)
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, nil))
}

func RunIpBlockLabelsList(c *core.CommandConfig) error {
	labelDcs, _, err := c.Labels().IpBlockList(viper.GetString(core.GetFlagName(c.NS, config.ArgIpBlockId)))
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, getLabelResources(labelDcs)))
}

func RunIpBlockLabelGet(c *core.CommandConfig) error {
	ipBlockId := viper.GetString(core.GetFlagName(c.NS, config.ArgIpBlockId))
	labelKey := viper.GetString(core.GetFlagName(c.NS, config.ArgLabelKey))
	c.Printer.Verbose("Getting label with key: %v for IpBlock with id: %v...", labelKey, ipBlockId)
	labelDc, _, err := c.Labels().IpBlockGet(ipBlockId, labelKey)
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, getLabelResource(labelDc)))
}

func RunIpBlockLabelAdd(c *core.CommandConfig) error {
	ipBlockId := viper.GetString(core.GetFlagName(c.NS, config.ArgIpBlockId))
	labelKey := viper.GetString(core.GetFlagName(c.NS, config.ArgLabelKey))
	labelValue := viper.GetString(core.GetFlagName(c.NS, config.ArgLabelValue))
	c.Printer.Verbose("Adding label with key: %v and value: %v to IpBlock with id: %v...", labelKey, labelValue, ipBlockId)
	labelDc, _, err := c.Labels().IpBlockCreate(ipBlockId, labelKey, labelValue)
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, getLabelResource(labelDc)))
}

func RunIpBlockLabelRemove(c *core.CommandConfig) error {
	ipBlockId := viper.GetString(core.GetFlagName(c.NS, config.ArgIpBlockId))
	labelKey := viper.GetString(core.GetFlagName(c.NS, config.ArgLabelKey))
	c.Printer.Verbose("Removing label with key: %v for IpBlock with id: %v...", labelKey, ipBlockId)
	_, err := c.Labels().IpBlockDelete(ipBlockId, labelKey)
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, nil))
}

func RunSnapshotLabelsList(c *core.CommandConfig) error {
	labelDcs, _, err := c.Labels().SnapshotList(viper.GetString(core.GetFlagName(c.NS, config.ArgSnapshotId)))
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, getLabelResources(labelDcs)))
}

func RunSnapshotLabelGet(c *core.CommandConfig) error {
	snapshotId := viper.GetString(core.GetFlagName(c.NS, config.ArgSnapshotId))
	labelKey := viper.GetString(core.GetFlagName(c.NS, config.ArgLabelKey))
	c.Printer.Verbose("Getting label with key: %v for Snapshot with id: %v...", labelKey, snapshotId)
	labelDc, _, err := c.Labels().SnapshotGet(snapshotId, labelKey)
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, getLabelResource(labelDc)))
}

func RunSnapshotLabelAdd(c *core.CommandConfig) error {
	snapshotId := viper.GetString(core.GetFlagName(c.NS, config.ArgSnapshotId))
	labelKey := viper.GetString(core.GetFlagName(c.NS, config.ArgLabelKey))
	labelValue := viper.GetString(core.GetFlagName(c.NS, config.ArgLabelValue))
	c.Printer.Verbose("Adding label with key: %v and value: %v to Snapshot with id: %v...", labelKey, labelValue, snapshotId)
	labelDc, _, err := c.Labels().SnapshotCreate(snapshotId, labelKey, labelValue)
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, getLabelResource(labelDc)))
}

func RunSnapshotLabelRemove(c *core.CommandConfig) error {
	snapshotId := viper.GetString(core.GetFlagName(c.NS, config.ArgSnapshotId))
	labelKey := viper.GetString(core.GetFlagName(c.NS, config.ArgLabelKey))
	c.Printer.Verbose("Removing label with key: %v for Snapshot with id: %v...", labelKey, snapshotId)
	_, err := c.Labels().SnapshotDelete(snapshotId, labelKey)
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

func getLabelResourcePrint(c *core.CommandConfig, s []v6.LabelResource) printer.Result {
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

func getLabelResources(labelResources v6.LabelResources) []v6.LabelResource {
	ss := make([]v6.LabelResource, 0)
	if items, ok := labelResources.GetItemsOk(); ok && items != nil {
		for _, s := range *items {
			ss = append(ss, v6.LabelResource{LabelResource: s})
		}
	}
	return ss
}

func getLabelResource(s *v6.LabelResource) []v6.LabelResource {
	ss := make([]v6.LabelResource, 0)
	if s != nil {
		ss = append(ss, v6.LabelResource{LabelResource: s.LabelResource})
	}
	return ss
}

func getLabelResourcesKVMaps(ss []v6.LabelResource) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(ss))
	for _, s := range ss {
		o := getLabelResourceKVMap(s)
		out = append(out, o)
	}
	return out
}

func getLabelResourceKVMap(s v6.LabelResource) map[string]interface{} {
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
