package commands

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"

	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
	compute "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	rq = compute.Request{
		Id: &testRequestVar,
		Metadata: &compute.RequestMetadata{
			RequestStatus: &compute.RequestStatus{
				Metadata: &compute.RequestStatusMetadata{
					Status:  &testRequestVar,
					Message: &testRequestVar,
					Targets: &[]compute.RequestTarget{
						{
							Target: &compute.ResourceReference{
								Id:   &testRequestVar,
								Type: &testTypeRequestVar,
							},
						},
					},
				},
			},
			CreatedDate: &testIonosTime,
			CreatedBy:   &testRequestVar,
		},
		Href: &testRequestPathVar,
		Properties: &compute.RequestProperties{
			Url:    &testRequestVar,
			Body:   &testRequestVar,
			Method: &testRequestVar,
		},
	}
	testRequestUpdated = compute.Request{
		Properties: &compute.RequestProperties{
			Method: &testRequestMethodPut,
		},
		Metadata: &compute.RequestMetadata{
			CreatedDate:   &testIonosTime,
			RequestStatus: &testRequestStatus.RequestStatus,
		},
	}
	testRequestUpdatedPatch = compute.Request{
		Properties: &compute.RequestProperties{
			Method: &testRequestMethodPatch,
		},
		Metadata: &compute.RequestMetadata{
			CreatedDate:   &testIonosTime,
			RequestStatus: &testRequestStatus.RequestStatus,
		},
	}
	testRequestDeleted = compute.Request{
		Properties: &compute.RequestProperties{
			Method: &testRequestMethodDelete,
		},
		Metadata: &compute.RequestMetadata{
			CreatedDate:   &testIonosTime,
			RequestStatus: &testRequestStatus.RequestStatus,
		},
	}
	testRequestCreated = compute.Request{
		Properties: &compute.RequestProperties{
			Method: &testRequestMethodPost,
		},
		Metadata: &compute.RequestMetadata{
			CreatedDate:   &testIonosTime,
			RequestStatus: &testRequestStatus.RequestStatus,
		},
	}
	testRequestStatus = resources.RequestStatus{
		RequestStatus: compute.RequestStatus{
			Id: &testRequestVar,
			Metadata: &compute.RequestStatusMetadata{
				Status:  &testRequestStatusVar,
				Message: &testRequestVar,
				Targets: &testRequestTargetsVar,
			},
		},
	}
	testRequests = resources.Requests{
		Requests: compute.Requests{
			Id:    &testRequestVar,
			Items: &[]compute.Request{rq, testRequestUpdated, testRequestUpdatedPatch, testRequestDeleted, testRequestCreated},
		},
	}
	testRequestTargetsVar = []compute.RequestTarget{
		{
			Target: &compute.ResourceReference{
				Id:   &testRequestVar,
				Type: &testTypeRequestVar,
			},
		},
	}
	testRequestStatusVar    = "DONE"
	testRequestVar          = "test-request"
	testRequestPathVar      = fmt.Sprintf("https://api.ionos.com/cloudapi/v6/requests/%s", testRequestVar)
	testRequestErr          = errors.New("request test: error occurred")
	testIonosTime           = compute.IonosTime{Time: time.Date(2021, 1, 1, 0, 0, 0, 0, time.Now().Location())}
	testRequestMethodPut    = "PUT"
	testRequestMethodPatch  = "PATCH"
	testRequestMethodDelete = "DELETE"
	testRequestMethodPost   = "POST"
	testTypeRequestVar      = compute.Type("datacenter")
)

func TestRequestCmd(t *testing.T) {
	var err error
	core.RootCmdTest.AddCommand(RequestCmd())
	if ok := RequestCmd().IsAvailableCommand(); !ok {
		err = errors.New("non-available cmd")
	}
	assert.NoError(t, err)
}

func TestPreRunRequestList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRequestId), testRequestVar)
		viper.Set(constants.ArgQuiet, false)
		err := PreRunRequestList(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunRequestListFilters(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRequestId), testRequestVar)
		viper.Set(constants.ArgQuiet, false)
		cfg.Command.Command.Flags().Set(cloudapiv6.ArgFilters, fmt.Sprintf("createdBy=%s", testQueryParamVar))
		err := PreRunRequestList(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunRequestListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRequestId), testRequestVar)
		viper.Set(constants.ArgQuiet, false)
		cfg.Command.Command.Flags().Set(cloudapiv6.ArgFilters, fmt.Sprintf("%s=%s", testQueryParamVar, testQueryParamVar))
		err := PreRunRequestList(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunRequestId(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRequestId), testRequestVar)
		viper.Set(constants.ArgQuiet, false)
		err := PreRunRequestId(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunRequestIdRequiredFlagErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		err := PreRunRequestId(cfg)
		assert.Error(t, err)
	})
}

func TestRunRequestList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.Namespace, constants.ArgCols), allRequestCols)
		rm.CloudApiV6Mocks.Request.EXPECT().List(gomock.AssignableToTypeOf(testListQueryParam)).Return(testRequests, &testResponse, nil)
		err := RunRequestList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunRequestListQueryParams(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.Namespace, constants.ArgCols), allRequestCols)
		cfg.Command.Command.Flags().Set(cloudapiv6.ArgFilters, fmt.Sprintf("%s=%s", testQueryParamVar, testQueryParamVar))
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgOrderBy), testQueryParamVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagMaxResults), testMaxResultsVar)
		rm.CloudApiV6Mocks.Request.EXPECT().List(gomock.AssignableToTypeOf(testListQueryParam)).Return(testRequests, &testResponse, nil)
		err := RunRequestList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunRequestListSortedUpdate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLatest), 10)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgMethod), "UPDATE")
		rm.CloudApiV6Mocks.Request.EXPECT().List(gomock.AssignableToTypeOf(testListQueryParam)).Return(testRequests, nil, nil)
		err := RunRequestList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunRequestListSortedCreate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLatest), 10)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgMethod), "CREATE")
		rm.CloudApiV6Mocks.Request.EXPECT().List(gomock.AssignableToTypeOf(testListQueryParam)).Return(testRequests, nil, nil)
		err := RunRequestList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunRequestListSortedDelete(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLatest), 1)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgMethod), "DELETE")
		rm.CloudApiV6Mocks.Request.EXPECT().List(gomock.AssignableToTypeOf(testListQueryParam)).Return(testRequests, nil, nil)
		err := RunRequestList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunRequestListSortedErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLatest), 10)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgMethod), "no method")
		rm.CloudApiV6Mocks.Request.EXPECT().List(gomock.AssignableToTypeOf(testListQueryParam)).Return(testRequests, nil, nil)
		err := RunRequestList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunRequestListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		rm.CloudApiV6Mocks.Request.EXPECT().List(gomock.AssignableToTypeOf(testListQueryParam)).Return(testRequests, nil, testRequestErr)
		err := RunRequestList(cfg)
		assert.Error(t, err)
		assert.True(t, err == testRequestErr)
	})
}

func TestRunRequestGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRequestId), testRequestVar)
		req := resources.Request{Request: rq}
		rm.CloudApiV6Mocks.Request.EXPECT().Get(testRequestVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&req, &testResponse, nil)
		err := RunRequestGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunRequestGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRequestId), testRequestVar)
		req := resources.Request{Request: rq}
		rm.CloudApiV6Mocks.Request.EXPECT().Get(testRequestVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&req, nil, testRequestErr)
		err := RunRequestGet(cfg)
		assert.Error(t, err)
		assert.True(t, err == testRequestErr)
	})
}

func TestRunRequestWait(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRequestId), testRequestVar)
		req := resources.Request{Request: rq}
		rm.CloudApiV6Mocks.Request.EXPECT().Get(testRequestVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&req, nil, nil)
		rm.CloudApiV6Mocks.Request.EXPECT().Wait(testRequestPathVar+"/status").Return(nil, nil)
		err := RunRequestWait(cfg)
		assert.NoError(t, err)
	})
}

func TestRunRequestWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRequestId), testRequestVar)
		req := resources.Request{Request: rq}
		rm.CloudApiV6Mocks.Request.EXPECT().Get(testRequestVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&req, nil, nil)
		rm.CloudApiV6Mocks.Request.EXPECT().Wait(testRequestPathVar+"/status").Return(nil, testRequestErr)
		err := RunRequestWait(cfg)
		assert.Error(t, err)
		assert.True(t, err == testRequestErr)
	})
}
