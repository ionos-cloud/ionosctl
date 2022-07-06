package commands

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"

	"go.uber.org/multierr"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/commands/cloudapi-v6/completer"
	"github.com/ionos-cloud/ionosctl/commands/cloudapi-v6/query"
	"github.com/ionos-cloud/ionosctl/commands/cloudapi-v6/waiter"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/printer"
	"github.com/ionos-cloud/ionosctl/pkg/utils"
	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/services/cloudapi-v6"
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
		LongDesc:   "Use this command to list all Volumes from a Data Center on your account.\n\nYou can filter the results using `--filters` option. Use the following format to set filters: `--filters KEY1=VALUE1,KEY2=VALUE2`.\n" + completer.VolumesFiltersUsage() + "\n\nRequired values to run command:\n\n* Data Center Id",
		Example:    listVolumeExample,
		PreCmdRun:  PreRunVolumeList,
		CmdRun:     RunVolumeList,
		InitClient: true,
	})
	list.AddStringFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddIntFlag(cloudapiv6.ArgMaxResults, cloudapiv6.ArgMaxResultsShort, 0, "The maximum number of elements to return")
	list.AddIntFlag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, config.DefaultListDepth, "Controls the detail depth of the response objects. Max depth is 10.")
	list.AddStringFlag(cloudapiv6.ArgOrderBy, "", "", "Limits results to those containing a matching value for a specific property")
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgOrderBy, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.VolumesFilters(), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddStringSliceFlag(cloudapiv6.ArgFilters, cloudapiv6.ArgFiltersShort, []string{""}, "Limits results to those containing a matching value for a specific property. Use the following format to set filters: --filters KEY1=VALUE1,KEY2=VALUE2")
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgFilters, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.VolumesFilters(), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddBoolFlag(config.ArgNoHeaders, "", false, "When using text output, don't print headers")

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
	get.AddStringFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddStringFlag(cloudapiv6.ArgVolumeId, cloudapiv6.ArgIdShort, "", cloudapiv6.VolumeId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgVolumeId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.VolumesIds(os.Stderr, viper.GetString(core.GetFlagName(get.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddBoolFlag(config.ArgNoHeaders, "", false, "When using text output, don't print headers")
	get.AddIntFlag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, config.DefaultGetDepth, "Controls the detail depth of the response objects. Max depth is 10.")

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
		PreCmdRun:  PreRunVolumeCreate,
		CmdRun:     RunVolumeCreate,
		InitClient: true,
	})
	create.AddStringFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "Unnamed Volume", "Name of the Volume")
	create.AddStringFlag(cloudapiv6.ArgSize, cloudapiv6.ArgSizeShort, strconv.Itoa(cloudapiv6.DefaultVolumeSize), "The size of the Volume in GB. e.g.: --size 10 or --size 10GB. The maximum Volume size is determined by your contract limit")
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgSize, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"10GB", "20GB", "50GB", "100GB", "1TB"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(cloudapiv6.ArgBus, "", "VIRTIO", "The bus type of the Volume")
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgBus, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"VIRTIO", "IDE"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(cloudapiv6.ArgLicenceType, "", "LINUX", "Licence Type of the Volume")
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgLicenceType, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"LINUX", "WINDOWS", "WINDOWS2016", "UNKNOWN", "OTHER"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(cloudapiv6.ArgType, "", "HDD", "Type of the Volume")
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgLicenceType, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"HDD", "SSD", "SSD Standard", "SSD Premium"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(cloudapiv6.ArgAvailabilityZone, cloudapiv6.ArgAvailabilityZoneShort, "AUTO", "Availability zone of the Volume. Storage zone can only be selected prior provisioning")
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgAvailabilityZone, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"AUTO", "ZONE_1", "ZONE_2", "ZONE_3"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(cloudapiv6.ArgBackupUnitId, "", "", "The unique Id of the Backup Unit that User has access to. It is mandatory to provide either 'public image' or 'imageAlias' in conjunction with this property")
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgBackupUnitId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.BackupUnitsIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(cloudapiv6.ArgImageId, "", "", "The Image Id or Snapshot Id to be used as template for the new Volume. A password or SSH Key need to be set")
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgImageId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ImageIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(cloudapiv6.ArgImageAlias, cloudapiv6.ArgImageAliasShort, "", "The Image Alias to set instead of Image Id. A password or SSH Key need to be set")
	create.AddStringFlag(cloudapiv6.ArgPassword, cloudapiv6.ArgPasswordShort, "", "Initial password to be set for installed OS. Works with public Images only. Not modifiable. Password rules allows all characters from a-z, A-Z, 0-9")
	create.AddStringFlag(cloudapiv6.ArgUserData, "", "", "The cloud-init configuration for the Volume as base64 encoded string. It is mandatory to provide either 'public image' or 'imageAlias' that has cloud-init compatibility in conjunction with this property")
	create.AddBoolFlag(cloudapiv6.ArgCpuHotPlug, "", false, "It is capable of CPU hot plug (no reboot required). E.g.: --cpu-hot-plug=true, --cpu-hot-plug=false")
	create.AddBoolFlag(cloudapiv6.ArgRamHotPlug, "", false, "It is capable of memory hot plug (no reboot required). E.g.: --ram-hot-plug=true, --ram-hot-plug=false")
	create.AddBoolFlag(cloudapiv6.ArgNicHotPlug, "", false, "It is capable of nic hot plug (no reboot required). E.g.: --nic-hot-plug=true, --nic-hot-plug=false")
	create.AddBoolFlag(cloudapiv6.ArgNicHotUnplug, "", false, "It is capable of nic hot unplug (no reboot required). E.g.: --nic-hot-unplug=true, --nic-hot-unplug=false")
	create.AddBoolFlag(cloudapiv6.ArgDiscVirtioHotPlug, "", false, "It is capable of Virt-IO drive hot plug (no reboot required). E.g.: --disc-virtio-plug=true, --disc-virtio-plug=false")
	create.AddBoolFlag(cloudapiv6.ArgDiscVirtioHotUnplug, "", false, "It is capable of Virt-IO drive hot unplug (no reboot required). This works only for non-Windows virtual Machines. E.g.: --disc-virtio-unplug=true, --disc-virtio-unplug=false")
	create.AddStringFlag(cloudapiv6.ArgSshKeyPaths, cloudapiv6.ArgSshKeyPathsShort, "", "Absolute paths of the SSH Keys for the Volume")
	create.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for Volume creation to be executed")
	create.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option for Request for Volume creation [seconds]")
	create.AddIntFlag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, config.DefaultCreateDepth, "Controls the detail depth of the response objects. Max depth is 10.")

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
	update.AddStringFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(cloudapiv6.ArgVolumeId, cloudapiv6.ArgIdShort, "", cloudapiv6.VolumeId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgVolumeId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.VolumesIds(os.Stderr, viper.GetString(core.GetFlagName(update.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "", "Name of the Volume")
	update.AddStringFlag(cloudapiv6.ArgSize, "", "", "The size of the Volume in GB. e.g. 10 or 10GB. The maximum volume size is determined by your contract limit")
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgSize, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"10GB", "20GB", "50GB", "100GB", "1TB"}, cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(cloudapiv6.ArgBus, "", "VIRTIO", "Bus of the Volume")
	update.AddBoolFlag(cloudapiv6.ArgCpuHotPlug, "", false, "It is capable of CPU hot plug (no reboot required). E.g.: --cpu-hot-plug=true, --cpu-hot-plug=false")
	update.AddBoolFlag(cloudapiv6.ArgRamHotPlug, "", false, "It is capable of memory hot plug (no reboot required). E.g.: --ram-hot-plug=true, --ram-hot-plug=false")
	update.AddBoolFlag(cloudapiv6.ArgNicHotPlug, "", false, "It is capable of nic hot plug (no reboot required). E.g.: --nic-hot-plug=true, --nic-hot-plug=false")
	update.AddBoolFlag(cloudapiv6.ArgNicHotUnplug, "", false, "It is capable of nic hot unplug (no reboot required). E.g.: --nic-hot-unplug=true, --nic-hot-unplug=false")
	update.AddBoolFlag(cloudapiv6.ArgDiscVirtioHotPlug, "", false, "It is capable of Virt-IO drive hot plug (no reboot required). E.g.: --disc-virtio-plug=true, --disc-virtio-plug=false")
	update.AddBoolFlag(cloudapiv6.ArgDiscVirtioHotUnplug, "", false, "It is capable of Virt-IO drive hot unplug (no reboot required). This works only for non-Windows virtual Machines. E.g.: --disc-virtio-unplug=true, --disc-virtio-unplug=false")
	update.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for Volume update to be executed")
	update.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option for Request for Volume update [seconds]")
	update.AddIntFlag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, config.DefaultUpdateDepth, "Controls the detail depth of the response objects. Max depth is 10.")

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
		PreCmdRun:  PreRunDcVolumeDelete,
		CmdRun:     RunVolumeDelete,
		InitClient: true,
	})
	deleteCmd.AddStringFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddStringFlag(cloudapiv6.ArgVolumeId, cloudapiv6.ArgIdShort, "", cloudapiv6.VolumeId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgVolumeId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.VolumesIds(os.Stderr, viper.GetString(core.GetFlagName(deleteCmd.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for Volume deletion to be executed")
	deleteCmd.AddBoolFlag(cloudapiv6.ArgAll, cloudapiv6.ArgAllShort, false, "Delete all Volumes from a virtual Datacenter.")
	deleteCmd.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option for Request for Volume deletion [seconds]")
	deleteCmd.AddIntFlag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, config.DefaultDeleteDepth, "Controls the detail depth of the response objects. Max depth is 10.")

	return volumeCmd
}

func PreRunVolumeList(c *core.PreCommandConfig) error {
	if err := core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgDataCenterId); err != nil {
		return err
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgFilters)) {
		return query.ValidateFilters(c, completer.VolumesFilters(), completer.VolumesFiltersUsage())
	}
	return nil
}

func PreRunDcVolumeIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgDataCenterId, cloudapiv6.ArgVolumeId)
}

func PreRunDcVolumeDelete(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{cloudapiv6.ArgDataCenterId, cloudapiv6.ArgVolumeId},
		[]string{cloudapiv6.ArgDataCenterId, cloudapiv6.ArgAll},
	)
}

func PreRunVolumeCreate(c *core.PreCommandConfig) error {
	err := core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgDataCenterId)
	if err != nil {
		return err
	}
	// Validate flags
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgImageId)) || viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgImageAlias)) {
		return core.CheckRequiredFlagsSets(c.Command, c.NS,
			[]string{cloudapiv6.ArgDataCenterId, cloudapiv6.ArgImageId, cloudapiv6.ArgPassword},
			[]string{cloudapiv6.ArgDataCenterId, cloudapiv6.ArgImageId, cloudapiv6.ArgSshKeyPaths},
			[]string{cloudapiv6.ArgDataCenterId, cloudapiv6.ArgImageAlias, cloudapiv6.ArgPassword},
			[]string{cloudapiv6.ArgDataCenterId, cloudapiv6.ArgImageAlias, cloudapiv6.ArgSshKeyPaths})
	}
	return nil
}

func RunVolumeList(c *core.CommandConfig) error {
	c.Printer.Verbose("Listing Volumes from Datacenter with ID: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)))
	// Add Query Parameters for GET Requests
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	if !structs.IsZero(listQueryParams) {
		if listQueryParams.Filters != nil {
			filters := *listQueryParams.Filters
			if val, ok := filters["size"]; ok {
				convertedSize, err := utils.ConvertSize(val, utils.GigaBytes)
				if err != nil {
					return err
				}
				filters["size"] = strconv.Itoa(convertedSize)
				listQueryParams.Filters = &filters
			}
		}
		c.Printer.Verbose("Query Parameters set: %v", utils.GetPropertiesKVSet(listQueryParams))
	}
	volumes, resp, err := c.CloudApiV6Services.Volumes().List(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)), listQueryParams)
	if resp != nil {
		c.Printer.Verbose(config.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getVolumePrint(nil, c, getVolumes(volumes)))
}

func RunVolumeGet(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	queryParams := listQueryParams.QueryParams
	if !structs.IsZero(queryParams) {
		c.Printer.Verbose("Query Parameters set: %v", utils.GetPropertiesKVSet(queryParams))
	}
	c.Printer.Verbose("Datacenter ID: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)))
	c.Printer.Verbose("Volume with id: %v is getting...", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgVolumeId)))
	vol, resp, err := c.CloudApiV6Services.Volumes().Get(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgVolumeId)),
	)
	if resp != nil {
		c.Printer.Verbose(config.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getVolumePrint(nil, c, getVolume(vol)))
}

func RunVolumeCreate(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	queryParams := listQueryParams.QueryParams
	if !structs.IsZero(queryParams) {
		c.Printer.Verbose("Query Parameters set: %v", utils.GetPropertiesKVSet(queryParams))
	}
	input, err := getNewVolume(c)
	if err != nil {
		return err
	}
	vol, resp, err := c.CloudApiV6Services.Volumes().Create(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)), *input)
	if resp != nil && printer.GetId(resp) != "" {
		c.Printer.Verbose(config.RequestInfoMessage, printer.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}
	if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getVolumePrint(resp, c, getVolume(vol)))
}

func RunVolumeUpdate(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	queryParams := listQueryParams.QueryParams
	if !structs.IsZero(queryParams) {
		c.Printer.Verbose("Query Parameters set: %v", utils.GetPropertiesKVSet(queryParams))
	}
	input, err := getVolumeInfo(c)
	if err != nil {
		return err
	}
	vol, resp, err := c.CloudApiV6Services.Volumes().Update(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgVolumeId)), *input)
	if resp != nil && printer.GetId(resp) != "" {
		c.Printer.Verbose(config.RequestInfoMessage, printer.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}
	if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getVolumePrint(resp, c, getVolume(vol)))
}

func RunVolumeDelete(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	queryParams := listQueryParams.QueryParams
	if !structs.IsZero(queryParams) {
		c.Printer.Verbose("Query Parameters set: %v", utils.GetPropertiesKVSet(queryParams))
	}
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))
	volumeId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgVolumeId))
	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAll)) {
		if err := DeleteAllVolumes(c); err != nil {
			return err
		}
		return c.Printer.Print(printer.Result{Resource: c.Resource, Verb: c.Verb})
	} else {
		if err := utils.AskForConfirm(c.Stdin, c.Printer, "delete volume"); err != nil {
			return err
		}
		c.Printer.Verbose("Starting deleting Volume with id: %v...", volumeId)
		resp, err := c.CloudApiV6Services.Volumes().Delete(dcId, volumeId)
		if resp != nil && printer.GetId(resp) != "" {
			c.Printer.Verbose(config.RequestInfoMessage, printer.GetId(resp), resp.RequestTime)
		}
		if err != nil {
			return err
		}
		if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
			return err
		}
		return c.Printer.Print(getVolumePrint(resp, c, nil))
	}
}

func getNewVolume(c *core.CommandConfig) (*resources.Volume, error) {
	proper := resources.VolumeProperties{}
	name := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName))
	bus := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgBus))
	volumeType := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgType))
	availabilityZone := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgAvailabilityZone))
	// It will get the default values, if flags not set
	proper.SetName(name)
	proper.SetBus(bus)
	proper.SetType(volumeType)
	proper.SetAvailabilityZone(availabilityZone)
	c.Printer.Verbose("Properties set for creating the Volume: Name: %v, Bus: %v, VolumeType: %v, AvailabilityZone: %v",
		name, bus, volumeType, availabilityZone)
	size, err := utils.ConvertSize(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgSize)),
		utils.GigaBytes,
	)
	if err != nil {
		return nil, err
	}
	proper.SetSize(float32(size))
	c.Printer.Verbose("Property Size set: %vGB", float32(size))
	// Check if flags are set and set options
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgBackupUnitId)) {
		backupUnitId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgBackupUnitId))
		proper.SetBackupunitId(backupUnitId)
		c.Printer.Verbose("Property BackupUnitId set: %v", backupUnitId)
	}
	if (!viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgImageId)) &&
		!viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgImageAlias))) ||
		viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgLicenceType)) {
		licenceType := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLicenceType))
		proper.SetLicenceType(licenceType)
		c.Printer.Verbose("Property LicenceType set: %v", licenceType)
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgImageId)) {
		imageId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgImageId))
		proper.SetImage(imageId)
		c.Printer.Verbose("Property Image set: %v", imageId)
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgImageAlias)) {
		imageAlias := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgImageAlias))
		proper.SetImageAlias(imageAlias)
		c.Printer.Verbose("Property ImageAlias set: %v", imageAlias)
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgPassword)) {
		imagePassword := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgPassword))
		proper.SetImagePassword(imagePassword)
		c.Printer.Verbose("Property ImagePassword set")
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgSshKeyPaths)) {
		sshKeysPaths := viper.GetStringSlice(core.GetFlagName(c.NS, cloudapiv6.ArgSshKeyPaths))
		c.Printer.Verbose("SSH Key Paths: %v", sshKeysPaths)
		sshKeys, err := getSshKeysFromPaths(sshKeysPaths)
		if err != nil {
			return nil, err
		}
		proper.SetSshKeys(sshKeys)
		c.Printer.Verbose("Property SshKeys set")
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgUserData)) {
		userData := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgUserData))
		proper.SetUserData(userData)
		c.Printer.Verbose("Property UserData set: %v", userData)
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgCpuHotPlug)) {
		cpuHotPlug := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgCpuHotPlug))
		proper.SetCpuHotPlug(cpuHotPlug)
		c.Printer.Verbose("Property CpuHotPlug set: %v", cpuHotPlug)
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgRamHotPlug)) {
		ramHotPlug := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgRamHotPlug))
		proper.SetRamHotPlug(ramHotPlug)
		c.Printer.Verbose("Property RamHotPlug set: %v", ramHotPlug)
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgNicHotPlug)) {
		nicHotPlug := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgNicHotPlug))
		proper.SetNicHotPlug(nicHotPlug)
		c.Printer.Verbose("Property NicHotPlug set: %v", nicHotPlug)
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgNicHotUnplug)) {
		nicHotUnplug := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgNicHotUnplug))
		proper.SetNicHotUnplug(nicHotUnplug)
		c.Printer.Verbose("Property NicHotUnplug set: %v", nicHotUnplug)
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgDiscVirtioHotPlug)) {
		discVirtioHotPlug := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgDiscVirtioHotPlug))
		proper.SetDiscVirtioHotPlug(discVirtioHotPlug)
		c.Printer.Verbose("Property DiscVirtioHotPlug set: %v", discVirtioHotPlug)
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgDiscVirtioHotUnplug)) {
		discVirtioHotUnplug := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgDiscVirtioHotUnplug))
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
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgName)) {
		name := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName))
		input.SetName(name)
		c.Printer.Verbose("Property Name set: %v", name)
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgBus)) {
		bus := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgBus))
		input.SetBus(bus)
		c.Printer.Verbose("Property Bus set: %v", bus)
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgSize)) {
		size, err := utils.ConvertSize(
			viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgSize)),
			utils.GigaBytes,
		)
		if err != nil {
			return nil, err
		}
		input.SetSize(float32(size))
		c.Printer.Verbose("Property Size set: %vGB", float32(size))
	}
	// Check if flags are set and set options
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgCpuHotPlug)) {
		cpuHotPlug := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgCpuHotPlug))
		input.SetCpuHotPlug(cpuHotPlug)
		c.Printer.Verbose("Property CpuHotPlug set: %v", cpuHotPlug)
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgRamHotPlug)) {
		ramHotPlug := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgRamHotPlug))
		input.SetRamHotPlug(ramHotPlug)
		c.Printer.Verbose("Property RamHotPlug set: %v", ramHotPlug)
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgNicHotPlug)) {
		nicHotPlug := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgNicHotPlug))
		input.SetNicHotPlug(nicHotPlug)
		c.Printer.Verbose("Property NicHotPlug set: %v", nicHotPlug)
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgNicHotUnplug)) {
		nicHotUnplug := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgNicHotUnplug))
		input.SetNicHotUnplug(nicHotUnplug)
		c.Printer.Verbose("Property NicHotUnplug set: %v", nicHotUnplug)
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgDiscVirtioHotPlug)) {
		discVirtioHotPlug := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgDiscVirtioHotPlug))
		input.SetDiscVirtioHotPlug(discVirtioHotPlug)
		c.Printer.Verbose("Property DiscVirtioHotPlug set: %v", discVirtioHotPlug)
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgDiscVirtioHotUnplug)) {
		discVirtioHotUnplug := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgDiscVirtioHotUnplug))
		input.SetDiscVirtioHotUnplug(discVirtioHotUnplug)
		c.Printer.Verbose("Property DiscVirtioHotUnplug set: %v", discVirtioHotUnplug)
	}
	return &input, nil
}

func DeleteAllVolumes(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	queryParams := listQueryParams.QueryParams
	if !structs.IsZero(queryParams) {
		c.Printer.Verbose("Query Parameters set: %v", utils.GetPropertiesKVSet(queryParams))
	}
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))
	c.Printer.Verbose("Datacenter ID: %v", dcId)
	c.Printer.Verbose("Getting Volumes...")
	volumes, resp, err := c.CloudApiV6Services.Volumes().List(dcId, resources.ListQueryParams{})
	if err != nil {
		return err
	}
	if volumesItems, ok := volumes.GetItemsOk(); ok && volumesItems != nil {
		if len(*volumesItems) > 0 {
			_ = c.Printer.Print("Volumes to be deleted:")
			for _, volume := range *volumesItems {
				if id, ok := volume.GetIdOk(); ok && id != nil {
					_ = c.Printer.Print("Volume Id: " + *id)
				}
				if properties, ok := volume.GetPropertiesOk(); ok && properties != nil {
					if name, ok := properties.GetNameOk(); ok && name != nil {
						_ = c.Printer.Print("Volume Name: " + *name)
					}
				}
			}
			if err = utils.AskForConfirm(c.Stdin, c.Printer, "delete all the Volumes"); err != nil {
				return err
			}
			c.Printer.Verbose("Deleting all the Volumes...")
			var multiErr error
			for _, volume := range *volumesItems {
				if id, ok := volume.GetIdOk(); ok && id != nil {
					c.Printer.Verbose("Starting deleting Volume with id: %v is...", *id)
					resp, err = c.CloudApiV6Services.Volumes().Delete(dcId, *id)
					if resp != nil && printer.GetId(resp) != "" {
						c.Printer.Verbose(config.RequestInfoMessage, printer.GetId(resp), resp.RequestTime)
					}
					if err != nil {
						multiErr = multierr.Append(multiErr, fmt.Errorf(config.DeleteAllAppendErr, c.Resource, *id, err))
						continue
					} else {
						_ = c.Printer.Print(fmt.Sprintf(config.StatusDeletingAll, c.Resource, *id))
					}
					if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
						multiErr = multierr.Append(multiErr, fmt.Errorf(config.WaitDeleteAllAppendErr, c.Resource, *id, err))
						continue
					}
				}
			}
			if multiErr != nil {
				return multiErr
			}
			return nil
		} else {
			return errors.New("no Volumes found")
		}
	} else {
		return errors.New("could not get items of Volumes")
	}
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
	attachVolume.AddStringFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = attachVolume.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	attachVolume.AddStringFlag(cloudapiv6.ArgVolumeId, cloudapiv6.ArgIdShort, "", cloudapiv6.VolumeId, core.RequiredFlagOption())
	_ = attachVolume.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgVolumeId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.VolumesIds(os.Stderr, viper.GetString(core.GetFlagName(attachVolume.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	attachVolume.AddStringFlag(cloudapiv6.ArgServerId, "", "", cloudapiv6.ServerId, core.RequiredFlagOption())
	_ = attachVolume.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgServerId, func(cmd *cobra.Command, ags []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ServersIds(os.Stderr, viper.GetString(core.GetFlagName(attachVolume.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
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
		LongDesc:   "Use this command to retrieve a list of Volumes attached to the Server.\n\nYou can filter the results using `--filters` option. Use the following format to set filters: `--filters KEY1=VALUE1,KEY2=VALUE2`.\n" + completer.VolumesFiltersUsage() + "\n\nRequired values to run command:\n\n* Data Center Id\n* Server Id",
		Example:    listVolumesServerExample,
		PreCmdRun:  PreRunServerVolumeList,
		CmdRun:     RunServerVolumesList,
		InitClient: true,
	})
	listVolumes.AddStringSliceFlag(config.ArgCols, "", defaultVolumeCols, printer.ColsMessage(allVolumeCols))
	_ = listVolumes.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allVolumeCols, cobra.ShellCompDirectiveNoFileComp
	})
	listVolumes.AddStringFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = listVolumes.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	listVolumes.AddStringFlag(cloudapiv6.ArgServerId, "", "", cloudapiv6.ServerId, core.RequiredFlagOption())
	_ = listVolumes.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgServerId, func(cmd *cobra.Command, ags []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ServersIds(os.Stderr, viper.GetString(core.GetFlagName(listVolumes.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	listVolumes.AddIntFlag(cloudapiv6.ArgMaxResults, cloudapiv6.ArgMaxResultsShort, 0, "The maximum number of elements to return")
	listVolumes.AddStringFlag(cloudapiv6.ArgOrderBy, "", "", "Limits results to those containing a matching value for a specific property")
	_ = listVolumes.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgOrderBy, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.VolumesFilters(), cobra.ShellCompDirectiveNoFileComp
	})
	listVolumes.AddStringSliceFlag(cloudapiv6.ArgFilters, cloudapiv6.ArgFiltersShort, []string{""}, "Limits results to those containing a matching value for a specific property. Use the following format to set filters: --filters KEY1=VALUE1,KEY2=VALUE2")
	_ = listVolumes.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgFilters, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.VolumesFilters(), cobra.ShellCompDirectiveNoFileComp
	})
	listVolumes.AddBoolFlag(config.ArgNoHeaders, "", false, "When using text output, don't print headers")

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
	getVolumeCmd.AddStringFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = getVolumeCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	getVolumeCmd.AddStringFlag(cloudapiv6.ArgServerId, "", "", cloudapiv6.ServerId, core.RequiredFlagOption())
	_ = getVolumeCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgServerId, func(cmd *cobra.Command, ags []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ServersIds(os.Stderr, viper.GetString(core.GetFlagName(getVolumeCmd.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	getVolumeCmd.AddStringFlag(cloudapiv6.ArgVolumeId, cloudapiv6.ArgIdShort, "", cloudapiv6.VolumeId, core.RequiredFlagOption())
	_ = getVolumeCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgVolumeId, func(cmd *cobra.Command, ags []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.AttachedVolumesIds(os.Stderr, viper.GetString(core.GetFlagName(getVolumeCmd.NS, cloudapiv6.ArgDataCenterId)), viper.GetString(core.GetFlagName(getVolumeCmd.NS, cloudapiv6.ArgServerId))), cobra.ShellCompDirectiveNoFileComp
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
		PreCmdRun:  PreRunDcServerVolumeDetach,
		CmdRun:     RunServerVolumeDetach,
		InitClient: true,
	})
	detachVolume.AddStringSliceFlag(config.ArgCols, "", defaultVolumeCols, printer.ColsMessage(allVolumeCols))
	_ = detachVolume.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allVolumeCols, cobra.ShellCompDirectiveNoFileComp
	})
	detachVolume.AddStringFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = detachVolume.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	detachVolume.AddStringFlag(cloudapiv6.ArgVolumeId, cloudapiv6.ArgIdShort, "", cloudapiv6.VolumeId, core.RequiredFlagOption())
	_ = detachVolume.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgVolumeId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.AttachedVolumesIds(os.Stderr, viper.GetString(core.GetFlagName(detachVolume.NS, cloudapiv6.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(detachVolume.NS, cloudapiv6.ArgServerId))), cobra.ShellCompDirectiveNoFileComp
	})
	detachVolume.AddStringFlag(cloudapiv6.ArgServerId, "", "", cloudapiv6.ServerId, core.RequiredFlagOption())
	_ = detachVolume.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ServersIds(os.Stderr, viper.GetString(core.GetFlagName(detachVolume.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	detachVolume.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for Volume detachment to be executed")
	detachVolume.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option for Request for Volume detachment [seconds]")
	detachVolume.AddBoolFlag(cloudapiv6.ArgAll, cloudapiv6.ArgAllShort, false, "Detach all Volumes.")

	return serverVolumeCmd
}

func PreRunServerVolumeList(c *core.PreCommandConfig) error {
	if err := core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgDataCenterId, cloudapiv6.ArgServerId); err != nil {
		return err
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgFilters)) {
		return query.ValidateFilters(c, completer.VolumesFilters(), completer.VolumesFiltersUsage())
	}
	return nil
}

func PreRunDcServerVolumeIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgDataCenterId, cloudapiv6.ArgServerId, cloudapiv6.ArgVolumeId)
}

func PreRunDcServerVolumeDetach(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{cloudapiv6.ArgDataCenterId, cloudapiv6.ArgServerId, cloudapiv6.ArgVolumeId},
		[]string{cloudapiv6.ArgDataCenterId, cloudapiv6.ArgServerId, cloudapiv6.ArgAll},
	)
}

func RunServerVolumeAttach(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	queryParams := listQueryParams.QueryParams
	if !structs.IsZero(queryParams) {
		c.Printer.Verbose("Query Parameters set: %v", utils.GetPropertiesKVSet(queryParams))
	}
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))
	serverId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgServerId))
	volumeId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgVolumeId))
	c.Printer.Verbose("Datacenter ID: %v", dcId)
	c.Printer.Verbose("Attaching Volume with ID: %v to Server with ID: %v...", volumeId, serverId)
	attachedVol, resp, err := c.CloudApiV6Services.Servers().AttachVolume(dcId, serverId, volumeId)
	if resp != nil && printer.GetId(resp) != "" {
		c.Printer.Verbose(config.RequestInfoMessage, printer.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}

	if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getVolumePrint(resp, c, getVolume(attachedVol)))
}

func RunServerVolumesList(c *core.CommandConfig) error {
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))
	serverId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgServerId))
	c.Printer.Verbose("Datacenter ID: %v", dcId)
	c.Printer.Verbose("Listing attached Volumes from Server with ID: %v...", serverId)
	// Add Query Parameters for GET Requests
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	if !structs.IsZero(listQueryParams) {
		if listQueryParams.Filters != nil {
			filters := *listQueryParams.Filters
			if val, ok := filters["size"]; ok {
				convertedSize, err := utils.ConvertSize(val, utils.GigaBytes)
				if err != nil {
					return err
				}
				filters["size"] = strconv.Itoa(convertedSize)
				listQueryParams.Filters = &filters
			}
		}
		c.Printer.Verbose("Query Parameters set: %v", utils.GetPropertiesKVSet(listQueryParams))
	}
	attachedVols, resp, err := c.CloudApiV6Services.Servers().ListVolumes(dcId, serverId, listQueryParams)
	if resp != nil {
		c.Printer.Verbose(config.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getVolumePrint(nil, c, getAttachedVolumes(attachedVols)))
}

func RunServerVolumeGet(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	queryParams := listQueryParams.QueryParams
	if !structs.IsZero(queryParams) {
		c.Printer.Verbose("Query Parameters set: %v", utils.GetPropertiesKVSet(queryParams))
	}
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))
	serverId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgServerId))
	volumeId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgVolumeId))
	c.Printer.Verbose("Datacenter ID: %v", dcId)
	c.Printer.Verbose("Getting attached Volume with ID: %v from Server with ID: %v...", volumeId, serverId)
	attachedVol, resp, err := c.CloudApiV6Services.Servers().GetVolume(dcId, serverId, volumeId)
	if resp != nil {
		c.Printer.Verbose(config.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getVolumePrint(nil, c, getVolume(attachedVol)))
}

func RunServerVolumeDetach(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	queryParams := listQueryParams.QueryParams
	if !structs.IsZero(queryParams) {
		c.Printer.Verbose("Query Parameters set: %v", utils.GetPropertiesKVSet(queryParams))
	}
	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAll)) {
		if err := DetachAllServerVolumes(c); err != nil {
			return err
		}
		return c.Printer.Print(printer.Result{Resource: c.Resource, Verb: c.Verb})
	} else {
		if err := utils.AskForConfirm(c.Stdin, c.Printer, "detach volume from server"); err != nil {
			return err
		}
		dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))
		serverId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgServerId))
		volumeId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgVolumeId))
		c.Printer.Verbose("Datacenter ID: %v", dcId)
		c.Printer.Verbose("Detaching Volume with ID: %v from Server with ID: %v...", volumeId, serverId)
		resp, err := c.CloudApiV6Services.Servers().DetachVolume(dcId, serverId, volumeId)
		if resp != nil && printer.GetId(resp) != "" {
			c.Printer.Verbose(config.RequestInfoMessage, printer.GetId(resp), resp.RequestTime)
		}
		if err != nil {
			return err
		}
		if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
			return err
		}
		return c.Printer.Print(getVolumePrint(resp, c, nil))
	}
}

func DetachAllServerVolumes(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	queryParams := listQueryParams.QueryParams
	if !structs.IsZero(queryParams) {
		c.Printer.Verbose("Query Parameters set: %v", utils.GetPropertiesKVSet(queryParams))
	}
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))
	serverId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgServerId))
	c.Printer.Verbose("Datacenter ID: %v", dcId)
	c.Printer.Verbose("Server ID: %v", serverId)
	c.Printer.Verbose("Getting Volumes...")
	volumes, resp, err := c.CloudApiV6Services.Servers().ListVolumes(dcId, serverId, resources.ListQueryParams{})
	if err != nil {
		return err
	}
	if volumesItems, ok := volumes.GetItemsOk(); ok && volumesItems != nil {
		if len(*volumesItems) > 0 {
			_ = c.Printer.Print("Volumes to be detached:")
			for _, volume := range *volumesItems {
				toPrint := ""
				if id, ok := volume.GetIdOk(); ok && id != nil {
					toPrint += "Volume Id: " + *id
				}
				if properties, ok := volume.GetPropertiesOk(); ok && properties != nil {
					if name, ok := properties.GetNameOk(); ok && name != nil {
						toPrint += " Volume Name: " + *name
					}
				}
				_ = c.Printer.Print(toPrint)
			}
			if err := utils.AskForConfirm(c.Stdin, c.Printer, "detach all the Volumes"); err != nil {
				return err
			}
			c.Printer.Verbose("Detaching all the Volumes...")
			var multiErr error
			for _, volume := range *volumesItems {
				if id, ok := volume.GetIdOk(); ok && id != nil {
					c.Printer.Verbose("Starting detaching Volume with id: %v...", *id)
					resp, err = c.CloudApiV6Services.Servers().DetachVolume(dcId, serverId, *id)
					if resp != nil && printer.GetId(resp) != "" {
						c.Printer.Verbose(config.RequestInfoMessage, printer.GetId(resp), resp.RequestTime)
					}
					if err != nil {
						multiErr = multierr.Append(multiErr, fmt.Errorf(config.DeleteAllAppendErr, c.Resource, *id, err))
						continue
					} else {
						_ = c.Printer.Print(fmt.Sprintf(config.StatusRemovingAll, c.Resource, *id))
					}
					if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
						multiErr = multierr.Append(multiErr, fmt.Errorf(config.WaitDeleteAllAppendErr, c.Resource, *id, err))
						continue
					}
				}
			}
			if multiErr != nil {
				return multiErr
			}
			return nil
		} else {
			return errors.New("no Volumes found")
		}
	} else {
		return errors.New("could not get items of Volumes")
	}
}

// Output Printing

var (
	defaultVolumeCols = []string{"VolumeId", "Name", "Size", "Type", "LicenceType", "State", "Image"}
	allVolumeCols     = []string{"VolumeId", "Name", "Size", "Type", "LicenceType", "State", "Image", "Bus", "AvailabilityZone", "BackupunitId",
		"DeviceNumber", "UserData", "BootServerId"}
)

type VolumePrint struct {
	VolumeId         string `json:"VolumeId,omitempty"`
	Name             string `json:"Name,omitempty"`
	Size             string `json:"Size,omitempty"`
	Type             string `json:"Type,omitempty"`
	LicenceType      string `json:"LicenceType,omitempty"`
	Bus              string `json:"Bus,omitempty"`
	AvailabilityZone string `json:"AvailabilityZone,omitempty"`
	State            string `json:"State,omitempty"`
	Image            string `json:"Image,omitempty"`
	DeviceNumber     int64  `json:"DeviceNumber,omitempty"`
	BackupunitId     string `json:"BackupunitId,omitempty"`
	UserData         string `json:"UserData,omitempty"`
	BootServerId     string `json:"BootServerId,omitempty"`
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
		"DeviceNumber":     "DeviceNumber",
		"BackupunitId":     "BackupunitId",
		"UserData":         "UserData",
		"BootServerId":     "BootServerId",
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
	volumeObjs := make([]resources.Volume, 0)
	if items, ok := volumes.GetItemsOk(); ok && items != nil {
		for _, volume := range *items {
			volumeObjs = append(volumeObjs, resources.Volume{Volume: volume})
		}
	}
	return volumeObjs
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
		var volumePrint VolumePrint
		if propertiesOk, ok := v.GetPropertiesOk(); ok && propertiesOk != nil {
			if idOk, ok := v.GetIdOk(); ok && idOk != nil {
				volumePrint.VolumeId = *idOk
			}
			if nameOk, ok := propertiesOk.GetNameOk(); ok && nameOk != nil {
				volumePrint.Name = *nameOk
			}
			if licenceTypeOk, ok := propertiesOk.GetLicenceTypeOk(); ok && licenceTypeOk != nil {
				volumePrint.LicenceType = *licenceTypeOk
			}
			if sizeOk, ok := propertiesOk.GetSizeOk(); ok && sizeOk != nil {
				volumePrint.Size = fmt.Sprintf("%vGB", *sizeOk)
			}
			if busOk, ok := propertiesOk.GetBusOk(); ok && busOk != nil {
				volumePrint.Bus = *busOk
			}
			if typeOk, ok := propertiesOk.GetTypeOk(); ok && typeOk != nil {
				volumePrint.Type = *typeOk
			}
			if availabilityZoneOk, ok := propertiesOk.GetAvailabilityZoneOk(); ok && availabilityZoneOk != nil {
				volumePrint.AvailabilityZone = *availabilityZoneOk
			}
			if backupunitIdOk, ok := propertiesOk.GetBackupunitIdOk(); ok && backupunitIdOk != nil {
				volumePrint.BackupunitId = *backupunitIdOk
			}
			if imageOk, ok := propertiesOk.GetImageOk(); ok && imageOk != nil {
				volumePrint.Image = *imageOk
			}
			if userDataOk, ok := propertiesOk.GetUserDataOk(); ok && userDataOk != nil {
				volumePrint.UserData = *userDataOk
			}
			if deviceNumberOk, ok := propertiesOk.GetDeviceNumberOk(); ok && deviceNumberOk != nil {
				volumePrint.DeviceNumber = *deviceNumberOk
			}
			if bootServerOk, ok := propertiesOk.GetBootServerOk(); ok && bootServerOk != nil {
				volumePrint.BootServerId = *bootServerOk
			}
		}
		if metadataOk, ok := v.GetMetadataOk(); ok && metadataOk != nil {
			if stateOk, ok := metadataOk.GetStateOk(); ok && stateOk != nil {
				volumePrint.State = *stateOk
			}
		}
		o := structs.Map(volumePrint)
		out = append(out, o)
	}
	return out
}

func getSshKeysFromPaths(paths []string) ([]string, error) {
	sshKeys := make([]string, 0)
	if len(paths) != 0 {
		for _, sshKeyPath := range paths {
			publicKey, err := utils.ReadPublicKey(sshKeyPath)
			if err != nil {
				return sshKeys, err
			}
			sshKeys = append(sshKeys, publicKey)
		}
	}
	return sshKeys, nil
}
