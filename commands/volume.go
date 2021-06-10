package commands

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/resources"
	"github.com/ionos-cloud/ionosctl/pkg/utils"
	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
	"github.com/ionos-cloud/ionosctl/pkg/utils/printer"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	multierror "go.uber.org/multierr"
)

func volume() *core.Command {
	ctx := context.TODO()
	volumeCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "volume",
			Aliases:          []string{"v", "vol"},
			Short:            "Volume Operations",
			Long:             `The sub-commands of ` + "`" + `ionosctl volume` + "`" + ` manage your block storage volumes by creating, updating, getting specific information, deleting Volumes. To attach a Volume to a Server, use the Server command ` + "`" + `ionosctl server volume attach` + "`" + `.`,
			TraverseChildren: true,
		},
	}
	globalFlags := volumeCmd.GlobalFlags()
	globalFlags.StringP(config.ArgDataCenterId, "", "", config.RequiredFlagDatacenterId)
	_ = viper.BindPFlag(core.GetGlobalFlagName(volumeCmd.Name(), config.ArgDataCenterId), globalFlags.Lookup(config.ArgDataCenterId))
	_ = volumeCmd.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	globalFlags.StringSliceP(config.ArgCols, "", defaultVolumeCols, utils.ColsMessage(allVolumeCols))
	_ = viper.BindPFlag(core.GetGlobalFlagName(volumeCmd.Name(), config.ArgCols), globalFlags.Lookup(config.ArgCols))
	_ = volumeCmd.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allVolumeCols, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		List Command
	*/
	core.NewCommand(ctx, volumeCmd, core.CommandBuilder{
		Namespace:  "volume",
		Resource:   "volume",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List Volumes",
		LongDesc:   "Use this command to list all Volumes from a Data Center on your account.\n\nRequired values to run command:\n\n* Data Center Id",
		Example:    listVolumeExample,
		PreCmdRun:  PreRunGlobalDcId,
		CmdRun:     RunVolumeList,
		InitClient: true,
	})

	/*
		Get Command
	*/
	get := core.NewCommand(ctx, volumeCmd, core.CommandBuilder{
		Namespace:  "volume",
		Resource:   "volume",
		Verb:       "get",
		Aliases:    []string{"g"},
		ShortDesc:  "Get a Volume",
		LongDesc:   "Use this command to retrieve information about a Volume using its ID.\n\nRequired values to run command:\n\n* Data Center Id\n* Volume Id",
		Example:    getVolumeExample,
		PreCmdRun:  PreRunGlobalDcIdVolumeId,
		CmdRun:     RunVolumeGet,
		InitClient: true,
	})
	get.AddStringFlag(config.ArgVolumeId, config.ArgIdShort, "", config.RequiredFlagVolumeId)
	_ = get.Command.RegisterFlagCompletionFunc(config.ArgVolumeId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getVolumesIds(os.Stderr, viper.GetString(core.GetGlobalFlagName(volumeCmd.Name(), config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Create Command
	*/
	create := core.NewCommand(ctx, volumeCmd, core.CommandBuilder{
		Namespace: "volume",
		Resource:  "volume",
		Verb:      "create",
		Aliases:   []string{"c"},
		ShortDesc: "Create a Volume",
		LongDesc: `Use this command to create a Volume on your account, within a Data Center. This will NOT attach the Volume to a Server. Please see the Servers commands for details on how to attach storage Volumes. You can specify the name, size, type, licence type, availability zone, image and other properties for the object.

Note: You will need to provide a valid value for either the Image, Image Alias, or the Licence Type options. The Licence Type is required, but if Image or Image Alias is supplied, then Licence Type is already set and cannot be changed. Similarly either the Image Password or SSH Keys attributes need to be defined when creating a Volume that uses an Image or Image Alias of an IONOS public HDD Image. You may wish to set a valid value for Image Password even when using SSH Keys so that it is possible to authenticate with a password when using the remote console feature of the DCD.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option.

Required values to run command:

* Data Center Id
* Licence Type/Image Id or Image Alias
* Size`,
		Example:    createVolumeExample,
		PreCmdRun:  PreRunGlobalDcIdVolumeProperties,
		CmdRun:     RunVolumeCreate,
		InitClient: true,
	})
	create.AddStringFlag(config.ArgName, config.ArgNameShort, "Unnamed Storage", "Name of the Volume")
	create.AddStringFlag(config.ArgSize, "", strconv.Itoa(config.DefaultVolumeSize), "The size of the Volume in GB. e.g.: --size 10 or --size 10GB. The maximum Volume size is determined by your contract limit "+config.RequiredFlag)
	_ = create.Command.RegisterFlagCompletionFunc(config.ArgSize, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"10GB", "20GB", "50GB", "100GB", "1TB"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(config.ArgBus, "", "VIRTIO", "The bus type of the Volume")
	_ = create.Command.RegisterFlagCompletionFunc(config.ArgBus, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"VIRTIO", "IDE"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(config.ArgLicenceType, "", "", "Licence Type of the Volume")
	_ = create.Command.RegisterFlagCompletionFunc(config.ArgLicenceType, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"LINUX", "WINDOWS", "WINDOWS2016", "UNKNOWN", "OTHER"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(config.ArgType, "", "HDD", "Type of the Volume")
	_ = create.Command.RegisterFlagCompletionFunc(config.ArgLicenceType, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"HDD", "SSD", "SSD Standard", "SSD Premium"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(config.ArgAvailabilityZone, config.ArgAvailabilityZoneShort, "AUTO", "Availability zone of the Volume. Storage zone can only be selected prior provisioning")
	_ = create.Command.RegisterFlagCompletionFunc(config.ArgAvailabilityZone, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"AUTO", "ZONE_1", "ZONE_2", "ZONE_3"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(config.ArgBackupUnitId, "", "", "The unique Id of the Backup Unit that User has access to. It is mandatory to provide either 'public image' or 'imageAlias' in conjunction with this property")
	_ = create.Command.RegisterFlagCompletionFunc(config.ArgBackupUnitId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getBackupUnitsIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(config.ArgImageId, "", "", "The Image Id or snapshot Id to be used as template for the new Volume")
	_ = create.Command.RegisterFlagCompletionFunc(config.ArgImageId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getImageIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(config.ArgImageAlias, "", "", "The Image Alias to set instead of Image Id")
	create.AddStringFlag(config.ArgPassword, config.ArgPasswordShort, "", "Initial password to be set for installed OS. Works with public Images only. Not modifiable. Password rules allows all characters from a-z, A-Z, 0-9")
	create.AddStringFlag(config.ArgUserData, "", "", "The cloud-init configuration for the Volume as base64 encoded string. It is mandatory to provide either 'public image' or 'imageAlias' that has cloud-init compatibility in conjunction with this property")
	create.AddBoolFlag(config.ArgCpuHotPlug, "", false, "It is capable of CPU hot plug (no reboot required)")
	create.AddBoolFlag(config.ArgRamHotPlug, "", false, "It is capable of memory hot plug (no reboot required)")
	create.AddBoolFlag(config.ArgNicHotPlug, "", false, "It is capable of nic hot plug (no reboot required)")
	create.AddBoolFlag(config.ArgNicHotUnplug, "", false, "It is capable of nic hot unplug (no reboot required)")
	create.AddBoolFlag(config.ArgDiscVirtioHotPlug, "", false, "It is capable of Virt-IO drive hot plug (no reboot required)")
	create.AddBoolFlag(config.ArgDiscVirtioHotUnplug, "", false, "It is capable of Virt-IO drive hot unplug (no reboot required). This works only for non-Windows virtual Machines")
	create.AddStringSliceFlag(config.ArgSshKeys, "", []string{""}, "SSH Keys of the Volume")
	create.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for Volume creation to be executed")
	create.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option for Request for Volume creation [seconds]")

	/*
		Update Command
	*/
	update := core.NewCommand(ctx, volumeCmd, core.CommandBuilder{
		Namespace: "volume",
		Resource:  "volume",
		Verb:      "update",
		Aliases:   []string{"u", "up"},
		ShortDesc: "Update a Volume",
		LongDesc: `Use this command to update a Volume. You may increase the size of an existing storage Volume. You cannot reduce the size of an existing storage Volume. The Volume size will be increased without reboot if the appropriate "hot plug" settings have been set to true. The additional capacity is not added to any partition therefore you will need to adjust the partition inside the operating system afterwards.

Once you have increased the Volume size you cannot decrease the Volume size using the Cloud API. Certain attributes can only be set when a Volume is created and are considered immutable once the Volume has been provisioned.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option.

Required values to run command:

* Data Center Id
* Volume Id`,
		Example:    updateVolumeExample,
		PreCmdRun:  PreRunGlobalDcIdVolumeId,
		CmdRun:     RunVolumeUpdate,
		InitClient: true,
	})
	update.AddStringFlag(config.ArgVolumeId, config.ArgIdShort, "", config.RequiredFlagVolumeId)
	_ = update.Command.RegisterFlagCompletionFunc(config.ArgVolumeId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getVolumesIds(os.Stderr, viper.GetString(core.GetGlobalFlagName(volumeCmd.Name(), config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(config.ArgName, config.ArgNameShort, "", "Name of the Volume")
	update.AddStringFlag(config.ArgSize, "", strconv.Itoa(config.DefaultVolumeSize), "The size of the Volume in GB. e.g. 10 or 10GB. The maximum volume size is determined by your contract limit "+config.RequiredFlag)
	_ = update.Command.RegisterFlagCompletionFunc(config.ArgSize, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"10GB", "20GB", "50GB", "100GB", "1TB"}, cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(config.ArgBus, "", "VIRTIO", "Bus of the Volume")
	update.AddBoolFlag(config.ArgCpuHotPlug, "", false, "It is capable of CPU hot plug (no reboot required)")
	update.AddBoolFlag(config.ArgRamHotPlug, "", false, "It is capable of memory hot plug (no reboot required)")
	update.AddBoolFlag(config.ArgNicHotPlug, "", false, "It is capable of nic hot plug (no reboot required)")
	update.AddBoolFlag(config.ArgNicHotUnplug, "", false, "It is capable of nic hot unplug (no reboot required)")
	update.AddBoolFlag(config.ArgDiscVirtioHotPlug, "", false, "It is capable of Virt-IO drive hot plug (no reboot required)")
	update.AddBoolFlag(config.ArgDiscVirtioHotUnplug, "", false, "It is capable of Virt-IO drive hot unplug (no reboot required). This works only for non-Windows virtual Machines")
	update.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for Volume update to be executed")
	update.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option for Request for Volume update [seconds]")

	/*
		Delete Command
	*/
	deleteCmd := core.NewCommand(ctx, volumeCmd, core.CommandBuilder{
		Namespace: "volume",
		Resource:  "volume",
		Verb:      "delete",
		Aliases:   []string{"d"},
		ShortDesc: "Delete a Volume",
		LongDesc: `Use this command to delete specified Volume. This will result in the Volume being removed from your Virtual Data Center. Please use this with caution!

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option. You can force the command to execute without user input using ` + "`" + `--force` + "`" + ` option.

Required values to run command:

* Data Center Id
* Volume Id`,
		Example:    deleteVolumeExample,
		PreCmdRun:  PreRunGlobalDcIdVolumeId,
		CmdRun:     RunVolumeDelete,
		InitClient: true,
	})
	deleteCmd.AddStringFlag(config.ArgVolumeId, config.ArgIdShort, "", config.RequiredFlagVolumeId)
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(config.ArgVolumeId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getVolumesIds(os.Stderr, viper.GetString(core.GetGlobalFlagName(volumeCmd.Name(), config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for Volume deletion to be executed")
	deleteCmd.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option for Request for Volume deletion [seconds]")

	return volumeCmd
}

func PreRunGlobalDcIdVolumeProperties(c *core.PreCommandConfig) error {
	var result error
	if err := core.CheckRequiredGlobalFlags(c.Resource, config.ArgDataCenterId); err != nil {
		result = multierror.Append(result, err)
	}
	if err := core.CheckRequiredFlags(c.NS, config.ArgSize); err != nil {
		result = multierror.Append(result, err)
	}
	// Check required flags
	if !viper.IsSet(core.GetFlagName(c.NS, config.ArgLicenceType)) {
		if !viper.IsSet(core.GetFlagName(c.NS, config.ArgImageId)) &&
			!viper.IsSet(core.GetFlagName(c.NS, config.ArgImageAlias)) {
			result = multierror.Append(result, errors.New("image-id, image-alias or licence-type option must be set"))
		} else {
			if !viper.IsSet(core.GetFlagName(c.NS, config.ArgPassword)) &&
				!viper.IsSet(core.GetFlagName(c.NS, config.ArgSshKeys)) {
				result = multierror.Append(result, errors.New("image-password or ssh-keys option must be set"))
			}
		}
	}
	if result != nil {
		return result
	}
	return nil
}

func PreRunGlobalDcIdVolumeId(c *core.PreCommandConfig) error {
	var result error
	if err := core.CheckRequiredGlobalFlags(c.Resource, config.ArgDataCenterId); err != nil {
		result = multierror.Append(result, err)
	}
	if err := core.CheckRequiredFlags(c.NS, config.ArgVolumeId); err != nil {
		result = multierror.Append(result, err)
	}
	if result != nil {
		return result
	}
	return nil
}

func RunVolumeList(c *core.CommandConfig) error {
	volumes, _, err := c.Volumes().List(viper.GetString(core.GetGlobalFlagName(c.Resource, config.ArgDataCenterId)))
	if err != nil {
		return err
	}
	return c.Printer.Print(getVolumePrint(nil, c, getVolumes(volumes)))
}

func RunVolumeGet(c *core.CommandConfig) error {
	vol, _, err := c.Volumes().Get(
		viper.GetString(core.GetGlobalFlagName(c.Resource, config.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgVolumeId)),
	)
	if err != nil {
		return err
	}
	return c.Printer.Print(getVolumePrint(nil, c, getVolume(vol)))
}

func RunVolumeCreate(c *core.CommandConfig) error {
	input, err := getNewVolume(c)
	if err != nil {
		return err
	}
	vol, resp, err := c.Volumes().Create(
		viper.GetString(core.GetGlobalFlagName(c.Resource, config.ArgDataCenterId)),
		*input,
	)
	if err != nil {
		return err
	}
	if err = utils.WaitForRequest(c, printer.GetRequestPath(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getVolumePrint(resp, c, getVolume(vol)))
}

func RunVolumeUpdate(c *core.CommandConfig) error {
	input, err := getVolumeInfo(c)
	if err != nil {
		return err
	}
	vol, resp, err := c.Volumes().Update(
		viper.GetString(core.GetGlobalFlagName(c.Resource, config.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgVolumeId)),
		*input,
	)
	if err != nil {
		return err
	}
	if err = utils.WaitForRequest(c, printer.GetRequestPath(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getVolumePrint(resp, c, getVolume(vol)))
}

func RunVolumeDelete(c *core.CommandConfig) error {
	if err := utils.AskForConfirm(c.Stdin, c.Printer, "delete volume"); err != nil {
		return err
	}
	resp, err := c.Volumes().Delete(
		viper.GetString(core.GetGlobalFlagName(c.Resource, config.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgVolumeId)),
	)
	if err != nil {
		return err
	}

	if err = utils.WaitForRequest(c, printer.GetRequestPath(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getVolumePrint(resp, c, nil))
}

func getNewVolume(c *core.CommandConfig) (*resources.Volume, error) {
	proper := resources.VolumeProperties{}
	// It will get the default values, if flags not set
	proper.SetName(viper.GetString(core.GetFlagName(c.NS, config.ArgName)))
	proper.SetBus(viper.GetString(core.GetFlagName(c.NS, config.ArgBus)))
	proper.SetType(viper.GetString(core.GetFlagName(c.NS, config.ArgType)))
	proper.SetAvailabilityZone(viper.GetString(core.GetFlagName(c.NS, config.ArgAvailabilityZone)))
	if viper.IsSet(core.GetFlagName(c.NS, config.ArgSize)) {
		size, err := utils.ConvertSize(
			viper.GetString(core.GetFlagName(c.NS, config.ArgSize)),
			utils.GigaBytes,
		)
		if err != nil {
			return nil, err
		}
		proper.SetSize(float32(size))
	}
	// Check if flags are set and set options
	if viper.IsSet(core.GetFlagName(c.NS, config.ArgBackupUnitId)) {
		proper.SetBackupunitId(viper.GetString(core.GetFlagName(c.NS, config.ArgBackupUnitId)))
	}
	if viper.IsSet(core.GetFlagName(c.NS, config.ArgLicenceType)) {
		proper.SetLicenceType(viper.GetString(core.GetFlagName(c.NS, config.ArgLicenceType)))
	}
	if viper.IsSet(core.GetFlagName(c.NS, config.ArgImageId)) {
		proper.SetImage(viper.GetString(core.GetFlagName(c.NS, config.ArgImageId)))
	}
	if viper.IsSet(core.GetFlagName(c.NS, config.ArgImageAlias)) {
		proper.SetImageAlias(viper.GetString(core.GetFlagName(c.NS, config.ArgImageAlias)))
	}
	if viper.IsSet(core.GetFlagName(c.NS, config.ArgPassword)) {
		proper.SetImagePassword(viper.GetString(core.GetFlagName(c.NS, config.ArgPassword)))
	}
	if viper.IsSet(core.GetFlagName(c.NS, config.ArgSshKeys)) {
		proper.SetSshKeys(viper.GetStringSlice(core.GetFlagName(c.NS, config.ArgSshKeys)))
	}
	if viper.IsSet(core.GetFlagName(c.NS, config.ArgUserData)) {
		proper.SetUserData(viper.GetString(core.GetFlagName(c.NS, config.ArgUserData)))
	}
	if viper.IsSet(core.GetFlagName(c.NS, config.ArgCpuHotPlug)) {
		proper.SetCpuHotPlug(viper.GetBool(core.GetFlagName(c.NS, config.ArgCpuHotPlug)))
	}
	if viper.IsSet(core.GetFlagName(c.NS, config.ArgRamHotPlug)) {
		proper.SetRamHotPlug(viper.GetBool(core.GetFlagName(c.NS, config.ArgRamHotPlug)))
	}
	if viper.IsSet(core.GetFlagName(c.NS, config.ArgNicHotPlug)) {
		proper.SetNicHotPlug(viper.GetBool(core.GetFlagName(c.NS, config.ArgNicHotPlug)))
	}
	if viper.IsSet(core.GetFlagName(c.NS, config.ArgNicHotUnplug)) {
		proper.SetNicHotUnplug(viper.GetBool(core.GetFlagName(c.NS, config.ArgNicHotUnplug)))
	}
	if viper.IsSet(core.GetFlagName(c.NS, config.ArgDiscVirtioHotPlug)) {
		proper.SetDiscVirtioHotPlug(viper.GetBool(core.GetFlagName(c.NS, config.ArgDiscVirtioHotPlug)))
	}
	if viper.IsSet(core.GetFlagName(c.NS, config.ArgDiscVirtioHotUnplug)) {
		proper.SetDiscVirtioHotUnplug(viper.GetBool(core.GetFlagName(c.NS, config.ArgDiscVirtioHotUnplug)))
	}
	return &resources.Volume{
		Volume: ionoscloud.Volume{
			Properties: &proper.VolumeProperties,
		},
	}, nil
}

func getVolumeInfo(c *core.CommandConfig) (*resources.VolumeProperties, error) {
	input := resources.VolumeProperties{}
	if viper.IsSet(core.GetFlagName(c.NS, config.ArgName)) {
		input.SetName(viper.GetString(core.GetFlagName(c.NS, config.ArgName)))
	}
	if viper.IsSet(core.GetFlagName(c.NS, config.ArgBus)) {
		input.SetBus(viper.GetString(core.GetFlagName(c.NS, config.ArgBus)))
	}
	if viper.IsSet(core.GetFlagName(c.NS, config.ArgSize)) {
		size, err := utils.ConvertSize(
			viper.GetString(core.GetFlagName(c.NS, config.ArgSize)),
			utils.GigaBytes,
		)
		if err != nil {
			return nil, err
		}
		input.SetSize(float32(size))
	}
	// Check if flags are set and set options
	if viper.IsSet(core.GetFlagName(c.NS, config.ArgCpuHotPlug)) {
		input.SetCpuHotPlug(viper.GetBool(core.GetFlagName(c.NS, config.ArgCpuHotPlug)))
	}
	if viper.IsSet(core.GetFlagName(c.NS, config.ArgRamHotPlug)) {
		input.SetRamHotPlug(viper.GetBool(core.GetFlagName(c.NS, config.ArgRamHotPlug)))
	}
	if viper.IsSet(core.GetFlagName(c.NS, config.ArgNicHotPlug)) {
		input.SetNicHotPlug(viper.GetBool(core.GetFlagName(c.NS, config.ArgNicHotPlug)))
	}
	if viper.IsSet(core.GetFlagName(c.NS, config.ArgNicHotUnplug)) {
		input.SetNicHotUnplug(viper.GetBool(core.GetFlagName(c.NS, config.ArgNicHotUnplug)))
	}
	if viper.IsSet(core.GetFlagName(c.NS, config.ArgDiscVirtioHotPlug)) {
		input.SetDiscVirtioHotPlug(viper.GetBool(core.GetFlagName(c.NS, config.ArgDiscVirtioHotPlug)))
	}
	if viper.IsSet(core.GetFlagName(c.NS, config.ArgDiscVirtioHotUnplug)) {
		input.SetDiscVirtioHotUnplug(viper.GetBool(core.GetFlagName(c.NS, config.ArgDiscVirtioHotUnplug)))
	}
	return &input, nil
}

// Server Volume Commands

func serverVolume() *core.Command {
	ctx := context.TODO()
	serverVolumeCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "volume",
			Aliases:          []string{"v", "vol"},
			Short:            "Server Volume Operations",
			Long:             `The sub-commands of ` + "`" + `ionosctl server volume` + "`" + ` allow you to attach, get, list, detach Volumes from Servers.`,
			TraverseChildren: true,
		},
	}

	/*
		Attach Volume Command
	*/
	attachVolume := core.NewCommand(ctx, serverVolumeCmd, core.CommandBuilder{
		Namespace: "server",
		Resource:  "volume",
		Verb:      "attach",
		Aliases:   []string{"a"},
		ShortDesc: "Attach a Volume to a Server",
		LongDesc: `Use this command to attach a pre-existing Volume to a Server.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option.

Required values to run command:

* Data Center Id
* Server Id
* Volume Id`,
		Example:    attachVolumeServerExample,
		PreCmdRun:  PreRunDcServerVolumeIds,
		CmdRun:     RunServerVolumeAttach,
		InitClient: true,
	})
	attachVolume.AddStringSliceFlag(config.ArgCols, "", defaultVolumeCols, utils.ColsMessage(allVolumeCols))
	_ = attachVolume.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allVolumeCols, cobra.ShellCompDirectiveNoFileComp
	})
	attachVolume.AddStringFlag(config.ArgDataCenterId, "", "", config.RequiredFlagDatacenterId)
	_ = attachVolume.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	attachVolume.AddStringFlag(config.ArgVolumeId, config.ArgIdShort, "", config.RequiredFlagVolumeId)
	_ = attachVolume.Command.RegisterFlagCompletionFunc(config.ArgVolumeId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getVolumesIds(os.Stderr, viper.GetString(core.GetFlagName(attachVolume.NS, config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	attachVolume.AddStringFlag(config.ArgServerId, "", "", config.RequiredFlagServerId)
	_ = attachVolume.Command.RegisterFlagCompletionFunc(config.ArgServerId, func(cmd *cobra.Command, ags []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getServersIds(os.Stderr, viper.GetString(core.GetFlagName(attachVolume.NS, config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	attachVolume.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for Volume attachment to be executed")
	attachVolume.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option for Request for Volume attachment [seconds]")

	/*
		List Volumes Command
	*/
	listVolumes := core.NewCommand(ctx, serverVolumeCmd, core.CommandBuilder{
		Namespace:  "server",
		Resource:   "volume",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List attached Volumes from a Server",
		LongDesc:   "Use this command to retrieve a list of Volumes attached to the Server.\n\nRequired values to run command:\n\n* Data Center Id\n* Server Id",
		Example:    listVolumesServerExample,
		PreCmdRun:  PreRunDcServerIds,
		CmdRun:     RunServerVolumesList,
		InitClient: true,
	})
	listVolumes.AddStringSliceFlag(config.ArgCols, "", defaultVolumeCols, utils.ColsMessage(allVolumeCols))
	_ = listVolumes.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allVolumeCols, cobra.ShellCompDirectiveNoFileComp
	})
	listVolumes.AddStringFlag(config.ArgDataCenterId, "", "", config.RequiredFlagDatacenterId)
	_ = listVolumes.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	listVolumes.AddStringFlag(config.ArgServerId, "", "", config.RequiredFlagServerId)
	_ = listVolumes.Command.RegisterFlagCompletionFunc(config.ArgServerId, func(cmd *cobra.Command, ags []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getServersIds(os.Stderr, viper.GetString(core.GetFlagName(listVolumes.NS, config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Get Volume Command
	*/
	getVolumeCmd := core.NewCommand(ctx, serverVolumeCmd, core.CommandBuilder{
		Namespace:  "server",
		Resource:   "volume",
		Verb:       "get",
		Aliases:    []string{"g"},
		ShortDesc:  "Get an attached Volume from a Server",
		LongDesc:   "Use this command to retrieve information about an attached Volume on Server.\n\nRequired values to run command:\n\n* Data Center Id\n* Server Id\n* Volume Id",
		Example:    getVolumeServerExample,
		InitClient: true,
		PreCmdRun:  PreRunDcServerVolumeIds,
		CmdRun:     RunServerVolumeGet,
	})
	getVolumeCmd.AddStringSliceFlag(config.ArgCols, "", defaultVolumeCols, utils.ColsMessage(allVolumeCols))
	_ = getVolumeCmd.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allVolumeCols, cobra.ShellCompDirectiveNoFileComp
	})
	getVolumeCmd.AddStringFlag(config.ArgDataCenterId, "", "", config.RequiredFlagDatacenterId)
	_ = getVolumeCmd.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	getVolumeCmd.AddStringFlag(config.ArgServerId, "", "", config.RequiredFlagServerId)
	_ = getVolumeCmd.Command.RegisterFlagCompletionFunc(config.ArgServerId, func(cmd *cobra.Command, ags []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getServersIds(os.Stderr, viper.GetString(core.GetFlagName(getVolumeCmd.NS, config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	getVolumeCmd.AddStringFlag(config.ArgVolumeId, config.ArgIdShort, "", config.RequiredFlagVolumeId)
	_ = getVolumeCmd.Command.RegisterFlagCompletionFunc(config.ArgVolumeId, func(cmd *cobra.Command, ags []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getAttachedVolumesIds(os.Stderr, viper.GetString(core.GetFlagName(getVolumeCmd.NS, config.ArgDataCenterId)), viper.GetString(core.GetFlagName(getVolumeCmd.NS, config.ArgServerId))), cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Detach Volume Command
	*/
	detachVolume := core.NewCommand(ctx, serverVolumeCmd, core.CommandBuilder{
		Namespace: "server",
		Resource:  "volume",
		Verb:      "detach",
		Aliases:   []string{"d"},
		ShortDesc: "Detach a Volume from a Server",
		LongDesc: `This will detach the Volume from the Server. Depending on the Volume HotUnplug settings, this may result in the Server being rebooted. This will NOT delete the Volume from your Virtual Data Center. You will need to use a separate command to delete a Volume.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option. You can force the command to execute without user input using ` + "`" + `--force` + "`" + ` option.

Required values to run command:

* Data Center Id
* Server Id
* Volume Id`,
		Example:    detachVolumeServerExample,
		PreCmdRun:  PreRunDcServerVolumeIds,
		CmdRun:     RunServerVolumeDetach,
		InitClient: true,
	})
	detachVolume.AddStringSliceFlag(config.ArgCols, "", defaultVolumeCols, utils.ColsMessage(allVolumeCols))
	_ = detachVolume.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allVolumeCols, cobra.ShellCompDirectiveNoFileComp
	})
	detachVolume.AddStringFlag(config.ArgDataCenterId, "", "", config.RequiredFlagDatacenterId)
	_ = detachVolume.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	detachVolume.AddStringFlag(config.ArgVolumeId, config.ArgIdShort, "", config.RequiredFlagVolumeId)
	_ = detachVolume.Command.RegisterFlagCompletionFunc(config.ArgVolumeId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getAttachedVolumesIds(os.Stderr, viper.GetString(core.GetFlagName(detachVolume.NS, config.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(detachVolume.NS, config.ArgServerId))), cobra.ShellCompDirectiveNoFileComp
	})
	detachVolume.AddStringFlag(config.ArgServerId, "", "", config.RequiredFlagServerId)
	_ = detachVolume.Command.RegisterFlagCompletionFunc(config.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getServersIds(os.Stderr, viper.GetString(core.GetFlagName(detachVolume.NS, config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	detachVolume.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for Volume detachment to be executed")
	detachVolume.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option for Request for Volume detachment [seconds]")

	return serverVolumeCmd
}

func PreRunDcServerVolumeIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.NS, config.ArgDataCenterId, config.ArgServerId, config.ArgVolumeId)
}

func RunServerVolumeAttach(c *core.CommandConfig) error {
	attachedVol, resp, err := c.Servers().AttachVolume(
		viper.GetString(core.GetFlagName(c.NS, config.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgServerId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgVolumeId)),
	)
	if err != nil {
		return err
	}

	if err = utils.WaitForRequest(c, printer.GetRequestPath(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getVolumePrint(resp, c, getVolume(attachedVol)))
}

func RunServerVolumesList(c *core.CommandConfig) error {
	attachedVols, _, err := c.Servers().ListVolumes(
		viper.GetString(core.GetFlagName(c.NS, config.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgServerId)),
	)
	if err != nil {
		return err
	}
	return c.Printer.Print(getVolumePrint(nil, c, getAttachedVolumes(attachedVols)))
}

func RunServerVolumeGet(c *core.CommandConfig) error {
	attachedVol, _, err := c.Servers().GetVolume(
		viper.GetString(core.GetFlagName(c.NS, config.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgServerId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgVolumeId)),
	)
	if err != nil {
		return err
	}
	return c.Printer.Print(getVolumePrint(nil, c, getVolume(attachedVol)))
}

func RunServerVolumeDetach(c *core.CommandConfig) error {
	if err := utils.AskForConfirm(c.Stdin, c.Printer, "detach volume from server"); err != nil {
		return err
	}
	resp, err := c.Servers().DetachVolume(
		viper.GetString(core.GetFlagName(c.NS, config.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgServerId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgVolumeId)),
	)
	if err != nil {
		return err
	}

	if err = utils.WaitForRequest(c, printer.GetRequestPath(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getVolumePrint(resp, c, nil))
}

// Output Printing

var (
	defaultVolumeCols = []string{"VolumeId", "Name", "Size", "Type", "LicenceType", "State", "Image"}
	allVolumeCols     = []string{"VolumeId", "Name", "Size", "Type", "LicenceType", "State", "Image", "Bus", "AvailabilityZone", "BackupunitId", "SshKeys",
		"ImageAlias", "DeviceNumber", "UserData"}
)

type VolumePrint struct {
	VolumeId         string   `json:"VolumeId,omitempty"`
	Name             string   `json:"Name,omitempty"`
	Size             string   `json:"Size,omitempty"`
	Type             string   `json:"Type,omitempty"`
	LicenceType      string   `json:"LicenceType,omitempty"`
	Bus              string   `json:"Bus,omitempty"`
	AvailabilityZone string   `json:"AvailabilityZone,omitempty"`
	State            string   `json:"State,omitempty"`
	Image            string   `json:"Image,omitempty"`
	ImageAlias       string   `json:"ImageAlias,omitempty"`
	SshKeys          []string `json:"SshKeys,omitempty"`
	DeviceNumber     int64    `json:"DeviceNumber,omitempty"`
	BackupunitId     string   `json:"BackupunitId,omitempty"`
	UserData         string   `json:"UserData,omitempty"`
}

func getVolumePrint(resp *resources.Response, c *core.CommandConfig, vols []resources.Volume) printer.Result {
	r := printer.Result{}
	if c != nil {
		if resp != nil {
			r.ApiResponse = resp
			r.Resource = c.Resource
			r.Verb = c.Verb
			r.WaitForRequest = viper.GetBool(core.GetFlagName(c.NS, config.ArgWaitForRequest))
		}
		if vols != nil {
			r.OutputJSON = vols
			r.KeyValue = getVolumesKVMaps(vols)
			if c.Resource != c.Namespace {
				r.Columns = getVolumesCols(core.GetFlagName(c.NS, config.ArgCols), c.Printer.GetStderr())
			} else {
				r.Columns = getVolumesCols(core.GetGlobalFlagName(c.Resource, config.ArgCols), c.Printer.GetStderr())
			}
		}
	}
	return r
}

func getVolumesCols(flagName string, outErr io.Writer) []string {
	var cols []string
	if viper.IsSet(flagName) {
		cols = viper.GetStringSlice(flagName)
	} else {
		return defaultVolumeCols
	}

	columnsMap := map[string]string{
		"VolumeId":         "VolumeId",
		"Name":             "Name",
		"Size":             "Size",
		"Type":             "Type",
		"LicenceType":      "LicenceType",
		"Bus":              "Bus",
		"AvailabilityZone": "AvailabilityZone",
		"State":            "State",
		"Image":            "Image",
		"ImageAlias":       "ImageAlias",
		"SshKeys":          "SshKeys",
		"DeviceNumber":     "DeviceNumber",
		"BackupunitId":     "BackupunitId",
		"UserData":         "UserData",
	}
	var volumeCols []string
	for _, k := range cols {
		col := columnsMap[k]
		if col != "" {
			volumeCols = append(volumeCols, col)
		} else {
			clierror.CheckError(errors.New("unknown column "+k), outErr)
		}
	}
	return volumeCols
}

func getVolumes(volumes resources.Volumes) []resources.Volume {
	vs := make([]resources.Volume, 0)
	for _, s := range *volumes.Items {
		vs = append(vs, resources.Volume{Volume: s})
	}
	return vs
}

func getVolume(vol *resources.Volume) []resources.Volume {
	vols := make([]resources.Volume, 0)
	if vol != nil {
		vols = append(vols, resources.Volume{Volume: vol.Volume})
	}
	return vols
}

func getAttachedVolumes(volumes resources.AttachedVolumes) []resources.Volume {
	vs := make([]resources.Volume, 0)
	for _, s := range *volumes.AttachedVolumes.Items {
		vs = append(vs, resources.Volume{Volume: s})
	}
	return vs
}

func getVolumesKVMaps(vs []resources.Volume) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(vs))
	for _, v := range vs {
		properties := v.GetProperties()
		metadata := v.GetMetadata()
		var volumePrint VolumePrint
		if id, ok := v.GetIdOk(); ok && id != nil {
			volumePrint.VolumeId = *id
		}
		if name, ok := properties.GetNameOk(); ok && name != nil {
			volumePrint.Name = *name
		}
		if licenceType, ok := properties.GetLicenceTypeOk(); ok && licenceType != nil {
			volumePrint.LicenceType = *licenceType
		}
		if size, ok := properties.GetSizeOk(); ok && size != nil {
			volumePrint.Size = fmt.Sprintf("%vGB", *size)
		}
		if bus, ok := properties.GetBusOk(); ok && bus != nil {
			volumePrint.Bus = *bus
		}
		if volumetype, ok := properties.GetTypeOk(); ok && volumetype != nil {
			volumePrint.Type = *volumetype
		}
		if zone, ok := properties.GetAvailabilityZoneOk(); ok && zone != nil {
			volumePrint.AvailabilityZone = *zone
		}
		if state, ok := metadata.GetStateOk(); ok && state != nil {
			volumePrint.State = *state
		}
		if backUpUnitId, ok := properties.GetBackupunitIdOk(); ok && backUpUnitId != nil {
			volumePrint.BackupunitId = *backUpUnitId
		}
		if img, ok := properties.GetImageOk(); ok && img != nil {
			volumePrint.Image = *img
		}
		if imgA, ok := properties.GetImageAliasOk(); ok && imgA != nil {
			volumePrint.ImageAlias = *imgA
		}
		if userData, ok := properties.GetUserDataOk(); ok && userData != nil {
			volumePrint.UserData = *userData
		}
		if no, ok := properties.GetDeviceNumberOk(); ok && no != nil {
			volumePrint.DeviceNumber = *no
		}
		if sshKeys, ok := properties.GetSshKeysOk(); ok && sshKeys != nil {
			volumePrint.SshKeys = *sshKeys
		}
		o := structs.Map(volumePrint)
		out = append(out, o)
	}
	return out
}

func getVolumesIds(outErr io.Writer, datacenterId string) []string {
	err := config.Load()
	clierror.CheckError(err, outErr)
	clientSvc, err := resources.NewClientService(
		viper.GetString(config.Username),
		viper.GetString(config.Password),
		viper.GetString(config.Token),
		viper.GetString(config.ArgServerUrl),
	)
	clierror.CheckError(err, outErr)
	volumeSvc := resources.NewVolumeService(clientSvc.Get(), context.TODO())
	volumes, _, err := volumeSvc.List(datacenterId)
	clierror.CheckError(err, outErr)
	volumesIds := make([]string, 0)
	if items, ok := volumes.Volumes.GetItemsOk(); ok && items != nil {
		for _, item := range *items {
			if itemId, ok := item.GetIdOk(); ok && itemId != nil {
				volumesIds = append(volumesIds, *itemId)
			}
		}
	} else {
		return nil
	}
	return volumesIds
}

func getAttachedVolumesIds(outErr io.Writer, datacenterId, serverId string) []string {
	err := config.Load()
	clierror.CheckError(err, outErr)
	clientSvc, err := resources.NewClientService(
		viper.GetString(config.Username),
		viper.GetString(config.Password),
		viper.GetString(config.Token),
		viper.GetString(config.ArgServerUrl),
	)
	clierror.CheckError(err, outErr)
	serverSvc := resources.NewServerService(clientSvc.Get(), context.TODO())
	volumes, _, err := serverSvc.ListVolumes(datacenterId, serverId)
	clierror.CheckError(err, outErr)
	attachedVolumesIds := make([]string, 0)
	if items, ok := volumes.AttachedVolumes.GetItemsOk(); ok && items != nil {
		for _, item := range *items {
			if itemId, ok := item.GetIdOk(); ok && itemId != nil {
				attachedVolumesIds = append(attachedVolumesIds, *itemId)
			}
		}
	} else {
		return nil
	}
	return attachedVolumesIds
}
