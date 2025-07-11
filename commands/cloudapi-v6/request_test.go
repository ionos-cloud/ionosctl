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
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	rq = ionoscloud.Request{
		Id: &testRequestVar,
		Metadata: &ionoscloud.RequestMetadata{
			RequestStatus: &ionoscloud.RequestStatus{
				Metadata: &ionoscloud.RequestStatusMetadata{
					Status:  &testRequestVar,
					Message: &testRequestVar,
					Targets: &[]ionoscloud.RequestTarget{
						{
							Target: &ionoscloud.ResourceReference{
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
		Properties: &ionoscloud.RequestProperties{
			Url:    &testRequestVar,
			Body:   &testRequestVar,
			Method: &testRequestVar,
		},
	}
	testRequestUpdated = ionoscloud.Request{
		Properties: &ionoscloud.RequestProperties{
			Method: &testRequestMethodPut,
		},
		Metadata: &ionoscloud.RequestMetadata{
			CreatedDate:   &testIonosTime,
			RequestStatus: &testRequestStatus.RequestStatus,
		},
	}
	testRequestUpdatedPatch = ionoscloud.Request{
		Properties: &ionoscloud.RequestProperties{
			Method: &testRequestMethodPatch,
		},
		Metadata: &ionoscloud.RequestMetadata{
			CreatedDate:   &testIonosTime,
			RequestStatus: &testRequestStatus.RequestStatus,
		},
	}
	testRequestDeleted = ionoscloud.Request{
		Properties: &ionoscloud.RequestProperties{
			Method: &testRequestMethodDelete,
		},
		Metadata: &ionoscloud.RequestMetadata{
			CreatedDate:   &testIonosTime,
			RequestStatus: &testRequestStatus.RequestStatus,
		},
	}
	testRequestCreated = ionoscloud.Request{
		Properties: &ionoscloud.RequestProperties{
			Method: &testRequestMethodPost,
		},
		Metadata: &ionoscloud.RequestMetadata{
			CreatedDate:   &testIonosTime,
			RequestStatus: &testRequestStatus.RequestStatus,
		},
	}
	testRequestStatus = resources.RequestStatus{
		RequestStatus: ionoscloud.RequestStatus{
			Id: &testRequestVar,
			Metadata: &ionoscloud.RequestStatusMetadata{
				Status:  &testRequestStatusVar,
				Message: &testRequestVar,
				Targets: &testRequestTargetsVar,
			},
		},
	}
	testRequests = resources.Requests{
		Requests: ionoscloud.Requests{
			Id:    &testRequestVar,
			Items: &[]ionoscloud.Request{rq, testRequestUpdated, testRequestUpdatedPatch, testRequestDeleted, testRequestCreated},
		},
	}
	testRequestTargetsVar = []ionoscloud.RequestTarget{
		{
			Target: &ionoscloud.ResourceReference{
				Id:   &testRequestVar,
				Type: &testTypeRequestVar,
			},
		},
	}
	testRequestStatusVar    = "DONE"
	testRequestVar          = "test-request"
	testRequestPathVar      = fmt.Sprintf("https://api.ionos.com/cloudapi/v6/requests/%s", testRequestVar)
	testRequestErr          = errors.New("request test: error occurred")
	testIonosTime           = ionoscloud.IonosTime{Time: time.Date(2021, 1, 1, 0, 0, 0, 0, time.Now().Location())}
	testRequestMethodPut    = "PUT"
	testRequestMethodPatch  = "PATCH"
	testRequestMethodDelete = "DELETE"
	testRequestMethodPost   = "POST"
	testTypeRequestVar      = ionoscloud.Type("datacenter")
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
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagRequestId), testRequestVar)
		viper.Set(constants.FlagQuiet, false)
		err := PreRunRequestList(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunRequestListFilters(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagRequestId), testRequestVar)
		viper.Set(constants.FlagQuiet, false)
		cfg.Command.Command.Flags().Set(cloudapiv6.FlagFilters, fmt.Sprintf("createdBy=%s", testQueryParamVar))
		err := PreRunRequestList(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunRequestListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagRequestId), testRequestVar)
		viper.Set(constants.FlagQuiet, false)
		cfg.Command.Command.Flags().Set(cloudapiv6.FlagFilters, fmt.Sprintf("%s=%s", testQueryParamVar, testQueryParamVar))
		err := PreRunRequestList(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunRequestId(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagRequestId), testRequestVar)
		viper.Set(constants.FlagQuiet, false)
		err := PreRunRequestId(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunRequestIdRequiredFlagErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		err := PreRunRequestId(cfg)
		assert.Error(t, err)
	})
}

func TestRunRequestList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagVerbose, false)
		viper.Set(core.GetFlagName(cfg.Namespace, constants.FlagCols), allRequestCols)
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
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagVerbose, false)
		viper.Set(core.GetFlagName(cfg.Namespace, constants.FlagCols), allRequestCols)
		cfg.Command.Command.Flags().Set(cloudapiv6.FlagFilters, fmt.Sprintf("%s=%s", testQueryParamVar, testQueryParamVar))
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagOrderBy), testQueryParamVar)
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
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagLatest), 10)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagMethod), "UPDATE")
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
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagLatest), 10)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagMethod), "CREATE")
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
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagLatest), 1)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagMethod), "DELETE")
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
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagLatest), 10)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagMethod), "no method")
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
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
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
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagRequestId), testRequestVar)
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
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagRequestId), testRequestVar)
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
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagRequestId), testRequestVar)
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
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagRequestId), testRequestVar)
		req := resources.Request{Request: rq}
		rm.CloudApiV6Mocks.Request.EXPECT().Get(testRequestVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&req, nil, nil)
		rm.CloudApiV6Mocks.Request.EXPECT().Wait(testRequestPathVar+"/status").Return(nil, testRequestErr)
		err := RunRequestWait(cfg)
		assert.Error(t, err)
		assert.True(t, err == testRequestErr)
	})
}
