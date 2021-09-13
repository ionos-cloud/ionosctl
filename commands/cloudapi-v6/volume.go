package commands

import (
	"context"
	"errors"
	"fmt"
	"github.com/ionos-cloud/ionosctl/commands/cloudapi-v6/completer"
	"github.com/ionos-cloud/ionosctl/commands/cloudapi-v6/waiter"
	cloudapi_v6 "github.com/ionos-cloud/ionosctl/services/cloudapi-v6"
	"io"
	"os"
	"strconv"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/internal/config"
	"github.com/ionos-cloud/ionosctl/internal/core"
	"github.com/ionos-cloud/ionosctl/internal/printer"
	"github.com/ionos-cloud/ionosctl/internal/utils"
	"github.com/ionos-cloud/ionosctl/internal/utils/clierror"
	"github.com/ionos-cloud/ionosctl/services/cloudapi-v6/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func VolumeCmd() *core.Command {
	ctx := context.TODO()
	volumeCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "volume",
			Aliases:          []string{"v", "vol"},
			Short:            "Volume Operations",
			Long:             "The sub-commands of `ionosctl volume` manage your block storage volumes by creating, updating, getting specific information, deleting Volumes. To attach a Volume to a Server, use the Server command `ionosctl server volume attach`.",
			TraverseChildren: true,
		},
	}
	globalFlags := volumeCmd.GlobalFlags()
	globalFlags.StringSliceP(config.ArgCols, "", defaultVolumeCols, printer.ColsMessage(allVolumeCols))
	_ = viper.BindPFlag(core.GetGlobalFlagName(volumeCmd.Name(), config.ArgCols), globalFlags.Lookup(config.ArgCols))
	_ = volumeCmd.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allVolumeCols, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		List Command
	*/
	list := core.NewCommand(ctx, volumeCmd, core.CommandBuilder{
		Namespace:  "volume",
		Resource:   "volume",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List Volumes",
		LongDesc:   "Use this command to list all Volumes from a Data Center on your account.\n\nRequired values to run command:\n\n* Data Center Id",
		Example:    listVolumeExample,
		PreCmdRun:  PreRunDataCenterId,
		CmdRun:     RunVolumeList,
		InitClient: true,
	})
	list.AddStringFlag(cloudapi_v6.ArgDataCenterId, "", "", cloudapi_v6.DatacenterId, core.RequiredFlagOption())
	_ = list.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
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
		PreCmdRun:  PreRunDcVolumeIds,
		CmdRun:     RunVolumeGet,
		InitClient: true,
	})
	get.AddStringFlag(cloudapi_v6.ArgDataCenterId, "", "", cloudapi_v6.DatacenterId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddStringFlag(cloudapi_v6.ArgVolumeId, cloudapi_v6.ArgIdShort, "", cloudapi_v6.VolumeId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgVolumeId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.VolumesIds(os.Stderr, viper.GetString(core.GetFlagName(get.NS, cloudapi_v6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
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

NNote: The Licence Type has a default value, but if Image ID or Image Alias is supplied, then Licence Type will be automatically set. The Image Password or SSH Keys attributes can be defined when creating a Volume that uses an Image ID or Image Alias of an IONOS public Image. You may wish to set a valid value for Image Password even when using SSH Keys so that it is possible to authenticate with a password when using the remote console feature of the DCD.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option.

Required values to run command:

* Data Center Id`,
		Example:    createVolumeExample,
		PreCmdRun:  PreRunDataCenterId,
		CmdRun:     RunVolumeCreate,
		InitClient: true,
	})
	create.AddStringFlag(cloudapi_v6.ArgDataCenterId, "", "", cloudapi_v6.DatacenterId, core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(cloudapi_v6.ArgName, cloudapi_v6.ArgNameShort, "Unnamed Volume", "Name of the Volume")
	create.AddStringFlag(cloudapi_v6.ArgSize, "", strconv.Itoa(cloudapi_v6.DefaultVolumeSize), "The size of the Volume in GB. e.g.: --size 10 or --size 10GB. The maximum Volume size is determined by your contract limit")
	_ = create.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgSize, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"10GB", "20GB", "50GB", "100GB", "1TB"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(cloudapi_v6.ArgBus, "", "VIRTIO", "The bus type of the Volume")
	_ = create.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgBus, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"VIRTIO", "IDE"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(cloudapi_v6.ArgLicenceType, "", "LINUX", "Licence Type of the Volume")
	_ = create.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgLicenceType, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"LINUX", "WINDOWS", "WINDOWS2016", "UNKNOWN", "OTHER"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(cloudapi_v6.ArgType, "", "HDD", "Type of the Volume")
	_ = create.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgLicenceType, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"HDD", "SSD", "SSD Standard", "SSD Premium"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(cloudapi_v6.ArgAvailabilityZone, cloudapi_v6.ArgAvailabilityZoneShort, "AUTO", "Availability zone of the Volume. Storage zone can only be selected prior provisioning")
	_ = create.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgAvailabilityZone, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"AUTO", "ZONE_1", "ZONE_2", "ZONE_3"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(cloudapi_v6.ArgBackupUnitId, "", "", "The unique Id of the Backup Unit that User has access to. It is mandatory to provide either 'public image' or 'imageAlias' in conjunction with this property")
	_ = create.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgBackupUnitId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.BackupUnitsIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(cloudapi_v6.ArgImageId, "", "", "The Image Id or Snapshot Id to be used as template for the new Volume")
	_ = create.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgImageId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ImageIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(cloudapi_v6.ArgImageAlias, "", "", "The Image Alias to set instead of Image Id")
	create.AddStringFlag(cloudapi_v6.ArgPassword, cloudapi_v6.ArgPasswordShort, "abcde12345", "Initial password to be set for installed OS. Works with public Images only. Not modifiable. Password rules allows all characters from a-z, A-Z, 0-9")
	create.AddStringFlag(cloudapi_v6.ArgUserData, "", "", "The cloud-init configuration for the Volume as base64 encoded string. It is mandatory to provide either 'public image' or 'imageAlias' that has cloud-init compatibility in conjunction with this property")
	create.AddBoolFlag(cloudapi_v6.ArgCpuHotPlug, "", false, "It is capable of CPU hot plug (no reboot required)")
	create.AddBoolFlag(cloudapi_v6.ArgRamHotPlug, "", false, "It is capable of memory hot plug (no reboot required)")
	create.AddBoolFlag(cloudapi_v6.ArgNicHotPlug, "", false, "It is capable of nic hot plug (no reboot required)")
	create.AddBoolFlag(cloudapi_v6.ArgNicHotUnplug, "", false, "It is capable of nic hot unplug (no reboot required)")
	create.AddBoolFlag(cloudapi_v6.ArgDiscVirtioHotPlug, "", false, "It is capable of Virt-IO drive hot plug (no reboot required)")
	create.AddBoolFlag(cloudapi_v6.ArgDiscVirtioHotUnplug, "", false, "It is capable of Virt-IO drive hot unplug (no reboot required). This works only for non-Windows virtual Machines")
	create.AddStringSliceFlag(cloudapi_v6.ArgSshKeys, "", []string{""}, "SSH Keys of the Volume")
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
		PreCmdRun:  PreRunDcVolumeIds,
		CmdRun:     RunVolumeUpdate,
		InitClient: true,
	})
	update.AddStringFlag(cloudapi_v6.ArgDataCenterId, "", "", cloudapi_v6.DatacenterId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(cloudapi_v6.ArgVolumeId, cloudapi_v6.ArgIdShort, "", cloudapi_v6.VolumeId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgVolumeId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.VolumesIds(os.Stderr, viper.GetString(core.GetFlagName(update.NS, cloudapi_v6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(cloudapi_v6.ArgName, cloudapi_v6.ArgNameShort, "", "Name of the Volume")
	update.AddStringFlag(cloudapi_v6.ArgSize, "", "", "The size of the Volume in GB. e.g. 10 or 10GB. The maximum volume size is determined by your contract limit")
	_ = update.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgSize, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"10GB", "20GB", "50GB", "100GB", "1TB"}, cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(cloudapi_v6.ArgBus, "", "VIRTIO", "Bus of the Volume")
	update.AddBoolFlag(cloudapi_v6.ArgCpuHotPlug, "", false, "It is capable of CPU hot plug (no reboot required)")
	update.AddBoolFlag(cloudapi_v6.ArgRamHotPlug, "", false, "It is capable of memory hot plug (no reboot required)")
	update.AddBoolFlag(cloudapi_v6.ArgNicHotPlug, "", false, "It is capable of nic hot plug (no reboot required)")
	update.AddBoolFlag(cloudapi_v6.ArgNicHotUnplug, "", false, "It is capable of nic hot unplug (no reboot required)")
	update.AddBoolFlag(cloudapi_v6.ArgDiscVirtioHotPlug, "", false, "It is capable of Virt-IO drive hot plug (no reboot required)")
	update.AddBoolFlag(cloudapi_v6.ArgDiscVirtioHotUnplug, "", false, "It is capable of Virt-IO drive hot unplug (no reboot required). This works only for non-Windows virtual Machines")
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
		PreCmdRun:  PreRunDcVolumeIds,
		CmdRun:     RunVolumeDelete,
		InitClient: true,
	})
	deleteCmd.AddStringFlag(cloudapi_v6.ArgDataCenterId, "", "", cloudapi_v6.DatacenterId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddStringFlag(cloudapi_v6.ArgVolumeId, cloudapi_v6.ArgIdShort, "", cloudapi_v6.VolumeId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgVolumeId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.VolumesIds(os.Stderr, viper.GetString(core.GetFlagName(deleteCmd.NS, cloudapi_v6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for Volume deletion to be executed")
	deleteCmd.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option for Request for Volume deletion [seconds]")

	return volumeCmd
}

func PreRunDcVolumeIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapi_v6.ArgDataCenterId, cloudapi_v6.ArgVolumeId)
}

func RunVolumeList(c *core.CommandConfig) error {
	c.Printer.Verbose("Listing Volumes from Datacenter with ID: %v", viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgDataCenterId)))
	volumes, _, err := c.CloudApiV6Services.Volumes().List(viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgDataCenterId)))
	if err != nil {
		return err
	}
	return c.Printer.Print(getVolumePrint(nil, c, getVolumes(volumes)))
}

func RunVolumeGet(c *core.CommandConfig) error {
	c.Printer.Verbose("Datacenter ID: %v", viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgDataCenterId)))
	c.Printer.Verbose("Volume with id: %v is getting...", viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgVolumeId)))
	vol, _, err := c.CloudApiV6Services.Volumes().Get(
		viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgVolumeId)),
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
	vol, resp, err := c.CloudApiV6Services.Volumes().Create(viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgDataCenterId)), *input)
	if resp != nil {
		c.Printer.Verbose("Request href: %v ", resp.Header.Get("location"))
	}
	if err != nil {
		return err
	}
	if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetRequestPath(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getVolumePrint(resp, c, getVolume(vol)))
}

func RunVolumeUpdate(c *core.CommandConfig) error {
	input, err := getVolumeInfo(c)
	if err != nil {
		return err
	}
	vol, resp, err := c.CloudApiV6Services.Volumes().Update(viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgVolumeId)), *input)
	if err != nil {
		return err
	}
	if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetRequestPath(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getVolumePrint(resp, c, getVolume(vol)))
}

func RunVolumeDelete(c *core.CommandConfig) error {
	if err := utils.AskForConfirm(c.Stdin, c.Printer, "delete volume"); err != nil {
		return err
	}
	c.Printer.Verbose("Volume with id: %v is deleting...", viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgVolumeId)))
	resp, err := c.CloudApiV6Services.Volumes().Delete(
		viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgVolumeId)),
	)
	if err != nil {
		return err
	}

	if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetRequestPath(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getVolumePrint(resp, c, nil))
}

func getNewVolume(c *core.CommandConfig) (*resources.Volume, error) {
	proper := resources.VolumeProperties{}
	name := viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgName))
	bus := viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgBus))
	volumeType := viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgType))
	availabilityZone := viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgAvailabilityZone))
	// It will get the default values, if flags not set
	proper.SetName(name)
	proper.SetBus(bus)
	proper.SetType(volumeType)
	proper.SetAvailabilityZone(availabilityZone)
	c.Printer.Verbose("Properties set for creating the Volume: Name: %v, Bus: %v, VolumeType: %v, AvailabilityZone: %v",
		name, bus, volumeType, availabilityZone)
	size, err := utils.ConvertSize(
		viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgSize)),
		utils.GigaBytes,
	)
	if err != nil {
		return nil, err
	}
	proper.SetSize(float32(size))
	c.Printer.Verbose("Property Size set: %vGB", float32(size))
	// Check if flags are set and set options
	if viper.IsSet(core.GetFlagName(c.NS, cloudapi_v6.ArgBackupUnitId)) {
		backupUnitId := viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgBackupUnitId))
		proper.SetBackupunitId(backupUnitId)
		c.Printer.Verbose("Property BackupUnitId set: %v", backupUnitId)
	}
	if !viper.IsSet(core.GetFlagName(c.NS, cloudapi_v6.ArgImageId)) &&
		!viper.IsSet(core.GetFlagName(c.NS, cloudapi_v6.ArgImageAlias)) ||
		viper.IsSet(core.GetFlagName(c.NS, cloudapi_v6.ArgLicenceType)) {
		licenceType := viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgLicenceType))
		proper.SetLicenceType(licenceType)
		c.Printer.Verbose("Property LicenceType set: %v", licenceType)
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapi_v6.ArgImageId)) {
		imageId := viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgImageId))
		proper.SetImage(imageId)
		c.Printer.Verbose("Property Image set: %v", imageId)
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapi_v6.ArgImageId)) ||
		viper.IsSet(core.GetFlagName(c.NS, cloudapi_v6.ArgImageAlias)) ||
		viper.IsSet(core.GetFlagName(c.NS, cloudapi_v6.ArgPassword)) {
		imagePassword := viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgPassword))
		proper.SetImagePassword(imagePassword)
		c.Printer.Verbose("Property ImagePassword set")
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapi_v6.ArgSshKeys)) {
		sshKeys := viper.GetStringSlice(core.GetFlagName(c.NS, cloudapi_v6.ArgSshKeys))
		proper.SetSshKeys(sshKeys)
		c.Printer.Verbose("Property SshKeys set: %v", sshKeys)
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapi_v6.ArgUserData)) {
		userData := viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgUserData))
		proper.SetUserData(userData)
		c.Printer.Verbose("Property UserData set: %v", userData)
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapi_v6.ArgCpuHotPlug)) {
		cpuHotPlug := viper.GetBool(core.GetFlagName(c.NS, cloudapi_v6.ArgCpuHotPlug))
		proper.SetCpuHotPlug(cpuHotPlug)
		c.Printer.Verbose("Property CpuHotPlug set: %v", cpuHotPlug)
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapi_v6.ArgRamHotPlug)) {
		ramHotPlug := viper.GetBool(core.GetFlagName(c.NS, cloudapi_v6.ArgRamHotPlug))
		proper.SetRamHotPlug(ramHotPlug)
		c.Printer.Verbose("Property RamHotPlug set: %v", ramHotPlug)
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapi_v6.ArgNicHotPlug)) {
		nicHotPlug := viper.GetBool(core.GetFlagName(c.NS, cloudapi_v6.ArgNicHotPlug))
		proper.SetNicHotPlug(nicHotPlug)
		c.Printer.Verbose("Property NicHotPlug set: %v", nicHotPlug)
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapi_v6.ArgNicHotUnplug)) {
		nicHotUnplug := viper.GetBool(core.GetFlagName(c.NS, cloudapi_v6.ArgNicHotUnplug))
		proper.SetNicHotUnplug(nicHotUnplug)
		c.Printer.Verbose("Property NicHotUnplug set: %v", nicHotUnplug)
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapi_v6.ArgDiscVirtioHotPlug)) {
		discVirtioHotPlug := viper.GetBool(core.GetFlagName(c.NS, cloudapi_v6.ArgDiscVirtioHotPlug))
		proper.SetDiscVirtioHotPlug(discVirtioHotPlug)
		c.Printer.Verbose("Property DiscVirtioHotPlug set: %v", discVirtioHotPlug)
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapi_v6.ArgDiscVirtioHotUnplug)) {
		discVirtioHotUnplug := viper.GetBool(core.GetFlagName(c.NS, cloudapi_v6.ArgDiscVirtioHotUnplug))
		proper.SetDiscVirtioHotUnplug(discVirtioHotUnplug)
		c.Printer.Verbose("Property DiscVirtioHotUnplug set: %v", discVirtioHotUnplug)
	}
	return &resources.Volume{
		Volume: ionoscloud.Volume{
			Properties: &proper.VolumeProperties,
		},
	}, nil
}

func getVolumeInfo(c *core.CommandConfig) (*resources.VolumeProperties, error) {
	input := resources.VolumeProperties{}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapi_v6.ArgName)) {
		name := viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgName))
		input.SetName(name)
		c.Printer.Verbose("Property Name set: %v", name)
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapi_v6.ArgBus)) {
		bus := viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgBus))
		input.SetBus(bus)
		c.Printer.Verbose("Property Bus set: %v", bus)
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapi_v6.ArgSize)) {
		size, err := utils.ConvertSize(
			viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgSize)),
			utils.GigaBytes,
		)
		if err != nil {
			return nil, err
		}
		input.SetSize(float32(size))
		c.Printer.Verbose("Property Size set: %vGB", float32(size))
	}
	// Check if flags are set and set options
	if viper.IsSet(core.GetFlagName(c.NS, cloudapi_v6.ArgCpuHotPlug)) {
		cpuHotPlug := viper.GetBool(core.GetFlagName(c.NS, cloudapi_v6.ArgCpuHotPlug))
		input.SetCpuHotPlug(cpuHotPlug)
		c.Printer.Verbose("Property CpuHotPlug set: %v", cpuHotPlug)
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapi_v6.ArgRamHotPlug)) {
		ramHotPlug := viper.GetBool(core.GetFlagName(c.NS, cloudapi_v6.ArgRamHotPlug))
		input.SetRamHotPlug(ramHotPlug)
		c.Printer.Verbose("Property RamHotPlug set: %v", ramHotPlug)
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapi_v6.ArgNicHotPlug)) {
		nicHotPlug := viper.GetBool(core.GetFlagName(c.NS, cloudapi_v6.ArgNicHotPlug))
		input.SetNicHotPlug(nicHotPlug)
		c.Printer.Verbose("Property NicHotPlug set: %v", nicHotPlug)
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapi_v6.ArgNicHotUnplug)) {
		nicHotUnplug := viper.GetBool(core.GetFlagName(c.NS, cloudapi_v6.ArgNicHotUnplug))
		input.SetNicHotUnplug(nicHotUnplug)
		c.Printer.Verbose("Property NicHotUnplug set: %v", nicHotUnplug)
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapi_v6.ArgDiscVirtioHotPlug)) {
		discVirtioHotPlug := viper.GetBool(core.GetFlagName(c.NS, cloudapi_v6.ArgDiscVirtioHotPlug))
		input.SetDiscVirtioHotPlug(discVirtioHotPlug)
		c.Printer.Verbose("Property DiscVirtioHotPlug set: %v", discVirtioHotPlug)
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapi_v6.ArgDiscVirtioHotUnplug)) {
		discVirtioHotUnplug := viper.GetBool(core.GetFlagName(c.NS, cloudapi_v6.ArgDiscVirtioHotUnplug))
		input.SetDiscVirtioHotUnplug(discVirtioHotUnplug)
		c.Printer.Verbose("Property DiscVirtioHotUnplug set: %v", discVirtioHotUnplug)
	}
	return &input, nil
}

// Server Volume Commands

func ServerVolumeCmd() *core.Command {
	ctx := context.TODO()
	serverVolumeCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "volume",
			Aliases:          []string{"v", "vol"},
			Short:            "Server Volume Operations",
			Long:             "The sub-commands of `ionosctl server volume` allow you to attach, get, list, detach Volumes from Servers.",
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
	attachVolume.AddStringSliceFlag(config.ArgCols, "", defaultVolumeCols, printer.ColsMessage(allVolumeCols))
	_ = attachVolume.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allVolumeCols, cobra.ShellCompDirectiveNoFileComp
	})
	attachVolume.AddStringFlag(cloudapi_v6.ArgDataCenterId, "", "", cloudapi_v6.DatacenterId, core.RequiredFlagOption())
	_ = attachVolume.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	attachVolume.AddStringFlag(cloudapi_v6.ArgVolumeId, cloudapi_v6.ArgIdShort, "", cloudapi_v6.VolumeId, core.RequiredFlagOption())
	_ = attachVolume.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgVolumeId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.VolumesIds(os.Stderr, viper.GetString(core.GetFlagName(attachVolume.NS, cloudapi_v6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	attachVolume.AddStringFlag(cloudapi_v6.ArgServerId, "", "", cloudapi_v6.ServerId, core.RequiredFlagOption())
	_ = attachVolume.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgServerId, func(cmd *cobra.Command, ags []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ServersIds(os.Stderr, viper.GetString(core.GetFlagName(attachVolume.NS, cloudapi_v6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
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
	listVolumes.AddStringSliceFlag(config.ArgCols, "", defaultVolumeCols, printer.ColsMessage(allVolumeCols))
	_ = listVolumes.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allVolumeCols, cobra.ShellCompDirectiveNoFileComp
	})
	listVolumes.AddStringFlag(cloudapi_v6.ArgDataCenterId, "", "", cloudapi_v6.DatacenterId, core.RequiredFlagOption())
	_ = listVolumes.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	listVolumes.AddStringFlag(cloudapi_v6.ArgServerId, "", "", cloudapi_v6.ServerId, core.RequiredFlagOption())
	_ = listVolumes.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgServerId, func(cmd *cobra.Command, ags []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ServersIds(os.Stderr, viper.GetString(core.GetFlagName(listVolumes.NS, cloudapi_v6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
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
	getVolumeCmd.AddStringSliceFlag(config.ArgCols, "", defaultVolumeCols, printer.ColsMessage(allVolumeCols))
	_ = getVolumeCmd.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allVolumeCols, cobra.ShellCompDirectiveNoFileComp
	})
	getVolumeCmd.AddStringFlag(cloudapi_v6.ArgDataCenterId, "", "", cloudapi_v6.DatacenterId, core.RequiredFlagOption())
	_ = getVolumeCmd.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	getVolumeCmd.AddStringFlag(cloudapi_v6.ArgServerId, "", "", cloudapi_v6.ServerId, core.RequiredFlagOption())
	_ = getVolumeCmd.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgServerId, func(cmd *cobra.Command, ags []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ServersIds(os.Stderr, viper.GetString(core.GetFlagName(getVolumeCmd.NS, cloudapi_v6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	getVolumeCmd.AddStringFlag(cloudapi_v6.ArgVolumeId, cloudapi_v6.ArgIdShort, "", cloudapi_v6.VolumeId, core.RequiredFlagOption())
	_ = getVolumeCmd.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgVolumeId, func(cmd *cobra.Command, ags []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.AttachedVolumesIds(os.Stderr, viper.GetString(core.GetFlagName(getVolumeCmd.NS, cloudapi_v6.ArgDataCenterId)), viper.GetString(core.GetFlagName(getVolumeCmd.NS, cloudapi_v6.ArgServerId))), cobra.ShellCompDirectiveNoFileComp
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
	detachVolume.AddStringSliceFlag(config.ArgCols, "", defaultVolumeCols, printer.ColsMessage(allVolumeCols))
	_ = detachVolume.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allVolumeCols, cobra.ShellCompDirectiveNoFileComp
	})
	detachVolume.AddStringFlag(cloudapi_v6.ArgDataCenterId, "", "", cloudapi_v6.DatacenterId, core.RequiredFlagOption())
	_ = detachVolume.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	detachVolume.AddStringFlag(cloudapi_v6.ArgVolumeId, cloudapi_v6.ArgIdShort, "", cloudapi_v6.VolumeId, core.RequiredFlagOption())
	_ = detachVolume.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgVolumeId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.AttachedVolumesIds(os.Stderr, viper.GetString(core.GetFlagName(detachVolume.NS, cloudapi_v6.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(detachVolume.NS, cloudapi_v6.ArgServerId))), cobra.ShellCompDirectiveNoFileComp
	})
	detachVolume.AddStringFlag(cloudapi_v6.ArgServerId, "", "", cloudapi_v6.ServerId, core.RequiredFlagOption())
	_ = detachVolume.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ServersIds(os.Stderr, viper.GetString(core.GetFlagName(detachVolume.NS, cloudapi_v6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	detachVolume.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for Volume detachment to be executed")
	detachVolume.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option for Request for Volume detachment [seconds]")

	return serverVolumeCmd
}

func PreRunDcServerVolumeIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapi_v6.ArgDataCenterId, cloudapi_v6.ArgServerId, cloudapi_v6.ArgVolumeId)
}

func RunServerVolumeAttach(c *core.CommandConfig) error {
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgDataCenterId))
	serverId := viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgServerId))
	volumeId := viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgVolumeId))
	c.Printer.Verbose("Datacenter ID: %v", dcId)
	c.Printer.Verbose("Attaching Volume with ID: %v to Server with ID: %v...", volumeId, serverId)
	attachedVol, resp, err := c.CloudApiV6Services.Servers().AttachVolume(dcId, serverId, volumeId)
	if err != nil {
		return err
	}

	if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetRequestPath(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getVolumePrint(resp, c, getVolume(attachedVol)))
}

func RunServerVolumesList(c *core.CommandConfig) error {
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgDataCenterId))
	serverId := viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgServerId))
	c.Printer.Verbose("Datacenter ID: %v", dcId)
	c.Printer.Verbose("Listing attached Volumes from Server with ID: %v...", serverId)
	attachedVols, _, err := c.CloudApiV6Services.Servers().ListVolumes(dcId, serverId)
	if err != nil {
		return err
	}
	return c.Printer.Print(getVolumePrint(nil, c, getAttachedVolumes(attachedVols)))
}

func RunServerVolumeGet(c *core.CommandConfig) error {
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgDataCenterId))
	serverId := viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgServerId))
	volumeId := viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgVolumeId))
	c.Printer.Verbose("Datacenter ID: %v", dcId)
	c.Printer.Verbose("Getting attached Volume with ID: %v from Server with ID: %v...", volumeId, serverId)
	attachedVol, _, err := c.CloudApiV6Services.Servers().GetVolume(dcId, serverId, volumeId)
	if err != nil {
		return err
	}
	return c.Printer.Print(getVolumePrint(nil, c, getVolume(attachedVol)))
}

func RunServerVolumeDetach(c *core.CommandConfig) error {
	if err := utils.AskForConfirm(c.Stdin, c.Printer, "detach volume from server"); err != nil {
		return err
	}
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgDataCenterId))
	serverId := viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgServerId))
	volumeId := viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgVolumeId))
	c.Printer.Verbose("Datacenter ID: %v", dcId)
	c.Printer.Verbose("Detaching Volume with ID: %v from Server with ID: %v...", volumeId, serverId)
	resp, err := c.CloudApiV6Services.Servers().DetachVolume(dcId, serverId, volumeId)
	if err != nil {
		return err
	}
	if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetRequestPath(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getVolumePrint(resp, c, nil))
}

// Output Printing

var (
	defaultVolumeCols = []string{"VolumeId", "Name", "Size", "Type", "LicenceType", "State", "Image"}
	allVolumeCols     = []string{"VolumeId", "Name", "Size", "Type", "LicenceType", "State", "Image", "Bus", "AvailabilityZone", "BackupunitId", "SshKeys",
		"DeviceNumber", "UserData"}
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
