package commands

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
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
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
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
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFilters), []string{fmt.Sprintf("createdBy=%s", testQueryParamVar)})
		err := PreRunVolumeList(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunVolumeListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFilters), []string{fmt.Sprintf("%s=%s", testQueryParamVar, testQueryParamVar)})
		err := PreRunVolumeList(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunDcVolumeIds(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
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
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		err := PreRunDcVolumeIds(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunVolumeCreate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testVolumeVar)
		err := PreRunVolumeCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunVolumeCreateImg(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgImageId), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPassword), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgSshKeyPaths), []string{testVolumeVar})
		err := PreRunVolumeCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunVolumeCreateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		err := PreRunVolumeCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunVolumeList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testVolumeVar)
		rm.CloudApiV6Mocks.Volume.EXPECT().List(testVolumeVar, resources.ListQueryParams{}).Return(vs, &testResponse, nil)
		err := RunVolumeList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunVolumeListQueryParams(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFilters), []string{fmt.Sprintf("%s=%s", testQueryParamVar, testQueryParamVar)})
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgOrderBy), testQueryParamVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgMaxResults), testMaxResultsVar)
		rm.CloudApiV6Mocks.Volume.EXPECT().List(testVolumeVar, testListQueryParam).Return(resources.Volumes{}, &testResponse, nil)
		err := RunVolumeList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunVolumeListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testVolumeVar)
		rm.CloudApiV6Mocks.Volume.EXPECT().List(testVolumeVar, resources.ListQueryParams{}).Return(vs, nil, testVolumeErr)
		err := RunVolumeList(cfg)
		assert.Error(t, err)
	})
}

func TestRunVolumeGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgVolumeId), testVolumeVar)
		rm.CloudApiV6Mocks.Volume.EXPECT().Get(testVolumeVar, testVolumeVar, resources.QueryParams{}).Return(&resources.Volume{Volume: v}, &testResponse, nil)
		err := RunVolumeGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunVolumeGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgVolumeId), testVolumeVar)
		rm.CloudApiV6Mocks.Volume.EXPECT().Get(testVolumeVar, testVolumeVar, resources.QueryParams{}).Return(&resources.Volume{Volume: v}, nil, testVolumeErr)
		err := RunVolumeGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunVolumeCreate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgSize), sizeVolume)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLicenceType), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgType), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgBus), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAvailabilityZone), zoneVolume)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgBackupUnitId), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgUserData), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCpuHotPlug), testVolumeBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRamHotPlug), testVolumeBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicHotPlug), testVolumeBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicHotUnplug), testVolumeBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDiscVirtioHotPlug), testVolumeBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDiscVirtioHotUnplug), testVolumeBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.Volume.EXPECT().Create(testVolumeVar, testVolume, resources.QueryParams{}).Return(&resources.Volume{Volume: v}, &testResponse, nil)
		err := RunVolumeCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunVolumeCreateImg(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgSize), sizeVolume)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgImageId), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgImageAlias), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPassword), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgType), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgBus), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAvailabilityZone), zoneVolume)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgBackupUnitId), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgUserData), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCpuHotPlug), testVolumeBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRamHotPlug), testVolumeBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicHotPlug), testVolumeBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicHotUnplug), testVolumeBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDiscVirtioHotPlug), testVolumeBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDiscVirtioHotUnplug), testVolumeBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.Volume.EXPECT().Create(testVolumeVar, testVolumeImg, resources.QueryParams{}).Return(&resources.Volume{Volume: v}, nil, nil)
		err := RunVolumeCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunVolumeCreateSshKeyErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgSize), sizeVolume)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgImageId), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgImageAlias), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgSshKeyPaths), []string{testVolumeVar})
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgType), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgBus), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAvailabilityZone), zoneVolume)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgBackupUnitId), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgUserData), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCpuHotPlug), testVolumeBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRamHotPlug), testVolumeBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicHotPlug), testVolumeBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicHotUnplug), testVolumeBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDiscVirtioHotPlug), testVolumeBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDiscVirtioHotUnplug), testVolumeBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		err := RunVolumeCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunVolumeCreateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgSize), sizeVolume)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLicenceType), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgType), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgBus), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAvailabilityZone), zoneVolume)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgBackupUnitId), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgUserData), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCpuHotPlug), testVolumeBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRamHotPlug), testVolumeBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicHotPlug), testVolumeBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicHotUnplug), testVolumeBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDiscVirtioHotPlug), testVolumeBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDiscVirtioHotUnplug), testVolumeBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.Volume.EXPECT().Create(testVolumeVar, testVolume, resources.QueryParams{}).Return(&resources.Volume{Volume: v}, nil, testVolumeErr)
		err := RunVolumeCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunVolumeCreateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgSize), sizeVolume)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLicenceType), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgType), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgBus), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAvailabilityZone), zoneVolume)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgBackupUnitId), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgUserData), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCpuHotPlug), testVolumeBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRamHotPlug), testVolumeBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicHotPlug), testVolumeBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicHotUnplug), testVolumeBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDiscVirtioHotPlug), testVolumeBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDiscVirtioHotUnplug), testVolumeBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		rm.CloudApiV6Mocks.Volume.EXPECT().Create(testVolumeVar, testVolume, resources.QueryParams{}).Return(&resources.Volume{Volume: v}, &testResponse, nil)
		rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testRequestIdVar, resources.QueryParams{}).Return(&testRequestStatus, nil, testRequestErr)
		err := RunVolumeCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunVolumeUpdate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, false)
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
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.Volume.EXPECT().Update(testVolumeVar, testVolumeVar, volumeProperties, resources.QueryParams{}).Return(&volumeNew, &testResponse, nil)
		err := RunVolumeUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunVolumeUpdateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
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
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.Volume.EXPECT().Update(testVolumeVar, testVolumeVar, volumeProperties, resources.QueryParams{}).Return(&volumeNew, nil, testVolumeErr)
		err := RunVolumeUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunVolumeUpdateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
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
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		rm.CloudApiV6Mocks.Volume.EXPECT().Update(testVolumeVar, testVolumeVar, volumeProperties, resources.QueryParams{}).Return(&volumeNew, &testResponse, nil)
		rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testRequestIdVar, resources.QueryParams{}).Return(&testRequestStatus, nil, testRequestErr)
		err := RunVolumeUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunVolumeDelete(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(config.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgVolumeId), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.Volume.EXPECT().Delete(testVolumeVar, testVolumeVar, resources.QueryParams{}).Return(&testResponse, nil)
		err := RunVolumeDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunVolumeDeleteAll(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.Volume.EXPECT().List(testVolumeVar, resources.ListQueryParams{}).Return(vsList, &testResponse, nil)
		rm.CloudApiV6Mocks.Volume.EXPECT().Delete(testVolumeVar, testVolumeVar, resources.QueryParams{}).Return(&testResponse, nil)
		rm.CloudApiV6Mocks.Volume.EXPECT().Delete(testVolumeVar, testVolumeVar, resources.QueryParams{}).Return(&testResponse, nil)
		err := RunVolumeDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunVolumeDeleteAllListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.Volume.EXPECT().List(testVolumeVar, resources.ListQueryParams{}).Return(vsList, nil, testVolumeErr)
		err := RunVolumeDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunVolumeDeleteAllItemsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.Volume.EXPECT().List(testVolumeVar, resources.ListQueryParams{}).Return(resources.Volumes{}, &testResponse, nil)
		err := RunVolumeDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunVolumeDeleteAllLenErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.Volume.EXPECT().List(testVolumeVar, resources.ListQueryParams{}).Return(
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
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.Volume.EXPECT().List(testVolumeVar, resources.ListQueryParams{}).Return(vsList, &testResponse, nil)
		rm.CloudApiV6Mocks.Volume.EXPECT().Delete(testVolumeVar, testVolumeVar, resources.QueryParams{}).Return(&testResponse, testVolumeErr)
		rm.CloudApiV6Mocks.Volume.EXPECT().Delete(testVolumeVar, testVolumeVar, resources.QueryParams{}).Return(&testResponse, nil)
		err := RunVolumeDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunVolumeDeleteErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgVolumeId), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.Volume.EXPECT().Delete(testVolumeVar, testVolumeVar, resources.QueryParams{}).Return(nil, testVolumeErr)
		err := RunVolumeDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunVolumeDeleteWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgVolumeId), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		rm.CloudApiV6Mocks.Volume.EXPECT().Delete(testVolumeVar, testVolumeVar, resources.QueryParams{}).Return(&testResponse, nil)
		rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testRequestIdVar, resources.QueryParams{}).Return(&testRequestStatus, nil, testRequestErr)
		err := RunVolumeDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunVolumeDeleteAskForConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, false)
		cfg.Stdin = bytes.NewReader([]byte("YES\n"))
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgVolumeId), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.Volume.EXPECT().Delete(testVolumeVar, testVolumeVar, resources.QueryParams{}).Return(nil, nil)
		err := RunVolumeDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunVolumeDeleteAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, false)
		cfg.Stdin = os.Stdin
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgVolumeId), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		err := RunVolumeDelete(cfg)
		assert.Error(t, err)
	})
}

func TestVolumesCols(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() { return }
	w := bufio.NewWriter(&b)
	viper.Set(core.GetGlobalFlagName("volume", config.ArgCols), []string{"Name"})
	getVolumesCols(core.GetGlobalFlagName("volume", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
}

func TestGetVolumesColsErr(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() { return }
	w := bufio.NewWriter(&b)
	viper.Set(core.GetGlobalFlagName("volume", config.ArgCols), []string{"Unknown"})
	getVolumesCols(core.GetGlobalFlagName("volume", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`unknown column Unknown`)
	assert.True(t, re.Match(b.Bytes()))
}

// Server Volume

func TestPreRunDcServerIdsRequiredFlagsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		err := PreRunDcServerIds(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunDcServerVolumeIds(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
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
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
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
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFilters), []string{fmt.Sprintf("createdBy=%s", testQueryParamVar)})
		err := PreRunServerVolumeList(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunServerVolumeListFiltersErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFilters), []string{fmt.Sprintf("%s=%s", testQueryParamVar, testQueryParamVar)})
		err := PreRunServerVolumeList(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunDcServerVolumeIdsRequiredFlagsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		err := PreRunDcServerVolumeIds(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerVolumeAttach(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgVolumeId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.Server.EXPECT().AttachVolume(testServerVar, testServerVar, testServerVar, resources.QueryParams{}).Return(&resources.Volume{Volume: v}, nil, nil)
		err := RunServerVolumeAttach(cfg)
		assert.NoError(t, err)
	})
}

func TestRunServerVolumeAttachErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgVolumeId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.Server.EXPECT().AttachVolume(testServerVar, testServerVar, testServerVar, resources.QueryParams{}).Return(&resources.Volume{Volume: v}, nil, testVolumeErr)
		err := RunServerVolumeAttach(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerVolumeAttachWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgVolumeId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		rm.CloudApiV6Mocks.Server.EXPECT().AttachVolume(testServerVar, testServerVar, testServerVar, resources.QueryParams{}).Return(&resources.Volume{Volume: v}, &testResponse, nil)
		rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testRequestIdVar, resources.QueryParams{}).Return(&testRequestStatus, nil, testRequestErr)
		err := RunServerVolumeAttach(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerVolumesList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar)
		rm.CloudApiV6Mocks.Server.EXPECT().ListVolumes(testServerVar, testServerVar, resources.ListQueryParams{}).Return(vsAttached, nil, nil)
		err := RunServerVolumesList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunServerVolumesListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar)
		rm.CloudApiV6Mocks.Server.EXPECT().ListVolumes(testServerVar, testServerVar, resources.ListQueryParams{}).Return(vsAttached, nil, testVolumeErr)
		err := RunServerVolumesList(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerVolumeGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgVolumeId), testServerVar)
		rm.CloudApiV6Mocks.Server.EXPECT().GetVolume(testServerVar, testServerVar, testServerVar, resources.QueryParams{}).Return(&resources.Volume{Volume: v}, nil, nil)
		err := RunServerVolumeGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunServerVolumeGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgVolumeId), testServerVar)
		rm.CloudApiV6Mocks.Server.EXPECT().GetVolume(testServerVar, testServerVar, testServerVar, resources.QueryParams{}).Return(&resources.Volume{Volume: v}, nil, testVolumeErr)
		err := RunServerVolumeGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerVolumeDetach(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgVolumeId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.Server.EXPECT().DetachVolume(testServerVar, testServerVar, testServerVar, resources.QueryParams{}).Return(nil, nil)
		err := RunServerVolumeDetach(cfg)
		assert.NoError(t, err)
	})
}

func TestRunServerVolumeDetachAll(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.Server.EXPECT().ListVolumes(testServerVar, testServerVar, resources.ListQueryParams{}).Return(vsAttachedList, nil, nil)
		rm.CloudApiV6Mocks.Server.EXPECT().DetachVolume(testServerVar, testServerVar, testServerVar, resources.QueryParams{}).Return(&testResponse, nil)
		rm.CloudApiV6Mocks.Server.EXPECT().DetachVolume(testServerVar, testServerVar, testServerVar, resources.QueryParams{}).Return(&testResponse, nil)
		err := RunServerVolumeDetach(cfg)
		assert.NoError(t, err)
	})
}

func TestRunServerVolumeDetachAllListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.Server.EXPECT().ListVolumes(testServerVar, testServerVar, resources.ListQueryParams{}).Return(vsAttachedList, nil, testVolumeErr)
		err := RunServerVolumeDetach(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerVolumeDetachAllItemsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.Server.EXPECT().ListVolumes(testServerVar, testServerVar, resources.ListQueryParams{}).Return(resources.AttachedVolumes{}, nil, nil)
		err := RunServerVolumeDetach(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerVolumeDetachAllLenErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.Server.EXPECT().ListVolumes(testServerVar, testServerVar, resources.ListQueryParams{}).Return(
			resources.AttachedVolumes{AttachedVolumes: ionoscloud.AttachedVolumes{Items: &[]ionoscloud.Volume{}}}, nil, nil)
		err := RunServerVolumeDetach(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerVolumeDetachAllErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.Server.EXPECT().ListVolumes(testServerVar, testServerVar, resources.ListQueryParams{}).Return(vsAttachedList, nil, nil)
		rm.CloudApiV6Mocks.Server.EXPECT().DetachVolume(testServerVar, testServerVar, testServerVar, resources.QueryParams{}).Return(&testResponse, testVolumeErr)
		rm.CloudApiV6Mocks.Server.EXPECT().DetachVolume(testServerVar, testServerVar, testServerVar, resources.QueryParams{}).Return(&testResponse, nil)
		err := RunServerVolumeDetach(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerVolumeDetachErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgVolumeId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.Server.EXPECT().DetachVolume(testServerVar, testServerVar, testServerVar, resources.QueryParams{}).Return(&testResponse, testVolumeErr)
		err := RunServerVolumeDetach(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerVolumeDetachWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgVolumeId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		rm.CloudApiV6Mocks.Server.EXPECT().DetachVolume(testServerVar, testServerVar, testServerVar, resources.QueryParams{}).Return(&testResponse, nil)
		rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testRequestIdVar, resources.QueryParams{}).Return(&testRequestStatus, nil, testRequestErr)
		err := RunServerVolumeDetach(cfg)
		assert.Error(t, err)
	})
}

func TestRunVolumeDetachAskForConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, false)
		cfg.Stdin = bytes.NewReader([]byte("YES\n"))
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgVolumeId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.Server.EXPECT().DetachVolume(testServerVar, testServerVar, testServerVar, resources.QueryParams{}).Return(nil, nil)
		err := RunServerVolumeDetach(cfg)
		assert.NoError(t, err)
	})
}

func TestRunServerVolumeDetachAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, false)
		cfg.Stdin = os.Stdin
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgVolumeId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		err := RunServerVolumeDetach(cfg)
		assert.Error(t, err)
	})
}
