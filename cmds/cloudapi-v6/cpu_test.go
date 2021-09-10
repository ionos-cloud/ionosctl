package commands

import (
	"bufio"
	"bytes"
	"errors"
	"regexp"
	"strings"
	"testing"

	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	testLocationCpuVar = "test/location"
	testLocationCpuErr = errors.New("location cpu test error occurred")
)

func TestRunLocationCpuList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLocationId), testLocationCpuVar)
		testIds := strings.Split(testLocationCpuVar, "/")
		rm.Location.EXPECT().GetByRegionAndLocationId(testIds[0], testIds[1]).Return(&loc, nil, nil)
		err := RunLocationCpuList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunLocationCpuListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLocationId), testLocationCpuVar)
		testIds := strings.Split(testLocationCpuVar, "/")
		rm.Location.EXPECT().GetByRegionAndLocationId(testIds[0], testIds[1]).Return(&loc, nil, testLocationCpuErr)
		err := RunLocationCpuList(cfg)
		assert.Error(t, err)
	})
}

func TestGetCpusCols(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() { return }
	w := bufio.NewWriter(&b)
	viper.Set(core.GetGlobalFlagName("cpu", config.ArgCols), []string{"Vendor"})
	getCpuCols(core.GetGlobalFlagName("cpu", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
}

func TestGetCpusColsErr(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() { return }
	w := bufio.NewWriter(&b)
	viper.Set(core.GetGlobalFlagName("cpu", config.ArgCols), []string{"Unknown"})
	getCpuCols(core.GetGlobalFlagName("cpu", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`unknown column Unknown`)
	assert.True(t, re.Match(b.Bytes()))
}
