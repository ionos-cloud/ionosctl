package commands

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/ionos-cloud/ionosctl/pkg/builder"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/resources"
	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
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
				},
			},
		},
		Href: &testRequestPathVar,
	}
	rqs = resources.Requests{
		Requests: ionoscloud.Requests{
			Id:    &testRequestVar,
			Items: &[]ionoscloud.Request{rq},
		},
	}
	testRequestVar     = "test-request"
	testRequestPathVar = fmt.Sprintf("https://api.ionos.com/cloudapi/v5/requests/%s", testRequestVar)
	testRequestErr     = errors.New("request test: error occurred")
)

func TestPreRunRequestIdValidate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgRequestId), testRequestVar)
		viper.Set(config.ArgQuiet, false)
		err := PreRunRequestIdValidate(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunRequestIdValidateRequiredFlagErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgRequestId), "")
		viper.Set(config.ArgQuiet, false)
		err := PreRunRequestIdValidate(cfg)
		assert.Error(t, err)
		assert.True(t, err.Error() == clierror.NewRequiredFlagErr(config.ArgRequestId).Error())
	})
}

func TestRunRequestList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		rm.Request.EXPECT().List().Return(rqs, nil, nil)
		viper.Set(config.ArgQuiet, false)
		err := RunRequestList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunRequestListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		rm.Request.EXPECT().List().Return(rqs, nil, testRequestErr)
		viper.Set(config.ArgQuiet, false)
		err := RunRequestList(cfg)
		assert.Error(t, err)
		assert.True(t, err == testRequestErr)
	})
}

func TestRunRequestGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgRequestId), testRequestVar)
		req := resources.Request{rq}
		rm.Request.EXPECT().Get(testRequestVar).Return(&req, nil, nil)
		viper.Set(config.ArgQuiet, false)
		err := RunRequestGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunRequestGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgRequestId), testRequestVar)
		req := resources.Request{rq}
		rm.Request.EXPECT().Get(testRequestVar).Return(&req, nil, testRequestErr)
		viper.Set(config.ArgQuiet, false)
		err := RunRequestGet(cfg)
		assert.Error(t, err)
		assert.True(t, err == testRequestErr)
	})
}

func TestRunRequestWait(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgRequestId), testRequestVar)
		req := resources.Request{rq}
		rm.Request.EXPECT().Get(testRequestVar).Return(&req, nil, nil)
		rm.Request.EXPECT().Wait(testRequestPathVar+"/status").Return(nil, nil)
		viper.Set(config.ArgQuiet, false)
		err := RunRequestWait(cfg)
		assert.NoError(t, err)
	})
}

func TestRunRequestWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgRequestId), testRequestVar)
		req := resources.Request{rq}
		rm.Request.EXPECT().Get(testRequestVar).Return(&req, nil, nil)
		rm.Request.EXPECT().Wait(testRequestPathVar+"/status").Return(nil, testRequestErr)
		viper.Set(config.ArgQuiet, false)
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
	viper.Set(builder.GetGlobalFlagName("request", config.ArgCols), []string{"RequestId"})
	getRequestsCols(builder.GetGlobalFlagName("request", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
}

func TestGetRequestsColsErr(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() { return }

	w := bufio.NewWriter(&b)
	viper.Set(builder.GetGlobalFlagName("request", config.ArgCols), []string{"Unknown"})
	getRequestsCols(builder.GetGlobalFlagName("request", config.ArgCols), w)
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
