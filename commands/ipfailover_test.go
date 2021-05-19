package commands

import (
	"bufio"
	"bytes"
	"errors"
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
	testIpFailoverVar     = "test-ip-failover"
	testIpFailoverBoolVar = false
	testIpFailoverErr     = errors.New("ip failover error test")
)

//func TestRunIpFailoverList(t *testing.T) {
//	var b bytes.Buffer
//	w := bufio.NewWriter(&b)
//	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
//		viper.Reset()
//		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
//		viper.Set(config.ArgQuiet, false)
//		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testIpFailoverVar)
//		viper.Set(core.GetFlagName(cfg.NS, config.ArgLanId), testIpFailoverVar)
//		rm.Lan.EXPECT().Get(testIpFailoverVar, testIpFailoverVar).Return(&testLanIpFailover, nil, nil)
//		err := RunIpFailoverList(cfg)
//		assert.NoError(t, err)
//	})
//}

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
