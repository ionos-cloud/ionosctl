package commands

import (
	"bufio"
	"bytes"
	"errors"
	"regexp"
	"testing"

	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/resources"
	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	loc = resources.Location{
		Location: ionoscloud.Location{
			Id: &testLocationVar,
			Properties: &ionoscloud.LocationProperties{
				Name:     &testLocationVar,
				Features: &[]string{testLocationVar},
			},
		},
	}
	locations = resources.Locations{
		Locations: ionoscloud.Locations{
			Id:    &testLocationVar,
			Items: &[]ionoscloud.Location{loc.Location},
		},
	}
	testLocationVar = "test-location"
	testLocationErr = errors.New("location test error occurred")
)

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
	viper.Set(config.ArgConfig, "../pkg/testdata/config.json")
	getLocationIds(w)
	err := w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`401 Unauthorized`)
	assert.True(t, re.Match(b.Bytes()))
}
