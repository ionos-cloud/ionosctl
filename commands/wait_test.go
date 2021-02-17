package commands

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/ionos-cloud/ionosctl/pkg/builder"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	pathRequest          = fmt.Sprintf("%s/%s/status/", config.DefaultApiURL, testWaitForActionVar)
	testWaitForActionVar = "test-waitForAction"
	testWaitForActionErr = errors.New("waitForAction test error occurred")
)

func TestWaitForAction(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), true)
		rm.Request.EXPECT().Wait(pathRequest).Return(nil, nil)
		err := waitForAction(cfg, pathRequest)
		assert.NoError(t, err)
	})
}

func TestWaitForAction_Err(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), true)
		rm.Request.EXPECT().Wait(pathRequest).Return(nil, testWaitForActionErr)
		err := waitForAction(cfg, pathRequest)
		assert.Error(t, err)
	})
}

func TestWaitForAction_PathErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), false)
		err := waitForAction(cfg, pathRequest)
		assert.NoError(t, err)
	})
}
