package commands

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/resources"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestRunIpConsumersList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgIpBlockId), testIpBlockVar)
		rm.IpBlocks.EXPECT().Get(testIpBlockVar).Return(&resources.IpBlock{IpBlock: testIpBlock}, nil, nil)
		err := RunIpConsumersList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunIpConsumersListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(core.GetFlagName(cfg.NS, config.ArgIpBlockId), testIpBlockVar)
		rm.IpBlocks.EXPECT().Get(testIpBlockVar).Return(&resources.IpBlock{IpBlock: testIpBlock}, nil, testIpBlockErr)
		err := RunIpConsumersList(cfg)
		assert.Error(t, err)
	})
}
