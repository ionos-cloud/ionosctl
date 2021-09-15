package commands

import (
	"bufio"
	"bytes"
	"errors"
	"regexp"
	"strings"
	"testing"

	"github.com/ionos-cloud/ionosctl/internal/config"
	"github.com/ionos-cloud/ionosctl/internal/core"
	"github.com/ionos-cloud/ionosctl/internal/utils/clierror"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/services/cloudapi-v6/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	loc = resources.Location{
		Location: ionoscloud.Location{
			Id: &testLocationVar,
			Properties: &ionoscloud.LocationProperties{
				Name:         &testLocationVar,
				Features:     &[]string{testLocationVar},
				ImageAliases: &[]string{testLocationVar},
				CpuArchitecture: &[]ionoscloud.CpuArchitectureProperties{
					{
						CpuFamily: &testLocationVar,
						MaxRam:    &testLocationIntVar,
						MaxCores:  &testLocationIntVar,
						Vendor:    &testLocationVar,
					},
				},
			},
		},
	}
	locations = resources.Locations{
		Locations: ionoscloud.Locations{
			Id:    &testLocationVar,
			Items: &[]ionoscloud.Location{loc.Location},
		},
	}
	testLocationIntVar = int32(1)
	testLocationVar    = "test/location"
	testLocationErr    = errors.New("location test error occurred")
)

func TestLocationCmd(t *testing.T) {
	var err error
	core.RootCmdTest.AddCommand(LocationCmd())
	if ok := LocationCmd().IsAvailableCommand(); !ok {
		err = errors.New("non-available cmd")
	}
	assert.NoError(t, err)
}

func TestPreLocationId(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLocationId), testLocationVar)
		err := PreRunLocationId(cfg)
		assert.NoError(t, err)
	})
}

func TestPreLocationIdErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		err := PreRunLocationId(cfg)
		assert.Error(t, err)
	})
}

func TestRunLocationList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		rm.CloudApiV6Mocks.Location.EXPECT().List().Return(locations, nil, nil)
		err := RunLocationList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunLocationListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		rm.CloudApiV6Mocks.Location.EXPECT().List().Return(locations, nil, testLocationErr)
		err := RunLocationList(cfg)
		assert.Error(t, err)
	})
}

func TestRunLocationGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLocationId), testLocationVar)
		testIds := strings.Split(testLocationVar, "/")
		rm.CloudApiV6Mocks.Location.EXPECT().GetByRegionAndLocationId(testIds[0], testIds[1]).Return(&loc, nil, nil)
		err := RunLocationGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunLocationGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLocationId), testLocationVar)
		testIds := strings.Split(testLocationVar, "/")
		rm.CloudApiV6Mocks.Location.EXPECT().GetByRegionAndLocationId(testIds[0], testIds[1]).Return(&loc, nil, testLocationErr)
		err := RunLocationGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunLocationGetLocationErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLocationId), "us")
		err := RunLocationGet(cfg)
		assert.Error(t, err)
	})
}

func TestGetLocationsCols(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() { return }
	w := bufio.NewWriter(&b)
	viper.Set(core.GetGlobalFlagName("location", config.ArgCols), []string{"Name"})
	getLocationCols(core.GetGlobalFlagName("location", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
}

func TestGetLocationsColsErr(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() { return }
	w := bufio.NewWriter(&b)
	viper.Set(core.GetGlobalFlagName("location", config.ArgCols), []string{"Unknown"})
	getLocationCols(core.GetGlobalFlagName("location", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`unknown column Unknown`)
	assert.True(t, re.Match(b.Bytes()))
}
