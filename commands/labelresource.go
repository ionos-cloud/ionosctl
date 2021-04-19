package commands

import (
	"context"
	"github.com/hashicorp/go-multierror"
	"io"
	"os"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/pkg/builder"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/resources"
	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
	"github.com/ionos-cloud/ionosctl/pkg/utils/printer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	listLabelsCmd  = "list-labels"
	getLabelCmd    = "get-label"
	addLabelCmd    = "add-label"
	removeLabelCmd = "remove-label"
)

func labelDatacenter(datacenterCmd *builder.Command) {
	ctx := context.TODO()

	/*
		List Labels Command
	*/
	listLabels := builder.NewCommand(ctx, datacenterCmd, PreRunDataCenterIdValidate, RunDataCenterListLabels, listLabelsCmd, "List Labels from a Data Center",
		"Use this command to list all Labels from a specified Data Center.\n\nRequired values to run command:\n\n* Data Center Id", listDataCenterLabelsExample, true)
	listLabels.AddStringFlag(config.ArgDataCenterId, "", "", config.RequiredFlagDatacenterId)
	_ = listLabels.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Get Label Command
	*/
	get := builder.NewCommand(ctx, datacenterCmd, PreRunDcIdLabelKeyValidate, RunDataCenterGetLabel, getLabelCmd, "Get a Label from a Data Center",
		"Use this command to get information about a specified Label from a Data Center.\n\nRequired values to run command:\n\n* Data Center Id\n* Label Key",
		getDataCenterLabelExample, true)
	get.AddStringFlag(config.ArgLabelKey, "", "", config.RequiredFlagLabelKey)
	_ = get.Command.RegisterFlagCompletionFunc(config.ArgLabelKey, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDatacenterLabelIds(os.Stderr, viper.GetString(builder.GetFlagName(datacenterCmd.Command.Name(), get.Command.Name(), config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddStringFlag(config.ArgDataCenterId, "", "", config.RequiredFlagDatacenterId)
	_ = get.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})

	/*
	 Add Label Command
	*/
	addLabel := builder.NewCommand(ctx, datacenterCmd, PreRunDcIdLabelKeyValueValidate, RunDataCenterAddLabel, addLabelCmd, "Add a Label to a Data Center",
		`Use this command to add a Label to a Data Center. You must specify the key and the value for the Label.

Required values to run command: 

* Data Center Id 
* Label Key
* Label Value`, addDataCenterLabelExample, true)
	addLabel.AddStringFlag(config.ArgLabelKey, "", "", config.RequiredFlagLabelKey)
	addLabel.AddStringFlag(config.ArgLabelValue, "", "", config.RequiredFlagLabelValue)
	addLabel.AddStringFlag(config.ArgDataCenterId, "", "", config.RequiredFlagDatacenterId)
	_ = addLabel.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Remove Label Command
	*/
	removeLabel := builder.NewCommand(ctx, datacenterCmd, PreRunDcIdLabelKeyValidate, RunDataCenterRemoveLabel, removeLabelCmd, "Remove a Label from a Data Center",
		`Use this command to remove/delete a specified Label from a Data Center.

Required values to run command:

* Data Center Id
* Label Key`, removeDataCenterLabelExample, true)
	removeLabel.AddStringFlag(config.ArgLabelKey, "", "", config.RequiredFlagLabelKey)
	_ = removeLabel.Command.RegisterFlagCompletionFunc(config.ArgLabelKey, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDatacenterLabelIds(os.Stderr, viper.GetString(builder.GetFlagName(datacenterCmd.Command.Name(), removeLabel.Command.Name(), config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	removeLabel.AddStringFlag(config.ArgDataCenterId, "", "", config.RequiredFlagDatacenterId)
	_ = removeLabel.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})

	return
}

func PreRunDcIdLabelKeyValidate(c *builder.PreCommandConfig) error {
	return builder.CheckRequiredFlags(c.ParentName, c.Name, config.ArgDataCenterId, config.ArgLabelKey)
}

func PreRunDcIdLabelKeyValueValidate(c *builder.PreCommandConfig) error {
	return builder.CheckRequiredFlags(c.ParentName, c.Name, config.ArgDataCenterId, config.ArgLabelKey, config.ArgLabelValue)
}

func RunDataCenterListLabels(c *builder.CommandConfig) error {
	labelDcs, _, err := c.Labels().DatacenterList(viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgDataCenterId)))
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, getLabelResources(labelDcs)))
}

func RunDataCenterGetLabel(c *builder.CommandConfig) error {
	labelDc, _, err := c.Labels().DatacenterGet(
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgDataCenterId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgLabelKey)))
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, getLabelResource(labelDc)))
}

func RunDataCenterAddLabel(c *builder.CommandConfig) error {
	labelDc, _, err := c.Labels().DatacenterCreate(
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgDataCenterId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgLabelKey)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgLabelValue)),
	)
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, getLabelResource(labelDc)))
}

func RunDataCenterRemoveLabel(c *builder.CommandConfig) error {
	_, err := c.Labels().DatacenterDelete(
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgDataCenterId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgLabelKey)),
	)
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, nil))
}

func getDatacenterLabelIds(outErr io.Writer, datacenterId string) []string {
	err := config.Load()
	clierror.CheckError(err, outErr)
	clientSvc, err := resources.NewClientService(
		viper.GetString(config.Username),
		viper.GetString(config.Password),
		viper.GetString(config.Token),
		viper.GetString(config.ArgServerUrl),
	)
	clierror.CheckError(err, outErr)
	labelSvc := resources.NewLabelResourceService(clientSvc.Get(), context.TODO())
	labels, _, err := labelSvc.DatacenterList(datacenterId)
	clierror.CheckError(err, outErr)
	labelsIds := make([]string, 0)
	if items, ok := labels.GetItemsOk(); ok && items != nil {
		for _, item := range *items {
			if itemId, ok := item.GetIdOk(); ok && itemId != nil {
				labelsIds = append(labelsIds, *itemId)
			}
		}
	} else {
		return nil
	}
	return labelsIds
}

func labelServer(serverCmd *builder.Command) {
	// Note: Data Center Id flag is Global Flag for serverCmd
	ctx := context.TODO()

	/*
		List Labels Command
	*/
	list := builder.NewCommand(ctx, serverCmd, PreRunGlobalDcIdServerIdValidate, RunServerListLabels, listLabelsCmd, "List Labels from a Server",
		"Use this command to list all Labels from a specified Server.\n\nRequired values to run command:\n\n* Data Center Id\n* Server Id", listServerLabelsExample, true)
	list.AddStringFlag(config.ArgServerId, "", "", config.RequiredFlagServerId)
	_ = list.Command.RegisterFlagCompletionFunc(config.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getServersIds(os.Stderr, viper.GetString(builder.GetGlobalFlagName(serverCmd.Command.Name(), config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Get Label Command
	*/
	get := builder.NewCommand(ctx, serverCmd, PreRunGlobalDcIdServerLabelKeyValidate, RunServerGetLabel, getLabelCmd, "Get a Label from a Server",
		"Use this command to get information about a specified Label from a Server.\n\nRequired values to run command:\n\n* Data Center Id\n* Server Id\n* Label Key",
		getServerLabelExample, true)
	get.AddStringFlag(config.ArgLabelKey, "", "", config.RequiredFlagLabelKey)
	_ = get.Command.RegisterFlagCompletionFunc(config.ArgLabelKey, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getServerLabelIds(os.Stderr,
			viper.GetString(builder.GetGlobalFlagName(serverCmd.Command.Name(), config.ArgDataCenterId)),
			viper.GetString(builder.GetFlagName(serverCmd.Command.Name(), get.Command.Name(), config.ArgServerId)),
		), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddStringFlag(config.ArgServerId, "", "", config.RequiredFlagServerId)
	_ = get.Command.RegisterFlagCompletionFunc(config.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getServersIds(os.Stderr, viper.GetString(builder.GetGlobalFlagName(serverCmd.Command.Name(), config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Add Label Command
	*/
	addLabel := builder.NewCommand(ctx, serverCmd, PreRunGlobalDcIdServerLabelKeyValueValidate, RunServerAddLabel, addLabelCmd, "Add a Label on a Server",
		`Use this command to add/create a Label on Server. You must specify the key and the value for the Label.

Required values to run command:

* Data Center Id
* Server Id
* Label Key
* Label Value`, addServerLabelExample, true)
	addLabel.AddStringFlag(config.ArgLabelKey, "", "", config.RequiredFlagLabelKey)
	addLabel.AddStringFlag(config.ArgLabelValue, "", "", config.RequiredFlagLabelValue)
	addLabel.AddStringFlag(config.ArgServerId, "", "", config.RequiredFlagServerId)
	_ = addLabel.Command.RegisterFlagCompletionFunc(config.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getServersIds(os.Stderr, viper.GetString(builder.GetGlobalFlagName(serverCmd.Command.Name(), config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Remove Label Command
	*/
	removeLabel := builder.NewCommand(ctx, serverCmd, PreRunGlobalDcIdServerLabelKeyValidate, RunServerRemoveLabel, removeLabelCmd, "Remove a Label from a Server",
		`Use this command to remove/delete a specified Label from a Server.

Required values to run command:

* Data Center Id
* Server Id
* Label Key`, removeServerLabelExample, true)
	removeLabel.AddStringFlag(config.ArgLabelKey, "", "", config.RequiredFlagLabelKey)
	_ = removeLabel.Command.RegisterFlagCompletionFunc(config.ArgLabelKey, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getServerLabelIds(os.Stderr,
			viper.GetString(builder.GetGlobalFlagName(serverCmd.Command.Name(), config.ArgDataCenterId)),
			viper.GetString(builder.GetFlagName(serverCmd.Command.Name(), removeLabel.Command.Name(), config.ArgServerId)),
		), cobra.ShellCompDirectiveNoFileComp
	})
	removeLabel.AddStringFlag(config.ArgServerId, "", "", config.RequiredFlagServerId)
	_ = removeLabel.Command.RegisterFlagCompletionFunc(config.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getServersIds(os.Stderr, viper.GetString(builder.GetGlobalFlagName(serverCmd.Command.Name(), config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	return
}

func PreRunGlobalDcIdServerLabelKeyValidate(c *builder.PreCommandConfig) error {
	var result *multierror.Error
	if err := builder.CheckRequiredGlobalFlags(c.ParentName, config.ArgDataCenterId); err != nil {
		result = multierror.Append(result, err)
	}
	if err := builder.CheckRequiredFlags(c.ParentName, c.Name, config.ArgServerId, config.ArgLabelKey); err != nil {
		result = multierror.Append(result, err)
	}
	if result != nil {
		return result
	}
	return nil
}

func PreRunGlobalDcIdServerLabelKeyValueValidate(c *builder.PreCommandConfig) error {
	var result *multierror.Error
	if err := builder.CheckRequiredGlobalFlags(c.ParentName, config.ArgDataCenterId); err != nil {
		result = multierror.Append(result, err)
	}
	if err := builder.CheckRequiredFlags(c.ParentName, c.Name, config.ArgServerId, config.ArgLabelKey, config.ArgLabelValue); err != nil {
		result = multierror.Append(result, err)
	}
	if result != nil {
		return result
	}
	return nil
}

func RunServerListLabels(c *builder.CommandConfig) error {
	labelDcs, _, err := c.Labels().ServerList(
		viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgDataCenterId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgServerId)),
	)
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, getLabelResources(labelDcs)))
}

func RunServerGetLabel(c *builder.CommandConfig) error {
	labelDc, _, err := c.Labels().ServerGet(
		viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgDataCenterId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgServerId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgLabelKey)))
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, getLabelResource(labelDc)))
}

func RunServerAddLabel(c *builder.CommandConfig) error {
	labelDc, _, err := c.Labels().ServerCreate(
		viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgDataCenterId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgServerId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgLabelKey)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgLabelValue)),
	)
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, getLabelResource(labelDc)))
}

func RunServerRemoveLabel(c *builder.CommandConfig) error {
	_, err := c.Labels().ServerDelete(
		viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgDataCenterId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgServerId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgLabelKey)),
	)
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, nil))
}

func getServerLabelIds(outErr io.Writer, datacenterId, serverId string) []string {
	err := config.Load()
	clierror.CheckError(err, outErr)
	clientSvc, err := resources.NewClientService(
		viper.GetString(config.Username),
		viper.GetString(config.Password),
		viper.GetString(config.Token),
		viper.GetString(config.ArgServerUrl),
	)
	clierror.CheckError(err, outErr)
	labelSvc := resources.NewLabelResourceService(clientSvc.Get(), context.TODO())
	labels, _, err := labelSvc.ServerList(datacenterId, serverId)
	clierror.CheckError(err, outErr)
	labelsIds := make([]string, 0)
	if items, ok := labels.GetItemsOk(); ok && items != nil {
		for _, item := range *items {
			if itemId, ok := item.GetIdOk(); ok && itemId != nil {
				labelsIds = append(labelsIds, *itemId)
			}
		}
	} else {
		return nil
	}
	return labelsIds
}

func labelVolume(volumeCmd *builder.Command) {
	// Note: Data Center Id flag is Global Flag for volumeCmd
	ctx := context.TODO()

	/*
		List Labels Command
	*/
	list := builder.NewCommand(ctx, volumeCmd, PreRunGlobalDcIdVolumeIdValidate, RunVolumeListLabels, listLabelsCmd, "List Labels from a Volume",
		"Use this command to list all Labels from a specified Volume.\n\nRequired values to run command:\n\n* Data Center Id\n* Volume Id", listVolumeLabelsExample, true)
	list.AddStringFlag(config.ArgVolumeId, "", "", config.RequiredFlagVolumeId)
	_ = list.Command.RegisterFlagCompletionFunc(config.ArgVolumeId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getVolumesIds(os.Stderr, viper.GetString(builder.GetGlobalFlagName(volumeCmd.Command.Name(), config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Get Label Command
	*/
	get := builder.NewCommand(ctx, volumeCmd, PreRunGlobalDcIdVolumeLabelKeyValidate, RunVolumeGetLabel, getLabelCmd, "Get a Label from a Volume",
		"Use this command to get information about a specified Label from a Volume.\n\nRequired values to run command:\n\n* Data Center Id\n* Volume Id\n* Label Key",
		getVolumeLabelExample, true)
	get.AddStringFlag(config.ArgLabelKey, "", "", config.RequiredFlagLabelKey)
	_ = get.Command.RegisterFlagCompletionFunc(config.ArgLabelKey, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getVolumeLabelIds(os.Stderr,
			viper.GetString(builder.GetGlobalFlagName(volumeCmd.Command.Name(), config.ArgDataCenterId)),
			viper.GetString(builder.GetFlagName(volumeCmd.Command.Name(), get.Command.Name(), config.ArgVolumeId)),
		), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddStringFlag(config.ArgVolumeId, "", "", config.RequiredFlagVolumeId)
	_ = get.Command.RegisterFlagCompletionFunc(config.ArgVolumeId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getVolumesIds(os.Stderr, viper.GetString(builder.GetGlobalFlagName(volumeCmd.Command.Name(), config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Add Label Command
	*/
	addLabel := builder.NewCommand(ctx, volumeCmd, PreRunGlobalDcIdVolumeLabelKeyValueValidate, RunVolumeAddLabel, addLabelCmd, "Add a Label on a Volume",
		`Use this command to add/create a Label on Volume. You must specify the key and the value for the Label.

Required values to run command:

* Data Center Id
* Volume Id
* Label Key
* Label Value`, addVolumeLabelExample, true)
	addLabel.AddStringFlag(config.ArgLabelKey, "", "", config.RequiredFlagLabelKey)
	addLabel.AddStringFlag(config.ArgLabelValue, "", "", config.RequiredFlagLabelValue)
	addLabel.AddStringFlag(config.ArgVolumeId, "", "", config.RequiredFlagVolumeId)
	_ = addLabel.Command.RegisterFlagCompletionFunc(config.ArgVolumeId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getVolumesIds(os.Stderr, viper.GetString(builder.GetGlobalFlagName(volumeCmd.Command.Name(), config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	/*
		Remove Label Command
	*/
	removeLabel := builder.NewCommand(ctx, volumeCmd, PreRunGlobalDcIdVolumeLabelKeyValidate, RunVolumeRemoveLabel, removeLabelCmd, "Remove a Label from a Volume",
		`Use this command to remove/delete a specified Label from a Volume.

Required values to run command:

* Data Center Id
* Volume Id
* Label Key`, removeVolumeLabelExample, true)
	removeLabel.AddStringFlag(config.ArgLabelKey, "", "", config.RequiredFlagLabelKey)
	_ = removeLabel.Command.RegisterFlagCompletionFunc(config.ArgLabelKey, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getVolumeLabelIds(os.Stderr,
			viper.GetString(builder.GetGlobalFlagName(volumeCmd.Command.Name(), config.ArgDataCenterId)),
			viper.GetString(builder.GetFlagName(volumeCmd.Command.Name(), removeLabel.Command.Name(), config.ArgVolumeId)),
		), cobra.ShellCompDirectiveNoFileComp
	})
	removeLabel.AddStringFlag(config.ArgVolumeId, "", "", config.RequiredFlagVolumeId)
	_ = removeLabel.Command.RegisterFlagCompletionFunc(config.ArgVolumeId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getVolumesIds(os.Stderr, viper.GetString(builder.GetGlobalFlagName(volumeCmd.Command.Name(), config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	return
}

func PreRunGlobalDcIdVolumeLabelKeyValidate(c *builder.PreCommandConfig) error {
	var result *multierror.Error
	if err := builder.CheckRequiredGlobalFlags(c.ParentName, config.ArgDataCenterId); err != nil {
		result = multierror.Append(result, err)
	}
	if err := builder.CheckRequiredFlags(c.ParentName, c.Name, config.ArgVolumeId, config.ArgLabelKey); err != nil {
		result = multierror.Append(result, err)
	}
	if result != nil {
		return result
	}
	return nil
}

func PreRunGlobalDcIdVolumeLabelKeyValueValidate(c *builder.PreCommandConfig) error {
	var result *multierror.Error
	if err := builder.CheckRequiredGlobalFlags(c.ParentName, config.ArgDataCenterId); err != nil {
		result = multierror.Append(result, err)
	}
	if err := builder.CheckRequiredFlags(c.ParentName, c.Name, config.ArgVolumeId, config.ArgLabelKey, config.ArgLabelValue); err != nil {
		result = multierror.Append(result, err)
	}
	if result != nil {
		return result
	}
	return nil
}

func RunVolumeListLabels(c *builder.CommandConfig) error {
	labelDcs, _, err := c.Labels().VolumeList(
		viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgDataCenterId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgVolumeId)),
	)
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, getLabelResources(labelDcs)))
}

func RunVolumeGetLabel(c *builder.CommandConfig) error {
	labelDc, _, err := c.Labels().VolumeGet(
		viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgDataCenterId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgVolumeId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgLabelKey)))
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, getLabelResource(labelDc)))
}

func RunVolumeAddLabel(c *builder.CommandConfig) error {
	labelDc, _, err := c.Labels().VolumeCreate(
		viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgDataCenterId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgVolumeId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgLabelKey)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgLabelValue)),
	)
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, getLabelResource(labelDc)))
}

func RunVolumeRemoveLabel(c *builder.CommandConfig) error {
	_, err := c.Labels().VolumeDelete(
		viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgDataCenterId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgVolumeId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgLabelKey)),
	)
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, nil))
}

func getVolumeLabelIds(outErr io.Writer, datacenterId, volumeId string) []string {
	err := config.Load()
	clierror.CheckError(err, outErr)
	clientSvc, err := resources.NewClientService(
		viper.GetString(config.Username),
		viper.GetString(config.Password),
		viper.GetString(config.Token),
		viper.GetString(config.ArgServerUrl),
	)
	clierror.CheckError(err, outErr)
	labelSvc := resources.NewLabelResourceService(clientSvc.Get(), context.TODO())
	labels, _, err := labelSvc.VolumeList(datacenterId, volumeId)
	clierror.CheckError(err, outErr)
	labelsIds := make([]string, 0)
	if items, ok := labels.GetItemsOk(); ok && items != nil {
		for _, item := range *items {
			if itemId, ok := item.GetIdOk(); ok && itemId != nil {
				labelsIds = append(labelsIds, *itemId)
			}
		}
	} else {
		return nil
	}
	return labelsIds
}

func labelIpBlock(ipBlockCmd *builder.Command) {
	ctx := context.TODO()

	/*
		List Labels Command
	*/
	list := builder.NewCommand(ctx, ipBlockCmd, PreRunIpBlockIdValidate, RunIpBlockListLabels, listLabelsCmd, "List Labels from a IPBlock",
		"Use this command to list all Labels from a specified IPBlock.\n\nRequired values to run command:\n\n* IPBlock Id", listIpBlockLabelsExample, true)
	list.AddStringFlag(config.ArgIpBlockId, "", "", config.RequiredFlagIpBlockId)
	_ = list.Command.RegisterFlagCompletionFunc(config.ArgIpBlockId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getIpBlocksIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Get Label Command
	*/
	get := builder.NewCommand(ctx, ipBlockCmd, PreRunIpBlockIdLabelKeyValidate, RunIpBlockGetLabel, getLabelCmd, "Get a Label from a IPBlock",
		"Use this command to get information about a specified Label from a IPBlock.\n\nRequired values to run command:\n\n* IPBlock Id\n* Label Key",
		getIpBlockLabelExample, true)
	get.AddStringFlag(config.ArgLabelKey, "", "", config.RequiredFlagLabelKey)
	_ = get.Command.RegisterFlagCompletionFunc(config.ArgLabelKey, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getIPBlockLabelIds(os.Stderr, viper.GetString(builder.GetFlagName(ipBlockCmd.Command.Name(), get.Command.Name(), config.ArgIpBlockId))), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddStringFlag(config.ArgIpBlockId, "", "", config.RequiredFlagIpBlockId)
	_ = get.Command.RegisterFlagCompletionFunc(config.ArgIpBlockId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getIpBlocksIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Add Label Command
	*/
	addLabel := builder.NewCommand(ctx, ipBlockCmd, PreRunIpBlockIdLabelKeyValueValidate, RunIpBlockAddLabel, addLabelCmd, "Add a Label on a IPBlock",
		`Use this command to add/create a Label on IPBlock. You must specify the key and the value for the Label.

Required values to run command: 

* IPBlock Id 
* Label Key
* Label Value`, addIpBlockLabelExample, true)
	addLabel.AddStringFlag(config.ArgLabelKey, "", "", config.RequiredFlagLabelKey)
	addLabel.AddStringFlag(config.ArgLabelValue, "", "", config.RequiredFlagLabelValue)
	addLabel.AddStringFlag(config.ArgIpBlockId, "", "", config.RequiredFlagIpBlockId)
	_ = addLabel.Command.RegisterFlagCompletionFunc(config.ArgIpBlockId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getIpBlocksIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Remove Label Command
	*/
	removeLabel := builder.NewCommand(ctx, ipBlockCmd, PreRunIpBlockIdLabelKeyValidate, RunIpBlockRemoveLabel, removeLabelCmd, "Remove a Label from a IPBlock",
		`Use this command to remove/delete a specified Label from a IPBlock.

Required values to run command:

* IPBlock Id
* Label Key`, removeIpBlockLabelExample, true)
	removeLabel.AddStringFlag(config.ArgLabelKey, "", "", config.RequiredFlagLabelKey)
	_ = removeLabel.Command.RegisterFlagCompletionFunc(config.ArgLabelKey, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getIPBlockLabelIds(os.Stderr, viper.GetString(builder.GetFlagName(ipBlockCmd.Command.Name(), removeLabel.Command.Name(), config.ArgIpBlockId))), cobra.ShellCompDirectiveNoFileComp
	})
	removeLabel.AddStringFlag(config.ArgIpBlockId, "", "", config.RequiredFlagIpBlockId)
	_ = removeLabel.Command.RegisterFlagCompletionFunc(config.ArgIpBlockId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getIpBlocksIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})

	return
}

func PreRunIpBlockIdLabelKeyValidate(c *builder.PreCommandConfig) error {
	return builder.CheckRequiredFlags(c.ParentName, c.Name, config.ArgIpBlockId, config.ArgLabelKey)
}

func PreRunIpBlockIdLabelKeyValueValidate(c *builder.PreCommandConfig) error {
	return builder.CheckRequiredFlags(c.ParentName, c.Name, config.ArgIpBlockId, config.ArgLabelKey, config.ArgLabelValue)
}

func RunIpBlockListLabels(c *builder.CommandConfig) error {
	labelDcs, _, err := c.Labels().IpBlockList(viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgIpBlockId)))
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, getLabelResources(labelDcs)))
}

func RunIpBlockGetLabel(c *builder.CommandConfig) error {
	labelDc, _, err := c.Labels().IpBlockGet(
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgIpBlockId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgLabelKey)))
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, getLabelResource(labelDc)))
}

func RunIpBlockAddLabel(c *builder.CommandConfig) error {
	labelDc, _, err := c.Labels().IpBlockCreate(
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgIpBlockId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgLabelKey)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgLabelValue)),
	)
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, getLabelResource(labelDc)))
}

func RunIpBlockRemoveLabel(c *builder.CommandConfig) error {
	_, err := c.Labels().IpBlockDelete(
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgIpBlockId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgLabelKey)),
	)
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, nil))
}

func getIPBlockLabelIds(outErr io.Writer, ipblockId string) []string {
	err := config.Load()
	clierror.CheckError(err, outErr)
	clientSvc, err := resources.NewClientService(
		viper.GetString(config.Username),
		viper.GetString(config.Password),
		viper.GetString(config.Token),
		viper.GetString(config.ArgServerUrl),
	)
	clierror.CheckError(err, outErr)
	labelSvc := resources.NewLabelResourceService(clientSvc.Get(), context.TODO())
	labels, _, err := labelSvc.IpBlockList(ipblockId)
	clierror.CheckError(err, outErr)
	labelsIds := make([]string, 0)
	if items, ok := labels.GetItemsOk(); ok && items != nil {
		for _, item := range *items {
			if itemId, ok := item.GetIdOk(); ok && itemId != nil {
				labelsIds = append(labelsIds, *itemId)
			}
		}
	} else {
		return nil
	}
	return labelsIds
}

func labelSnapshot(snapshotCmd *builder.Command) {
	ctx := context.TODO()

	/*
		List Labels Command
	*/
	list := builder.NewCommand(ctx, snapshotCmd, PreRunSnapshotIdValidate, RunSnapshotListLabels, listLabelsCmd, "List Labels from a Snapshot",
		"Use this command to list all Labels from a specified Snapshot.\n\nRequired values to run command:\n\n* Snapshot Id", listSnapshotLabelsExample, true)
	list.AddStringFlag(config.ArgSnapshotId, "", "", config.RequiredFlagSnapshotId)
	_ = list.Command.RegisterFlagCompletionFunc(config.ArgSnapshotId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getSnapshotIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Get Label Command
	*/
	get := builder.NewCommand(ctx, snapshotCmd, PreRunSnapshotIdLabelKeyValidate, RunSnapshotGetLabel, getLabelCmd, "Get a Label from a Snapshot",
		"Use this command to get information about a specified Label from a Snapshot.\n\nRequired values to run command:\n\n* Snapshot Id\n* Label Key",
		getSnapshotLabelExample, true)
	get.AddStringFlag(config.ArgLabelKey, "", "", config.RequiredFlagLabelKey)
	_ = get.Command.RegisterFlagCompletionFunc(config.ArgLabelKey, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getSnapshotLabelIds(os.Stderr, viper.GetString(builder.GetFlagName(snapshotCmd.Command.Name(), get.Command.Name(), config.ArgSnapshotId))), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddStringFlag(config.ArgSnapshotId, "", "", config.RequiredFlagSnapshotId)
	_ = get.Command.RegisterFlagCompletionFunc(config.ArgSnapshotId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getSnapshotIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	/*
		Add Label Command
	*/
	addLabel := builder.NewCommand(ctx, snapshotCmd, PreRunSnapshotIdLabelKeyValueValidate, RunSnapshotAddLabel, addLabelCmd, "Add a Label on a Snapshot",
		`Use this command to create a Label on Snapshot. You must specify the key and the value for the Label.

Required values to run command: 

* Snapshot Id 
* Label Key
* Label Value`, addSnapshotLabelExample, true)
	addLabel.AddStringFlag(config.ArgLabelKey, "", "", config.RequiredFlagLabelKey)
	addLabel.AddStringFlag(config.ArgLabelValue, "", "", config.RequiredFlagLabelValue)
	addLabel.AddStringFlag(config.ArgSnapshotId, "", "", config.RequiredFlagSnapshotId)
	_ = addLabel.Command.RegisterFlagCompletionFunc(config.ArgSnapshotId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getSnapshotIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	/*
		Remove Label Command
	*/
	removeLabel := builder.NewCommand(ctx, snapshotCmd, PreRunSnapshotIdLabelKeyValidate, RunSnapshotRemoveLabel, removeLabelCmd, "Remove a Label from a Snapshot",
		`Use this command to remove/delete a specified Label from a Snapshot.

Required values to run command:

* Snapshot Id
* Label Key`, removeSnapshotLabelExample, true)
	removeLabel.AddStringFlag(config.ArgLabelKey, "", "", config.RequiredFlagLabelKey)
	_ = removeLabel.Command.RegisterFlagCompletionFunc(config.ArgLabelKey, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getSnapshotLabelIds(os.Stderr, viper.GetString(builder.GetFlagName(snapshotCmd.Command.Name(), removeLabel.Command.Name(), config.ArgSnapshotId))), cobra.ShellCompDirectiveNoFileComp
	})
	removeLabel.AddStringFlag(config.ArgSnapshotId, "", "", config.RequiredFlagSnapshotId)
	_ = removeLabel.Command.RegisterFlagCompletionFunc(config.ArgSnapshotId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getSnapshotIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	return
}

func PreRunSnapshotIdLabelKeyValidate(c *builder.PreCommandConfig) error {
	return builder.CheckRequiredFlags(c.ParentName, c.Name, config.ArgSnapshotId, config.ArgLabelKey)
}

func PreRunSnapshotIdLabelKeyValueValidate(c *builder.PreCommandConfig) error {
	return builder.CheckRequiredFlags(c.ParentName, c.Name, config.ArgSnapshotId, config.ArgLabelKey, config.ArgLabelValue)
}

func RunSnapshotListLabels(c *builder.CommandConfig) error {
	labelDcs, _, err := c.Labels().SnapshotList(viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgSnapshotId)))
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, getLabelResources(labelDcs)))
}

func RunSnapshotGetLabel(c *builder.CommandConfig) error {
	labelDc, _, err := c.Labels().SnapshotGet(
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgSnapshotId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgLabelKey)))
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, getLabelResource(labelDc)))
}

func RunSnapshotAddLabel(c *builder.CommandConfig) error {
	labelDc, _, err := c.Labels().SnapshotCreate(
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgSnapshotId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgLabelKey)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgLabelValue)),
	)
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, getLabelResource(labelDc)))
}

func RunSnapshotRemoveLabel(c *builder.CommandConfig) error {
	_, err := c.Labels().SnapshotDelete(
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgSnapshotId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgLabelKey)),
	)
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(c, nil))
}

func getSnapshotLabelIds(outErr io.Writer, snapshotId string) []string {
	err := config.Load()
	clierror.CheckError(err, outErr)
	clientSvc, err := resources.NewClientService(
		viper.GetString(config.Username),
		viper.GetString(config.Password),
		viper.GetString(config.Token),
		viper.GetString(config.ArgServerUrl),
	)
	clierror.CheckError(err, outErr)
	labelSvc := resources.NewLabelResourceService(clientSvc.Get(), context.TODO())
	labels, _, err := labelSvc.SnapshotList(snapshotId)
	clierror.CheckError(err, outErr)
	labelsIds := make([]string, 0)
	if items, ok := labels.GetItemsOk(); ok && items != nil {
		for _, item := range *items {
			if itemId, ok := item.GetIdOk(); ok && itemId != nil {
				labelsIds = append(labelsIds, *itemId)
			}
		}
	} else {
		return nil
	}
	return labelsIds
}

// Output Printing

var defaultLabelResourceCols = []string{"Key", "Value"}

type LabelResourcePrint struct {
	Key   string `json:"Key,omitempty"`
	Value string `json:"Value,omitempty"`
}

func getLabelResourcePrint(c *builder.CommandConfig, s []resources.LabelResource) printer.Result {
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
