package commands

import (
	"bufio"
	"bytes"
	"errors"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/resources/v5"
	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	loc = v5.Location{
		Location: ionoscloud.Location{
			Id: &testLocationVar,
			Properties: &ionoscloud.LocationProperties{
				Name:         &testLocationVar,
				Features:     &[]string{testLocationVar},
				ImageAliases: &[]string{testLocationVar},
			},
		},
	}
	locations = v5.Locations{
		Locations: ionoscloud.Locations{
			Id:    &testLocationVar,
			Items: &[]ionoscloud.Location{loc.Location},
		},
	}
	testLocationVar = "test/location"
	testLocationErr = errors.New("location test error occurred")
)

func TestPreLocationId(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLocationId), testLocationVar)
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
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		rm.Location.EXPECT().List().Return(locations, nil, nil)
		err := RunLocationList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunLocationListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		rm.Location.EXPECT().List().Return(locations, nil, testLocationErr)
		err := RunLocationList(cfg)
		assert.Error(t, err)
	})
}

func TestRunLocationGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLocationId), testLocationVar)
		testIds := strings.Split(testLocationVar, "/")
		rm.Location.EXPECT().GetByRegionAndLocationId(testIds[0], testIds[1]).Return(&loc, nil, nil)
		err := RunLocationGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunLocationGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLocationId), testLocationVar)
		testIds := strings.Split(testLocationVar, "/")
		rm.Location.EXPECT().GetByRegionAndLocationId(testIds[0], testIds[1]).Return(&loc, nil, testLocationErr)
		err := RunLocationGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunLocationGetLocationErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLocationId), "us")
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

func TestGetLocationsIds(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() { return }
	w := bufio.NewWriter(&b)
	err := os.Setenv(ionoscloud.IonosUsernameEnvVar, "user")
	assert.NoError(t, err)
	err = os.Setenv(ionoscloud.IonosPasswordEnvVar, "pass")
	assert.NoError(t, err)
	err = os.Setenv(ionoscloud.IonosTokenEnvVar, "tok")
	assert.NoError(t, err)
	viper.Set(config.ArgServerUrl, config.DefaultApiURL)
	getLocationIds(w)
	err = w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`401 Unauthorized`)
	assert.True(t, re.Match(b.Bytes()))
}
