package commands

import (
	"context"
	"errors"
	"io"
	"os"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/resources"
	"github.com/ionos-cloud/ionosctl/pkg/utils"
	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
	"github.com/ionos-cloud/ionosctl/pkg/utils/printer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ipblock() *core.Command {
	ctx := context.TODO()
	ipblockCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "ipblock",
			Short:            "IpBlock Operations",
			Long:             `The sub-commands of ` + "`" + `ionosctl ipblock` + "`" + ` allow you to create/reserve, list, get, update, delete IpBlocks.`,
			TraverseChildren: true,
		},
	}
	globalFlags := ipblockCmd.GlobalFlags()
	globalFlags.StringSlice(config.ArgCols, defaultIpBlockCols, "Columns to be printed in the standard output")
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
		ShortDesc:  "List IpBlocks",
		LongDesc:   "Use this command to list IpBlocks.",
		Example:    listIpBlockExample,
		PreCmdRun:  noPreRun,
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
		ShortDesc:  "Get an IpBlock",
		LongDesc:   "Use this command to retrieve the attributes of a specific IpBlock.\n\nRequired values to run command:\n\n* IpBlock Id",
		Example:    getIpBlockExample,
		PreCmdRun:  PreRunIpBlockId,
		CmdRun:     RunIpBlockGet,
		InitClient: true,
	})
	get.AddStringFlag(config.ArgIpBlockId, "", "", config.RequiredFlagIpBlockId)
	_ = get.Command.RegisterFlagCompletionFunc(config.ArgIpBlockId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getIpBlocksIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Create Command
	*/
	create := core.NewCommand(ctx, ipblockCmd, core.CommandBuilder{
		Namespace: "ipblock",
		Resource:  "ipblock",
		Verb:      "create",
		ShortDesc: "Create/Reserve an IpBlock",
		LongDesc: `Use this command to create/reserve an IpBlock in a specified location that can be used by resources within any Virtual Data Centers provisioned in that same location. An IpBlock consists of one or more static IP addresses. The name, size of the IpBlock can be set.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option.

Required values to run command:

* IpBlock Location`,
		Example:    createIpBlockExample,
		PreCmdRun:  PreRunIpBlockLocation,
		CmdRun:     RunIpBlockCreate,
		InitClient: true,
	})
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
	update := core.NewCommand(ctx, ipblockCmd, core.CommandBuilder{
		Namespace: "ipblock",
		Resource:  "ipblock",
		Verb:      "update",
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
	deleteCmd := core.NewCommand(ctx, ipblockCmd, core.CommandBuilder{
		Namespace: "ipblock",
		Resource:  "ipblock",
		Verb:      "delete",
		ShortDesc: "Delete an IpBlock",
		LongDesc: `Use this command to delete a specified IpBlock.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option. You can force the command to execute without user input using ` + "`" + `--force` + "`" + ` option.

Required values to run command:

* IpBlock Id`,
		Example:    deleteIpBlockExample,
		PreCmdRun:  PreRunIpBlockId,
		CmdRun:     RunIpBlockDelete,
		InitClient: true,
	})
	deleteCmd.AddStringFlag(config.ArgIpBlockId, "", "", config.RequiredFlagIpBlockId)
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(config.ArgIpBlockId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getIpBlocksIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddBoolFlag(config.ArgWaitForRequest, "", config.DefaultWait, "Wait for the Request for IpBlock deletion to be executed")
	deleteCmd.AddIntFlag(config.ArgTimeout, "", config.DefaultTimeoutSeconds, "Timeout option for Request for IpBlock deletion [seconds]")

	ipblockCmd.AddCommand(ipconsumers())

	return ipblockCmd
}

func PreRunIpBlockLocation(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.NS, config.ArgIpBlockLocation)
}

func PreRunIpBlockId(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.NS, config.ArgIpBlockId)
}

func RunIpBlockList(c *core.CommandConfig) error {
	ipblocks, _, err := c.IpBlocks().List()
	if err != nil {
		return err
	}
	return c.Printer.Print(getIpBlockPrint(nil, c, getIpBlocks(ipblocks)))
}

func RunIpBlockGet(c *core.CommandConfig) error {
	i, _, err := c.IpBlocks().Get(viper.GetString(core.GetFlagName(c.NS, config.ArgIpBlockId)))
	if err != nil {
		return err
	}
	return c.Printer.Print(getIpBlockPrint(nil, c, getIpBlock(i)))
}

func RunIpBlockCreate(c *core.CommandConfig) error {
	i, resp, err := c.IpBlocks().Create(
		viper.GetString(core.GetFlagName(c.NS, config.ArgIpBlockName)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgIpBlockLocation)),
		viper.GetInt32(core.GetFlagName(c.NS, config.ArgIpBlockSize)),
	)
	if err != nil {
		return err
	}

	if err = utils.WaitForRequest(c, printer.GetRequestPath(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getIpBlockPrint(resp, c, getIpBlock(i)))
}

func RunIpBlockUpdate(c *core.CommandConfig) error {
	input := resources.IpBlockProperties{}
	if viper.IsSet(core.GetFlagName(c.NS, config.ArgIpBlockName)) {
		input.SetName(viper.GetString(core.GetFlagName(c.NS, config.ArgIpBlockName)))
	}
	i, resp, err := c.IpBlocks().Update(
		viper.GetString(core.GetFlagName(c.NS, config.ArgIpBlockId)),
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

func RunIpBlockDelete(c *core.CommandConfig) error {
	if err := utils.AskForConfirm(c.Stdin, c.Printer, "delete ipblock"); err != nil {
		return err
	}
	resp, err := c.IpBlocks().Delete(viper.GetString(core.GetFlagName(c.NS, config.ArgIpBlockId)))
	if err != nil {
		return err
	}

	if err = utils.WaitForRequest(c, printer.GetRequestPath(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getIpBlockPrint(resp, c, nil))
}

// IpBlock Consumer Commands

func ipconsumers() *core.Command {
	ctx := context.TODO()
	resourceCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "consumers",
			Short:            "Ip Consumers Operations",
			Long:             `The sub-command of ` + "`" + `ionosctl ipblock consumers` + "`" + ` allows you to list information about where each IP address from an IpBlock is being used.`,
			TraverseChildren: true,
		},
	}
	globalFlags := resourceCmd.GlobalFlags()
	globalFlags.StringSlice(config.ArgCols, defaultIpConsumerCols, "Columns to be printed in the standard output")
	_ = viper.BindPFlag(core.GetGlobalFlagName(resourceCmd.Name(), config.ArgCols), globalFlags.Lookup(config.ArgCols))
	_ = resourceCmd.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allIpConsumerCols, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		List Consumers Command
	*/
	listResources := core.NewCommand(ctx, resourceCmd, core.CommandBuilder{
		Namespace:  "ipblock",
		Resource:   "consumers",
		Verb:       "list",
		ShortDesc:  "List IpConsumers",
		LongDesc:   "Use this command to get a list of Resources where each IP address from an IpBlock is being used.\n\nRequired values to run command:\n\n* IpBlock Id",
		Example:    consumersListIpBlockExample,
		PreCmdRun:  PreRunIpBlockId,
		CmdRun:     RunIpConsumersList,
		InitClient: true,
	})
	listResources.AddStringFlag(config.ArgIpBlockId, "", "", config.RequiredFlagIpBlockId)
	_ = listResources.Command.RegisterFlagCompletionFunc(config.ArgIpBlockId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getIpBlocksIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})

	return resourceCmd
}

func RunIpConsumersList(c *core.CommandConfig) error {
	ipsConsumers := make([]resources.IpConsumer, 0)
	ipBlock, _, err := c.IpBlocks().Get(viper.GetString(core.GetFlagName(c.NS, config.ArgIpBlockId)))
	if err != nil {
		return err
	}
	if properties, ok := ipBlock.GetPropertiesOk(); ok && properties != nil {
		if ipCons, ok := properties.GetIpConsumersOk(); ok && ipCons != nil {
			for _, ip := range *ipCons {
				ipsConsumers = append(ipsConsumers, resources.IpConsumer{IpConsumer: ip})
			}
		}
	}
	return c.Printer.Print(getIpConsumerPrint(c, ipsConsumers))
}

// Output Printing

var (
	defaultIpConsumerCols = []string{"Ip", "NicId", "ServerId", "DatacenterId", "K8sNodePoolId", "K8sClusterId"}
	allIpConsumerCols     = []string{"Ip", "Mac", "NicId", "ServerId", "ServerName", "DatacenterId", "DatacenterName", "K8sNodePoolId", "K8sClusterId"}
)

type IpConsumerPrint struct {
	Ip             string `json:"Ip,omitempty"`
	Mac            string `json:"Mac,omitempty"`
	NicId          string `json:"NicId,omitempty"`
	ServerId       string `json:"ServerId,omitempty"`
	ServerName     string `json:"ServerName,omitempty"`
	DatacenterId   string `json:"DatacenterId,omitempty"`
	DatacenterName string `json:"DatacenterName,omitempty"`
	K8sNodePoolId  string `json:"K8sNodePoolId,omitempty"`
	K8sClusterId   string `json:"K8sClusterId,omitempty"`
}

func getIpConsumerPrint(c *core.CommandConfig, groups []resources.IpConsumer) printer.Result {
	r := printer.Result{}
	if c != nil {
		if groups != nil {
			r.OutputJSON = groups
			r.KeyValue = getIpConsumersKVMaps(groups)
			r.Columns = getIpConsumerCols(core.GetGlobalFlagName(c.Resource, config.ArgCols), c.Printer.GetStderr())
		}
	}
	return r
}

func getIpConsumerCols(flagName string, outErr io.Writer) []string {
	if viper.IsSet(flagName) {
		var groupCols []string
		columnsMap := map[string]string{
			"Ip":             "Ip",
			"Mac":            "Mac",
			"NicId":          "NicId",
			"ServerId":       "ServerId",
			"ServerName":     "ServerName",
			"DatacenterId":   "DatacenterId",
			"DatacenterName": "DatacenterName",
			"K8sNodePoolId":  "K8sNodePoolId",
			"K8sClusterId":   "K8sClusterId",
		}
		for _, k := range viper.GetStringSlice(flagName) {
			col := columnsMap[k]
			if col != "" {
				groupCols = append(groupCols, col)
			} else {
				clierror.CheckError(errors.New("unknown column "+k), outErr)
			}
		}
		return groupCols
	} else {
		return defaultIpConsumerCols
	}
}

func getIpConsumersKVMaps(rs []resources.IpConsumer) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(rs))
	for _, r := range rs {
		var rPrint IpConsumerPrint
		if ip, ok := r.GetIpOk(); ok && ip != nil {
			rPrint.Ip = *ip
		}
		if mac, ok := r.GetMacOk(); ok && mac != nil {
			rPrint.Mac = *mac
		}
		if nicId, ok := r.GetNicIdOk(); ok && nicId != nil {
			rPrint.NicId = *nicId
		}
		if serverId, ok := r.GetServerIdOk(); ok && serverId != nil {
			rPrint.ServerId = *serverId
		}
		if serverName, ok := r.GetServerNameOk(); ok && serverName != nil {
			rPrint.ServerName = *serverName
		}
		if datacenterId, ok := r.GetDatacenterIdOk(); ok && datacenterId != nil {
			rPrint.DatacenterId = *datacenterId
		}
		if datacenterName, ok := r.GetDatacenterNameOk(); ok && datacenterName != nil {
			rPrint.DatacenterName = *datacenterName
		}
		if nodepoolId, ok := r.GetK8sNodePoolUuidOk(); ok && nodepoolId != nil {
			rPrint.K8sNodePoolId = *nodepoolId
		}
		if clusterId, ok := r.GetK8sClusterUuidOk(); ok && clusterId != nil {
			rPrint.K8sClusterId = *clusterId
		}
		o := structs.Map(rPrint)
		out = append(out, o)
	}
	return out
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
