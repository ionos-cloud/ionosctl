package commands

import (
	"bufio"
	"bytes"
	"errors"
	"regexp"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/utils/clierror"
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

func TestGetCpusCols(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() { return }
	w := bufio.NewWriter(&b)
	viper.Set(core.GetFlagName("cpu", constants.ArgCols), []string{"Vendor"})
	getCpuCols(core.GetFlagName("cpu", constants.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
}

func TestGetCpusColsErr(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() { return }
	w := bufio.NewWriter(&b)
	viper.Set(core.GetFlagName("cpu", constants.ArgCols), []string{"Unknown"})
	getCpuCols(core.GetFlagName("cpu", constants.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`unknown column Unknown`)
	assert.True(t, re.Match(b.Bytes()))
}
