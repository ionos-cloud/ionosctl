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
	"github.com/ionos-cloud/ionosctl/pkg/utils"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	cores    = int32(2)
	coresNew = int32(4)
	ram      = int32(256)
	ramNew   = int32(256)
	state    = "ACTIVE"
	s        = ionoscloud.Server{
		Id: &testServerVar,
		Metadata: &ionoscloud.DatacenterElementMetadata{
			State: &state,
		},
		Properties: &ionoscloud.ServerProperties{
			Name:             &testServerVar,
			Cores:            &cores,
			Ram:              &ram,
			CpuFamily:        &testServerVar,
			AvailabilityZone: &testServerVar,
		},
	}
	ss = resources.Servers{
		Servers: ionoscloud.Servers{
			Id:    &testServerVar,
			Items: &[]ionoscloud.Server{s},
		},
	}
	serverProperties = resources.ServerProperties{
		ServerProperties: ionoscloud.ServerProperties{
			Name:             &testServerNewVar,
			Cores:            &coresNew,
			Ram:              &ramNew,
			CpuFamily:        &testServerNewVar,
			AvailabilityZone: &testServerNewVar,
		},
	}
	serverNew = resources.Server{
		Server: ionoscloud.Server{
			Metadata: &ionoscloud.DatacenterElementMetadata{
				State: &state,
			},
			Id: &testServerVar,
			Properties: &ionoscloud.ServerProperties{
				Name:             serverProperties.ServerProperties.Name,
				Cores:            serverProperties.ServerProperties.Cores,
				Ram:              serverProperties.ServerProperties.Ram,
				CpuFamily:        serverProperties.ServerProperties.CpuFamily,
				AvailabilityZone: serverProperties.ServerProperties.AvailabilityZone,
			},
		},
	}
	dcServer = ionoscloud.Datacenter{
		Id: &testServerVar,
		Properties: &ionoscloud.DatacenterProperties{
			Name:        &testServerVar,
			Description: &testServerVar,
			Location:    &testServerVar,
		},
		Entities: &ionoscloud.DataCenterEntities{
			Servers: &ss.Servers,
		},
	}
	testServerVar    = "test-server"
	testServerNewVar = "test-new-server"
	testServerErr    = errors.New("server test: error occurred")
)

func TestPreRunGlobalDcIdServerIdValidate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testServerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgServerId), testServerVar)
		err := PreRunGlobalDcIdServerIdValidate(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunGlobalDcIdServerIdValidate_RequiredFlagErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), "")
		err := PreRunGlobalDcIdServerIdValidate(cfg)
		assert.Error(t, err)
		assert.True(t, err.Error() == utils.NewRequiredFlagErr(config.ArgDataCenterId).Error())

		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testServerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgServerId), "")
		err = PreRunGlobalDcIdServerIdValidate(cfg)
		assert.Error(t, err)
		assert.True(t, err.Error() == utils.NewRequiredFlagErr(config.ArgServerId).Error())
	})
}

func TestRunServerList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testServerVar)
		rm.Server.EXPECT().List(testServerVar).Return(ss, nil, nil)
		err := RunServerList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunServerList_Err(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testServerVar)
		rm.Server.EXPECT().List(testServerVar).Return(ss, nil, testServerErr)
		err := RunServerList(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testServerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgServerId), testServerVar)
		rm.Server.EXPECT().Get(testServerVar, testServerVar).Return(&resources.Server{s}, nil, nil)
		err := RunServerGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunServerGet_Err(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testServerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgServerId), testServerVar)
		rm.Server.EXPECT().Get(testServerVar, testServerVar).Return(&resources.Server{s}, nil, testServerErr)
		err := RunServerGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerCreate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testServerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgServerName), testServerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgServerCPUFamily), testServerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgServerCores), cores)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgServerZone), testServerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgServerRAM), ram)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), false)
		rm.Server.EXPECT().Create(testServerVar, testServerVar, testServerVar, testServerVar, cores, ram).Return(&resources.Server{s}, nil, nil)
		err := RunServerCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunServerCreate_Err(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testServerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgServerName), testServerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgServerCPUFamily), testServerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgServerCores), cores)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgServerRAM), ram)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), false)
		rm.Server.EXPECT().Create(testServerVar, testServerVar, testServerVar, testServerVar, cores, ram).Return(&resources.Server{s}, nil, testServerErr)
		err := RunServerCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerCreate_WaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testServerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgServerName), testServerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgServerCPUFamily), testServerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgServerCores), cores)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgServerRAM), ram)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), true)
		rm.Server.EXPECT().Create(testServerVar, testServerVar, testServerVar, testServerVar, cores, ram).Return(&resources.Server{s}, nil, nil)
		err := RunServerCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerUpdate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testServerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgServerId), testServerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgServerName), testServerNewVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgServerCores), coresNew)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgServerRAM), ramNew)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgServerZone), testServerNewVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgServerCPUFamily), testServerNewVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), false)
		rm.Server.EXPECT().Update(testServerVar, testServerVar, serverProperties).Return(&serverNew, nil, nil)
		err := RunServerUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunServerUpdate_Err(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testServerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgServerId), testServerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgServerName), testServerNewVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgServerCores), coresNew)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgServerRAM), ramNew)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgServerZone), testServerNewVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgServerCPUFamily), testServerNewVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), false)
		rm.Server.EXPECT().Update(testServerVar, testServerVar, serverProperties).Return(&serverNew, nil, testServerErr)
		err := RunServerUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerUpdate_WaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testServerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgServerId), testServerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgServerName), testServerNewVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgServerCores), coresNew)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgServerRAM), ramNew)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgServerZone), testServerNewVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgServerCPUFamily), testServerNewVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), true)
		rm.Server.EXPECT().Update(testServerVar, testServerVar, serverProperties).Return(&serverNew, nil, nil)
		err := RunServerUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerDelete(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgIgnoreStdin, true)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testServerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgServerId), testServerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), false)
		rm.Server.EXPECT().Delete(testServerVar, testServerVar).Return(nil, nil)
		err := RunServerDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunServerDelete_Err(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgIgnoreStdin, true)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testServerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgServerId), testServerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), false)
		rm.Server.EXPECT().Delete(testServerVar, testServerVar).Return(nil, testServerErr)
		err := RunServerDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerDelete_WaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgIgnoreStdin, true)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testServerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgServerId), testServerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), true)
		rm.Server.EXPECT().Delete(testServerVar, testServerVar).Return(nil, nil)
		err := RunServerDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerDelete_AskForConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgIgnoreStdin, false)
		cfg.Stdin = bytes.NewReader([]byte("YES\n"))
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testServerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgServerId), testServerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), false)
		rm.Server.EXPECT().Delete(testServerVar, testServerVar).Return(nil, nil)
		err := RunServerDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunServerDelete_AskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgIgnoreStdin, false)
		cfg.Stdin = os.Stdin
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testServerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgServerId), testServerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), false)
		err := RunServerDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerStart(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgIgnoreStdin, true)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testServerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgServerId), testServerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), false)
		rm.Server.EXPECT().Start(testServerVar, testServerVar).Return(nil, nil)
		err := RunServerStart(cfg)
		assert.NoError(t, err)
	})
}

func TestRunServerStart_Err(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgIgnoreStdin, true)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testServerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgServerId), testServerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), false)
		rm.Server.EXPECT().Start(testServerVar, testServerVar).Return(nil, testServerErr)
		err := RunServerStart(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerStart_WaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgIgnoreStdin, true)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testServerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgServerId), testServerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), true)
		rm.Server.EXPECT().Start(testServerVar, testServerVar).Return(nil, nil)
		err := RunServerStart(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerStart_AskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgIgnoreStdin, false)
		cfg.Stdin = os.Stdin
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testServerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgServerId), testServerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), false)
		err := RunServerStart(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerStop(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgIgnoreStdin, true)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testServerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgServerId), testServerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), false)
		rm.Server.EXPECT().Stop(testServerVar, testServerVar).Return(nil, nil)
		err := RunServerStop(cfg)
		assert.NoError(t, err)
	})
}

func TestRunServerStop_Err(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgIgnoreStdin, true)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testServerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgServerId), testServerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), false)
		rm.Server.EXPECT().Stop(testServerVar, testServerVar).Return(nil, testServerErr)
		err := RunServerStop(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerStop_WaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgIgnoreStdin, true)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testServerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgServerId), testServerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), true)
		rm.Server.EXPECT().Stop(testServerVar, testServerVar).Return(nil, nil)
		err := RunServerStop(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerStop_AskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgIgnoreStdin, false)
		cfg.Stdin = os.Stdin
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testServerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgServerId), testServerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), false)
		err := RunServerStop(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerReboot(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgIgnoreStdin, true)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testServerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgServerId), testServerVar)
		rm.Server.EXPECT().Reboot(testServerVar, testServerVar).Return(nil, nil)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), false)
		err := RunServerReboot(cfg)
		assert.NoError(t, err)
	})
}

func TestRunServerReboot_Err(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgIgnoreStdin, true)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testServerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgServerId), testServerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), false)
		rm.Server.EXPECT().Reboot(testServerVar, testServerVar).Return(nil, testServerErr)
		err := RunServerReboot(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerReboot_WaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgIgnoreStdin, true)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testServerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgServerId), testServerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), true)
		rm.Server.EXPECT().Reboot(testServerVar, testServerVar).Return(nil, nil)
		err := RunServerReboot(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerReboot_AskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgIgnoreStdin, false)
		cfg.Stdin = os.Stdin
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testServerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgServerId), testServerVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), false)
		err := RunServerReboot(cfg)
		assert.Error(t, err)
	})
}

func TestGetServersCols(t *testing.T) {
	defer func(a func()) { utils.ErrAction = a }(utils.ErrAction)
	var b bytes.Buffer
	utils.ErrAction = func() {}

	w := bufio.NewWriter(&b)
	viper.Set(builder.GetGlobalFlagName("server", config.ArgCols), []string{"Name"})
	getServersCols(builder.GetGlobalFlagName("server", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
}

func TestGetServersCols_Err(t *testing.T) {
	defer func(a func()) { utils.ErrAction = a }(utils.ErrAction)
	var b bytes.Buffer
	utils.ErrAction = func() {}

	w := bufio.NewWriter(&b)
	viper.Set(builder.GetGlobalFlagName("server", config.ArgCols), []string{"Unknown"})
	getServersCols(builder.GetGlobalFlagName("server", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`unknown column Unknown`)
	assert.True(t, re.Match(b.Bytes()))
}

func TestGetServersIds(t *testing.T) {
	defer func(a func()) { utils.ErrAction = a }(utils.ErrAction)
	var b bytes.Buffer
	utils.ErrAction = func() {}

	w := bufio.NewWriter(&b)
	viper.Set(config.ArgConfig, "../pkg/testdata/config.json")
	getServersIds(w, "server")
	err := w.Flush()
	assert.NoError(t, err)

	re := regexp.MustCompile(`401 Unauthorized`)
	assert.True(t, re.Match(b.Bytes()))
}
