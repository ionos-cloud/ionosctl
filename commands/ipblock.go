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

func ipblock() *builder.Command {
	ctx := context.TODO()
	ipblockCmd := &builder.Command{
		Command: &cobra.Command{
			Use:              "ipblock",
			Short:            "IpBlock Operations",
			Long:             `The sub-commands of ` + "`" + `ionosctl ipblock` + "`" + ` allow you to create/reserve, list, get, update, delete IpBlocks.`,
			TraverseChildren: true,
		},
	}
	globalFlags := ipblockCmd.GlobalFlags()
	globalFlags.StringSlice(config.ArgCols, defaultIpBlockCols, "Columns to be printed in the standard output")
	_ = viper.BindPFlag(builder.GetGlobalFlagName(ipblockCmd.Name(), config.ArgCols), globalFlags.Lookup(config.ArgCols))

	/*
		List Command
	*/
	builder.NewCommand(ctx, ipblockCmd, noPreRun, RunIpBlockList, "list", "List IpBlocks",
		"Use this command to list IpBlocks.",
		listIpBlockExample, true)

	/*
		Get Command
	*/
	get := builder.NewCommand(ctx, ipblockCmd, PreRunIpBlockId, RunIpBlockGet, "get", "Get an IpBlock",
		"Use this command to retrieve the attributes of a specific IpBlock.\n\nRequired values to run command:\n\n* IpBlock Id",
		getIpBlockExample, true)
	get.AddStringFlag(config.ArgIpBlockId, "", "", config.RequiredFlagIpBlockId)
	_ = get.Command.RegisterFlagCompletionFunc(config.ArgIpBlockId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getIpBlocksIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Create Command
	*/
	create := builder.NewCommand(ctx, ipblockCmd, PreRunIpBlockLocation, RunIpBlockCreate, "create", "Create/Reserve an IpBlock",
		`Use this command to create/reserve an IpBlock in a specified location that can be used by resources within any Virtual Data Centers provisioned in that same location. An IpBlock consists of one or more static IP addresses. The name, size of the IpBlock can be set.

You can wait for the Request to be executed using `+"`"+`--wait-for-request`+"`"+` option.

Required values to run command:

* IpBlock Location`, createIpBlockExample, true)
	create.AddStringFlag(config.ArgIpBlockName, "", "", "Name of the IpBlock")
	create.AddStringFlag(config.ArgIpBlockLocation, "", "", "Location of the IpBlock "+config.RequiredFlag)
	_ = create.Command.RegisterFlagCompletionFunc(config.ArgIpBlockLocation, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getLocationIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddIntFlag(config.ArgIpBlockSize, "", 2, "Size of the IpBlock")
	create.AddBoolFlag(config.ArgWaitForRequest, "", config.DefaultWait, "Wait for the Request for IpBlock creation to be executed")
	create.AddIntFlag(config.ArgTimeout, "", config.DefaultTimeoutSeconds, "Timeout option for Request for IpBlock creation [seconds]")

	/*
		Update Command
	*/
	update := builder.NewCommand(ctx, ipblockCmd, PreRunIpBlockId, RunIpBlockUpdate, "update", "Update an IpBlock",
		`Use this command to update the properties of an existing IpBlock.

You can wait for the Request to be executed using `+"`"+`--wait-for-request`+"`"+` option.

Required values to run command:

* IpBlock Id`, updateIpBlockExample, true)
	update.AddStringFlag(config.ArgIpBlockId, "", "", config.RequiredFlagIpBlockId)
	_ = update.Command.RegisterFlagCompletionFunc(config.ArgIpBlockId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getIpBlocksIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(config.ArgIpBlockName, "", "", "Name of the IpBlock")
	update.AddBoolFlag(config.ArgWaitForRequest, "", config.DefaultWait, "Wait for the Request for IpBlock update to be executed")
	update.AddIntFlag(config.ArgTimeout, "", config.DefaultTimeoutSeconds, "Timeout option for Request for IpBlock update [seconds]")

	/*
		Delete Command
	*/
	deleteCmd := builder.NewCommand(ctx, ipblockCmd, PreRunIpBlockId, RunIpBlockDelete, "delete", "Delete an IpBlock",
		`Use this command to delete a specified IpBlock.

You can wait for the Request to be executed using `+"`"+`--wait-for-request`+"`"+` option. You can force the command to execute without user input using `+"`"+`--force`+"`"+` option.

Required values to run command:

* IpBlock Id`, deleteIpBlockExample, true)
	deleteCmd.AddStringFlag(config.ArgIpBlockId, "", "", config.RequiredFlagIpBlockId)
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(config.ArgIpBlockId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getIpBlocksIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddBoolFlag(config.ArgWaitForRequest, "", config.DefaultWait, "Wait for the Request for IpBlock deletion to be executed")
	deleteCmd.AddIntFlag(config.ArgTimeout, "", config.DefaultTimeoutSeconds, "Timeout option for Request for IpBlock deletion [seconds]")

	return ipblockCmd
}

func PreRunIpBlockLocation(c *builder.PreCommandConfig) error {
	return builder.CheckRequiredFlags(c.ParentName, c.Name, config.ArgIpBlockLocation)
}

func PreRunIpBlockId(c *builder.PreCommandConfig) error {
	return builder.CheckRequiredFlags(c.ParentName, c.Name, config.ArgIpBlockId)
}

func RunIpBlockList(c *builder.CommandConfig) error {
	ipblocks, _, err := c.IpBlocks().List()
	if err != nil {
		return err
	}
	return c.Printer.Print(getIpBlockPrint(nil, c, getIpBlocks(ipblocks)))
}

func RunIpBlockGet(c *builder.CommandConfig) error {
	i, _, err := c.IpBlocks().Get(viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgIpBlockId)))
	if err != nil {
		return err
	}
	return c.Printer.Print(getIpBlockPrint(nil, c, getIpBlock(i)))
}

func RunIpBlockCreate(c *builder.CommandConfig) error {
	i, resp, err := c.IpBlocks().Create(
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgIpBlockName)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgIpBlockLocation)),
		viper.GetInt32(builder.GetFlagName(c.ParentName, c.Name, config.ArgIpBlockSize)),
	)
	if err != nil {
		return err
	}

	if err = utils.WaitForRequest(c, printer.GetRequestPath(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getIpBlockPrint(resp, c, getIpBlock(i)))
}

func RunIpBlockUpdate(c *builder.CommandConfig) error {
	input := resources.IpBlockProperties{}
	if viper.IsSet(builder.GetFlagName(c.ParentName, c.Name, config.ArgIpBlockName)) {
		input.SetName(viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgIpBlockName)))
	}
	i, resp, err := c.IpBlocks().Update(
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgIpBlockId)),
		input,
	)
	if err != nil {
		return err
	}

	if err = utils.WaitForRequest(c, printer.GetRequestPath(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getIpBlockPrint(resp, c, getIpBlock(i)))
}

func RunIpBlockDelete(c *builder.CommandConfig) error {
	if err := utils.AskForConfirm(c.Stdin, c.Printer, "delete ipblock"); err != nil {
		return err
	}
	resp, err := c.IpBlocks().Delete(viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgIpBlockId)))
	if err != nil {
		return err
	}

	if err = utils.WaitForRequest(c, printer.GetRequestPath(resp)); err != nil {
		return err
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

func getIpBlockPrint(resp *resources.Response, c *builder.CommandConfig, ipBlocks []resources.IpBlock) printer.Result {
	r := printer.Result{}
	if c != nil {
		if resp != nil {
			r.ApiResponse = resp
			r.Resource = c.ParentName
			r.Verb = c.Name
			r.WaitForRequest = viper.GetBool(builder.GetFlagName(c.ParentName, c.Name, config.ArgWaitForRequest))
		}
		if ipBlocks != nil {
			r.OutputJSON = ipBlocks
			r.KeyValue = getIpBlocksKVMaps(ipBlocks)
			r.Columns = getIpBlocksCols(builder.GetGlobalFlagName(c.ParentName, config.ArgCols), c.Printer.GetStderr())
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

func getIpBlocksIds(outErr io.Writer) []string {
	err := config.Load()
	clierror.CheckError(err, outErr)
	clientSvc, err := resources.NewClientService(
		viper.GetString(config.Username),
		viper.GetString(config.Password),
		viper.GetString(config.Token),
		viper.GetString(config.ArgServerUrl),
	)
	clierror.CheckError(err, outErr)
	ipBlockSvc := resources.NewIpBlockService(clientSvc.Get(), context.TODO())
	ipBlocks, _, err := ipBlockSvc.List()
	clierror.CheckError(err, outErr)
	ssIds := make([]string, 0)
	if items, ok := ipBlocks.IpBlocks.GetItemsOk(); ok && items != nil {
		for _, item := range *items {
			if itemId, ok := item.GetIdOk(); ok && itemId != nil {
				ssIds = append(ssIds, *itemId)
			}
		}
	} else {
		return nil
	}
	return ssIds
}
