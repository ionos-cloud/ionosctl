package commands

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/resources"
	"github.com/ionos-cloud/ionosctl/pkg/utils"
	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
	"github.com/ionos-cloud/ionosctl/pkg/utils/printer"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
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
			Long:             "The sub-commands of `ionosctl request` allow you to see information about requests on your account. With the `ionosctl request` command, you can list, get or wait for a Request.",
			TraverseChildren: true,
		},
	}
	globalFlags := reqCmd.GlobalFlags()
	globalFlags.StringSliceP(config.ArgCols, "", defaultRequestCols, utils.ColsMessage(allRequestCols))
	_ = viper.BindPFlag(core.GetGlobalFlagName(reqCmd.Name(), config.ArgCols), globalFlags.Lookup(config.ArgCols))
	_ = reqCmd.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allRequestCols, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		List Command
	*/
	list := core.NewCommand(ctx, reqCmd, core.CommandBuilder{
		Namespace: "request",
		Resource:  "request",
		Verb:      "list",
		Aliases:   []string{"l", "ls"},
		ShortDesc: "List Requests",
		LongDesc: `Use this command to list all Requests on your account.

Use flags to retrieve a list of Requests:

* sorting by the time the Request was created, starting from now in descending order, take the first N Requests: ` + "`" + `ionosctl request list --latest N` + "`" + `
* sorting by method: ` + "`" + `ionosctl request list --method REQUEST_METHOD` + "`" + `, where REQUEST_METHOD can be CREATE or POST, UPDATE or PATCH, PUT and DELETE
* sorting by both of the above options: ` + "`" + `ionosctl request list --method REQUEST_METHOD --latest N` + "`" + ``,
		Example:    listRequestsExample,
		PreCmdRun:  noPreRun,
		CmdRun:     RunRequestList,
		InitClient: true,
	})
	list.AddIntFlag(config.ArgLatest, "", 0, "Show latest N Requests. If it is not set, all Requests will be printed")
	list.AddStringFlag(config.ArgMethod, "", "", "Show only the Requests with this method. E.g CREATE, UPDATE, DELETE")
	_ = list.Command.RegisterFlagCompletionFunc(config.ArgMethod, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"POST", "PUT", "DELETE", "PATCH", "CREATE", "UPDATE"}, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Get Command
	*/
	get := core.NewCommand(ctx, reqCmd, core.CommandBuilder{
		Namespace:  "request",
		Resource:   "request",
		Verb:       "get",
		Aliases:    []string{"g"},
		ShortDesc:  "Get a Request",
		LongDesc:   "Use this command to get information about a specified Request.\n\nRequired values to run command:\n\n* Request Id",
		Example:    getRequestExample,
		PreCmdRun:  PreRunRequestId,
		CmdRun:     RunRequestGet,
		InitClient: true,
	})
	get.AddStringFlag(config.ArgRequestId, config.ArgIdShort, "", config.RequiredFlagRequestId)
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
		Aliases:   []string{"w"},
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
	wait.AddStringFlag(config.ArgRequestId, config.ArgIdShort, "", config.RequiredFlagRequestId)
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
	if viper.IsSet(core.GetFlagName(c.NS, config.ArgMethod)) {
		switch strings.ToUpper(viper.GetString(core.GetFlagName(c.NS, config.ArgMethod))) {
		case "CREATE":
			requests = sortRequestsByMethod(requests, "POST")
		case "UPDATE":
			// On UPDATE, take Requests with PATCH and PUT methods
			sortReqsUpdated := make([]ionoscloud.Request, 0)
			requestsPatch := sortRequestsByMethod(requests, "PATCH")
			requestsPut := sortRequestsByMethod(requests, "PUT")
			if len(getRequests(requestsPatch)) > 0 {
				for _, requestPatch := range getRequests(requestsPatch) {
					sortReqsUpdated = append(sortReqsUpdated, requestPatch.Request)
				}
			}
			if len(getRequests(requestsPut)) > 0 {
				for _, requestPut := range getRequests(requestsPut) {
					sortReqsUpdated = append(sortReqsUpdated, requestPut.Request)
				}
			}
			requests.Items = &sortReqsUpdated
		default:
			requests = sortRequestsByMethod(requests, strings.ToUpper(viper.GetString(core.GetFlagName(c.NS, config.ArgMethod))))
		}
	}
	if viper.IsSet(core.GetFlagName(c.NS, config.ArgLatest)) {
		requests = sortRequestsByTime(requests, viper.GetInt(core.GetFlagName(c.NS, config.ArgLatest)))
	}

	if itemsOk, ok := requests.GetItemsOk(); ok && itemsOk != nil {
		if len(*itemsOk) == 0 {
			return errors.New("error getting requests based on given criteria")
		}
	}
	return c.Printer.Print(getRequestPrint(c, getRequests(requests)))
}

func RunRequestGet(c *core.CommandConfig) error {
	req, _, err := c.Requests().Get(viper.GetString(core.GetFlagName(c.NS, config.ArgRequestId)))
	if err != nil {
		return err
	}
	return c.Printer.Print(getRequestPrint(c, []resources.Request{*req}))
}

func RunRequestWait(c *core.CommandConfig) error {
	req, _, err := c.Requests().Get(viper.GetString(core.GetFlagName(c.NS, config.ArgRequestId)))
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
	if _, err = c.Requests().Wait(fmt.Sprintf("%s/status", *req.GetHref())); err != nil {
		return err
	}
	return c.Printer.Print(getRequestPrint(c, []resources.Request{*req}))
}

// Output Printing

var (
	defaultRequestCols = []string{"RequestId", "CreatedDate", "Method", "Status", "Message"}
	allRequestCols     = []string{"RequestId", "CreatedDate", "CreatedBy", "Method", "Status", "Message", "Url", "Body"}
)

type RequestPrint struct {
	RequestId   string    `json:"RequestId,omitempty"`
	Status      string    `json:"Status,omitempty"`
	Message     string    `json:"Message,omitempty"`
	Method      string    `json:"Method,omitempty"`
	Url         string    `json:"Url,omitempty"`
	Body        string    `json:"Body,omitempty"`
	CreatedBy   string    `json:"CreatedBy,omitempty"`
	CreatedDate time.Time `json:"CreatedDate,omitempty"`
}

func getRequestPrint(c *core.CommandConfig, reqs []resources.Request) printer.Result {
	r := printer.Result{}
	if c != nil {
		if reqs != nil {
			r.OutputJSON = reqs
			r.KeyValue = getRequestsKVMaps(reqs)
			r.Columns = getRequestsCols(core.GetGlobalFlagName(c.Resource, config.ArgCols), c.Printer.GetStderr())
		}
	}
	return r
}

func getRequestsCols(flagName string, outErr io.Writer) []string {
	var cols []string
	if viper.IsSet(flagName) {
		cols = viper.GetStringSlice(flagName)
	} else {
		return defaultRequestCols
	}

	columnsMap := map[string]string{
		"RequestId":   "RequestId",
		"Status":      "Status",
		"Message":     "Message",
		"Method":      "Method",
		"Url":         "Url",
		"Body":        "Body",
		"CreatedDate": "CreatedDate",
		"CreatedBy":   "CreatedBy",
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

func sortRequestsByMethod(requests resources.Requests, method string) resources.Requests {
	var sortedRequests resources.Requests
	if items, ok := requests.GetItemsOk(); ok && items != nil {
		requestsItems := make([]ionoscloud.Request, 0)
		for _, item := range *items {
			properties := item.GetProperties()
			if methodOk, ok := properties.GetMethodOk(); ok && methodOk != nil {
				if *methodOk == method {
					requestsItems = append(requestsItems, item)
				}
			}
		}
		sortedRequests.Items = &requestsItems
	}
	return sortedRequests
}

func sortRequestsByTime(requests resources.Requests, n int) resources.Requests {
	var sortedRequests resources.Requests
	if items, ok := requests.GetItemsOk(); ok && items != nil {
		reqItems := *items
		if len(reqItems) > 0 {
			// Sort requests using time.Time, starting from now in descending order
			sort.SliceStable(reqItems, func(i, j int) bool {
				return reqItems[i].Metadata.CreatedDate.Time.After(reqItems[j].Metadata.CreatedDate.Time)
			})
		}
		if len(reqItems) >= n {
			// Take the first N requests
			reqItems = reqItems[:n]
		}
		sortedRequests.Items = &reqItems
	}
	return sortedRequests
}

func getRequestsKVMaps(requests []resources.Request) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(requests))
	for _, req := range requests {
		var reqPrint RequestPrint
		if id, ok := req.GetIdOk(); ok && id != nil {
			reqPrint.RequestId = *id
		}
		if propertiesOk, ok := req.GetPropertiesOk(); ok && propertiesOk != nil {
			if method, ok := propertiesOk.GetMethodOk(); ok && method != nil {
				reqPrint.Method = *method
			}
			if url, ok := propertiesOk.GetUrlOk(); ok && url != nil {
				reqPrint.Url = *url
			}
			if bodyOk, ok := propertiesOk.GetBodyOk(); ok && bodyOk != nil {
				reqPrint.Body = *bodyOk
			}
		}
		if metadataOk, ok := req.GetMetadataOk(); ok && metadataOk != nil {
			if createdDateOk, ok := metadataOk.GetCreatedDateOk(); ok && createdDateOk != nil {
				reqPrint.CreatedDate = *createdDateOk
			}
			if createdByOk, ok := metadataOk.GetCreatedByOk(); ok && createdByOk != nil {
				reqPrint.CreatedBy = *createdByOk
			}
			if requestStatusOk, ok := metadataOk.GetRequestStatusOk(); ok && requestStatusOk != nil {
				if statusMetadata, ok := requestStatusOk.GetMetadataOk(); ok && statusMetadata != nil {
					if statusOk, ok := statusMetadata.GetStatusOk(); ok && statusOk != nil {
						reqPrint.Status = *statusOk
					}
					if messageOk, ok := statusMetadata.GetMessageOk(); ok && messageOk != nil {
						reqPrint.Message = *messageOk
					}
				}
			}
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
