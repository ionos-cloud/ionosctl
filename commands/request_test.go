package commands

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/resources"
	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
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
	testRequests = resources.Requests{
		Requests: ionoscloud.Requests{
			Id:    &testRequestVar,
			Items: &[]ionoscloud.Request{rq, testRequestUpdated, testRequestUpdatedPatch, testRequestDeleted, testRequestCreated},
		},
	}
	testRequestVar          = "test-request"
	testRequestPathVar      = fmt.Sprintf("https://api.ionos.com/cloudapi/v6/requests/%s", testRequestVar)
	testRequestErr          = errors.New("request test: error occurred")
	testIonosTime           = ionoscloud.IonosTime{Time: time.Date(2021, 1, 1, 0, 0, 0, 0, time.Now().Location())}
	testRequestMethodPut    = "PUT"
	testRequestMethodPatch  = "PATCH"
	testRequestMethodDelete = "DELETE"
	testRequestMethodPost   = "POST"
)

func TestPreRunRequestId(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgRequestId), testRequestVar)
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
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetGlobalFlagName(cfg.Namespace, config.ArgCols), allRequestCols)
		rm.Request.EXPECT().List().Return(testRequests, nil, nil)
		err := RunRequestList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunRequestListSortedUpdate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLatest), 10)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgMethod), "UPDATE")
		rm.Request.EXPECT().List().Return(testRequests, nil, nil)
		err := RunRequestList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunRequestListSortedCreate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLatest), 10)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgMethod), "CREATE")
		rm.Request.EXPECT().List().Return(testRequests, nil, nil)
		err := RunRequestList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunRequestListSortedDelete(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLatest), 1)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgMethod), "DELETE")
		rm.Request.EXPECT().List().Return(testRequests, nil, nil)
		err := RunRequestList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunRequestListSortedErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLatest), 10)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgMethod), "no method")
		rm.Request.EXPECT().List().Return(testRequests, nil, nil)
		err := RunRequestList(cfg)
		assert.Error(t, err)
	})
}

func TestRunRequestListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		rm.Request.EXPECT().List().Return(testRequests, nil, testRequestErr)
		err := RunRequestList(cfg)
		assert.Error(t, err)
		assert.True(t, err == testRequestErr)
	})
}

func TestRunRequestGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgRequestId), testRequestVar)
		req := resources.Request{Request: rq}
		rm.Request.EXPECT().Get(testRequestVar).Return(&req, nil, nil)
		err := RunRequestGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunRequestGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgRequestId), testRequestVar)
		req := resources.Request{Request: rq}
		rm.Request.EXPECT().Get(testRequestVar).Return(&req, nil, testRequestErr)
		err := RunRequestGet(cfg)
		assert.Error(t, err)
		assert.True(t, err == testRequestErr)
	})
}

func TestRunRequestWait(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgRequestId), testRequestVar)
		req := resources.Request{Request: rq}
		rm.Request.EXPECT().Get(testRequestVar).Return(&req, nil, nil)
		rm.Request.EXPECT().Wait(testRequestPathVar+"/status").Return(nil, nil)
		err := RunRequestWait(cfg)
		assert.NoError(t, err)
	})
}

func TestRunRequestWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgRequestId), testRequestVar)
		req := resources.Request{Request: rq}
		rm.Request.EXPECT().Get(testRequestVar).Return(&req, nil, nil)
		rm.Request.EXPECT().Wait(testRequestPathVar+"/status").Return(nil, testRequestErr)
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

func TestGetRequestsIds(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() { return }
	w := bufio.NewWriter(&b)
	viper.Set(config.ArgConfig, "../pkg/testdata/config.json")
	getRequestsIds(w)
	err := w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`401 Unauthorized`)
	assert.True(t, re.Match(b.Bytes()))
}
