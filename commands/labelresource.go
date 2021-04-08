package commands

import (
	"context"
	"errors"
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

func labelDatacenter() *builder.Command {
	ctx := context.TODO()
	labelDatacenterCmd := &builder.Command{
		Command: &cobra.Command{
			Use:              "datacenter",
			Short:            "Label Data Center Operations",
			Long:             `The sub-commands of ` + "`" + `ionosctl label datacenter` + "`" + ` allow you to create, get, list, delete a Label from a Data Center`,
			TraverseChildren: true,
		},
	}
	globalFlags := labelDatacenterCmd.Command.PersistentFlags()
	globalFlags.StringP(config.ArgDataCenterId, "", "", config.RequiredFlagDatacenterId)
	_ = viper.BindPFlag(builder.GetGlobalFlagName(labelDatacenterCmd.Command.Name(), config.ArgDataCenterId), globalFlags.Lookup(config.ArgDataCenterId))
	_ = labelDatacenterCmd.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	globalFlags.StringSlice(config.ArgCols, defaultLabelResourceCols, "Columns to be printed in the standard output")
	_ = viper.BindPFlag(builder.GetGlobalFlagName(labelDatacenterCmd.Command.Name(), config.ArgCols), globalFlags.Lookup(config.ArgCols))

	/*
		List Command
	*/
	builder.NewCommand(ctx, labelDatacenterCmd, PreRunGlobalDcIdValidate, RunLabelDataCenterList, "list", "List Labels from a Data Center",
		"Use this command to list all Labels from a specified Data Center.\n\nRequired values to run command:\n\n* Data Center Id", listDatacenterExample, true)

	/*
		Get Command
	*/
	get := builder.NewCommand(ctx, labelDatacenterCmd, PreRunGlobalDcIdLabelKeyValidate, RunLabelDataCenterGet, "get", "Get a Label from a Data Center",
		"Use this command to get information about a specified Label from a Data Center.\n\nRequired values to run command:\n\n* Data Center Id\n* Label Key", getDatacenterExample, true)
	get.AddStringFlag(config.ArgLabelKey, "", "", config.RequiredFlagLabelKey)
	_ = get.Command.RegisterFlagCompletionFunc(config.ArgLabelKey, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getLabelResourceKeys(os.Stderr, viper.GetString(builder.GetGlobalFlagName(labelDatacenterCmd.Command.Name(), config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Create Command
	*/
	create := builder.NewCommand(ctx, labelDatacenterCmd, PreRunGlobalDcIdLabelKeyValidate, RunLabelDataCenterCreate, "create", "Create a Label on a Data Center",
		`Use this command to create a Label on Data Center. You must specify the key and the value for the Label.

You can wait for the action to be executed using `+"`"+`--wait`+"`"+` option.

Required values to run command: 

* Data Center Id 
* Label Key`, createDatacenterExample, true)
	create.AddStringFlag(config.ArgLabelKey, "", "", config.RequiredFlagLabelKey)
	create.AddStringFlag(config.ArgLabelValue, "", "", "Value of the Label in the Data Center. If not set, it will take the value of the key.")

	/*
		Delete Command
	*/
	deleteCmd := builder.NewCommand(ctx, labelDatacenterCmd, PreRunGlobalDcIdLabelKeyValidate, RunLabelDataCenterDelete, "delete", "Delete a Label from a Data Center",
		`Use this command to delete a specified Label from a Data Center.

You can wait for the action to be executed using `+"`"+`--wait`+"`"+` option. You can force the command to execute without user input using `+"`"+`--ignore-stdin`+"`"+` option.

Required values to run command:

* Data Center Id
* Label Key`, deleteDatacenterExample, true)
	deleteCmd.AddStringFlag(config.ArgLabelKey, "", "", config.RequiredFlagLabelKey)
	return labelDatacenterCmd
}

func PreRunGlobalDcIdLabelKeyValidate(c *builder.PreCommandConfig) error {
	err := builder.CheckRequiredGlobalFlags(c.ParentName, config.ArgDataCenterId)
	if err != nil {
		return err
	}
	return builder.CheckRequiredFlags(c.ParentName, c.Name, config.ArgLabelKey)
}

func RunLabelDataCenterList(c *builder.CommandConfig) error {
	labelDcs, _, err := c.Labels().DatacenterList(viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgDataCenterId)))
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(nil, c, getLabelResources(labelDcs)))
}

func RunLabelDataCenterGet(c *builder.CommandConfig) error {
	labelDc, _, err := c.Labels().DatacenterGet(
		viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgDataCenterId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgLabelKey)))
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(nil, c, getLabelResource(labelDc)))
}

func RunLabelDataCenterCreate(c *builder.CommandConfig) error {
	labelDc, _, err := c.Labels().DatacenterCreate(
		viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgDataCenterId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgLabelKey)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgLabelKey)),
	)
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(nil, c, getLabelResource(labelDc)))
}

func RunLabelDataCenterDelete(c *builder.CommandConfig) error {
	err := utils.AskForConfirm(c.Stdin, c.Printer, "delete label from data center")
	if err != nil {
		return err
	}
	_, err = c.Labels().DatacenterDelete(
		viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgDataCenterId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgLabelKey)),
	)
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(nil, c, nil))
}

func labelServer() *builder.Command {
	ctx := context.TODO()
	labelServerCmd := &builder.Command{
		Command: &cobra.Command{
			Use:              "server",
			Short:            "Label Server Operations",
			Long:             `The sub-commands of ` + "`" + `ionosctl label server` + "`" + ` allow you to create, get, list, delete a Label from a Server`,
			TraverseChildren: true,
		},
	}
	globalFlags := labelServerCmd.Command.PersistentFlags()
	globalFlags.StringP(config.ArgDataCenterId, "", "", config.RequiredFlagDatacenterId)
	_ = viper.BindPFlag(builder.GetGlobalFlagName(labelServerCmd.Command.Name(), config.ArgDataCenterId), globalFlags.Lookup(config.ArgDataCenterId))
	_ = labelServerCmd.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	globalFlags.StringP(config.ArgServerId, "", "", config.RequiredFlagServerId)
	_ = viper.BindPFlag(builder.GetGlobalFlagName(labelServerCmd.Command.Name(), config.ArgServerId), globalFlags.Lookup(config.ArgServerId))
	_ = labelServerCmd.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getServersIds(os.Stderr, viper.GetString(builder.GetGlobalFlagName(labelServerCmd.Command.Name(), config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	globalFlags.StringSlice(config.ArgCols, defaultLabelResourceCols, "Columns to be printed in the standard output")
	_ = viper.BindPFlag(builder.GetGlobalFlagName(labelServerCmd.Command.Name(), config.ArgCols), globalFlags.Lookup(config.ArgCols))

	/*
		List Command
	*/
	builder.NewCommand(ctx, labelServerCmd, PreRunGlobalDcServerIdsValidate, RunLabelServerList, "list", "List Labels from a Server",
		"Use this command to list all Labels from a specified Server.\n\nRequired values to run command:\n\n* Data Center Id\n* Server Id", listServerExample, true)

	/*
		Get Command
	*/
	get := builder.NewCommand(ctx, labelServerCmd, PreRunGlobalDcServerIdLabelKeyValidate, RunLabelServerGet, "get", "Get a Label from a Server",
		"Use this command to get information about a specified Label from a Server.\n\nRequired values to run command:\n\n* Data Center Id\n* Server Id\n* Label Key", getServerExample, true)
	get.AddStringFlag(config.ArgLabelKey, "", "", config.RequiredFlagLabelKey)
	_ = get.Command.RegisterFlagCompletionFunc(config.ArgLabelKey, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getLabelResourceKeys(os.Stderr, viper.GetString(builder.GetGlobalFlagName(labelServerCmd.Command.Name(), config.ArgServerId))), cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Create Command
	*/
	create := builder.NewCommand(ctx, labelServerCmd, PreRunGlobalDcServerIdLabelKeyValidate, RunLabelServerCreate, "create", "Create a Label on a Server",
		`Use this command to create a Label on Server. You must specify the key and the value for the Label.

You can wait for the action to be executed using `+"`"+`--wait`+"`"+` option.

Required values to run command: 

* Data Center Id
* Server Id 
* Label Key`, createServerExample, true)
	create.AddStringFlag(config.ArgLabelKey, "", "", config.RequiredFlagLabelKey)
	create.AddStringFlag(config.ArgLabelValue, "", "", "Value of the Label in the Server. If not set, it will take the value of the key.")

	/*
		Delete Command
	*/
	deleteCmd := builder.NewCommand(ctx, labelServerCmd, PreRunGlobalDcServerIdLabelKeyValidate, RunLabelServerDelete, "delete", "Delete a Label from a Server",
		`Use this command to delete a specified Label from a Server.

You can wait for the action to be executed using `+"`"+`--wait`+"`"+` option. You can force the command to execute without user input using `+"`"+`--ignore-stdin`+"`"+` option.

Required values to run command:

* Data Center Id
* Server Id
* Label Key`, deleteServerExample, true)
	deleteCmd.AddStringFlag(config.ArgLabelKey, "", "", config.RequiredFlagLabelKey)
	return labelServerCmd
}

func PreRunGlobalDcServerIdLabelKeyValidate(c *builder.PreCommandConfig) error {
	err := builder.CheckRequiredGlobalFlags(c.ParentName, config.ArgDataCenterId, config.ArgServerId)
	if err != nil {
		return err
	}
	return builder.CheckRequiredFlags(c.ParentName, c.Name, config.ArgLabelKey)
}

func RunLabelServerList(c *builder.CommandConfig) error {
	labelDcs, _, err := c.Labels().ServerList(
		viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgDataCenterId)),
		viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgDataCenterId)),
	)
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(nil, c, getLabelResources(labelDcs)))
}

func RunLabelServerGet(c *builder.CommandConfig) error {
	labelDc, _, err := c.Labels().ServerGet(
		viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgDataCenterId)),
		viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgServerId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgLabelKey)))
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(nil, c, getLabelResource(labelDc)))
}

func RunLabelServerCreate(c *builder.CommandConfig) error {
	labelDc, _, err := c.Labels().ServerCreate(
		viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgDataCenterId)),
		viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgServerId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgLabelKey)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgLabelKey)),
	)
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(nil, c, getLabelResource(labelDc)))
}

func RunLabelServerDelete(c *builder.CommandConfig) error {
	err := utils.AskForConfirm(c.Stdin, c.Printer, "delete label from server")
	if err != nil {
		return err
	}
	_, err = c.Labels().ServerDelete(
		viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgDataCenterId)),
		viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgServerId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgLabelKey)),
	)
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(nil, c, nil))
}

func labelVolume() *builder.Command {
	ctx := context.TODO()
	labelVolumeCmd := &builder.Command{
		Command: &cobra.Command{
			Use:              "volume",
			Short:            "Label Volume Operations",
			Long:             `The sub-commands of ` + "`" + `ionosctl label volume` + "`" + ` allow you to create, get, list, delete a Label from a Volume`,
			TraverseChildren: true,
		},
	}
	globalFlags := labelVolumeCmd.Command.PersistentFlags()
	globalFlags.StringP(config.ArgDataCenterId, "", "", config.RequiredFlagDatacenterId)
	_ = viper.BindPFlag(builder.GetGlobalFlagName(labelVolumeCmd.Command.Name(), config.ArgDataCenterId), globalFlags.Lookup(config.ArgDataCenterId))
	_ = labelVolumeCmd.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	globalFlags.StringP(config.ArgVolumeId, "", "", config.RequiredFlagVolumeId)
	_ = viper.BindPFlag(builder.GetGlobalFlagName(labelVolumeCmd.Command.Name(), config.ArgVolumeId), globalFlags.Lookup(config.ArgVolumeId))
	_ = labelVolumeCmd.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getVolumesIds(os.Stderr, viper.GetString(builder.GetGlobalFlagName(labelVolumeCmd.Command.Name(), config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	globalFlags.StringSlice(config.ArgCols, defaultLabelResourceCols, "Columns to be printed in the standard output")
	_ = viper.BindPFlag(builder.GetGlobalFlagName(labelVolumeCmd.Command.Name(), config.ArgCols), globalFlags.Lookup(config.ArgCols))

	/*
		List Command
	*/
	builder.NewCommand(ctx, labelVolumeCmd, PreRunGlobalDcVolumeIdsValidate, RunLabelVolumeList, "list", "List Labels from a Volume",
		"Use this command to list all Labels from a specified Volume.\n\nRequired values to run command:\n\n* Data Center Id\n* Volume Id", listVolumeExample, true)

	/*
		Get Command
	*/
	get := builder.NewCommand(ctx, labelVolumeCmd, PreRunGlobalDcVolumeIdLabelKeyValidate, RunLabelVolumeGet, "get", "Get a Label from a Volume",
		"Use this command to get information about a specified Label from a Volume.\n\nRequired values to run command:\n\n* Data Center Id\n* Volume Id\n* Label Key", getVolumeExample, true)
	get.AddStringFlag(config.ArgLabelKey, "", "", config.RequiredFlagLabelKey)
	_ = get.Command.RegisterFlagCompletionFunc(config.ArgLabelKey, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getLabelResourceKeys(os.Stderr, viper.GetString(builder.GetGlobalFlagName(labelVolumeCmd.Command.Name(), config.ArgVolumeId))), cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Create Command
	*/
	create := builder.NewCommand(ctx, labelVolumeCmd, PreRunGlobalDcVolumeIdLabelKeyValidate, RunLabelVolumeCreate, "create", "Create a Label on a Volume",
		`Use this command to create a Label on Volume. You must specify the key and the value for the Label.

You can wait for the action to be executed using `+"`"+`--wait`+"`"+` option.

Required values to run command:

* Data Center Id
* Volume Id
* Label Key`, createVolumeExample, true)
	create.AddStringFlag(config.ArgLabelKey, "", "", config.RequiredFlagLabelKey)
	create.AddStringFlag(config.ArgLabelValue, "", "", "Value of the Label in the Volume. If not set, it will take the value of the key.")

	/*
		Delete Command
	*/
	deleteCmd := builder.NewCommand(ctx, labelVolumeCmd, PreRunGlobalDcVolumeIdLabelKeyValidate, RunLabelVolumeDelete, "delete", "Delete a Label from a Volume",
		`Use this command to delete a specified Label from a Volume.

You can wait for the action to be executed using `+"`"+`--wait`+"`"+` option. You can force the command to execute without user input using `+"`"+`--ignore-stdin`+"`"+` option.

Required values to run command:

* Data Center Id
* Volume Id
* Label Key`, deleteVolumeExample, true)
	deleteCmd.AddStringFlag(config.ArgLabelKey, "", "", config.RequiredFlagLabelKey)
	return labelVolumeCmd
}

func PreRunGlobalDcVolumeIdsValidate(c *builder.PreCommandConfig) error {
	return builder.CheckRequiredGlobalFlags(c.ParentName, config.ArgDataCenterId, config.ArgVolumeId)
}

func PreRunGlobalDcVolumeIdLabelKeyValidate(c *builder.PreCommandConfig) error {
	err := builder.CheckRequiredGlobalFlags(c.ParentName, config.ArgDataCenterId, config.ArgVolumeId)
	if err != nil {
		return err
	}
	return builder.CheckRequiredFlags(c.ParentName, c.Name, config.ArgLabelKey)
}

func RunLabelVolumeList(c *builder.CommandConfig) error {
	labelDcs, _, err := c.Labels().VolumeList(
		viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgDataCenterId)),
		viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgDataCenterId)),
	)
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(nil, c, getLabelResources(labelDcs)))
}

func RunLabelVolumeGet(c *builder.CommandConfig) error {
	labelDc, _, err := c.Labels().VolumeGet(
		viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgDataCenterId)),
		viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgVolumeId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgLabelKey)))
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(nil, c, getLabelResource(labelDc)))
}

func RunLabelVolumeCreate(c *builder.CommandConfig) error {
	labelDc, _, err := c.Labels().VolumeCreate(
		viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgDataCenterId)),
		viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgVolumeId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgLabelKey)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgLabelKey)),
	)
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(nil, c, getLabelResource(labelDc)))
}

func RunLabelVolumeDelete(c *builder.CommandConfig) error {
	err := utils.AskForConfirm(c.Stdin, c.Printer, "delete label from volume")
	if err != nil {
		return err
	}
	_, err = c.Labels().VolumeDelete(
		viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgDataCenterId)),
		viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgVolumeId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgLabelKey)),
	)
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(nil, c, nil))
}

func labelIpBlock() *builder.Command {
	ctx := context.TODO()
	labelIpBlockCmd := &builder.Command{
		Command: &cobra.Command{
			Use:              "ipblock",
			Short:            "Label IPBlock Operations",
			Long:             `The sub-commands of ` + "`" + `ionosctl label ipblock` + "`" + ` allow you to create, get, list, delete a Label from a IPBlock`,
			TraverseChildren: true,
		},
	}
	globalFlags := labelIpBlockCmd.Command.PersistentFlags()
	globalFlags.StringP(config.ArgIpBlockId, "", "", config.RequiredFlagIpBlockId)
	_ = viper.BindPFlag(builder.GetGlobalFlagName(labelIpBlockCmd.Command.Name(), config.ArgIpBlockId), globalFlags.Lookup(config.ArgIpBlockId))
	_ = labelIpBlockCmd.Command.RegisterFlagCompletionFunc(config.ArgIpBlockId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getIpBlocksIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	globalFlags.StringSlice(config.ArgCols, defaultLabelResourceCols, "Columns to be printed in the standard output")
	_ = viper.BindPFlag(builder.GetGlobalFlagName(labelIpBlockCmd.Command.Name(), config.ArgCols), globalFlags.Lookup(config.ArgCols))

	/*
		List Command
	*/
	builder.NewCommand(ctx, labelIpBlockCmd, PreRunGlobalIpBlockIdValidate, RunLabelIpBlockList, "list", "List Labels from a IPBlock",
		"Use this command to list all Labels from a specified IPBlock.\n\nRequired values to run command:\n\n* IPBlock Id", listIpBlockExample, true)

	/*
		Get Command
	*/
	get := builder.NewCommand(ctx, labelIpBlockCmd, PreRunGlobalIpBlockIdLabelKeyValidate, RunLabelIpBlockGet, "get", "Get a Label from a IPBlock",
		"Use this command to get information about a specified Label from a IPBlock.\n\nRequired values to run command:\n\n* IPBlock Id\n* Label Key", getIpBlockExample, true)
	get.AddStringFlag(config.ArgLabelKey, "", "", config.RequiredFlagLabelKey)
	_ = get.Command.RegisterFlagCompletionFunc(config.ArgLabelKey, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getLabelResourceKeys(os.Stderr, viper.GetString(builder.GetGlobalFlagName(labelIpBlockCmd.Command.Name(), config.ArgIpBlockId))), cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Create Command
	*/
	create := builder.NewCommand(ctx, labelIpBlockCmd, PreRunGlobalIpBlockIdLabelKeyValidate, RunLabelIpBlockCreate, "create", "Create a Label on a IPBlock",
		`Use this command to create a Label on IPBlock. You must specify the key and the value for the Label.

You can wait for the action to be executed using `+"`"+`--wait`+"`"+` option.

Required values to run command: 

* IPBlock Id 
* Label Key`, createIpBlockExample, true)
	create.AddStringFlag(config.ArgLabelKey, "", "", config.RequiredFlagLabelKey)
	create.AddStringFlag(config.ArgLabelValue, "", "", "Value of the Label in the IPBlock. If not set, it will take the value of the key.")

	/*
		Delete Command
	*/
	deleteCmd := builder.NewCommand(ctx, labelIpBlockCmd, PreRunGlobalIpBlockIdLabelKeyValidate, RunLabelIpBlockDelete, "delete", "Delete a Label from a IPBlock",
		`Use this command to delete a specified Label from a IPBlock.

You can wait for the action to be executed using `+"`"+`--wait`+"`"+` option. You can force the command to execute without user input using `+"`"+`--ignore-stdin`+"`"+` option.

Required values to run command:

* IPBlock Id
* Label Key`, deleteIpBlockExample, true)
	deleteCmd.AddStringFlag(config.ArgLabelKey, "", "", config.RequiredFlagLabelKey)
	return labelIpBlockCmd
}

func PreRunGlobalIpBlockIdValidate(c *builder.PreCommandConfig) error {
	return builder.CheckRequiredGlobalFlags(c.ParentName, config.ArgIpBlockId)
}

func PreRunGlobalIpBlockIdLabelKeyValidate(c *builder.PreCommandConfig) error {
	err := builder.CheckRequiredGlobalFlags(c.ParentName, config.ArgIpBlockId)
	if err != nil {
		return err
	}
	return builder.CheckRequiredFlags(c.ParentName, c.Name, config.ArgLabelKey)
}

func RunLabelIpBlockList(c *builder.CommandConfig) error {
	labelDcs, _, err := c.Labels().IpBlockList(viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgIpBlockId)))
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(nil, c, getLabelResources(labelDcs)))
}

func RunLabelIpBlockGet(c *builder.CommandConfig) error {
	labelDc, _, err := c.Labels().IpBlockGet(
		viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgIpBlockId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgLabelKey)))
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(nil, c, getLabelResource(labelDc)))
}

func RunLabelIpBlockCreate(c *builder.CommandConfig) error {
	labelDc, _, err := c.Labels().IpBlockCreate(
		viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgIpBlockId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgLabelKey)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgLabelKey)),
	)
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(nil, c, getLabelResource(labelDc)))
}

func RunLabelIpBlockDelete(c *builder.CommandConfig) error {
	err := utils.AskForConfirm(c.Stdin, c.Printer, "delete label from ip block")
	if err != nil {
		return err
	}
	_, err = c.Labels().IpBlockDelete(
		viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgIpBlockId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgLabelKey)),
	)
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(nil, c, nil))
}

func labelSnapshot() *builder.Command {
	ctx := context.TODO()
	labelSnapshotCmd := &builder.Command{
		Command: &cobra.Command{
			Use:              "snapshot",
			Short:            "Label Snapshot Operations",
			Long:             `The sub-commands of ` + "`" + `ionosctl label snapshot` + "`" + ` allow you to create, get, list, delete a Label from a Snapshot`,
			TraverseChildren: true,
		},
	}
	globalFlags := labelSnapshotCmd.Command.PersistentFlags()
	globalFlags.StringP(config.ArgSnapshotId, "", "", config.RequiredFlagSnapshotId)
	_ = viper.BindPFlag(builder.GetGlobalFlagName(labelSnapshotCmd.Command.Name(), config.ArgSnapshotId), globalFlags.Lookup(config.ArgSnapshotId))
	_ = labelSnapshotCmd.Command.RegisterFlagCompletionFunc(config.ArgSnapshotId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getSnapshotIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	globalFlags.StringSlice(config.ArgCols, defaultLabelResourceCols, "Columns to be printed in the standard output")
	_ = viper.BindPFlag(builder.GetGlobalFlagName(labelSnapshotCmd.Command.Name(), config.ArgCols), globalFlags.Lookup(config.ArgCols))

	/*
		List Command
	*/
	builder.NewCommand(ctx, labelSnapshotCmd, PreRunGlobalSnapshotIdValidate, RunLabelSnapshotList, "list", "List Labels from a Snapshot",
		"Use this command to list all Labels from a specified Snapshot.\n\nRequired values to run command:\n\n* Snapshot Id", "", true)

	/*
		Get Command
	*/
	get := builder.NewCommand(ctx, labelSnapshotCmd, PreRunGlobalSnapshotIdLabelKeyValidate, RunLabelSnapshotGet, "get", "Get a Label from a Snapshot",
		"Use this command to get information about a specified Label from a Snapshot.\n\nRequired values to run command:\n\n* Snapshot Id\n* Label Key", getSnapshotExample, true)
	get.AddStringFlag(config.ArgLabelKey, "", "", config.RequiredFlagLabelKey)
	_ = get.Command.RegisterFlagCompletionFunc(config.ArgLabelKey, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getLabelResourceKeys(os.Stderr, viper.GetString(builder.GetGlobalFlagName(labelSnapshotCmd.Command.Name(), config.ArgSnapshotId))), cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Create Command
	*/
	create := builder.NewCommand(ctx, labelSnapshotCmd, PreRunGlobalSnapshotIdLabelKeyValidate, RunLabelSnapshotCreate, "create", "Create a Label on a Snapshot",
		`Use this command to create a Label on Snapshot. You must specify the key and the value for the Label.

You can wait for the action to be executed using `+"`"+`--wait`+"`"+` option.

Required values to run command: 

* Snapshot Id 
* Label Key`, createSnapshotExample, true)
	create.AddStringFlag(config.ArgLabelKey, "", "", config.RequiredFlagLabelKey)
	create.AddStringFlag(config.ArgLabelValue, "", "", "Value of the Label in the Snapshot. If not set, it will take the value of the key.")

	/*
		Delete Command
	*/
	deleteCmd := builder.NewCommand(ctx, labelSnapshotCmd, PreRunGlobalSnapshotIdLabelKeyValidate, RunLabelSnapshotDelete, "delete", "Delete a Label from a Snapshot",
		`Use this command to delete a specified Label from a Snapshot.

You can wait for the action to be executed using `+"`"+`--wait`+"`"+` option. You can force the command to execute without user input using `+"`"+`--ignore-stdin`+"`"+` option.

Required values to run command:

* Snapshot Id
* Label Key`, deleteSnapshotExample, true)
	deleteCmd.AddStringFlag(config.ArgLabelKey, "", "", config.RequiredFlagLabelKey)
	return labelSnapshotCmd
}

func PreRunGlobalSnapshotIdValidate(c *builder.PreCommandConfig) error {
	return builder.CheckRequiredGlobalFlags(c.ParentName, config.ArgSnapshotId)
}

func PreRunGlobalSnapshotIdLabelKeyValidate(c *builder.PreCommandConfig) error {
	err := builder.CheckRequiredGlobalFlags(c.ParentName, config.ArgSnapshotId)
	if err != nil {
		return err
	}
	return builder.CheckRequiredFlags(c.ParentName, c.Name, config.ArgLabelKey)
}

func RunLabelSnapshotList(c *builder.CommandConfig) error {
	labelDcs, _, err := c.Labels().SnapshotList(viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgSnapshotId)))
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(nil, c, getLabelResources(labelDcs)))
}

func RunLabelSnapshotGet(c *builder.CommandConfig) error {
	labelDc, _, err := c.Labels().SnapshotGet(
		viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgSnapshotId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgLabelKey)))
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(nil, c, getLabelResource(labelDc)))
}

func RunLabelSnapshotCreate(c *builder.CommandConfig) error {
	labelDc, _, err := c.Labels().SnapshotCreate(
		viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgSnapshotId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgLabelKey)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgLabelKey)),
	)
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(nil, c, getLabelResource(labelDc)))
}

func RunLabelSnapshotDelete(c *builder.CommandConfig) error {
	err := utils.AskForConfirm(c.Stdin, c.Printer, "delete label from ip block")
	if err != nil {
		return err
	}
	_, err = c.Labels().SnapshotDelete(
		viper.GetString(builder.GetGlobalFlagName(c.ParentName, config.ArgSnapshotId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgLabelKey)),
	)
	if err != nil {
		return err
	}
	return c.Printer.Print(getLabelResourcePrint(nil, c, nil))
}

// Output Printing

var defaultLabelResourceCols = []string{"Key", "Value"}

type LabelResourcePrint struct {
	Key   string `json:"Key,omitempty"`
	Value string `json:"Value,omitempty"`
}

func getLabelResourcePrint(resp *resources.Response, c *builder.CommandConfig, s []resources.LabelResource) printer.Result {
	r := printer.Result{}
	if c != nil {
		if resp != nil {
			r.ApiResponse = resp
			r.Resource = "label " + c.ParentName
			r.Verb = c.Name
		}
		if s != nil {
			r.OutputJSON = s
			r.KeyValue = getLabelResourcesKVMaps(s)
			r.Columns = getLabelResourceCols(builder.GetGlobalFlagName(c.ParentName, config.ArgCols), c.Printer.GetStderr())
		}
	}
	return r
}

func getLabelResourceCols(flagName string, outErr io.Writer) []string {
	if viper.IsSet(flagName) {
		var labelCols []string
		columnsMap := map[string]string{
			"Key":   "Key",
			"Value": "Value",
		}
		for _, k := range viper.GetStringSlice(flagName) {
			col := columnsMap[k]
			if col != "" {
				labelCols = append(labelCols, col)
			} else {
				clierror.CheckError(errors.New("unknown column "+k), outErr)
			}
		}
		return labelCols
	} else {
		return defaultLabelResourceCols
	}
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

func getLabelResourceKeys(outErr io.Writer, datacenterId string) []string {
	err := config.LoadFile()
	clierror.CheckError(err, outErr)
	clientSvc, err := resources.NewClientService(
		viper.GetString(config.Username),
		viper.GetString(config.Password),
		viper.GetString(config.Token),
		viper.GetString(config.ArgServerUrl),
	)
	clierror.CheckError(err, outErr)
	labelResourceSvc := resources.NewLabelResourceService(clientSvc.Get(), context.TODO())
	labelResources, _, err := labelResourceSvc.DatacenterList(datacenterId)
	clierror.CheckError(err, outErr)
	labelsDc := make([]string, 0)
	if items, ok := labelResources.LabelResources.GetItemsOk(); ok && items != nil {
		for _, item := range *items {
			if properties, ok := item.GetPropertiesOk(); ok && properties != nil {
				if key, ok := properties.GetKeyOk(); ok && key != nil {
					labelsDc = append(labelsDc, *key)
				}
			}
		}
	} else {
		return nil
	}
	return labelsDc
}
