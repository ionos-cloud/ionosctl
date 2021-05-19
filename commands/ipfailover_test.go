package commands

import (
	"bufio"
	"bytes"
	"errors"
	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
	"regexp"
	"testing"

	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
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
	testLanIpFailoverRemove = resources.Lan{
		Lan: ionoscloud.Lan{
			Id: &testIpFailoverVar,
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

func TestRunIpFailoverListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testIpFailoverVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLanId), testIpFailoverVar)
		rm.Lan.EXPECT().Get(testIpFailoverVar, testIpFailoverVar).Return(&testLanIpFailover, nil, testIpFailoverErr)
		err := RunIpFailoverList(cfg)
		assert.Error(t, err)
	})
}

func TestRunIpFailoverAddErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
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

func TestRunIpFailoverAddWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
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

func TestGetIpFailoverCols(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	viper.Set(core.GetGlobalFlagName("ipfailover", config.ArgCols), []string{"Ip"})
	getIpFailoverCols(core.GetGlobalFlagName("ipfailover", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
}

func TestGetIpFailoverColsErr(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	viper.Set(core.GetGlobalFlagName("ipfailover", config.ArgCols), []string{"Unknown"})
	getIpFailoverCols(core.GetGlobalFlagName("ipfailover", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`unknown column Unknown`)
	assert.True(t, re.Match(b.Bytes()))
}
