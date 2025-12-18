package commands

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
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
	sizeVolume    = float32(12)
	sizeVolumeNew = float32(12)
	zoneVolume    = "ZONE_1"
	v             = ionoscloud.Volume{
		Id: &testVolumeVar,
		Properties: &ionoscloud.VolumeProperties{
			Name:                &testVolumeVar,
			Size:                &sizeVolume,
			LicenceType:         &testVolumeVar,
			Type:                &testVolumeVar,
			Bus:                 &testVolumeVar,
			Image:               &testVolumeVar,
			ImageAlias:          &testVolumeVar,
			AvailabilityZone:    &zoneVolume,
			SshKeys:             &testVolumeSliceVar,
			BackupunitId:        &testVolumeVar,
			UserData:            &testVolumeVar,
			CpuHotPlug:          &testVolumeBoolVar,
			RamHotPlug:          &testVolumeBoolVar,
			NicHotPlug:          &testVolumeBoolVar,
			NicHotUnplug:        &testVolumeBoolVar,
			DiscVirtioHotPlug:   &testVolumeBoolVar,
			DiscVirtioHotUnplug: &testVolumeBoolVar,
			DeviceNumber:        &testDeviceNumberVolumeVar,
			BootServer:          &testVolumeVar,
		},
		Metadata: &ionoscloud.DatacenterElementMetadata{
			State: &testVolumeVar,
		},
	}
	serverVolume = ionoscloud.Volume{
		Id: &testServerVar,
		Properties: &ionoscloud.VolumeProperties{
			Name:                &testVolumeVar,
			Size:                &sizeVolume,
			LicenceType:         &testVolumeVar,
			Type:                &testVolumeVar,
			Bus:                 &testVolumeVar,
			Image:               &testVolumeVar,
			ImageAlias:          &testVolumeVar,
			AvailabilityZone:    &zoneVolume,
			BackupunitId:        &testVolumeVar,
			UserData:            &testVolumeVar,
			CpuHotPlug:          &testVolumeBoolVar,
			RamHotPlug:          &testVolumeBoolVar,
			NicHotPlug:          &testVolumeBoolVar,
			NicHotUnplug:        &testVolumeBoolVar,
			DiscVirtioHotPlug:   &testVolumeBoolVar,
			DiscVirtioHotUnplug: &testVolumeBoolVar,
			BootServer:          &testVolumeVar,
		},
		Metadata: &ionoscloud.DatacenterElementMetadata{
			State: &testVolumeVar,
		},
	}
	testVolume = resources.Volume{
		Volume: ionoscloud.Volume{
			Properties: &ionoscloud.VolumeProperties{
				Name:                &testVolumeVar,
				Size:                &sizeVolume,
				LicenceType:         &testVolumeVar,
				Type:                &testVolumeVar,
				Bus:                 &testVolumeVar,
				AvailabilityZone:    &zoneVolume,
				BackupunitId:        &testVolumeVar,
				UserData:            &testVolumeVar,
				CpuHotPlug:          &testVolumeBoolVar,
				RamHotPlug:          &testVolumeBoolVar,
				NicHotPlug:          &testVolumeBoolVar,
				NicHotUnplug:        &testVolumeBoolVar,
				DiscVirtioHotPlug:   &testVolumeBoolVar,
				DiscVirtioHotUnplug: &testVolumeBoolVar,
			},
		},
	}
	testVolumeImg = resources.Volume{
		Volume: ionoscloud.Volume{
			Properties: &ionoscloud.VolumeProperties{
				Name:                &testVolumeVar,
				Size:                &sizeVolume,
				Image:               &testVolumeVar,
				ImageAlias:          &testVolumeVar,
				ImagePassword:       &testVolumeVar,
				Type:                &testVolumeVar,
				Bus:                 &testVolumeVar,
				AvailabilityZone:    &zoneVolume,
				BackupunitId:        &testVolumeVar,
				UserData:            &testVolumeVar,
				CpuHotPlug:          &testVolumeBoolVar,
				RamHotPlug:          &testVolumeBoolVar,
				NicHotPlug:          &testVolumeBoolVar,
				NicHotUnplug:        &testVolumeBoolVar,
				DiscVirtioHotPlug:   &testVolumeBoolVar,
				DiscVirtioHotUnplug: &testVolumeBoolVar,
			},
		},
	}
	vs = resources.Volumes{
		Volumes: ionoscloud.Volumes{
			Id:    &testVolumeVar,
			Items: &[]ionoscloud.Volume{v},
		},
	}
	vsList = resources.Volumes{
		Volumes: ionoscloud.Volumes{
			Id: &testVolumeVar,
			Items: &[]ionoscloud.Volume{
				v,
				v,
			},
		},
	}
	vsAttachedList = resources.AttachedVolumes{
		AttachedVolumes: ionoscloud.AttachedVolumes{
			Id:    &testVolumeVar,
			Items: &[]ionoscloud.Volume{serverVolume, serverVolume},
		},
	}
	volumeProperties = resources.VolumeProperties{
		VolumeProperties: ionoscloud.VolumeProperties{
			Name:                &testVolumeNewVar,
			Bus:                 &testVolumeNewVar,
			Size:                &sizeVolumeNew,
			CpuHotPlug:          &testVolumeBoolVar,
			RamHotPlug:          &testVolumeBoolVar,
			NicHotPlug:          &testVolumeBoolVar,
			NicHotUnplug:        &testVolumeBoolVar,
			DiscVirtioHotPlug:   &testVolumeBoolVar,
			DiscVirtioHotUnplug: &testVolumeBoolVar,
		},
	}
	volumeNew = resources.Volume{
		Volume: ionoscloud.Volume{
			Id: &testVolumeVar,
			Properties: &ionoscloud.VolumeProperties{
				Name:                volumeProperties.VolumeProperties.Name,
				Size:                volumeProperties.VolumeProperties.Size,
				LicenceType:         &testVolumeVar,
				Type:                &testVolumeVar,
				Bus:                 volumeProperties.VolumeProperties.Bus,
				AvailabilityZone:    &zoneVolume,
				CpuHotPlug:          &testVolumeBoolVar,
				RamHotPlug:          &testVolumeBoolVar,
				NicHotPlug:          &testVolumeBoolVar,
				NicHotUnplug:        &testVolumeBoolVar,
				DiscVirtioHotPlug:   &testVolumeBoolVar,
				DiscVirtioHotUnplug: &testVolumeBoolVar,
			},
		},
	}
	vsAttached = resources.AttachedVolumes{
		AttachedVolumes: ionoscloud.AttachedVolumes{
			Id:    &testVolumeVar,
			Items: &[]ionoscloud.Volume{v},
		},
	}
	testDeviceNumberVolumeVar = int64(1)
	testVolumeVar             = "test-volume"
	testVolumeBoolVar         = false
	testVolumeSliceVar        = []string{"test-volume"}
	testVolumeNewVar          = "test-new-volume"
	testVolumeErr             = errors.New("volume test: error occurred")
)

func TestVolumeCmd(t *testing.T) {
	var err error
	core.RootCmdTest.AddCommand(VolumeCmd())
	if ok := VolumeCmd().IsAvailableCommand(); !ok {
		err = errors.New("non-available cmd")
	}
	assert.NoError(t, err)
}

func TestPreRunVolumeList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testVolumeVar)
		err := PreRunVolumeList(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunVolumeListFilters(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testVolumeVar)
		cfg.Command.Command.Flags().Set(constants.FlagFilters, fmt.Sprintf("createdBy=%s", testQueryParamVar))
		err := PreRunVolumeList(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunVolumeListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testVolumeVar)
		cfg.Command.Command.Flags().Set(constants.FlagFilters, fmt.Sprintf("%s=%s", testQueryParamVar, testQueryParamVar))
		err := PreRunVolumeList(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunDcVolumeIds(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgVolumeId), testVolumeVar)
		err := PreRunDcVolumeIds(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunDcVolumeIdsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		err := PreRunDcVolumeIds(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunVolumeCreate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testVolumeVar)
		err := PreRunVolumeCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunVolumeCreateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		err := PreRunVolumeCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunVolumeListAll(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.Datacenter.EXPECT().List().Return(dcs, &testResponse, nil)
		rm.CloudApiV6Mocks.Volume.EXPECT().List(testDatacenterVar).Return(vsList, &testResponse, nil).Times(len(getDataCenters(dcs)))
		err := RunVolumeListAll(cfg)
		assert.NoError(t, err)
	})
}

func TestRunVolumeList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testVolumeVar)
		rm.CloudApiV6Mocks.Volume.EXPECT().List(testVolumeVar).Return(vs, &testResponse, nil)
		err := RunVolumeList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunVolumeListQueryParams(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testVolumeVar)
		cfg.Command.Command.Flags().Set(constants.FlagFilters, fmt.Sprintf("%s=%s", testQueryParamVar, testQueryParamVar))
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagOrderBy), testQueryParamVar)
		rm.CloudApiV6Mocks.Volume.EXPECT().List(testVolumeVar).Return(resources.Volumes{}, &testResponse, nil)
		err := RunVolumeList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunVolumeListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testVolumeVar)
		rm.CloudApiV6Mocks.Volume.EXPECT().List(testVolumeVar).Return(vs, nil, testVolumeErr)
		err := RunVolumeList(cfg)
		assert.Error(t, err)
	})
}

func TestRunVolumeGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgVolumeId), testVolumeVar)
		rm.CloudApiV6Mocks.Volume.EXPECT().Get(testVolumeVar, testVolumeVar).Return(&resources.Volume{Volume: v}, &testResponse, nil)
		err := RunVolumeGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunVolumeGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgVolumeId), testVolumeVar)
		rm.CloudApiV6Mocks.Volume.EXPECT().Get(testVolumeVar, testVolumeVar).Return(&resources.Volume{Volume: v}, nil, testVolumeErr)
		err := RunVolumeGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunVolumeCreate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgSize), sizeVolume)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLicenceType), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagType), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgBus), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagAvailabilityZone), zoneVolume)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgBackupUnitId), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgUserData), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCpuHotPlug), testVolumeBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRamHotPlug), testVolumeBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicHotPlug), testVolumeBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicHotUnplug), testVolumeBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDiscVirtioHotPlug), testVolumeBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDiscVirtioHotUnplug), testVolumeBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.Volume.EXPECT().Create(testVolumeVar, testVolume).Return(&resources.Volume{Volume: v}, &testResponse, nil)
		err := RunVolumeCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunVolumeCreateImg(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgSize), sizeVolume)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgImageId), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgImageAlias), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPassword), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagType), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgBus), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagAvailabilityZone), zoneVolume)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgBackupUnitId), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgUserData), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCpuHotPlug), testVolumeBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRamHotPlug), testVolumeBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicHotPlug), testVolumeBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicHotUnplug), testVolumeBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDiscVirtioHotPlug), testVolumeBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDiscVirtioHotUnplug), testVolumeBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.Volume.EXPECT().Create(testVolumeVar, testVolumeImg).Return(&resources.Volume{Volume: v}, nil, nil)
		err := RunVolumeCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunVolumeCreateSshKeyErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgSize), sizeVolume)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgImageId), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgImageAlias), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgSshKeyPaths), []string{testVolumeVar})
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagType), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgBus), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagAvailabilityZone), zoneVolume)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgBackupUnitId), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgUserData), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCpuHotPlug), testVolumeBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRamHotPlug), testVolumeBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicHotPlug), testVolumeBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicHotUnplug), testVolumeBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDiscVirtioHotPlug), testVolumeBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDiscVirtioHotUnplug), testVolumeBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		err := RunVolumeCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunVolumeCreateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgSize), sizeVolume)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLicenceType), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagType), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgBus), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagAvailabilityZone), zoneVolume)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgBackupUnitId), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgUserData), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCpuHotPlug), testVolumeBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRamHotPlug), testVolumeBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicHotPlug), testVolumeBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicHotUnplug), testVolumeBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDiscVirtioHotPlug), testVolumeBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDiscVirtioHotUnplug), testVolumeBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.Volume.EXPECT().Create(testVolumeVar, testVolume).Return(&resources.Volume{Volume: v}, nil, testVolumeErr)
		err := RunVolumeCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunVolumeCreateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgSize), sizeVolume)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLicenceType), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagType), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgBus), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagAvailabilityZone), zoneVolume)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgBackupUnitId), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgUserData), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCpuHotPlug), testVolumeBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRamHotPlug), testVolumeBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicHotPlug), testVolumeBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicHotUnplug), testVolumeBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDiscVirtioHotPlug), testVolumeBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDiscVirtioHotUnplug), testVolumeBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), true)
		rm.CloudApiV6Mocks.Volume.EXPECT().Create(testVolumeVar, testVolume).Return(&resources.Volume{Volume: v}, &testResponse, nil)
		rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testRequestIdVar).Return(&testRequestStatus, nil, testRequestErr)
		err := RunVolumeCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunVolumeUpdate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgVolumeId), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testVolumeNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgSize), sizeVolumeNew)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgBus), testVolumeNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCpuHotPlug), testVolumeBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRamHotPlug), testVolumeBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicHotPlug), testVolumeBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicHotUnplug), testVolumeBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDiscVirtioHotPlug), testVolumeBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDiscVirtioHotUnplug), testVolumeBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.Volume.EXPECT().Update(testVolumeVar, testVolumeVar, volumeProperties).Return(&volumeNew, &testResponse, nil)
		err := RunVolumeUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunVolumeUpdateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgVolumeId), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testVolumeNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgSize), sizeVolumeNew)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgBus), testVolumeNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCpuHotPlug), testVolumeBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRamHotPlug), testVolumeBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicHotPlug), testVolumeBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicHotUnplug), testVolumeBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDiscVirtioHotPlug), testVolumeBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDiscVirtioHotUnplug), testVolumeBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.Volume.EXPECT().Update(testVolumeVar, testVolumeVar, volumeProperties).Return(&volumeNew, nil, testVolumeErr)
		err := RunVolumeUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunVolumeUpdateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgVolumeId), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testVolumeNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgSize), sizeVolumeNew)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgBus), testVolumeNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCpuHotPlug), testVolumeBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRamHotPlug), testVolumeBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicHotPlug), testVolumeBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicHotUnplug), testVolumeBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDiscVirtioHotPlug), testVolumeBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDiscVirtioHotUnplug), testVolumeBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), true)
		rm.CloudApiV6Mocks.Volume.EXPECT().Update(testVolumeVar, testVolumeVar, volumeProperties).Return(&volumeNew, &testResponse, nil)
		rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testRequestIdVar).Return(&testRequestStatus, nil, testRequestErr)
		err := RunVolumeUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunVolumeDelete(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgVolumeId), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.Volume.EXPECT().Delete(testVolumeVar, testVolumeVar).Return(&testResponse, nil)
		err := RunVolumeDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunVolumeDeleteAll(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.Volume.EXPECT().List(testVolumeVar).Return(vsList, &testResponse, nil)
		rm.CloudApiV6Mocks.Volume.EXPECT().Delete(testVolumeVar, testVolumeVar).Return(&testResponse, nil)
		rm.CloudApiV6Mocks.Volume.EXPECT().Delete(testVolumeVar, testVolumeVar).Return(&testResponse, nil)
		err := RunVolumeDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunVolumeDeleteAllListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.Volume.EXPECT().List(testVolumeVar).Return(vsList, nil, testVolumeErr)
		err := RunVolumeDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunVolumeDeleteAllItemsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.Volume.EXPECT().List(testVolumeVar).Return(resources.Volumes{}, &testResponse, nil)
		err := RunVolumeDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunVolumeDeleteAllLenErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.Volume.EXPECT().List(testVolumeVar).Return(
			resources.Volumes{Volumes: ionoscloud.Volumes{Items: &[]ionoscloud.Volume{}}}, &testResponse, nil)
		err := RunVolumeDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunVolumeDeleteAllErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.Volume.EXPECT().List(testVolumeVar).Return(vsList, &testResponse, nil)
		rm.CloudApiV6Mocks.Volume.EXPECT().Delete(testVolumeVar, testVolumeVar).Return(&testResponse, testVolumeErr)
		rm.CloudApiV6Mocks.Volume.EXPECT().Delete(testVolumeVar, testVolumeVar).Return(&testResponse, nil)
		err := RunVolumeDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunVolumeDeleteErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgVolumeId), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.Volume.EXPECT().Delete(testVolumeVar, testVolumeVar).Return(nil, testVolumeErr)
		err := RunVolumeDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunVolumeDeleteWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgVolumeId), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), true)
		rm.CloudApiV6Mocks.Volume.EXPECT().Delete(testVolumeVar, testVolumeVar).Return(&testResponse, nil)
		rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testRequestIdVar).Return(&testRequestStatus, nil, testRequestErr)
		err := RunVolumeDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunVolumeDeleteAskForConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, false)
		cfg.Command.Command.SetIn(bytes.NewReader([]byte("YES\n")))
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgVolumeId), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.Volume.EXPECT().Delete(testVolumeVar, testVolumeVar).Return(nil, nil)
		err := RunVolumeDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunVolumeDeleteAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgVolumeId), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		cfg.Command.Command.SetIn(bytes.NewReader([]byte("\n")))
		err := RunVolumeDelete(cfg)
		assert.Error(t, err)
	})
}

// Server Volume

func TestPreRunDcServerIdsRequiredFlagsErr(t *testing.T) {
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

func TestPreRunDcServerVolumeIds(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgVolumeId), testServerVar)
		err := PreRunDcServerVolumeIds(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunServerVolumeList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar)
		err := PreRunServerVolumeList(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunServerVolumeListFilters(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar)
		cfg.Command.Command.Flags().Set(constants.FlagFilters, fmt.Sprintf("createdby=%s", testQueryParamVar))
		err := PreRunServerVolumeList(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunServerVolumeListFiltersErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar)
		cfg.Command.Command.Flags().Set(constants.FlagFilters, fmt.Sprintf("%s=%s", testQueryParamVar, testQueryParamVar))
		err := PreRunServerVolumeList(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunDcServerVolumeIdsRequiredFlagsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		err := PreRunDcServerVolumeIds(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerVolumeAttach(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgVolumeId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.Server.EXPECT().AttachVolume(testServerVar, testServerVar, testServerVar).Return(&resources.Volume{Volume: v}, nil, nil)
		err := RunServerVolumeAttach(cfg)
		assert.NoError(t, err)
	})
}

func TestRunServerVolumeAttachErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgVolumeId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.Server.EXPECT().AttachVolume(testServerVar, testServerVar, testServerVar).Return(&resources.Volume{Volume: v}, nil, testVolumeErr)
		err := RunServerVolumeAttach(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerVolumeAttachWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgVolumeId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), true)
		rm.CloudApiV6Mocks.Server.EXPECT().AttachVolume(testServerVar, testServerVar, testServerVar).Return(&resources.Volume{Volume: v}, &testResponse, nil)
		rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testRequestIdVar).Return(&testRequestStatus, nil, testRequestErr)
		err := RunServerVolumeAttach(cfg)
		assert.Error(t, err)
	})
}

// / NEW TESTS
func TestServerVolumeList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	var funcToTest = RunServerVolumesList
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		var tests = []core.TestCase{
			{
				Name: "server volume list",
				Args: []core.FlagValuePair{
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar},
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar},
				},
				Calls: func(...*gomock.Call) {
					gomock.InOrder(
						rm.CloudApiV6Mocks.Server.EXPECT().ListVolumes(gomock.AssignableToTypeOf(testResourceVar), gomock.AssignableToTypeOf(testResourceVar)).Return(vsAttached, nil, nil),
					)
				},
				ExpectedErr: false,
			},
			{
				Name: "server volume list error",
				Args: []core.FlagValuePair{
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar},
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar},
				},
				Calls: func(...*gomock.Call) {
					gomock.InOrder(
						rm.CloudApiV6Mocks.Server.EXPECT().ListVolumes(gomock.AssignableToTypeOf(testResourceVar), gomock.AssignableToTypeOf(testResourceVar)).Return(vsAttached, nil, testVolumeErr),
					)
				},
				ExpectedErr: true,
			},
		}
		core.ExecuteTestCases(t, funcToTest, tests, cfg)
	})
}

func TestServerVolumeGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	var funcToTest = RunServerVolumeGet
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		var tests = []core.TestCase{
			{
				Name: "server volume get",
				Args: []core.FlagValuePair{
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar},
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar},
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgVolumeId), testServerVar},
				},
				Calls: func(...*gomock.Call) {
					gomock.InOrder(
						rm.CloudApiV6Mocks.Server.EXPECT().GetVolume(testServerVar, testServerVar, testServerVar).Return(&resources.Volume{Volume: v}, nil, nil),
					)
				},
				ExpectedErr: false,
			},
			{
				Name: "server volume get error",
				Args: []core.FlagValuePair{
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar},
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar},
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgVolumeId), testServerVar},
				},
				Calls: func(...*gomock.Call) {
					gomock.InOrder(
						rm.CloudApiV6Mocks.Server.EXPECT().GetVolume(testServerVar, testServerVar, testServerVar).Return(&resources.Volume{Volume: v}, nil, testVolumeErr),
					)
				},
				ExpectedErr: true,
			},
		}
		core.ExecuteTestCases(t, funcToTest, tests, cfg)
	})
}

func TestServerVolumeDetach(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	var funcToTest = RunServerVolumeDetach
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		var tests = []core.TestCase{
			{
				Name: "server volume detach",
				Args: []core.FlagValuePair{
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar},
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar},
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgVolumeId), testServerVar},
					{core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false},
					{constants.ArgForce, true},
				},
				Calls: func(...*gomock.Call) {
					gomock.InOrder(
						rm.CloudApiV6Mocks.Server.EXPECT().DetachVolume(testServerVar, testServerVar, testServerVar).Return(nil, nil),
					)
				},
				ExpectedErr: false,
			},
			{
				Name: "server volume detach error",
				Args: []core.FlagValuePair{
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar},
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar},
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgVolumeId), testServerVar},
					{core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false},
					{constants.ArgForce, true},
				},
				Calls: func(...*gomock.Call) {
					gomock.InOrder(
						rm.CloudApiV6Mocks.Server.EXPECT().DetachVolume(testServerVar, testServerVar, testServerVar).Return(&testResponse, testVolumeErr),
					)
				},
				ExpectedErr: true,
			},
			{
				Name: "server volume detach all",
				Args: []core.FlagValuePair{
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar},
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar},
					{core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false},
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true},
					{constants.ArgForce, true},
				},
				Calls: func(...*gomock.Call) {
					gomock.InOrder(
						rm.CloudApiV6Mocks.Server.EXPECT().ListVolumes(testServerVar, testServerVar).Return(vsAttachedList, nil, nil),
						rm.CloudApiV6Mocks.Server.EXPECT().DetachVolume(testServerVar, testServerVar, testServerVar).Return(&testResponse, nil),
						rm.CloudApiV6Mocks.Server.EXPECT().DetachVolume(testServerVar, testServerVar, testServerVar).Return(&testResponse, nil),
					)
				},
				ExpectedErr: false,
			},
			{
				Name: "server volume detach all (error)",
				Args: []core.FlagValuePair{
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar},
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar},
					{core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false},
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true},
					{constants.ArgForce, true},
				},
				Calls: func(...*gomock.Call) {
					rm.CloudApiV6Mocks.Server.EXPECT().ListVolumes(testServerVar, testServerVar).Return(vsAttachedList, nil, nil)
					rm.CloudApiV6Mocks.Server.EXPECT().DetachVolume(testServerVar, testServerVar, testServerVar).Return(&testResponse, testVolumeErr)
					rm.CloudApiV6Mocks.Server.EXPECT().DetachVolume(testServerVar, testServerVar, testServerVar).Return(&testResponse, nil)
				},
				ExpectedErr: true,
			},
			{
				Name: "server volume detach all (list error)",
				Args: []core.FlagValuePair{
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar},
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar},
					{core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false},
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true},
					{constants.ArgForce, true},
				},
				Calls: func(...*gomock.Call) {
					gomock.InOrder(
						rm.CloudApiV6Mocks.Server.EXPECT().ListVolumes(
							testServerVar,
							testServerVar,
						).Return(
							vsAttachedList,
							nil,
							testVolumeErr,
						),
					)
				},
				ExpectedErr: true,
			},
			{
				Name: "server volume detach all (wrong items error)",
				Args: []core.FlagValuePair{
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar},
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar},
					{core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false},
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true},
					{constants.ArgForce, true},
				},
				Calls: func(...*gomock.Call) {
					gomock.InOrder(
						rm.CloudApiV6Mocks.Server.EXPECT().ListVolumes(
							testServerVar,
							testServerVar,
						).Return(
							vsAttachedList,
							nil,
							testVolumeErr,
						),
					)
				},
				ExpectedErr: true,
			},
			{
				Name: "server volume detach all (length error)",
				Args: []core.FlagValuePair{
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar},
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar},
					{core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false},
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true},
					{constants.ArgForce, true},
				},
				Calls: func(...*gomock.Call) {
					gomock.InOrder(
						rm.CloudApiV6Mocks.Server.EXPECT().ListVolumes(
							testServerVar,
							testServerVar,
						).Return(
							resources.AttachedVolumes{AttachedVolumes: ionoscloud.AttachedVolumes{Items: &[]ionoscloud.Volume{}}},
							nil,
							nil,
						),
					)
				},
				ExpectedErr: true,
			},
			{
				Name: "server volume detach (user confirm)",
				Args: []core.FlagValuePair{
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar},
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar},
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgVolumeId), testServerVar},
					{core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false},
					{constants.ArgForce, false},
				},
				Calls: func(...*gomock.Call) {
					gomock.InOrder(
						rm.CloudApiV6Mocks.Server.EXPECT().DetachVolume(
							testServerVar,
							testServerVar,
							testServerVar,
						).Return(
							nil,
							nil,
						),
					)
				},
				UserInput:   bytes.NewReader([]byte("YES\n")),
				ExpectedErr: false,
			},
			{
				Name: "server volume detach (wait error)",
				Args: []core.FlagValuePair{
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar},
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar},
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgVolumeId), testServerVar},
					{core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), true},
					{constants.ArgForce, true},
				},
				Calls: func(...*gomock.Call) {
					gomock.InOrder(
						rm.CloudApiV6Mocks.Server.EXPECT().DetachVolume(testServerVar, testServerVar, testServerVar).Return(&testResponse, nil),
						rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testRequestIdVar).Return(&testRequestStatus, nil, testRequestErr),
					)
				},
				ExpectedErr: true,
			},
			{
				Name: "server volume detach (user confirm error)",
				Args: []core.FlagValuePair{
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar},
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar},
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgVolumeId), testServerVar},
					{core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false},
					{constants.ArgForce, false},
				},
				Calls: func(...*gomock.Call) {
					gomock.InOrder()
				},
				UserInput:   bytes.NewReader([]byte("\n")),
				ExpectedErr: true,
			},
		}
		core.ExecuteTestCases(t, funcToTest, tests, cfg)
	})
}

// END NEW TESTS
