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
	snapshotsList = resources.Snapshots{
		Snapshots: ionoscloud.Snapshots{
			Id: &testSnapshotVar,
			Items: &[]ionoscloud.Snapshot{
				snapshotTest.Snapshot,
				snapshotTest.Snapshot,
			},
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

func TestSnapshotCmd(t *testing.T) {
	var err error
	core.RootCmdTest.AddCommand(SnapshotCmd())
	if ok := SnapshotCmd().IsAvailableCommand(); !ok {
		err = errors.New("non-available cmd")
	}
	assert.NoError(t, err)
}

func TestPreRunSnapshotList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		err := PreRunSnapshotList(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunSnapshotListFilters(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFilters), []string{fmt.Sprintf("createdBy=%s", testQueryParamVar)})
		err := PreRunSnapshotList(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunSnapshotListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFilters), []string{fmt.Sprintf("%s=%s", testQueryParamVar, testQueryParamVar)})
		err := PreRunSnapshotList(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunSnapshotId(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgSnapshotId), testSnapshotVar)
		err := PreRunSnapshotId(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunSnapshotIdErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		err := PreRunSnapshotId(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunSnapshotIdDcIdVolumeId(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgSnapshotId), testSnapshotVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testSnapshotVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgVolumeId), testSnapshotVar)
		err := PreRunSnapshotIdDcIdVolumeId(cfg)
		assert.NoError(t, err)
	})
}

func TestRunSnapshotList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		rm.CloudApiV6Mocks.Snapshot.EXPECT().List(resources.ListQueryParams{}).Return(snapshots, &testResponse, nil)
		err := RunSnapshotList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunSnapshotListQueryParams(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFilters), []string{fmt.Sprintf("%s=%s", testQueryParamVar, testQueryParamVar)})
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgOrderBy), testQueryParamVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgMaxResults), testMaxResultsVar)
		rm.CloudApiV6Mocks.Snapshot.EXPECT().List(testListQueryParam).Return(resources.Snapshots{}, &testResponse, nil)
		err := RunSnapshotList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunSnapshotListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		rm.CloudApiV6Mocks.Snapshot.EXPECT().List(resources.ListQueryParams{}).Return(snapshots, nil, testSnapshotErr)
		err := RunSnapshotList(cfg)
		assert.Error(t, err)
	})
}

func TestRunSnapshotListSort(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLicenceType), testSnapshotVar)
		rm.CloudApiV6Mocks.Snapshot.EXPECT().List(resources.ListQueryParams{}).Return(snapshots, nil, nil)
		err := RunSnapshotList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunSnapshotGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgSnapshotId), testSnapshotVar)
		rm.CloudApiV6Mocks.Snapshot.EXPECT().Get(testSnapshotVar).Return(&snapshotTest, &testResponse, nil)
		err := RunSnapshotGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunSnapshotGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgSnapshotId), testSnapshotVar)
		rm.CloudApiV6Mocks.Snapshot.EXPECT().Get(testSnapshotVar).Return(&snapshotTest, nil, testSnapshotErr)
		err := RunSnapshotGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunSnapshotCreate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testSnapshotVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testSnapshotVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgVolumeId), testSnapshotVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDescription), testSnapshotVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLicenceType), testSnapshotVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testSnapshotVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgSecAuthProtection), false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.Snapshot.EXPECT().Create(testSnapshotVar, testSnapshotVar, testSnapshotVar, testSnapshotVar, testSnapshotVar, false).Return(&snapshotTest, &testResponse, nil)
		err := RunSnapshotCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunSnapshotCreateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testSnapshotVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testSnapshotVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgVolumeId), testSnapshotVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDescription), testSnapshotVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLicenceType), testSnapshotVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testSnapshotVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgSecAuthProtection), false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.Snapshot.EXPECT().Create(testSnapshotVar, testSnapshotVar, testSnapshotVar, testSnapshotVar, testSnapshotVar, false).Return(&snapshotTest, &testResponse, testSnapshotErr)
		err := RunSnapshotCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunSnapshotUpdate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgSnapshotId), testSnapshotVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDescription), testSnapshotNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testSnapshotNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCpuHotPlug), testSnapshotBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCpuHotUnplug), testSnapshotBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRamHotPlug), testSnapshotBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRamHotUnplug), testSnapshotBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicHotPlug), testSnapshotBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicHotUnplug), testSnapshotBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDiscVirtioHotPlug), testSnapshotBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDiscVirtioHotUnplug), testSnapshotBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDiscScsiHotPlug), testSnapshotBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDiscScsiHotUnplug), testSnapshotBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgSecAuthProtection), testSnapshotBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLicenceType), testSnapshotVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.Snapshot.EXPECT().Update(testSnapshotVar, snapshotProperties).Return(&snapshotNew, &testResponse, nil)
		err := RunSnapshotUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunSnapshotUpdateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgSnapshotId), testSnapshotVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDescription), testSnapshotNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testSnapshotNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCpuHotPlug), testSnapshotBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCpuHotUnplug), testSnapshotBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRamHotPlug), testSnapshotBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRamHotUnplug), testSnapshotBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicHotPlug), testSnapshotBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicHotUnplug), testSnapshotBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDiscVirtioHotPlug), testSnapshotBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDiscVirtioHotUnplug), testSnapshotBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDiscScsiHotPlug), testSnapshotBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDiscScsiHotUnplug), testSnapshotBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgSecAuthProtection), testSnapshotBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLicenceType), testSnapshotVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		rm.CloudApiV6Mocks.Snapshot.EXPECT().Update(testSnapshotVar, snapshotProperties).Return(&snapshotNew, nil, testSnapshotErr)
		err := RunSnapshotUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunSnapshotRestore(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(config.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgSnapshotId), testSnapshotVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgVolumeId), testSnapshotVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testSnapshotVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.Snapshot.EXPECT().Restore(testSnapshotVar, testSnapshotVar, testSnapshotVar).Return(&testResponse, nil)
		err := RunSnapshotRestore(cfg)
		assert.NoError(t, err)
	})
}

func TestRunSnapshotRestoreErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgSnapshotId), testSnapshotVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgVolumeId), testSnapshotVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testSnapshotVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.Snapshot.EXPECT().Restore(testSnapshotVar, testSnapshotVar, testSnapshotVar).Return(nil, testSnapshotErr)
		err := RunSnapshotRestore(cfg)
		assert.Error(t, err)
	})
}

func TestRunSnapshotRestoreAskForConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgSnapshotId), testSnapshotVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgVolumeId), testSnapshotVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testSnapshotVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.Snapshot.EXPECT().Restore(testSnapshotVar, testSnapshotVar, testSnapshotVar).Return(nil, nil)
		cfg.Stdin = bytes.NewReader([]byte("YES\n"))
		err := RunSnapshotRestore(cfg)
		assert.NoError(t, err)
	})
}

func TestRunSnapshotDelete(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(config.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgSnapshotId), testSnapshotVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.Snapshot.EXPECT().Delete(testSnapshotVar).Return(&testResponse, nil)
		err := RunSnapshotDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunSnapshotDeleteAll(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.Snapshot.EXPECT().List(resources.ListQueryParams{}).Return(snapshotsList, &testResponse, nil)
		rm.CloudApiV6Mocks.Snapshot.EXPECT().Delete(testSnapshotVar).Return(&testResponse, nil)
		rm.CloudApiV6Mocks.Snapshot.EXPECT().Delete(testSnapshotVar).Return(&testResponse, nil)
		err := RunSnapshotDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunSnapshotDeleteAllListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.Snapshot.EXPECT().List(resources.ListQueryParams{}).Return(snapshotsList, nil, testSnapshotErr)
		err := RunSnapshotDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunSnapshotDeleteAllItemsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.Snapshot.EXPECT().List(resources.ListQueryParams{}).Return(
			resources.Snapshots{}, &testResponse, nil)
		err := RunSnapshotDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunSnapshotDeleteAllLenErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.Snapshot.EXPECT().List(resources.ListQueryParams{}).Return(
			resources.Snapshots{Snapshots: ionoscloud.Snapshots{Items: &[]ionoscloud.Snapshot{}}}, &testResponse, nil)
		err := RunSnapshotDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunSnapshotDeleteAllErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.Snapshot.EXPECT().List(resources.ListQueryParams{}).Return(snapshotsList, &testResponse, nil)
		rm.CloudApiV6Mocks.Snapshot.EXPECT().Delete(testSnapshotVar).Return(&testResponse, testSnapshotErr)
		rm.CloudApiV6Mocks.Snapshot.EXPECT().Delete(testSnapshotVar).Return(&testResponse, nil)
		err := RunSnapshotDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunSnapshotDeleteErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgSnapshotId), testSnapshotVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.Snapshot.EXPECT().Delete(testSnapshotVar).Return(nil, testSnapshotErr)
		err := RunSnapshotDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunSnapshotDeleteAskForConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgSnapshotId), testSnapshotVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		cfg.Stdin = bytes.NewReader([]byte("YES\n"))
		rm.CloudApiV6Mocks.Snapshot.EXPECT().Delete(testSnapshotVar).Return(nil, nil)
		err := RunSnapshotDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunSnapshotDeleteAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgSnapshotId), testSnapshotVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
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
	viper.Set(core.GetGlobalFlagName("snapshot", config.ArgCols), []string{"Name"})
	getSnapshotCols(core.GetGlobalFlagName("snapshot", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
}

func TestGetSnapshotsColsErr(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	viper.Set(core.GetGlobalFlagName("snapshot", config.ArgCols), []string{"Unknown"})
	getSnapshotCols(core.GetGlobalFlagName("snapshot", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`unknown column Unknown`)
	assert.True(t, re.Match(b.Bytes()))
}
