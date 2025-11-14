package commands

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/completer"
	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/query"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/resource2table"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	defaultRequestCols = []string{"RequestId", "CreatedDate", "Method", "Status", "Message", "Targets"}
	allRequestCols     = []string{"RequestId", "CreatedDate", "CreatedBy", "Method", "Status", "Message", "Url", "Body", "Targets"}
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
	globalFlags.StringSliceP(constants.ArgCols, "", defaultRequestCols, tabheaders.ColsMessage(allRequestCols))
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
	list.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, int32(2), cloudapiv6.ArgDepthDescription)
	list.AddStringFlag(cloudapiv6.ArgOrderBy, "", "", cloudapiv6.ArgOrderByDescription)
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgOrderBy, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.RequestsFilters(), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddStringSliceFlag(cloudapiv6.ArgFilters, cloudapiv6.ArgFiltersShort, []string{""}, cloudapiv6.ArgFiltersDescription)
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgFilters, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.RequestsFilters(), cobra.ShellCompDirectiveNoFileComp
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
	get.AddUUIDFlag(cloudapiv6.ArgRequestId, cloudapiv6.ArgIdShort, "", cloudapiv6.RequestId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgRequestId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.RequestsIds(), cobra.ShellCompDirectiveNoFileComp
	})
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
		return completer.RequestsIds(), cobra.ShellCompDirectiveNoFileComp
	})
	wait.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, constants.DefaultTimeoutSeconds, "Timeout option waiting for Request [seconds]")
	wait.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultMiscDepth, cloudapiv6.ArgDepthDescription)

	return core.WithConfigOverride(reqCmd, []string{fileconfiguration.Cloud, "compute"}, "")
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

	requests, resp, err := c.CloudApiV6Services.Requests().List(listQueryParams)
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
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

	if itemsOk, ok := requests.GetItemsOk(); !ok || itemsOk == nil {
		return nil
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	convertedRequests, err := resource2table.ConvertRequestsToTable(requests.Requests)
	if err != nil {
		return err
	}

	out, err := jsontabwriter.GenerateOutputPreconverted(requests.Requests, convertedRequests,
		tabheaders.GetHeaders(allRequestCols, defaultRequestCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)

	return nil
}

func RunRequestGet(c *core.CommandConfig) error {
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
		"Request with id: %v is getting...", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgRequestId))))

	req, resp, err := c.CloudApiV6Services.Requests().Get(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgRequestId)))
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	convertedReq, err := resource2table.ConvertRequestToTable(req.Request)
	if err != nil {
		return err
	}

	out, err := jsontabwriter.GenerateOutputPreconverted(req.Request, convertedReq,
		tabheaders.GetHeaders(allRequestCols, defaultRequestCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)

	return nil
}

func RunRequestWait(c *core.CommandConfig) error {
	req, _, err := c.CloudApiV6Services.Requests().Get(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgRequestId)))
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

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	convertedReq, err := resource2table.ConvertRequestToTable(req.Request)
	if err != nil {
		return err
	}

	out, err := jsontabwriter.GenerateOutputPreconverted(req.Request, convertedReq,
		tabheaders.GetHeaders(allRequestCols, defaultRequestCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)

	return nil
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
