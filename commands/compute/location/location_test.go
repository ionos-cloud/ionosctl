package location

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/testutil"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"

	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
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
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
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
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		err := PreRunLocationId(cfg)
		assert.Error(t, err)
	})
}
func TestRunLocationList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		rm.CloudApiV6Mocks.Location.EXPECT().List().Return(locations, &testutil.TestResponse, nil)
		err := RunLocationList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunLocationListQueryParams(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		cfg.Command.Command.Flags().Set(constants.FlagFilters, fmt.Sprintf("%s=%s", testutil.TestQueryParamVar, testutil.TestQueryParamVar))
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagOrderBy), testutil.TestQueryParamVar)
		rm.CloudApiV6Mocks.Location.EXPECT().List().Return(resources.Locations{}, &testutil.TestResponse, nil)
		err := RunLocationList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunLocationListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
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
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLocationId), testLocationVar)
		testIds := strings.Split(testLocationVar, "/")
		rm.CloudApiV6Mocks.Location.EXPECT().GetByRegionAndLocationId(testIds[0], testIds[1]).Return(&loc, &testutil.TestResponse, nil)
		err := RunLocationGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunLocationGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
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
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLocationId), "us")
		err := RunLocationGet(cfg)
		assert.Error(t, err)
	})
}
