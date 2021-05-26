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
	"github.com/ionos-cloud/ionosctl/pkg/resources"
	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
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
			VmState:          &state,
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

func TestPreRunDcServerIds(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgServerId), testServerVar)
		err := PreRunDcServerIds(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunDcServerIdsRequiredFlagErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), "")
		viper.Set(core.GetFlagName(cfg.NS, config.ArgServerId), "")
		err := PreRunDcServerIds(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testServerVar)
		rm.Server.EXPECT().List(testServerVar).Return(ss, nil, nil)
		err := RunServerList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunServerListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testServerVar)
		rm.Server.EXPECT().List(testServerVar).Return(ss, nil, testServerErr)
		err := RunServerList(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgServerId), testServerVar)
		rm.Server.EXPECT().Get(testServerVar, testServerVar).Return(&resources.Server{Server: s}, nil, nil)
		err := RunServerGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunServerGetWait(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForState), true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgServerId), testServerVar)
		rm.Server.EXPECT().Get(testServerVar, testServerVar).Return(&resources.Server{Server: s}, nil, nil)
		rm.Server.EXPECT().Get(testServerVar, testServerVar).Return(&resources.Server{Server: s}, nil, nil)
		err := RunServerGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunServerGetWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForState), true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgServerId), testServerVar)
		rm.Server.EXPECT().Get(testServerVar, testServerVar).Return(&resources.Server{Server: s}, nil, testServerErr)
		err := RunServerGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgServerId), testServerVar)
		rm.Server.EXPECT().Get(testServerVar, testServerVar).Return(&resources.Server{Server: s}, nil, testServerErr)
		err := RunServerGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerCreate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgCPUFamily), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgCores), cores)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgAvailabilityZone), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgRamSize), ram)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.Server.EXPECT().Create(testServerVar, testServerVar, testServerVar, testServerVar, cores, ram).Return(&resources.Server{Server: s}, nil, nil)
		err := RunServerCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunServerCreateWaitState(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgCPUFamily), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgAvailabilityZone), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgCores), cores)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgRamSize), ram)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForState), true)
		rm.Server.EXPECT().Create(testServerVar, testServerVar, testServerVar, testServerVar, cores, ram).Return(&resources.Server{Server: s}, nil, nil)
		rm.Server.EXPECT().Get(testServerVar, testServerVar).Return(&resources.Server{Server: s}, nil, nil)
		rm.Server.EXPECT().Get(testServerVar, testServerVar).Return(&resources.Server{Server: s}, nil, nil)
		err := RunServerCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunServerCreateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgCPUFamily), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgAvailabilityZone), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgCores), cores)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgRamSize), ram)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.Server.EXPECT().Create(testServerVar, testServerVar, testServerVar, testServerVar, cores, ram).Return(&resources.Server{Server: s}, nil, testServerErr)
		err := RunServerCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerCreateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgCPUFamily), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgAvailabilityZone), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgCores), cores)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgRamSize), ram)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		rm.Server.EXPECT().Create(testServerVar, testServerVar, testServerVar, testServerVar, cores, ram).Return(&resources.Server{Server: s}, nil, nil)
		err := RunServerCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerUpdate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgServerId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), testServerNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgCores), coresNew)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgRamSize), ramNew)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgAvailabilityZone), testServerNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgCPUFamily), testServerNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.Server.EXPECT().Update(testServerVar, testServerVar, serverProperties).Return(&serverNew, nil, nil)
		err := RunServerUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunServerUpdateWaitStateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgServerId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), testServerNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgCores), coresNew)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgRamSize), ramNew)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgAvailabilityZone), testServerNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgCPUFamily), testServerNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForState), true)
		rm.Server.EXPECT().Update(testServerVar, testServerVar, serverProperties).Return(&serverNew, nil, nil)
		rm.Server.EXPECT().Get(testServerVar, testServerVar).Return(&serverNew, nil, testServerErr)
		err := RunServerUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerUpdateWaitState(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgServerId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), testServerNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgCores), coresNew)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgRamSize), ramNew)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgAvailabilityZone), testServerNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgCPUFamily), testServerNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForState), true)
		rm.Server.EXPECT().Update(testServerVar, testServerVar, serverProperties).Return(&serverNew, nil, nil)
		rm.Server.EXPECT().Get(testServerVar, testServerVar).Return(&serverNew, nil, nil)
		rm.Server.EXPECT().Get(testServerVar, testServerVar).Return(&serverNew, nil, nil)
		err := RunServerUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunServerUpdateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgServerId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), testServerNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgCores), coresNew)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgRamSize), ramNew)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgAvailabilityZone), testServerNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgCPUFamily), testServerNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.Server.EXPECT().Update(testServerVar, testServerVar, serverProperties).Return(&serverNew, nil, testServerErr)
		err := RunServerUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerUpdateResponseErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgServerId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), testServerNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgCores), coresNew)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgRamSize), ramNew)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgAvailabilityZone), testServerNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgCPUFamily), testServerNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.Server.EXPECT().Update(testServerVar, testServerVar, serverProperties).Return(&serverNew, &testResponse, nil)
		err := RunServerUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerUpdateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgServerId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), testServerNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgCores), coresNew)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgRamSize), ramNew)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgAvailabilityZone), testServerNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgCPUFamily), testServerNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		rm.Server.EXPECT().Update(testServerVar, testServerVar, serverProperties).Return(&serverNew, nil, nil)
		err := RunServerUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerDelete(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgServerId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.Server.EXPECT().Delete(testServerVar, testServerVar).Return(nil, nil)
		err := RunServerDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunServerDeleteErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgServerId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.Server.EXPECT().Delete(testServerVar, testServerVar).Return(nil, testServerErr)
		err := RunServerDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerDeleteWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgServerId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		rm.Server.EXPECT().Delete(testServerVar, testServerVar).Return(nil, nil)
		err := RunServerDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerDeleteAskForConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, false)
		cfg.Stdin = bytes.NewReader([]byte("YES\n"))
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgServerId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.Server.EXPECT().Delete(testServerVar, testServerVar).Return(nil, nil)
		err := RunServerDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunServerDeleteAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, false)
		cfg.Stdin = os.Stdin
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgServerId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		err := RunServerDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerStart(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgServerId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.Server.EXPECT().Start(testServerVar, testServerVar).Return(nil, nil)
		err := RunServerStart(cfg)
		assert.NoError(t, err)
	})
}

func TestRunServerStartErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgServerId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.Server.EXPECT().Start(testServerVar, testServerVar).Return(nil, testServerErr)
		err := RunServerStart(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerStartWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgServerId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		rm.Server.EXPECT().Start(testServerVar, testServerVar).Return(nil, nil)
		err := RunServerStart(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerStartAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, false)
		cfg.Stdin = os.Stdin
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgServerId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		err := RunServerStart(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerStop(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgServerId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.Server.EXPECT().Stop(testServerVar, testServerVar).Return(nil, nil)
		err := RunServerStop(cfg)
		assert.NoError(t, err)
	})
}

func TestRunServerStopErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgServerId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.Server.EXPECT().Stop(testServerVar, testServerVar).Return(nil, testServerErr)
		err := RunServerStop(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerStopWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgServerId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		rm.Server.EXPECT().Stop(testServerVar, testServerVar).Return(nil, nil)
		err := RunServerStop(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerStopAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, false)
		cfg.Stdin = os.Stdin
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgServerId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		err := RunServerStop(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerReboot(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgServerId), testServerVar)
		rm.Server.EXPECT().Reboot(testServerVar, testServerVar).Return(nil, nil)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		err := RunServerReboot(cfg)
		assert.NoError(t, err)
	})
}

func TestRunServerRebootErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgServerId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.Server.EXPECT().Reboot(testServerVar, testServerVar).Return(nil, testServerErr)
		err := RunServerReboot(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerRebootWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgServerId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		rm.Server.EXPECT().Reboot(testServerVar, testServerVar).Return(nil, nil)
		err := RunServerReboot(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerRebootAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, false)
		cfg.Stdin = os.Stdin
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgServerId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		err := RunServerReboot(cfg)
		assert.Error(t, err)
	})
}

func TestGetServersCols(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() { return }
	w := bufio.NewWriter(&b)
	viper.Set(core.GetGlobalFlagName("server", config.ArgCols), []string{"Name"})
	getServersCols(core.GetGlobalFlagName("server", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
}

func TestGetServersColsErr(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() { return }
	w := bufio.NewWriter(&b)
	viper.Set(core.GetGlobalFlagName("server", config.ArgCols), []string{"Unknown"})
	getServersCols(core.GetGlobalFlagName("server", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`unknown column Unknown`)
	assert.True(t, re.Match(b.Bytes()))
}

func TestGetServersIds(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() { return }
	w := bufio.NewWriter(&b)
	viper.Set(config.ArgConfig, "../pkg/testdata/config.json")
	getServersIds(w, testServerVar)
	err := w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`401 Unauthorized`)
	assert.True(t, re.Match(b.Bytes()))
}
