package commands

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/commands/cloudapi-v6/completer"
	"github.com/ionos-cloud/ionosctl/commands/cloudapi-v6/waiter"
	"github.com/ionos-cloud/ionosctl/internal/config"
	"github.com/ionos-cloud/ionosctl/internal/core"
	"github.com/ionos-cloud/ionosctl/internal/printer"
	"github.com/ionos-cloud/ionosctl/internal/utils"
	"github.com/ionos-cloud/ionosctl/internal/utils/clierror"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/services/cloudapi-v6/resources"
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
	globalFlags.StringSliceP(config.ArgCols, "", defaultIpBlockCols, printer.ColsMessage(defaultIpBlockCols))
	_ = viper.BindPFlag(core.GetGlobalFlagName(ipblockCmd.Name(), config.ArgCols), globalFlags.Lookup(config.ArgCols))
	_ = ipblockCmd.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultIpBlockCols, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		List Command
	*/
	core.NewCommand(ctx, ipblockCmd, core.CommandBuilder{
		Namespace:  "ipblock",
		Resource:   "ipblock",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List IpBlocks",
		LongDesc:   "Use this command to list IpBlocks.",
		Example:    listIpBlockExample,
		PreCmdRun:  core.NoPreRun,
		CmdRun:     RunIpBlockList,
		InitClient: true,
	})

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
	get.AddStringFlag(cloudapiv6.ArgIpBlockId, cloudapiv6.ArgIdShort, "", cloudapiv6.IpBlockId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgIpBlockId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.IpBlocksIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})

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
	create.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for IpBlock creation to be executed")
	create.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option for Request for IpBlock creation [seconds]")

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
	update.AddStringFlag(cloudapiv6.ArgIpBlockId, cloudapiv6.ArgIdShort, "", cloudapiv6.IpBlockId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgIpBlockId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.IpBlocksIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "", "Name of the IpBlock")
	update.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for IpBlock update to be executed")
	update.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option for Request for IpBlock update [seconds]")

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
	deleteCmd.AddStringFlag(cloudapiv6.ArgIpBlockId, cloudapiv6.ArgIdShort, "", cloudapiv6.IpBlockId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgIpBlockId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.IpBlocksIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for IpBlock deletion to be executed")
	deleteCmd.AddBoolFlag(cloudapiv6.ArgAll, cloudapiv6.ArgAllShort, false, "delete all the IpBlocks.")
	deleteCmd.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option for Request for IpBlock deletion [seconds]")

	return ipblockCmd
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
	ipblocks, resp, err := c.CloudApiV6Services.IpBlocks().List()
	if resp != nil {
		c.Printer.Verbose(cloudapiv6.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getIpBlockPrint(nil, c, getIpBlocks(ipblocks)))
}

func RunIpBlockGet(c *core.CommandConfig) error {
	c.Printer.Verbose("Ip block with id: %v is getting...", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgIpBlockId)))
	i, resp, err := c.CloudApiV6Services.IpBlocks().Get(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgIpBlockId)))
	if resp != nil {
		c.Printer.Verbose(cloudapiv6.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getIpBlockPrint(nil, c, getIpBlock(i)))
}

func RunIpBlockCreate(c *core.CommandConfig) error {
	name := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName))
	loc := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLocation))
	size := viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgSize))
	c.Printer.Verbose("Properties set for creating the Ip block: Name: %v, Location: %v, Size: %v", name, loc, size)
	i, resp, err := c.CloudApiV6Services.IpBlocks().Create(name, loc, size)
	if resp != nil {
		c.Printer.Verbose("Request href: %v ", resp.Header.Get("location"))
		c.Printer.Verbose(cloudapiv6.RequestTimeMessage, resp.RequestTime)
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
	input := resources.IpBlockProperties{}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgName)) {
		name := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName))
		input.SetName(name)
		c.Printer.Verbose("Property Name set: %v", name)
	}
	i, resp, err := c.CloudApiV6Services.IpBlocks().Update(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgIpBlockId)),
		input,
	)
	if resp != nil {
		c.Printer.Verbose(cloudapiv6.RequestTimeMessage, resp.RequestTime)
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
	var resp *resources.Response
	var err error
	var ipBlocks resources.IpBlocks
	allFlag := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAll))
	if allFlag {
		fmt.Printf("IpBlocks to be deleted:\n")
		ipBlocks, resp, err = c.CloudApiV6Services.IpBlocks().List()
		if err != nil {
			return err
		}
		if ipBlocksItems, ok := ipBlocks.GetItemsOk(); ok && ipBlocksItems != nil {
			for _, dc := range *ipBlocksItems {
				if id, ok := dc.GetIdOk(); ok && id != nil {
					fmt.Printf("IpBlock Id: " + *id)
				}
				if properties, ok := dc.GetPropertiesOk(); ok && properties != nil {
					if name, ok := properties.GetNameOk(); ok && name != nil {
						fmt.Printf(" IpBlock Name: " + *name + "\n")
					}
				}
			}

			if err := utils.AskForConfirm(c.Stdin, c.Printer, "delete all the IpBlocks"); err != nil {
				return err
			}
			c.Printer.Verbose("Deleting all the IpBlocks...")
			for _, dc := range *ipBlocksItems {
				if id, ok := dc.GetIdOk(); ok && id != nil {
					c.Printer.Verbose("Deleting IpBlock with id: %v...", *id)
					resp, err = c.CloudApiV6Services.IpBlocks().Delete(*id)
					if resp != nil {
						c.Printer.Verbose(cloudapiv6.RequestTimeMessage, resp.RequestTime)
					}
					if err != nil {
						return err
					}
					if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
						return err
					}
				}
			}
		}
	} else {
		if err := utils.AskForConfirm(c.Stdin, c.Printer, "delete ipblock"); err != nil {
			return err
		}
		c.Printer.Verbose("Ip block with id: %v is deleting...", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgIpBlockId)))
		resp, err := c.CloudApiV6Services.IpBlocks().Delete(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgIpBlockId)))
		if resp != nil {
			c.Printer.Verbose(cloudapiv6.RequestTimeMessage, resp.RequestTime)
		}
		if err != nil {
			return err
		}

		if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
			return err
		}
	}
	return c.Printer.Print(getIpBlockPrint(resp, c, nil))
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
			r.WaitForRequest = viper.GetBool(core.GetFlagName(c.NS, config.ArgWaitForRequest))
		}
		if ipBlocks != nil {
			r.OutputJSON = ipBlocks
			r.KeyValue = getIpBlocksKVMaps(ipBlocks)
			r.Columns = getIpBlocksCols(core.GetGlobalFlagName(c.Resource, config.ArgCols), c.Printer.GetStderr())
		}
	}
	return r
}

func getIpBlocksCols(flagName string, outErr io.Writer) []string {
	var cols []string
	if viper.IsSet(flagName) {
		cols = viper.GetStringSlice(flagName)
	} else {
		return defaultIpBlockCols
	}
	columnsMap := map[string]string{
		"IpBlockId": "IpBlockId",
		"Name":      "Name",
		"Location":  "Location",
		"Size":      "Size",
		"Ips":       "Ips",
		"State":     "State",
	}
	var ipBlockCols []string
	for _, k := range cols {
		col := columnsMap[k]
		if col != "" {
			ipBlockCols = append(ipBlockCols, col)
		} else {
			clierror.CheckError(errors.New("unknown column "+k), outErr)
		}
	}
	return ipBlockCols
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
