package commands

import (
	"bufio"
	"bytes"
	"errors"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/services/cloudapi-v6"
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
	testTargetGroupTargetProperties = resources.TargetGroupProperties{
		TargetGroupProperties: ionoscloud.TargetGroupProperties{
			Targets: &[]ionoscloud.TargetGroupTarget{
				{
					Ip:     &testTargetGroupTargetVar,
					Port:   &testTargetGroupTargetIntVar,
					Weight: &testTargetGroupTargetIntVar,
					HealthCheck: &ionoscloud.TargetGroupTargetHealthCheck{
						Check:         &testTargetGroupTargetBoolVar,
						CheckInterval: &testTargetGroupTargetIntVar,
						Maintenance:   &testTargetGroupTargetBoolVar,
					},
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

func TestPreRunTargetGroupIdTargetIpPort(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetGroupId), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetIp), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetPort), testTargetGroupTargetIntVar)
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

func TestRunTargetGroupTargetList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetGroupId), testTargetGroupTargetVar)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().Get(testTargetGroupTargetVar).Return(&testTargetGroupTargetGetUpdated, nil, nil)
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
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
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
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
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
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetGroupId), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetIp), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetPort), testTargetGroupTargetIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgWeight), testTargetGroupTargetIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCheckInterval), testTargetGroupTargetIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCheck), testTargetGroupTargetBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgMaintenance), testTargetGroupTargetBoolVar)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().Get(testTargetGroupTargetVar).Return(&testTargetGroupTargetGet, nil, nil)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().Update(testTargetGroupTargetVar, &testTargetGroupTargetProperties).Return(&testTargetGroupTargetGetUpdated, nil, nil)
		err := RunTargetGroupTargetAdd(cfg)
		assert.NoError(t, err)
	})
}

func TestRunTargetGroupTargetAddResponseErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetGroupId), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetIp), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetPort), testTargetGroupTargetIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgWeight), testTargetGroupTargetIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCheckInterval), testTargetGroupTargetIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCheck), testTargetGroupTargetBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgMaintenance), testTargetGroupTargetBoolVar)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().Get(testTargetGroupTargetVar).Return(&testTargetGroupTargetGet, nil, nil)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().Update(testTargetGroupTargetVar, &testTargetGroupTargetProperties).Return(&testTargetGroupTargetGetUpdated, &testResponse, nil)
		err := RunTargetGroupTargetAdd(cfg)
		assert.Error(t, err)
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
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetIp), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetPort), testTargetGroupTargetIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgWeight), testTargetGroupTargetIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCheckInterval), testTargetGroupTargetIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCheck), testTargetGroupTargetBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgMaintenance), testTargetGroupTargetBoolVar)
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
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetGroupId), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetIp), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetPort), testTargetGroupTargetIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgWeight), testTargetGroupTargetIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCheckInterval), testTargetGroupTargetIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCheck), testTargetGroupTargetBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgMaintenance), testTargetGroupTargetBoolVar)
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
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetGroupId), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetIp), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetPort), testTargetGroupTargetIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgWeight), testTargetGroupTargetIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCheckInterval), testTargetGroupTargetIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCheck), testTargetGroupTargetBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgMaintenance), testTargetGroupTargetBoolVar)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().Get(testTargetGroupTargetVar).Return(&testTargetGroupTargetGet, nil, nil)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().Update(testTargetGroupTargetVar, &testTargetGroupTargetProperties).Return(&testTargetGroupTargetGetUpdated, nil, nil)
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
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetIp), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetPort), testTargetGroupTargetIntVar)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().Get(testTargetGroupTargetVar).Return(&testTargetGroupTargetGetUpdated, nil, nil)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().Update(testTargetGroupTargetVar,
			&resources.TargetGroupProperties{
				TargetGroupProperties: ionoscloud.TargetGroupProperties{
					Targets: &[]ionoscloud.TargetGroupTarget{},
				},
			}).Return(&testTargetGroupTargetGet, nil, nil)
		err := RunTargetGroupTargetRemove(cfg)
		assert.NoError(t, err)
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
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetIp), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetPort), testTargetGroupTargetIntVar)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().Get(testTargetGroupTargetVar).Return(&testTargetGroupTargetGetUpdated, nil, nil)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().Update(testTargetGroupTargetVar,
			&resources.TargetGroupProperties{
				TargetGroupProperties: ionoscloud.TargetGroupProperties{
					Targets: &[]ionoscloud.TargetGroupTarget{},
				},
			}).Return(&testTargetGroupTargetGet, nil, testTargetGroupTargetErr)
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
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetIp), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetPort), testTargetGroupTargetIntVar)
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
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetGroupId), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetIp), "x.x.x.x")
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetPort), testTargetGroupTargetIntVar)
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
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetGroupId), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetIp), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetPort), int32(2))
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
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetGroupId), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetIp), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetPort), testTargetGroupTargetIntVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().Get(testTargetGroupTargetVar).Return(&testTargetGroupTargetGetUpdated, nil, nil)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().Update(testTargetGroupTargetVar,
			&resources.TargetGroupProperties{
				TargetGroupProperties: ionoscloud.TargetGroupProperties{
					Targets: &[]ionoscloud.TargetGroupTarget{},
				},
			}).Return(&testTargetGroupTargetGet, nil, nil)
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
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetIp), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetPort), testTargetGroupTargetIntVar)
		cfg.Stdin = bytes.NewReader([]byte("YES\n"))
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().Get(testTargetGroupTargetVar).Return(&testTargetGroupTargetGetUpdated, nil, nil)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().Update(testTargetGroupTargetVar,
			&resources.TargetGroupProperties{
				TargetGroupProperties: ionoscloud.TargetGroupProperties{
					Targets: &[]ionoscloud.TargetGroupTarget{},
				},
			}).Return(&testTargetGroupTargetGet, nil, nil)
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
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetIp), testTargetGroupTargetVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetPort), testTargetGroupTargetIntVar)
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
