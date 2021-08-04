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
	labelDcs, _, err := c.Labels().DatacenterList(viper.GetString(core.GetGlobalFlagName(c.Resource, config.ArgDataCenterId)))
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, getLabelResources(labelDcs)))
}

func RunDataCenterLabelGet(c *core.CommandConfig) error {
	dcId := viper.GetString(core.GetGlobalFlagName(c.Resource, config.ArgDataCenterId))
	labelKey := viper.GetString(core.GetFlagName(c.NS, config.ArgLabelKey))
	c.Printer.Verbose("Getting label from Datacenter with id: %v...", dcId)
	labelDc, _, err := c.Labels().DatacenterGet(dcId, labelKey)
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, getLabelResource(labelDc)))
}

func RunDataCenterLabelAdd(c *core.CommandConfig) error {
	labelDc, _, err := c.Labels().DatacenterCreate(
		viper.GetString(core.GetGlobalFlagName(c.Resource, config.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgLabelKey)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgLabelValue)),
	)
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, getLabelResource(labelDc)))
}

func RunDataCenterLabelRemove(c *core.CommandConfig) error {
	_, err := c.Labels().DatacenterDelete(
		viper.GetString(core.GetGlobalFlagName(c.Resource, config.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgLabelKey)),
	)
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, nil))
}

func RunServerLabelsList(c *core.CommandConfig) error {
	labelDcs, _, err := c.Labels().ServerList(
		viper.GetString(core.GetGlobalFlagName(c.Resource, config.ArgDataCenterId)),
		viper.GetString(core.GetGlobalFlagName(c.Resource, config.ArgServerId)),
	)
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, getLabelResources(labelDcs)))
}

func RunServerLabelGet(c *core.CommandConfig) error {
	dcId := viper.GetString(core.GetGlobalFlagName(c.Resource, config.ArgDataCenterId))
	serverId := viper.GetString(core.GetGlobalFlagName(c.Resource, config.ArgServerId))
	labelkey := viper.GetString(core.GetFlagName(c.NS, config.ArgLabelKey))
	c.Printer.Verbose("Getting server with id: %v from Datacenter with id: %v...", serverId, dcId)
	labelDc, _, err := c.Labels().ServerGet(dcId, serverId, labelkey)
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, getLabelResource(labelDc)))
}

func RunServerLabelAdd(c *core.CommandConfig) error {
	labelDc, _, err := c.Labels().ServerCreate(
		viper.GetString(core.GetGlobalFlagName(c.Resource, config.ArgDataCenterId)),
		viper.GetString(core.GetGlobalFlagName(c.Resource, config.ArgServerId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgLabelKey)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgLabelValue)),
	)
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, getLabelResource(labelDc)))
}

func RunServerLabelRemove(c *core.CommandConfig) error {
	_, err := c.Labels().ServerDelete(
		viper.GetString(core.GetGlobalFlagName(c.Resource, config.ArgDataCenterId)),
		viper.GetString(core.GetGlobalFlagName(c.Resource, config.ArgServerId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgLabelKey)),
	)
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, nil))
}

func RunVolumeLabelsList(c *core.CommandConfig) error {
	labelDcs, _, err := c.Labels().VolumeList(
		viper.GetString(core.GetGlobalFlagName(c.Resource, config.ArgDataCenterId)),
		viper.GetString(core.GetGlobalFlagName(c.Resource, config.ArgVolumeId)),
	)
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, getLabelResources(labelDcs)))
}

func RunVolumeLabelGet(c *core.CommandConfig) error {
	dcId := viper.GetString(core.GetGlobalFlagName(c.Resource, config.ArgDataCenterId))
	volumeId := viper.GetString(core.GetGlobalFlagName(c.Resource, config.ArgVolumeId))
	labelKey := viper.GetString(core.GetFlagName(c.NS, config.ArgLabelKey))
	c.Printer.Verbose("Getting Volume with id: %v from Datacenter with id: %v...", volumeId, dcId)
	labelDc, _, err := c.Labels().VolumeGet(dcId, volumeId, labelKey)
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, getLabelResource(labelDc)))
}

func RunVolumeLabelAdd(c *core.CommandConfig) error {
	labelDc, _, err := c.Labels().VolumeCreate(
		viper.GetString(core.GetGlobalFlagName(c.Resource, config.ArgDataCenterId)),
		viper.GetString(core.GetGlobalFlagName(c.Resource, config.ArgVolumeId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgLabelKey)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgLabelValue)),
	)
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, getLabelResource(labelDc)))
}

func RunVolumeLabelRemove(c *core.CommandConfig) error {
	_, err := c.Labels().VolumeDelete(
		viper.GetString(core.GetGlobalFlagName(c.Resource, config.ArgDataCenterId)),
		viper.GetString(core.GetGlobalFlagName(c.Resource, config.ArgVolumeId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgLabelKey)),
	)
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, nil))
}

func RunIpBlockLabelsList(c *core.CommandConfig) error {
	labelDcs, _, err := c.Labels().IpBlockList(viper.GetString(core.GetGlobalFlagName(c.Resource, config.ArgIpBlockId)))
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, getLabelResources(labelDcs)))
}

func RunIpBlockLabelGet(c *core.CommandConfig) error {
	ipBlockId := viper.GetString(core.GetGlobalFlagName(c.Resource, config.ArgIpBlockId))
	labelKey := viper.GetString(core.GetFlagName(c.NS, config.ArgLabelKey))
	c.Printer.Verbose("Getting IpBlock with id: %v...", ipBlockId)
	labelDc, _, err := c.Labels().IpBlockGet(ipBlockId, labelKey)
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, getLabelResource(labelDc)))
}

func RunIpBlockLabelAdd(c *core.CommandConfig) error {
	labelDc, _, err := c.Labels().IpBlockCreate(
		viper.GetString(core.GetGlobalFlagName(c.Resource, config.ArgIpBlockId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgLabelKey)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgLabelValue)),
	)
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, getLabelResource(labelDc)))
}

func RunIpBlockLabelRemove(c *core.CommandConfig) error {
	_, err := c.Labels().IpBlockDelete(
		viper.GetString(core.GetGlobalFlagName(c.Resource, config.ArgIpBlockId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgLabelKey)),
	)
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, nil))
}

func RunSnapshotLabelsList(c *core.CommandConfig) error {
	labelDcs, _, err := c.Labels().SnapshotList(viper.GetString(core.GetGlobalFlagName(c.Resource, config.ArgSnapshotId)))
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, getLabelResources(labelDcs)))
}

func RunSnapshotLabelGet(c *core.CommandConfig) error {
	snapshotId := viper.GetString(core.GetGlobalFlagName(c.Resource, config.ArgSnapshotId))
	labelKey := viper.GetString(core.GetFlagName(c.NS, config.ArgLabelKey))
	c.Printer.Verbose("Getting Snapshot with id: %v...", snapshotId)
	labelDc, _, err := c.Labels().SnapshotGet(snapshotId, labelKey)
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, getLabelResource(labelDc)))
}

func RunSnapshotLabelAdd(c *core.CommandConfig) error {
	labelDc, _, err := c.Labels().SnapshotCreate(
		viper.GetString(core.GetGlobalFlagName(c.Resource, config.ArgSnapshotId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgLabelKey)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgLabelValue)),
	)
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, getLabelResource(labelDc)))
}

func RunSnapshotLabelRemove(c *core.CommandConfig) error {
	_, err := c.Labels().SnapshotDelete(
		viper.GetString(core.GetGlobalFlagName(c.Resource, config.ArgSnapshotId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgLabelKey)),
	)
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
