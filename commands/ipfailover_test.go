package commands

import (
	"bufio"
	"bytes"
	"errors"
	"os"
	"regexp"
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
	testLanIpFailover = v5.Lan{
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
	testLanIpFailoverRemove = v5.Lan{
		Lan: ionoscloud.Lan{
			Id: &testIpFailoverVar,
		},
	}
	testLanIpFailoverProperties = v5.Lan{
		Lan: ionoscloud.Lan{
			Id: &testIpFailoverVar,
		},
	}
	testLanIpFailoverGet = v5.Lan{
		Lan: ionoscloud.Lan{
			Id: &testIpFailoverVar,
			Properties: &ionoscloud.LanProperties{
				Name: &testIpFailoverVar,
			},
		},
	}
	testLanPropertiesIpFailover = v5.LanProperties{
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

func TestPreRunDcLanIds(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLanId), testIpFailoverVar)
		err := PreRunDcLanIds(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunDcLanIdsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		err := PreRunDcLanIds(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunDcLanServerNicIdsIp(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgServerId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNicId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLanId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgIp), testIpFailoverVar)
		err := PreRunDcLanServerNicIdsIp(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunDcLanServerNicIdsIpErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		err := PreRunDcLanServerNicIdsIp(cfg)
		assert.Error(t, err)
	})
}

func TestRunIpFailoverList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLanId), testIpFailoverVar)
		rm.Lan.EXPECT().Get(testIpFailoverVar, testIpFailoverVar).Return(&testLanIpFailover, nil, nil)
		err := RunIpFailoverList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunIpFailoverListPropertiesErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLanId), testIpFailoverVar)
		rm.Lan.EXPECT().Get(testIpFailoverVar, testIpFailoverVar).Return(&testLanIpFailoverProperties, nil, nil)
		err := RunIpFailoverList(cfg)
		assert.Error(t, err)
	})
}

func TestRunIpFailoverListGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLanId), testIpFailoverVar)
		rm.Lan.EXPECT().Get(testIpFailoverVar, testIpFailoverVar).Return(&testLanIpFailoverGet, nil, nil)
		err := RunIpFailoverList(cfg)
		assert.Error(t, err)
	})
}

func TestRunIpFailoverListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLanId), testIpFailoverVar)
		rm.Lan.EXPECT().Get(testIpFailoverVar, testIpFailoverVar).Return(&testLanIpFailover, nil, testIpFailoverErr)
		err := RunIpFailoverList(cfg)
		assert.Error(t, err)
	})
}

func TestRunIpFailoverAdd(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLanId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgServerId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNicId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgIp), testIpFailoverVar)
		rm.Lan.EXPECT().Update(testIpFailoverVar, testIpFailoverVar, testLanPropertiesIpFailover).Return(&testLanIpFailover, nil, nil)
		err := RunIpFailoverAdd(cfg)
		assert.NoError(t, err)
	})
}

func TestRunIpFailoverAddErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLanId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgServerId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNicId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgIp), testIpFailoverVar)
		rm.Lan.EXPECT().Update(testIpFailoverVar, testIpFailoverVar, testLanPropertiesIpFailover).Return(&testLanIpFailover, nil, testIpFailoverErr)
		err := RunIpFailoverAdd(cfg)
		assert.Error(t, err)
	})
}

func TestRunIpFailoverAddPropertiesErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLanId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgServerId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNicId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgIp), testIpFailoverVar)
		rm.Lan.EXPECT().Update(testIpFailoverVar, testIpFailoverVar, testLanPropertiesIpFailover).Return(&testLanIpFailoverProperties, nil, nil)
		err := RunIpFailoverAdd(cfg)
		assert.Error(t, err)
	})
}

func TestRunIpFailoverAddGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLanId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgServerId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNicId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgIp), testIpFailoverVar)
		rm.Lan.EXPECT().Update(testIpFailoverVar, testIpFailoverVar, testLanPropertiesIpFailover).Return(&testLanIpFailoverGet, nil, nil)
		err := RunIpFailoverAdd(cfg)
		assert.Error(t, err)
	})
}

func TestRunIpFailoverAddWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLanId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgServerId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNicId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgIp), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		rm.Lan.EXPECT().Update(testIpFailoverVar, testIpFailoverVar, testLanPropertiesIpFailover).Return(&testLanIpFailover, nil, nil)
		err := RunIpFailoverAdd(cfg)
		assert.Error(t, err)
	})
}

func TestRunIpFailoverRemove(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLanId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgServerId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNicId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgIp), testIpFailoverVar)
		rm.Lan.EXPECT().Get(testIpFailoverVar, testIpFailoverVar).Return(&testLanIpFailover, nil, nil)
		rm.Lan.EXPECT().Update(testIpFailoverVar, testIpFailoverVar, v5.LanProperties{
			LanProperties: ionoscloud.LanProperties{
				IpFailover: &[]ionoscloud.IPFailover{},
			},
		}).Return(&testLanIpFailoverRemove, nil, nil)
		err := RunIpFailoverRemove(cfg)
		assert.NoError(t, err)
	})
}

func TestRunIpFailoverRemoveResponseErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLanId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgServerId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNicId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgIp), testIpFailoverVar)
		rm.Lan.EXPECT().Get(testIpFailoverVar, testIpFailoverVar).Return(&testLanIpFailover, nil, nil)
		rm.Lan.EXPECT().Update(testIpFailoverVar, testIpFailoverVar, v5.LanProperties{
			LanProperties: ionoscloud.LanProperties{
				IpFailover: &[]ionoscloud.IPFailover{},
			},
		}).Return(&testLanIpFailoverRemove, &testResponseErr, nil)
		err := RunIpFailoverRemove(cfg)
		assert.Error(t, err)
	})
}

func TestRunIpFailoverRemovePropertiesErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLanId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgServerId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNicId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgIp), testIpFailoverVar)
		rm.Lan.EXPECT().Get(testIpFailoverVar, testIpFailoverVar).Return(&testLanIpFailoverProperties, nil, nil)
		err := RunIpFailoverRemove(cfg)
		assert.Error(t, err)
	})
}

func TestRunIpFailoverRemoveGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLanId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgServerId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNicId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgIp), testIpFailoverVar)
		rm.Lan.EXPECT().Get(testIpFailoverVar, testIpFailoverVar).Return(&testLanIpFailoverGet, nil, nil)
		err := RunIpFailoverRemove(cfg)
		assert.Error(t, err)
	})
}

func TestRunIpFailoverRemoveWaitReqErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLanId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgServerId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNicId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgIp), testIpFailoverVar)
		rm.Lan.EXPECT().Get(testIpFailoverVar, testIpFailoverVar).Return(&testLanIpFailover, nil, nil)
		rm.Lan.EXPECT().Update(testIpFailoverVar, testIpFailoverVar, v5.LanProperties{
			LanProperties: ionoscloud.LanProperties{
				IpFailover: &[]ionoscloud.IPFailover{},
			},
		}).Return(&testLanIpFailoverRemove, nil, nil)
		err := RunIpFailoverRemove(cfg)
		assert.Error(t, err)
	})
}

func TestRunIpFailoverRemoveErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLanId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgServerId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNicId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgIp), testIpFailoverVar)
		rm.Lan.EXPECT().Get(testIpFailoverVar, testIpFailoverVar).Return(&testLanIpFailover, nil, testIpFailoverErr)
		err := RunIpFailoverRemove(cfg)
		assert.Error(t, err)
	})
}

func TestRunIpFailoverRemoveUpdateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLanId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgServerId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNicId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgIp), testIpFailoverVar)
		rm.Lan.EXPECT().Get(testIpFailoverVar, testIpFailoverVar).Return(&testLanIpFailover, nil, nil)
		rm.Lan.EXPECT().Update(testIpFailoverVar, testIpFailoverVar, v5.LanProperties{
			LanProperties: ionoscloud.LanProperties{
				IpFailover: &[]ionoscloud.IPFailover{},
			},
		}).Return(&testLanIpFailoverRemove, nil, testIpFailoverErr)
		err := RunIpFailoverRemove(cfg)
		assert.Error(t, err)
	})
}

func TestRunIpFailoverRemoveAskForConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, false)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLanId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgServerId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNicId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgIp), testIpFailoverVar)
		cfg.Stdin = bytes.NewReader([]byte("YES\n"))
		rm.Lan.EXPECT().Get(testIpFailoverVar, testIpFailoverVar).Return(&testLanIpFailover, nil, nil)
		rm.Lan.EXPECT().Update(testIpFailoverVar, testIpFailoverVar, v5.LanProperties{
			LanProperties: ionoscloud.LanProperties{
				IpFailover: &[]ionoscloud.IPFailover{},
			},
		}).Return(&testLanIpFailoverRemove, nil, nil)
		err := RunIpFailoverRemove(cfg)
		assert.NoError(t, err)
	})
}

func TestRunIpFailoverRemoveAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLanId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgServerId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgNicId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgIp), testIpFailoverVar)
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
	viper.Set(config.ArgOutput, config.DefaultOutputFormat)
	viper.Set(config.ArgQuiet, false)
	viper.Set(core.GetGlobalFlagName("ipfailover", config.ArgCols), []string{"NicId"})
	getIpFailoverCols(core.GetGlobalFlagName("ipfailover", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
}

func TestGetIpFailoverColsErr(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	viper.Reset()
	viper.Set(config.ArgOutput, config.DefaultOutputFormat)
	viper.Set(config.ArgQuiet, false)
	viper.Set(core.GetGlobalFlagName("ipfailover", config.ArgCols), []string{"Unknown"})
	getIpFailoverCols(core.GetGlobalFlagName("ipfailover", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`unknown column Unknown`)
	assert.True(t, re.Match(b.Bytes()))
}
