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
	snapshotTest = resources.Snapshot{
		Snapshot: ionoscloud.Snapshot{
			Id: &testSnapshotVar,
			Properties: &ionoscloud.SnapshotProperties{
				Name:        &testSnapshotVar,
				Location:    &testSnapshotVar,
				Description: &testSnapshotVar,
				Size:        &testSnapshotSize,
				LicenceType: &testSnapshotVar,
			},
			Metadata: &ionoscloud.DatacenterElementMetadata{
				State: &snapshotState,
			},
		},
	}
	snapshotState = "BUSY"
	snapshots     = resources.Snapshots{
		Snapshots: ionoscloud.Snapshots{
			Id:    &testSnapshotVar,
			Items: &[]ionoscloud.Snapshot{snapshotTest.Snapshot},
		},
	}
	snapshotProperties = resources.SnapshotProperties{
		SnapshotProperties: ionoscloud.SnapshotProperties{
			Name:                &testSnapshotNewVar,
			Description:         &testSnapshotNewVar,
			CpuHotPlug:          &testSnapshotBoolVar,
			CpuHotUnplug:        &testSnapshotBoolVar,
			RamHotPlug:          &testSnapshotBoolVar,
			RamHotUnplug:        &testSnapshotBoolVar,
			NicHotPlug:          &testSnapshotBoolVar,
			NicHotUnplug:        &testSnapshotBoolVar,
			DiscVirtioHotPlug:   &testSnapshotBoolVar,
			DiscVirtioHotUnplug: &testSnapshotBoolVar,
			DiscScsiHotPlug:     &testSnapshotBoolVar,
			DiscScsiHotUnplug:   &testSnapshotBoolVar,
			SecAuthProtection:   &testSnapshotBoolVar,
			LicenceType:         &testSnapshotVar,
		},
	}
	snapshotNew = resources.Snapshot{
		Snapshot: ionoscloud.Snapshot{
			Id:         &testSnapshotVar,
			Properties: &snapshotProperties.SnapshotProperties,
		},
	}
	testSnapshotBoolVar = false
	testSnapshotSize    = float32(2)
	testSnapshotVar     = "test-snapshot"
	testSnapshotNewVar  = "test-new-snapshot"
	testSnapshotErr     = errors.New("snapshot test error")
)

func TestPreRunSnapshotIdValidate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgSnapshotId), testSnapshotVar)
		err := PreRunSnapshotIdValidate(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunSnapshotIdValidateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgSnapshotId), "")
		err := PreRunSnapshotIdValidate(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunSnapNameLicenceDcIdVolumeIdValidate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgSnapshotName), testSnapshotVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgSnapshotLicenceType), testSnapshotVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgDataCenterId), testSnapshotVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgVolumeId), testSnapshotVar)
		err := PreRunSnapNameLicenceDcIdVolumeIdValidate(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunSnapNameLicenceDcIdVolumeIdValidateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgSnapshotName), "")
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgSnapshotLicenceType), "")
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgDataCenterId), "")
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgVolumeId), "")
		err := PreRunSnapNameLicenceDcIdVolumeIdValidate(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunSnapshotIdDcIdVolumeIdValidate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgSnapshotId), testSnapshotVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgDataCenterId), testSnapshotVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgVolumeId), testSnapshotVar)
		err := PreRunSnapshotIdDcIdVolumeIdValidate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunSnapshotList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		rm.Snapshot.EXPECT().List().Return(snapshots, nil, nil)
		err := RunSnapshotList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunSnapshotListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		rm.Snapshot.EXPECT().List().Return(snapshots, nil, testSnapshotErr)
		err := RunSnapshotList(cfg)
		assert.Error(t, err)
	})
}

func TestRunSnapshotListSort(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgSnapshotLicenceType), testSnapshotVar)
		rm.Snapshot.EXPECT().List().Return(snapshots, nil, nil)
		err := RunSnapshotList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunSnapshotGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgSnapshotId), testSnapshotVar)
		rm.Snapshot.EXPECT().Get(testSnapshotVar).Return(&snapshotTest, nil, nil)
		err := RunSnapshotGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunSnapshotGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgSnapshotId), testSnapshotVar)
		rm.Snapshot.EXPECT().Get(testSnapshotVar).Return(&snapshotTest, nil, testSnapshotErr)
		err := RunSnapshotGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunSnapshotCreate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgSnapshotName), testSnapshotVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgDataCenterId), testSnapshotVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgVolumeId), testSnapshotVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgSnapshotDescription), testSnapshotVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgSnapshotLicenceType), testSnapshotVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgSnapshotName), testSnapshotVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgSnapshotSecAuthProtection), false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), false)
		rm.Snapshot.EXPECT().Create(testSnapshotVar, testSnapshotVar, testSnapshotVar, testSnapshotVar, testSnapshotVar, false).Return(&snapshotTest, nil, nil)
		err := RunSnapshotCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunSnapshotCreateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgSnapshotName), testSnapshotVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgDataCenterId), testSnapshotVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgVolumeId), testSnapshotVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgSnapshotDescription), testSnapshotVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgSnapshotLicenceType), testSnapshotVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgSnapshotName), testSnapshotVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgSnapshotSecAuthProtection), false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), false)
		rm.Snapshot.EXPECT().Create(testSnapshotVar, testSnapshotVar, testSnapshotVar, testSnapshotVar, testSnapshotVar, false).Return(&snapshotTest, &testResponse, nil)
		err := RunSnapshotCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunSnapshotUpdate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgSnapshotId), testSnapshotVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgSnapshotDescription), testSnapshotNewVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgSnapshotName), testSnapshotNewVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgSnapshotCpuHotPlug), testSnapshotBoolVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgSnapshotCpuHotUnplug), testSnapshotBoolVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgSnapshotRamHotPlug), testSnapshotBoolVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgSnapshotRamHotUnplug), testSnapshotBoolVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgSnapshotNicHotPlug), testSnapshotBoolVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgSnapshotNicHotUnplug), testSnapshotBoolVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgSnapshotDiscVirtioHotPlug), testSnapshotBoolVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgSnapshotDiscVirtioHotUnplug), testSnapshotBoolVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgSnapshotDiscScsiHotPlug), testSnapshotBoolVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgSnapshotDiscScsiHotUnplug), testSnapshotBoolVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgSnapshotSecAuthProtection), testSnapshotBoolVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgSnapshotLicenceType), testSnapshotVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), false)
		rm.Snapshot.EXPECT().Update(testSnapshotVar, snapshotProperties).Return(&snapshotNew, nil, nil)
		err := RunSnapshotUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunSnapshotUpdateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgSnapshotId), testSnapshotVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgSnapshotDescription), testSnapshotNewVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgSnapshotName), testSnapshotNewVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgSnapshotCpuHotPlug), testSnapshotBoolVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgSnapshotCpuHotUnplug), testSnapshotBoolVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgSnapshotRamHotPlug), testSnapshotBoolVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgSnapshotRamHotUnplug), testSnapshotBoolVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgSnapshotNicHotPlug), testSnapshotBoolVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgSnapshotNicHotUnplug), testSnapshotBoolVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgSnapshotDiscVirtioHotPlug), testSnapshotBoolVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgSnapshotDiscVirtioHotUnplug), testSnapshotBoolVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgSnapshotDiscScsiHotPlug), testSnapshotBoolVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgSnapshotDiscScsiHotUnplug), testSnapshotBoolVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgSnapshotSecAuthProtection), testSnapshotBoolVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgSnapshotLicenceType), testSnapshotVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), true)
		rm.Snapshot.EXPECT().Update(testSnapshotVar, snapshotProperties).Return(&snapshotNew, nil, testSnapshotErr)
		err := RunSnapshotUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunSnapshotRestore(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgSnapshotId), testSnapshotVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgVolumeId), testSnapshotVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgDataCenterId), testSnapshotVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), false)
		rm.Snapshot.EXPECT().Restore(testSnapshotVar, testSnapshotVar, testSnapshotVar).Return(nil, nil)
		err := RunSnapshotRestore(cfg)
		assert.NoError(t, err)
	})
}

func TestRunSnapshotRestoreErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgSnapshotId), testSnapshotVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgVolumeId), testSnapshotVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgDataCenterId), testSnapshotVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), false)
		rm.Snapshot.EXPECT().Restore(testSnapshotVar, testSnapshotVar, testSnapshotVar).Return(nil, testSnapshotErr)
		err := RunSnapshotRestore(cfg)
		assert.Error(t, err)
	})
}

func TestRunSnapshotRestoreAskForConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgSnapshotId), testSnapshotVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgVolumeId), testSnapshotVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgDataCenterId), testSnapshotVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), false)
		rm.Snapshot.EXPECT().Restore(testSnapshotVar, testSnapshotVar, testSnapshotVar).Return(nil, nil)
		cfg.Stdin = bytes.NewReader([]byte("YES\n"))
		err := RunSnapshotRestore(cfg)
		assert.NoError(t, err)
	})
}

func TestRunSnapshotDelete(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgSnapshotId), testSnapshotVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), false)
		rm.Snapshot.EXPECT().Delete(testSnapshotVar).Return(nil, nil)
		err := RunSnapshotDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunSnapshotDeleteErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgSnapshotId), testSnapshotVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), false)
		rm.Snapshot.EXPECT().Delete(testSnapshotVar).Return(nil, testSnapshotErr)
		err := RunSnapshotDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunSnapshotDeleteAskForConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgSnapshotId), testSnapshotVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), false)
		cfg.Stdin = bytes.NewReader([]byte("YES\n"))
		rm.Snapshot.EXPECT().Delete(testSnapshotVar).Return(nil, nil)
		err := RunSnapshotDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunSnapshotDeleteAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgSnapshotId), testSnapshotVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), false)
		cfg.Stdin = os.Stdin
		err := RunSnapshotDelete(cfg)
		assert.Error(t, err)
	})
}

func TestGetSnapshotsCols(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	viper.Set(builder.GetGlobalFlagName("snapshot", config.ArgCols), []string{"Name"})
	getSnapshotCols(builder.GetGlobalFlagName("snapshot", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
}

func TestGetSnapshotsColsErr(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	viper.Set(builder.GetGlobalFlagName("snapshot", config.ArgCols), []string{"Unknown"})
	getSnapshotCols(builder.GetGlobalFlagName("snapshot", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`unknown column Unknown`)
	assert.True(t, re.Match(b.Bytes()))
}

func TestGetSnapshotsIds(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	viper.Set(config.ArgConfig, "../pkg/testdata/config.json")
	getSnapshotIds(w)
	err := w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`401 Unauthorized`)
	assert.True(t, re.Match(b.Bytes()))
}
