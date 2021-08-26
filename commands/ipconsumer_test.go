package commands

import (
	"bufio"
	"bytes"
	"errors"
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
	testIpConsumer = ionoscloud.IpBlock{
		Id: &testIpConsumerVar,
		Properties: &ionoscloud.IpBlockProperties{
			IpConsumers: &[]ionoscloud.IpConsumer{
				{
					Ip:              &testIpConsumerVar,
					Mac:             &testIpConsumerVar,
					NicId:           &testIpConsumerVar,
					ServerId:        &testIpConsumerVar,
					ServerName:      &testIpConsumerVar,
					DatacenterId:    &testIpConsumerVar,
					DatacenterName:  &testIpConsumerVar,
					K8sNodePoolUuid: &testIpConsumerVar,
					K8sClusterUuid:  &testIpConsumerVar,
				},
			},
		},
	}
	testIpConsumerProperties = ionoscloud.IpBlock{
		Id: &testIpConsumerVar,
	}
	testIpConsumerGet = ionoscloud.IpBlock{
		Id: &testIpConsumerVar,
		Properties: &ionoscloud.IpBlockProperties{
			Name: &testIpConsumerVar,
		},
	}
	testIpConsumerVar = "test-ip-consumer"
	testIpConsumerErr = errors.New("ip consumer test: error occurred")
)

func TestRunIpConsumersList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, true)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgIpBlockId), testIpConsumerVar)
		rm.IpBlocks.EXPECT().Get(testIpConsumerVar).Return(&v5.IpBlock{IpBlock: testIpConsumer}, &testResponse, nil)
		err := RunIpConsumersList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunIpConsumersListPropertiesErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgIpBlockId), testIpConsumerVar)
		rm.IpBlocks.EXPECT().Get(testIpConsumerVar).Return(&v5.IpBlock{IpBlock: testIpConsumerProperties}, nil, nil)
		err := RunIpConsumersList(cfg)
		assert.Error(t, err)
	})
}

func TestRunIpConsumersListGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgIpBlockId), testIpConsumerVar)
		rm.IpBlocks.EXPECT().Get(testIpConsumerVar).Return(&v5.IpBlock{IpBlock: testIpConsumerGet}, nil, nil)
		err := RunIpConsumersList(cfg)
		assert.Error(t, err)
	})
}

func TestRunIpConsumersListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgIpBlockId), testIpConsumerVar)
		rm.IpBlocks.EXPECT().Get(testIpConsumerVar).Return(&v5.IpBlock{IpBlock: testIpConsumer}, nil, testIpConsumerErr)
		err := RunIpConsumersList(cfg)
		assert.Error(t, err)
	})
}

func TestGetIpConsumersCols(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() { return }
	w := bufio.NewWriter(&b)
	viper.Reset()
	viper.Set(config.ArgOutput, config.DefaultOutputFormat)
	viper.Set(config.ArgQuiet, false)
	viper.Set(config.ArgServerUrl, config.DefaultApiURL)
	viper.Set(core.GetGlobalFlagName("consumer", config.ArgCols), []string{"Ip"})
	getIpConsumerCols(core.GetGlobalFlagName("consumer", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
}

func TestGetIpConsumerColsErr(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() { return }
	w := bufio.NewWriter(&b)
	viper.Reset()
	viper.Set(config.ArgOutput, config.DefaultOutputFormat)
	viper.Set(config.ArgQuiet, false)
	viper.Set(config.ArgServerUrl, config.DefaultApiURL)
	viper.Set(core.GetGlobalFlagName("consumer", config.ArgCols), []string{"Unknown"})
	getIpConsumerCols(core.GetGlobalFlagName("consumer", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`unknown column Unknown`)
	assert.True(t, re.Match(b.Bytes()))
}
