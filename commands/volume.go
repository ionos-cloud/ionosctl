package commands

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/fatih/structs"
	"github.com/hashicorp/go-multierror"
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
	ctx := context.TODO()
	volumeCmd := &builder.Command{
		Command: &cobra.Command{
			Use:              "volume",
			Short:            "Volume Operations",
			Long:             `The sub-commands of ` + "`" + `ionosctl volume` + "`" + ` manage your block storage volumes by creating, updating, getting specific information, deleting Volumes. To attach a Volume to a Server, use the Server command ` + "`" + `ionosctl server attach-volume` + "`" + `.`,
			TraverseChildren: true,
		},
	}
	globalFlags := volumeCmd.GlobalFlags()
	globalFlags.StringP(config.ArgDataCenterId, "", "", config.RequiredFlagDatacenterId)
	_ = viper.BindPFlag(builder.GetGlobalFlagName(volumeCmd.Name(), config.ArgDataCenterId), globalFlags.Lookup(config.ArgDataCenterId))
	_ = volumeCmd.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	globalFlags.StringSlice(config.ArgCols, defaultVolumeCols, "Columns to be printed in the standard output")
	_ = viper.BindPFlag(builder.GetGlobalFlagName(volumeCmd.Name(), config.ArgCols), globalFlags.Lookup(config.ArgCols))

	/*
		List Command
	*/
	builder.NewCommand(ctx, volumeCmd, PreRunGlobalDcId, RunVolumeList, "list", "List Volumes",
		"Use this command to list all Volumes from a Data Center on your account.\n\nRequired values to run command:\n\n* Data Center Id",
		listVolumeExample, true)

	/*
		Get Command
	*/
	get := builder.NewCommand(ctx, volumeCmd, PreRunGlobalDcIdVolumeId, RunVolumeGet, "get", "Get a Volume",
		"Use this command to retrieve information about a Volume using its ID.\n\nRequired values to run command:\n\n* Data Center Id\n* Volume Id",
		getVolumeExample, true)
	get.AddStringFlag(config.ArgVolumeId, "", "", config.RequiredFlagVolumeId)
	_ = get.Command.RegisterFlagCompletionFunc(config.ArgVolumeId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getVolumesIds(os.Stderr, viper.GetString(builder.GetGlobalFlagName(volumeCmd.Name(), config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Create Command
	*/
	create := builder.NewCommand(ctx, volumeCmd, PreRunGlobalDcId, RunVolumeCreate, "create", "Create a Volume",
		`Use this command to create a Volume on your account. Creates a volume within the data center. This will NOT attach the Volume to a Server. Please see the Servers commands for details on how to attach storage Volumes. You can specify the name, size, type, licence type and availability zone for the object.

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
	update := builder.NewCommand(ctx, volumeCmd, PreRunGlobalDcIdVolumeId, RunVolumeUpdate, "update", "Update a Volume",
		`Use this command to update a Volume. You may increase the size of an existing storage Volume. You cannot reduce the size of an existing storage Volume. The Volume size will be increased without reboot if the appropriate "hot plug" settings have been set to true. The additional capacity is not added to any partition therefore you will need to adjust the partition inside the operating system afterwards.

Once you have increased the Volume size you cannot decrease the Volume size using the Cloud API. Certain attributes can only be set when a Volume is created and are considered immutable once the Volume has been provisioned.

You can wait for the action to be executed using `+"`"+`--wait`+"`"+` option.

Required values to run command:

* Data Center Id
* Volume Id`, updateVolumeExample, true)
	update.AddStringFlag(config.ArgVolumeId, "", "", config.RequiredFlagVolumeId)
	_ = update.Command.RegisterFlagCompletionFunc(config.ArgVolumeId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getVolumesIds(os.Stderr, viper.GetString(builder.GetGlobalFlagName(volumeCmd.Name(), config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
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
	deleteCmd := builder.NewCommand(ctx, volumeCmd, PreRunGlobalDcIdVolumeId, RunVolumeDelete, "delete", "Delete a Volume",
		`Use this command to delete specified Volume. This will result in the Volume being removed from your Virtual Data Center. Please use this with caution!

You can wait for the action to be executed using `+"`"+`--wait`+"`"+` option. You can force the command to execute without user input using `+"`"+`--force`+"`"+` option.

Required values to run command:

* Data Center Id
* Volume Id`, deleteVolumeExample, true)
	deleteCmd.AddStringFlag(config.ArgVolumeId, "", "", config.RequiredFlagVolumeId)
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(config.ArgVolumeId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getVolumesIds(os.Stderr, viper.GetString(builder.GetGlobalFlagName(volumeCmd.Name(), config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddBoolFlag(config.ArgWait, "", config.DefaultWait, "Wait for Volume to be deleted")
	deleteCmd.AddIntFlag(config.ArgTimeout, "", config.DefaultTimeoutSeconds, "Timeout option for Volume to be deleted [seconds]")

	return volumeCmd
}

func PreRunGlobalDcIdVolumeId(c *builder.PreCommandConfig) error {
	var result *multierror.Error
	if err := builder.CheckRequiredGlobalFlags(c.ParentName, config.ArgDataCenterId); err != nil {
		result = multierror.Append(result, err)
	}
	if err := builder.CheckRequiredFlags(c.ParentName, c.Name, config.ArgVolumeId); err != nil {
		result = multierror.Append(result, err)
	}
	if result != nil {
		return result
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

	if err = waitForAction(c, printer.GetRequestPath(resp)); err != nil {
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

	if err = waitForAction(c, printer.GetRequestPath(resp)); err != nil {
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
	if err := utils.AskForConfirm(c.Stdin, c.Printer, "delete volume"); err != nil {
		return err
	}
	resp, err := c.Volumes().Delete(
		viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgDataCenterId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgVolumeId)),
	)
	if err != nil {
		return err
	}

	if err = waitForAction(c, printer.GetRequestPath(resp)); err != nil {
		return err
	}
	return c.Printer.Print(printer.Result{
		ApiResponse: resp,
		Resource:    "volume",
		Verb:        "delete",
		WaitFlag:    viper.GetBool(builder.GetFlagName(c.ParentName, c.Name, config.ArgWait)),
	})
}

// Server Volume Commands

func serverVolume() *builder.Command {
	ctx := context.TODO()
	serverVolumeCmd := &builder.Command{
		Command: &cobra.Command{
			Use:              "volume",
			Short:            "Server Volume Operations",
			Long:             `The sub-commands of ` + "`" + `ionosctl server volume` + "`" + ` allow you to attach, get, list, detach Volumes from Servers.`,
			TraverseChildren: true,
		},
	}
	globalFlags := serverVolumeCmd.GlobalFlags()
	globalFlags.StringSlice(config.ArgCols, defaultVolumeCols, "Columns to be printed in the standard output")
	_ = viper.BindPFlag(builder.GetGlobalFlagName(serverVolumeCmd.Name(), config.ArgCols), globalFlags.Lookup(config.ArgCols))

	/*
		Attach Volume Command
	*/
	attachVolume := builder.NewCommand(ctx, serverVolumeCmd, PreRunDcServerVolumeIds, RunServerVolumeAttach, "attach", "Attach a Volume to a Server",
		`Use this command to attach a pre-existing Volume to a Server.

You can wait for the action to be executed using `+"`"+`--wait`+"`"+` option.

Required values to run command:

* Data Center Id
* Server Id
* Volume Id`, attachVolumeServerExample, true)
	attachVolume.AddStringFlag(config.ArgDataCenterId, "", "", config.RequiredFlagDatacenterId)
	_ = attachVolume.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	attachVolume.AddStringFlag(config.ArgVolumeId, "", "", config.RequiredFlagVolumeId)
	_ = attachVolume.Command.RegisterFlagCompletionFunc(config.ArgVolumeId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getVolumesIds(os.Stderr, viper.GetString(builder.GetFlagName(serverVolumeCmd.Name(), attachVolume.Name(), config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	attachVolume.AddStringFlag(config.ArgServerId, "", "", config.RequiredFlagServerId)
	_ = attachVolume.Command.RegisterFlagCompletionFunc(config.ArgServerId, func(cmd *cobra.Command, ags []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getServersIds(os.Stderr, viper.GetString(builder.GetFlagName(serverVolumeCmd.Name(), attachVolume.Name(), config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	attachVolume.AddBoolFlag(config.ArgWait, "", config.DefaultWait, "Wait for Volume to be attached to Server")
	attachVolume.AddIntFlag(config.ArgTimeout, "", config.DefaultTimeoutSeconds, "Timeout option for Volume to be attached to a Server [seconds]")

	/*
		List Volumes Command
	*/
	listVolumes := builder.NewCommand(ctx, serverVolumeCmd, PreRunDcServerIds, RunServerVolumesList, "list", "List attached Volumes from a Server",
		"Use this command to retrieve a list of Volumes attached to the Server.\n\nRequired values to run command:\n\n* Data Center Id\n* Server Id",
		listVolumesServerExample, true)
	listVolumes.AddStringFlag(config.ArgDataCenterId, "", "", config.RequiredFlagDatacenterId)
	_ = listVolumes.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	listVolumes.AddStringFlag(config.ArgServerId, "", "", config.RequiredFlagServerId)
	_ = listVolumes.Command.RegisterFlagCompletionFunc(config.ArgServerId, func(cmd *cobra.Command, ags []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getServersIds(os.Stderr, viper.GetString(builder.GetFlagName(serverVolumeCmd.Name(), listVolumes.Name(), config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Describe Volume Command
	*/
	describeVolumeCmd := builder.NewCommand(ctx, serverVolumeCmd, PreRunDcServerVolumeIds, RunServerVolumeDescribe, "describe", "Retrieve an attached Volume from a Server",
		"Use this command to retrieve information about an attached Volume on Server.\n\nRequired values to run command:\n\n* Data Center Id\n* Server Id\n* Volume Id",
		describeVolumeServerExample, true)
	describeVolumeCmd.AddStringFlag(config.ArgDataCenterId, "", "", config.RequiredFlagDatacenterId)
	_ = describeVolumeCmd.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	describeVolumeCmd.AddStringFlag(config.ArgServerId, "", "", config.RequiredFlagServerId)
	_ = describeVolumeCmd.Command.RegisterFlagCompletionFunc(config.ArgServerId, func(cmd *cobra.Command, ags []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getServersIds(os.Stderr, viper.GetString(builder.GetFlagName(serverVolumeCmd.Name(), describeVolumeCmd.Name(), config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	describeVolumeCmd.AddStringFlag(config.ArgVolumeId, "", "", config.RequiredFlagVolumeId)
	_ = describeVolumeCmd.Command.RegisterFlagCompletionFunc(config.ArgVolumeId, func(cmd *cobra.Command, ags []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getAttachedVolumesIds(os.Stderr, viper.GetString(builder.GetFlagName(serverVolumeCmd.Name(), describeVolumeCmd.Name(), config.ArgDataCenterId)), viper.GetString(builder.GetFlagName(serverVolumeCmd.Name(), describeVolumeCmd.Name(), config.ArgServerId))), cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Detach Volume Command
	*/
	detachVolume := builder.NewCommand(ctx, serverVolumeCmd, PreRunDcServerVolumeIds, RunServerVolumeDetach, "detach", "Detach a Volume from a Server",
		`This will detach the Volume from the Server. Depending on the Volume HotUnplug settings, this may result in the Server being rebooted. This will NOT delete the Volume from your Virtual Data Center. You will need to use a separate command to delete a Volume.

You can wait for the action to be executed using `+"`"+`--wait`+"`"+` option. You can force the command to execute without user input using `+"`"+`--ignore-stdin`+"`"+` option.

Required values to run command:

* Data Center Id
* Server Id
* Volume Id`, detachVolumeServerExample, true)
	detachVolume.AddStringFlag(config.ArgDataCenterId, "", "", config.RequiredFlagDatacenterId)
	_ = detachVolume.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	detachVolume.AddStringFlag(config.ArgVolumeId, "", "", config.RequiredFlagVolumeId)
	_ = detachVolume.Command.RegisterFlagCompletionFunc(config.ArgVolumeId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getAttachedVolumesIds(os.Stderr, viper.GetString(builder.GetFlagName(serverVolumeCmd.Name(), detachVolume.Name(), config.ArgDataCenterId)), viper.GetString(builder.GetFlagName(serverVolumeCmd.Name(), detachVolume.Name(), config.ArgServerId))), cobra.ShellCompDirectiveNoFileComp
	})
	detachVolume.AddStringFlag(config.ArgServerId, "", "", config.RequiredFlagServerId)
	_ = detachVolume.Command.RegisterFlagCompletionFunc(config.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getServersIds(os.Stderr, viper.GetString(builder.GetFlagName(serverVolumeCmd.Name(), detachVolume.Name(), config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	detachVolume.AddBoolFlag(config.ArgWait, "", config.DefaultWait, "Wait for Volume to be detached from Server")
	detachVolume.AddIntFlag(config.ArgTimeout, "", config.DefaultTimeoutSeconds, "Timeout option for Volume to be detached from a Server [seconds]")

	return serverVolumeCmd
}

func PreRunDcServerVolumeIds(c *builder.PreCommandConfig) error {
	return builder.CheckRequiredFlags(c.ParentName, c.Name, config.ArgDataCenterId, config.ArgServerId, config.ArgVolumeId)
}

func RunServerVolumeAttach(c *builder.CommandConfig) error {
	attachedVol, resp, err := c.Servers().AttachVolume(
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgDataCenterId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgServerId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgVolumeId)),
	)
	if err != nil {
		return err
	}

	if err = waitForAction(c, printer.GetRequestPath(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getVolumePrint(resp, c, getVolume(attachedVol)))
}

func RunServerVolumesList(c *builder.CommandConfig) error {
	attachedVols, _, err := c.Servers().ListVolumes(
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgDataCenterId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgServerId)),
	)
	if err != nil {
		return err
	}
	return c.Printer.Print(getVolumePrint(nil, c, getAttachedVolumes(attachedVols)))
}

func RunServerVolumeDescribe(c *builder.CommandConfig) error {
	attachedVol, _, err := c.Servers().GetVolume(
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgDataCenterId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgServerId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgVolumeId)),
	)
	if err != nil {
		return err
	}
	return c.Printer.Print(getVolumePrint(nil, c, getVolume(attachedVol)))
}

func RunServerVolumeDetach(c *builder.CommandConfig) error {
	if err := utils.AskForConfirm(c.Stdin, c.Printer, "detach volume from server"); err != nil {
		return err
	}
	resp, err := c.Servers().DetachVolume(
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgDataCenterId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgServerId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgVolumeId)),
	)
	if err != nil {
		return err
	}

	if err = waitForAction(c, printer.GetRequestPath(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getVolumePrint(resp, c, nil))
}

// Output Printing

var defaultVolumeCols = []string{"VolumeId", "Name", "Size", "Type", "LicenceType", "State", "Image"}

type VolumePrint struct {
	VolumeId         string   `json:"VolumeId,omitempty"`
	Name             string   `json:"Name,omitempty"`
	Size             string   `json:"Size,omitempty"`
	Type             string   `json:"Type,omitempty"`
	LicenceType      string   `json:"LicenceType,omitempty"`
	Bus              string   `json:"Bus,omitempty"`
	AvailabilityZone string   `json:"AvailabilityZone,omitempty"`
	State            string   `json:"State,omitempty"`
	BackUpUnitId     string   `json:"BackUpUnitId,omitempty"`
	Image            string   `json:"Image,omitempty"`
	SshKeys          []string `json:"SshKeys,omitempty"`
}

func getVolumePrint(resp *resources.Response, c *builder.CommandConfig, vols []resources.Volume) printer.Result {
	r := printer.Result{}
	if c != nil {
		if resp != nil {
			r.ApiResponse = resp
			r.Resource = c.ParentName
			r.Verb = c.Name
			r.WaitFlag = viper.GetBool(builder.GetFlagName(c.ParentName, c.Name, config.ArgWait))
		}
		if vols != nil {
			r.OutputJSON = vols
			r.KeyValue = getVolumesKVMaps(vols)
			r.Columns = getVolumesCols(builder.GetGlobalFlagName(c.ParentName, config.ArgCols), c.Printer.GetStderr())
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
			volumePrint.BackUpUnitId = *backUpUnitId
		}
		if img, ok := properties.GetImageOk(); ok && img != nil {
			volumePrint.Image = *img
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
