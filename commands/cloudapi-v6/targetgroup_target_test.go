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
	testTargetGroupTargetProperties = resources.TargetGroupProperties{
		TargetGroupProperties: ionoscloud.TargetGroupProperties{
			Targets: &[]ionoscloud.TargetGroupTarget{
				{
					Ip:                 &testTargetGroupTargetVar,
					Port:               &testTargetGroupTargetIntVar,
					Weight:             &testTargetGroupTargetIntVar,
					HealthCheckEnabled: &testTargetGroupTargetBoolVar,
					MaintenanceEnabled: &testTargetGroupTargetBoolVar,
				},
			},
		},
	}
	testTargetGroupTargetGet = resources.TargetGroup{
		TargetGroup: ionoscloud.TargetGroup{
			Id:         &testTargetGroupTargetVar,
			Properties: &ionoscloud.TargetGroupProperties{},
		},
	}
	testTargetGroupTargetGetUpdated = resources.TargetGroup{
		TargetGroup: ionoscloud.TargetGroup{
			Id:         &testTargetGroupTargetVar,
			Properties: &testTargetGroupTargetProperties.TargetGroupProperties,
		},
	}
	testTargetGroupTargetIntVar  = int32(1)
	testTargetGroupTargetBoolVar = false
	testTargetGroupTargetVar     = "test-targetgroup-target"
	testTargetGroupTargetErr     = errors.New("targetgroup-target test error")
)

func TestTargetGroupTargetCmd(t *testing.T) {
	var err error
	core.RootCmdTest.AddCommand(TargetGroupTargetCmd())
	if ok := TargetGroupTargetCmd().IsAvailableCommand(); !ok {
		err = errors.New("non-available cmd")
	}
	assert.NoError(t, err)
}

func TestPreRunTargetGroupIdTargetIpPort(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagTargetGroupId), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagIp), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagPort), testTargetGroupTargetIntVar)
		err := PreRunTargetGroupIdTargetIpPort(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunTargetGroupIdTargetIpPortErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		err := PreRunTargetGroupIdTargetIpPort(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunTargetGroupTargetRemove(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagTargetGroupId), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagIp), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagPort), testTargetGroupTargetIntVar)
		err := PreRunTargetGroupTargetRemove(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunTargetGroupTargetRemoveAll(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagTargetGroupId), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagAll), true)
		err := PreRunTargetGroupTargetRemove(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunTargetGroupTargetRemoveErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagTargetGroupId), testTargetGroupTargetVar)
		err := PreRunTargetGroupTargetRemove(cfg)
		assert.Error(t, err)
	})
}

func TestRunTargetGroupTargetList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagTargetGroupId), testTargetGroupTargetVar)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().Get(testTargetGroupTargetVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testTargetGroupTargetGetUpdated, nil, nil)
		err := RunTargetGroupTargetList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunTargetGroupTargetListResponse(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagVerbose, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagTargetGroupId), testTargetGroupTargetVar)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().Get(testTargetGroupTargetVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testTargetGroupTargetGetUpdated, &testResponse, nil)
		err := RunTargetGroupTargetList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunTargetGroupTargetListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagTargetGroupId), testTargetGroupTargetVar)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().Get(testTargetGroupTargetVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testTargetGroupTargetGetUpdated, nil, testTargetGroupTargetErr)
		err := RunTargetGroupTargetList(cfg)
		assert.Error(t, err)
	})
}

func TestRunTargetGroupTargetListGetTargetsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagTargetGroupId), testTargetGroupTargetVar)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().Get(testTargetGroupTargetVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testTargetGroupTargetGet, nil, nil)
		err := RunTargetGroupTargetList(cfg)
		assert.Error(t, err)
	})
}

func TestRunTargetGroupTargetListGetPropertiesErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagTargetGroupId), testTargetGroupTargetVar)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().Get(testTargetGroupTargetVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&resources.TargetGroup{}, nil, nil)
		err := RunTargetGroupTargetList(cfg)
		assert.Error(t, err)
	})
}

func TestRunTargetGroupTargetAdd(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagTargetGroupId), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagIp), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagPort), testTargetGroupTargetIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagWeight), testTargetGroupTargetIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagMaintenanceEnabled), testTargetGroupTargetBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagHealthCheckEnabled), testTargetGroupTargetBoolVar)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().Get(testTargetGroupTargetVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testTargetGroupTargetGet, nil, nil)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().Update(testTargetGroupTargetVar, &testTargetGroupTargetProperties, testQueryParamOther).Return(&testTargetGroupTargetGetUpdated, nil, nil)
		err := RunTargetGroupTargetAdd(cfg)
		assert.NoError(t, err)
	})
}

func TestRunTargetGroupTargetAddResponse(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagTargetGroupId), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagIp), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagPort), testTargetGroupTargetIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagWeight), testTargetGroupTargetIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagMaintenanceEnabled), testTargetGroupTargetBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagHealthCheckEnabled), testTargetGroupTargetBoolVar)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().Get(testTargetGroupTargetVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testTargetGroupTargetGet, nil, nil)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().Update(testTargetGroupTargetVar, &testTargetGroupTargetProperties, testQueryParamOther).Return(&testTargetGroupTargetGetUpdated, &testResponse, nil)
		err := RunTargetGroupTargetAdd(cfg)
		assert.NoError(t, err)
	})
}

func TestRunTargetGroupTargetAddErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagTargetGroupId), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagIp), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagPort), testTargetGroupTargetIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagWeight), testTargetGroupTargetIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagMaintenanceEnabled), testTargetGroupTargetBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagHealthCheckEnabled), testTargetGroupTargetBoolVar)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().Get(testTargetGroupTargetVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testTargetGroupTargetGet, nil, nil)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().Update(testTargetGroupTargetVar, &testTargetGroupTargetProperties, testQueryParamOther).Return(&testTargetGroupTargetGetUpdated, nil, testTargetGroupTargetErr)
		err := RunTargetGroupTargetAdd(cfg)
		assert.Error(t, err)
	})
}

func TestRunTargetGroupTargetAddGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagTargetGroupId), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagIp), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagPort), testTargetGroupTargetIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagWeight), testTargetGroupTargetIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagMaintenanceEnabled), testTargetGroupTargetBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagHealthCheckEnabled), testTargetGroupTargetBoolVar)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().Get(testTargetGroupTargetVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testTargetGroupTargetGet, nil, testTargetGroupTargetErr)
		err := RunTargetGroupTargetAdd(cfg)
		assert.Error(t, err)
	})
}

func TestRunTargetGroupTargetAddWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagWaitForRequest), true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagTargetGroupId), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagIp), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagPort), testTargetGroupTargetIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagWeight), testTargetGroupTargetIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagMaintenanceEnabled), testTargetGroupTargetBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagHealthCheckEnabled), testTargetGroupTargetBoolVar)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().Get(testTargetGroupTargetVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testTargetGroupTargetGet, nil, nil)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().Update(testTargetGroupTargetVar, &testTargetGroupTargetProperties, testQueryParamOther).Return(&testTargetGroupTargetGetUpdated, &testResponse, nil)
		rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testRequestIdVar).Return(&testRequestStatus, nil, testRequestErr)
		err := RunTargetGroupTargetAdd(cfg)
		assert.Error(t, err)
	})
}

func TestRunTargetGroupTargetRemove(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagTargetGroupId), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagIp), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagPort), testTargetGroupTargetIntVar)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().Get(testTargetGroupTargetVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testTargetGroupTargetGetUpdated, nil, nil)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().Update(testTargetGroupTargetVar,
			&resources.TargetGroupProperties{
				TargetGroupProperties: ionoscloud.TargetGroupProperties{
					Targets: &[]ionoscloud.TargetGroupTarget{},
				},
			},
			testQueryParamOther,
		).Return(&testTargetGroupTargetGet, nil, nil)
		err := RunTargetGroupTargetRemove(cfg)
		assert.NoError(t, err)
	})
}

func TestRunTargetGroupTargetRemoveAllErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagTargetGroupId), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagAll), true)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().Get(testTargetGroupTargetVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testTargetGroupTargetGetUpdated, nil, testTargetGroupTargetErr)
		err := RunTargetGroupTargetRemove(cfg)
		assert.Error(t, err)
	})
}

func TestRunTargetGroupTargetRemoveErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagTargetGroupId), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagIp), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagPort), testTargetGroupTargetIntVar)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().Get(testTargetGroupTargetVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testTargetGroupTargetGetUpdated, nil, nil)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().Update(testTargetGroupTargetVar,
			&resources.TargetGroupProperties{
				TargetGroupProperties: ionoscloud.TargetGroupProperties{
					Targets: &[]ionoscloud.TargetGroupTarget{},
				},
			},
			testQueryParamOther,
		).Return(&testTargetGroupTargetGet, nil, testTargetGroupTargetErr)
		err := RunTargetGroupTargetRemove(cfg)
		assert.Error(t, err)
	})
}

func TestRunTargetGroupTargetRemoveGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagTargetGroupId), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagIp), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagPort), testTargetGroupTargetIntVar)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().Get(testTargetGroupTargetVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testTargetGroupTargetGetUpdated, nil, testTargetGroupTargetErr)
		err := RunTargetGroupTargetRemove(cfg)
		assert.Error(t, err)
	})
}

func TestRunTargetGroupTargetRemoveIpErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagTargetGroupId), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagIp), "x.x.x.x")
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagPort), testTargetGroupTargetIntVar)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().Get(testTargetGroupTargetVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testTargetGroupTargetGetUpdated, nil, nil)
		err := RunTargetGroupTargetRemove(cfg)
		assert.Error(t, err)
	})
}

func TestRunTargetGroupTargetRemovePortErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagTargetGroupId), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagIp), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagPort), int32(2))
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().Get(testTargetGroupTargetVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testTargetGroupTargetGetUpdated, nil, nil)
		err := RunTargetGroupTargetRemove(cfg)
		assert.Error(t, err)
	})
}

func TestRunTargetGroupTargetRemoveWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagTargetGroupId), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagIp), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagPort), testTargetGroupTargetIntVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagWaitForRequest), true)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().Get(testTargetGroupTargetVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testTargetGroupTargetGetUpdated, nil, nil)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().Update(testTargetGroupTargetVar,
			&resources.TargetGroupProperties{
				TargetGroupProperties: ionoscloud.TargetGroupProperties{
					Targets: &[]ionoscloud.TargetGroupTarget{},
				},
			},
			testQueryParamOther,
		).Return(&testTargetGroupTargetGet, &testResponse, nil)
		rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testRequestIdVar).Return(&testRequestStatus, nil, testRequestErr)
		err := RunTargetGroupTargetRemove(cfg)
		assert.Error(t, err)
	})
}

func TestRunTargetGroupTargetRemoveAskForConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagForce, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagTargetGroupId), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagIp), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagPort), testTargetGroupTargetIntVar)
		cfg.Command.Command.SetIn(bytes.NewReader([]byte("YES\n")))
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().Get(testTargetGroupTargetVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testTargetGroupTargetGetUpdated, nil, nil)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().Update(testTargetGroupTargetVar,
			&resources.TargetGroupProperties{
				TargetGroupProperties: ionoscloud.TargetGroupProperties{
					Targets: &[]ionoscloud.TargetGroupTarget{},
				},
			},
			testQueryParamOther,
		).Return(&testTargetGroupTargetGet, nil, nil)
		err := RunTargetGroupTargetRemove(cfg)
		assert.NoError(t, err)
	})
}

func TestRunTargetGroupTargetRemoveAll(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagForce, true)
		viper.Set(constants.FlagVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagTargetGroupId), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagAll), true)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().Get(testTargetGroupTargetVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testTargetGroupTargetGetUpdated, nil, nil)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().Update(testTargetGroupTargetVar, gomock.AssignableToTypeOf(&testTargetGroupTargetProperties), gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testTargetGroupTargetGet, &testResponse, nil)
		err := RunTargetGroupTargetRemove(cfg)
		assert.NoError(t, err)
	})
}

func TestRunTargetGroupTargetRemoveAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagForce, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagTargetGroupId), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagIp), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagPort), testTargetGroupTargetIntVar)
		cfg.Command.Command.SetIn(bytes.NewReader([]byte("\n")))
		err := RunTargetGroupTargetRemove(cfg)
		assert.Error(t, err)
	})
}
