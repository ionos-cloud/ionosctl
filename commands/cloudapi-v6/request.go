package commands

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/completer"
	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/query"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/printer"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func RequestCmd() *core.Command {
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
	globalFlags.StringSliceP(constants.ArgCols, "", defaultRequestCols, printer.ColsMessage(allRequestCols))
	_ = viper.BindPFlag(core.GetFlagName(reqCmd.Name(), constants.ArgCols), globalFlags.Lookup(constants.ArgCols))
	_ = reqCmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allRequestCols, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		List Command
	*/
	list := core.NewCommand(ctx, reqCmd, core.CommandBuilder{
		Namespace:  "request",
		Resource:   "request",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List Requests",
		LongDesc:   "Use this command to list all Requests on your account.\n\nYou can filter the results using `--filters` option. Use the following format to set filters: `--filters KEY1=VALUE1,KEY2=VALUE2`.\n" + completer.RequestsFiltersUsage(),
		Example:    listRequestExample,
		PreCmdRun:  PreRunRequestList,
		CmdRun:     RunRequestList,
		InitClient: true,
	})
	list.AddIntFlag(cloudapiv6.ArgLatest, "", 0, "Show latest N Requests. If it is not set, all Requests will be printed", core.DeprecatedFlagOption("Use --filters --order-by --max-results options instead!"))
	list.AddStringFlag(cloudapiv6.ArgMethod, "", "", "Show only the Requests with this method. E.g CREATE, UPDATE, DELETE", core.DeprecatedFlagOption("Use --filters --order-by --max-results options instead!"))
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgMethod, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"POST", "PUT", "DELETE", "PATCH", "CREATE", "UPDATE"}, cobra.ShellCompDirectiveNoFileComp
	})
	list.AddInt32Flag(constants.FlagMaxResults, constants.FlagMaxResultsShort, cloudapiv6.DefaultMaxResults, constants.DescMaxResults)
	list.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultListDepth, cloudapiv6.ArgDepthDescription)
	list.AddStringFlag(cloudapiv6.ArgOrderBy, "", "", cloudapiv6.ArgOrderByDescription)
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgOrderBy, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.RequestsFilters(), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddStringSliceFlag(cloudapiv6.ArgFilters, cloudapiv6.ArgFiltersShort, []string{""}, cloudapiv6.ArgFiltersDescription)
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgFilters, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.RequestsFilters(), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddBoolFlag(constants.ArgNoHeaders, "", false, cloudapiv6.ArgNoHeadersDescription)

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
	get.AddUUIDFlag(cloudapiv6.ArgRequestId, cloudapiv6.ArgIdShort, "", cloudapiv6.RequestId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgRequestId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.RequestsIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddBoolFlag(constants.ArgNoHeaders, "", false, cloudapiv6.ArgNoHeadersDescription)
	get.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultGetDepth, cloudapiv6.ArgDepthDescription)

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
	wait.AddUUIDFlag(cloudapiv6.ArgRequestId, cloudapiv6.ArgIdShort, "", cloudapiv6.RequestId, core.RequiredFlagOption())
	_ = wait.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgRequestId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.RequestsIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	wait.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, constants.DefaultTimeoutSeconds, "Timeout option waiting for Request [seconds]")
	wait.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultMiscDepth, cloudapiv6.ArgDepthDescription)

	return reqCmd
}

func PreRunRequestList(c *core.PreCommandConfig) error {
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgFilters)) {
		return query.ValidateFilters(c, completer.RequestsFilters(), completer.RequestsFiltersUsage())
	}
	return nil
}

func PreRunRequestId(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgRequestId)
}

func RunRequestList(c *core.CommandConfig) error {
	// Add Query Parameters for GET Requests
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	requests, resp, err := c.CloudApiV6Services.Requests().List(listQueryParams)
	if resp != nil {
		c.Printer.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgMethod)) {
		switch strings.ToUpper(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgMethod))) {
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
			requests = sortRequestsByMethod(requests, strings.ToUpper(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgMethod))))
		}
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgLatest)) {
		requests = sortRequestsByTime(requests, viper.GetInt(core.GetFlagName(c.NS, cloudapiv6.ArgLatest)))
	}

	if itemsOk, ok := requests.GetItemsOk(); ok && itemsOk != nil {
		if len(*itemsOk) == 0 {
			return errors.New("error getting requests based on given criteria")
		}
	}
	return c.Printer.Print(getRequestPrint(c, getRequests(requests)))
}

func RunRequestGet(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	queryParams := listQueryParams.QueryParams
	c.Printer.Verbose("Request with id: %v is getting...", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgRequestId)))
	req, resp, err := c.CloudApiV6Services.Requests().Get(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgRequestId)), queryParams)
	if resp != nil {
		c.Printer.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getRequestPrint(c, []resources.Request{*req}))
}

func RunRequestWait(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	queryParams := listQueryParams.QueryParams
	req, _, err := c.CloudApiV6Services.Requests().Get(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgRequestId)), queryParams)
	if err != nil {
		return err
	}

	// Default timeout: 60s
	timeout := viper.GetInt(core.GetFlagName(c.NS, constants.ArgTimeout))
	ctxTimeout, cancel := context.WithTimeout(
		c.Context,
		time.Duration(timeout)*time.Second,
	)
	defer cancel()

	c.Context = ctxTimeout
	if _, err = c.CloudApiV6Services.Requests().Wait(fmt.Sprintf("%s/status", *req.GetHref())); err != nil {
		return err
	}
	return c.Printer.Print(getRequestPrint(c, []resources.Request{*req}))
}

// Output Printing

var (
	defaultRequestCols = []string{"RequestId", "CreatedDate", "Method", "Status", "Message", "Targets"}
	allRequestCols     = []string{"RequestId", "CreatedDate", "CreatedBy", "Method", "Status", "Message", "Url", "Body", "Targets"}
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
	Targets     []string  `json:"Targets,omitempty"`
}

func getRequestPrint(c *core.CommandConfig, reqs []resources.Request) printer.Result {
	r := printer.Result{}
	if c != nil {
		if reqs != nil {
			r.OutputJSON = reqs
			r.KeyValue = getRequestsKVMaps(reqs)
			r.Columns = printer.GetHeaders(allRequestCols, defaultRequestCols, viper.GetStringSlice(core.GetFlagName(c.NS, cloudapiv6.ArgCols)))
		}
	}
	return r
}

func getRequests(requests resources.Requests) []resources.Request {
	requestObjs := make([]resources.Request, 0)
	if items, ok := requests.GetItemsOk(); ok && items != nil {
		for _, request := range *items {
			requestObjs = append(requestObjs, resources.Request{Request: request})
		}
	}
	return requestObjs
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
					if targetsOk, ok := statusMetadata.GetTargetsOk(); ok && targetsOk != nil {
						for _, target := range *targetsOk {
							if targetOk, ok := target.GetTargetOk(); ok && targetOk != nil {
								if typeOk, ok := targetOk.GetTypeOk(); ok && typeOk != nil {
									reqPrint.Targets = append(reqPrint.Targets, string(*typeOk))
									reqPrint.Targets = append(reqPrint.Targets, " ")
								}
								if idOk, ok := targetOk.GetIdOk(); ok && idOk != nil {
									reqPrint.Targets = append(reqPrint.Targets, *idOk)
									reqPrint.Targets = append(reqPrint.Targets, " ")
								}
							}
						}
					}
				}
			}
		}
		o := structs.Map(reqPrint)
		out = append(out, o)
	}
	return out
}
