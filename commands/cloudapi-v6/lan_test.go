package commands

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/services/cloudapi-v6/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
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
		viper.Set(config.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testLanVar)
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
		viper.Set(config.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testLanVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFilters), []string{fmt.Sprintf("createdBy=%s", testQueryParamVar)})
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
		viper.Set(config.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testLanVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFilters), []string{fmt.Sprintf("%s=%s", testQueryParamVar, testQueryParamVar)})
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
		viper.Set(config.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testLanVar)
		rm.CloudApiV6Mocks.Lan.EXPECT().List(testLanVar, resources.ListQueryParams{}).Return(ls, &testResponse, nil)
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
		viper.Set(config.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testLanVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFilters), []string{fmt.Sprintf("%s=%s", testQueryParamVar, testQueryParamVar)})
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgOrderBy), testQueryParamVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgMaxResults), testMaxResultsVar)
		rm.CloudApiV6Mocks.Lan.EXPECT().List(testLanVar, testListQueryParam).Return(resources.Lans{}, &testResponse, nil)
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
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testLanVar)
		rm.CloudApiV6Mocks.Lan.EXPECT().List(testLanVar, resources.ListQueryParams{}).Return(ls, nil, testLanErr)
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
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testLanVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLanId), testLanVar)
		rm.CloudApiV6Mocks.Lan.EXPECT().Get(testLanVar, testLanVar, resources.QueryParams{}).Return(&resources.Lan{Lan: l}, &testResponse, nil)
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
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testLanVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLanId), testLanVar)
		rm.CloudApiV6Mocks.Lan.EXPECT().Get(testLanVar, testLanVar, resources.QueryParams{}).Return(&resources.Lan{Lan: l}, nil, testLanErr)
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
		viper.Set(config.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testLanVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testLanVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPccId), testLanVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPublic), publicLan)
		rm.CloudApiV6Mocks.Lan.EXPECT().Create(testLanVar, resources.LanPost{LanPost: lanPostTest}, resources.QueryParams{}).Return(&resources.LanPost{LanPost: lp}, &testResponse, nil)
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
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testLanVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testLanVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPublic), publicLan)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPccId), testLanVar)
		rm.CloudApiV6Mocks.Lan.EXPECT().Create(testLanVar, resources.LanPost{LanPost: lanPostTest}, resources.QueryParams{}).Return(&resources.LanPost{LanPost: lp}, nil, testLanErr)
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
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testLanVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testLanVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testLanVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPublic), publicLan)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPccId), testLanVar)
		rm.CloudApiV6Mocks.Lan.EXPECT().Create(testLanVar, resources.LanPost{LanPost: lanPostTest}, resources.QueryParams{}).Return(&resources.LanPost{LanPost: lp}, &testResponse, nil)
		rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testRequestIdVar, resources.QueryParams{}).Return(&testRequestStatus, nil, testRequestErr)
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
		viper.Set(config.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testLanVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLanId), testLanVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testLanNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPublic), publicNewLan)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPccId), testLanNewVar)
		rm.CloudApiV6Mocks.Lan.EXPECT().Update(testLanVar, testLanVar, lanProperties, resources.QueryParams{}).Return(&lanNew, &testResponse, nil)
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
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testLanVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testLanVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLanId), testLanVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testLanNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPublic), publicNewLan)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPccId), testLanNewVar)
		rm.CloudApiV6Mocks.Lan.EXPECT().Update(testLanVar, testLanVar, lanProperties, resources.QueryParams{}).Return(&lanNew, nil, testLanErr)
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
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testLanVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testLanVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLanId), testLanVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testLanNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPublic), publicNewLan)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPccId), testLanNewVar)
		rm.CloudApiV6Mocks.Lan.EXPECT().Update(testLanVar, testLanVar, lanProperties, resources.QueryParams{}).Return(&lanNew, &testResponse, testLanErr)
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
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testLanVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testLanVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLanId), testLanVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testLanNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPublic), publicNewLan)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPccId), testLanNewVar)
		rm.CloudApiV6Mocks.Lan.EXPECT().Update(testLanVar, testLanVar, lanProperties, resources.QueryParams{}).Return(&lanNew, &testResponse, nil)
		rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testRequestIdVar, resources.QueryParams{}).Return(&testRequestStatus, nil, testRequestErr)
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
		viper.Set(config.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testLanVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLanId), testLanVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.Lan.EXPECT().Delete(testLanVar, testLanVar, resources.QueryParams{}).Return(&testResponse, nil)
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
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testLanVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.Lan.EXPECT().List(testLanVar, resources.ListQueryParams{}).Return(lansList, &testResponse, nil)
		rm.CloudApiV6Mocks.Lan.EXPECT().Delete(testLanVar, testLanVar, resources.QueryParams{}).Return(&testResponse, nil)
		rm.CloudApiV6Mocks.Lan.EXPECT().Delete(testLanVar, testLanVar, resources.QueryParams{}).Return(&testResponse, nil)
		err := RunLanDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunLanDeleteAllListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testLanVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.Lan.EXPECT().List(testLanVar, resources.ListQueryParams{}).Return(lansList, nil, testLanErr)
		err := RunLanDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunLanDeleteAllItemsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testLanVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.Lan.EXPECT().List(testLanVar, resources.ListQueryParams{}).Return(resources.Lans{}, &testResponse, nil)
		err := RunLanDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunLanDeleteAllLenErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testLanVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.Lan.EXPECT().List(testLanVar, resources.ListQueryParams{}).Return(
			resources.Lans{Lans: ionoscloud.Lans{Items: &[]ionoscloud.Lan{}}}, &testResponse, nil)
		err := RunLanDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunLanDeleteAllErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testLanVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.Lan.EXPECT().List(testLanVar, resources.ListQueryParams{}).Return(lansList, &testResponse, nil)
		rm.CloudApiV6Mocks.Lan.EXPECT().Delete(testLanVar, testLanVar, resources.QueryParams{}).Return(&testResponse, testLanErr)
		rm.CloudApiV6Mocks.Lan.EXPECT().Delete(testLanVar, testLanVar, resources.QueryParams{}).Return(&testResponse, nil)
		err := RunLanDelete(cfg)
		assert.Error(t, err)
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
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testLanVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLanId), testLanVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.Lan.EXPECT().Delete(testLanVar, testLanVar, resources.QueryParams{}).Return(nil, testLanErr)
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
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testLanVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLanId), testLanVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		rm.CloudApiV6Mocks.Lan.EXPECT().Delete(testLanVar, testLanVar, resources.QueryParams{}).Return(&testResponse, nil)
		rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testRequestIdVar, resources.QueryParams{}).Return(&testRequestStatus, nil, testRequestErr)
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
		cfg.Stdin = bytes.NewReader([]byte("YES\n"))
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testLanVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLanId), testLanVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.Lan.EXPECT().Delete(testLanVar, testLanVar, resources.QueryParams{}).Return(nil, nil)
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
		cfg.Stdin = os.Stdin
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testLanVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLanId), testLanVar)
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
