package utils

//import (
//	"bufio"
//	"bytes"
//	"errors"
//	"fmt"
//	"testing"
//
//	"github.com/golang/mock/gomock"
//	"github.com/ionos-cloud/ionosctl/pkg/builder"
//	"github.com/ionos-cloud/ionosctl/pkg/config"
//	mockprinter "github.com/ionos-cloud/ionosctl/pkg/utils/printer/mocks"
//	"github.com/spf13/viper"
//	"github.com/stretchr/testify/assert"
//)
//
//var (
//	pathRequest           = fmt.Sprintf("%s/%s/status/test/test", config.DefaultApiURL, testWaitForRequestVar)
//	testWaitForRequestVar = "test-wait-for-action"
//	testWaitForRequestErr = errors.New("wait-for-action test error occurred")
//)
//
//func TestWaitForRequest(t *testing.T) {
//	var b bytes.Buffer
//	w := bufio.NewWriter(&b)
//	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
//		ctrl := gomock.NewController(t)
//		defer ctrl.Finish()
//		p := mockprinter.NewMockPrintService(ctrl)
//		p.EXPECT().Print("Waiting for request: test").Return(nil)
//
//		cfg.Printer = p
//		viper.Set(config.ArgQuiet, false)
//		viper.Set(config.ArgOutput, "text")
//		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWaitForRequest), true)
//		rm.Request.EXPECT().Wait(pathRequest).Return(nil, nil)
//		err := WaitForRequest(cfg, pathRequest)
//		assert.NoError(t, err)
//	})
//}
//
//func TestWaitForRequestErr(t *testing.T) {
//	var b bytes.Buffer
//	w := bufio.NewWriter(&b)
//	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
//		ctrl := gomock.NewController(t)
//		defer ctrl.Finish()
//		p := mockprinter.NewMockPrintService(ctrl)
//		p.EXPECT().Print("Waiting for request: test").Return(nil)
//
//		cfg.Printer = p
//		viper.Set(config.ArgQuiet, false)
//		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWaitForRequest), true)
//		rm.Request.EXPECT().Wait(pathRequest).Return(nil, testWaitForRequestErr)
//		err := WaitForRequest(cfg, pathRequest)
//		assert.Error(t, err)
//	})
//}
//
//func TestWaitForRequestIdErr(t *testing.T) {
//	var b bytes.Buffer
//	w := bufio.NewWriter(&b)
//	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
//		viper.Set(config.ArgQuiet, false)
//		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWaitForRequest), true)
//		err := WaitForRequest(cfg, "")
//		assert.Error(t, err)
//	})
//}
//
//func TestWaitForRequestPathErr(t *testing.T) {
//	var b bytes.Buffer
//	w := bufio.NewWriter(&b)
//	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
//		viper.Set(config.ArgQuiet, false)
//		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWaitForRequest), false)
//		err := WaitForRequest(cfg, pathRequest)
//		assert.NoError(t, err)
//	})
//}
