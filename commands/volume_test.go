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
	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	sizeVolume    = float32(12)
	sizeVolumeNew = float32(12)
	zoneVolume    = "ZONE 1"
	v             = ionoscloud.Volume{
		Id: &testVolumeVar,
		Properties: &ionoscloud.VolumeProperties{
			Name:             &testVolumeVar,
			Size:             &sizeVolume,
			LicenceType:      &testVolumeVar,
			Type:             &testVolumeVar,
			Bus:              &testVolumeVar,
			AvailabilityZone: &zoneVolume,
		},
	}
	vs = resources.Volumes{
		Volumes: ionoscloud.Volumes{
			Id:    &testVolumeVar,
			Items: &[]ionoscloud.Volume{v},
		},
	}
	volumeProperties = resources.VolumeProperties{
		VolumeProperties: ionoscloud.VolumeProperties{
			Name: &testVolumeNewVar,
			Bus:  &testVolumeNewVar,
			Size: &sizeVolumeNew,
		},
	}
	volumeNew = resources.Volume{
		Volume: ionoscloud.Volume{
			Id: &testVolumeVar,
			Properties: &ionoscloud.VolumeProperties{
				Name:             volumeProperties.VolumeProperties.Name,
				Size:             volumeProperties.VolumeProperties.Size,
				LicenceType:      &testVolumeVar,
				Type:             &testVolumeVar,
				Bus:              volumeProperties.VolumeProperties.Bus,
				AvailabilityZone: &zoneVolume,
			},
		},
	}
	vsAttached = resources.AttachedVolumes{
		AttachedVolumes: ionoscloud.AttachedVolumes{
			Id:    &testVolumeVar,
			Items: &[]ionoscloud.Volume{v},
		},
	}
	testVolumeVar    = "test-volume"
	testVolumeNewVar = "test-new-volume"
	testVolumeErr    = errors.New("volume test: error occurred")
)

func TestPreRunGlobalDcIdVolumeIdValidate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testVolumeVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgVolumeId), testVolumeVar)
		err := PreRunGlobalDcIdVolumeIdValidate(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunGlobalDcIdVolumeIdValidate_RequiredFlagsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), "")
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgVolumeId), testVolumeVar)
		err := PreRunGlobalDcIdVolumeIdValidate(cfg)
		assert.Error(t, err)
		assert.True(t, err.Error() == clierror.NewRequiredFlagErr(config.ArgDataCenterId).Error())

		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testVolumeVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgVolumeId), "")
		err = PreRunGlobalDcIdVolumeIdValidate(cfg)
		assert.Error(t, err)
		assert.True(t, err.Error() == clierror.NewRequiredFlagErr(config.ArgVolumeId).Error())
	})
}

func TestPreRunGlobalDcIdServerVolumeIdsValidate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testVolumeVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgServerId), testVolumeVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgVolumeId), testVolumeVar)
		err := PreRunGlobalDcIdServerVolumeIdsValidate(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunGlobalDcIdServerVolumeIdsValidate_RequiredFlagsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), "")
		err := PreRunGlobalDcIdServerVolumeIdsValidate(cfg)
		assert.Error(t, err)
		assert.True(t, err.Error() == clierror.NewRequiredFlagErr(config.ArgDataCenterId).Error())

		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testVolumeVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgServerId), "")
		err = PreRunGlobalDcIdServerVolumeIdsValidate(cfg)
		assert.Error(t, err)
		assert.True(t, err.Error() == clierror.NewRequiredFlagErr(config.ArgServerId).Error())

		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testVolumeVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgServerId), testVolumeVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgVolumeId), "")
		err = PreRunGlobalDcIdServerVolumeIdsValidate(cfg)
		assert.Error(t, err)
		assert.True(t, err.Error() == clierror.NewRequiredFlagErr(config.ArgVolumeId).Error())
	})
}

func TestRunVolumeList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testVolumeVar)
		rm.Volume.EXPECT().List(testVolumeVar).Return(vs, nil, nil)
		err := RunVolumeList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunVolumeList_Err(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testVolumeVar)
		rm.Volume.EXPECT().List(testVolumeVar).Return(vs, nil, testVolumeErr)
		err := RunVolumeList(cfg)
		assert.Error(t, err)
	})
}

func TestRunVolumeGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testVolumeVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgVolumeId), testVolumeVar)
		rm.Volume.EXPECT().Get(testVolumeVar, testVolumeVar).Return(&resources.Volume{v}, nil, nil)
		err := RunVolumeGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunVolumeGet_Err(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testVolumeVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgVolumeId), testVolumeVar)
		rm.Volume.EXPECT().Get(testVolumeVar, testVolumeVar).Return(&resources.Volume{v}, nil, testVolumeErr)
		err := RunVolumeGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunVolumeCreate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testVolumeVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgVolumeName), testVolumeVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgVolumeSize), sizeVolume)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgVolumeLicenceType), testVolumeVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgVolumeType), testVolumeVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgVolumeBus), testVolumeVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgVolumeZone), zoneVolume)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), false)
		rm.Volume.EXPECT().Create(testVolumeVar, testVolumeVar, testVolumeVar, testVolumeVar, testVolumeVar, zoneVolume, sizeVolume).Return(&resources.Volume{v}, nil, nil)
		err := RunVolumeCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunVolumeCreate_Err(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testVolumeVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgVolumeName), testVolumeVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgVolumeSize), sizeVolume)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgVolumeLicenceType), testVolumeVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgVolumeType), testVolumeVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgVolumeBus), testVolumeVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgVolumeZone), zoneVolume)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), false)
		rm.Volume.EXPECT().Create(testVolumeVar, testVolumeVar, testVolumeVar, testVolumeVar, testVolumeVar, zoneVolume, sizeVolume).Return(&resources.Volume{v}, nil, testVolumeErr)
		err := RunVolumeCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunVolumeCreate_WaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testVolumeVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgVolumeName), testVolumeVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgVolumeSize), sizeVolume)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgVolumeLicenceType), testVolumeVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgVolumeType), testVolumeVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgVolumeBus), testVolumeVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgVolumeZone), zoneVolume)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), true)
		rm.Volume.EXPECT().Create(testVolumeVar, testVolumeVar, testVolumeVar, testVolumeVar, testVolumeVar, zoneVolume, sizeVolume).Return(&resources.Volume{v}, nil, nil)
		err := RunVolumeCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunVolumeUpdate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testVolumeVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgVolumeId), testVolumeVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgVolumeName), testVolumeNewVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgVolumeSize), sizeVolumeNew)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgVolumeBus), testVolumeNewVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), false)
		rm.Volume.EXPECT().Update(testVolumeVar, testVolumeVar, volumeProperties).Return(&volumeNew, nil, nil)
		err := RunVolumeUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunVolumeUpdate_Err(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testVolumeVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgVolumeId), testVolumeVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgVolumeName), testVolumeNewVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgVolumeSize), sizeVolumeNew)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgVolumeBus), testVolumeNewVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), false)
		rm.Volume.EXPECT().Update(testVolumeVar, testVolumeVar, volumeProperties).Return(&volumeNew, nil, testVolumeErr)
		err := RunVolumeUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunVolumeUpdate_WaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testVolumeVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgVolumeId), testVolumeVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgVolumeName), testVolumeNewVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgVolumeSize), sizeVolumeNew)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgVolumeBus), testVolumeNewVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), true)
		rm.Volume.EXPECT().Update(testVolumeVar, testVolumeVar, volumeProperties).Return(&volumeNew, nil, nil)
		err := RunVolumeUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunVolumeDelete(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgIgnoreStdin, true)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testVolumeVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgVolumeId), testVolumeVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), false)
		rm.Volume.EXPECT().Delete(testVolumeVar, testVolumeVar).Return(nil, nil)
		err := RunVolumeDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunVolumeDelete_Err(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgIgnoreStdin, true)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testVolumeVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgVolumeId), testVolumeVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), false)
		rm.Volume.EXPECT().Delete(testVolumeVar, testVolumeVar).Return(nil, testVolumeErr)
		err := RunVolumeDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunVolumeDelete_WaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgIgnoreStdin, true)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testVolumeVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgVolumeId), testVolumeVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), true)
		rm.Volume.EXPECT().Delete(testVolumeVar, testVolumeVar).Return(nil, nil)
		err := RunVolumeDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunVolumeDelete_AskForConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgIgnoreStdin, false)
		cfg.Stdin = bytes.NewReader([]byte("YES\n"))
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testVolumeVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgVolumeId), testVolumeVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), false)
		rm.Volume.EXPECT().Delete(testVolumeVar, testVolumeVar).Return(nil, nil)
		err := RunVolumeDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunVolumeDelete_AskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgIgnoreStdin, false)
		cfg.Stdin = os.Stdin
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testVolumeVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgVolumeId), testVolumeVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), false)
		err := RunVolumeDelete(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunAttachGlobalDcIdServerIdValidate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName("volume", config.ArgDataCenterId), testVolumeVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgServerId), testVolumeVar)
		err := PreRunAttachGlobalDcIdServerIdValidate(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunAttachGlobalDcIdServerIdValidate_RequiredFlagsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName("volume", config.ArgDataCenterId), "")
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgServerId), testVolumeVar)
		err := PreRunAttachGlobalDcIdServerIdValidate(cfg)
		assert.Error(t, err)
		assert.True(t, err.Error() == clierror.NewRequiredFlagErr(config.ArgDataCenterId).Error())

		viper.Set(builder.GetGlobalFlagName("volume", config.ArgDataCenterId), testVolumeVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgServerId), "")
		err = PreRunAttachGlobalDcIdServerIdValidate(cfg)
		assert.Error(t, err)
		assert.True(t, err.Error() == clierror.NewRequiredFlagErr(config.ArgServerId).Error())
	})
}

func TestPreRunAttachGlobalDcIdServerVolumeIdsValidate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName("volume", config.ArgDataCenterId), testVolumeVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgServerId), testVolumeVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgVolumeId), testVolumeVar)
		err := PreRunAttachGlobalDcIdServerVolumeIdsValidate(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunAttachGlobalDcIdServerVolumeIdsValidate_RequiredFlagsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName("volume", config.ArgDataCenterId), "")
		err := PreRunAttachGlobalDcIdServerVolumeIdsValidate(cfg)
		assert.Error(t, err)
		assert.True(t, err.Error() == clierror.NewRequiredFlagErr(config.ArgDataCenterId).Error())

		viper.Set(builder.GetGlobalFlagName("volume", config.ArgDataCenterId), testVolumeVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgServerId), "")
		err = PreRunAttachGlobalDcIdServerVolumeIdsValidate(cfg)
		assert.Error(t, err)
		assert.True(t, err.Error() == clierror.NewRequiredFlagErr(config.ArgServerId).Error())

		viper.Set(builder.GetGlobalFlagName("volume", config.ArgDataCenterId), testVolumeVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgServerId), testVolumeVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgVolumeId), "")
		err = PreRunAttachGlobalDcIdServerVolumeIdsValidate(cfg)
		assert.Error(t, err)
		assert.True(t, err.Error() == clierror.NewRequiredFlagErr(config.ArgVolumeId).Error())
	})
}

func TestRunVolumeAttach(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testVolumeVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgVolumeId), testVolumeVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgServerId), testVolumeVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), false)
		rm.Volume.EXPECT().Attach(testVolumeVar, testVolumeVar, testVolumeVar).Return(&resources.Volume{v}, nil, nil)
		err := RunVolumeAttach(cfg)
		assert.NoError(t, err)
	})
}

func TestRunVolumeAttach_Err(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testVolumeVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgVolumeId), testVolumeVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgServerId), testVolumeVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), false)
		rm.Volume.EXPECT().Attach(testVolumeVar, testVolumeVar, testVolumeVar).Return(&resources.Volume{v}, nil, testVolumeErr)
		err := RunVolumeAttach(cfg)
		assert.Error(t, err)
	})
}

func TestRunVolumeAttach_WaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testVolumeVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgVolumeId), testVolumeVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgServerId), testVolumeVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), true)
		rm.Volume.EXPECT().Attach(testVolumeVar, testVolumeVar, testVolumeVar).Return(&resources.Volume{v}, nil, nil)
		err := RunVolumeAttach(cfg)
		assert.Error(t, err)
	})
}

func TestRunVolumeAttachList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName("volume", config.ArgDataCenterId), testVolumeVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgServerId), testVolumeVar)
		rm.Volume.EXPECT().ListAttached(testVolumeVar, testVolumeVar).Return(vsAttached, nil, nil)
		err := RunVolumesAttachList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunVolumeAttachList_Err(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName("volume", config.ArgDataCenterId), testVolumeVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgServerId), testVolumeVar)
		rm.Volume.EXPECT().ListAttached(testVolumeVar, testVolumeVar).Return(vsAttached, nil, testVolumeErr)
		err := RunVolumesAttachList(cfg)
		assert.Error(t, err)
	})
}

func TestRunVolumeAttachGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName("volume", config.ArgDataCenterId), testVolumeVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgServerId), testVolumeVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgVolumeId), testVolumeVar)
		rm.Volume.EXPECT().GetAttached(testVolumeVar, testVolumeVar, testVolumeVar).Return(&resources.Volume{v}, nil, nil)
		err := RunVolumeAttachGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunVolumeAttachGet_Err(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName("volume", config.ArgDataCenterId), testVolumeVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgServerId), testVolumeVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgVolumeId), testVolumeVar)
		rm.Volume.EXPECT().GetAttached(testVolumeVar, testVolumeVar, testVolumeVar).Return(&resources.Volume{v}, nil, testVolumeErr)
		err := RunVolumeAttachGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunVolumeDetach(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgIgnoreStdin, true)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testVolumeVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgServerId), testVolumeVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgVolumeId), testVolumeVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), false)
		rm.Volume.EXPECT().Detach(testVolumeVar, testVolumeVar, testVolumeVar).Return(nil, nil)
		err := RunVolumeDetach(cfg)
		assert.NoError(t, err)
	})
}

func TestRunVolumeDetach_Err(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgIgnoreStdin, true)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testVolumeVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgServerId), testVolumeVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgVolumeId), testVolumeVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), false)
		rm.Volume.EXPECT().Detach(testVolumeVar, testVolumeVar, testVolumeVar).Return(nil, testVolumeErr)
		err := RunVolumeDetach(cfg)
		assert.Error(t, err)
	})
}

func TestRunVolumeDetach_WaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgIgnoreStdin, true)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testVolumeVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgServerId), testVolumeVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgVolumeId), testVolumeVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), true)
		rm.Volume.EXPECT().Detach(testVolumeVar, testVolumeVar, testVolumeVar).Return(nil, nil)
		err := RunVolumeDetach(cfg)
		assert.Error(t, err)
	})
}

func TestRunVolumeDetach_AskForConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgIgnoreStdin, false)
		cfg.Stdin = bytes.NewReader([]byte("YES\n"))
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testVolumeVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgServerId), testVolumeVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgVolumeId), testVolumeVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), false)
		rm.Volume.EXPECT().Detach(testVolumeVar, testVolumeVar, testVolumeVar).Return(nil, nil)
		err := RunVolumeDetach(cfg)
		assert.NoError(t, err)
	})
}

func TestRunVolumeDetach_AskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgIgnoreStdin, false)
		cfg.Stdin = os.Stdin
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testVolumeVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgServerId), testVolumeVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgVolumeId), testVolumeVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), false)
		err := RunVolumeDetach(cfg)
		assert.Error(t, err)
	})
}

func TestVolumesCols(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}

	w := bufio.NewWriter(&b)
	viper.Set(builder.GetGlobalFlagName("volume", config.ArgCols), []string{"Name"})
	getVolumesCols(builder.GetGlobalFlagName("volume", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
}

func TestGetVolumesCols_Err(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}

	w := bufio.NewWriter(&b)
	viper.Set(builder.GetGlobalFlagName("volume", config.ArgCols), []string{"Unknown"})
	getVolumesCols(builder.GetGlobalFlagName("volume", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`unknown column Unknown`)
	assert.True(t, re.Match(b.Bytes()))
}

func TestGetVolumesIds(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}

	w := bufio.NewWriter(&b)
	viper.Set(config.ArgConfig, "../pkg/testdata/config.json")
	getVolumesIds(w, "volume")
	err := w.Flush()
	assert.NoError(t, err)

	re := regexp.MustCompile(`401 Unauthorized`)
	assert.True(t, re.Match(b.Bytes()))
}

func TestGetAttachedVolumesIds(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}

	w := bufio.NewWriter(&b)
	viper.Set(config.ArgConfig, "../pkg/testdata/config.json")
	getAttachedVolumesIds(w, "volume", "volume", "test")
	err := w.Flush()
	assert.NoError(t, err)

	re := regexp.MustCompile(`401 Unauthorized`)
	assert.True(t, re.Match(b.Bytes()))
}
