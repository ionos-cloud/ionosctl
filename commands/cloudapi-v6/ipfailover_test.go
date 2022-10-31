package commands

import (
	"bufio"
	"bytes"
	"errors"
	"os"
	"regexp"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/ionos-cloud/ionosctl/pkg/constants"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/services/cloudapi-v6/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	testLanIpFailover = resources.Lan{
		Lan: ionoscloud.Lan{
			Id: &testIpFailoverVar,
			Properties: &ionoscloud.LanProperties{
				IpFailover: &[]ionoscloud.IPFailover{
					{
						Ip:      &testIpFailoverVar,
						NicUuid: &testIpFailoverVar,
					},
				},
			},
		},
	}
	lansIpFailover = resources.Lans{
		Lans: ionoscloud.Lans{
			Id:    &testIpFailoverVar,
			Items: &[]ionoscloud.Lan{l},
		},
	}
	testLanIpFailoverRemove = resources.Lan{
		Lan: ionoscloud.Lan{
			Id: &testIpFailoverVar,
		},
	}
	testLanIpFailoverProperties = resources.Lan{
		Lan: ionoscloud.Lan{
			Id: &testIpFailoverVar,
		},
	}
	testLanIpFailoverGet = resources.Lan{
		Lan: ionoscloud.Lan{
			Id: &testIpFailoverVar,
			Properties: &ionoscloud.LanProperties{
				Name: &testIpFailoverVar,
			},
		},
	}
	testLanPropertiesIpFailover = resources.LanProperties{
		LanProperties: ionoscloud.LanProperties{
			IpFailover: &[]ionoscloud.IPFailover{
				{
					Ip:      &testIpFailoverVar,
					NicUuid: &testIpFailoverVar,
				},
			},
		},
	}
	testIpFailoverVar = "test-ip-failover"
	testIpFailoverErr = errors.New("ip failover error test")
)

func TestIpfailoverCmd(t *testing.T) {
	var err error
	core.RootCmdTest.AddCommand(IpfailoverCmd())
	if ok := IpfailoverCmd().IsAvailableCommand(); !ok {
		err = errors.New("non-available cmd")
	}
	assert.NoError(t, err)
}

func TestPreRunDcLanIds(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLanId), testIpFailoverVar)
		err := PreRunDcLanIds(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunDcLanIdsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		err := PreRunDcLanIds(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunDcLanServerNicIdsIp(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLanId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIp), testIpFailoverVar)
		err := PreRunDcLanServerNicIdsIp(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunDcLanServerNicIdsIpErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		err := PreRunDcLanServerNicIdsIp(cfg)
		assert.Error(t, err)
	})
}

func TestRunIpFailoverList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLanId), testIpFailoverVar)
		rm.CloudApiV6Mocks.Lan.EXPECT().Get(testIpFailoverVar, testIpFailoverVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testLanIpFailover, &testResponse, nil)
		err := RunIpFailoverList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunIpFailoverListPropertiesErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLanId), testIpFailoverVar)
		rm.CloudApiV6Mocks.Lan.EXPECT().Get(testIpFailoverVar, testIpFailoverVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testLanIpFailoverProperties, nil, nil)
		err := RunIpFailoverList(cfg)
		assert.Error(t, err)
	})
}

func TestRunIpFailoverListGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLanId), testIpFailoverVar)
		rm.CloudApiV6Mocks.Lan.EXPECT().Get(testIpFailoverVar, testIpFailoverVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testLanIpFailoverGet, nil, nil)
		err := RunIpFailoverList(cfg)
		assert.Error(t, err)
	})
}

func TestRunIpFailoverListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLanId), testIpFailoverVar)
		rm.CloudApiV6Mocks.Lan.EXPECT().Get(testIpFailoverVar, testIpFailoverVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testLanIpFailover, nil, testIpFailoverErr)
		err := RunIpFailoverList(cfg)
		assert.Error(t, err)
	})
}

func TestRunIpFailoverAdd(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLanId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIp), testIpFailoverVar)
		rm.CloudApiV6Mocks.Lan.EXPECT().Update(testIpFailoverVar, testIpFailoverVar, testLanPropertiesIpFailover, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testLanIpFailover, &testResponse, nil)
		err := RunIpFailoverAdd(cfg)
		assert.NoError(t, err)
	})
}

func TestRunIpFailoverAddErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLanId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIp), testIpFailoverVar)
		rm.CloudApiV6Mocks.Lan.EXPECT().Update(testIpFailoverVar, testIpFailoverVar, testLanPropertiesIpFailover, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testLanIpFailover, nil, testIpFailoverErr)
		err := RunIpFailoverAdd(cfg)
		assert.Error(t, err)
	})
}

func TestRunIpFailoverAddPropertiesErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLanId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIp), testIpFailoverVar)
		rm.CloudApiV6Mocks.Lan.EXPECT().Update(testIpFailoverVar, testIpFailoverVar, testLanPropertiesIpFailover, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testLanIpFailoverProperties, nil, nil)
		err := RunIpFailoverAdd(cfg)
		assert.Error(t, err)
	})
}

func TestRunIpFailoverAddGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLanId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIp), testIpFailoverVar)
		rm.CloudApiV6Mocks.Lan.EXPECT().Update(testIpFailoverVar, testIpFailoverVar, testLanPropertiesIpFailover, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testLanIpFailoverGet, nil, nil)
		err := RunIpFailoverAdd(cfg)
		assert.Error(t, err)
	})
}

func TestRunIpFailoverAddWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLanId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIp), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), true)
		rm.CloudApiV6Mocks.Lan.EXPECT().Update(testIpFailoverVar, testIpFailoverVar, testLanPropertiesIpFailover, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testLanIpFailover, &testResponse, nil)
		rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testRequestIdVar).Return(&testRequestStatus, nil, testRequestErr)
		err := RunIpFailoverAdd(cfg)
		assert.Error(t, err)
	})
}

func TestRunIpFailoverRemove(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLanId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIp), testIpFailoverVar)
		rm.CloudApiV6Mocks.Lan.EXPECT().Get(testIpFailoverVar, testIpFailoverVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testLanIpFailover, &testResponse, nil)
		rm.CloudApiV6Mocks.Lan.EXPECT().Update(testIpFailoverVar, testIpFailoverVar, resources.LanProperties{
			LanProperties: ionoscloud.LanProperties{
				IpFailover: &[]ionoscloud.IPFailover{},
			},
		},
			testQueryParamOther,
		).Return(&testLanIpFailoverRemove, nil, nil)
		err := RunIpFailoverRemove(cfg)
		assert.NoError(t, err)
	})
}

func TestRunIpFailoverRemoveAll(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLanId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.Lan.EXPECT().List(testIpFailoverVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(lansIpFailover, &testResponse, nil)
		rm.CloudApiV6Mocks.Lan.EXPECT().Get(testIpFailoverVar, testIpFailoverVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testLanIpFailover, &testResponse, nil)
		rm.CloudApiV6Mocks.Lan.EXPECT().Update(testIpFailoverVar, testIpFailoverVar, resources.LanProperties{
			LanProperties: ionoscloud.LanProperties{
				IpFailover: &[]ionoscloud.IPFailover{},
			},
		},
			testQueryParamOther,
		).Return(&testLanIpFailoverRemove, nil, nil)
		err := RunIpFailoverRemove(cfg)
		assert.NoError(t, err)
	})
}

func TestRunIpFailoverRemoveResponseErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLanId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIp), testIpFailoverVar)
		rm.CloudApiV6Mocks.Lan.EXPECT().Get(testIpFailoverVar, testIpFailoverVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testLanIpFailover, nil, nil)
		rm.CloudApiV6Mocks.Lan.EXPECT().Update(testIpFailoverVar, testIpFailoverVar, resources.LanProperties{
			LanProperties: ionoscloud.LanProperties{
				IpFailover: &[]ionoscloud.IPFailover{},
			},
		},
			testQueryParamOther,
		).Return(&testLanIpFailoverRemove, &testResponse, testIpFailoverErr)
		err := RunIpFailoverRemove(cfg)
		assert.Error(t, err)
	})
}

func TestRunIpFailoverRemovePropertiesErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLanId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIp), testIpFailoverVar)
		rm.CloudApiV6Mocks.Lan.EXPECT().Get(testIpFailoverVar, testIpFailoverVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testLanIpFailoverProperties, nil, nil)
		err := RunIpFailoverRemove(cfg)
		assert.Error(t, err)
	})
}

func TestRunIpFailoverRemoveGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLanId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIp), testIpFailoverVar)
		rm.CloudApiV6Mocks.Lan.EXPECT().Get(testIpFailoverVar, testIpFailoverVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testLanIpFailoverGet, nil, nil)
		err := RunIpFailoverRemove(cfg)
		assert.Error(t, err)
	})
}

func TestRunIpFailoverRemoveWaitReqErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLanId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIp), testIpFailoverVar)
		rm.CloudApiV6Mocks.Lan.EXPECT().Get(testIpFailoverVar, testIpFailoverVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testLanIpFailover, nil, nil)
		rm.CloudApiV6Mocks.Lan.EXPECT().Update(testIpFailoverVar, testIpFailoverVar, resources.LanProperties{
			LanProperties: ionoscloud.LanProperties{
				IpFailover: &[]ionoscloud.IPFailover{},
			},
		},
			testQueryParamOther,
		).Return(&testLanIpFailoverRemove, &testResponse, nil)
		rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testRequestIdVar).Return(&testRequestStatus, nil, testRequestErr)
		err := RunIpFailoverRemove(cfg)
		assert.Error(t, err)
	})
}

func TestRunIpFailoverRemoveErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLanId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIp), testIpFailoverVar)
		rm.CloudApiV6Mocks.Lan.EXPECT().Get(testIpFailoverVar, testIpFailoverVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testLanIpFailover, nil, testIpFailoverErr)
		err := RunIpFailoverRemove(cfg)
		assert.Error(t, err)
	})
}

func TestRunIpFailoverRemoveUpdateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLanId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIp), testIpFailoverVar)
		rm.CloudApiV6Mocks.Lan.EXPECT().Get(testIpFailoverVar, testIpFailoverVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testLanIpFailover, nil, nil)
		rm.CloudApiV6Mocks.Lan.EXPECT().Update(testIpFailoverVar, testIpFailoverVar, resources.LanProperties{
			LanProperties: ionoscloud.LanProperties{
				IpFailover: &[]ionoscloud.IPFailover{},
			},
		},
			testQueryParamOther,
		).Return(&testLanIpFailoverRemove, nil, testIpFailoverErr)
		err := RunIpFailoverRemove(cfg)
		assert.Error(t, err)
	})
}

func TestRunIpFailoverRemoveAskForConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLanId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIp), testIpFailoverVar)
		cfg.Stdin = bytes.NewReader([]byte("YES\n"))
		rm.CloudApiV6Mocks.Lan.EXPECT().Get(testIpFailoverVar, testIpFailoverVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testLanIpFailover, nil, nil)
		rm.CloudApiV6Mocks.Lan.EXPECT().Update(testIpFailoverVar, testIpFailoverVar, resources.LanProperties{
			LanProperties: ionoscloud.LanProperties{
				IpFailover: &[]ionoscloud.IPFailover{},
			},
		},
			testQueryParamOther,
		).Return(&testLanIpFailoverRemove, nil, nil)
		err := RunIpFailoverRemove(cfg)
		assert.NoError(t, err)
	})
}

func TestRunIpFailoverRemoveAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLanId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIp), testIpFailoverVar)
		cfg.Stdin = os.Stdin
		err := RunIpFailoverRemove(cfg)
		assert.Error(t, err)
	})
}

func TestGetIpFailoverCols(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	viper.Reset()
	viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
	viper.Set(constants.ArgQuiet, false)
	viper.Set(core.GetGlobalFlagName("ipfailover", constants.ArgCols), []string{"NicId"})
	getIpFailoverCols(core.GetGlobalFlagName("ipfailover", constants.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
}

func TestGetIpFailoverColsErr(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	viper.Reset()
	viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
	viper.Set(constants.ArgQuiet, false)
	viper.Set(core.GetGlobalFlagName("ipfailover", constants.ArgCols), []string{"Unknown"})
	getIpFailoverCols(core.GetGlobalFlagName("ipfailover", constants.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`unknown column Unknown`)
	assert.True(t, re.Match(b.Bytes()))
}
