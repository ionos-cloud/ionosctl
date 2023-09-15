package commands

import (
	"errors"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/confirm"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/tabheaders"
	"github.com/ionos-cloud/ionosctl/v6/pkg/jsontabwriter"

	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/query"
	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/waiter"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/printer"
	"github.com/ionos-cloud/ionosctl/v6/pkg/utils"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
	"github.com/spf13/viper"
)

var (
	allLabelResourceJSONPaths = map[string]string{
		"Key":   "properties.key",
		"Value": "properties.value",
	}

	defaultLabelResourceCols = []string{"Key", "Value"}
)

func RunDataCenterLabelsList(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	labelDcs, resp, err := c.CloudApiV6Services.Labels().DatacenterList(listQueryParams, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)))
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	cols, err := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
	if err != nil {
		return err
	}

	out, err := jsontabwriter.GenerateOutput("items", allLabelResourceJSONPaths, labelDcs.LabelResources,
		tabheaders.GetHeadersAllDefault(defaultLabelResourceCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)

	return nil
}

func RunDataCenterLabelGet(c *core.CommandConfig) error {
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))
	labelKey := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLabelKey))

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Getting label with key: %v for Datacenter with id: %v...", labelKey, dcId))

	labelDc, resp, err := c.CloudApiV6Services.Labels().DatacenterGet(dcId, labelKey)
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	cols, err := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
	if err != nil {
		return err
	}

	out, err := jsontabwriter.GenerateOutput("", allLabelResourceJSONPaths, labelDc.LabelResource,
		tabheaders.GetHeadersAllDefault(defaultLabelResourceCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)

	return nil
}

func RunDataCenterLabelAdd(c *core.CommandConfig) error {
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))
	labelKey := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLabelKey))
	labelValue := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLabelValue))

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
		"Adding label with key: %v and value: %v to Datacenter with id: %v...", labelKey, labelValue, dcId))

	labelDc, resp, err := c.CloudApiV6Services.Labels().DatacenterCreate(dcId, labelKey, labelValue)
	if resp != nil && printer.GetId(resp) != "" {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, printer.GetId(resp), resp.RequestTime))
	}
	if err != nil {
		return err
	}

	cols, err := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
	if err != nil {
		return err
	}

	out, err := jsontabwriter.GenerateOutput("", allLabelResourceJSONPaths, labelDc.LabelResource,
		tabheaders.GetHeadersAllDefault(defaultLabelResourceCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)

	return nil
}

func RunDataCenterLabelRemove(c *core.CommandConfig) error {
	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAll)) {
		if err := RemoveAllDatacenterLabels(c); err != nil {
			return err
		}

		fmt.Fprintf(c.Command.Command.OutOrStdout(), jsontabwriter.GenerateLogOutput("Datacenter Labels successfully deleted"))
		return nil
	}

	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))
	labelKey := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLabelKey))

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Removing label with key: %v for Datacenter with id: %v...", labelKey, dcId))

	resp, err := c.CloudApiV6Services.Labels().DatacenterDelete(dcId, labelKey)
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, printer.GetId(resp), resp.RequestTime))
	}
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), jsontabwriter.GenerateLogOutput("Datacenter Label successfully deleted"))
	return nil
}

func RemoveAllDatacenterLabels(c *core.CommandConfig) error {
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Datacenter ID: %v", dcId))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Getting Labels from Datacenter..."))

	labels, resp, err := c.CloudApiV6Services.Labels().DatacenterList(resources.ListQueryParams{}, dcId)
	if err != nil {
		return err
	}

	labelsItems, ok := labels.GetItemsOk()
	if !ok || labelsItems == nil {
		return fmt.Errorf("could not get items of Datacenter Labels")
	}

	if len(*labelsItems) <= 0 {
		return fmt.Errorf("no Datacenter Labels found")
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), jsontabwriter.GenerateLogOutput("Labels to be removed from Datacenter with ID: ", dcId))

	for _, label := range *labelsItems {
		delIdAndName := ""

		if properties, ok := label.GetPropertiesOk(); ok && properties != nil {
			if key, ok := properties.GetKeyOk(); ok && key != nil {
				delIdAndName += "Label Key: " + *key
			}

			if value, ok := properties.GetValueOk(); ok && value != nil {
				delIdAndName += " Label Value: " + *value
			}
		}

		fmt.Fprintf(c.Command.Command.OutOrStdout(), jsontabwriter.GenerateLogOutput(delIdAndName))
	}

	if !confirm.Ask("delete all the Labels from Datacenter with Id: " + dcId) {
		return nil
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Deleting all the Labels from Datacenter with Id: %v...", dcId))

	var multiErr error
	for _, label := range *labelsItems {
		properties, ok := label.GetPropertiesOk()
		if !ok || properties == nil {
			continue
		}

		key, ok := properties.GetKeyOk()
		if !ok || key == nil {
			continue
		}

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Starting deleting Label with id: %v...", *key))

		resp, err = c.CloudApiV6Services.Labels().DatacenterDelete(dcId, *key)
		if resp != nil && printer.GetId(resp) != "" {
			fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, printer.GetId(resp), resp.RequestTime))
		}
		if err != nil {
			multiErr = errors.Join(multiErr, fmt.Errorf(constants.ErrDeleteAll, c.Resource, *key, err))
			continue
		}

		fmt.Fprintf(c.Command.Command.OutOrStdout(), jsontabwriter.GenerateLogOutput(constants.MessageDeletingAll, c.Resource, *key))

		if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
			multiErr = errors.Join(multiErr, fmt.Errorf(constants.ErrWaitDeleteAll, c.Resource, *key, err))
		}
	}

	if multiErr != nil {
		return multiErr
	}

	return nil
}

func RunServerLabelsList(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	labelDcs, resp, err := c.CloudApiV6Services.Labels().ServerList(
		listQueryParams,
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgServerId)),
	)
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	cols, err := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
	if err != nil {
		return err
	}

	out, err := jsontabwriter.GenerateOutput("items", allLabelResourceJSONPaths, labelDcs.LabelResources,
		tabheaders.GetHeadersAllDefault(defaultLabelResourceCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)

	return nil
}

func RunServerLabelGet(c *core.CommandConfig) error {
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))
	serverId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgServerId))
	labelkey := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLabelKey))

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Getting label with key: %v for Server with id: %v...", labelkey, serverId))

	labelDc, resp, err := c.CloudApiV6Services.Labels().ServerGet(dcId, serverId, labelkey)
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	cols, err := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
	if err != nil {
		return err
	}

	out, err := jsontabwriter.GenerateOutput("", allLabelResourceJSONPaths, labelDc.LabelResource,
		tabheaders.GetHeadersAllDefault(defaultLabelResourceCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)

	return nil
}

func RunServerLabelAdd(c *core.CommandConfig) error {
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))
	serverId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgServerId))
	labelKey := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLabelKey))
	labelValue := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLabelValue))

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
		"Adding label with key: %v and value: %v to Server with id: %v...", labelKey, labelValue, serverId))

	labelDc, resp, err := c.CloudApiV6Services.Labels().ServerCreate(dcId, serverId, labelKey, labelValue)
	if resp != nil && printer.GetId(resp) != "" {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, printer.GetId(resp), resp.RequestTime))
	}
	if err != nil {
		return err
	}

	cols, err := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
	if err != nil {
		return err
	}

	out, err := jsontabwriter.GenerateOutput("", allLabelResourceJSONPaths, labelDc.LabelResource,
		tabheaders.GetHeadersAllDefault(defaultLabelResourceCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)

	return nil
}

func RunServerLabelRemove(c *core.CommandConfig) error {
	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAll)) {
		if err := RemoveAllServerLabels(c); err != nil {
			return err
		}

		fmt.Fprintf(c.Command.Command.OutOrStdout(), jsontabwriter.GenerateLogOutput("Server Labels successfully deleted"))
		return nil
	}

	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))
	serverId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgServerId))
	labelKey := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLabelKey))

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
		"Removing label with key: %v for Server with id: %v...", labelKey, serverId))

	resp, err := c.CloudApiV6Services.Labels().ServerDelete(dcId, serverId, labelKey)
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, printer.GetId(resp), resp.RequestTime))
	}
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), jsontabwriter.GenerateLogOutput("Server Label successfully deleted"))
	return nil
}

func RemoveAllServerLabels(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))
	serverId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgServerId))

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Datacenter ID: %v", dcId))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Server ID: %v", serverId))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Getting Labels from Server..."))

	labels, resp, err := c.CloudApiV6Services.Labels().ServerList(listQueryParams, dcId, serverId)
	if err != nil {
		return err
	}

	labelsItems, ok := labels.GetItemsOk()
	if !ok || labelsItems == nil {
		return fmt.Errorf("could not get items of Server Labels")
	}

	if len(*labelsItems) <= 0 {
		return fmt.Errorf("no Server Labels found")
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), jsontabwriter.GenerateLogOutput("Labels to be removed from Server with Id: ", serverId))

	for _, label := range *labelsItems {
		delIdAndName := ""

		if properties, ok := label.GetPropertiesOk(); ok && properties != nil {
			if key, ok := properties.GetKeyOk(); ok && key != nil {
				delIdAndName += "Label Key: " + *key
			}

			if value, ok := properties.GetValueOk(); ok && value != nil {
				delIdAndName += " Label Value: " + *value
			}
		}

		fmt.Fprintf(c.Command.Command.OutOrStdout(), jsontabwriter.GenerateLogOutput(delIdAndName))
	}

	if !confirm.Ask("delete all the Labels from Server with Id: " + serverId) {
		return nil
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Deleting all the Labels from Server with Id: %v...", serverId))

	var multiErr error
	for _, label := range *labelsItems {
		properties, ok := label.GetPropertiesOk()
		if !ok || properties == nil {
			continue
		}

		key, ok := properties.GetKeyOk()
		if !ok || key == nil {
			continue
		}

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Starting deleting Label with id: %v...", *key))

		resp, err = c.CloudApiV6Services.Labels().ServerDelete(dcId, serverId, *key)
		if resp != nil && printer.GetId(resp) != "" {
			fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, printer.GetId(resp), resp.RequestTime))
		}
		if err != nil {
			multiErr = errors.Join(multiErr, fmt.Errorf(constants.ErrDeleteAll, c.Resource, *key, err))
			continue
		}

		fmt.Fprintf(c.Command.Command.OutOrStdout(), jsontabwriter.GenerateLogOutput(constants.MessageDeletingAll, c.Resource, *key))

		if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
			multiErr = errors.Join(multiErr, fmt.Errorf(constants.ErrWaitDeleteAll, c.Resource, *key, err))
		}
	}

	if multiErr != nil {
		return multiErr
	}

	return nil
}

func RunVolumeLabelsList(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	labelDcs, resp, err := c.CloudApiV6Services.Labels().VolumeList(
		listQueryParams,
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgVolumeId)),
	)
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	cols, err := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
	if err != nil {
		return err
	}

	out, err := jsontabwriter.GenerateOutput("items", allLabelResourceJSONPaths, labelDcs.LabelResources,
		tabheaders.GetHeadersAllDefault(defaultLabelResourceCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)

	return nil
}

func RunVolumeLabelGet(c *core.CommandConfig) error {
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))
	volumeId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgVolumeId))
	labelKey := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLabelKey))

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Getting label with key: %v for Volume with id: %v...", labelKey, volumeId))

	labelDc, resp, err := c.CloudApiV6Services.Labels().VolumeGet(dcId, volumeId, labelKey)
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	cols, err := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
	if err != nil {
		return err
	}

	out, err := jsontabwriter.GenerateOutput("", allLabelResourceJSONPaths, labelDc.LabelResource,
		tabheaders.GetHeadersAllDefault(defaultLabelResourceCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)

	return nil
}

func RunVolumeLabelAdd(c *core.CommandConfig) error {
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))
	volumeId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgVolumeId))
	labelKey := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLabelKey))
	labelValue := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLabelValue))

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
		"Adding label with key: %v and value: %v to Volume with id: %v...", labelKey, labelValue, volumeId))

	labelDc, resp, err := c.CloudApiV6Services.Labels().VolumeCreate(dcId, volumeId, labelKey, labelValue)
	if resp != nil && printer.GetId(resp) != "" {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, printer.GetId(resp), resp.RequestTime))
	}
	if err != nil {
		return err
	}

	cols, err := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
	if err != nil {
		return err
	}

	out, err := jsontabwriter.GenerateOutput("", allLabelResourceJSONPaths, labelDc.LabelResource,
		tabheaders.GetHeadersAllDefault(defaultLabelResourceCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)

	return nil
}

func RunVolumeLabelRemove(c *core.CommandConfig) error {
	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAll)) {
		if err := RemoveAllVolumeLabels(c); err != nil {
			return err
		}

		fmt.Fprintf(c.Command.Command.OutOrStdout(), jsontabwriter.GenerateLogOutput("Volume Labels successfully deleted"))
		return nil
	}

	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))
	volumeId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgVolumeId))
	labelKey := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLabelKey))

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
		"Removing label with key: %v for Volume with id: %v...", labelKey, volumeId))

	resp, err := c.CloudApiV6Services.Labels().VolumeDelete(dcId, volumeId, labelKey)
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, printer.GetId(resp), resp.RequestTime))
	}
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), jsontabwriter.GenerateLogOutput("Volume Label successfully deleted"))
	return nil

}

func RemoveAllVolumeLabels(c *core.CommandConfig) error {
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))
	volumeId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgVolumeId))

	fmt.Fprintf(c.Command.Command.OutOrStdout(), jsontabwriter.GenerateLogOutput("Labels to be removed from Volume with Id: ", volumeId))

	labels, resp, err := c.CloudApiV6Services.Labels().VolumeList(resources.ListQueryParams{}, dcId, volumeId)
	if err != nil {
		return err
	}

	labelsItems, ok := labels.GetItemsOk()
	if !ok || labelsItems == nil {
		return fmt.Errorf("could not get items of Volume Labels")
	}

	if len(*labelsItems) <= 0 {
		return fmt.Errorf("no Volume Labels found")
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateLogOutput("Labels to be removed from Volume with Id: ", volumeId))

	for _, label := range *labelsItems {
		delIdAndName := ""

		if properties, ok := label.GetPropertiesOk(); ok && properties != nil {
			if key, ok := properties.GetKeyOk(); ok && key != nil {
				delIdAndName += "Label Key: " + *key
			}

			if value, ok := properties.GetValueOk(); ok && value != nil {
				delIdAndName += " Label Value: " + *value
			}
		}

		fmt.Fprintf(c.Command.Command.OutOrStdout(), jsontabwriter.GenerateLogOutput(delIdAndName))
	}

	if !confirm.Ask("delete all the Labels from Volume with Id: " + volumeId) {
		return nil
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Deleting all the Labels from Volume with Id: %v...", volumeId))

	var multiErr error
	for _, label := range *labelsItems {
		properties, ok := label.GetPropertiesOk()
		if !ok || properties == nil {
			continue
		}

		key, ok := properties.GetKeyOk()
		if !ok || key == nil {
			continue
		}

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Starting deleting Label with id: %v...", *key))

		resp, err = c.CloudApiV6Services.Labels().VolumeDelete(dcId, volumeId, *key)
		if resp != nil && printer.GetId(resp) != "" {
			fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, printer.GetId(resp), resp.RequestTime))
		}
		if err != nil {
			multiErr = errors.Join(multiErr, fmt.Errorf(constants.ErrDeleteAll, c.Resource, *key, err))
			continue
		}

		fmt.Fprintf(c.Command.Command.OutOrStdout(), jsontabwriter.GenerateLogOutput(constants.MessageDeletingAll, c.Resource, *key))

		if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
			multiErr = errors.Join(multiErr, fmt.Errorf(constants.ErrWaitDeleteAll, c.Resource, *key, err))
		}
	}

	if multiErr != nil {
		return multiErr
	}

	return nil
}

func RunIpBlockLabelsList(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	labelDcs, resp, err := c.CloudApiV6Services.Labels().IpBlockList(listQueryParams, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgIpBlockId)))
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	cols, err := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
	if err != nil {
		return err
	}

	out, err := jsontabwriter.GenerateOutput("items", allLabelResourceJSONPaths, labelDcs.LabelResources,
		tabheaders.GetHeadersAllDefault(defaultLabelResourceCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)

	return nil
}

func RunIpBlockLabelGet(c *core.CommandConfig) error {
	ipBlockId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgIpBlockId))
	labelKey := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLabelKey))

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
		"Getting label with key: %v for IpBlock with id: %v...", labelKey, ipBlockId))

	labelDc, resp, err := c.CloudApiV6Services.Labels().IpBlockGet(ipBlockId, labelKey)
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	cols, err := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
	if err != nil {
		return err
	}

	out, err := jsontabwriter.GenerateOutput("", allLabelResourceJSONPaths, labelDc.LabelResource,
		tabheaders.GetHeadersAllDefault(defaultLabelResourceCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)

	return nil
}

func RunIpBlockLabelAdd(c *core.CommandConfig) error {
	ipBlockId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgIpBlockId))
	labelKey := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLabelKey))
	labelValue := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLabelValue))

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
		"Adding label with key: %v and value: %v to IpBlock with id: %v...", labelKey, labelValue, ipBlockId))

	labelDc, resp, err := c.CloudApiV6Services.Labels().IpBlockCreate(ipBlockId, labelKey, labelValue)
	if resp != nil && printer.GetId(resp) != "" {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, printer.GetId(resp), resp.RequestTime))
	}
	if err != nil {
		return err
	}

	cols, err := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
	if err != nil {
		return err
	}

	out, err := jsontabwriter.GenerateOutput("", allLabelResourceJSONPaths, labelDc.LabelResource,
		tabheaders.GetHeadersAllDefault(defaultLabelResourceCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)

	return nil
}

func RunIpBlockLabelRemove(c *core.CommandConfig) error {
	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAll)) {
		if err := RemoveAllIpBlockLabels(c); err != nil {
			return err
		}

		fmt.Fprintf(c.Command.Command.OutOrStdout(), jsontabwriter.GenerateLogOutput("IP Block Labels successfully deleted"))
		return nil
	}

	ipBlockId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgIpBlockId))
	labelKey := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLabelKey))

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Removing label with key: %v for IpBlock with id: %v...", labelKey, ipBlockId))

	resp, err := c.CloudApiV6Services.Labels().IpBlockDelete(ipBlockId, labelKey)
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, printer.GetId(resp), resp.RequestTime))
	}
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), jsontabwriter.GenerateLogOutput("IP Block Label successfully deleted"))
	return nil
}

func RemoveAllIpBlockLabels(c *core.CommandConfig) error {
	ipBlockId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgIpBlockId))

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("IpBlock ID: %v", ipBlockId))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Getting Labels from IpBlock..."))

	labels, resp, err := c.CloudApiV6Services.Labels().IpBlockList(resources.ListQueryParams{}, ipBlockId)
	if err != nil {
		return err
	}

	labelsItems, ok := labels.GetItemsOk()
	if !ok || labelsItems == nil {
		return fmt.Errorf("could not get items of IP Block Labels")
	}

	if len(*labelsItems) <= 0 {
		return fmt.Errorf("no IP Block Labels found")
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), jsontabwriter.GenerateLogOutput("Labels to be removed from IpBlock with Id: ", ipBlockId))

	for _, label := range *labelsItems {
		delIdAndName := ""

		if properties, ok := label.GetPropertiesOk(); ok && properties != nil {
			if key, ok := properties.GetKeyOk(); ok && key != nil {
				delIdAndName += "Label Key: " + *key
			}

			if value, ok := properties.GetValueOk(); ok && value != nil {
				delIdAndName += " Label Value: " + *value
			}
		}

		fmt.Fprintf(c.Command.Command.OutOrStdout(), jsontabwriter.GenerateLogOutput(delIdAndName))
	}

	if !confirm.Ask("delete all the Labels from IpBlock with Id: " + ipBlockId) {
		return nil
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Deleting all the Labels from IpBlock with Id: %v...", ipBlockId))

	var multiErr error
	for _, label := range *labelsItems {
		properties, ok := label.GetPropertiesOk()
		if !ok || properties == nil {
			continue
		}

		key, ok := properties.GetKeyOk()
		if !ok || key == nil {
			continue
		}

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Starting deleting Label with id: %v...", *key))

		resp, err = c.CloudApiV6Services.Labels().IpBlockDelete(ipBlockId, *key)
		if resp != nil && printer.GetId(resp) != "" {
			fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, printer.GetId(resp), resp.RequestTime))
		}
		if err != nil {
			multiErr = errors.Join(multiErr, fmt.Errorf(constants.ErrDeleteAll, c.Resource, *key, err))
			continue
		}

		fmt.Fprintf(c.Command.Command.OutOrStdout(), jsontabwriter.GenerateLogOutput(constants.MessageDeletingAll, c.Resource, *key))

		if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
			multiErr = errors.Join(multiErr, fmt.Errorf(constants.ErrWaitDeleteAll, c.Resource, *key, err))
		}
	}

	if multiErr != nil {
		return multiErr
	}

	return nil
}

func RunSnapshotLabelsList(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	labelDcs, resp, err := c.CloudApiV6Services.Labels().SnapshotList(listQueryParams, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgSnapshotId)))
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	cols, err := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
	if err != nil {
		return err
	}

	out, err := jsontabwriter.GenerateOutput("items", allLabelResourceJSONPaths, labelDcs.LabelResources,
		tabheaders.GetHeadersAllDefault(defaultLabelResourceCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)

	return nil
}

func RunSnapshotLabelGet(c *core.CommandConfig) error {
	snapshotId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgSnapshotId))
	labelKey := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLabelKey))

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
		"Getting label with key: %v for Snapshot with id: %v...", labelKey, snapshotId))

	labelDc, resp, err := c.CloudApiV6Services.Labels().SnapshotGet(snapshotId, labelKey)
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	cols, err := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
	if err != nil {
		return err
	}

	out, err := jsontabwriter.GenerateOutput("", allLabelResourceJSONPaths, labelDc.LabelResource,
		tabheaders.GetHeadersAllDefault(defaultLabelResourceCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)

	return nil
}

func RunSnapshotLabelAdd(c *core.CommandConfig) error {
	snapshotId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgSnapshotId))
	labelKey := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLabelKey))
	labelValue := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLabelValue))

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
		"Adding label with key: %v and value: %v to Snapshot with id: %v...", labelKey, labelValue, snapshotId))

	labelDc, resp, err := c.CloudApiV6Services.Labels().SnapshotCreate(snapshotId, labelKey, labelValue)
	if resp != nil && printer.GetId(resp) != "" {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, printer.GetId(resp), resp.RequestTime))
	}
	if err != nil {
		return err
	}

	cols, err := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
	if err != nil {
		return err
	}

	out, err := jsontabwriter.GenerateOutput("", allLabelResourceJSONPaths, labelDc.LabelResource,
		tabheaders.GetHeadersAllDefault(defaultLabelResourceCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)

	return nil
}

func RunSnapshotLabelRemove(c *core.CommandConfig) error {
	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAll)) {
		if err := RemoveAllSnapshotLabels(c); err != nil {
			return err
		}

		fmt.Fprintf(c.Command.Command.OutOrStdout(), jsontabwriter.GenerateLogOutput("Snapshot Labels successfully deleted"))
		return nil
	}

	snapshotId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgSnapshotId))
	labelKey := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLabelKey))

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
		"Removing label with key: %v for Snapshot with id: %v...", labelKey, snapshotId))

	resp, err := c.CloudApiV6Services.Labels().SnapshotDelete(snapshotId, labelKey)
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, printer.GetId(resp), resp.RequestTime))
	}
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), jsontabwriter.GenerateLogOutput("Snapshot Label successfully deleted"))
	return nil
}

func RemoveAllSnapshotLabels(c *core.CommandConfig) error {
	snapshotId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgSnapshotId))

	fmt.Fprintf(c.Command.Command.OutOrStdout(), jsontabwriter.GenerateLogOutput("Labels to be removed from Snapshot with Id: ", snapshotId))

	labels, resp, err := c.CloudApiV6Services.Labels().SnapshotList(resources.ListQueryParams{}, snapshotId)
	if err != nil {
		return err
	}

	labelsItems, ok := labels.GetItemsOk()
	if !ok || labelsItems == nil {
		return fmt.Errorf("could not get items of Snapshot Labels")
	}

	if len(*labelsItems) <= 0 {
		return fmt.Errorf("no Snapshot Labels found")
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), jsontabwriter.GenerateLogOutput("Labels to be removed from Snapshot with Id: ", snapshotId))

	for _, label := range *labelsItems {
		delIdAndName := ""

		if properties, ok := label.GetPropertiesOk(); ok && properties != nil {
			if key, ok := properties.GetKeyOk(); ok && key != nil {
				delIdAndName += "Label Key: " + *key
			}

			if value, ok := properties.GetValueOk(); ok && value != nil {
				delIdAndName += " Label Value: " + *value
			}
		}

		fmt.Fprintf(c.Command.Command.OutOrStdout(), jsontabwriter.GenerateLogOutput(delIdAndName))
	}

	if !confirm.Ask("delete all the Labels from Snapshot with Id: " + snapshotId) {
		return nil
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Deleting all the Labels from Snapshot with Id: %v...", snapshotId))

	var multiErr error
	for _, label := range *labelsItems {
		properties, ok := label.GetPropertiesOk()
		if !ok || properties == nil {
			continue
		}

		key, ok := properties.GetKeyOk()
		if !ok || key == nil {
			continue
		}

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Starting deleting Label with id: %v...", *key))

		resp, err = c.CloudApiV6Services.Labels().SnapshotDelete(snapshotId, *key)
		if resp != nil && printer.GetId(resp) != "" {
			fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, printer.GetId(resp), resp.RequestTime))
		}
		if err != nil {
			multiErr = errors.Join(multiErr, fmt.Errorf(constants.ErrDeleteAll, c.Resource, *key, err))
			continue
		}

		fmt.Fprintf(c.Command.Command.OutOrStdout(), jsontabwriter.GenerateLogOutput(constants.MessageDeletingAll, c.Resource, *key))

		if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
			multiErr = errors.Join(multiErr, fmt.Errorf(constants.ErrWaitDeleteAll, c.Resource, *key, err))
		}
	}

	if multiErr != nil {
		return multiErr
	}

	return nil
}
