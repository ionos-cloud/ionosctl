package cloudapi_v5

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/ionos-cloud/ionosctl/internal/config"
	"github.com/ionos-cloud/ionosctl/internal/core"
	"github.com/ionos-cloud/ionosctl/internal/utils/clierror"
	cloudapiv5 "github.com/ionos-cloud/ionosctl/services/cloudapi-v5"
	"github.com/ionos-cloud/ionosctl/services/cloudapi-v5/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	publicLan    = true
	publicNewLan = false
	lanPostTest  = ionoscloud.LanPost{
		Properties: &ionoscloud.LanPropertiesPost{
			Name:       &testLanVar,
			IpFailover: nil,
			Pcc:        &testLanVar,
			Public:     &publicLan,
		},
	}
	lp = ionoscloud.LanPost{
		Id:         &testLanVar,
		Properties: lanPostTest.Properties,
		Metadata:   &ionoscloud.DatacenterElementMetadata{State: &testStateVar},
	}
	l = ionoscloud.Lan{
		Id: &testLanVar,
		Properties: &ionoscloud.LanProperties{
			Name: &testLanVar,
			Pcc:  &testLanVar,
		},
	}
	lanProperties = resources.LanProperties{
		LanProperties: ionoscloud.LanProperties{
			Name:   &testLanNewVar,
			Pcc:    &testLanNewVar,
			Public: &publicNewLan,
		},
	}
	lanNew = resources.Lan{
		Lan: ionoscloud.Lan{
			Id: &testLanVar,
			Properties: &ionoscloud.LanProperties{
				Name:       lanProperties.LanProperties.Name,
				Public:     lanProperties.LanProperties.Public,
				IpFailover: nil,
				Pcc:        &testLanNewVar,
			},
		},
	}
	ls = resources.Lans{
		Lans: ionoscloud.Lans{
			Id:    &testLanVar,
			Items: &[]ionoscloud.Lan{l},
		},
	}
	lansList = resources.Lans{
		Lans: ionoscloud.Lans{
			Id: &testLanVar,
			Items: &[]ionoscloud.Lan{
				l,
				l,
			},
		},
	}
	testLanVar    = "test-lan"
	testLanNewVar = "test-new-lan"
	testLanErr    = errors.New("lan test: error occurred")
)

func TestLanCmd(t *testing.T) {
	var err error
	core.RootCmdTest.AddCommand(LanCmd())
	if ok := LanCmd().IsAvailableCommand(); !ok {
		err = errors.New("non-available cmd")
	}
	assert.NoError(t, err)
}

func TestPreRunLansList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgDataCenterId), testLanVar)
		err := PreRunLansList(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunLansListFilters(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgDataCenterId), testLanVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgFilters), []string{fmt.Sprintf("createdBy=%s", testQueryParamVar)})
		err := PreRunLansList(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunLansListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgDataCenterId), testLanVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgFilters), []string{fmt.Sprintf("%s=%s", testQueryParamVar, testQueryParamVar)})
		err := PreRunLansList(cfg)
		assert.Error(t, err)
	})
}

func TestRunLanList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, true)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgDataCenterId), testLanVar)
		rm.CloudApiV5Mocks.Lan.EXPECT().List(testLanVar, resources.ListQueryParams{}).Return(ls, &testResponse, nil)
		err := RunLanList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunLanListQueryParams(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgDataCenterId), testLanVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgFilters), []string{fmt.Sprintf("%s=%s", testQueryParamVar, testQueryParamVar)})
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgOrderBy), testQueryParamVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgMaxResults), testMaxResultsVar)
		rm.CloudApiV5Mocks.Lan.EXPECT().List(testLanVar, testListQueryParam).Return(resources.Lans{}, &testResponse, nil)
		err := RunLanList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunLanListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgDataCenterId), testLanVar)
		rm.CloudApiV5Mocks.Lan.EXPECT().List(testLanVar, resources.ListQueryParams{}).Return(ls, nil, testLanErr)
		err := RunLanList(cfg)
		assert.Error(t, err)
	})
}

func TestRunLanGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, true)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgDataCenterId), testLanVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgLanId), testLanVar)
		rm.CloudApiV5Mocks.Lan.EXPECT().Get(testLanVar, testLanVar).Return(&resources.Lan{Lan: l}, &testResponse, nil)
		err := RunLanGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunLanGet_Err(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgDataCenterId), testLanVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgLanId), testLanVar)
		rm.CloudApiV5Mocks.Lan.EXPECT().Get(testLanVar, testLanVar).Return(&resources.Lan{Lan: l}, nil, testLanErr)
		err := RunLanGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunLanCreate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgDataCenterId), testLanVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgName), testLanVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgPccId), testLanVar)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgPublic), publicLan)
		rm.CloudApiV5Mocks.Lan.EXPECT().Create(testLanVar, resources.LanPost{LanPost: lanPostTest}).Return(&resources.LanPost{LanPost: lp}, &testResponse, nil)
		err := RunLanCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunLanCreateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgDataCenterId), testLanVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgName), testLanVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgPublic), publicLan)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgPccId), testLanVar)
		rm.CloudApiV5Mocks.Lan.EXPECT().Create(testLanVar, resources.LanPost{LanPost: lanPostTest}).Return(&resources.LanPost{LanPost: lp}, nil, testLanErr)
		err := RunLanCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunLanCreateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgDataCenterId), testLanVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgServerId), testLanVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgName), testLanVar)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgPublic), publicLan)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgPccId), testLanVar)
		rm.CloudApiV5Mocks.Lan.EXPECT().Create(testLanVar, resources.LanPost{LanPost: lanPostTest}).Return(&resources.LanPost{LanPost: lp}, &testResponse, nil)
		rm.CloudApiV5Mocks.Request.EXPECT().GetStatus(testRequestIdVar).Return(&testRequestStatus, nil, testRequestErr)
		err := RunLanCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunLanUpdate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, true)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgDataCenterId), testLanVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgLanId), testLanVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgName), testLanNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgPublic), publicNewLan)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgPccId), testLanNewVar)
		rm.CloudApiV5Mocks.Lan.EXPECT().Update(testLanVar, testLanVar, lanProperties).Return(&lanNew, &testResponse, nil)
		err := RunLanUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunLanUpdateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgDataCenterId), testLanVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgServerId), testLanVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgLanId), testLanVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgName), testLanNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgPublic), publicNewLan)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgPccId), testLanNewVar)
		rm.CloudApiV5Mocks.Lan.EXPECT().Update(testLanVar, testLanVar, lanProperties).Return(&lanNew, nil, testLanErr)
		err := RunLanUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunLanUpdateResponseErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgDataCenterId), testLanVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgServerId), testLanVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgLanId), testLanVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgName), testLanNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgPublic), publicNewLan)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgPccId), testLanNewVar)
		rm.CloudApiV5Mocks.Lan.EXPECT().Update(testLanVar, testLanVar, lanProperties).Return(&lanNew, &testResponseErr, nil)
		err := RunLanUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunLanUpdateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgDataCenterId), testLanVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgServerId), testLanVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgLanId), testLanVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgName), testLanNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgPublic), publicNewLan)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgPccId), testLanNewVar)
		rm.CloudApiV5Mocks.Lan.EXPECT().Update(testLanVar, testLanVar, lanProperties).Return(&lanNew, &testResponse, nil)
		rm.CloudApiV5Mocks.Request.EXPECT().GetStatus(testRequestIdVar).Return(&testRequestStatus, nil, testRequestErr)
		err := RunLanUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunLanDelete(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(config.ArgVerbose, true)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgDataCenterId), testLanVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgLanId), testLanVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.CloudApiV5Mocks.Lan.EXPECT().Delete(testLanVar, testLanVar).Return(&testResponse, nil)
		err := RunLanDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunLanDeleteAll(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(config.ArgVerbose, true)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgDataCenterId), testLanVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgAll), true)
		rm.CloudApiV5Mocks.Lan.EXPECT().List(testLanVar, resources.ListQueryParams{}).Return(lansList, &testResponse, nil)
		rm.CloudApiV5Mocks.Lan.EXPECT().Delete(testLanVar, testLanVar).Return(&testResponse, nil)
		rm.CloudApiV5Mocks.Lan.EXPECT().Delete(testLanVar, testLanVar).Return(&testResponse, nil)
		err := RunLanDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunLanDeleteErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgDataCenterId), testLanVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgLanId), testLanVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.CloudApiV5Mocks.Lan.EXPECT().Delete(testLanVar, testLanVar).Return(nil, testLanErr)
		err := RunLanDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunLanDeleteWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgDataCenterId), testLanVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgLanId), testLanVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		rm.CloudApiV5Mocks.Lan.EXPECT().Delete(testLanVar, testLanVar).Return(&testResponse, nil)
		rm.CloudApiV5Mocks.Request.EXPECT().GetStatus(testRequestIdVar).Return(&testRequestStatus, nil, testRequestErr)
		err := RunLanDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunLanDeleteAskForConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, false)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		cfg.Stdin = bytes.NewReader([]byte("YES\n"))
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgDataCenterId), testLanVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgLanId), testLanVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.CloudApiV5Mocks.Lan.EXPECT().Delete(testLanVar, testLanVar).Return(nil, nil)
		err := RunLanDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunLanDeleteAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, false)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		cfg.Stdin = os.Stdin
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgDataCenterId), testLanVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgLanId), testLanVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		err := RunLanDelete(cfg)
		assert.Error(t, err)
	})
}

func TestGetLansCols(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() { return }
	w := bufio.NewWriter(&b)
	viper.Set(core.GetGlobalFlagName("lan", config.ArgCols), []string{"Name"})
	getLansCols(core.GetGlobalFlagName("lan", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
}

func TestGetLansColsErr(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() { return }
	w := bufio.NewWriter(&b)
	viper.Set(core.GetGlobalFlagName("lan", config.ArgCols), []string{"Unknown"})
	getLansCols(core.GetGlobalFlagName("lan", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`unknown column Unknown`)
	assert.True(t, re.Match(b.Bytes()))
}
