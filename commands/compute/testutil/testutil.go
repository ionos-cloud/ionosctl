// Package testutil provides shared test fixtures for compute command tests.
package testutil

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
)

var (
	// TestRequestIdVar is a test request ID used across test files.
	TestRequestIdVar = "f2354da4-83e3-4e92-9d23-f3cb1ffecc31"

	// TestResponse is a standard mock API response used by nearly all test files.
	TestResponse = resources.Response{
		APIResponse: ionoscloud.APIResponse{
			Response: &http.Response{
				Header: map[string][]string{
					"Location": {"https://api.ionos.com/cloudapi/v6/requests/f2354da4-83e3-4e92-9d23-f3cb1ffecc31/status"},
				},
			},
			RequestTime: time.Duration(50),
		},
	}

	// TestResponseErr is a mock API response with an empty location header.
	TestResponseErr = resources.Response{
		APIResponse: ionoscloud.APIResponse{
			Response: &http.Response{
				Header: map[string][]string{
					"Location": {""},
				},
			},
			RequestTime: time.Duration(50),
		},
	}

	// TestStateVar is a common test state value.
	TestStateVar = "ACTIVE"

	// TestQueryParamVar is a common test filter value.
	TestQueryParamVar = "test-filter"

	// TestOrderByVar is a default test orderBy value.
	TestOrderByVar = ""

	// TestDepthListVar is a common test depth for list operations.
	TestDepthListVar = int32(1)

	// TestDepthOtherVar is a common test depth for non-list operations.
	TestDepthOtherVar = int32(0)

	// TestListQueryParamFilters contains test list query params with filters.
	TestListQueryParamFilters = resources.ListQueryParams{
		Filters: &map[string][]string{
			TestQueryParamVar: {TestQueryParamVar},
		},
		OrderBy: &TestQueryParamVar,
		QueryParams: resources.QueryParams{
			Depth: &TestDepthListVar,
		},
	}

	// TestListQueryParam contains default test list query params.
	TestListQueryParam = resources.ListQueryParams{
		OrderBy: &TestOrderByVar,
		QueryParams: resources.QueryParams{
			Depth: &TestDepthListVar,
		},
	}

	// TestQueryParamOther contains test query params for non-list operations.
	TestQueryParamOther = resources.QueryParams{
		Depth: &TestDepthOtherVar,
	}

	// Request test fixtures (from request_test.go, used across many packages)

	TestRequestVar       = "test-request"
	TestRequestStatusVar = "DONE"
	TestTypeRequestVar   = ionoscloud.Type("datacenter")
	TestRequestErr       = errors.New("request test: error occurred")

	TestRequestTargetsVar = []ionoscloud.RequestTarget{
		{
			Target: &ionoscloud.ResourceReference{
				Id:   &TestRequestVar,
				Type: &TestTypeRequestVar,
			},
		},
	}

	TestRequestStatus = resources.RequestStatus{
		RequestStatus: ionoscloud.RequestStatus{
			Id: &TestRequestVar,
			Metadata: &ionoscloud.RequestStatusMetadata{
				Status:  &TestRequestStatusVar,
				Message: &TestRequestVar,
				Targets: &TestRequestTargetsVar,
			},
		},
	}

	// Flowlog test fixtures (from flowlog_test.go, used by natgateway/nlb/alb flowlog tests)

	TestFlowLogState       = "AVAILABLE"
	TestFlowLogVar         = "test-flowlog"
	TestFlowLogUpperVar    = strings.ToUpper(TestFlowLogVar)
	TestFlowLogNewVar      = "test-new-flowlog"
	TestFlowLogNewUpperVar = strings.ToUpper(TestFlowLogNewVar)
	TestFlowLogErr         = errors.New("flowlog test error")

	TestFlowLog = resources.FlowLog{
		FlowLog: ionoscloud.FlowLog{
			Id: &TestFlowLogVar,
			Properties: &ionoscloud.FlowLogProperties{
				Name:      &TestFlowLogVar,
				Action:    &TestFlowLogUpperVar,
				Direction: &TestFlowLogUpperVar,
				Bucket:    &TestFlowLogVar,
			},
			Metadata: &ionoscloud.DatacenterElementMetadata{
				State: &TestFlowLogState,
			},
		},
	}

	TestFlowLogsList = resources.FlowLogs{
		FlowLogs: ionoscloud.FlowLogs{
			Id: &TestFlowLogVar,
			Items: &[]ionoscloud.FlowLog{
				TestFlowLog.FlowLog,
				TestFlowLog.FlowLog,
			},
		},
	}

	TestInputFlowLog = resources.FlowLog{
		FlowLog: ionoscloud.FlowLog{
			Properties: TestFlowLog.FlowLog.Properties,
		},
	}

	TestFlowLogProperties = resources.FlowLogProperties{
		FlowLogProperties: ionoscloud.FlowLogProperties{
			Name:      &TestFlowLogNewVar,
			Action:    &TestFlowLogNewUpperVar,
			Direction: &TestFlowLogNewUpperVar,
			Bucket:    &TestFlowLogNewVar,
		},
	}

	TestFlowLogUpdated = resources.FlowLog{
		FlowLog: ionoscloud.FlowLog{
			Properties: &TestFlowLogProperties.FlowLogProperties,
		},
	}

	TestFlowLogs = resources.FlowLogs{
		FlowLogs: ionoscloud.FlowLogs{
			Id:    &TestFlowLogVar,
			Items: &[]ionoscloud.FlowLog{TestFlowLog.FlowLog},
		},
	}

	// Datacenter test fixtures (used by list-all tests in multiple packages)

	TestDatacenterVar = "test-datacenter"
	TestDcVersion     = int32(1)

	TestDc = ionoscloud.Datacenter{
		Id: &TestDatacenterVar,
		Properties: &ionoscloud.DatacenterProperties{
			Name:     &TestDatacenterVar,
			Location: &TestDatacenterVar,
			Version:  &TestDcVersion,
		},
	}

	TestDcs = resources.Datacenters{
		Datacenters: ionoscloud.Datacenters{
			Id:    &TestDatacenterVar,
			Items: &[]ionoscloud.Datacenter{TestDc, TestDc},
		},
	}
)
