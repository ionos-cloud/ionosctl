package commands

import (
	"context"
	"errors"
	"github.com/ionos-cloud/ionosctl/commands/cloudapi-v6/completer"
	"github.com/ionos-cloud/ionosctl/commands/cloudapi-v6/waiter"
	cloudapi_v6 "github.com/ionos-cloud/ionosctl/services/cloudapi-v6"
	"io"
	"os"
	"strings"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/internal/config"
	"github.com/ionos-cloud/ionosctl/internal/core"
	"github.com/ionos-cloud/ionosctl/internal/printer"
	"github.com/ionos-cloud/ionosctl/internal/utils"
	"github.com/ionos-cloud/ionosctl/internal/utils/clierror"
	"github.com/ionos-cloud/ionosctl/services/cloudapi-v6/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func FlowlogCmd() *core.Command {
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
	globalFlags.StringSliceP(config.ArgCols, "", defaultFlowLogCols, printer.ColsMessage(defaultFlowLogCols))
	_ = viper.BindPFlag(core.GetGlobalFlagName(flowLogCmd.Name(), config.ArgCols), globalFlags.Lookup(config.ArgCols))
	_ = flowLogCmd.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultFlowLogCols, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		List Command
	*/
	list := core.NewCommand(ctx, flowLogCmd, core.CommandBuilder{
		Namespace:  "flowlog",
		Resource:   "flowlog",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List FlowLogs",
		LongDesc:   "Use this command to get a list of FlowLogs from a specified NIC from a Server.\n\nRequired values to run command:\n\n* Data Center Id\n* Server Id\n* Nic Id",
		Example:    listFlowLogExample,
		PreCmdRun:  PreRunDcServerNicIds,
		CmdRun:     RunFlowLogList,
		InitClient: true,
	})
	list.AddStringFlag(cloudapi_v6.ArgDataCenterId, "", "", cloudapi_v6.DatacenterId, core.RequiredFlagOption())
	_ = list.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddStringFlag(cloudapi_v6.ArgServerId, "", "", cloudapi_v6.ServerId, core.RequiredFlagOption())
	_ = list.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ServersIds(os.Stderr, viper.GetString(core.GetFlagName(list.NS, cloudapi_v6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddStringFlag(cloudapi_v6.ArgNicId, "", "", cloudapi_v6.NicId, core.RequiredFlagOption())
	_ = list.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgNicId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.NicsIds(os.Stderr, viper.GetString(core.GetFlagName(list.NS, cloudapi_v6.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(list.NS, cloudapi_v6.ArgServerId))), cobra.ShellCompDirectiveNoFileComp
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
		LongDesc:   "Use this command to retrieve information of a specified FlowLog from a NIC.\n\nRequired values to run command:\n\n* Data Center Id\n* Server Id\n* Nic Id\n* FlowLog Id",
		Example:    getFlowLogExample,
		PreCmdRun:  PreRunDcServerNicFlowLogIds,
		CmdRun:     RunFlowLogGet,
		InitClient: true,
	})
	get.AddStringFlag(cloudapi_v6.ArgDataCenterId, "", "", cloudapi_v6.DatacenterId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddStringFlag(cloudapi_v6.ArgServerId, "", "", cloudapi_v6.ServerId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ServersIds(os.Stderr, viper.GetString(core.GetFlagName(get.NS, cloudapi_v6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddStringFlag(cloudapi_v6.ArgNicId, "", "", cloudapi_v6.NicId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgNicId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.NicsIds(os.Stderr, viper.GetString(core.GetFlagName(get.NS, cloudapi_v6.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(get.NS, cloudapi_v6.ArgServerId))), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddStringFlag(cloudapi_v6.ArgFlowLogId, cloudapi_v6.ArgIdShort, "", cloudapi_v6.FlowLogId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgFlowLogId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.FlowLogsIds(os.Stderr,
			viper.GetString(core.GetFlagName(get.NS, cloudapi_v6.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(get.NS, cloudapi_v6.ArgServerId)),
			viper.GetString(core.GetFlagName(get.NS, cloudapi_v6.ArgNicId))), cobra.ShellCompDirectiveNoFileComp
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

NOTE: Please disable the FlowLog before deleting the existing Bucket.

Required values to run command:

* Data Center Id
* Server Id
* Nic Id
* Target S3 Bucket Name`,
		Example:    createFlowLogExample,
		PreCmdRun:  PreRunFlowLogCreate,
		CmdRun:     RunFlowLogCreate,
		InitClient: true,
	})
	create.AddStringFlag(cloudapi_v6.ArgDataCenterId, "", "", cloudapi_v6.DatacenterId, core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(cloudapi_v6.ArgServerId, "", "", cloudapi_v6.ServerId, core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ServersIds(os.Stderr, viper.GetString(core.GetFlagName(create.NS, cloudapi_v6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(cloudapi_v6.ArgNicId, "", "", cloudapi_v6.NicId, core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgNicId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.NicsIds(os.Stderr, viper.GetString(core.GetFlagName(create.NS, cloudapi_v6.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(create.NS, cloudapi_v6.ArgServerId))), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(cloudapi_v6.ArgName, cloudapi_v6.ArgNameShort, "Unnamed FlowLog", "The name for the FlowLog")
	create.AddStringFlag(cloudapi_v6.ArgAction, cloudapi_v6.ArgActionShort, "ALL", "Specifies the traffic Action pattern")
	_ = create.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgAction, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"ALL", "REJECTED", "ACCEPTED"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(cloudapi_v6.ArgDirection, cloudapi_v6.ArgDirectionShort, "BIDIRECTIONAL", "Specifies the traffic Direction pattern")
	_ = create.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgDirection, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"BIDIRECTIONAL", "INGRESS", "EGRESS"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(cloudapi_v6.ArgS3Bucket, cloudapi_v6.ArgS3BucketShort, "", "S3 Bucket name of an existing IONOS Cloud S3 Bucket", core.RequiredFlagOption())
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
		PreCmdRun:  PreRunDcServerNicFlowLogIds,
		CmdRun:     RunFlowLogDelete,
		InitClient: true,
	})
	deleteCmd.AddStringFlag(cloudapi_v6.ArgDataCenterId, "", "", cloudapi_v6.DatacenterId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddStringFlag(cloudapi_v6.ArgServerId, "", "", cloudapi_v6.ServerId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ServersIds(os.Stderr, viper.GetString(core.GetFlagName(deleteCmd.NS, cloudapi_v6.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddStringFlag(cloudapi_v6.ArgNicId, "", "", cloudapi_v6.NicId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgNicId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.NicsIds(os.Stderr, viper.GetString(core.GetFlagName(deleteCmd.NS, cloudapi_v6.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(deleteCmd.NS, cloudapi_v6.ArgServerId))), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddStringFlag(cloudapi_v6.ArgFlowLogId, cloudapi_v6.ArgIdShort, "", cloudapi_v6.FlowLogId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapi_v6.ArgFlowLogId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.FlowLogsIds(os.Stderr, viper.GetString(core.GetFlagName(deleteCmd.NS, cloudapi_v6.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(deleteCmd.NS, cloudapi_v6.ArgServerId)),
			viper.GetString(core.GetFlagName(deleteCmd.NS, cloudapi_v6.ArgNicId))), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for Request for FlowLog deletion to be executed")
	deleteCmd.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option for Request for FlowLog deletion [seconds]")

	return flowLogCmd
}

func PreRunFlowLogCreate(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapi_v6.ArgDataCenterId, cloudapi_v6.ArgServerId, cloudapi_v6.ArgNicId, cloudapi_v6.ArgS3Bucket)
}

func PreRunDcServerNicFlowLogIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapi_v6.ArgDataCenterId, cloudapi_v6.ArgServerId, cloudapi_v6.ArgNicId, cloudapi_v6.ArgFlowLogId)
}

func RunFlowLogList(c *core.CommandConfig) error {
	flowLogs, _, err := c.CloudApiV6Services.FlowLogs().List(
		viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgServerId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgNicId)),
	)
	if err != nil {
		return err
	}
	return c.Printer.Print(getFlowLogPrint(nil, c, getFlowLogs(flowLogs)))
}

func RunFlowLogGet(c *core.CommandConfig) error {
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgDataCenterId))
	serverId := viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgServerId))
	nicId := viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgNicId))
	flowLogId := viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgFlowLogId))
	c.Printer.Verbose("FlowLog with id: %v from Nic with id: %v is getting...", flowLogId, nicId)
	flowLog, _, err := c.CloudApiV6Services.FlowLogs().Get(
		dcId,
		serverId,
		nicId,
		flowLogId,
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
	flowLog, resp, err := c.CloudApiV6Services.FlowLogs().Create(
		viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgServerId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgNicId)),
		input,
	)
	if resp != nil {
		c.Printer.Verbose("Request href: %v ", resp.Header.Get("location"))
	}
	if err != nil {
		return err
	}

	if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetRequestPath(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getFlowLogPrint(resp, c, getFlowLog(flowLog)))
}

func RunFlowLogDelete(c *core.CommandConfig) error {
	if err := utils.AskForConfirm(c.Stdin, c.Printer, "delete flow log"); err != nil {
		return err
	}
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgDataCenterId))
	flowLogId := viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgFlowLogId))
	c.Printer.Verbose("FlowLog with id: %v is deleting...", flowLogId)
	resp, err := c.CloudApiV6Services.FlowLogs().Delete(dcId,
		viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgServerId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgNicId)),
		flowLogId,
	)
	if err != nil {
		return err
	}

	if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetRequestPath(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getFlowLogPrint(resp, c, nil))
}

// Get FlowLog Properties set used for create commands
func getFlowLogPropertiesSet(c *core.CommandConfig) resources.FlowLogProperties {
	properties := resources.FlowLogProperties{}
	name := viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgName))
	properties.SetName(name)
	c.Printer.Verbose("Property Name set: %v", name)
	action := viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgAction))
	properties.SetAction(strings.ToUpper(action))
	c.Printer.Verbose("Property Action set: %v", action)
	direction := viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgDirection))
	properties.SetDirection(strings.ToUpper(direction))
	c.Printer.Verbose("Property Direction set: %v", direction)
	if viper.IsSet(core.GetFlagName(c.NS, cloudapi_v6.ArgS3Bucket)) {
		bucketName := viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgS3Bucket))
		properties.SetBucket(bucketName)
		c.Printer.Verbose("Property Bucket set: %v", bucketName)
	}
	return properties
}

// Get FlowLog Properties set used for update commands
func getFlowLogPropertiesUpdate(c *core.CommandConfig) resources.FlowLogProperties {
	properties := resources.FlowLogProperties{}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapi_v6.ArgName)) {
		name := viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgName))
		properties.SetName(name)
		c.Printer.Verbose("Property Name set: %v", name)
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapi_v6.ArgAction)) {
		action := viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgAction))
		properties.SetAction(strings.ToUpper(action))
		c.Printer.Verbose("Property Action set: %v", action)
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapi_v6.ArgDirection)) {
		direction := viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgDirection))
		properties.SetDirection(strings.ToUpper(direction))
		c.Printer.Verbose("Property Direction set: %v", direction)
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapi_v6.ArgS3Bucket)) {
		bucketName := viper.GetString(core.GetFlagName(c.NS, cloudapi_v6.ArgS3Bucket))
		properties.SetBucket(bucketName)
		c.Printer.Verbose("Property Bucket set: %v", bucketName)
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
			if c.Resource != c.Namespace {
				r.Columns = getFlowLogsCols(core.GetFlagName(c.NS, config.ArgCols), c.Printer.GetStderr())
			} else {
				r.Columns = getFlowLogsCols(core.GetGlobalFlagName(c.Resource, config.ArgCols), c.Printer.GetStderr())
			}
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
