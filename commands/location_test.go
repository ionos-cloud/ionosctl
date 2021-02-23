package commands

import (
	"bufio"
	"bytes"
	"errors"
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
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		rm.Location.EXPECT().List().Return(locations, nil, nil)
		viper.Set(config.ArgQuiet, false)
		err := RunLocationList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunLocationList_Err(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		rm.Location.EXPECT().List().Return(locations, nil, testLocationErr)
		viper.Set(config.ArgQuiet, false)
		err := RunLocationList(cfg)
		assert.Error(t, err)
	})
}

func TestGetLocationsCols(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}

	w := bufio.NewWriter(&b)
	viper.Set(builder.GetGlobalFlagName("location", config.ArgCols), []string{"Name"})
	getLocationCols(builder.GetGlobalFlagName("location", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
}

func TestGetLocationsCols_Err(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}

	w := bufio.NewWriter(&b)
	viper.Set(builder.GetGlobalFlagName("location", config.ArgCols), []string{"Unknown"})
	getLocationCols(builder.GetGlobalFlagName("location", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`unknown column Unknown`)
	assert.True(t, re.Match(b.Bytes()))
}

func TestGetLocationsIds(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}

	w := bufio.NewWriter(&b)
	viper.Set(config.ArgConfig, "../pkg/testdata/config.json")
	getLocationIds(w)
	err := w.Flush()
	assert.NoError(t, err)

	re := regexp.MustCompile(`401 Unauthorized`)
	assert.True(t, re.Match(b.Bytes()))
}
