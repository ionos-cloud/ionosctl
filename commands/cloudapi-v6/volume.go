package commands

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/completer"
	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/query"
	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/waiter"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/confirm"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/json2table"
	"github.com/ionos-cloud/ionosctl/v6/pkg/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/pkg/tabheaders"
	"github.com/ionos-cloud/ionosctl/v6/pkg/utils"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	allVolumeJSONPaths = map[string]string{
		"VolumeId":         "id",
		"Name":             "properties.name",
		"Size":             "properties.size",
		"Type":             "properties.type",
		"LicenceType":      "properties.licenceType",
		"Bus":              "properties.bus",
		"AvailabilityZone": "properties.availabilityZone",
		"State":            "metadata.state",
		"Image":            "properties.image",
		"DeviceNumber":     "properties.deviceNumber",
		"BackupunitId":     "properties.backupunitId",
		"UserData":         "properties.userData",
		"BootServerId":     "properties.bootServer",
	}

	defaultVolumeCols = []string{"VolumeId", "Name", "Size", "Type", "LicenceType", "State", "Image"}
	allVolumeCols     = []string{"VolumeId", "Name", "Size", "Type", "LicenceType", "State", "Image", "Bus", "AvailabilityZone", "BackupunitId",
		"DeviceNumber", "UserData", "BootServerId", "DatacenterId"}
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
	globalFlags.StringSliceP(constants.ArgCols, "", defaultVolumeCols, tabheaders.ColsMessage(allVolumeCols))
	_ = viper.BindPFlag(core.GetFlagName(volumeCmd.Name(), constants.ArgCols), globalFlags.Lookup(constants.ArgCols))
	_ = volumeCmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
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
	list.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddInt32Flag(constants.FlagMaxResults, constants.FlagMaxResultsShort, cloudapiv6.DefaultMaxResults, constants.DescMaxResults)
	list.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultListDepth, cloudapiv6.ArgDepthDescription)
	list.AddStringFlag(cloudapiv6.ArgOrderBy, "", "", cloudapiv6.ArgOrderByDescription)
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgOrderBy, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.VolumesFilters(), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddStringSliceFlag(cloudapiv6.ArgFilters, cloudapiv6.ArgFiltersShort, []string{""}, cloudapiv6.ArgFiltersDescription)
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgFilters, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.VolumesFilters(), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddBoolFlag(constants.ArgNoHeaders, "", false, cloudapiv6.ArgNoHeadersDescription)
	list.AddBoolFlag(cloudapiv6.ArgAll, cloudapiv6.ArgAllShort, false, cloudapiv6.ArgListAllDescription)

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
	get.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddUUIDFlag(cloudapiv6.ArgVolumeId, cloudapiv6.ArgIdShort, "", cloudapiv6.VolumeId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgVolumeId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.VolumesIds(viper.GetString(core.GetFlagName(get.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddBoolFlag(constants.ArgNoHeaders, "", false, cloudapiv6.ArgNoHeadersDescription)
	get.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultGetDepth, cloudapiv6.ArgDepthDescription)

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
	create.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
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
	create.AddSetFlag(cloudapiv6.ArgLicenceType, "", "LINUX", constants.EnumLicenceType, "Licence Type of the Volume")
	create.AddStringFlag(constants.FlagType, "", "HDD", "Type of the Volume")
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgLicenceType, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"HDD", "SSD", "SSD Standard", "SSD Premium"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(constants.FlagAvailabilityZone, constants.FlagAvailabilityZoneShort, "AUTO", "Availability zone of the Volume. Storage zone can only be selected prior provisioning")
	_ = create.Command.RegisterFlagCompletionFunc(constants.FlagAvailabilityZone, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"AUTO", "ZONE_1", "ZONE_2", "ZONE_3"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddUUIDFlag(cloudapiv6.ArgBackupUnitId, "", "", "The unique Id of the Backup Unit that User has access to. It is mandatory to provide either 'public image' or 'imageAlias' in conjunction with this property")
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgBackupUnitId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.BackupUnitsIds(), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddUUIDFlag(cloudapiv6.ArgImageId, "", "", "The Image Id or Snapshot Id to be used as template for the new Volume. A password or SSH Key need to be set")
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgImageId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ImageIds(func(r ionoscloud.ApiImagesGetRequest) ionoscloud.ApiImagesGetRequest {
			// Completer for HDD images that are in the same location as the datacenter
			chosenDc, _, err := client.Must().CloudClient.DataCentersApi.DatacentersFindById(context.Background(),
				viper.GetString(core.GetFlagName(create.NS, cloudapiv6.ArgDataCenterId))).Execute()
			if err != nil || chosenDc.Properties == nil || chosenDc.Properties.Location == nil {
				return ionoscloud.ApiImagesGetRequest{}
			}

			return r.Filter("location", *chosenDc.Properties.Location).Filter("imageType", "HDD")
		}), cobra.ShellCompDirectiveNoFileComp
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
	create.AddBoolFlag(constants.ArgWaitForRequest, constants.ArgWaitForRequestShort, constants.DefaultWait, "Wait for the Request for Volume creation to be executed")
	create.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, constants.DefaultTimeoutSeconds, "Timeout option for Request for Volume creation [seconds]")
	create.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultCreateDepth, cloudapiv6.ArgDepthDescription)

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
	update.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddUUIDFlag(cloudapiv6.ArgVolumeId, cloudapiv6.ArgIdShort, "", cloudapiv6.VolumeId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgVolumeId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.VolumesIds(viper.GetString(core.GetFlagName(update.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
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
	update.AddBoolFlag(constants.ArgWaitForRequest, constants.ArgWaitForRequestShort, constants.DefaultWait, "Wait for the Request for Volume update to be executed")
	update.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, constants.DefaultTimeoutSeconds, "Timeout option for Request for Volume update [seconds]")
	update.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultUpdateDepth, cloudapiv6.ArgDepthDescription)

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
	deleteCmd.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddUUIDFlag(cloudapiv6.ArgVolumeId, cloudapiv6.ArgIdShort, "", cloudapiv6.VolumeId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgVolumeId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.VolumesIds(viper.GetString(core.GetFlagName(deleteCmd.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddBoolFlag(constants.ArgWaitForRequest, constants.ArgWaitForRequestShort, constants.DefaultWait, "Wait for the Request for Volume deletion to be executed")
	deleteCmd.AddBoolFlag(cloudapiv6.ArgAll, cloudapiv6.ArgAllShort, false, "Delete all Volumes from a virtual Datacenter.")
	deleteCmd.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, constants.DefaultTimeoutSeconds, "Timeout option for Request for Volume deletion [seconds]")
	deleteCmd.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultDeleteDepth, cloudapiv6.ArgDepthDescription)

	return volumeCmd
}

func PreRunVolumeList(c *core.PreCommandConfig) error {
	if err := core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{cloudapiv6.ArgDataCenterId},
		[]string{cloudapiv6.ArgAll},
	); err != nil {
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

	setRequiredFlagsPublicImage := [][]string{
		{cloudapiv6.ArgDataCenterId, cloudapiv6.ArgImageId, constants.ArgPassword},
		{cloudapiv6.ArgDataCenterId, cloudapiv6.ArgImageId, cloudapiv6.ArgSshKeyPaths},
		{cloudapiv6.ArgDataCenterId, cloudapiv6.ArgImageAlias, constants.ArgPassword},
		{cloudapiv6.ArgDataCenterId, cloudapiv6.ArgImageAlias, cloudapiv6.ArgSshKeyPaths},
	}
	publicImageAsImageId := false

	// image-id and image-alias create more flag requirements
	if fn := core.GetFlagName(c.NS, cloudapiv6.ArgImageId); viper.IsSet(fn) {
		setRequiredFlagsPrivateImage := [][]string{
			{cloudapiv6.ArgDataCenterId, cloudapiv6.ArgImageId},
			{cloudapiv6.ArgDataCenterId, cloudapiv6.ArgImageAlias},
		}

		img, _, err := client.Must().CloudClient.ImagesApi.ImagesFindById(context.Background(),
			viper.GetString(fn)).Execute()
		if err != nil {
			return fmt.Errorf("failed getting image %s: %w", viper.GetString(fn), err)
		}

		if !*img.Properties.Public {
			return core.CheckRequiredFlagsSets(c.Command, c.NS, setRequiredFlagsPrivateImage...)
		}
		publicImageAsImageId = true
	}

	// Public image requirements. Private images cannot have aliases
	if fn := core.GetFlagName(c.NS, cloudapiv6.ArgImageAlias); publicImageAsImageId || viper.IsSet(fn) {
		return core.CheckRequiredFlagsSets(c.Command, c.NS, setRequiredFlagsPublicImage...)
	}
	return nil

}

func RunVolumeListAll(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	datacenters, _, err := c.CloudApiV6Services.DataCenters().List(cloudapiv6.ParentResourceListQueryParams)
	if err != nil {
		return err
	}

	allDcs := getDataCenters(datacenters)
	var allVolumes []ionoscloud.Volumes
	var allVolumesConverted []map[string]interface{}
	totalTime := time.Duration(0)

	for _, dc := range allDcs {
		id, ok := dc.GetIdOk()
		if !ok || id == nil {
			return fmt.Errorf("could not retrieve Datacenter Id")
		}

		volumes, resp, err := c.CloudApiV6Services.Volumes().List(*dc.GetId(), listQueryParams)
		if err != nil {
			return err
		}

		items, ok := volumes.GetItemsOk()
		if !ok || items == nil {
			continue
		}

		for _, item := range *items {
			temp, err := json2table.ConvertJSONToTable("", allVolumeJSONPaths, item)
			if err != nil {
				fmt.Errorf("could not convert from JSON to Table format: %w", err)
			}

			temp[0]["DatacenterId"] = id
			allVolumesConverted = append(allVolumesConverted, temp[0])
		}

		allVolumes = append(allVolumes, volumes.Volumes)
		totalTime += resp.RequestTime
	}

	if totalTime != time.Duration(0) {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, totalTime))
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	out, err := jsontabwriter.GenerateOutputPreconverted(allVolumes, allVolumesConverted,
		tabheaders.GetHeaders(allVolumeCols, defaultVolumeCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)
	return nil
}

func RunVolumeList(c *core.CommandConfig) error {
	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAll)) {
		return RunVolumeListAll(c)
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
		"Listing Volumes from Datacenter with ID: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))))
	// Add Query Parameters for GET Requests
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	if !structs.IsZero(listQueryParams) {
		if listQueryParams.Filters != nil {
			filters := *listQueryParams.Filters

			if val, ok := filters["size"]; ok {
				convertedSize, err := utils.ConvertSize(val[0], utils.GigaBytes)
				if err != nil {
					return err
				}

				filters["size"] = []string{strconv.Itoa(convertedSize)}
				listQueryParams.Filters = &filters
			}
		}
	}

	volumes, resp, err := c.CloudApiV6Services.Volumes().List(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)), listQueryParams)
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	out, err := jsontabwriter.GenerateOutput("items", allVolumeJSONPaths, volumes.Volumes,
		tabheaders.GetHeaders(allVolumeCols, defaultVolumeCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)
	return nil
}

func RunVolumeGet(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	queryParams := listQueryParams.QueryParams

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
		constants.DatacenterId, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
		"Volume with id: %v is getting...", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgVolumeId))))

	vol, resp, err := c.CloudApiV6Services.Volumes().Get(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgVolumeId)),
		queryParams,
	)
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	out, err := jsontabwriter.GenerateOutput("", allVolumeJSONPaths, vol.Volume,
		tabheaders.GetHeaders(allVolumeCols, defaultVolumeCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)
	return nil
}

func RunVolumeCreate(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	queryParams := listQueryParams.QueryParams

	input, err := getNewVolume(c)
	if err != nil {
		return err
	}

	vol, resp, err := c.CloudApiV6Services.Volumes().Create(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)), *input, queryParams)
	if resp != nil && utils.GetId(resp) != "" {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, utils.GetId(resp), resp.RequestTime))
	}
	if err != nil {
		return err
	}

	if err = utils.WaitForRequest(c, waiter.RequestInterrogator, utils.GetId(resp)); err != nil {
		return err
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	out, err := jsontabwriter.GenerateOutput("", allVolumeJSONPaths, vol.Volume,
		tabheaders.GetHeaders(allVolumeCols, defaultVolumeCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)
	return nil
}

func RunVolumeUpdate(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	queryParams := listQueryParams.QueryParams

	input, err := getVolumeInfo(c)
	if err != nil {
		return err
	}

	vol, resp, err := c.CloudApiV6Services.Volumes().Update(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgVolumeId)), *input, queryParams)
	if resp != nil && utils.GetId(resp) != "" {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, utils.GetId(resp), resp.RequestTime))
	}
	if err != nil {
		return err
	}

	if err = utils.WaitForRequest(c, waiter.RequestInterrogator, utils.GetId(resp)); err != nil {
		return err
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	out, err := jsontabwriter.GenerateOutput("", allVolumeJSONPaths, vol.Volume,
		tabheaders.GetHeaders(allVolumeCols, defaultVolumeCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)
	return nil
}

func RunVolumeDelete(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	queryParams := listQueryParams.QueryParams
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))
	volumeId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgVolumeId))

	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAll)) {
		if err := DeleteAllVolumes(c); err != nil {
			return err
		}

		return nil
	}

	if !confirm.FAsk(c.Command.Command.InOrStdin(), "delete volume", viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Starting deleting Volume with id: %v...", volumeId))

	resp, err := c.CloudApiV6Services.Volumes().Delete(dcId, volumeId, queryParams)
	if resp != nil && utils.GetId(resp) != "" {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, utils.GetId(resp), resp.RequestTime))
	}
	if err != nil {
		return err
	}

	if err = utils.WaitForRequest(c, waiter.RequestInterrogator, utils.GetId(resp)); err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), jsontabwriter.GenerateLogOutput("Volume successfully deleted"))
	return nil
}

func getNewVolume(c *core.CommandConfig) (*resources.Volume, error) {
	proper := resources.VolumeProperties{}

	name := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName))
	bus := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgBus))
	volumeType := viper.GetString(core.GetFlagName(c.NS, constants.FlagType))
	availabilityZone := viper.GetString(core.GetFlagName(c.NS, constants.FlagAvailabilityZone))

	// It will get the default values, if flags not set
	proper.SetName(name)
	proper.SetBus(bus)
	proper.SetType(volumeType)
	proper.SetAvailabilityZone(availabilityZone)

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
		"Properties set for creating the Volume: Name: %v, Bus: %v, VolumeType: %v, AvailabilityZone: %v",
		name, bus, volumeType, availabilityZone))

	size, err := utils.ConvertSize(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgSize)),
		utils.GigaBytes,
	)
	if err != nil {
		return nil, err
	}

	proper.SetSize(float32(size))

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property Size set: %vGB", float32(size)))

	// Check if flags are set and set options
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgBackupUnitId)) {
		backupUnitId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgBackupUnitId))
		proper.SetBackupunitId(backupUnitId)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property BackupUnitId set: %v", backupUnitId))
	}

	if (!viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgImageId)) &&
		!viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgImageAlias))) ||
		viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgLicenceType)) {
		licenceType := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLicenceType))
		proper.SetLicenceType(licenceType)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property LicenceType set: %v", licenceType))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgImageId)) {
		imageId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgImageId))
		proper.SetImage(imageId)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property Image set: %v", imageId))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgImageAlias)) {
		imageAlias := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgImageAlias))
		proper.SetImageAlias(imageAlias)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property ImageAlias set: %v", imageAlias))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgPassword)) {
		imagePassword := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgPassword))
		proper.SetImagePassword(imagePassword)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property ImagePassword set"))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgSshKeyPaths)) {
		sshKeysPaths := viper.GetStringSlice(core.GetFlagName(c.NS, cloudapiv6.ArgSshKeyPaths))

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("SSH Key Paths: %v", sshKeysPaths))

		sshKeys, err := getSshKeysFromPaths(sshKeysPaths)
		if err != nil {
			return nil, err
		}

		proper.SetSshKeys(sshKeys)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property SshKeys set"))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgUserData)) {
		userData := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgUserData))
		proper.SetUserData(userData)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property UserData set: %v", userData))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgCpuHotPlug)) {
		cpuHotPlug := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgCpuHotPlug))
		proper.SetCpuHotPlug(cpuHotPlug)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property CpuHotPlug set: %v", cpuHotPlug))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgRamHotPlug)) {
		ramHotPlug := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgRamHotPlug))
		proper.SetRamHotPlug(ramHotPlug)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property RamHotPlug set: %v", ramHotPlug))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgNicHotPlug)) {
		nicHotPlug := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgNicHotPlug))
		proper.SetNicHotPlug(nicHotPlug)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property NicHotPlug set: %v", nicHotPlug))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgNicHotUnplug)) {
		nicHotUnplug := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgNicHotUnplug))
		proper.SetNicHotUnplug(nicHotUnplug)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property NicHotUnplug set: %v", nicHotUnplug))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgDiscVirtioHotPlug)) {
		discVirtioHotPlug := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgDiscVirtioHotPlug))
		proper.SetDiscVirtioHotPlug(discVirtioHotPlug)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property DiscVirtioHotPlug set: %v", discVirtioHotPlug))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgDiscVirtioHotUnplug)) {
		discVirtioHotUnplug := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgDiscVirtioHotUnplug))
		proper.SetDiscVirtioHotUnplug(discVirtioHotUnplug)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property DiscVirtioHotUnplug set: %v", discVirtioHotUnplug))
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

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property Name set: %v", name))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgBus)) {
		bus := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgBus))
		input.SetBus(bus)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property Bus set: %v", bus))
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

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property Size set: %vGB", float32(size)))
	}

	// Check if flags are set and set options
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgCpuHotPlug)) {
		cpuHotPlug := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgCpuHotPlug))
		input.SetCpuHotPlug(cpuHotPlug)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property CpuHotPlug set: %v", cpuHotPlug))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgRamHotPlug)) {
		ramHotPlug := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgRamHotPlug))
		input.SetRamHotPlug(ramHotPlug)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property RamHotPlug set: %v", ramHotPlug))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgNicHotPlug)) {
		nicHotPlug := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgNicHotPlug))
		input.SetNicHotPlug(nicHotPlug)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property NicHotPlug set: %v", nicHotPlug))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgNicHotUnplug)) {
		nicHotUnplug := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgNicHotUnplug))
		input.SetNicHotUnplug(nicHotUnplug)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property NicHotUnplug set: %v", nicHotUnplug))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgDiscVirtioHotPlug)) {
		discVirtioHotPlug := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgDiscVirtioHotPlug))
		input.SetDiscVirtioHotPlug(discVirtioHotPlug)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property DiscVirtioHotPlug set: %v", discVirtioHotPlug))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgDiscVirtioHotUnplug)) {
		discVirtioHotUnplug := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgDiscVirtioHotUnplug))
		input.SetDiscVirtioHotUnplug(discVirtioHotUnplug)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property DiscVirtioHotUnplug set: %v", discVirtioHotUnplug))
	}

	return &input, nil
}

func DeleteAllVolumes(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	queryParams := listQueryParams.QueryParams
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.DatacenterId, dcId))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Getting Volumes..."))

	volumes, resp, err := c.CloudApiV6Services.Volumes().List(dcId, cloudapiv6.ParentResourceListQueryParams)
	if err != nil {
		return err
	}

	volumesItems, ok := volumes.GetItemsOk()
	if !ok || volumesItems == nil {
		return fmt.Errorf("could not get items of Volumes")
	}

	if len(*volumesItems) <= 0 {
		return fmt.Errorf("no Volumes found")
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateLogOutput("Volumes to be deleted:"))

	for _, volume := range *volumesItems {
		if id, ok := volume.GetIdOk(); ok && id != nil {
			fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateLogOutput("Volume Id: %v", *id))
		}

		if properties, ok := volume.GetPropertiesOk(); ok && properties != nil {
			if name, ok := properties.GetNameOk(); ok && name != nil {
				fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateLogOutput("Volume Name: %v", *name))
			}
		}
	}

	if !confirm.FAsk(c.Command.Command.InOrStdin(), "delete all the Volumes", viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Deleting all the Volumes..."))

	var multiErr error
	for _, volume := range *volumesItems {
		id, ok := volume.GetIdOk()
		if !ok || id == nil {
			continue
		}

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Starting deleting Volume with id: %v is...", *id))

		resp, err = c.CloudApiV6Services.Volumes().Delete(dcId, *id, queryParams)
		if resp != nil && utils.GetId(resp) != "" {
			fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, utils.GetId(resp), resp.RequestTime))
		}
		if err != nil {
			multiErr = errors.Join(multiErr, fmt.Errorf(constants.ErrDeleteAll, c.Resource, *id, err))
			continue
		}

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateLogOutput(constants.MessageDeletingAll, c.Resource, *id))

		if err = utils.WaitForRequest(c, waiter.RequestInterrogator, utils.GetId(resp)); err != nil {
			multiErr = errors.Join(multiErr, fmt.Errorf(constants.ErrWaitDeleteAll, c.Resource, *id, err))
		}
	}

	if multiErr != nil {
		return multiErr
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), jsontabwriter.GenerateLogOutput("Volumes successfully deleted"))
	return nil
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
	attachVolume.AddStringSliceFlag(constants.ArgCols, "", defaultVolumeCols, tabheaders.ColsMessage(allVolumeCols))
	_ = attachVolume.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allVolumeCols, cobra.ShellCompDirectiveNoFileComp
	})
	attachVolume.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = attachVolume.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	attachVolume.AddUUIDFlag(cloudapiv6.ArgVolumeId, cloudapiv6.ArgIdShort, "", cloudapiv6.VolumeId, core.RequiredFlagOption())
	_ = attachVolume.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgVolumeId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.VolumesIds(viper.GetString(core.GetFlagName(attachVolume.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	attachVolume.AddUUIDFlag(cloudapiv6.ArgServerId, "", "", cloudapiv6.ServerId, core.RequiredFlagOption())
	_ = attachVolume.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgServerId, func(cmd *cobra.Command, ags []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ServersIds(viper.GetString(core.GetFlagName(attachVolume.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	attachVolume.AddBoolFlag(constants.ArgWaitForRequest, constants.ArgWaitForRequestShort, constants.DefaultWait, "Wait for the Request for Volume attachment to be executed")
	attachVolume.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, constants.DefaultTimeoutSeconds, "Timeout option for Request for Volume attachment [seconds]")

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
	listVolumes.AddStringSliceFlag(constants.ArgCols, "", defaultVolumeCols, tabheaders.ColsMessage(allVolumeCols))
	_ = listVolumes.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allVolumeCols, cobra.ShellCompDirectiveNoFileComp
	})
	listVolumes.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = listVolumes.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	listVolumes.AddUUIDFlag(cloudapiv6.ArgServerId, "", "", cloudapiv6.ServerId, core.RequiredFlagOption())
	_ = listVolumes.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgServerId, func(cmd *cobra.Command, ags []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ServersIds(viper.GetString(core.GetFlagName(listVolumes.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	listVolumes.AddInt32Flag(constants.FlagMaxResults, constants.FlagMaxResultsShort, cloudapiv6.DefaultMaxResults, constants.DescMaxResults)
	listVolumes.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultListDepth, cloudapiv6.ArgDepthDescription)
	listVolumes.AddStringFlag(cloudapiv6.ArgOrderBy, "", "", cloudapiv6.ArgOrderByDescription)
	_ = listVolumes.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgOrderBy, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.VolumesFilters(), cobra.ShellCompDirectiveNoFileComp
	})
	listVolumes.AddStringSliceFlag(cloudapiv6.ArgFilters, cloudapiv6.ArgFiltersShort, []string{""}, cloudapiv6.ArgFiltersDescription)
	_ = listVolumes.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgFilters, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.VolumesFilters(), cobra.ShellCompDirectiveNoFileComp
	})
	listVolumes.AddBoolFlag(constants.ArgNoHeaders, "", false, cloudapiv6.ArgNoHeadersDescription)

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
	getVolumeCmd.AddStringSliceFlag(constants.ArgCols, "", defaultVolumeCols, tabheaders.ColsMessage(allVolumeCols))
	_ = getVolumeCmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allVolumeCols, cobra.ShellCompDirectiveNoFileComp
	})
	getVolumeCmd.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = getVolumeCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	getVolumeCmd.AddUUIDFlag(cloudapiv6.ArgServerId, "", "", cloudapiv6.ServerId, core.RequiredFlagOption())
	_ = getVolumeCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgServerId, func(cmd *cobra.Command, ags []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ServersIds(viper.GetString(core.GetFlagName(getVolumeCmd.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	getVolumeCmd.AddUUIDFlag(cloudapiv6.ArgVolumeId, cloudapiv6.ArgIdShort, "", cloudapiv6.VolumeId, core.RequiredFlagOption())
	_ = getVolumeCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgVolumeId, func(cmd *cobra.Command, ags []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.AttachedVolumesIds(viper.GetString(core.GetFlagName(getVolumeCmd.NS, cloudapiv6.ArgDataCenterId)), viper.GetString(core.GetFlagName(getVolumeCmd.NS, cloudapiv6.ArgServerId))), cobra.ShellCompDirectiveNoFileComp
	})
	getVolumeCmd.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultMiscDepth, cloudapiv6.ArgDepthDescription)

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
	detachVolume.AddStringSliceFlag(constants.ArgCols, "", defaultVolumeCols, tabheaders.ColsMessage(allVolumeCols))
	_ = detachVolume.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allVolumeCols, cobra.ShellCompDirectiveNoFileComp
	})
	detachVolume.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = detachVolume.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	detachVolume.AddUUIDFlag(cloudapiv6.ArgVolumeId, cloudapiv6.ArgIdShort, "", cloudapiv6.VolumeId, core.RequiredFlagOption())
	_ = detachVolume.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgVolumeId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.AttachedVolumesIds(viper.GetString(core.GetFlagName(detachVolume.NS, cloudapiv6.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(detachVolume.NS, cloudapiv6.ArgServerId))), cobra.ShellCompDirectiveNoFileComp
	})
	detachVolume.AddUUIDFlag(cloudapiv6.ArgServerId, "", "", cloudapiv6.ServerId, core.RequiredFlagOption())
	_ = detachVolume.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ServersIds(viper.GetString(core.GetFlagName(detachVolume.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	detachVolume.AddBoolFlag(constants.ArgWaitForRequest, constants.ArgWaitForRequestShort, constants.DefaultWait, "Wait for the Request for Volume detachment to be executed")
	detachVolume.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, constants.DefaultTimeoutSeconds, "Timeout option for Request for Volume detachment [seconds]")
	detachVolume.AddBoolFlag(cloudapiv6.ArgAll, cloudapiv6.ArgAllShort, false, "Detach all Volumes.")
	detachVolume.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultMiscDepth, cloudapiv6.ArgDepthDescription)

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
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))
	serverId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgServerId))
	volumeId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgVolumeId))

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.DatacenterId, dcId))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
		"Attaching Volume with ID: %v to Server with ID: %v...", volumeId, serverId))

	attachedVol, resp, err := c.CloudApiV6Services.Servers().AttachVolume(dcId, serverId, volumeId, queryParams)
	if resp != nil && utils.GetId(resp) != "" {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, utils.GetId(resp), resp.RequestTime))
	}
	if err != nil {
		return err
	}

	if err = utils.WaitForRequest(c, waiter.RequestInterrogator, utils.GetId(resp)); err != nil {
		return err
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	out, err := jsontabwriter.GenerateOutput("", allVolumeJSONPaths, attachedVol.Volume,
		tabheaders.GetHeaders(allVolumeCols, defaultVolumeCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)
	return nil
}

func RunServerVolumesList(c *core.CommandConfig) error {
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))
	serverId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgServerId))

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.DatacenterId, dcId))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Listing attached Volumes from Server with ID: %v...", serverId))

	// Add Query Parameters for GET Requests
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	if !structs.IsZero(listQueryParams) {
		if listQueryParams.Filters != nil {
			filters := *listQueryParams.Filters

			if val, ok := filters["size"]; ok {
				convertedSize, err := utils.ConvertSize(val[0], utils.GigaBytes)
				if err != nil {
					return err
				}

				filters["size"] = []string{strconv.Itoa(convertedSize)}
				listQueryParams.Filters = &filters
			}
		}
	}

	attachedVols, resp, err := c.CloudApiV6Services.Servers().ListVolumes(dcId, serverId, listQueryParams)
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	out, err := jsontabwriter.GenerateOutput("items", allVolumeJSONPaths, attachedVols.AttachedVolumes,
		tabheaders.GetHeaders(allVolumeCols, defaultVolumeCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)
	return nil
}

func RunServerVolumeGet(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	queryParams := listQueryParams.QueryParams
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))
	serverId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgServerId))
	volumeId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgVolumeId))

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.DatacenterId, dcId))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
		"Getting attached Volume with ID: %v from Server with ID: %v...", volumeId, serverId))

	attachedVol, resp, err := c.CloudApiV6Services.Servers().GetVolume(dcId, serverId, volumeId, queryParams)
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	out, err := jsontabwriter.GenerateOutput("", allVolumeJSONPaths, attachedVol.Volume,
		tabheaders.GetHeaders(allVolumeCols, defaultVolumeCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)
	return nil
}

func RunServerVolumeDetach(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	queryParams := listQueryParams.QueryParams

	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAll)) {
		if err := DetachAllServerVolumes(c); err != nil {
			return err
		}

		return nil
	}
	if !confirm.FAsk(c.Command.Command.InOrStdin(), "detach volume from server", viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))
	serverId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgServerId))
	volumeId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgVolumeId))

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.DatacenterId, dcId))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
		"Detaching Volume with ID: %v from Server with ID: %v...", volumeId, serverId))

	resp, err := c.CloudApiV6Services.Servers().DetachVolume(dcId, serverId, volumeId, queryParams)
	if resp != nil && utils.GetId(resp) != "" {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, utils.GetId(resp), resp.RequestTime))
	}
	if err != nil {
		return err
	}

	if err = utils.WaitForRequest(c, waiter.RequestInterrogator, utils.GetId(resp)); err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), jsontabwriter.GenerateLogOutput("Volume successfully detached"))
	return nil
}

func DetachAllServerVolumes(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	queryParams := listQueryParams.QueryParams
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))
	serverId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgServerId))

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.DatacenterId, dcId))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Server ID: %v", serverId))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Getting Volumes..."))

	volumes, resp, err := c.CloudApiV6Services.Servers().ListVolumes(dcId, serverId, cloudapiv6.ParentResourceListQueryParams)
	if err != nil {
		return err
	}

	volumesItems, ok := volumes.GetItemsOk()
	if !ok || volumesItems == nil {
		return fmt.Errorf("could not get items of Volumes")
	}

	if len(*volumesItems) <= 0 {
		return fmt.Errorf("no Volumes found")
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateLogOutput("Volumes to be detached:"))

	for _, volume := range *volumesItems {
		delIdAndName := ""
		if id, ok := volume.GetIdOk(); ok && id != nil {
			delIdAndName += "Volume Id: " + *id
		}

		if properties, ok := volume.GetPropertiesOk(); ok && properties != nil {
			if name, ok := properties.GetNameOk(); ok && name != nil {
				delIdAndName += " Volume Name: " + *name
			}
		}

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateLogOutput(delIdAndName))
	}

	if !confirm.FAsk(c.Command.Command.InOrStdin(), "detach all the Volumes", viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Detaching all the Volumes..."))

	var multiErr error
	for _, volume := range *volumesItems {
		id, ok := volume.GetIdOk()
		if !ok || id == nil {
			continue
		}

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Starting detaching Volume with id: %v...", *id))

		resp, err = c.CloudApiV6Services.Servers().DetachVolume(dcId, serverId, *id, queryParams)
		if resp != nil && utils.GetId(resp) != "" {
			fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, utils.GetId(resp), resp.RequestTime))
		}
		if err != nil {
			multiErr = errors.Join(multiErr, fmt.Errorf(constants.ErrDeleteAll, c.Resource, *id, err))
			continue
		}

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateLogOutput(constants.MessageRemovingAll, c.Resource, *id))

		if err = utils.WaitForRequest(c, waiter.RequestInterrogator, utils.GetId(resp)); err != nil {
			multiErr = errors.Join(multiErr, fmt.Errorf(constants.ErrWaitDeleteAll, c.Resource, *id, err))
		}
	}

	if multiErr != nil {
		return multiErr
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), jsontabwriter.GenerateLogOutput("Volumes successfully detached"))
	return nil
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
