package commands

import (
	"bufio"
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
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
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLocationId), testLocationCpuVar)
		testIds := strings.Split(testLocationCpuVar, "/")
		rm.CloudApiV6Mocks.Location.EXPECT().GetByRegionAndLocationId(testIds[0], testIds[1], gomock.AssignableToTypeOf(testQueryParamOther)).Return(&loc, &testResponse, nil)
		err := RunLocationCpuList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunLocationCpuListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLocationId), testLocationCpuVar)
		testIds := strings.Split(testLocationCpuVar, "/")
		rm.CloudApiV6Mocks.Location.EXPECT().GetByRegionAndLocationId(testIds[0], testIds[1], gomock.AssignableToTypeOf(testQueryParamOther)).Return(&loc, nil, testLocationCpuErr)
		err := RunLocationCpuList(cfg)
		assert.Error(t, err)
	})
}
