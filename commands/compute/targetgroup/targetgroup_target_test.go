package targetgroup

import (
	"bufio"
	"bytes"
	"errors"
	"testing"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/testutil"

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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		err := PreRunTargetGroupIdTargetIpPort(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunTargetGroupTargetRemove(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetGroupId), testTargetGroupTargetVar)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().Get(testTargetGroupTargetVar).Return(&testTargetGroupTargetGetUpdated, nil, nil)
		err := RunTargetGroupTargetList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunTargetGroupTargetListResponse(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgVerbose, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetGroupId), testTargetGroupTargetVar)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().Get(testTargetGroupTargetVar).Return(&testTargetGroupTargetGetUpdated, &testutil.TestResponse, nil)
		err := RunTargetGroupTargetList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunTargetGroupTargetListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetGroupId), testTargetGroupTargetVar)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().Get(testTargetGroupTargetVar).Return(&testTargetGroupTargetGetUpdated, nil, testTargetGroupTargetErr)
		err := RunTargetGroupTargetList(cfg)
		assert.Error(t, err)
	})
}

func TestRunTargetGroupTargetListGetTargetsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetGroupId), testTargetGroupTargetVar)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().Get(testTargetGroupTargetVar).Return(&testTargetGroupTargetGet, nil, nil)
		err := RunTargetGroupTargetList(cfg)
		assert.Error(t, err)
	})
}

func TestRunTargetGroupTargetListGetPropertiesErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetGroupId), testTargetGroupTargetVar)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().Get(testTargetGroupTargetVar).Return(&resources.TargetGroup{}, nil, nil)
		err := RunTargetGroupTargetList(cfg)
		assert.Error(t, err)
	})
}

func TestRunTargetGroupTargetAdd(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetGroupId), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIp), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPort), testTargetGroupTargetIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgWeight), testTargetGroupTargetIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgMaintenanceEnabled), testTargetGroupTargetBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgHealthCheckEnabled), testTargetGroupTargetBoolVar)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().Get(testTargetGroupTargetVar).Return(&testTargetGroupTargetGet, nil, nil)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().Update(testTargetGroupTargetVar, &testTargetGroupTargetProperties).Return(&testTargetGroupTargetGetUpdated, nil, nil)
		err := RunTargetGroupTargetAdd(cfg)
		assert.NoError(t, err)
	})
}

func TestRunTargetGroupTargetAddResponse(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetGroupId), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIp), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPort), testTargetGroupTargetIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgWeight), testTargetGroupTargetIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgMaintenanceEnabled), testTargetGroupTargetBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgHealthCheckEnabled), testTargetGroupTargetBoolVar)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().Get(testTargetGroupTargetVar).Return(&testTargetGroupTargetGet, nil, nil)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().Update(testTargetGroupTargetVar, &testTargetGroupTargetProperties).Return(&testTargetGroupTargetGetUpdated, &testutil.TestResponse, nil)
		err := RunTargetGroupTargetAdd(cfg)
		assert.NoError(t, err)
	})
}

func TestRunTargetGroupTargetAddErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetGroupId), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIp), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPort), testTargetGroupTargetIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgWeight), testTargetGroupTargetIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgMaintenanceEnabled), testTargetGroupTargetBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgHealthCheckEnabled), testTargetGroupTargetBoolVar)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().Get(testTargetGroupTargetVar).Return(&testTargetGroupTargetGet, nil, nil)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().Update(testTargetGroupTargetVar, &testTargetGroupTargetProperties).Return(&testTargetGroupTargetGetUpdated, nil, testTargetGroupTargetErr)
		err := RunTargetGroupTargetAdd(cfg)
		assert.Error(t, err)
	})
}

func TestRunTargetGroupTargetAddGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetGroupId), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIp), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPort), testTargetGroupTargetIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgWeight), testTargetGroupTargetIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgMaintenanceEnabled), testTargetGroupTargetBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgHealthCheckEnabled), testTargetGroupTargetBoolVar)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().Get(testTargetGroupTargetVar).Return(&testTargetGroupTargetGet, nil, testTargetGroupTargetErr)
		err := RunTargetGroupTargetAdd(cfg)
		assert.Error(t, err)
	})
}

func TestRunTargetGroupTargetAddWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetGroupId), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIp), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPort), testTargetGroupTargetIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgWeight), testTargetGroupTargetIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgMaintenanceEnabled), testTargetGroupTargetBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgHealthCheckEnabled), testTargetGroupTargetBoolVar)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().Get(testTargetGroupTargetVar).Return(&testTargetGroupTargetGet, nil, nil)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().Update(testTargetGroupTargetVar, &testTargetGroupTargetProperties).Return(&testTargetGroupTargetGetUpdated, &testutil.TestResponse, nil)
		rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testutil.TestRequestIdVar).Return(&testutil.TestRequestStatus, nil, testutil.TestRequestErr)
		err := RunTargetGroupTargetAdd(cfg)
		assert.Error(t, err)
	})
}

func TestRunTargetGroupTargetRemove(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetGroupId), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIp), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPort), testTargetGroupTargetIntVar)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().Get(testTargetGroupTargetVar).Return(&testTargetGroupTargetGetUpdated, nil, nil)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().Update(testTargetGroupTargetVar,
			&resources.TargetGroupProperties{
				TargetGroupProperties: ionoscloud.TargetGroupProperties{
					Targets: &[]ionoscloud.TargetGroupTarget{},
				},
			},
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetGroupId), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().Get(testTargetGroupTargetVar).Return(&testTargetGroupTargetGetUpdated, nil, testTargetGroupTargetErr)
		err := RunTargetGroupTargetRemove(cfg)
		assert.Error(t, err)
	})
}

func TestRunTargetGroupTargetRemoveErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetGroupId), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIp), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPort), testTargetGroupTargetIntVar)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().Get(testTargetGroupTargetVar).Return(&testTargetGroupTargetGetUpdated, nil, nil)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().Update(testTargetGroupTargetVar,
			&resources.TargetGroupProperties{
				TargetGroupProperties: ionoscloud.TargetGroupProperties{
					Targets: &[]ionoscloud.TargetGroupTarget{},
				},
			},
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetGroupId), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIp), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPort), testTargetGroupTargetIntVar)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().Get(testTargetGroupTargetVar).Return(&testTargetGroupTargetGetUpdated, nil, testTargetGroupTargetErr)
		err := RunTargetGroupTargetRemove(cfg)
		assert.Error(t, err)
	})
}

func TestRunTargetGroupTargetRemoveIpErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetGroupId), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIp), "x.x.x.x")
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPort), testTargetGroupTargetIntVar)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().Get(testTargetGroupTargetVar).Return(&testTargetGroupTargetGetUpdated, nil, nil)
		err := RunTargetGroupTargetRemove(cfg)
		assert.Error(t, err)
	})
}

func TestRunTargetGroupTargetRemovePortErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetGroupId), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIp), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPort), int32(2))
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().Get(testTargetGroupTargetVar).Return(&testTargetGroupTargetGetUpdated, nil, nil)
		err := RunTargetGroupTargetRemove(cfg)
		assert.Error(t, err)
	})
}

func TestRunTargetGroupTargetRemoveWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetGroupId), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIp), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPort), testTargetGroupTargetIntVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), true)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().Get(testTargetGroupTargetVar).Return(&testTargetGroupTargetGetUpdated, nil, nil)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().Update(testTargetGroupTargetVar,
			&resources.TargetGroupProperties{
				TargetGroupProperties: ionoscloud.TargetGroupProperties{
					Targets: &[]ionoscloud.TargetGroupTarget{},
				},
			},
		).Return(&testTargetGroupTargetGet, &testutil.TestResponse, nil)
		rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testutil.TestRequestIdVar).Return(&testutil.TestRequestStatus, nil, testutil.TestRequestErr)
		err := RunTargetGroupTargetRemove(cfg)
		assert.Error(t, err)
	})
}

func TestRunTargetGroupTargetRemoveAskForConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetGroupId), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIp), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPort), testTargetGroupTargetIntVar)
		cfg.Command.Command.SetIn(bytes.NewReader([]byte("YES\n")))
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().Get(testTargetGroupTargetVar).Return(&testTargetGroupTargetGetUpdated, nil, nil)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().Update(testTargetGroupTargetVar,
			&resources.TargetGroupProperties{
				TargetGroupProperties: ionoscloud.TargetGroupProperties{
					Targets: &[]ionoscloud.TargetGroupTarget{},
				},
			},
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetGroupId), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().Get(testTargetGroupTargetVar).Return(&testTargetGroupTargetGetUpdated, nil, nil)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().Update(testTargetGroupTargetVar, gomock.AssignableToTypeOf(&testTargetGroupTargetProperties)).Return(&testTargetGroupTargetGet, &testutil.TestResponse, nil)
		err := RunTargetGroupTargetRemove(cfg)
		assert.NoError(t, err)
	})
}

func TestRunTargetGroupTargetRemoveAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetGroupId), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIp), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPort), testTargetGroupTargetIntVar)
		cfg.Command.Command.SetIn(bytes.NewReader([]byte("\n")))
		err := RunTargetGroupTargetRemove(cfg)
		assert.Error(t, err)
	})
}
