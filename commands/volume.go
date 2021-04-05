package commands

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/pkg/builder"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/resources"
	"github.com/ionos-cloud/ionosctl/pkg/utils"
	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
	"github.com/ionos-cloud/ionosctl/pkg/utils/printer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func volume() *builder.Command {
	volumeCmd := &builder.Command{
		Command: &cobra.Command{
			Use:              "volume",
			Aliases:          []string{"vol"},
			Short:            "Volume Operations",
			Long:             `The sub-commands of ` + "`" + `ionosctl volume` + "`" + ` manage your block storage volumes by creating, updating, getting specific information, deleting Volumes or attaching, detaching a Volume from a Server.`,
			TraverseChildren: true,
		},
	}
	globalFlags := volumeCmd.Command.PersistentFlags()
	globalFlags.StringP(config.ArgDataCenterId, "", "", config.RequiredFlagDatacenterId)
	_ = viper.BindPFlag(builder.GetGlobalFlagName(volumeCmd.Command.Use, config.ArgDataCenterId), globalFlags.Lookup(config.ArgDataCenterId))
	_ = volumeCmd.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	globalFlags.StringSlice(config.ArgCols, defaultDatacenterCols, "Columns to be printed in the standard output")
	_ = viper.BindPFlag(builder.GetGlobalFlagName(volumeCmd.Command.Name(), config.ArgCols), globalFlags.Lookup(config.ArgCols))

	/*
		List Command
	*/
	builder.NewCommand(context.TODO(), volumeCmd, PreRunGlobalDcIdValidate, RunVolumeList, "list", "List Volumes",
		"Use this command to list all Volumes from a Data Center on your account.\n\nRequired values to run command:\n\n* Data Center Id",
		listVolumeExample, true)

	/*
		Get Command
	*/
	get := builder.NewCommand(context.TODO(), volumeCmd, PreRunGlobalDcIdVolumeIdValidate, RunVolumeGet, "get", "Get a Volume",
		"Use this command to retrieve information about a Volume using its ID.\n\nRequired values to run command:\n\n* Data Center Id\n* Volume Id",
		getVolumeExample, true)
	get.AddStringFlag(config.ArgVolumeId, "", "", config.RequiredFlagVolumeId)
	_ = get.Command.RegisterFlagCompletionFunc(config.ArgVolumeId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getVolumesIds(os.Stderr, viper.GetString(builder.GetGlobalFlagName(volumeCmd.Command.Name(), config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Create Command
	*/
	create := builder.NewCommand(context.TODO(), volumeCmd, PreRunGlobalDcIdValidate, RunVolumeCreate, "create", "Create a Volume",
		`Use this command to create a Volume on your account. You can specify the name, size, type, licence type and availability zone for the object.

You can wait for the action to be executed using `+"`"+`--wait`+"`"+` option.

Required values to run command:

* Data Center Id`,
		createVolumeExample, true)
	create.AddStringFlag(config.ArgVolumeName, "", "", "Name of the Volume")
	create.AddFloat32Flag(config.ArgVolumeSize, "", config.DefaultVolumeSize, "Size in GB of the Volume")
	create.AddStringFlag(config.ArgVolumeBus, "", "VIRTIO", "Bus for the Volume")
	create.AddStringFlag(config.ArgVolumeLicenceType, "", "LINUX", "Licence Type of the Volume")
	create.AddStringFlag(config.ArgVolumeType, "", "HDD", "Type of the Volume")
	create.AddStringFlag(config.ArgVolumeZone, "", "AUTO", "Availability zone of the Volume. Storage zone can only be selected prior provisioning")
	create.AddStringFlag(config.ArgVolumeSshKey, "", "", "Ssh Key of the Volume")
	create.AddBoolFlag(config.ArgWait, "", config.DefaultWait, "Wait for Volume to be created")
	create.AddIntFlag(config.ArgTimeout, "", config.DefaultTimeoutSeconds, "Timeout option for Volume to be created [seconds]")

	/*
		Update Command
	*/
	update := builder.NewCommand(context.TODO(), volumeCmd, PreRunGlobalDcIdVolumeIdValidate, RunVolumeUpdate, "update", "Update a Volume",
		`Use this command to update a Volume. You may increase the size of an existing storage Volume. You cannot reduce the size of an existing storage Volume. The Volume size will be increased without reboot if the appropriate "hot plug" settings have been set to true. The additional capacity is not added to any partition therefore you will need to adjust the partition inside the operating system afterwards.

Once you have increased the Volume size you cannot decrease the Volume size using the Cloud API. Certain attributes can only be set when a Volume is created and are considered immutable once the Volume has been provisioned.

You can wait for the action to be executed using `+"`"+`--wait`+"`"+` option.

Required values to run command:

* Data Center Id
* Volume Id`, updateVolumeExample, true)
	update.AddStringFlag(config.ArgVolumeId, "", "", config.RequiredFlagVolumeId)
	_ = update.Command.RegisterFlagCompletionFunc(config.ArgVolumeId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getVolumesIds(os.Stderr, viper.GetString(builder.GetGlobalFlagName(volumeCmd.Command.Name(), config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(config.ArgVolumeName, "", "", "Name of the Volume")
	update.AddFloat32Flag(config.ArgVolumeSize, "", config.DefaultVolumeSize, "Size in GB of the Volume")
	update.AddStringFlag(config.ArgVolumeBus, "", "VIRTIO", "Bus of the Volume")
	update.AddStringFlag(config.ArgVolumeSshKey, "", "", "Ssh Key of the Volume")
	update.AddBoolFlag(config.ArgWait, "", config.DefaultWait, "Wait for Volume to be updated")
	update.AddIntFlag(config.ArgTimeout, "", config.DefaultTimeoutSeconds, "Timeout option for Volume to be updated [seconds]")

	/*
		Delete Command
	*/
	deleteCmd := builder.NewCommand(context.TODO(), volumeCmd, PreRunGlobalDcIdVolumeIdValidate, RunVolumeDelete, "delete", "Delete a Volume",
		`Use this command to delete specified Volume. This will result in the Volume being removed from your Virtual Data Center. Please use this with caution!

You can wait for the action to be executed using `+"`"+`--wait`+"`"+` option. You can force the command to execute without user input using `+"`"+`--ignore-stdin`+"`"+` option.

Required values to run command:

* Data Center Id
* Volume Id`, deleteVolumeExample, true)
	deleteCmd.AddStringFlag(config.ArgVolumeId, "", "", config.RequiredFlagVolumeId)
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(config.ArgVolumeId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getVolumesIds(os.Stderr, viper.GetString(builder.GetGlobalFlagName(volumeCmd.Command.Name(), config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddBoolFlag(config.ArgWait, "", config.DefaultWait, "Wait for Volume to be deleted")
	deleteCmd.AddIntFlag(config.ArgTimeout, "", config.DefaultTimeoutSeconds, "Timeout option for Volume to be deleted [seconds]")

	/*
		Attach Command
	*/
	attachVolume := builder.NewCommand(context.TODO(), volumeCmd, PreRunGlobalDcIdServerVolumeIdsValidate, RunVolumeAttach, "attach", "Attach a Volume to a Server",
		`Use this command to attach a Volume to a Server from a Data Center.

You can wait for the action to be executed using `+"`"+`--wait`+"`"+` option.

Required values to run command:

* Data Center Id
* Server Id
* Volume Id

The sub-commands of `+"`"+`ionosctl volume attach`+"`"+` allow you to retrieve information about attached Volumes or about a specified attached Volume.`, attachVolumeExample, true)
	attachVolume.AddStringFlag(config.ArgVolumeId, "", "", config.RequiredFlagVolumeId)
	_ = attachVolume.Command.RegisterFlagCompletionFunc(config.ArgVolumeId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getVolumesIds(os.Stderr, viper.GetString(builder.GetGlobalFlagName(volumeCmd.Command.Name(), config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	attachVolume.AddStringFlag(config.ArgServerId, "", "", config.RequiredFlagServerId)
	_ = attachVolume.Command.RegisterFlagCompletionFunc(config.ArgServerId, func(cmd *cobra.Command, ags []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getServersIds(os.Stderr, viper.GetString(builder.GetGlobalFlagName(volumeCmd.Command.Name(), config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	attachVolume.AddBoolFlag(config.ArgWait, "", config.DefaultWait, "Wait for Volume to attach to Server")
	attachVolume.AddIntFlag(config.ArgTimeout, "", config.DefaultTimeoutSeconds, "Timeout option for Volume to be attached to a Server [seconds]")

	/*
		Attach List Command
	*/
	listAttached := builder.NewCommand(context.TODO(), attachVolume, PreRunAttachGlobalDcIdServerIdValidate, RunVolumesAttachList, "list", "List attached Volumes from a Server",
		"Use this command to get a list of attached Volumes to a Server from a Data Center.\n\nRequired values to run command:\n\n* Data Center Id\n* Server Id",
		attachListVolumeExample, true)
	listAttached.AddStringFlag(config.ArgServerId, "", "", config.RequiredFlagServerId)
	_ = listAttached.Command.RegisterFlagCompletionFunc(config.ArgServerId, func(cmd *cobra.Command, ags []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getServersIds(os.Stderr, viper.GetString(builder.GetGlobalFlagName(volumeCmd.Command.Name(), config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Attach Get Command
	*/
	getAttached := builder.NewCommand(context.TODO(), attachVolume, PreRunAttachGlobalDcIdServerVolumeIdsValidate, RunVolumeAttachGet, "get", "Get an attached Volume from a Server",
		"Use this command to retrieve information about an attached Volume.\n\nRequired values to run command:\n\n* Data Center Id\n* Server Id\n* Volume Id",
		attachGetVolumeExample, true)
	getAttached.AddStringFlag(config.ArgServerId, "", "", config.RequiredFlagServerId)
	_ = getAttached.Command.RegisterFlagCompletionFunc(config.ArgServerId, func(cmd *cobra.Command, ags []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getServersIds(os.Stderr, viper.GetString(builder.GetGlobalFlagName(volumeCmd.Command.Name(), config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	getAttached.AddStringFlag(config.ArgVolumeId, "", "", config.RequiredFlagVolumeId)
	_ = getAttached.Command.RegisterFlagCompletionFunc(config.ArgVolumeId, func(cmd *cobra.Command, ags []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getAttachedVolumesIds(os.Stderr, viper.GetString(builder.GetGlobalFlagName(volumeCmd.Command.Name(), config.ArgDataCenterId)), viper.GetString(builder.GetFlagName(attachVolume.Command.Name(), getAttached.Command.Name(), config.ArgServerId))), cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Detach Command
	*/
	detachVolume := builder.NewCommand(context.TODO(), volumeCmd, PreRunGlobalDcIdServerVolumeIdsValidate, RunVolumeDetach, "detach", "Detach a Volume from a Server",
		`Use this command to detach a Volume from a Server.

You can wait for the action to be executed using `+"`"+`--wait`+"`"+` option. You can force the command to execute without user input using `+"`"+`--ignore-stdin`+"`"+` option.

Required values to run command:

* Data Center Id
* Server Id
* Volume Id`, detachVolumeExample, true)
	detachVolume.AddStringFlag(config.ArgVolumeId, "", "", config.RequiredFlagVolumeId)
	_ = detachVolume.Command.RegisterFlagCompletionFunc(config.ArgVolumeId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getAttachedVolumesIds(os.Stderr, viper.GetString(builder.GetGlobalFlagName(volumeCmd.Command.Name(), config.ArgDataCenterId)), viper.GetString(builder.GetFlagName(volumeCmd.Command.Name(), detachVolume.Command.Name(), config.ArgServerId))), cobra.ShellCompDirectiveNoFileComp
	})
	detachVolume.AddStringFlag(config.ArgServerId, "", "", config.RequiredFlagServerId)
	_ = detachVolume.Command.RegisterFlagCompletionFunc(config.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getServersIds(os.Stderr, viper.GetString(builder.GetGlobalFlagName(volumeCmd.Command.Name(), config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	detachVolume.AddBoolFlag(config.ArgWait, "", config.DefaultWait, "Wait for Volume to detach from Server")
	detachVolume.AddIntFlag(config.ArgTimeout, "", config.DefaultTimeoutSeconds, "Timeout option for Server to be detached from a Server [seconds]")

	return volumeCmd
}

func PreRunGlobalDcIdVolumeIdValidate(c *builder.PreCommandConfig) error {
	err := builder.CheckRequiredGlobalFlags(c.ParentName, config.ArgDataCenterId)
	if err != nil {
		return err
	}
	err = builder.CheckRequiredFlags(c.ParentName, c.Name, config.ArgVolumeId)
	if err != nil {
		return err
	}
	return nil
}

func PreRunGlobalDcIdServerVolumeIdsValidate(c *builder.PreCommandConfig) error {
	err := builder.CheckRequiredGlobalFlags(c.ParentName, config.ArgDataCenterId)
	if err != nil {
		return err
	}
	err = builder.CheckRequiredFlags(c.ParentName, c.Name, config.ArgServerId, config.ArgVolumeId)
	if err != nil {
		return err
	}
	return nil
}

func RunVolumeList(c *builder.CommandConfig) error {
	volumes, _, err := c.Volumes().List(viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgDataCenterId)))
	if err != nil {
		return err
	}
	ss := getVolumes(volumes)
	return c.Printer.Print(printer.Result{
		OutputJSON: volumes,
		KeyValue:   getVolumesKVMaps(ss),
		Columns:    getVolumesCols(builder.GetGlobalFlagName(c.ParentName, config.ArgCols), c.Printer.GetStderr()),
	})
}

func RunVolumeGet(c *builder.CommandConfig) error {
	volume, _, err := c.Volumes().Get(
		viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgDataCenterId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgVolumeId)),
	)
	if err != nil {
		return err
	}
	return c.Printer.Print(printer.Result{
		OutputJSON: volume,
		KeyValue:   getVolumesKVMaps([]resources.Volume{*volume}),
		Columns:    getVolumesCols(builder.GetGlobalFlagName(c.ParentName, config.ArgCols), c.Printer.GetStderr()),
	})
}

func RunVolumeCreate(c *builder.CommandConfig) error {
	volume, resp, err := c.Volumes().Create(
		viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgDataCenterId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgVolumeName)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgVolumeBus)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgVolumeType)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgVolumeLicenceType)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgVolumeZone)),
		float32(viper.GetFloat64(builder.GetFlagName(c.ParentName, c.Name, config.ArgVolumeSize))),
	)
	if err != nil {
		return err
	}
	err = waitForAction(c, printer.GetRequestPath(resp))
	if err != nil {
		return err
	}
	return c.Printer.Print(printer.Result{
		OutputJSON:  volume,
		KeyValue:    getVolumesKVMaps([]resources.Volume{*volume}),
		Columns:     getVolumesCols(builder.GetGlobalFlagName(c.ParentName, config.ArgCols), c.Printer.GetStderr()),
		ApiResponse: resp,
		Resource:    "volume",
		Verb:        "create",
		WaitFlag:    viper.GetBool(builder.GetFlagName(c.ParentName, c.Name, config.ArgWait)),
	})
}

func RunVolumeUpdate(c *builder.CommandConfig) error {
	input := resources.VolumeProperties{}
	if viper.IsSet(builder.GetFlagName(c.ParentName, c.Name, config.ArgVolumeName)) {
		input.SetName(viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgVolumeName)))
	}
	if viper.IsSet(builder.GetFlagName(c.ParentName, c.Name, config.ArgVolumeBus)) {
		input.SetBus(viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgVolumeBus)))
	}
	if viper.IsSet(builder.GetFlagName(c.ParentName, c.Name, config.ArgVolumeSize)) {
		input.SetSize(float32(viper.GetFloat64(builder.GetFlagName(c.ParentName, c.Name, config.ArgVolumeSize))))
	}
	volume, resp, err := c.Volumes().Update(
		viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgDataCenterId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgVolumeId)),
		input,
	)
	if err != nil {
		return err
	}
	err = waitForAction(c, printer.GetRequestPath(resp))
	if err != nil {
		return err
	}
	return c.Printer.Print(printer.Result{
		OutputJSON:  volume,
		KeyValue:    getVolumesKVMaps([]resources.Volume{*volume}),
		Columns:     getVolumesCols(builder.GetGlobalFlagName(c.ParentName, config.ArgCols), c.Printer.GetStderr()),
		ApiResponse: resp,
		Resource:    "volume",
		Verb:        "update",
		WaitFlag:    viper.GetBool(builder.GetFlagName(c.ParentName, c.Name, config.ArgWait)),
	})
}

func RunVolumeDelete(c *builder.CommandConfig) error {
	err := utils.AskForConfirm(c.Stdin, c.Printer, "delete volume")
	if err != nil {
		return err
	}
	resp, err := c.Volumes().Delete(
		viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgDataCenterId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgVolumeId)),
	)
	if err != nil {
		return err
	}
	err = waitForAction(c, printer.GetRequestPath(resp))
	if err != nil {
		return err
	}
	return c.Printer.Print(printer.Result{
		ApiResponse: resp,
		Resource:    "volume",
		Verb:        "delete",
		WaitFlag:    viper.GetBool(builder.GetFlagName(c.ParentName, c.Name, config.ArgWait)),
	})
}

func PreRunAttachGlobalDcIdServerVolumeIdsValidate(c *builder.PreCommandConfig) error {
	// Data Center Id is inherited from command "volume"
	err := builder.CheckRequiredGlobalFlags("volume", config.ArgDataCenterId)
	if err != nil {
		return err
	}
	err = builder.CheckRequiredFlags(c.ParentName, c.Name, config.ArgServerId, config.ArgVolumeId)
	if err != nil {
		return err
	}
	return nil
}

func PreRunAttachGlobalDcIdServerIdValidate(c *builder.PreCommandConfig) error {
	// Data Center Id is inherited from command "volume"
	err := builder.CheckRequiredGlobalFlags("volume", config.ArgDataCenterId)
	if err != nil {
		return err
	}
	err = builder.CheckRequiredFlags(c.ParentName, c.Name, config.ArgServerId)
	if err != nil {
		return err
	}
	return nil
}

func RunVolumeAttach(c *builder.CommandConfig) error {
	attachedvol, resp, err := c.Volumes().Attach(
		viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgDataCenterId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgServerId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgVolumeId)),
	)
	if err != nil {
		return err
	}
	err = waitForAction(c, printer.GetRequestPath(resp))
	if err != nil {
		return err
	}
	return c.Printer.Print(printer.Result{
		OutputJSON:  attachedvol,
		KeyValue:    getVolumesKVMaps([]resources.Volume{*attachedvol}),
		Columns:     getVolumesCols(builder.GetGlobalFlagName(c.ParentName, config.ArgCols), c.Printer.GetStderr()),
		ApiResponse: resp,
		Resource:    "volume",
		Verb:        "attach",
		WaitFlag:    viper.GetBool(builder.GetFlagName(c.ParentName, c.Name, config.ArgWait)),
	})
}

func RunVolumesAttachList(c *builder.CommandConfig) error {
	attachedvols, _, err := c.Volumes().ListAttached(
		viper.GetString(builder.GetGlobalFlagName("volume", config.ArgDataCenterId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgServerId)),
	)
	if err != nil {
		return err
	}
	vs := getAttachedVolumes(attachedvols)
	return c.Printer.Print(printer.Result{
		OutputJSON: attachedvols,
		KeyValue:   getVolumesKVMaps(vs),
		Columns:    getVolumesCols(builder.GetGlobalFlagName(c.ParentName, config.ArgCols), c.Printer.GetStderr()),
	})
}

func RunVolumeAttachGet(c *builder.CommandConfig) error {
	volume, _, err := c.Volumes().GetAttached(
		viper.GetString(builder.GetGlobalFlagName("volume", config.ArgDataCenterId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgServerId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgVolumeId)),
	)
	if err != nil {
		return err
	}
	return c.Printer.Print(printer.Result{
		OutputJSON: volume,
		KeyValue:   getVolumesKVMaps([]resources.Volume{*volume}),
		Columns:    getVolumesCols(builder.GetGlobalFlagName(c.ParentName, config.ArgCols), c.Printer.GetStderr()),
	})
}

func RunVolumeDetach(c *builder.CommandConfig) error {
	err := utils.AskForConfirm(c.Stdin, c.Printer, "detach volume")
	if err != nil {
		return err
	}
	resp, err := c.Volumes().Detach(
		viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgDataCenterId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgServerId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgVolumeId)),
	)
	if err != nil {
		return err
	}

	err = waitForAction(c, printer.GetRequestPath(resp))
	if err != nil {
		return err
	}
	return c.Printer.Print(printer.Result{
		ApiResponse: resp,
		Resource:    "volume",
		Verb:        "detach",
		WaitFlag:    viper.GetBool(builder.GetFlagName(c.ParentName, c.Name, config.ArgWait)),
	})
}

var defaultVolumeCols = []string{"VolumeId", "Name", "Size", "Type", "LicenseType", "State", "Image"}

type VolumePrint struct {
	VolumeId         string   `json:"VolumeId,omitempty"`
	Name             string   `json:"Name,omitempty"`
	Size             string   `json:"Size,omitempty"`
	Type             string   `json:"Type,omitempty"`
	LicenseType      string   `json:"LicenseType,omitempty"`
	Bus              string   `json:"Bus,omitempty"`
	AvailabilityZone string   `json:"AvailabilityZone,omitempty"`
	State            string   `json:"State,omitempty"`
	BackUpUnitId     string   `json:"BackUpUnitId,omitempty"`
	Image            string   `json:"Image,omitempty"`
	SshKeys          []string `json:"SshKeys,omitempty"`
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
		"LicenseType":      "LicenseType",
		"Bus":              "Bus",
		"AvailabilityZone": "AvailabilityZone",
		"State":            "State",
		"BackUpUnitId":     "BackUpUnitId",
		"Image":            "Image",
		"SshKeys":          "SshKeys",
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
		if licenseType, ok := properties.GetLicenceTypeOk(); ok && licenseType != nil {
			volumePrint.LicenseType = *licenseType
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
			volumePrint.BackUpUnitId = *backUpUnitId
		}
		if image, ok := properties.GetImageOk(); ok && image != nil {
			volumePrint.Image = *image
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
	volumeSvc := resources.NewVolumeService(clientSvc.Get(), context.TODO())
	volumes, _, err := volumeSvc.ListAttached(datacenterId, serverId)
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
