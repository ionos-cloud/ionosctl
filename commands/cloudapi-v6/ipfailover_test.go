package commands

import (
	"bufio"
	"bytes"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"

	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
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
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagLanId), testIpFailoverVar)
		err := PreRunDcLanIds(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunDcLanIdsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		err := PreRunDcLanIds(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunDcLanServerNicIdsIp(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagServerId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagNicId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagLanId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagIp), testIpFailoverVar)
		err := PreRunDcLanServerNicIdsIp(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunDcLanServerNicIdsIpErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		err := PreRunDcLanServerNicIdsIp(cfg)
		assert.Error(t, err)
	})
}

func TestRunIpFailoverList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagLanId), testIpFailoverVar)
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
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagLanId), testIpFailoverVar)
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
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagLanId), testIpFailoverVar)
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
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagLanId), testIpFailoverVar)
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
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagLanId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagServerId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagNicId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagIp), testIpFailoverVar)
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
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagLanId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagServerId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagNicId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagIp), testIpFailoverVar)
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
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagLanId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagServerId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagNicId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagIp), testIpFailoverVar)
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
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagLanId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagServerId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagNicId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagIp), testIpFailoverVar)
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
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagLanId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagServerId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagNicId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagIp), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagWaitForRequest), true)
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
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagForce, true)
		viper.Set(constants.FlagVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagWaitForRequest), false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagLanId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagServerId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagNicId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagIp), testIpFailoverVar)
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
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagVerbose, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagForce, true)
		viper.Set(constants.FlagServerUrl, constants.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagWaitForRequest), false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagLanId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagServerId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagNicId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagAll), true)
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
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagForce, true)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagWaitForRequest), false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagLanId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagServerId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagNicId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagIp), testIpFailoverVar)
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
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagForce, true)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagWaitForRequest), false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagLanId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagServerId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagNicId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagIp), testIpFailoverVar)
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
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagForce, true)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagWaitForRequest), false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagLanId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagServerId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagNicId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagIp), testIpFailoverVar)
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
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagForce, true)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagWaitForRequest), true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagLanId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagServerId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagNicId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagIp), testIpFailoverVar)
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
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagWaitForRequest), false)
		viper.Set(constants.FlagForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagLanId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagServerId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagNicId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagIp), testIpFailoverVar)
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
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagWaitForRequest), false)
		viper.Set(constants.FlagForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagLanId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagServerId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagNicId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagIp), testIpFailoverVar)
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
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagForce, false)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagWaitForRequest), false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagLanId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagServerId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagNicId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagIp), testIpFailoverVar)
		cfg.Command.Command.SetIn(bytes.NewReader([]byte("YES\n")))
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
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagForce, false)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagWaitForRequest), false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagLanId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagServerId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagNicId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagIp), testIpFailoverVar)
		cfg.Command.Command.SetIn(bytes.NewReader([]byte("\n")))
		err := RunIpFailoverRemove(cfg)
		assert.Error(t, err)
	})
}
