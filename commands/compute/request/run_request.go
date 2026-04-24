package request

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/viper"
)

func PreRunRequestId(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgRequestId)
}

func RunRequestList(c *core.CommandConfig) error {
	requests, resp, err := c.CloudApiV6Services.Requests().List()
	if resp != nil {
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
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

	return c.Printer(allRequestCols).Prefix("items").Print(requests.Requests)
}

func RunRequestGet(c *core.CommandConfig) error {
	c.Verbose("Request with id: %v is getting...", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgRequestId)))

	req, resp, err := c.CloudApiV6Services.Requests().Get(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgRequestId)))
	if resp != nil {
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allRequestCols).Print(req.Request)
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

	return c.Printer(allRequestCols).Print(req.Request)
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
