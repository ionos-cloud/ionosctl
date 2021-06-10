package commands

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/resources"
	"github.com/ionos-cloud/ionosctl/pkg/utils"
	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
	"github.com/ionos-cloud/ionosctl/pkg/utils/printer"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	multierror "go.uber.org/multierr"
)

func flowlog() *core.Command {
	ctx := context.TODO()
	flowLogCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "flowlog",
			Aliases:          []string{"fl"},
			Short:            "FlowLog Operations",
			Long:             `The sub-commands of ` + "`" + `ionosctl flowlog` + "`" + ` allow you to create, list, get, delete FlowLogs on specific NICs.`,
			TraverseChildren: true,
		},
	}
	globalFlags := flowLogCmd.GlobalFlags()
	globalFlags.StringP(config.ArgDataCenterId, "", "", config.RequiredFlagDatacenterId)
	_ = viper.BindPFlag(core.GetGlobalFlagName(flowLogCmd.Name(), config.ArgDataCenterId), globalFlags.Lookup(config.ArgDataCenterId))
	_ = flowLogCmd.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	globalFlags.StringP(config.ArgServerId, "", "", config.RequiredFlagServerId)
	_ = viper.BindPFlag(core.GetGlobalFlagName(flowLogCmd.Name(), config.ArgServerId), globalFlags.Lookup(config.ArgServerId))
	_ = flowLogCmd.Command.RegisterFlagCompletionFunc(config.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getServersIds(os.Stderr, viper.GetString(core.GetGlobalFlagName(flowLogCmd.Name(), config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	globalFlags.StringP(config.ArgNicId, "", "", config.RequiredFlagNicId)
	_ = viper.BindPFlag(core.GetGlobalFlagName(flowLogCmd.Name(), config.ArgNicId), globalFlags.Lookup(config.ArgNicId))
	_ = flowLogCmd.Command.RegisterFlagCompletionFunc(config.ArgNicId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getNicsIds(os.Stderr,
			viper.GetString(core.GetGlobalFlagName(flowLogCmd.Name(), config.ArgDataCenterId)),
			viper.GetString(core.GetGlobalFlagName(flowLogCmd.Name(), config.ArgServerId)),
		), cobra.ShellCompDirectiveNoFileComp
	})
	globalFlags.StringSliceP(config.ArgCols, "", defaultFlowLogCols,
		fmt.Sprintf("Set of columns to be printed on output \nAvailable columns: %v", defaultFlowLogCols))
	_ = viper.BindPFlag(core.GetGlobalFlagName(flowLogCmd.Name(), config.ArgCols), globalFlags.Lookup(config.ArgCols))
	_ = flowLogCmd.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultFlowLogCols, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		List Command
	*/
	core.NewCommand(ctx, flowLogCmd, core.CommandBuilder{
		Namespace:  "flowlog",
		Resource:   "flowlog",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List FlowLogs",
		LongDesc:   "Use this command to get a list of FlowLogs from a specified NIC from a Server.\n\nRequired values to run command:\n\n* Data Center Id\n* Server Id\n*Nic Id",
		Example:    listFlowLogExample,
		PreCmdRun:  PreRunGlobalDcServerNicIds,
		CmdRun:     RunFlowLogList,
		InitClient: true,
	})

	/*
		Get Command
	*/
	get := core.NewCommand(ctx, flowLogCmd, core.CommandBuilder{
		Namespace:  "flowlog",
		Resource:   "flowlog",
		Verb:       "get",
		Aliases:    []string{"g"},
		ShortDesc:  "Get a FlowLog",
		LongDesc:   "Use this command to retrieve information of a specified FlowLog from a NIC.\n\nRequired values to run command:\n\n* Data Center Id\n* Server Id\n*Nic Id\n* FlowLog Id",
		Example:    getFlowLogExample,
		PreCmdRun:  PreRunGlobalDcServerNicIdsFlowLogId,
		CmdRun:     RunFlowLogGet,
		InitClient: true,
	})
	get.AddStringFlag(config.ArgFlowLogId, config.ArgIdShort, "", config.RequiredFlagFlowLogId)
	_ = get.Command.RegisterFlagCompletionFunc(config.ArgFlowLogId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getFlowLogsIds(os.Stderr,
			viper.GetString(core.GetGlobalFlagName(flowLogCmd.Name(), config.ArgDataCenterId)),
			viper.GetString(core.GetGlobalFlagName(flowLogCmd.Name(), config.ArgServerId)),
			viper.GetString(core.GetGlobalFlagName(flowLogCmd.Name(), config.ArgNicId))), cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Create Command
	*/
	create := core.NewCommand(ctx, flowLogCmd, core.CommandBuilder{
		Namespace: "flowlog",
		Resource:  "flowlog",
		Verb:      "create",
		Aliases:   []string{"c"},
		ShortDesc: "Create a FlowLog on a NIC",
		LongDesc: `Use this command to create a new FlowLog to the specified NIC.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option.

NOTE: Please disable the Firewall Active using before deleting the existing Bucket, so that FlowLogs know to not upload to a non-existing Bucket. To disable the Firewall, you can use ` + "`" + `ionosctl nic update` + "`" + ` command with ` + "`" + `--firewall-active=false` + "`" + ` option set.

Required values to run command:

* Data Center Id
* Server Id
* Nic Id 
* Name
* Direction
* Action
* Target S3 Bucket Name`,
		Example:    createFlowLogExample,
		PreCmdRun:  PreRunFlowLogCreate,
		CmdRun:     RunFlowLogCreate,
		InitClient: true,
	})
	create.AddStringFlag(config.ArgName, config.ArgNameShort, "", "The name for the FlowLog "+config.RequiredFlag)
	create.AddStringFlag(config.ArgAction, config.ArgActionShort, "", "Specifies the traffic Action pattern "+config.RequiredFlag)
	_ = create.Command.RegisterFlagCompletionFunc(config.ArgAction, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"ALL", "REJECTED", "ACCEPTED"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(config.ArgDirection, config.ArgDirectionShort, "", "Specifies the traffic Direction pattern "+config.RequiredFlag)
	_ = create.Command.RegisterFlagCompletionFunc(config.ArgDirection, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"BIDIRECTIONAL", "INGRESS", "EGRESS"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(config.ArgBucketName, "", "", "S3 Bucket name of an existing IONOS Cloud S3 Bucket "+config.RequiredFlag)
	create.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for Request for FlowLog creation to be executed")
	create.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option for Request for FlowLog creation [seconds]")

	/*
		Delete Command
	*/
	deleteCmd := core.NewCommand(ctx, flowLogCmd, core.CommandBuilder{
		Namespace: "flowlog",
		Resource:  "flowlog",
		Verb:      "delete",
		Aliases:   []string{"d"},
		ShortDesc: "Delete a FlowLog from a NIC",
		LongDesc: `Use this command to delete a specified FlowLog from a NIC.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option. You can force the command to execute without user input using ` + "`" + `--force` + "`" + ` option.

Required values to run command:

* Data Center Id
* Server Id
* Nic Id
* FlowLog Id`,
		Example:    deleteFlowLogExample,
		PreCmdRun:  PreRunGlobalDcServerNicIdsFlowLogId,
		CmdRun:     RunFlowLogDelete,
		InitClient: true,
	})
	deleteCmd.AddStringFlag(config.ArgFlowLogId, config.ArgIdShort, "", config.RequiredFlagFlowLogId)
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(config.ArgFlowLogId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getFlowLogsIds(os.Stderr,
			viper.GetString(core.GetGlobalFlagName(flowLogCmd.Name(), config.ArgDataCenterId)),
			viper.GetString(core.GetGlobalFlagName(flowLogCmd.Name(), config.ArgServerId)),
			viper.GetString(core.GetGlobalFlagName(flowLogCmd.Name(), config.ArgNicId))), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for Request for FlowLog deletion to be executed")
	deleteCmd.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option for Request for FlowLog deletion [seconds]")

	return flowLogCmd
}

func PreRunFlowLogCreate(c *core.PreCommandConfig) error {
	var result error
	if err := core.CheckRequiredGlobalFlags(c.Resource, config.ArgDataCenterId, config.ArgServerId, config.ArgNicId); err != nil {
		result = multierror.Append(result, err)
	}
	if err := core.CheckRequiredFlags(c.NS, config.ArgName, config.ArgAction, config.ArgDirection, config.ArgBucketName); err != nil {
		result = multierror.Append(result, err)
	}
	if result != nil {
		return result
	}
	return nil
}

func PreRunGlobalDcServerNicIdsFlowLogId(c *core.PreCommandConfig) error {
	var result error
	if err := core.CheckRequiredGlobalFlags(c.Resource, config.ArgDataCenterId, config.ArgServerId, config.ArgNicId); err != nil {
		result = multierror.Append(result, err)
	}
	if err := core.CheckRequiredFlags(c.NS, config.ArgFlowLogId); err != nil {
		result = multierror.Append(result, err)
	}
	if result != nil {
		return result
	}
	return nil
}

func RunFlowLogList(c *core.CommandConfig) error {
	flowLogs, _, err := c.FlowLogs().List(
		viper.GetString(core.GetGlobalFlagName(c.Resource, config.ArgDataCenterId)),
		viper.GetString(core.GetGlobalFlagName(c.Resource, config.ArgServerId)),
		viper.GetString(core.GetGlobalFlagName(c.Resource, config.ArgNicId)),
	)
	if err != nil {
		return err
	}
	return c.Printer.Print(getFlowLogPrint(nil, c, getFlowLogs(flowLogs)))
}

func RunFlowLogGet(c *core.CommandConfig) error {
	flowLog, _, err := c.FlowLogs().Get(
		viper.GetString(core.GetGlobalFlagName(c.Resource, config.ArgDataCenterId)),
		viper.GetString(core.GetGlobalFlagName(c.Resource, config.ArgServerId)),
		viper.GetString(core.GetGlobalFlagName(c.Resource, config.ArgNicId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgFlowLogId)),
	)
	if err != nil {
		return err
	}
	return c.Printer.Print(getFlowLogPrint(nil, c, getFlowLog(flowLog)))
}

func RunFlowLogCreate(c *core.CommandConfig) error {
	properties := getFlowLogPropertiesSet(c)
	input := resources.FlowLog{
		FlowLog: ionoscloud.FlowLog{
			Properties: &properties.FlowLogProperties,
		},
	}
	flowLog, resp, err := c.FlowLogs().Create(
		viper.GetString(core.GetGlobalFlagName(c.Resource, config.ArgDataCenterId)),
		viper.GetString(core.GetGlobalFlagName(c.Resource, config.ArgServerId)),
		viper.GetString(core.GetGlobalFlagName(c.Resource, config.ArgNicId)),
		input,
	)
	if err != nil {
		return err
	}

	if err = utils.WaitForRequest(c, printer.GetRequestPath(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getFlowLogPrint(resp, c, getFlowLog(flowLog)))
}

func RunFlowLogDelete(c *core.CommandConfig) error {
	if err := utils.AskForConfirm(c.Stdin, c.Printer, "delete flow log"); err != nil {
		return err
	}
	resp, err := c.FlowLogs().Delete(
		viper.GetString(core.GetGlobalFlagName(c.Resource, config.ArgDataCenterId)),
		viper.GetString(core.GetGlobalFlagName(c.Resource, config.ArgServerId)),
		viper.GetString(core.GetGlobalFlagName(c.Resource, config.ArgNicId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgFlowLogId)),
	)
	if err != nil {
		return err
	}

	if err = utils.WaitForRequest(c, printer.GetRequestPath(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getFlowLogPrint(resp, c, nil))
}

// Get FlowLog Properties set used for create and update commands
func getFlowLogPropertiesSet(c *core.CommandConfig) resources.FlowLogProperties {
	properties := resources.FlowLogProperties{}
	if viper.IsSet(core.GetFlagName(c.NS, config.ArgName)) {
		properties.SetName(viper.GetString(core.GetFlagName(c.NS, config.ArgName)))
	}
	if viper.IsSet(core.GetFlagName(c.NS, config.ArgAction)) {
		properties.SetAction(strings.ToUpper(viper.GetString(core.GetFlagName(c.NS, config.ArgAction))))
	}
	if viper.IsSet(core.GetFlagName(c.NS, config.ArgDirection)) {
		properties.SetDirection(strings.ToUpper(viper.GetString(core.GetFlagName(c.NS, config.ArgDirection))))
	}
	if viper.IsSet(core.GetFlagName(c.NS, config.ArgBucketName)) {
		properties.SetBucket(viper.GetString(core.GetFlagName(c.NS, config.ArgBucketName)))
	}
	return properties
}

// Output Printing

var defaultFlowLogCols = []string{"FlowLogId", "Name", "Action", "Direction", "Bucket", "State"}

type FlowLogPrint struct {
	FlowLogId string `json:"FlowLogId,omitempty"`
	Name      string `json:"Name,omitempty"`
	Action    string `json:"Action,omitempty"`
	Direction string `json:"Direction,omitempty"`
	Bucket    string `json:"Bucket,omitempty"`
	State     string `json:"State,omitempty"`
}

func getFlowLogPrint(resp *resources.Response, c *core.CommandConfig, rule []resources.FlowLog) printer.Result {
	var r printer.Result
	if c != nil {
		if resp != nil {
			r.ApiResponse = resp
			r.Resource = c.Resource
			r.Verb = c.Verb
			r.WaitForRequest = viper.GetBool(core.GetFlagName(c.NS, config.ArgWaitForRequest))
		}
		if rule != nil {
			r.OutputJSON = rule
			r.KeyValue = getFlowLogsKVMaps(rule)
			r.Columns = getFlowLogsCols(core.GetGlobalFlagName(c.Resource, config.ArgCols), c.Printer.GetStderr())
		}
	}
	return r
}

func getFlowLogsCols(flagName string, outErr io.Writer) []string {
	if viper.IsSet(flagName) && len(viper.GetStringSlice(flagName)) > 0 {
		var flowLogCols []string
		columnsMap := map[string]string{
			"FlowLogId": "FlowLogId",
			"Name":      "Name",
			"Action":    "Action",
			"Direction": "Direction",
			"Bucket":    "Bucket",
			"State":     "State",
		}
		for _, k := range viper.GetStringSlice(flagName) {
			col := columnsMap[k]
			if col != "" {
				flowLogCols = append(flowLogCols, col)
			} else {
				clierror.CheckError(errors.New("unknown column "+k), outErr)
			}
		}
		return flowLogCols
	} else {
		return defaultFlowLogCols
	}
}

func getFlowLogs(flowLogs resources.FlowLogs) []resources.FlowLog {
	ls := make([]resources.FlowLog, 0)
	if items, ok := flowLogs.GetItemsOk(); ok && items != nil {
		for _, s := range *items {
			ls = append(ls, resources.FlowLog{FlowLog: s})
		}
	}
	return ls
}

func getFlowLog(flowLog *resources.FlowLog) []resources.FlowLog {
	ss := make([]resources.FlowLog, 0)
	if flowLog != nil {
		ss = append(ss, resources.FlowLog{FlowLog: flowLog.FlowLog})
	}
	return ss
}

func getFlowLogsKVMaps(ls []resources.FlowLog) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(ls))
	if len(ls) > 0 {
		for _, l := range ls {
			o := getFlowLogKVMap(l)
			out = append(out, o)
		}
	}
	return out
}

func getFlowLogKVMap(l resources.FlowLog) map[string]interface{} {
	var flowLogPrint FlowLogPrint
	if id, ok := l.GetIdOk(); ok && id != nil {
		flowLogPrint.FlowLogId = *id
	}
	if properties, ok := l.GetPropertiesOk(); ok && properties != nil {
		if name, ok := properties.GetNameOk(); ok && name != nil {
			flowLogPrint.Name = *name
		}
		if action, ok := properties.GetActionOk(); ok && action != nil {
			flowLogPrint.Action = *action
		}
		if direction, ok := properties.GetDirectionOk(); ok && direction != nil {
			flowLogPrint.Direction = *direction
		}
		if bucket, ok := properties.GetBucketOk(); ok && bucket != nil {
			flowLogPrint.Bucket = *bucket
		}
	}
	if metadata, ok := l.GetMetadataOk(); ok && metadata != nil {
		if state, ok := metadata.GetStateOk(); ok && state != nil {
			flowLogPrint.State = *state
		}
	}
	return structs.Map(flowLogPrint)
}

func getFlowLogsIds(outErr io.Writer, datacenterId, serverId, nicId string) []string {
	err := config.Load()
	clierror.CheckError(err, outErr)
	clientSvc, err := resources.NewClientService(
		viper.GetString(config.Username),
		viper.GetString(config.Password),
		viper.GetString(config.Token),
		viper.GetString(config.ArgServerUrl),
	)
	clierror.CheckError(err, outErr)
	flowLogSvc := resources.NewFlowLogService(clientSvc.Get(), context.TODO())
	flowLogs, _, err := flowLogSvc.List(datacenterId, serverId, nicId)
	clierror.CheckError(err, outErr)
	flowLogsIds := make([]string, 0)
	if items, ok := flowLogs.FlowLogs.GetItemsOk(); ok && items != nil {
		for _, item := range *items {
			if itemId, ok := item.GetIdOk(); ok && itemId != nil {
				flowLogsIds = append(flowLogsIds, *itemId)
			}
		}
	} else {
		return nil
	}
	return flowLogsIds
}
