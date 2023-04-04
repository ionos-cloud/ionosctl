package commands

import (
	"context"
	"errors"
	"fmt"
	"os"

	"go.uber.org/multierr"

	"github.com/fatih/structs"
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

func IpblockCmd() *core.Command {
	ctx := context.TODO()
	ipblockCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "ipblock",
			Aliases:          []string{"ip", "ipb"},
			Short:            "IpBlock Operations",
			Long:             "The sub-commands of `ionosctl ipblock` allow you to create/reserve, list, get, update, delete IpBlocks.",
			TraverseChildren: true,
		},
	}
	globalFlags := ipblockCmd.GlobalFlags()
	globalFlags.StringSliceP(constants.ArgCols, "", defaultIpBlockCols, printer.ColsMessage(defaultIpBlockCols))
	_ = viper.BindPFlag(core.GetFlagName(ipblockCmd.Name(), constants.ArgCols), globalFlags.Lookup(constants.ArgCols))
	_ = ipblockCmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultIpBlockCols, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		List Command
	*/
	list := core.NewCommand(ctx, ipblockCmd, core.CommandBuilder{
		Namespace:  "ipblock",
		Resource:   "ipblock",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List IpBlocks",
		LongDesc:   "Use this command to list IpBlocks.\n\nYou can filter the results using `--filters` option. Use the following format to set filters: `--filters KEY1=VALUE1,KEY2=VALUE2`.\n" + completer.IpBlocksFiltersUsage(),
		Example:    listIpBlockExample,
		PreCmdRun:  PreRunIpblockList,
		CmdRun:     RunIpBlockList,
		InitClient: true,
	})
	list.AddInt32Flag(constants.FlagMaxResults, constants.FlagMaxResultsShort, cloudapiv6.DefaultMaxResults, constants.DescMaxResults)
	list.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultListDepth, cloudapiv6.ArgDepthDescription)
	list.AddStringFlag(cloudapiv6.ArgOrderBy, "", "", cloudapiv6.ArgOrderByDescription)
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgOrderBy, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.IpBlocksFilters(), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddStringSliceFlag(cloudapiv6.ArgFilters, cloudapiv6.ArgFiltersShort, []string{""}, cloudapiv6.ArgFiltersDescription)
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgFilters, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.IpBlocksFilters(), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddBoolFlag(constants.ArgNoHeaders, "", false, cloudapiv6.ArgNoHeadersDescription)

	/*
		Get Command
	*/
	get := core.NewCommand(ctx, ipblockCmd, core.CommandBuilder{
		Namespace:  "ipblock",
		Resource:   "ipblock",
		Verb:       "get",
		Aliases:    []string{"g"},
		ShortDesc:  "Get an IpBlock",
		LongDesc:   "Use this command to retrieve the attributes of a specific IpBlock.\n\nRequired values to run command:\n\n* IpBlock Id",
		Example:    getIpBlockExample,
		PreCmdRun:  PreRunIpBlockId,
		CmdRun:     RunIpBlockGet,
		InitClient: true,
	})
	get.AddUUIDFlag(cloudapiv6.ArgIpBlockId, cloudapiv6.ArgIdShort, "", cloudapiv6.IpBlockId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgIpBlockId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.IpBlocksIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddBoolFlag(constants.ArgNoHeaders, "", false, cloudapiv6.ArgNoHeadersDescription)
	get.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultGetDepth, cloudapiv6.ArgDepthDescription)

	/*
		Create Command
	*/
	create := core.NewCommand(ctx, ipblockCmd, core.CommandBuilder{
		Namespace: "ipblock",
		Resource:  "ipblock",
		Verb:      "create",
		Aliases:   []string{"c"},
		ShortDesc: "Create/Reserve an IpBlock",
		LongDesc: `Use this command to create/reserve an IpBlock in a specified location that can be used by resources within any Virtual Data Centers provisioned in that same location. An IpBlock consists of one or more static IP addresses. The name, size of the IpBlock can be set.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option.`,
		Example:    createIpBlockExample,
		PreCmdRun:  core.NoPreRun,
		CmdRun:     RunIpBlockCreate,
		InitClient: true,
	})
	create.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "", "Name of the IpBlock. If not set, it will automatically be set")
	create.AddStringFlag(cloudapiv6.ArgLocation, cloudapiv6.ArgLocationShort, "de/txl", "Location of the IpBlock")
	_ = create.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgLocation, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.LocationIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddIntFlag(cloudapiv6.ArgSize, "", 2, "Size of the IpBlock")
	create.AddBoolFlag(constants.ArgWaitForRequest, constants.ArgWaitForRequestShort, constants.DefaultWait, "Wait for the Request for IpBlock creation to be executed")
	create.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, constants.DefaultTimeoutSeconds, "Timeout option for Request for IpBlock creation [seconds]")
	create.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultCreateDepth, cloudapiv6.ArgDepthDescription)

	/*
		Update Command
	*/
	update := core.NewCommand(ctx, ipblockCmd, core.CommandBuilder{
		Namespace: "ipblock",
		Resource:  "ipblock",
		Verb:      "update",
		Aliases:   []string{"u", "up"},
		ShortDesc: "Update an IpBlock",
		LongDesc: `Use this command to update the properties of an existing IpBlock.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option.

Required values to run command:

* IpBlock Id`,
		Example:    updateIpBlockExample,
		PreCmdRun:  PreRunIpBlockId,
		CmdRun:     RunIpBlockUpdate,
		InitClient: true,
	})
	update.AddUUIDFlag(cloudapiv6.ArgIpBlockId, cloudapiv6.ArgIdShort, "", cloudapiv6.IpBlockId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgIpBlockId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.IpBlocksIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "", "Name of the IpBlock")
	update.AddBoolFlag(constants.ArgWaitForRequest, constants.ArgWaitForRequestShort, constants.DefaultWait, "Wait for the Request for IpBlock update to be executed")
	update.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, constants.DefaultTimeoutSeconds, "Timeout option for Request for IpBlock update [seconds]")
	update.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultUpdateDepth, cloudapiv6.ArgDepthDescription)

	/*
		Delete Command
	*/
	deleteCmd := core.NewCommand(ctx, ipblockCmd, core.CommandBuilder{
		Namespace: "ipblock",
		Resource:  "ipblock",
		Verb:      "delete",
		Aliases:   []string{"d"},
		ShortDesc: "Delete an IpBlock",
		LongDesc: `Use this command to delete a specified IpBlock.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option. You can force the command to execute without user input using ` + "`" + `--force` + "`" + ` option.

Required values to run command:

* IpBlock Id`,
		Example:    deleteIpBlockExample,
		PreCmdRun:  PreRunIpBlockDelete,
		CmdRun:     RunIpBlockDelete,
		InitClient: true,
	})
	deleteCmd.AddUUIDFlag(cloudapiv6.ArgIpBlockId, cloudapiv6.ArgIdShort, "", cloudapiv6.IpBlockId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgIpBlockId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.IpBlocksIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddBoolFlag(constants.ArgWaitForRequest, constants.ArgWaitForRequestShort, constants.DefaultWait, "Wait for the Request for IpBlock deletion to be executed")
	deleteCmd.AddBoolFlag(cloudapiv6.ArgAll, cloudapiv6.ArgAllShort, false, "Delete all the IpBlocks.")
	deleteCmd.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, constants.DefaultTimeoutSeconds, "Timeout option for Request for IpBlock deletion [seconds]")
	deleteCmd.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultDeleteDepth, cloudapiv6.ArgDepthDescription)

	return ipblockCmd
}

func PreRunIpblockList(c *core.PreCommandConfig) error {
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgFilters)) {
		return query.ValidateFilters(c, completer.IpBlocksFilters(), completer.IpBlocksFiltersUsage())
	}
	return nil
}

func PreRunIpBlockId(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgIpBlockId)
}

func PreRunIpBlockDelete(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{cloudapiv6.ArgIpBlockId},
		[]string{cloudapiv6.ArgAll},
	)
}

func RunIpBlockList(c *core.CommandConfig) error {
	// Add Query Parameters for GET Requests
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	ipblocks, resp, err := c.CloudApiV6Services.IpBlocks().List(listQueryParams)
	if resp != nil {
		c.Printer.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getIpBlockPrint(nil, c, getIpBlocks(ipblocks)))
}

func RunIpBlockGet(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	queryParams := listQueryParams.QueryParams
	c.Printer.Verbose("Ip block with id: %v is getting...", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgIpBlockId)))
	i, resp, err := c.CloudApiV6Services.IpBlocks().Get(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgIpBlockId)), queryParams)
	if resp != nil {
		c.Printer.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getIpBlockPrint(nil, c, getIpBlock(i)))
}

func RunIpBlockCreate(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	queryParams := listQueryParams.QueryParams
	name := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName))
	loc := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLocation))
	size := viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgSize))
	c.Printer.Verbose("Properties set for creating the Ip block: Name: %v, Location: %v, Size: %v", name, loc, size)
	i, resp, err := c.CloudApiV6Services.IpBlocks().Create(name, loc, size, queryParams)
	if resp != nil && printer.GetId(resp) != "" {
		c.Printer.Verbose(constants.MessageRequestInfo, printer.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}
	if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getIpBlockPrint(resp, c, getIpBlock(i)))
}

func RunIpBlockUpdate(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	queryParams := listQueryParams.QueryParams
	input := resources.IpBlockProperties{}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgName)) {
		name := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName))
		input.SetName(name)
		c.Printer.Verbose("Property Name set: %v", name)
	}
	i, resp, err := c.CloudApiV6Services.IpBlocks().Update(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgIpBlockId)), input, queryParams)
	if resp != nil && printer.GetId(resp) != "" {
		c.Printer.Verbose(constants.MessageRequestInfo, printer.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}
	if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getIpBlockPrint(resp, c, getIpBlock(i)))
}

func RunIpBlockDelete(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	queryParams := listQueryParams.QueryParams
	ipBlockId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgIpBlockId))
	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAll)) {
		if err := DeleteAllIpBlocks(c); err != nil {
			return err
		}
		return c.Printer.Print(printer.Result{Resource: c.Resource, Verb: c.Verb})
	} else {
		if err := utils.AskForConfirm(c.Stdin, c.Printer, "delete ipblock"); err != nil {
			return err
		}
		c.Printer.Verbose("Starting deleting Ip block with ID: %v...", ipBlockId)
		resp, err := c.CloudApiV6Services.IpBlocks().Delete(ipBlockId, queryParams)
		if resp != nil && printer.GetId(resp) != "" {
			c.Printer.Verbose(constants.MessageRequestInfo, printer.GetId(resp), resp.RequestTime)
		}
		if err != nil {
			return err
		}
		if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
			return err
		}
		return c.Printer.Print(getIpBlockPrint(resp, c, nil))
	}
}

func DeleteAllIpBlocks(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	queryParams := listQueryParams.QueryParams
	c.Printer.Verbose("Getting all IpBlocks...")
	ipBlocks, resp, err := c.CloudApiV6Services.IpBlocks().List(cloudapiv6.ParentResourceListQueryParams)
	if err != nil {
		return err
	}
	if ipBlocksItems, ok := ipBlocks.GetItemsOk(); ok && ipBlocksItems != nil {
		if len(*ipBlocksItems) > 0 {
			_ = c.Printer.Warn("IpBlocks to be deleted:")
			for _, dc := range *ipBlocksItems {
				delIdAndName := ""
				if id, ok := dc.GetIdOk(); ok && id != nil {
					delIdAndName += "IpBlock Id: " + *id
				}
				if properties, ok := dc.GetPropertiesOk(); ok && properties != nil {
					if name, ok := properties.GetNameOk(); ok && name != nil {
						delIdAndName += " IpBlock Name: " + *name
					}
				}
				_ = c.Printer.Warn(delIdAndName)
			}
			if err = utils.AskForConfirm(c.Stdin, c.Printer, "delete all the IpBlocks"); err != nil {
				return err
			}
			c.Printer.Verbose("Deleting all the IpBlocks...")
			var multiErr error
			for _, dc := range *ipBlocksItems {
				if id, ok := dc.GetIdOk(); ok && id != nil {
					c.Printer.Verbose("Starting deleting IpBlock with id: %v...", *id)
					resp, err = c.CloudApiV6Services.IpBlocks().Delete(*id, queryParams)
					if resp != nil && printer.GetId(resp) != "" {
						c.Printer.Verbose(constants.MessageRequestInfo, printer.GetId(resp), resp.RequestTime)
					}
					if err != nil {
						multiErr = multierr.Append(multiErr, fmt.Errorf(constants.ErrDeleteAll, c.Resource, *id, err))
						continue
					} else {
						_ = c.Printer.Warn(fmt.Sprintf(constants.MessageDeletingAll, c.Resource, *id))
					}
					if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
						multiErr = multierr.Append(multiErr, fmt.Errorf(constants.ErrWaitDeleteAll, c.Resource, *id, err))
						continue
					}
				}
			}
			if multiErr != nil {
				return multiErr
			}
			return nil
		} else {
			return errors.New("no IpBlocks found")
		}
	} else {
		return errors.New("could not get items of IpBlocks")
	}
}

// Output Printing

var defaultIpBlockCols = []string{"IpBlockId", "Name", "Location", "Size", "Ips", "State"}

type IpBlockPrint struct {
	IpBlockId string   `json:"IpBlockId,omitempty"`
	Name      string   `json:"Name,omitempty"`
	Location  string   `json:"Location,omitempty"`
	Size      int32    `json:"Size,omitempty"`
	Ips       []string `json:"Ips,omitempty"`
	State     string   `json:"State,omitempty"`
}

func getIpBlockPrint(resp *resources.Response, c *core.CommandConfig, ipBlocks []resources.IpBlock) printer.Result {
	r := printer.Result{}
	if c != nil {
		if resp != nil {
			r.ApiResponse = resp
			r.Resource = c.Resource
			r.Verb = c.Verb
			r.WaitForRequest = viper.GetBool(core.GetFlagName(c.NS, constants.ArgWaitForRequest))
		}
		if ipBlocks != nil {
			r.OutputJSON = ipBlocks
			r.KeyValue = getIpBlocksKVMaps(ipBlocks)
			r.Columns = printer.GetHeadersAllDefault(defaultIpBlockCols, viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols)))
		}
	}
	return r
}

func getIpBlocks(ipBlocks resources.IpBlocks) []resources.IpBlock {
	ss := make([]resources.IpBlock, 0)
	if items, ok := ipBlocks.GetItemsOk(); ok && items != nil {
		for _, item := range *items {
			ss = append(ss, resources.IpBlock{IpBlock: item})
		}
	}
	return ss
}

func getIpBlock(ipBlock *resources.IpBlock) []resources.IpBlock {
	ss := make([]resources.IpBlock, 0)
	if ipBlock != nil {
		ss = append(ss, resources.IpBlock{IpBlock: ipBlock.IpBlock})
	}
	return ss
}

func getIpBlocksKVMaps(ss []resources.IpBlock) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(ss))
	for _, s := range ss {
		o := getIpBlockKVMap(s)
		out = append(out, o)
	}
	return out
}

func getIpBlockKVMap(s resources.IpBlock) map[string]interface{} {
	var ipblockPrint IpBlockPrint
	if id, ok := s.GetIdOk(); ok && id != nil {
		ipblockPrint.IpBlockId = *id
	}
	if properties, ok := s.GetPropertiesOk(); ok && properties != nil {
		if name, ok := properties.GetNameOk(); ok && name != nil {
			ipblockPrint.Name = *name
		}
		if loc, ok := properties.GetLocationOk(); ok && loc != nil {
			ipblockPrint.Location = *loc
		}
		if size, ok := properties.GetSizeOk(); ok && size != nil {
			ipblockPrint.Size = *size
		}
		if ips, ok := properties.GetIpsOk(); ok && ips != nil {
			ipblockPrint.Ips = *ips
		}
	}
	if metadata, ok := s.GetMetadataOk(); ok && metadata != nil {
		if state, ok := metadata.GetStateOk(); ok && state != nil {
			ipblockPrint.State = *state
		}
	}
	return structs.Map(ipblockPrint)
}
