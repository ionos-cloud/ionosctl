package commands

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/utils/clierror"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go-bundle/products/compute"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	// Resources
	serverCreate = resources.Server{
		Server: ionoscloud.Server{
			Properties: &ionoscloud.ServerProperties{
				Name:             &testServerVar,
				Cores:            &cores,
				Ram:              &ram,
				CpuFamily:        &testServerVar,
				AvailabilityZone: &testServerVar,
				Type:             &testServerEnterpriseType,
			},
		},
	}
	serverCubeCreate = resources.Server{
		Server: ionoscloud.Server{
			Properties: &ionoscloud.ServerProperties{
				Name:             &testServerVar,
				Type:             &testServerCubeType,
				TemplateUuid:     &testServerVar,
				CpuFamily:        &testCpuFamilyType,
				AvailabilityZone: &testServerVar,
			},
			Entities: &ionoscloud.ServerEntities{
				Volumes: &ionoscloud.AttachedVolumes{
					Items: &[]ionoscloud.Volume{
						{
							Properties: &ionoscloud.VolumeProperties{
								Name:        &testServerVar,
								Bus:         &testServerVar,
								Type:        &testVolumeType,
								LicenceType: &testLicenceType,
							},
						},
					},
				},
			},
		},
	}
	serverCubeCreateImg = resources.Server{
		Server: ionoscloud.Server{
			Properties: &ionoscloud.ServerProperties{
				Name:             &testServerVar,
				Type:             &testServerCubeType,
				TemplateUuid:     &testServerVar,
				CpuFamily:        &testCpuFamilyType,
				AvailabilityZone: &testServerVar,
			},
			Entities: &ionoscloud.ServerEntities{
				Volumes: &ionoscloud.AttachedVolumes{
					Items: &[]ionoscloud.Volume{
						{
							Properties: &ionoscloud.VolumeProperties{
								Name:          &testServerVar,
								Bus:           &testServerVar,
								Type:          &testVolumeType,
								ImageAlias:    &testServerVar,
								ImagePassword: &testServerVar,
							},
						},
					},
				},
			},
		},
	}
	s = ionoscloud.Server{
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
	ssList = resources.Servers{
		Servers: ionoscloud.Servers{
			Id: &testServerVar,
			Items: &[]ionoscloud.Server{
				s,
				s,
			},
		},
	}
	serverProperties = resources.ServerProperties{
		ServerProperties: ionoscloud.ServerProperties{
			Name:             &testServerNewVar,
			Cores:            &coresNew,
			Ram:              &ramNew,
			CpuFamily:        &testServerNewVar,
			AvailabilityZone: &testServerNewVar,
			BootVolume: &ionoscloud.ResourceReference{
				Id: &testServerVar,
			},
			BootCdrom: &ionoscloud.ResourceReference{
				Id: &testServerVar,
			},
		},
	}
	serverNew = resources.Server{
		Server: ionoscloud.Server{
			Metadata: &ionoscloud.DatacenterElementMetadata{
				State: &state,
			},
			Id:         &testServerVar,
			Properties: &serverProperties.ServerProperties,
		},
	}
	// Resources Attributes
	cores                    = int32(2)
	coresNew                 = int32(4)
	ram                      = int32(256)
	ramNew                   = int32(256)
	state                    = "ACTIVE"
	testServerVar            = "test-server"
	testServerNewVar         = "test-new-server"
	testVolumeType           = "DAS"
	testLicenceType          = "UNKNOWN"
	testCpuFamilyType        = "INTEL_SKYLAKE"
	testServerCubeType       = serverCubeType
	testServerEnterpriseType = serverEnterpriseType
	testServerErr            = errors.New("server test: error occurred")
)

func TestServerCmd(t *testing.T) {
	var err error
	core.RootCmdTest.AddCommand(ServerCmd())
	if ok := ServerCmd().IsAvailableCommand(); !ok {
		err = errors.New("non-available cmd")
	}
	assert.NoError(t, err)
}

func TestPreRunServerList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar)
		err := PreRunServerList(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunServerListFilters(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFilters), []string{fmt.Sprintf("createdBy=%s", testQueryParamVar)})
		err := PreRunServerList(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunServerListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFilters), []string{fmt.Sprintf("%s=%s", testQueryParamVar, testQueryParamVar)})
		err := PreRunServerList(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunDcServerIds(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar)
		err := PreRunDcServerIds(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunDcServerIdsRequiredFlagErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		err := PreRunDcServerIds(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunServerCreate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagCores), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagRam), testServerVar)
		err := PreRunServerCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunServerCreateImageAlias(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagCores), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagRam), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgImageAlias), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgSshKeyPaths), []string{testServerVar})
		err := PreRunServerCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunServerCreateCube(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgType), testServerCubeType)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTemplateId), testServerVar)
		err := PreRunServerCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunServerCreateCubeImgId(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgType), testServerCubeType)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTemplateId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgImageId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPassword), testServerVar)
		err := PreRunServerCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunServerCreateCubeImgErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgType), testServerCubeType)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTemplateId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgImageId), testServerVar)
		err := PreRunServerCreate(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunServerCreateCubeErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		err := PreRunServerCreate(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunServerCreateCubeImg(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgType), testServerCubeType)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTemplateId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgImageId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPassword), testServerVar)
		err := PreRunServerCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunServerListAll(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.Datacenter.EXPECT().List(gomock.AssignableToTypeOf(testListQueryParam)).Return(dcs, &testResponse, nil)
		rm.CloudApiV6Mocks.Server.EXPECT().List(testDatacenterVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(ss, &testResponse, nil).Times(len(getDataCenters(dcs)))
		err := RunServerListAll(cfg)
		assert.NoError(t, err)
	})
}

func TestRunServerList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar)
		rm.CloudApiV6Mocks.Server.EXPECT().List(testServerVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(ss, &testResponse, nil)
		err := RunServerList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunServerListQueryParams(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFilters), []string{fmt.Sprintf("%s=%s", testQueryParamVar, testQueryParamVar)})
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgOrderBy), testQueryParamVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagMaxResults), testMaxResultsVar)
		rm.CloudApiV6Mocks.Server.EXPECT().List(testServerVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(resources.Servers{}, &testResponse, nil)
		err := RunServerList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunServerListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar)
		rm.CloudApiV6Mocks.Server.EXPECT().List(testServerVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(ss, nil, testServerErr)
		err := RunServerList(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar)
		rm.CloudApiV6Mocks.Server.EXPECT().Get(testServerVar, testServerVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&resources.Server{Server: s}, &testResponse, nil)
		err := RunServerGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunServerGetWait(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForState), true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar)
		rm.CloudApiV6Mocks.Server.EXPECT().Get(testServerVar, testServerVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&resources.Server{Server: s}, nil, nil)
		rm.CloudApiV6Mocks.Server.EXPECT().Get(testServerVar, testServerVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&resources.Server{Server: s}, nil, nil)
		err := RunServerGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunServerGetWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForState), true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar)
		rm.CloudApiV6Mocks.Server.EXPECT().Get(testServerVar, testServerVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&resources.Server{Server: s}, nil, testServerErr)
		err := RunServerGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar)
		rm.CloudApiV6Mocks.Server.EXPECT().Get(testServerVar, testServerVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&resources.Server{Server: s}, nil, testServerErr)
		err := RunServerGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerCreate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagCpuFamily), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagAvailabilityZone), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgType), testServerEnterpriseType)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagCores), cores)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagRam), strconv.Itoa(int(ram)))
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.Server.EXPECT().Create(testServerVar, serverCreate, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&resources.Server{Server: s}, &testResponse, nil)
		err := RunServerCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunServerCreateCube(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgVolumeName), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgType), testServerCubeType)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgBus), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTemplateId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagAvailabilityZone), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLicenceType), testLicenceType)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.Server.EXPECT().Create(testServerVar, serverCubeCreate, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&resources.Server{Server: s}, nil, nil)
		err := RunServerCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunServerCreateCubeImgAlias(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgVolumeName), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgType), testServerCubeType)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgBus), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTemplateId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagAvailabilityZone), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgImageAlias), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPassword), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.Server.EXPECT().Create(testServerVar, serverCubeCreateImg, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&resources.Server{Server: s}, nil, nil)
		err := RunServerCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunServerCreateCubeImgSshErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgVolumeName), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgType), testServerCubeType)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgBus), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTemplateId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagAvailabilityZone), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgImageAlias), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPassword), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgSshKeyPaths), []string{testServerVar})
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		err := RunServerCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerCreateWaitState(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagCpuFamily), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagAvailabilityZone), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgType), testServerEnterpriseType)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagCores), cores)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagRam), ram)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForState), true)
		rm.CloudApiV6Mocks.Server.EXPECT().Create(testServerVar, serverCreate, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&resources.Server{Server: s}, nil, nil)
		rm.CloudApiV6Mocks.Server.EXPECT().Get(testServerVar, testServerVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&resources.Server{Server: s}, nil, nil)
		rm.CloudApiV6Mocks.Server.EXPECT().Get(testServerVar, testServerVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&resources.Server{Server: s}, nil, nil)
		err := RunServerCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunServerCreateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagCpuFamily), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagAvailabilityZone), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgType), testServerEnterpriseType)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagCores), cores)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagRam), ram)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.Server.EXPECT().Create(testServerVar, serverCreate, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&resources.Server{Server: s}, nil, testServerErr)
		err := RunServerCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerCreateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagCpuFamily), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagAvailabilityZone), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgType), testServerEnterpriseType)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagCores), cores)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagRam), ram)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), true)
		rm.CloudApiV6Mocks.Server.EXPECT().Create(testServerVar, serverCreate, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&resources.Server{Server: s}, &testResponse, nil)
		rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testRequestIdVar).Return(&testRequestStatus, nil, testRequestErr)
		err := RunServerCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerUpdate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testServerNewVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagCores), coresNew)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagRam), ramNew)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagAvailabilityZone), testServerNewVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagCpuFamily), testServerNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgVolumeId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCdromId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.Server.EXPECT().Update(testServerVar, testServerVar, serverProperties, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&serverNew, &testResponse, nil)
		err := RunServerUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunServerUpdateWaitStateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testServerNewVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagCores), coresNew)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagRam), ramNew)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagAvailabilityZone), testServerNewVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagCpuFamily), testServerNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgVolumeId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCdromId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForState), true)
		rm.CloudApiV6Mocks.Server.EXPECT().Update(testServerVar, testServerVar, serverProperties, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&serverNew, nil, nil)
		rm.CloudApiV6Mocks.Server.EXPECT().Get(testServerVar, testServerVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&serverNew, nil, testServerErr)
		err := RunServerUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerUpdateWaitState(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testServerNewVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagCores), coresNew)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagRam), ramNew)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagAvailabilityZone), testServerNewVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagCpuFamily), testServerNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgVolumeId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCdromId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForState), true)
		rm.CloudApiV6Mocks.Server.EXPECT().Update(testServerVar, testServerVar, serverProperties, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&serverNew, nil, nil)
		rm.CloudApiV6Mocks.Server.EXPECT().Get(testServerVar, testServerVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&serverNew, nil, nil)
		rm.CloudApiV6Mocks.Server.EXPECT().Get(testServerVar, testServerVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&serverNew, nil, nil)
		err := RunServerUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunServerUpdateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testServerNewVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagCores), coresNew)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagRam), ramNew)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagAvailabilityZone), testServerNewVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagCpuFamily), testServerNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgVolumeId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCdromId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.Server.EXPECT().Update(testServerVar, testServerVar, serverProperties, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&serverNew, nil, testServerErr)
		err := RunServerUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerUpdateResponseErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testServerNewVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagCores), coresNew)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagRam), ramNew)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagAvailabilityZone), testServerNewVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagCpuFamily), testServerNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgVolumeId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCdromId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.Server.EXPECT().Update(testServerVar, testServerVar, serverProperties, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&serverNew, &testResponse, testServerErr)
		err := RunServerUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerUpdateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testServerNewVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagCores), coresNew)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagRam), ramNew)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagAvailabilityZone), testServerNewVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagCpuFamily), testServerNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgVolumeId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCdromId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), true)
		rm.CloudApiV6Mocks.Server.EXPECT().Update(testServerVar, testServerVar, serverProperties, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&serverNew, &testResponse, nil)
		rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testRequestIdVar).Return(&testRequestStatus, nil, testRequestErr)
		err := RunServerUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerDelete(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.Server.EXPECT().Delete(testServerVar, testServerVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, nil)
		err := RunServerDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunServerDeleteAll(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.Server.EXPECT().List(testServerVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(ssList, &testResponse, nil)
		rm.CloudApiV6Mocks.Server.EXPECT().Delete(testServerVar, testServerVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, nil)
		rm.CloudApiV6Mocks.Server.EXPECT().Delete(testServerVar, testServerVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, nil)
		err := RunServerDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunServerDeleteAllListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.Server.EXPECT().List(testServerVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(ssList, nil, testServerErr)
		err := RunServerDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerDeleteAllLenErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.Server.EXPECT().List(testServerVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(resources.Servers{}, &testResponse, nil)
		err := RunServerDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerDeleteAllErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.Server.EXPECT().List(testServerVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(ssList, &testResponse, nil)
		rm.CloudApiV6Mocks.Server.EXPECT().Delete(testServerVar, testServerVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, testServerErr)
		rm.CloudApiV6Mocks.Server.EXPECT().Delete(testServerVar, testServerVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, nil)
		err := RunServerDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerDeleteErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.Server.EXPECT().Delete(testServerVar, testServerVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(nil, testServerErr)
		err := RunServerDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerDeleteWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), true)
		rm.CloudApiV6Mocks.Server.EXPECT().Delete(testServerVar, testServerVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, nil)
		rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testRequestIdVar).Return(&testRequestStatus, nil, testRequestErr)
		err := RunServerDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerDeleteAskForConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, false)
		cfg.Stdin = bytes.NewReader([]byte("YES\n"))
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.Server.EXPECT().Delete(testServerVar, testServerVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(nil, nil)
		err := RunServerDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunServerDeleteAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, false)
		cfg.Stdin = os.Stdin
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		err := RunServerDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerSuspend(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.Server.EXPECT().Suspend(testServerVar, testServerVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, nil)
		err := RunServerSuspend(cfg)
		assert.NoError(t, err)
	})
}

func TestRunServerSuspendErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.Server.EXPECT().Suspend(testServerVar, testServerVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(nil, testServerErr)
		err := RunServerSuspend(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerSuspendWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), true)
		rm.CloudApiV6Mocks.Server.EXPECT().Suspend(testServerVar, testServerVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, nil)
		rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testRequestIdVar).Return(&testRequestStatus, nil, testRequestErr)
		err := RunServerSuspend(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerSuspendAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, false)
		cfg.Stdin = os.Stdin
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		err := RunServerSuspend(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerStart(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.Server.EXPECT().Start(testServerVar, testServerVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, nil)
		err := RunServerStart(cfg)
		assert.NoError(t, err)
	})
}

func TestRunServerStartErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.Server.EXPECT().Start(testServerVar, testServerVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(nil, testServerErr)
		err := RunServerStart(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerStartWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), true)
		rm.CloudApiV6Mocks.Server.EXPECT().Start(testServerVar, testServerVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, nil)
		rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testRequestIdVar).Return(&testRequestStatus, nil, testRequestErr)
		err := RunServerStart(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerStartAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, false)
		cfg.Stdin = os.Stdin
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		err := RunServerStart(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerStop(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.Server.EXPECT().Stop(testServerVar, testServerVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, nil)
		err := RunServerStop(cfg)
		assert.NoError(t, err)
	})
}

func TestRunServerStopErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.Server.EXPECT().Stop(testServerVar, testServerVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(nil, testServerErr)
		err := RunServerStop(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerStopWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), true)
		rm.CloudApiV6Mocks.Server.EXPECT().Stop(testServerVar, testServerVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, nil)
		rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testRequestIdVar).Return(&testRequestStatus, nil, testRequestErr)
		err := RunServerStop(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerStopAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, false)
		cfg.Stdin = os.Stdin
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		err := RunServerStop(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerReboot(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar)
		rm.CloudApiV6Mocks.Server.EXPECT().Reboot(testServerVar, testServerVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, nil)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		err := RunServerReboot(cfg)
		assert.NoError(t, err)
	})
}

func TestRunServerRebootErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.Server.EXPECT().Reboot(testServerVar, testServerVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(nil, testServerErr)
		err := RunServerReboot(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerRebootWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), true)
		rm.CloudApiV6Mocks.Server.EXPECT().Reboot(testServerVar, testServerVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, nil)
		rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testRequestIdVar).Return(&testRequestStatus, nil, testRequestErr)
		err := RunServerReboot(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerRebootAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, false)
		cfg.Stdin = os.Stdin
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		err := RunServerReboot(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerResume(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.Server.EXPECT().Resume(testServerVar, testServerVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, nil)
		err := RunServerResume(cfg)
		assert.NoError(t, err)
	})
}

func TestRunServerResumeErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.Server.EXPECT().Resume(testServerVar, testServerVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(nil, testServerErr)
		err := RunServerResume(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerResumeWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), true)
		rm.CloudApiV6Mocks.Server.EXPECT().Resume(testServerVar, testServerVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, nil)
		rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testRequestIdVar).Return(&testRequestStatus, nil, testRequestErr)
		err := RunServerResume(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerResumeAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, false)
		cfg.Stdin = os.Stdin
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		err := RunServerResume(cfg)
		assert.Error(t, err)
	})
}

func TestGetServersCols(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() { return }
	w := bufio.NewWriter(&b)
	viper.Set(core.GetFlagName("server", constants.ArgCols), []string{"Name"})
	getServersCols(core.GetFlagName("server", constants.ArgCols), core.GetFlagName("server", cloudapiv6.ArgAll), w)
	err := w.Flush()
	assert.NoError(t, err)
}

func TestGetServersColsErr(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() { return }
	w := bufio.NewWriter(&b)
	viper.Set(core.GetFlagName("server", constants.ArgCols), []string{"Unknown"})
	getServersCols(core.GetFlagName("server", constants.ArgCols), core.GetFlagName("server", cloudapiv6.ArgAll), w)
	err := w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`unknown column Unknown`)
	assert.True(t, re.Match(b.Bytes()))
}
