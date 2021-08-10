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
	"github.com/ionos-cloud/ionosctl/pkg/resources/v6"
	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
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
	testIpBlocks = v6.IpBlocks{
		IpBlocks: ionoscloud.IpBlocks{
			Id:    &testIpBlockVar,
			Items: &[]ionoscloud.IpBlock{testIpBlock},
		},
	}
	newTestIpBlockProperties = v6.IpBlockProperties{
		IpBlockProperties: ionoscloud.IpBlockProperties{
			Name: &newTestIpBlockVar,
		},
	}
	newTestIpBlock = v6.IpBlock{
		IpBlock: ionoscloud.IpBlock{
			Id:         &testIpBlockVar,
			Properties: &newTestIpBlockProperties.IpBlockProperties,
		},
	}
	resTestIpBlock      = v6.IpBlock{IpBlock: testIpBlock}
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
		viper.Set(core.GetFlagName(cfg.NS, config.ArgIpBlockId), testIpBlockVar)
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

func TestPreRunIpBlockLocation(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLocation), testIpBlockLocation)
		viper.Set(config.ArgQuiet, false)
		err := PreRunIpBlockLocation(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunIpBlockLocationErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		err := PreRunIpBlockLocation(cfg)
		assert.Error(t, err)
	})
}

func TestRunIpBlockList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		rm.IpBlocks.EXPECT().List().Return(testIpBlocks, nil, nil)
		err := RunIpBlockList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunIpBlockListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		rm.IpBlocks.EXPECT().List().Return(testIpBlocks, nil, testIpBlockErr)
		err := RunIpBlockList(cfg)
		assert.Error(t, err)
		assert.True(t, err == testIpBlockErr)
	})
}

func TestRunIpBlockGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgIpBlockId), testIpBlockVar)
		rm.IpBlocks.EXPECT().Get(testIpBlockVar).Return(&resTestIpBlock, nil, nil)
		err := RunIpBlockGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunIpBlockGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgIpBlockId), testIpBlockVar)
		rm.IpBlocks.EXPECT().Get(testIpBlockVar).Return(&resTestIpBlock, nil, testIpBlockErr)
		err := RunIpBlockGet(cfg)
		assert.Error(t, err)
		assert.True(t, err == testIpBlockErr)
	})
}

func TestRunIpBlockCreate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), testIpBlockVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLocation), testIpBlockLocation)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgSize), testIpBlockSize)
		rm.IpBlocks.EXPECT().Create(testIpBlockVar, testIpBlockLocation, testIpBlockSize).Return(&resTestIpBlock, nil, nil)
		err := RunIpBlockCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunIpBlockCreateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), testIpBlockVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLocation), testIpBlockLocation)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgSize), testIpBlockSize)
		rm.IpBlocks.EXPECT().Create(testIpBlockVar, testIpBlockLocation, testIpBlockSize).Return(&resTestIpBlock, &testResponse, nil)
		err := RunIpBlockCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunIpBlockUpdate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgIpBlockId), testIpBlockVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), newTestIpBlockVar)
		rm.IpBlocks.EXPECT().Update(testIpBlockVar, newTestIpBlockProperties).Return(&newTestIpBlock, nil, nil)
		err := RunIpBlockUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunIpBlockUpdateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgIpBlockId), testIpBlockVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), newTestIpBlockVar)
		rm.IpBlocks.EXPECT().Update(testIpBlockVar, newTestIpBlockProperties).Return(&newTestIpBlock, nil, testIpBlockErr)
		err := RunIpBlockUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunIpBlockDelete(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgIpBlockId), testIpBlockVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.IpBlocks.EXPECT().Delete(testIpBlockVar).Return(nil, nil)
		err := RunIpBlockDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunIpBlockDeleteErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgIpBlockId), testIpBlockVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.IpBlocks.EXPECT().Delete(testIpBlockVar).Return(nil, testIpBlockErr)
		err := RunIpBlockDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunIpBlockDeleteAskForConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgIpBlockId), testIpBlockVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		cfg.Stdin = bytes.NewReader([]byte("YES\n"))
		rm.IpBlocks.EXPECT().Delete(testIpBlockVar).Return(nil, nil)
		err := RunIpBlockDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunIpBlockDeleteAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgIpBlockId), testIpBlockVar)
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

func TestGetIpBlocksIds(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() { return }
	w := bufio.NewWriter(&b)
	err := os.Setenv(ionoscloud.IonosUsernameEnvVar, "user")
	assert.NoError(t, err)
	err = os.Setenv(ionoscloud.IonosPasswordEnvVar, "pass")
	assert.NoError(t, err)
	viper.Set(config.ArgServerUrl, config.DefaultApiURL)
	getIpBlocksIds(w)
	err = w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`401 Unauthorized`)
	assert.True(t, re.Match(b.Bytes()))
}
