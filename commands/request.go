package commands

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/resources"
	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
	"github.com/ionos-cloud/ionosctl/pkg/utils/printer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func request() *core.Command {
	ctx := context.TODO()
	reqCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "request",
			Aliases:          []string{"req"},
			Short:            "Request Operations",
			Long:             `The sub-commands of ` + "`" + `ionosctl request` + "`" + ` allow you to see information about requests on your account. With the ` + "`" + `ionosctl request` + "`" + ` command, you can list, get or wait for a Request.`,
			TraverseChildren: true,
		},
	}
	globalFlags := reqCmd.GlobalFlags()
	globalFlags.StringSliceP(config.ArgCols, config.ArgColsShort, defaultRequestCols,
		fmt.Sprintf("Set of columns to be printed on output \nAvailable columns: %v", defaultRequestCols))
	_ = viper.BindPFlag(core.GetGlobalFlagName(reqCmd.Name(), config.ArgCols), globalFlags.Lookup(config.ArgCols))
	_ = reqCmd.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultRequestCols, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		List Command
	*/
	core.NewCommand(ctx, reqCmd, core.CommandBuilder{
		Namespace:  "request",
		Resource:   "request",
		Verb:       "list",
		ShortDesc:  "List Requests",
		LongDesc:   "Use this command to list all Requests on your account",
		Example:    "",
		PreCmdRun:  noPreRun,
		CmdRun:     RunRequestList,
		InitClient: true,
	})

	/*
		Get Command
	*/
	get := core.NewCommand(ctx, reqCmd, core.CommandBuilder{
		Namespace:  "request",
		Resource:   "request",
		Verb:       "get",
		ShortDesc:  "Get a Request",
		LongDesc:   "Use this command to get information about a specified Request.\n\nRequired values to run command:\n\n* Request Id",
		Example:    getRequestExample,
		PreCmdRun:  PreRunRequestId,
		CmdRun:     RunRequestGet,
		InitClient: true,
	})
	get.AddStringFlag(config.ArgRequestId, "", "", config.RequiredFlagRequestId)
	_ = get.Command.RegisterFlagCompletionFunc(config.ArgRequestId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getRequestsIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Wait Command
	*/
	wait := core.NewCommand(ctx, reqCmd, core.CommandBuilder{
		Namespace: "request",
		Resource:  "request",
		Verb:      "wait",
		ShortDesc: "Wait a Request",
		LongDesc: `Use this command to wait for a specified Request to execute.

You can specify a timeout for the Request to be executed using ` + "`" + `--timeout` + "`" + ` option.

Required values to run command:

* Request Id`,
		Example:    waitRequestExample,
		PreCmdRun:  PreRunRequestId,
		CmdRun:     RunRequestWait,
		InitClient: true,
	})
	wait.AddStringFlag(config.ArgRequestId, "", "", config.RequiredFlagRequestId)
	_ = wait.Command.RegisterFlagCompletionFunc(config.ArgRequestId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getRequestsIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	wait.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option waiting for Request [seconds]")

	return reqCmd
}

func PreRunRequestId(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.NS, config.ArgRequestId)
}

func RunRequestList(c *core.CommandConfig) error {
	requests, _, err := c.Requests().List()
	if err != nil {
		return err
	}
	rqs := getRequests(requests)
	return c.Printer.Print(printer.Result{
		OutputJSON: requests,
		Columns:    getRequestsCols(core.GetGlobalFlagName(c.Resource, config.ArgCols), c.Printer.GetStderr()),
		KeyValue:   getRequestsKVMaps(rqs),
	})
}

func RunRequestGet(c *core.CommandConfig) error {
	request, _, err := c.Requests().Get(viper.GetString(core.GetFlagName(c.NS, config.ArgRequestId)))
	if err != nil {
		return err
	}
	return c.Printer.Print(printer.Result{
		OutputJSON: request,
		KeyValue:   getRequestsKVMaps([]resources.Request{*request}),
		Columns:    getRequestsCols(core.GetGlobalFlagName(c.Resource, config.ArgCols), c.Printer.GetStderr()),
	})
}

func RunRequestWait(c *core.CommandConfig) error {
	request, _, err := c.Requests().Get(viper.GetString(core.GetFlagName(c.NS, config.ArgRequestId)))
	if err != nil {
		return err
	}

	// Default timeout: 60s
	timeout := viper.GetInt(core.GetFlagName(c.NS, config.ArgTimeout))
	ctxTimeout, cancel := context.WithTimeout(
		c.Context,
		time.Duration(timeout)*time.Second,
	)
	defer cancel()

	c.Context = ctxTimeout
	if _, err = c.Requests().Wait(fmt.Sprintf("%s/status", *request.GetHref())); err != nil {
		return err
	}
	return c.Printer.Print(printer.Result{
		OutputJSON: request,
		KeyValue:   getRequestsKVMaps([]resources.Request{*request}),
		Columns:    getRequestsCols(core.GetGlobalFlagName(c.Resource, config.ArgCols), c.Printer.GetStderr()),
	})
}

// Output Printing

var defaultRequestCols = []string{"RequestId", "Status", "Message"}

type RequestPrint struct {
	RequestId string `json:"RequestId,omitempty"`
	Status    string `json:"Status,omitempty"`
	Message   string `json:"Message,omitempty"`
}

func getRequestsCols(flagName string, outErr io.Writer) []string {
	var cols []string
	if viper.IsSet(flagName) {
		cols = viper.GetStringSlice(flagName)
	} else {
		return defaultRequestCols
	}

	columnsMap := map[string]string{
		"RequestId": "RequestId",
		"Status":    "Status",
		"Message":   "Message",
	}
	var requestCols []string
	for _, k := range cols {
		col := columnsMap[k]
		if col != "" {
			requestCols = append(requestCols, col)
		} else {
			clierror.CheckError(errors.New("unknown column "+k), outErr)
		}
	}
	return requestCols
}

func getRequests(requests resources.Requests) []resources.Request {
	req := make([]resources.Request, 0)
	for _, r := range *requests.Items {
		req = append(req, resources.Request{Request: r})
	}
	return req
}

func getRequestsKVMaps(requests []resources.Request) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(requests))
	for _, request := range requests {
		var reqPrint RequestPrint
		if id, ok := request.GetIdOk(); ok && id != nil {
			reqPrint.RequestId = *id
		}
		if status, ok := request.GetMetadata().GetRequestStatus().GetMetadata().GetStatusOk(); ok && status != nil {
			reqPrint.Status = *status
		}
		if msg, ok := request.GetMetadata().GetRequestStatus().GetMetadata().GetMessageOk(); ok && msg != nil {
			reqPrint.Message = *msg
		}
		o := structs.Map(reqPrint)
		out = append(out, o)
	}
	return out
}

func getRequestsIds(outErr io.Writer) []string {
	err := config.Load()
	clierror.CheckError(err, outErr)
	clientSvc, err := resources.NewClientService(
		viper.GetString(config.Username),
		viper.GetString(config.Password),
		viper.GetString(config.Token),
		viper.GetString(config.ArgServerUrl),
	)
	clierror.CheckError(err, outErr)
	reqSvc := resources.NewRequestService(clientSvc.Get(), context.TODO())
	requests, _, err := reqSvc.List()
	clierror.CheckError(err, outErr)
	reqIds := make([]string, 0)
	if items, ok := requests.Requests.GetItemsOk(); ok && items != nil {
		for _, item := range *items {
			if itemId, ok := item.GetIdOk(); ok && itemId != nil {
				reqIds = append(reqIds, *itemId)
			}
		}
	} else {
		return nil
	}
	return reqIds
}
