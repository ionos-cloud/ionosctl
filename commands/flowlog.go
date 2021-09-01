package commands

import (
	"context"
	"errors"
	"io"
	"os"
	"strings"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/resources/v6"
	"github.com/ionos-cloud/ionosctl/pkg/utils"
	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
	"github.com/ionos-cloud/ionosctl/pkg/utils/printer"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	globalFlags.StringSliceP(config.ArgCols, "", defaultFlowLogCols, utils.ColsMessage(defaultFlowLogCols))
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
	list.AddStringFlag(config.ArgDataCenterId, "", "", config.DatacenterId, core.RequiredFlagOption())
	_ = list.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddStringFlag(config.ArgServerId, "", "", config.ServerId, core.RequiredFlagOption())
	_ = list.Command.RegisterFlagCompletionFunc(config.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getServersIds(os.Stderr, viper.GetString(core.GetFlagName(list.NS, config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddStringFlag(config.ArgNicId, "", "", config.NicId, core.RequiredFlagOption())
	_ = list.Command.RegisterFlagCompletionFunc(config.ArgNicId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getNicsIds(os.Stderr, viper.GetString(core.GetFlagName(list.NS, config.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(list.NS, config.ArgServerId))), cobra.ShellCompDirectiveNoFileComp
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
	get.AddStringFlag(config.ArgDataCenterId, "", "", config.DatacenterId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddStringFlag(config.ArgServerId, "", "", config.ServerId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(config.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getServersIds(os.Stderr, viper.GetString(core.GetFlagName(get.NS, config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddStringFlag(config.ArgNicId, "", "", config.NicId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(config.ArgNicId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getNicsIds(os.Stderr, viper.GetString(core.GetFlagName(get.NS, config.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(get.NS, config.ArgServerId))), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddStringFlag(config.ArgFlowLogId, config.ArgIdShort, "", config.FlowLogId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(config.ArgFlowLogId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getFlowLogsIds(os.Stderr,
			viper.GetString(core.GetFlagName(get.NS, config.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(get.NS, config.ArgServerId)),
			viper.GetString(core.GetFlagName(get.NS, config.ArgNicId))), cobra.ShellCompDirectiveNoFileComp
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
	create.AddStringFlag(config.ArgDataCenterId, "", "", config.DatacenterId, core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(config.ArgServerId, "", "", config.ServerId, core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(config.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getServersIds(os.Stderr, viper.GetString(core.GetFlagName(create.NS, config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(config.ArgNicId, "", "", config.NicId, core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(config.ArgNicId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getNicsIds(os.Stderr, viper.GetString(core.GetFlagName(create.NS, config.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(create.NS, config.ArgServerId))), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(config.ArgName, config.ArgNameShort, "Unnamed FlowLog", "The name for the FlowLog")
	create.AddStringFlag(config.ArgAction, config.ArgActionShort, "ALL", "Specifies the traffic Action pattern")
	_ = create.Command.RegisterFlagCompletionFunc(config.ArgAction, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"ALL", "REJECTED", "ACCEPTED"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(config.ArgDirection, config.ArgDirectionShort, "BIDIRECTIONAL", "Specifies the traffic Direction pattern")
	_ = create.Command.RegisterFlagCompletionFunc(config.ArgDirection, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"BIDIRECTIONAL", "INGRESS", "EGRESS"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(config.ArgS3Bucket, config.ArgS3BucketShort, "", "S3 Bucket name of an existing IONOS Cloud S3 Bucket", core.RequiredFlagOption())
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
	deleteCmd.AddStringFlag(config.ArgDataCenterId, "", "", config.DatacenterId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(config.ArgDataCenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getDataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddStringFlag(config.ArgServerId, "", "", config.ServerId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(config.ArgServerId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getServersIds(os.Stderr, viper.GetString(core.GetFlagName(deleteCmd.NS, config.ArgDataCenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddStringFlag(config.ArgNicId, "", "", config.NicId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(config.ArgNicId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getNicsIds(os.Stderr, viper.GetString(core.GetFlagName(deleteCmd.NS, config.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(deleteCmd.NS, config.ArgServerId))), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddStringFlag(config.ArgFlowLogId, config.ArgIdShort, "", config.FlowLogId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(config.ArgFlowLogId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getFlowLogsIds(os.Stderr, viper.GetString(core.GetFlagName(deleteCmd.NS, config.ArgDataCenterId)),
			viper.GetString(core.GetFlagName(deleteCmd.NS, config.ArgServerId)),
			viper.GetString(core.GetFlagName(deleteCmd.NS, config.ArgNicId))), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for Request for FlowLog deletion to be executed")
	deleteCmd.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option for Request for FlowLog deletion [seconds]")

	return flowLogCmd
}

func PreRunFlowLogCreate(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, config.ArgDataCenterId, config.ArgServerId, config.ArgNicId, config.ArgS3Bucket)
}

func PreRunDcServerNicFlowLogIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, config.ArgDataCenterId, config.ArgServerId, config.ArgNicId, config.ArgFlowLogId)
}

func RunFlowLogList(c *core.CommandConfig) error {
	flowLogs, _, err := c.FlowLogs().List(
		viper.GetString(core.GetFlagName(c.NS, config.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgServerId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgNicId)),
	)
	if err != nil {
		return err
	}
	return c.Printer.Print(getFlowLogPrint(nil, c, getFlowLogs(flowLogs)))
}

func RunFlowLogGet(c *core.CommandConfig) error {
	dcId := viper.GetString(core.GetFlagName(c.NS, config.ArgDataCenterId))
	serverId := viper.GetString(core.GetFlagName(c.NS, config.ArgServerId))
	nicId := viper.GetString(core.GetFlagName(c.NS, config.ArgNicId))
	flowLogId := viper.GetString(core.GetFlagName(c.NS, config.ArgFlowLogId))
	c.Printer.Verbose("FlowLog with id: %v from Nic with id: %v is getting...", flowLogId, nicId)
	flowLog, _, err := c.FlowLogs().Get(
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
	input := v6.FlowLog{
		FlowLog: ionoscloud.FlowLog{
			Properties: &properties.FlowLogProperties,
		},
	}
	flowLog, resp, err := c.FlowLogs().Create(
		viper.GetString(core.GetFlagName(c.NS, config.ArgDataCenterId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgServerId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgNicId)),
		input,
	)
	if resp != nil {
		c.Printer.Verbose("Request href: %v ", resp.Header.Get("location"))
	}
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
	dcId := viper.GetString(core.GetFlagName(c.NS, config.ArgDataCenterId))
	flowLogId := viper.GetString(core.GetFlagName(c.NS, config.ArgFlowLogId))
	c.Printer.Verbose("FlowLog with id: %v is deleting...", flowLogId)
	resp, err := c.FlowLogs().Delete(dcId,
		viper.GetString(core.GetFlagName(c.NS, config.ArgServerId)),
		viper.GetString(core.GetFlagName(c.NS, config.ArgNicId)),
		flowLogId,
	)
	if err != nil {
		return err
	}

	if err = utils.WaitForRequest(c, printer.GetRequestPath(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getFlowLogPrint(resp, c, nil))
}

// Get FlowLog Properties set used for create commands
func getFlowLogPropertiesSet(c *core.CommandConfig) v6.FlowLogProperties {
	properties := v6.FlowLogProperties{}
	name := viper.GetString(core.GetFlagName(c.NS, config.ArgName))
	properties.SetName(name)
	c.Printer.Verbose("Property Name set: %v", name)
	action := viper.GetString(core.GetFlagName(c.NS, config.ArgAction))
	properties.SetAction(strings.ToUpper(action))
	c.Printer.Verbose("Property Action set: %v", action)
	direction := viper.GetString(core.GetFlagName(c.NS, config.ArgDirection))
	properties.SetDirection(strings.ToUpper(direction))
	c.Printer.Verbose("Property Direction set: %v", direction)
	if viper.IsSet(core.GetFlagName(c.NS, config.ArgS3Bucket)) {
		bucketName := viper.GetString(core.GetFlagName(c.NS, config.ArgS3Bucket))
		properties.SetBucket(bucketName)
		c.Printer.Verbose("Property Bucket set: %v", bucketName)
	}
	return properties
}

// Get FlowLog Properties set used for update commands
func getFlowLogPropertiesUpdate(c *core.CommandConfig) v6.FlowLogProperties {
	properties := v6.FlowLogProperties{}
	if viper.IsSet(core.GetFlagName(c.NS, config.ArgName)) {
		name := viper.GetString(core.GetFlagName(c.NS, config.ArgName))
		properties.SetName(name)
		c.Printer.Verbose("Property Name set: %v", name)
	}
	if viper.IsSet(core.GetFlagName(c.NS, config.ArgAction)) {
		action := viper.GetString(core.GetFlagName(c.NS, config.ArgAction))
		properties.SetAction(strings.ToUpper(action))
		c.Printer.Verbose("Property Action set: %v", action)
	}
	if viper.IsSet(core.GetFlagName(c.NS, config.ArgDirection)) {
		direction := viper.GetString(core.GetFlagName(c.NS, config.ArgDirection))
		properties.SetDirection(strings.ToUpper(direction))
		c.Printer.Verbose("Property Direction set: %v", direction)
	}
	if viper.IsSet(core.GetFlagName(c.NS, config.ArgS3Bucket)) {
		bucketName := viper.GetString(core.GetFlagName(c.NS, config.ArgS3Bucket))
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

func getFlowLogPrint(resp *v6.Response, c *core.CommandConfig, rule []v6.FlowLog) printer.Result {
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

func getFlowLogs(flowLogs v6.FlowLogs) []v6.FlowLog {
	ls := make([]v6.FlowLog, 0)
	if items, ok := flowLogs.GetItemsOk(); ok && items != nil {
		for _, s := range *items {
			ls = append(ls, v6.FlowLog{FlowLog: s})
		}
	}
	return ls
}

func getFlowLog(flowLog *v6.FlowLog) []v6.FlowLog {
	ss := make([]v6.FlowLog, 0)
	if flowLog != nil {
		ss = append(ss, v6.FlowLog{FlowLog: flowLog.FlowLog})
	}
	return ss
}

func getFlowLogsKVMaps(ls []v6.FlowLog) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(ls))
	if len(ls) > 0 {
		for _, l := range ls {
			o := getFlowLogKVMap(l)
			out = append(out, o)
		}
	}
	return out
}

func getFlowLogKVMap(l v6.FlowLog) map[string]interface{} {
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
	clientSvc, err := v6.NewClientService(
		viper.GetString(config.Username),
		viper.GetString(config.Password),
		viper.GetString(config.Token),
		config.GetServerUrl(),
	)
	clierror.CheckError(err, outErr)
	flowLogSvc := v6.NewFlowLogService(clientSvc.Get(), context.TODO())
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
