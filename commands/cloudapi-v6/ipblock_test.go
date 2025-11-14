package commands

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"testing"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"

	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
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
	testIpBlocksList = resources.IpBlocks{
		IpBlocks: ionoscloud.IpBlocks{
			Id: &testIpBlockVar,
			Items: &[]ionoscloud.IpBlock{
				testIpBlock,
				testIpBlock,
			},
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

func TestIpblockCmd(t *testing.T) {
	var err error
	core.RootCmdTest.AddCommand(IpblockCmd())
	if ok := IpblockCmd().IsAvailableCommand(); !ok {
		err = errors.New("non-available cmd")
	}
	assert.NoError(t, err)
}

func TestPreRunIpBlockId(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIpBlockId), testIpBlockVar)
		viper.Set(constants.ArgQuiet, false)
		err := PreRunIpBlockId(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunIpBlockIdErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		err := PreRunIpBlockId(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunIpBlockList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIpBlockId), testIpBlockVar)
		viper.Set(constants.ArgQuiet, false)
		err := PreRunIpblockList(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunIpBlockListFilters(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIpBlockId), testIpBlockVar)
		viper.Set(constants.ArgQuiet, false)
		cfg.Command.Command.Flags().Set(cloudapiv6.ArgFilters, fmt.Sprintf("createdBy=%s", testQueryParamVar))
		err := PreRunIpblockList(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunIpBlockListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIpBlockId), testIpBlockVar)
		viper.Set(constants.ArgQuiet, false)
		cfg.Command.Command.Flags().Set(cloudapiv6.ArgFilters, fmt.Sprintf("%s=%s", testQueryParamVar, testQueryParamVar))
		err := PreRunIpblockList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunIpBlockList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		rm.CloudApiV6Mocks.IpBlocks.EXPECT().List().Return(testIpBlocks, &testResponse, nil)
		err := RunIpBlockList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunIpBlockListQueryParams(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		cfg.Command.Command.Flags().Set(cloudapiv6.ArgFilters, fmt.Sprintf("%s=%s", testQueryParamVar, testQueryParamVar))
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgOrderBy), testQueryParamVar)
		rm.CloudApiV6Mocks.IpBlocks.EXPECT().List().Return(resources.IpBlocks{}, &testResponse, nil)
		err := RunIpBlockList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunIpBlockListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
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
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIpBlockId), testIpBlockVar)
		rm.CloudApiV6Mocks.IpBlocks.EXPECT().Get(testIpBlockVar).Return(&resTestIpBlock, &testResponse, nil)
		err := RunIpBlockGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunIpBlockGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIpBlockId), testIpBlockVar)
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
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testIpBlockVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLocation), testIpBlockLocation)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgSize), testIpBlockSize)
		rm.CloudApiV6Mocks.IpBlocks.EXPECT().Create(testIpBlockVar, testIpBlockLocation, testIpBlockSize).Return(&resTestIpBlock, &testResponse, nil)
		err := RunIpBlockCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunIpBlockCreateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testIpBlockVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLocation), testIpBlockLocation)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgSize), testIpBlockSize)
		rm.CloudApiV6Mocks.IpBlocks.EXPECT().Create(testIpBlockVar, testIpBlockLocation, testIpBlockSize).Return(&resTestIpBlock, &testResponse, testIpBlockErr)
		err := RunIpBlockCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunIpBlockUpdate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIpBlockId), testIpBlockVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), newTestIpBlockVar)
		rm.CloudApiV6Mocks.IpBlocks.EXPECT().Update(testIpBlockVar, newTestIpBlockProperties).Return(&newTestIpBlock, &testResponse, nil)
		err := RunIpBlockUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunIpBlockUpdateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIpBlockId), testIpBlockVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), newTestIpBlockVar)
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
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIpBlockId), testIpBlockVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.IpBlocks.EXPECT().Delete(testIpBlockVar).Return(&testResponse, nil)
		err := RunIpBlockDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunIpBlockDeleteAll(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIpBlockId), testIpBlockVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.IpBlocks.EXPECT().List().Return(testIpBlocksList, &testResponse, nil)
		rm.CloudApiV6Mocks.IpBlocks.EXPECT().Delete(testIpBlockVar).Return(&testResponse, nil)
		rm.CloudApiV6Mocks.IpBlocks.EXPECT().Delete(testIpBlockVar).Return(&testResponse, nil)
		err := RunIpBlockDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunIpBlockDeleteAllListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIpBlockId), testIpBlockVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.IpBlocks.EXPECT().List().Return(testIpBlocksList, nil, testIpBlockErr)
		err := RunIpBlockDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunIpBlockDeleteAllItemsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIpBlockId), testIpBlockVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.IpBlocks.EXPECT().List().Return(resources.IpBlocks{}, &testResponse, nil)
		err := RunIpBlockDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunIpBlockDeleteAllLenErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIpBlockId), testIpBlockVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.IpBlocks.EXPECT().List().Return(
			resources.IpBlocks{IpBlocks: ionoscloud.IpBlocks{Items: &[]ionoscloud.IpBlock{}}}, &testResponse, nil)
		err := RunIpBlockDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunIpBlockDeleteAllErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIpBlockId), testIpBlockVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.IpBlocks.EXPECT().List().Return(testIpBlocksList, &testResponse, nil)
		rm.CloudApiV6Mocks.IpBlocks.EXPECT().Delete(testIpBlockVar).Return(&testResponse, testIpBlockErr)
		rm.CloudApiV6Mocks.IpBlocks.EXPECT().Delete(testIpBlockVar).Return(&testResponse, testIpBlockErr)
		err := RunIpBlockDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunIpBlockDeleteErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIpBlockId), testIpBlockVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
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
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIpBlockId), testIpBlockVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		cfg.Command.Command.SetIn(bytes.NewReader([]byte("YES\n")))
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
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIpBlockId), testIpBlockVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		cfg.Command.Command.SetIn(bytes.NewReader([]byte("\n")))
		err := RunIpBlockDelete(cfg)
		assert.Error(t, err)
	})
}
