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
			Image:            &testVolumeVar,
			BackupunitId:     &testVolumeVar,
			SshKeys:          &testVolumeSliceVar,
		},
		Metadata: &ionoscloud.DatacenterElementMetadata{
			State: &testVolumeVar,
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
	testVolumeVar      = "test-volume"
	testVolumeSliceVar = []string{"test-volume"}
	testVolumeNewVar   = "test-new-volume"
	testVolumeErr      = errors.New("volume test: error occurred")
)

func TestPreRunGlobalDcIdVolumeId(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetGlobalFlagName(cfg.Namespace, config.ArgDataCenterId), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgVolumeId), testVolumeVar)
		err := PreRunGlobalDcIdVolumeId(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunGlobalDcIdVolumeIdRequiredFlagsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetGlobalFlagName(cfg.Namespace, config.ArgDataCenterId), "")
		viper.Set(core.GetFlagName(cfg.NS, config.ArgVolumeId), "")
		err := PreRunGlobalDcIdVolumeId(cfg)
		assert.Error(t, err)
	})
}

func TestRunVolumeList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetGlobalFlagName(cfg.Namespace, config.ArgDataCenterId), testVolumeVar)
		rm.Volume.EXPECT().List(testVolumeVar).Return(vs, nil, nil)
		err := RunVolumeList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunVolumeListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetGlobalFlagName(cfg.Namespace, config.ArgDataCenterId), testVolumeVar)
		rm.Volume.EXPECT().List(testVolumeVar).Return(vs, nil, testVolumeErr)
		err := RunVolumeList(cfg)
		assert.Error(t, err)
	})
}

func TestRunVolumeGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetGlobalFlagName(cfg.Namespace, config.ArgDataCenterId), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgVolumeId), testVolumeVar)
		rm.Volume.EXPECT().Get(testVolumeVar, testVolumeVar).Return(&resources.Volume{Volume: v}, nil, nil)
		err := RunVolumeGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunVolumeGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetGlobalFlagName(cfg.Namespace, config.ArgDataCenterId), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgVolumeId), testVolumeVar)
		rm.Volume.EXPECT().Get(testVolumeVar, testVolumeVar).Return(&resources.Volume{Volume: v}, nil, testVolumeErr)
		err := RunVolumeGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunVolumeCreate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetGlobalFlagName(cfg.Namespace, config.ArgDataCenterId), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgVolumeName), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgVolumeSize), sizeVolume)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgVolumeLicenceType), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgVolumeType), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgVolumeBus), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgVolumeZone), zoneVolume)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.Volume.EXPECT().Create(testVolumeVar, testVolumeVar, testVolumeVar, testVolumeVar, testVolumeVar, zoneVolume, sizeVolume).Return(&resources.Volume{Volume: v}, nil, nil)
		err := RunVolumeCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunVolumeCreateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetGlobalFlagName(cfg.Namespace, config.ArgDataCenterId), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgVolumeName), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgVolumeSize), sizeVolume)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgVolumeLicenceType), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgVolumeType), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgVolumeBus), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgVolumeZone), zoneVolume)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.Volume.EXPECT().Create(testVolumeVar, testVolumeVar, testVolumeVar, testVolumeVar, testVolumeVar, zoneVolume, sizeVolume).Return(&resources.Volume{Volume: v}, nil, testVolumeErr)
		err := RunVolumeCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunVolumeCreateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetGlobalFlagName(cfg.Namespace, config.ArgDataCenterId), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgVolumeName), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgVolumeSize), sizeVolume)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgVolumeLicenceType), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgVolumeType), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgVolumeBus), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgVolumeZone), zoneVolume)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		rm.Volume.EXPECT().Create(testVolumeVar, testVolumeVar, testVolumeVar, testVolumeVar, testVolumeVar, zoneVolume, sizeVolume).Return(&resources.Volume{Volume: v}, nil, nil)
		err := RunVolumeCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunVolumeUpdate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetGlobalFlagName(cfg.Namespace, config.ArgDataCenterId), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgVolumeId), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgVolumeName), testVolumeNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgVolumeSize), sizeVolumeNew)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgVolumeBus), testVolumeNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.Volume.EXPECT().Update(testVolumeVar, testVolumeVar, volumeProperties).Return(&volumeNew, nil, nil)
		err := RunVolumeUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunVolumeUpdateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetGlobalFlagName(cfg.Namespace, config.ArgDataCenterId), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgVolumeId), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgVolumeName), testVolumeNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgVolumeSize), sizeVolumeNew)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgVolumeBus), testVolumeNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.Volume.EXPECT().Update(testVolumeVar, testVolumeVar, volumeProperties).Return(&volumeNew, nil, testVolumeErr)
		err := RunVolumeUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunVolumeUpdateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetGlobalFlagName(cfg.Namespace, config.ArgDataCenterId), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgVolumeId), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgVolumeName), testVolumeNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgVolumeSize), sizeVolumeNew)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgVolumeBus), testVolumeNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		rm.Volume.EXPECT().Update(testVolumeVar, testVolumeVar, volumeProperties).Return(&volumeNew, nil, nil)
		err := RunVolumeUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunVolumeDelete(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetGlobalFlagName(cfg.Namespace, config.ArgDataCenterId), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgVolumeId), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.Volume.EXPECT().Delete(testVolumeVar, testVolumeVar).Return(nil, nil)
		err := RunVolumeDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunVolumeDeleteErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetGlobalFlagName(cfg.Namespace, config.ArgDataCenterId), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgVolumeId), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.Volume.EXPECT().Delete(testVolumeVar, testVolumeVar).Return(nil, testVolumeErr)
		err := RunVolumeDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunVolumeDeleteWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetGlobalFlagName(cfg.Namespace, config.ArgDataCenterId), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgVolumeId), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		rm.Volume.EXPECT().Delete(testVolumeVar, testVolumeVar).Return(nil, nil)
		err := RunVolumeDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunVolumeDeleteAskForConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, false)
		cfg.Stdin = bytes.NewReader([]byte("YES\n"))
		viper.Set(core.GetGlobalFlagName(cfg.Namespace, config.ArgDataCenterId), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgVolumeId), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.Volume.EXPECT().Delete(testVolumeVar, testVolumeVar).Return(nil, nil)
		err := RunVolumeDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunVolumeDeleteAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, false)
		cfg.Stdin = os.Stdin
		viper.Set(core.GetGlobalFlagName(cfg.Namespace, config.ArgDataCenterId), testVolumeVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgVolumeId), testVolumeVar)
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

func TestGetAttachedVolumesIds(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() { return }
	w := bufio.NewWriter(&b)
	viper.Set(config.ArgConfig, "../pkg/testdata/config.json")
	getAttachedVolumesIds(w, testVolumeVar, testVolumeVar)
	err := w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`401 Unauthorized`)
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
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), "")
		viper.Set(core.GetFlagName(cfg.NS, config.ArgServerId), "")
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
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgServerId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgVolumeId), testServerVar)
		err := PreRunDcServerVolumeIds(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunDcServerVolumeIdsRequiredFlagsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), "")
		viper.Set(core.GetFlagName(cfg.NS, config.ArgServerId), "")
		viper.Set(core.GetFlagName(cfg.NS, config.ArgVolumeId), "")
		err := PreRunDcServerVolumeIds(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerVolumeAttach(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgVolumeId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgServerId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.Server.EXPECT().AttachVolume(testServerVar, testServerVar, testServerVar).Return(&resources.Volume{Volume: v}, nil, nil)
		err := RunServerVolumeAttach(cfg)
		assert.NoError(t, err)
	})
}

func TestRunServerVolumeAttachErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgVolumeId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgServerId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.Server.EXPECT().AttachVolume(testServerVar, testServerVar, testServerVar).Return(&resources.Volume{Volume: v}, nil, testVolumeErr)
		err := RunServerVolumeAttach(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerVolumeAttachWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgVolumeId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgServerId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		rm.Server.EXPECT().AttachVolume(testServerVar, testServerVar, testServerVar).Return(&resources.Volume{Volume: v}, nil, nil)
		err := RunServerVolumeAttach(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerVolumesList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgServerId), testServerVar)
		rm.Server.EXPECT().ListVolumes(testServerVar, testServerVar).Return(vsAttached, nil, nil)
		err := RunServerVolumesList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunServerVolumesListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgServerId), testServerVar)
		rm.Server.EXPECT().ListVolumes(testServerVar, testServerVar).Return(vsAttached, nil, testVolumeErr)
		err := RunServerVolumesList(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerVolumeDescribe(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgServerId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgVolumeId), testServerVar)
		rm.Server.EXPECT().GetVolume(testServerVar, testServerVar, testServerVar).Return(&resources.Volume{Volume: v}, nil, nil)
		err := RunServerVolumeDescribe(cfg)
		assert.NoError(t, err)
	})
}

func TestRunServerVolumeDescribeErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgServerId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgVolumeId), testServerVar)
		rm.Server.EXPECT().GetVolume(testServerVar, testServerVar, testServerVar).Return(&resources.Volume{Volume: v}, nil, testVolumeErr)
		err := RunServerVolumeDescribe(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerVolumeDetach(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgServerId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgVolumeId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.Server.EXPECT().DetachVolume(testServerVar, testServerVar, testServerVar).Return(nil, nil)
		err := RunServerVolumeDetach(cfg)
		assert.NoError(t, err)
	})
}

func TestRunServerVolumeDetachErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgServerId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgVolumeId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.Server.EXPECT().DetachVolume(testServerVar, testServerVar, testServerVar).Return(&testResponse, nil)
		err := RunServerVolumeDetach(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerVolumeDetachWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgServerId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgVolumeId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		rm.Server.EXPECT().DetachVolume(testServerVar, testServerVar, testServerVar).Return(nil, nil)
		err := RunServerVolumeDetach(cfg)
		assert.Error(t, err)
	})
}

func TestRunVolumeDetachAskForConfirm(t *testing.T) {
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
		viper.Set(core.GetFlagName(cfg.NS, config.ArgVolumeId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.Server.EXPECT().DetachVolume(testServerVar, testServerVar, testServerVar).Return(nil, nil)
		err := RunServerVolumeDetach(cfg)
		assert.NoError(t, err)
	})
}

func TestRunServerVolumeDetachAskForConfirmErr(t *testing.T) {
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
		viper.Set(core.GetFlagName(cfg.NS, config.ArgVolumeId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		err := RunServerVolumeDetach(cfg)
		assert.Error(t, err)
	})
}
