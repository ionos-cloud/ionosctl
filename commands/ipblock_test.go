package commands

import (
	"bufio"
	"bytes"
	"errors"
	"os"
	"regexp"
	"testing"

	"github.com/ionos-cloud/ionosctl/pkg/builder"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/resources"
	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
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

func TestPreRunIpBlockIdValidate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgIpBlockId), testIpBlockVar)
		viper.Set(config.ArgQuiet, false)
		err := PreRunIpBlockIdValidate(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunIpBlockIdValidateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgIpBlockId), "")
		viper.Set(config.ArgQuiet, false)
		err := PreRunIpBlockIdValidate(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunIpBlockLocationValidate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgIpBlockLocation), testIpBlockLocation)
		viper.Set(config.ArgQuiet, false)
		err := PreRunIpBlockLocationValidate(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunIpBlockLocationValidateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgIpBlockLocation), "")
		viper.Set(config.ArgQuiet, false)
		err := PreRunIpBlockLocationValidate(cfg)
		assert.Error(t, err)
	})
}

func TestRunIpBlockList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
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
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
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
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgIpBlockId), testIpBlockVar)
		rm.IpBlocks.EXPECT().Get(testIpBlockVar).Return(&resTestIpBlock, nil, nil)
		err := RunIpBlockGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunIpBlockGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgIpBlockId), testIpBlockVar)
		rm.IpBlocks.EXPECT().Get(testIpBlockVar).Return(&resTestIpBlock, nil, testIpBlockErr)
		err := RunIpBlockGet(cfg)
		assert.Error(t, err)
		assert.True(t, err == testIpBlockErr)
	})
}

func TestRunIpBlockCreate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgIpBlockName), testIpBlockVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgIpBlockLocation), testIpBlockLocation)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgIpBlockSize), testIpBlockSize)
		rm.IpBlocks.EXPECT().Create(testIpBlockVar, testIpBlockLocation, testIpBlockSize).Return(&resTestIpBlock, nil, nil)
		err := RunIpBlockCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunIpBlockCreateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgIpBlockName), testIpBlockVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgIpBlockLocation), testIpBlockLocation)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgIpBlockSize), testIpBlockSize)
		rm.IpBlocks.EXPECT().Create(testIpBlockVar, testIpBlockLocation, testIpBlockSize).Return(&resTestIpBlock, &testResponse, nil)
		err := RunIpBlockCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunIpBlockUpdate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgIpBlockId), testIpBlockVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgIpBlockName), newTestIpBlockVar)
		rm.IpBlocks.EXPECT().Update(testIpBlockVar, newTestIpBlockProperties).Return(&newTestIpBlock, nil, nil)
		err := RunIpBlockUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunIpBlockUpdateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgIpBlockId), testIpBlockVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgIpBlockName), newTestIpBlockVar)
		rm.IpBlocks.EXPECT().Update(testIpBlockVar, newTestIpBlockProperties).Return(&newTestIpBlock, nil, testIpBlockErr)
		err := RunIpBlockUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunIpBlockDelete(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgIpBlockId), testIpBlockVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), false)
		rm.IpBlocks.EXPECT().Delete(testIpBlockVar).Return(nil, nil)
		err := RunIpBlockDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunIpBlockDeleteErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgIpBlockId), testIpBlockVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), false)
		rm.IpBlocks.EXPECT().Delete(testIpBlockVar).Return(nil, testIpBlockErr)
		err := RunIpBlockDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunIpBlockDeleteAskForConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgIpBlockId), testIpBlockVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), false)
		cfg.Stdin = bytes.NewReader([]byte("YES\n"))
		rm.IpBlocks.EXPECT().Delete(testIpBlockVar).Return(nil, nil)
		err := RunIpBlockDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunIpBlockDeleteAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgIpBlockId), testIpBlockVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), false)
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
	viper.Set(builder.GetGlobalFlagName("ipblock", config.ArgCols), []string{"IpBlockId"})
	getIpBlocksCols(builder.GetGlobalFlagName("ipblock", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
}

func TestGetIpBlocksColsErr(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() { return }
	w := bufio.NewWriter(&b)
	viper.Set(builder.GetGlobalFlagName("ipblock", config.ArgCols), []string{"Unknown"})
	getIpBlocksCols(builder.GetGlobalFlagName("ipblock", config.ArgCols), w)
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
	viper.Set(config.ArgConfig, "../pkg/testdata/config.json")
	getIpBlocksIds(w)
	err := w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`401 Unauthorized`)
	assert.True(t, re.Match(b.Bytes()))
}
