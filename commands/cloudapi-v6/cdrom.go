package commands

import (
	"context"
	"github.com/ionos-cloud/ionosctl/commands/cloudapi-v6/waiter"
	cloudapi_v6 "github.com/ionos-cloud/ionosctl/services/cloudapi-v6"
	"io"
	"os"

	"github.com/ionos-cloud/ionosctl/internal/config"
	"github.com/ionos-cloud/ionosctl/internal/core"
	"github.com/ionos-cloud/ionosctl/internal/printer"
	"github.com/ionos-cloud/ionosctl/internal/utils"
	"github.com/ionos-cloud/ionosctl/internal/utils/clierror"
	"github.com/ionos-cloud/ionosctl/services/cloudapi-v6/resources"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Server Cdrom Commands

func ServerCdromCmd() *core.Command {
	ctx := context.TODO()
	serverCdromCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "cdrom",
			Aliases:          []string{"cd"},
			Short:            "Server CD-ROM Operations",
			Long:             "The sub-commands of `ionosctl server cdrom` allow you to attach, get, list, detach CD-ROMs from Servers.",
			TraverseChildren: true,
		},
	}
	globalFlags := serverCdromCmd.GlobalFlags()
	globalFlags.StringSliceP(config.ArgCols, "", defaultImageCols, printer.ColsMessage(allImageCols))
	_ = viper.BindPFlag(core.GetGlobalFlagName(serverCdromCmd.Name(), config.ArgCols), globalFlags.Lookup(config.ArgCols))
	_ = serverCdromCmd.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allImageCols, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Attach Cdrom Command
	*/
	attachCdrom := core.NewCommand(ctx, serverCdromCmd, core.CommandBuilder{
		Namespace: "server",
		Resource:  "cdrom",
		Verb:      "attach",
		Aliases:   []string{"a"},
		ShortDesc: "Attach a CD-ROM to a Server",
		LongDesc: `Use this command to attach a CD-ROM to an existing Server.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option.

Required values to run command:

* Data Center Id
* Server Id
* Cdrom Id`,
		Example:    attachCdromServerExample,
		PreCmdRun:  PreRunDcServerCdromIds,
		CmdRun:     RunServerCdromAttach,
		InitClient: true,
	})
	attachCdrom.AddStringFlag(cloudapi_v6.ArgDataCenterId, "", "", cloudapi_v6.DatacenterId, core.RequiredFlagOption())
	_ = attachCdrom.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	attachCdrom.AddStringFlag(cloudapi_v6.ArgCdromId, cloudapi_v6.ArgIdShort, "", cloudapi_v6.CdromId, core.RequiredFlagOption())
	_ = attachCdrom.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgCdromId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getImagesCdromIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	attachCdrom.AddStringFlag(cloudapi_v6.ArgServerId, "", "", cloudapi_v6.ServerId, core.RequiredFlagOption())
	_ = attachCdrom.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgServerId, func(cmd *cobra.Command, ags []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getServersIds(os.Stderr, viper.GetString(core.GetFlagName(attachCdrom.NS, cloudapi_v6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	attachCdrom.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for CD-ROM attachment to be executed")
	attachCdrom.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option for Request for Cdrom attachment [seconds]")

	/*
		List Cdroms Command
	*/
	listCdroms := core.NewCommand(ctx, serverCdromCmd, core.CommandBuilder{
		Namespace:  "server",
		Resource:   "cdrom",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List attached CD-ROMs from a Server",
		LongDesc:   "Use this command to retrieve a list of CD-ROMs attached to the Server.\n\nRequired values to run command:\n\n* Data Center Id\n* Server Id",
		Example:    listCdromServerExample,
		PreCmdRun:  PreRunDcServerIds,
		CmdRun:     RunServerCdromsList,
		InitClient: true,
	})
	listCdroms.AddStringFlag(cloudapi_v6.ArgDataCenterId, "", "", cloudapi_v6.DatacenterId, core.RequiredFlagOption())
	_ = listCdroms.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	listCdroms.AddStringFlag(cloudapi_v6.ArgServerId, "", "", cloudapi_v6.ServerId, core.RequiredFlagOption())
	_ = listCdroms.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgServerId, func(cmd *cobra.Command, ags []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getServersIds(os.Stderr, viper.GetString(core.GetFlagName(listCdroms.NS, cloudapi_v6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Get Cdrom Command
	*/
	getCdromCmd := core.NewCommand(ctx, serverCdromCmd, core.CommandBuilder{
		Namespace:  "server",
		Resource:   "cdrom",
		Verb:       "get",
		Aliases:    []string{"g"},
		ShortDesc:  "Get a specific attached CD-ROM from a Server",
		LongDesc:   "Use this command to retrieve information about an attached CD-ROM on Server.\n\nRequired values to run command:\n\n* Data Center Id\n* Server Id\n* Cdrom Id",
		Example:    getCdromServerExample,
		InitClient: true,
		PreCmdRun:  PreRunDcServerCdromIds,
		CmdRun:     RunServerCdromGet,
	})
	getCdromCmd.AddStringFlag(cloudapi_v6.ArgDataCenterId, "", "", cloudapi_v6.DatacenterId, core.RequiredFlagOption())
	_ = getCdromCmd.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	getCdromCmd.AddStringFlag(cloudapi_v6.ArgServerId, "", "", cloudapi_v6.ServerId, core.RequiredFlagOption())
	_ = getCdromCmd.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgServerId, func(cmd *cobra.Command, ags []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getServersIds(os.Stderr, viper.GetString(core.GetFlagName(getCdromCmd.NS, cloudapi_v6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	getCdromCmd.AddStringFlag(cloudapi_v6.ArgCdromId, cloudapi_v6.ArgIdShort, "", cloudapi_v6.CdromId, core.RequiredFlagOption())
	_ = getCdromCmd.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgCdromId, func(cmd *cobra.Command, ags []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getAttachedCdromsIds(os.Stderr, viper.GetString(core.GetFlagName(getCdromCmd.NS, cloudapi_v6.ArgDataCenterId)), viper.GetString(core.GetFlagName(getCdromCmd.NS, cloudapi_v6.ArgServerId))), cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Detach Cdrom Command
	*/
	detachCdrom := core.NewCommand(ctx, serverCdromCmd, core.CommandBuilder{
		Namespace: "server",
		Resource:  "cdrom",
		Verb:      "detach",
		Aliases:   []string{"d"},
		ShortDesc: "Detach a CD-ROM from a Server",
		LongDesc: `This will detach the CD-ROM from the Server.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option. You can force the command to execute without user input using ` + "`" + `--force` + "`" + ` option.

Required values to run command:

* Data Center Id
* Server Id
* Cdrom Id`,
		Example:    detachCdromServerExample,
		PreCmdRun:  PreRunDcServerCdromIds,
		CmdRun:     RunServerCdromDetach,
		InitClient: true,
	})
	detachCdrom.AddStringFlag(cloudapi_v6.ArgDataCenterId, "", "", cloudapi_v6.DatacenterId, core.RequiredFlagOption())
	_ = detachCdrom.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	detachCdrom.AddStringFlag(cloudapi_v6.ArgCdromId, cloudapi_v6.ArgIdShort, "", cloudapi_v6.CdromId, core.RequiredFlagOption())
	_ = detachCdrom.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgCdromId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getAttachedCdromsIds(os.Stderr, viper.GetString(core.GetFlagName(detachCdrom.NS, cloudapi_v6.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(detachCdrom.NS, cloudapi_v6.ArgServerId))), cobra.ShellCompDirectiveNoFileComp
	})
	detachCdrom.AddStringFlag(cloudapi_v6.ArgServerId, "", "", cloudapi_v6.ServerId, core.RequiredFlagOption())
	_ = detachCdrom.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getServersIds(os.Stderr, viper.GetString(core.GetFlagName(detachCdrom.NS, cloudapi_v6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	detachCdrom.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for CD-ROM detachment to be executed")
	detachCdrom.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option for Request for CD-ROM detachment [seconds]")

	return serverCdromCmd
}

func PreRunDcServerCdromIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapi_v6.ArgDataCenterId, cloudapi_v6.ArgServerId, cloudapi_v6.ArgCdromId)
}

func RunServerCdromAttach(c *core.CommandConfig) error {
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgDataCenterId))
	serverId := viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgServerId))
	cdRomId := viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgCdromId))
	c.Printer.Verbose("CD-ROM with id: %v is attaching to server with id: %v from Datacenter with id: %v... ", cdRomId, serverId, dcId)
	attachedCdrom, resp, err := c.CloudApiV6Services.Servers().AttachCdrom(dcId, serverId, cdRomId)
	if err != nil {
		return err
	}

	if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetRequestPath(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getImagePrint(resp, c, getImage(attachedCdrom)))
}

func RunServerCdromsList(c *core.CommandConfig) error {
	attachedCdroms, _, err := c.CloudApiV6Services.Servers().ListCdroms(
		viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgServerId)),
	)
	if err != nil {
		return err
	}
	return c.Printer.Print(getImagePrint(nil, c, getCdroms(attachedCdroms)))
}

func RunServerCdromGet(c *core.CommandConfig) error {
	c.Printer.Verbose("CD-ROM with id: %v is getting... ", viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgCdromId)))
	attachedCdrom, _, err := c.CloudApiV6Services.Servers().GetCdrom(
		viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgServerId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgCdromId)),
	)
	if err != nil {
		return err
	}
	return c.Printer.Print(getImagePrint(nil, c, getImage(attachedCdrom)))
}

func RunServerCdromDetach(c *core.CommandConfig) error {
	if err := utils.AskForConfirm(c.Stdin, c.Printer, "detach cdrom from server"); err != nil {
		return err
	}
	c.Printer.Verbose("CD-ROM with id: %v is detaching... ", viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgCdromId)))
	resp, err := c.CloudApiV6Services.Servers().DetachCdrom(
		viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgServerId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgCdromId)),
	)

	if err != nil {
		return err
	}

	if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetRequestPath(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getImagePrint(resp, c, nil))
}

// Output Printing

func getCdroms(cdroms resources.Cdroms) []resources.Image {
	imgs := make([]resources.Image, 0)
	if items, ok := cdroms.GetItemsOk(); ok && items != nil {
		for _, d := range *items {
			imgs = append(imgs, resources.Image{Image: d})
		}
	}
	return imgs
}

func getAttachedCdromsIds(outErr io.Writer, datacenterId, serverId string) []string {
	err := config.Load()
	clierror.CheckError(err, outErr)
	clientSvc, err := resources.NewClientService(
		viper.GetString(config.Username),
		viper.GetString(config.Password),
		viper.GetString(config.Token),
		config.GetServerUrl(),
	)
	clierror.CheckError(err, outErr)
	serverSvc := resources.NewServerService(clientSvc.Get(), context.TODO())
	cdroms, _, err := serverSvc.ListCdroms(datacenterId, serverId)
	clierror.CheckError(err, outErr)
	attachedCdromsIds := make([]string, 0)
	if items, ok := cdroms.GetItemsOk(); ok && items != nil {
		for _, item := range *items {
			if itemId, ok := item.GetIdOk(); ok && itemId != nil {
				attachedCdromsIds = append(attachedCdromsIds, *itemId)
			}
		}
	} else {
		return nil
	}
	return attachedCdromsIds
}

func getImagesCdromIds(outErr io.Writer) []string {
	err := config.Load()
	clierror.CheckError(err, outErr)
	clientSvc, err := resources.NewClientService(
		viper.GetString(config.Username),
		viper.GetString(config.Password),
		viper.GetString(config.Token),
		config.GetServerUrl(),
	)
	clierror.CheckError(err, outErr)
	imageSvc := resources.NewImageService(clientSvc.Get(), context.TODO())
	images, _, err := imageSvc.List()
	clierror.CheckError(err, outErr)
	imgsIds := make([]string, 0)
	images = sortImagesByType(images, "CDROM")
	if items, ok := images.Images.GetItemsOk(); ok && items != nil {
		for _, item := range *items {
			if itemId, ok := item.GetIdOk(); ok && itemId != nil {
				imgsIds = append(imgsIds, *itemId)
			}
		}
	} else {
		return nil
	}
	return imgsIds
}
