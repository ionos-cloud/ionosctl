package commands

import (
	"bufio"
	"bytes"
	"errors"
	cloudapi_v6 "github.com/ionos-cloud/ionosctl/services/cloudapi-v6"
	"os"
	"regexp"
	"testing"

	"github.com/ionos-cloud/ionosctl/internal/config"
	"github.com/ionos-cloud/ionosctl/internal/core"
	"github.com/ionos-cloud/ionosctl/internal/utils/clierror"
	"github.com/ionos-cloud/ionosctl/services/cloudapi-v6/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	testIpBlock = ionoscloud.IpBlock{
		Id: &testIpBlockVar,
		Properties: &ionoscloud.IpBlockProperties{
			Location: &testIpBlockLocation,
			Size:     &testIpBlockSize,
			Name:     &testIpBlockVar,
			Ips:      &testIpBlockIpsVar,
			IpConsumers: &[]ionoscloud.IpConsumer{
				{
					Ip:              &testIpBlockVar,
					Mac:             &testIpBlockVar,
					NicId:           &testIpBlockVar,
					ServerId:        &testIpBlockVar,
					ServerName:      &testIpBlockVar,
					DatacenterId:    &testIpBlockVar,
					DatacenterName:  &testIpBlockVar,
					K8sNodePoolUuid: &testIpBlockVar,
					K8sClusterUuid:  &testIpBlockVar,
				},
			},
		},
		Metadata: &ionoscloud.DatacenterElementMetadata{
			State: &testIpBlockStateVar,
		},
	}
	testIpBlocks = resources.IpBlocks{
		IpBlocks: ionoscloud.IpBlocks{
			Id:    &testIpBlockVar,
			Items: &[]ionoscloud.IpBlock{testIpBlock},
		},
	}
	newTestIpBlockProperties = resources.IpBlockProperties{
		IpBlockProperties: ionoscloud.IpBlockProperties{
			Name: &newTestIpBlockVar,
		},
	}
	newTestIpBlock = resources.IpBlock{
		IpBlock: ionoscloud.IpBlock{
			Id:         &testIpBlockVar,
			Properties: &newTestIpBlockProperties.IpBlockProperties,
		},
	}
	resTestIpBlock      = resources.IpBlock{IpBlock: testIpBlock}
	testIpBlockVar      = "test-ip-block"
	testIpBlockStateVar = "AVAILABLE"
	testIpBlockIpsVar   = []string{"x.x.x.x"}
	newTestIpBlockVar   = "new-test-ip-block"
	testIpBlockLocation = "us/las"
	testIpBlockSize     = int32(1)
	testIpBlockErr      = errors.New("ip block test: error occurred")
)

func TestPreRunIpBlockId(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgIpBlockId), testIpBlockVar)
		viper.Set(config.ArgQuiet, false)
		err := PreRunIpBlockId(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunIpBlockIdErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		err := PreRunIpBlockId(cfg)
		assert.Error(t, err)
	})
}

func TestRunIpBlockList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		rm.CloudApiV6Mocks.IpBlocks.EXPECT().List().Return(testIpBlocks, nil, nil)
		err := RunIpBlockList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunIpBlockListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		rm.CloudApiV6Mocks.IpBlocks.EXPECT().List().Return(testIpBlocks, nil, testIpBlockErr)
		err := RunIpBlockList(cfg)
		assert.Error(t, err)
		assert.True(t, err == testIpBlockErr)
	})
}

func TestRunIpBlockGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgIpBlockId), testIpBlockVar)
		rm.CloudApiV6Mocks.IpBlocks.EXPECT().Get(testIpBlockVar).Return(&resTestIpBlock, nil, nil)
		err := RunIpBlockGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunIpBlockGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgIpBlockId), testIpBlockVar)
		rm.CloudApiV6Mocks.IpBlocks.EXPECT().Get(testIpBlockVar).Return(&resTestIpBlock, nil, testIpBlockErr)
		err := RunIpBlockGet(cfg)
		assert.Error(t, err)
		assert.True(t, err == testIpBlockErr)
	})
}

func TestRunIpBlockCreate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgName), testIpBlockVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgLocation), testIpBlockLocation)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgSize), testIpBlockSize)
		rm.CloudApiV6Mocks.IpBlocks.EXPECT().Create(testIpBlockVar, testIpBlockLocation, testIpBlockSize).Return(&resTestIpBlock, nil, nil)
		err := RunIpBlockCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunIpBlockCreateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgName), testIpBlockVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgLocation), testIpBlockLocation)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgSize), testIpBlockSize)
		rm.CloudApiV6Mocks.IpBlocks.EXPECT().Create(testIpBlockVar, testIpBlockLocation, testIpBlockSize).Return(&resTestIpBlock, &testResponse, nil)
		err := RunIpBlockCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunIpBlockUpdate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgIpBlockId), testIpBlockVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgName), newTestIpBlockVar)
		rm.CloudApiV6Mocks.IpBlocks.EXPECT().Update(testIpBlockVar, newTestIpBlockProperties).Return(&newTestIpBlock, nil, nil)
		err := RunIpBlockUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunIpBlockUpdateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgIpBlockId), testIpBlockVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgName), newTestIpBlockVar)
		rm.CloudApiV6Mocks.IpBlocks.EXPECT().Update(testIpBlockVar, newTestIpBlockProperties).Return(&newTestIpBlock, nil, testIpBlockErr)
		err := RunIpBlockUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunIpBlockDelete(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgIpBlockId), testIpBlockVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.IpBlocks.EXPECT().Delete(testIpBlockVar).Return(nil, nil)
		err := RunIpBlockDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunIpBlockDeleteErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgIpBlockId), testIpBlockVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.IpBlocks.EXPECT().Delete(testIpBlockVar).Return(nil, testIpBlockErr)
		err := RunIpBlockDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunIpBlockDeleteAskForConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgIpBlockId), testIpBlockVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		cfg.Stdin = bytes.NewReader([]byte("YES\n"))
		rm.CloudApiV6Mocks.IpBlocks.EXPECT().Delete(testIpBlockVar).Return(nil, nil)
		err := RunIpBlockDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunIpBlockDeleteAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapi_v6.ArgIpBlockId), testIpBlockVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		cfg.Stdin = os.Stdin
		err := RunIpBlockDelete(cfg)
		assert.Error(t, err)
	})
}

func TestGetIpBlocksCols(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() { return }
	w := bufio.NewWriter(&b)
	viper.Set(core.GetGlobalFlagName("ipblock", config.ArgCols), []string{"IpBlockId"})
	getIpBlocksCols(core.GetGlobalFlagName("ipblock", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
}

func TestGetIpBlocksColsErr(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() { return }
	w := bufio.NewWriter(&b)
	viper.Set(core.GetGlobalFlagName("ipblock", config.ArgCols), []string{"Unknown"})
	getIpBlocksCols(core.GetGlobalFlagName("ipblock", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`unknown column Unknown`)
	assert.True(t, re.Match(b.Bytes()))
}
