package commands

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/completer"
	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/query"
	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/waiter"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/printer"
	"github.com/ionos-cloud/ionosctl/v6/pkg/utils"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
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
	globalFlags.StringSliceP(constants.ArgCols, "", defaultImageCols, printer.ColsMessage(allImageCols))
	_ = viper.BindPFlag(core.GetFlagName(serverCdromCmd.Name(), constants.ArgCols), globalFlags.Lookup(constants.ArgCols))
	_ = serverCdromCmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
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
	attachCdrom.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = attachCdrom.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	attachCdrom.AddUUIDFlag(cloudapiv6.ArgCdromId, cloudapiv6.ArgIdShort, "", cloudapiv6.CdromId, core.RequiredFlagOption())
	_ = attachCdrom.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgCdromId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ImagesIdsCustom(os.Stderr, resources.ListQueryParams{Filters: &map[string][]string{
			"type": {"CDROM"},
		}}), cobra.ShellCompDirectiveNoFileComp
	})
	attachCdrom.AddUUIDFlag(cloudapiv6.ArgServerId, "", "", cloudapiv6.ServerId, core.RequiredFlagOption())
	_ = attachCdrom.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgServerId, func(cmd *cobra.Command, ags []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ServersIds(os.Stderr, viper.GetString(core.GetFlagName(attachCdrom.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	attachCdrom.AddBoolFlag(constants.ArgWaitForRequest, constants.ArgWaitForRequestShort, constants.DefaultWait, "Wait for the Request for CD-ROM attachment to be executed")
	attachCdrom.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, constants.DefaultTimeoutSeconds, "Timeout option for Request for Cdrom attachment [seconds]")
	attachCdrom.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultMiscDepth, cloudapiv6.ArgDepthDescription)

	/*
		List Cdroms Command
	*/
	listCdroms := core.NewCommand(ctx, serverCdromCmd, core.CommandBuilder{
		Namespace:  "server",
		Resource:   "cdrom",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List attached CD-ROMs from a Server",
		LongDesc:   "Use this command to retrieve a list of CD-ROMs attached to the Server.\n\nYou can filter the results using `--filters` option. Use the following format to set filters: `--filters KEY1=VALUE1,KEY2=VALUE2`.\n" + completer.ImagesFiltersUsage() + "\n\nRequired values to run command:\n\n* Data Center Id\n* Server Id",
		Example:    listCdromServerExample,
		PreCmdRun:  PreRunServerCdromList,
		CmdRun:     RunServerCdromsList,
		InitClient: true,
	})
	listCdroms.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = listCdroms.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	listCdroms.AddUUIDFlag(cloudapiv6.ArgServerId, "", "", cloudapiv6.ServerId, core.RequiredFlagOption())
	_ = listCdroms.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgServerId, func(cmd *cobra.Command, ags []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ServersIds(os.Stderr, viper.GetString(core.GetFlagName(listCdroms.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	listCdroms.AddInt32Flag(constants.FlagMaxResults, constants.FlagMaxResultsShort, cloudapiv6.DefaultMaxResults, constants.DescMaxResults)
	listCdroms.AddStringFlag(cloudapiv6.ArgOrderBy, "", "", cloudapiv6.ArgOrderByDescription)
	_ = listCdroms.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgOrderBy, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ImagesFilters(), cobra.ShellCompDirectiveNoFileComp
	})
	listCdroms.AddStringSliceFlag(cloudapiv6.ArgFilters, cloudapiv6.ArgFiltersShort, []string{""}, cloudapiv6.ArgFiltersDescription)
	_ = listCdroms.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgFilters, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ImagesFilters(), cobra.ShellCompDirectiveNoFileComp
	})
	listCdroms.AddBoolFlag(constants.ArgNoHeaders, "", false, cloudapiv6.ArgNoHeadersDescription)
	listCdroms.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultListDepth, cloudapiv6.ArgDepthDescription)

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
	getCdromCmd.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = getCdromCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	getCdromCmd.AddUUIDFlag(cloudapiv6.ArgServerId, "", "", cloudapiv6.ServerId, core.RequiredFlagOption())
	_ = getCdromCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgServerId, func(cmd *cobra.Command, ags []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ServersIds(os.Stderr, viper.GetString(core.GetFlagName(getCdromCmd.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	getCdromCmd.AddUUIDFlag(cloudapiv6.ArgCdromId, cloudapiv6.ArgIdShort, "", cloudapiv6.CdromId, core.RequiredFlagOption())
	_ = getCdromCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgCdromId, func(cmd *cobra.Command, ags []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.AttachedCdromsIds(os.Stderr, viper.GetString(core.GetFlagName(getCdromCmd.NS, cloudapiv6.ArgDataCenterId)), viper.GetString(core.GetFlagName(getCdromCmd.NS, cloudapiv6.ArgServerId))), cobra.ShellCompDirectiveNoFileComp
	})
	getCdromCmd.AddBoolFlag(constants.ArgNoHeaders, "", false, cloudapiv6.ArgNoHeadersDescription)
	getCdromCmd.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultGetDepth, cloudapiv6.ArgDepthDescription)

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
		PreCmdRun:  PreRunDcServerCdromDetach,
		CmdRun:     RunServerCdromDetach,
		InitClient: true,
	})
	detachCdrom.AddUUIDFlag(cloudapiv6.ArgDataCenterId, "", "", cloudapiv6.DatacenterId, core.RequiredFlagOption())
	_ = detachCdrom.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	detachCdrom.AddUUIDFlag(cloudapiv6.ArgCdromId, cloudapiv6.ArgIdShort, "", cloudapiv6.CdromId, core.RequiredFlagOption())
	_ = detachCdrom.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgCdromId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.AttachedCdromsIds(os.Stderr, viper.GetString(core.GetFlagName(detachCdrom.NS, cloudapiv6.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(detachCdrom.NS, cloudapiv6.ArgServerId))), cobra.ShellCompDirectiveNoFileComp
	})
	detachCdrom.AddUUIDFlag(cloudapiv6.ArgServerId, "", "", cloudapiv6.ServerId, core.RequiredFlagOption())
	_ = detachCdrom.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ServersIds(os.Stderr, viper.GetString(core.GetFlagName(detachCdrom.NS, cloudapiv6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	detachCdrom.AddBoolFlag(constants.ArgWaitForRequest, constants.ArgWaitForRequestShort, constants.DefaultWait, "Wait for the Request for CD-ROM detachment to be executed")
	detachCdrom.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, constants.DefaultTimeoutSeconds, "Timeout option for Request for CD-ROM detachment [seconds]")
	detachCdrom.AddBoolFlag(cloudapiv6.ArgAll, cloudapiv6.ArgAllShort, false, "Detach all CD-ROMS from a Server.")
	detachCdrom.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultMiscDepth, cloudapiv6.ArgDepthDescription)

	return serverCdromCmd
}

func PreRunServerCdromList(c *core.PreCommandConfig) error {
	if err := core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgDataCenterId, cloudapiv6.ArgServerId); err != nil {
		return err
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgFilters)) {
		return query.ValidateFilters(c, completer.ImagesFilters(), completer.ImagesFiltersUsage())
	}
	return nil
}

func PreRunDcServerCdromIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgDataCenterId, cloudapiv6.ArgServerId, cloudapiv6.ArgCdromId)
}

func PreRunDcServerCdromDetach(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{cloudapiv6.ArgDataCenterId, cloudapiv6.ArgServerId, cloudapiv6.ArgCdromId},
		[]string{cloudapiv6.ArgDataCenterId, cloudapiv6.ArgServerId, cloudapiv6.ArgAll},
	)
}

func RunServerCdromAttach(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	queryParams := listQueryParams.QueryParams
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))
	serverId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgServerId))
	cdRomId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgCdromId))
	c.Printer.Verbose("CD-ROM with id: %v is attaching to server with id: %v from Datacenter with id: %v... ", cdRomId, serverId, dcId)
	attachedCdrom, resp, err := c.CloudApiV6Services.Servers().AttachCdrom(dcId, serverId, cdRomId, queryParams)
	if resp != nil && printer.GetId(resp) != "" {
		c.Printer.Verbose(constants.MessageRequestInfo, printer.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}

	if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getImagePrint(resp, c, getImage(attachedCdrom)))
}

func RunServerCdromsList(c *core.CommandConfig) error {
	// Add Query Parameters for GET Requests
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	attachedCdroms, resp, err := c.CloudApiV6Services.Servers().ListCdroms(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgServerId)),
		listQueryParams,
	)
	if resp != nil {
		c.Printer.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getImagePrint(nil, c, getCdroms(attachedCdroms)))
}

func RunServerCdromGet(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	queryParams := listQueryParams.QueryParams
	c.Printer.Verbose("CD-ROM with id: %v is getting... ", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgCdromId)))
	attachedCdrom, resp, err := c.CloudApiV6Services.Servers().GetCdrom(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgServerId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgCdromId)),
		queryParams,
	)
	if resp != nil {
		c.Printer.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getImagePrint(nil, c, getImage(attachedCdrom)))
}

func RunServerCdromDetach(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	queryParams := listQueryParams.QueryParams
	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAll)) {
		if err := DetachAllCdRoms(c); err != nil {
			return err
		}
		return c.Printer.Print(printer.Result{Resource: c.Resource, Verb: c.Verb})
	} else {
		if err := utils.AskForConfirm(c.Stdin, c.Printer, "detach cdrom from server"); err != nil {
			return err
		}
		c.Printer.Verbose("CD-ROM with id: %v is detaching... ", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgCdromId)))
		resp, err := c.CloudApiV6Services.Servers().DetachCdrom(
			viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgServerId)),
			viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgCdromId)),
			queryParams,
		)
		if resp != nil && printer.GetId(resp) != "" {
			c.Printer.Verbose(constants.MessageRequestInfo, printer.GetId(resp), resp.RequestTime)
		}
		if err != nil {
			return err
		}
		if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
			return err
		}
		return c.Printer.Print(getImagePrint(resp, c, nil))
	}
}

func DetachAllCdRoms(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	queryParams := listQueryParams.QueryParams
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))
	serverId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgServerId))
	c.Printer.Verbose("Getting CD-ROMs...")
	cdRoms, resp, err := c.CloudApiV6Services.Servers().ListCdroms(dcId, serverId, cloudapiv6.ParentResourceListQueryParams)
	if err != nil {
		return err
	}
	if cdRomsItems, ok := cdRoms.GetItemsOk(); ok && cdRomsItems != nil {
		if len(*cdRomsItems) > 0 {
			_ = c.Printer.Warn("CD-ROMS to be detached:")
			delIdAndName := ""
			for _, cdRom := range *cdRomsItems {
				if id, ok := cdRom.GetIdOk(); ok && id != nil {
					delIdAndName += "CD-ROM Id: " + *id
				}
				if properties, ok := cdRom.GetPropertiesOk(); ok && properties != nil {
					if name, ok := properties.GetNameOk(); ok && name != nil {
						delIdAndName += " CD-ROM Name: " + *name
					}
				}
				_ = c.Printer.Warn(delIdAndName)
			}
			if err := utils.AskForConfirm(c.Stdin, c.Printer, "detach all the CD-ROMS"); err != nil {
				return err
			}
			c.Printer.Verbose("Detaching all the CD-ROM...")
			var multiErr error
			for _, cdRom := range *cdRomsItems {
				if id, ok := cdRom.GetIdOk(); ok && id != nil {
					c.Printer.Verbose("Starting detaching CD-ROM with id: %v...", *id)
					resp, err = c.CloudApiV6Services.Servers().DetachCdrom(dcId, serverId, *id, queryParams)
					if resp != nil && printer.GetId(resp) != "" {
						c.Printer.Verbose(constants.MessageRequestInfo, printer.GetId(resp), resp.RequestTime)
					}
					if err != nil {
						multiErr = errors.Join(multiErr, fmt.Errorf(constants.ErrDeleteAll, c.Resource, *id, err))
						continue
					} else {
						_ = c.Printer.Warn(fmt.Sprintf(constants.MessageDeletingAll, c.Resource, *id))
					}
					if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
						multiErr = errors.Join(multiErr, fmt.Errorf(constants.ErrWaitDeleteAll, c.Resource, *id, err))
						continue
					}
				}
			}
			if multiErr != nil {
				return multiErr
			}
			return nil
		} else {
			return errors.New("no CD-ROMs found")
		}
	} else {
		return errors.New("could not get items of CD-ROMs")
	}
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
