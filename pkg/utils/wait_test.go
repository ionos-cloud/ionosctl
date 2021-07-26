package utils

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/ionos-cloud/ionosctl/pkg/resources/v5"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	mockprinter "github.com/ionos-cloud/ionosctl/pkg/utils/printer/mocks"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	pathRequest              = fmt.Sprintf("%s/%s/status/test/test", config.DefaultApiURL, testWaitForRequestVar)
	testWaitForRequestVar    = "test-wait-for-action"
	testRunningRequestStatus = &v5.RequestStatus{
		RequestStatus: ionoscloud.RequestStatus{
			Id: &testVar,
			Metadata: &ionoscloud.RequestStatusMetadata{
				Status: &testRunningStateVar,
			},
		},
	}
	testDoneRequestStatus = &v5.RequestStatus{
		RequestStatus: ionoscloud.RequestStatus{
			Id: &testVar,
			Metadata: &ionoscloud.RequestStatusMetadata{
				Status: &testDoneStateVar,
			},
		},
	}
	testQueuedRequestStatus = &v5.RequestStatus{
		RequestStatus: ionoscloud.RequestStatus{
			Id: &testVar,
			Metadata: &ionoscloud.RequestStatusMetadata{
				Status: &testQueuedStateVar,
			},
		},
	}
	testFailedRequestStatus = &v5.RequestStatus{
		RequestStatus: ionoscloud.RequestStatus{
			Id: &testVar,
			Metadata: &ionoscloud.RequestStatusMetadata{
				Status: &testFailedStateVar,
			},
		},
	}
	testRunningStateVar = "RUNNING"
	testDoneStateVar    = "DONE"
	testQueuedStateVar  = "QUEUED"
	testFailedStateVar  = "FAILED"
)

func TestNoWaitForRequest(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		p := mockprinter.NewMockPrintService(ctrl)

		cfg.Printer = p
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, "text")
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		err := WaitForRequest(cfg, pathRequest)
		assert.NoError(t, err)
	})
}

func TestWaitForRequestIdErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		p := mockprinter.NewMockPrintService(ctrl)

		cfg.Printer = p
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, "text")
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		err := WaitForRequest(cfg, testVar)
		assert.Error(t, err)
	})
}

func TestWaitForRequest(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		p := mockprinter.NewMockPrintService(ctrl)
		p.EXPECT().GetStdout().Return(nil)

		cfg.Printer = p
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, "text")
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		rm.Request.EXPECT().GetStatus(testVar).Return(testQueuedRequestStatus, nil, nil)
		rm.Request.EXPECT().GetStatus(testVar).Return(testRunningRequestStatus, nil, nil)
		rm.Request.EXPECT().GetStatus(testVar).Return(testDoneRequestStatus, nil, nil)
		err := WaitForRequest(cfg, pathRequest)
		assert.NoError(t, err)
	})
}

func TestWaitForRequestErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		p := mockprinter.NewMockPrintService(ctrl)
		p.EXPECT().GetStdout().Return(os.Stdout)

		cfg.Printer = p
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, "text")
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		rm.Request.EXPECT().GetStatus(testVar).Return(testQueuedRequestStatus, nil, nil)
		rm.Request.EXPECT().GetStatus(testVar).Return(testRunningRequestStatus, nil, nil)
		rm.Request.EXPECT().GetStatus(testVar).Return(testFailedRequestStatus, nil, nil)
		err := WaitForRequest(cfg, pathRequest)
		assert.Error(t, err)
	})
}

func TestWaitForRequestJson(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		p := mockprinter.NewMockPrintService(ctrl)
		p.EXPECT().Print(waitingForRequestMsg).Return(nil)
		p.EXPECT().Print(done).Return(nil)

		cfg.Printer = p
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, "json")
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		rm.Request.EXPECT().GetStatus(testVar).Return(testQueuedRequestStatus, nil, nil)
		rm.Request.EXPECT().GetStatus(testVar).Return(testRunningRequestStatus, nil, nil)
		rm.Request.EXPECT().GetStatus(testVar).Return(testDoneRequestStatus, nil, nil)
		err := WaitForRequest(cfg, pathRequest)
		assert.NoError(t, err)
	})
}

func TestWaitForRequestJsonStatusErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		p := mockprinter.NewMockPrintService(ctrl)
		p.EXPECT().Print(waitingForRequestMsg).Return(nil)
		p.EXPECT().Print(failed).Return(nil)

		cfg.Printer = p
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, "json")
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		rm.Request.EXPECT().GetStatus(testVar).Return(nil, nil, nil)
		err := WaitForRequest(cfg, pathRequest)
		assert.Error(t, err)
	})
}

var (
	testInterrogatorFuncErr       = func(c *core.CommandConfig, resourceId string) (*string, error) { return nil, nil }
	testInterrogatorFailedFunc    = func(c *core.CommandConfig, resourceId string) (*string, error) { return &testFailedStateVar, nil }
	testInterrogatorAvailableFunc = func(c *core.CommandConfig, resourceId string) (*string, error) { return &testAvailableStateVar, nil }
	testInterrogatorReadyFunc     = func(c *core.CommandConfig, resourceId string) (*string, error) { return &testReadyStateVar, nil }
	testInterrogatorActiveFunc    = func(c *core.CommandConfig, resourceId string) (*string, error) { return &testActiveStateVar, nil }
	testReadyStateVar             = stateReadyStatus
	testAvailableStateVar         = stateAvailableStatus
	testActiveStateVar            = stateActiveStatus
)

func TestNoWaitForState(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		p := mockprinter.NewMockPrintService(ctrl)
		cfg.Printer = p
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, "text")
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForState), false)
		err := WaitForState(cfg, testInterrogatorFuncErr, pathRequest)
		assert.NoError(t, err)
	})
}

func TestWaitForStateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		p := mockprinter.NewMockPrintService(ctrl)
		p.EXPECT().GetStdout().Return(nil)

		cfg.Printer = p
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, "text")
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForState), true)
		err := WaitForState(cfg, testInterrogatorFuncErr, pathRequest)
		assert.Error(t, err)
	})
}

func TestWaitForStateFailedErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		p := mockprinter.NewMockPrintService(ctrl)
		p.EXPECT().GetStdout().Return(nil)

		cfg.Printer = p
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, "text")
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForState), true)
		err := WaitForState(cfg, testInterrogatorFailedFunc, pathRequest)
		assert.Error(t, err)
	})
}

func TestWaitForState(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		p := mockprinter.NewMockPrintService(ctrl)
		p.EXPECT().GetStdout().Return(nil)

		cfg.Printer = p
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, "text")
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForState), true)
		err := WaitForState(cfg, testInterrogatorActiveFunc, pathRequest)
		assert.NoError(t, err)
	})
}

func TestWaitForStateAvailable(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		p := mockprinter.NewMockPrintService(ctrl)
		p.EXPECT().GetStdout().Return(nil)

		cfg.Printer = p
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, "text")
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForState), true)
		err := WaitForState(cfg, testInterrogatorAvailableFunc, pathRequest)
		assert.NoError(t, err)
	})
}

func TestWaitForStateReady(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		p := mockprinter.NewMockPrintService(ctrl)
		p.EXPECT().GetStdout().Return(nil)

		cfg.Printer = p
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, "text")
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForState), true)
		err := WaitForState(cfg, testInterrogatorReadyFunc, pathRequest)
		assert.NoError(t, err)
	})
}

func TestWaitForStateJson(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		p := mockprinter.NewMockPrintService(ctrl)
		p.EXPECT().Print(waitingForStateMsg).Return(nil)
		p.EXPECT().Print(done).Return(nil)

		cfg.Printer = p
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, "json")
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForState), true)
		err := WaitForState(cfg, testInterrogatorActiveFunc, pathRequest)
		assert.NoError(t, err)
	})
}

func TestWaitForStateJsonErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		p := mockprinter.NewMockPrintService(ctrl)
		p.EXPECT().Print(waitingForStateMsg).Return(nil)
		p.EXPECT().Print(failed).Return(nil)

		cfg.Printer = p
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, "json")
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForState), true)
		err := WaitForState(cfg, testInterrogatorFailedFunc, pathRequest)
		assert.Error(t, err)
	})
}
