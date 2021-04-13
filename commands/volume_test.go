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
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testVolumeVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgVolumeId), testVolumeVar)
		err := PreRunGlobalDcIdVolumeIdValidate(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunGlobalDcIdVolumeIdValidateRequiredFlagsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
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

func TestRunVolumeList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testVolumeVar)
		rm.Volume.EXPECT().List(testVolumeVar).Return(vs, nil, nil)
		err := RunVolumeList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunVolumeListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
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
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testVolumeVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgVolumeId), testVolumeVar)
		rm.Volume.EXPECT().Get(testVolumeVar, testVolumeVar).Return(&resources.Volume{v}, nil, nil)
		err := RunVolumeGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunVolumeGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
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
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
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

func TestRunVolumeCreateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
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

func TestRunVolumeCreateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
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
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
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

func TestRunVolumeUpdateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
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

func TestRunVolumeUpdateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
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
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
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

func TestRunVolumeDeleteErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
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

func TestRunVolumeDeleteWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
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

func TestRunVolumeDeleteAskForConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
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

func TestRunVolumeDeleteAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
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

func TestVolumesCols(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() { return }
	w := bufio.NewWriter(&b)
	viper.Set(builder.GetGlobalFlagName("volume", config.ArgCols), []string{"Name"})
	getVolumesCols(builder.GetGlobalFlagName("volume", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
}

func TestGetVolumesColsErr(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() { return }
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
	clierror.ErrAction = func() { return }
	w := bufio.NewWriter(&b)
	viper.Set(config.ArgConfig, "../pkg/testdata/config.json")
	getVolumesIds(w, testVolumeVar)
	err := w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`401 Unauthorized`)
	assert.True(t, re.Match(b.Bytes()))
}
