package cloudapi_v5

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/ionos-cloud/ionosctl/internal/config"
	"github.com/ionos-cloud/ionosctl/internal/core"
	"github.com/ionos-cloud/ionosctl/internal/utils/clierror"
	cloudapiv5 "github.com/ionos-cloud/ionosctl/services/cloudapi-v5"
	"github.com/ionos-cloud/ionosctl/services/cloudapi-v5/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
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
			CreatedDate: &testIonosTime,
		},
	}
	testRequestUpdatedPatch = ionoscloud.Request{
		Properties: &ionoscloud.RequestProperties{
			Method: &testRequestMethodPatch,
		},
		Metadata: &ionoscloud.RequestMetadata{
			CreatedDate: &testIonosTime,
		},
	}
	testRequestDeleted = ionoscloud.Request{
		Properties: &ionoscloud.RequestProperties{
			Method: &testRequestMethodDelete,
		},
		Metadata: &ionoscloud.RequestMetadata{
			CreatedDate: &testIonosTime,
		},
	}
	testRequestCreated = ionoscloud.Request{
		Properties: &ionoscloud.RequestProperties{
			Method: &testRequestMethodPost,
		},
		Metadata: &ionoscloud.RequestMetadata{
			CreatedDate: &testIonosTime,
		},
	}
	testRequestStatus = resources.RequestStatus{
		RequestStatus: ionoscloud.RequestStatus{
			Id: &testRequestVar,
			Metadata: &ionoscloud.RequestStatusMetadata{
				Status:  &testRequestStatusVar,
				Message: &testRequestVar,
			},
		},
	}
	testRequests = resources.Requests{
		Requests: ionoscloud.Requests{
			Id:    &testRequestVar,
			Items: &[]ionoscloud.Request{rq, testRequestUpdated, testRequestUpdatedPatch, testRequestDeleted, testRequestCreated},
		},
	}
	testRequestStatusVar    = "DONE"
	testIonosTime           = ionoscloud.IonosTime{Time: time.Date(2021, 1, 1, 0, 0, 0, 0, time.Now().Location())}
	testRequestVar          = "test-request"
	testRequestMethodPut    = "PUT"
	testRequestMethodPatch  = "PATCH"
	testRequestMethodDelete = "DELETE"
	testRequestMethodPost   = "POST"
	testRequestPathVar      = fmt.Sprintf("https://api.ionos.com/cloudapi/v5/requests/%s", testRequestVar)
	testRequestErr          = errors.New("request test: error occurred")
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
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgRequestId), testRequestVar)
		viper.Set(config.ArgQuiet, false)
		err := PreRunRequestList(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunRequestListFilters(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgRequestId), testRequestVar)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgFilters), []string{fmt.Sprintf("createdBy=%s", testQueryParamVar)})
		err := PreRunRequestList(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunRequestListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgRequestId), testRequestVar)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgFilters), []string{fmt.Sprintf("%s=%s", testQueryParamVar, testQueryParamVar)})
		err := PreRunRequestList(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunRequestId(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgRequestId), testRequestVar)
		viper.Set(config.ArgQuiet, false)
		err := PreRunRequestId(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunRequestIdRequiredFlagErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		err := PreRunRequestId(cfg)
		assert.Error(t, err)
	})
}

func TestRunRequestList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, true)
		viper.Set(core.GetGlobalFlagName(cfg.Namespace, config.ArgCols), allRequestCols)
		rm.CloudApiV5Mocks.Request.EXPECT().List(resources.ListQueryParams{}).Return(testRequests, &testResponse, nil)
		err := RunRequestList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunRequestListQueryParams(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, true)
		viper.Set(core.GetGlobalFlagName(cfg.Namespace, config.ArgCols), allRequestCols)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgFilters), []string{fmt.Sprintf("%s=%s", testQueryParamVar, testQueryParamVar)})
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgOrderBy), testQueryParamVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgMaxResults), testMaxResultsVar)
		rm.CloudApiV5Mocks.Request.EXPECT().List(testListQueryParam).Return(resources.Requests{}, &testResponse, nil)
		err := RunRequestList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunRequestListSortedUpdate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgLatest), 10)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgMethod), "UPDATE")
		rm.CloudApiV5Mocks.Request.EXPECT().List(resources.ListQueryParams{}).Return(testRequests, nil, nil)
		err := RunRequestList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunRequestListSortedCreate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgLatest), 10)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgMethod), "CREATE")
		rm.CloudApiV5Mocks.Request.EXPECT().List(resources.ListQueryParams{}).Return(testRequests, nil, nil)
		err := RunRequestList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunRequestListSortedDelete(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgLatest), 1)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgMethod), "DELETE")
		rm.CloudApiV5Mocks.Request.EXPECT().List(resources.ListQueryParams{}).Return(testRequests, nil, nil)
		err := RunRequestList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunRequestListSortedErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgLatest), 10)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgMethod), "no method")
		rm.CloudApiV5Mocks.Request.EXPECT().List(resources.ListQueryParams{}).Return(testRequests, nil, nil)
		err := RunRequestList(cfg)
		assert.Error(t, err)
	})
}

func TestRunRequestListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		rm.CloudApiV5Mocks.Request.EXPECT().List(resources.ListQueryParams{}).Return(testRequests, nil, testRequestErr)
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
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgRequestId), testRequestVar)
		req := resources.Request{Request: rq}
		rm.CloudApiV5Mocks.Request.EXPECT().Get(testRequestVar).Return(&req, &testResponse, nil)
		err := RunRequestGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunRequestGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgRequestId), testRequestVar)
		req := resources.Request{Request: rq}
		rm.CloudApiV5Mocks.Request.EXPECT().Get(testRequestVar).Return(&req, nil, testRequestErr)
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
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgRequestId), testRequestVar)
		req := resources.Request{Request: rq}
		rm.CloudApiV5Mocks.Request.EXPECT().Get(testRequestVar).Return(&req, nil, nil)
		rm.CloudApiV5Mocks.Request.EXPECT().Wait(testRequestPathVar+"/status").Return(nil, nil)
		err := RunRequestWait(cfg)
		assert.NoError(t, err)
	})
}

func TestRunRequestWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgRequestId), testRequestVar)
		req := resources.Request{Request: rq}
		rm.CloudApiV5Mocks.Request.EXPECT().Get(testRequestVar).Return(&req, nil, nil)
		rm.CloudApiV5Mocks.Request.EXPECT().Wait(testRequestPathVar+"/status").Return(nil, testRequestErr)
		err := RunRequestWait(cfg)
		assert.Error(t, err)
		assert.True(t, err == testRequestErr)
	})
}

func TestGetRequestsCols(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() { return }
	w := bufio.NewWriter(&b)
	viper.Set(core.GetGlobalFlagName("request", config.ArgCols), []string{"RequestId"})
	getRequestsCols(core.GetGlobalFlagName("request", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
}

func TestGetRequestsColsErr(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() { return }
	w := bufio.NewWriter(&b)
	viper.Set(core.GetGlobalFlagName("request", config.ArgCols), []string{"Unknown"})
	getRequestsCols(core.GetGlobalFlagName("request", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`unknown column Unknown`)
	assert.True(t, re.Match(b.Bytes()))
}
