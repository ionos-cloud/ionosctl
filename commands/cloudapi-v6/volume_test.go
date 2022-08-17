package commands

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"io"
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

func TestRunVolumeListAll(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.Datacenter.EXPECT().List(testListQueryParam).Return(dcs, &testResponse, nil)
		rm.CloudApiV6Mocks.Volume.EXPECT().List(testDatacenterVar, testListQueryParam).Return(vsList, &testResponse, nil).Times(len(getDataCenters(dcs)))
		err := RunVolumeListAll(cfg)
		assert.NoError(t, err)
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
		rm.CloudApiV6Mocks.Volume.EXPECT().List(testVolumeVar, testListQueryParam).Return(vs, &testResponse, nil)
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
		rm.CloudApiV6Mocks.Volume.EXPECT().List(testVolumeVar, testListQueryParam).Return(vs, nil, testVolumeErr)
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
		rm.CloudApiV6Mocks.Volume.EXPECT().Get(testVolumeVar, testVolumeVar, testQueryParamOther).Return(&resources.Volume{Volume: v}, &testResponse, nil)
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
		rm.CloudApiV6Mocks.Volume.EXPECT().Get(testVolumeVar, testVolumeVar, testQueryParamOther).Return(&resources.Volume{Volume: v}, nil, testVolumeErr)
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
		rm.CloudApiV6Mocks.Volume.EXPECT().Create(testVolumeVar, testVolume, testQueryParamOther).Return(&resources.Volume{Volume: v}, &testResponse, nil)
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
		rm.CloudApiV6Mocks.Volume.EXPECT().Create(testVolumeVar, testVolumeImg, testQueryParamOther).Return(&resources.Volume{Volume: v}, nil, nil)
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
		rm.CloudApiV6Mocks.Volume.EXPECT().Create(testVolumeVar, testVolume, testQueryParamOther).Return(&resources.Volume{Volume: v}, nil, testVolumeErr)
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
		rm.CloudApiV6Mocks.Volume.EXPECT().Create(testVolumeVar, testVolume, testQueryParamOther).Return(&resources.Volume{Volume: v}, &testResponse, nil)
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
		rm.CloudApiV6Mocks.Volume.EXPECT().Update(testVolumeVar, testVolumeVar, volumeProperties, testQueryParamOther).Return(&volumeNew, &testResponse, nil)
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
		rm.CloudApiV6Mocks.Volume.EXPECT().Update(testVolumeVar, testVolumeVar, volumeProperties, testQueryParamOther).Return(&volumeNew, nil, testVolumeErr)
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
		rm.CloudApiV6Mocks.Volume.EXPECT().Update(testVolumeVar, testVolumeVar, volumeProperties, testQueryParamOther).Return(&volumeNew, &testResponse, nil)
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
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(config.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgVolumeId), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.Volume.EXPECT().Delete(testVolumeVar, testVolumeVar, testQueryParamOther).Return(&testResponse, nil)
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
		rm.CloudApiV6Mocks.Volume.EXPECT().List(testVolumeVar, testListQueryParam).Return(vsList, &testResponse, nil)
		rm.CloudApiV6Mocks.Volume.EXPECT().Delete(testVolumeVar, testVolumeVar, testQueryParamOther).Return(&testResponse, nil)
		rm.CloudApiV6Mocks.Volume.EXPECT().Delete(testVolumeVar, testVolumeVar, testQueryParamOther).Return(&testResponse, nil)
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
		rm.CloudApiV6Mocks.Volume.EXPECT().List(testVolumeVar, testListQueryParam).Return(vsList, nil, testVolumeErr)
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
		rm.CloudApiV6Mocks.Volume.EXPECT().List(testVolumeVar, testListQueryParam).Return(resources.Volumes{}, &testResponse, nil)
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
		rm.CloudApiV6Mocks.Volume.EXPECT().List(testVolumeVar, testListQueryParam).Return(
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
		rm.CloudApiV6Mocks.Volume.EXPECT().List(testVolumeVar, testListQueryParam).Return(vsList, &testResponse, nil)
		rm.CloudApiV6Mocks.Volume.EXPECT().Delete(testVolumeVar, testVolumeVar, testQueryParamOther).Return(&testResponse, testVolumeErr)
		rm.CloudApiV6Mocks.Volume.EXPECT().Delete(testVolumeVar, testVolumeVar, testQueryParamOther).Return(&testResponse, nil)
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
		rm.CloudApiV6Mocks.Volume.EXPECT().Delete(testVolumeVar, testVolumeVar, testQueryParamOther).Return(nil, testVolumeErr)
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
		rm.CloudApiV6Mocks.Volume.EXPECT().Delete(testVolumeVar, testVolumeVar, testQueryParamOther).Return(&testResponse, nil)
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
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, false)
		cfg.Stdin = bytes.NewReader([]byte("YES\n"))
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgVolumeId), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.Volume.EXPECT().Delete(testVolumeVar, testVolumeVar, testQueryParamOther).Return(nil, nil)
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
	getVolumesCols(core.GetGlobalFlagName("volume", config.ArgCols), core.GetFlagName("volume", cloudapiv6.ArgAll), w)
	err := w.Flush()
	assert.NoError(t, err)
}

func TestGetVolumesColsErr(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() { return }
	w := bufio.NewWriter(&b)
	viper.Set(core.GetGlobalFlagName("volume", config.ArgCols), []string{"Unknown"})
	getVolumesCols(core.GetGlobalFlagName("volume", config.ArgCols), core.GetFlagName("volume", cloudapiv6.ArgAll), w)
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
		rm.CloudApiV6Mocks.Server.EXPECT().AttachVolume(testServerVar, testServerVar, testServerVar, testQueryParamOther).Return(&resources.Volume{Volume: v}, nil, nil)
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
		rm.CloudApiV6Mocks.Server.EXPECT().AttachVolume(testServerVar, testServerVar, testServerVar, testQueryParamOther).Return(&resources.Volume{Volume: v}, nil, testVolumeErr)
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
		rm.CloudApiV6Mocks.Server.EXPECT().AttachVolume(testServerVar, testServerVar, testServerVar, testQueryParamOther).Return(&resources.Volume{Volume: v}, &testResponse, nil)
		rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testRequestIdVar).Return(&testRequestStatus, nil, testRequestErr)
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
		rm.CloudApiV6Mocks.Server.EXPECT().ListVolumes(testServerVar, testServerVar, testListQueryParam).Return(vsAttached, nil, nil)
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
		rm.CloudApiV6Mocks.Server.EXPECT().ListVolumes(testServerVar, testServerVar, testListQueryParam).Return(vsAttached, nil, testVolumeErr)
		err := RunServerVolumesList(cfg)
		assert.Error(t, err)
	})
}

/// NEW TESTS

type FlagValuePair struct {
	flag  string
	value interface{}
}

type testCase struct {
	name      string
	userInput io.Reader
	args      []FlagValuePair
	calls     func(...*gomock.Call)
	err       bool // To be replaced by `error` type once it makes sense to do so (currently only one type of error is thrown)
}

func ExecuteTestCases(t *testing.T, funcToTest func(c *core.CommandConfig) error, testCases []testCase, cfg *core.CommandConfig) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			viper.Reset()
			for _, argPair := range tc.args {
				viper.Set(argPair.flag, argPair.value)
			}

			if tc.userInput != nil {
				cfg.Stdin = tc.userInput
			}

			// Expected gomock calls, call order, call counts and returned values
			tc.calls()

			err := funcToTest(cfg)

			if tc.err {
				assert.Error(t, err, fmt.Sprintf("Test case '%s' expected an error but got none", tc.name))
			} else {
				assert.NoError(t, err, fmt.Sprintf("Test case '%s' expected no error but got %s", tc.name, err))
			}
		})
	}
}

func TestServerVolumeGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	var funcToTest = RunServerVolumeGet
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		var tests = []testCase{
			{
				name: "server volume get",
				args: []FlagValuePair{
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar},
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar},
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgVolumeId), testServerVar},
				},
				calls: func(...*gomock.Call) {
					gomock.InOrder(
						rm.CloudApiV6Mocks.Server.EXPECT().GetVolume(testServerVar, testServerVar, testServerVar, testQueryParamOther).Return(&resources.Volume{Volume: v}, nil, nil),
					)
				},
				err: false,
			},
			{
				name: "server volume get error",
				args: []FlagValuePair{
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar},
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar},
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgVolumeId), testServerVar},
				},
				calls: func(...*gomock.Call) {
					gomock.InOrder(
						rm.CloudApiV6Mocks.Server.EXPECT().GetVolume(testServerVar, testServerVar, testServerVar, testQueryParamOther).Return(&resources.Volume{Volume: v}, nil, testVolumeErr),
					)
				},
				err: true,
			},
		}
		ExecuteTestCases(t, funcToTest, tests, cfg)
	})
}

func TestServerVolumeDetach(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	var funcToTest = RunServerVolumeDetach
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		var tests = []testCase{
			{
				name: "server volume detach",
				args: []FlagValuePair{
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar},
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar},
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgVolumeId), testServerVar},
					{core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false},
					{config.ArgForce, true},
				},
				calls: func(...*gomock.Call) {
					gomock.InOrder(
						rm.CloudApiV6Mocks.Server.EXPECT().DetachVolume(testServerVar, testServerVar, testServerVar, testQueryParamOther).Return(nil, nil),
					)
				},
				err: false,
			},
			{
				name: "server volume detach error",
				args: []FlagValuePair{
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar},
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar},
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgVolumeId), testServerVar},
					{core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false},
					{config.ArgForce, true},
				},
				calls: func(...*gomock.Call) {
					gomock.InOrder(
						rm.CloudApiV6Mocks.Server.EXPECT().DetachVolume(testServerVar, testServerVar, testServerVar, testQueryParamOther).Return(&testResponse, testVolumeErr),
					)
				},
				err: true,
			},
			{
				name: "server volume detach all",
				args: []FlagValuePair{
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar},
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar},
					{core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false},
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true},
					{config.ArgForce, true},
				},
				calls: func(...*gomock.Call) {
					gomock.InOrder(
						rm.CloudApiV6Mocks.Server.EXPECT().ListVolumes(testServerVar, testServerVar, cloudapiv6.ParentResourceListQueryParams).Return(vsAttachedList, nil, nil),
						rm.CloudApiV6Mocks.Server.EXPECT().DetachVolume(testServerVar, testServerVar, testServerVar, testQueryParamOther).Return(&testResponse, nil),
						rm.CloudApiV6Mocks.Server.EXPECT().DetachVolume(testServerVar, testServerVar, testServerVar, testQueryParamOther).Return(&testResponse, nil),
					)
				},
				err: false,
			},
			{
				name: "server volume detach all (error)",
				args: []FlagValuePair{
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar},
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar},
					{core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false},
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true},
					{config.ArgForce, true},
				},
				calls: func(...*gomock.Call) {
					rm.CloudApiV6Mocks.Server.EXPECT().ListVolumes(testServerVar, testServerVar, cloudapiv6.ParentResourceListQueryParams).Return(vsAttachedList, nil, nil)
					rm.CloudApiV6Mocks.Server.EXPECT().DetachVolume(testServerVar, testServerVar, testServerVar, testQueryParamOther).Return(&testResponse, testVolumeErr)
					rm.CloudApiV6Mocks.Server.EXPECT().DetachVolume(testServerVar, testServerVar, testServerVar, testQueryParamOther).Return(&testResponse, nil)
				},
				err: true,
			},
			{
				name: "server volume detach all (list error)",
				args: []FlagValuePair{
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar},
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar},
					{core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false},
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true},
					{config.ArgForce, true},
				},
				calls: func(...*gomock.Call) {
					gomock.InOrder(
						rm.CloudApiV6Mocks.Server.EXPECT().ListVolumes(
							testServerVar,
							testServerVar,
							cloudapiv6.ParentResourceListQueryParams,
						).Return(
							vsAttachedList,
							nil,
							testVolumeErr,
						),
					)
				},
				err: true,
			},
			{
				name: "server volume detach all (wrong items error)",
				args: []FlagValuePair{
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar},
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar},
					{core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false},
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true},
					{config.ArgForce, true},
				},
				calls: func(...*gomock.Call) {
					gomock.InOrder(
						rm.CloudApiV6Mocks.Server.EXPECT().ListVolumes(
							testServerVar,
							testServerVar,
							cloudapiv6.ParentResourceListQueryParams,
						).Return(
							vsAttachedList,
							nil,
							testVolumeErr,
						),
					)
				},
				err: true,
			},
			{
				name: "server volume detach all (length error)",
				args: []FlagValuePair{
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar},
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar},
					{core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false},
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true},
					{config.ArgForce, true},
				},
				calls: func(...*gomock.Call) {
					gomock.InOrder(
						rm.CloudApiV6Mocks.Server.EXPECT().ListVolumes(
							testServerVar,
							testServerVar,
							cloudapiv6.ParentResourceListQueryParams,
						).Return(
							resources.AttachedVolumes{AttachedVolumes: ionoscloud.AttachedVolumes{Items: &[]ionoscloud.Volume{}}},
							nil,
							nil,
						),
					)
				},
				err: true,
			},
			{
				name: "server volume detach (user confirm)",
				args: []FlagValuePair{
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar},
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar},
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgVolumeId), testServerVar},
					{core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false},
					{config.ArgForce, false},
				},
				calls: func(...*gomock.Call) {
					gomock.InOrder(
						rm.CloudApiV6Mocks.Server.EXPECT().DetachVolume(
							testServerVar,
							testServerVar,
							testServerVar,
							testQueryParamOther,
						).Return(
							nil,
							nil,
						),
					)
				},
				userInput: bytes.NewReader([]byte("YES\n")),
				err:       false,
			},
			{
				name: "server volume detach (wait error)",
				args: []FlagValuePair{
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar},
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar},
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgVolumeId), testServerVar},
					{core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true},
					{config.ArgForce, true},
				},
				calls: func(...*gomock.Call) {
					gomock.InOrder(
						rm.CloudApiV6Mocks.Server.EXPECT().DetachVolume(testServerVar, testServerVar, testServerVar, testQueryParamOther).Return(&testResponse, nil),
						rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testRequestIdVar).Return(&testRequestStatus, nil, testRequestErr),
					)
				},
				err: true,
			},
			{
				name: "server volume detach (user confirm error)",
				args: []FlagValuePair{
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar},
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar},
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgVolumeId), testServerVar},
					{core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false},
					{config.ArgForce, false},
				},
				calls: func(...*gomock.Call) {
					gomock.InOrder()
				},
				userInput: os.Stdin,
				err:       true,
			},
		}
		ExecuteTestCases(t, funcToTest, tests, cfg)
	})
}

// END NEW TESTS
