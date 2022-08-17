package commands

import (
	"github.com/golang/mock/gomock"
	"bufio"
	"bytes"
	"errors"
	"os"
	"regexp"
	"testing"

	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/services/cloudapi-v6/resources"
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
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetGroupId), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIp), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPort), testTargetGroupTargetIntVar)
		err := PreRunTargetGroupIdTargetIpPort(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunTargetGroupIdTargetIpPortErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		err := PreRunTargetGroupIdTargetIpPort(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunTargetGroupTargetRemove(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetGroupId), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIp), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPort), testTargetGroupTargetIntVar)
		err := PreRunTargetGroupTargetRemove(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunTargetGroupTargetRemoveAll(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetGroupId), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		err := PreRunTargetGroupTargetRemove(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunTargetGroupTargetRemoveErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetGroupId), testTargetGroupTargetVar)
		err := PreRunTargetGroupTargetRemove(cfg)
		assert.Error(t, err)
	})
}

func TestRunTargetGroupTargetList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetGroupId), testTargetGroupTargetVar)
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
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgVerbose, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetGroupId), testTargetGroupTargetVar)
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
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetGroupId), testTargetGroupTargetVar)
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
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetGroupId), testTargetGroupTargetVar)
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
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetGroupId), testTargetGroupTargetVar)
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
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetGroupId), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIp), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPort), testTargetGroupTargetIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgWeight), testTargetGroupTargetIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgMaintenanceEnabled), testTargetGroupTargetBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgHealthCheckEnabled), testTargetGroupTargetBoolVar)
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
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetGroupId), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIp), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPort), testTargetGroupTargetIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgWeight), testTargetGroupTargetIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgMaintenanceEnabled), testTargetGroupTargetBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgHealthCheckEnabled), testTargetGroupTargetBoolVar)
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
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetGroupId), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIp), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPort), testTargetGroupTargetIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgWeight), testTargetGroupTargetIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgMaintenanceEnabled), testTargetGroupTargetBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgHealthCheckEnabled), testTargetGroupTargetBoolVar)
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
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetGroupId), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIp), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPort), testTargetGroupTargetIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgWeight), testTargetGroupTargetIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgMaintenanceEnabled), testTargetGroupTargetBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgHealthCheckEnabled), testTargetGroupTargetBoolVar)
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
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetGroupId), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIp), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPort), testTargetGroupTargetIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgWeight), testTargetGroupTargetIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgMaintenanceEnabled), testTargetGroupTargetBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgHealthCheckEnabled), testTargetGroupTargetBoolVar)
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
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetGroupId), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIp), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPort), testTargetGroupTargetIntVar)
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
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetGroupId), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
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
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetGroupId), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIp), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPort), testTargetGroupTargetIntVar)
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
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetGroupId), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIp), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPort), testTargetGroupTargetIntVar)
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
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetGroupId), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIp), "x.x.x.x")
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPort), testTargetGroupTargetIntVar)
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
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetGroupId), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIp), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPort), int32(2))
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
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetGroupId), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIp), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPort), testTargetGroupTargetIntVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
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
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetGroupId), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIp), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPort), testTargetGroupTargetIntVar)
		cfg.Stdin = bytes.NewReader([]byte("YES\n"))
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
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(config.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetGroupId), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().Get(testTargetGroupTargetVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testTargetGroupTargetGetUpdated, nil, nil)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().Update(testTargetGroupTargetVar,
			&resources.TargetGroupProperties{
				TargetGroupProperties: ionoscloud.TargetGroupProperties{
					Targets: &[]ionoscloud.TargetGroupTarget{},
				},
			},
			testQueryParamOther,
		).Return(&testTargetGroupTargetGet, &testResponse, nil)
		err := RunTargetGroupTargetRemove(cfg)
		assert.NoError(t, err)
	})
}

func TestRunTargetGroupTargetRemoveAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetGroupId), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIp), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPort), testTargetGroupTargetIntVar)
		cfg.Stdin = os.Stdin
		err := RunTargetGroupTargetRemove(cfg)
		assert.Error(t, err)
	})
}

func TestGetTargetGroupTargetCols(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	viper.Set(core.GetGlobalFlagName("target", config.ArgCols), []string{"TargetIp"})
	getTargetGroupTargetCols(core.GetGlobalFlagName("target", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
}

func TestGetTargetGroupTargetColsErr(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	viper.Set(core.GetGlobalFlagName("target", config.ArgCols), []string{"Unknown"})
	getTargetGroupTargetCols(core.GetGlobalFlagName("target", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`unknown column Unknown`)
	assert.True(t, re.Match(b.Bytes()))
}
