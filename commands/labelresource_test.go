package commands

import (
	"bufio"
	"bytes"
	"errors"
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
	testLabelResource = ionoscloud.LabelResource{
		Id: &testLabelVar,
		Properties: &ionoscloud.LabelResourceProperties{
			Key:   &testLabelResourceVar,
			Value: &testLabelResourceVar,
		},
	}
	testLabelResources = resources.LabelResources{
		LabelResources: ionoscloud.LabelResources{
			Id:    &testLabelVar,
			Items: &[]ionoscloud.LabelResource{testLabelResource},
		},
	}
	testLabelResourceRes = resources.LabelResource{LabelResource: testLabelResource}
	testLabelResourceVar = "test-label-resource"
	testLabelResourceErr = errors.New("label resource test error")
)

func TestPreRunDcIdLabelKey(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgDataCenterId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelKey), testLabelResourceVar)
		viper.Set(config.ArgQuiet, false)
		err := PreRunDcIdLabelKey(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunDcIdLabelKeyErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgDataCenterId), "")
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelKey), "")
		viper.Set(config.ArgQuiet, false)
		err := PreRunDcIdLabelKey(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunInheritedDcIdVolume(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetGlobalFlagName(cfg.GrandParentName, config.ArgDataCenterId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgVolumeId), testLabelResourceVar)
		viper.Set(config.ArgQuiet, false)
		err := PreRunInheritedDcIdVolume(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunInheritedDcIdVolumeErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetGlobalFlagName(cfg.GrandParentName, config.ArgDataCenterId), "")
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgVolumeId), "")
		viper.Set(config.ArgQuiet, false)
		err := PreRunInheritedDcIdVolume(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunDcIdLabelKeyValue(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgDataCenterId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelKey), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelValue), testLabelResourceVar)
		viper.Set(config.ArgQuiet, false)
		err := PreRunDcIdLabelKeyValue(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunDcIdLabelKeyValueErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgDataCenterId), "")
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelKey), "")
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelValue), "")
		viper.Set(config.ArgQuiet, false)
		err := PreRunDcIdLabelKeyValue(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunIpBlockIdLabelKey(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgIpBlockId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelKey), testLabelResourceVar)
		viper.Set(config.ArgQuiet, false)
		err := PreRunIpBlockIdLabelKey(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunIpBlockIdLabelKeyErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgIpBlockId), "")
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelKey), "")
		viper.Set(config.ArgQuiet, false)
		err := PreRunIpBlockIdLabelKey(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunIpBlockIdLabelKeyValue(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgIpBlockId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelKey), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelValue), testLabelResourceVar)
		viper.Set(config.ArgQuiet, false)
		err := PreRunIpBlockIdLabelKeyValue(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunIpBlockIdLabelKeyValueErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgIpBlockId), "")
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelKey), "")
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelValue), "")
		viper.Set(config.ArgQuiet, false)
		err := PreRunIpBlockIdLabelKeyValue(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunSnapshotIdLabelKey(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgSnapshotId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelKey), testLabelResourceVar)
		viper.Set(config.ArgQuiet, false)
		err := PreRunSnapshotIdLabelKey(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunSnapshotIdLabelKeyErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgSnapshotId), "")
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelKey), "")
		viper.Set(config.ArgQuiet, false)
		err := PreRunSnapshotIdLabelKey(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunSnapshotIdLabelKeyValue(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgSnapshotId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelKey), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelValue), testLabelResourceVar)
		viper.Set(config.ArgQuiet, false)
		err := PreRunSnapshotIdLabelKeyValue(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunSnapshotIdLabelKeyValueErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgSnapshotId), "")
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelKey), "")
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelValue), "")
		viper.Set(config.ArgQuiet, false)
		err := PreRunSnapshotIdLabelKeyValue(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunDcServerIdsLabelKey(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgDataCenterId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgServerId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelKey), testLabelResourceVar)
		viper.Set(config.ArgQuiet, false)
		err := PreRunDcServerIdsLabelKey(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunDcServerIdsLabelKeyErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgDataCenterId), "")
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgServerId), "")
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelKey), "")
		err := PreRunDcServerIdsLabelKey(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunDcServerIdsLabelKeyGlobalErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgDataCenterId), "")
		err := PreRunDcServerIdsLabelKey(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunDcServerIdsLabelKeyValueGlobalErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgDataCenterId), "")
		err := PreRunDcServerIdsLabelKeyValue(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunDcServerIdsLabelKeyValue(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgDataCenterId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgServerId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelKey), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelValue), testLabelResourceVar)
		viper.Set(config.ArgQuiet, false)
		err := PreRunDcServerIdsLabelKeyValue(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunDcServerIdsLabelKeyValueErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgDataCenterId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgServerId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelKey), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelValue), "")
		viper.Set(config.ArgQuiet, false)
		err := PreRunDcServerIdsLabelKeyValue(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunInheritedDcIdVolumeLabelKey(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetGlobalFlagName(cfg.GrandParentName, config.ArgDataCenterId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgVolumeId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelKey), testLabelResourceVar)
		viper.Set(config.ArgQuiet, false)
		err := PreRunInheritedDcIdVolumeLabelKey(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunInheritedDcIdVolumeLabelKeyErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetGlobalFlagName(cfg.GrandParentName, config.ArgDataCenterId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgVolumeId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelKey), "")
		viper.Set(config.ArgQuiet, false)
		err := PreRunInheritedDcIdVolumeLabelKey(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunInheritedDcIdVolumeLabelKeyGlobalErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetGlobalFlagName(cfg.GrandParentName, config.ArgDataCenterId), "")
		viper.Set(config.ArgQuiet, false)
		err := PreRunInheritedDcIdVolumeLabelKey(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunInheritedDcIdVolumeLabelKeyValue(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetGlobalFlagName(cfg.GrandParentName, config.ArgDataCenterId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgVolumeId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelKey), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelValue), testLabelResourceVar)
		viper.Set(config.ArgQuiet, false)
		err := PreRunInheritedDcIdVolumeLabelKeyValue(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunInheritedDcIdVolumeLabelKeyValueErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetGlobalFlagName(cfg.GrandParentName, config.ArgDataCenterId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgVolumeId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelKey), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelValue), "")
		viper.Set(config.ArgQuiet, false)
		err := PreRunInheritedDcIdVolumeLabelKeyValue(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunInheritedDcIdVolumeLabelKeyValueGlobalErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetGlobalFlagName(cfg.GrandParentName, config.ArgDataCenterId), "")
		viper.Set(config.ArgQuiet, false)
		err := PreRunInheritedDcIdVolumeLabelKeyValue(cfg)
		assert.Error(t, err)
	})
}

func TestRunDataCenterLabelsList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgDataCenterId), testLabelResourceVar)
		rm.Label.EXPECT().DatacenterList(testLabelResourceVar).Return(testLabelResources, nil, nil)
		err := RunDataCenterLabelsList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunDataCenterLabelsListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgDataCenterId), testLabelResourceVar)
		rm.Label.EXPECT().DatacenterList(testLabelResourceVar).Return(testLabelResources, nil, testLabelResourceErr)
		err := RunDataCenterLabelsList(cfg)
		assert.Error(t, err)
	})
}

func TestRunDataCenterLabelGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgDataCenterId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelKey), testLabelResourceVar)
		rm.Label.EXPECT().DatacenterGet(testLabelResourceVar, testLabelResourceVar).Return(&testLabelResourceRes, nil, nil)
		err := RunDataCenterLabelGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunDataCenterLabelGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgDataCenterId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelKey), testLabelResourceVar)
		rm.Label.EXPECT().DatacenterGet(testLabelResourceVar, testLabelResourceVar).Return(&testLabelResourceRes, nil, testLabelResourceErr)
		err := RunDataCenterLabelGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunDatacenterLabelAdd(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgDataCenterId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelKey), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelValue), testLabelResourceVar)
		rm.Label.EXPECT().DatacenterCreate(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar).Return(&testLabelResourceRes, nil, nil)
		err := RunDataCenterLabelAdd(cfg)
		assert.NoError(t, err)
	})
}

func TestRunDatacenterLabelAddErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgDataCenterId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelKey), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelValue), testLabelResourceVar)
		rm.Label.EXPECT().DatacenterCreate(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar).Return(&testLabelResourceRes, nil, testLabelResourceErr)
		err := RunDataCenterLabelAdd(cfg)
		assert.Error(t, err)
	})
}

func TestRunDatacenterLabelRemove(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgDataCenterId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelKey), testLabelResourceVar)
		rm.Label.EXPECT().DatacenterDelete(testLabelResourceVar, testLabelResourceVar).Return(nil, nil)
		err := RunDataCenterLabelRemove(cfg)
		assert.NoError(t, err)
	})
}

func TestRunDatacenterLabelRemoveErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgDataCenterId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelKey), testLabelResourceVar)
		rm.Label.EXPECT().DatacenterDelete(testLabelResourceVar, testLabelResourceVar).Return(nil, testLabelResourceErr)
		err := RunDataCenterLabelRemove(cfg)
		assert.Error(t, err)
	})
}

func TestRunIpBlockLabelsList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgIpBlockId), testLabelResourceVar)
		rm.Label.EXPECT().IpBlockList(testLabelResourceVar).Return(testLabelResources, nil, nil)
		err := RunIpBlockLabelsList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunIpBlockLabelsListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgIpBlockId), testLabelResourceVar)
		rm.Label.EXPECT().IpBlockList(testLabelResourceVar).Return(testLabelResources, nil, testLabelResourceErr)
		err := RunIpBlockLabelsList(cfg)
		assert.Error(t, err)
	})
}

func TestRunIpBlockLabelGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgIpBlockId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelKey), testLabelResourceVar)
		rm.Label.EXPECT().IpBlockGet(testLabelResourceVar, testLabelResourceVar).Return(&testLabelResourceRes, nil, nil)
		err := RunIpBlockLabelGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunIpBlockLabelGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgIpBlockId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelKey), testLabelResourceVar)
		rm.Label.EXPECT().IpBlockGet(testLabelResourceVar, testLabelResourceVar).Return(&testLabelResourceRes, nil, testLabelResourceErr)
		err := RunIpBlockLabelGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunIpBlockLabelAdd(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgIpBlockId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelKey), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelValue), testLabelResourceVar)
		rm.Label.EXPECT().IpBlockCreate(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar).Return(&testLabelResourceRes, nil, nil)
		err := RunIpBlockLabelAdd(cfg)
		assert.NoError(t, err)
	})
}

func TestRunIpBlockLabelAddErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgIpBlockId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelKey), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelValue), testLabelResourceVar)
		rm.Label.EXPECT().IpBlockCreate(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar).Return(&testLabelResourceRes, nil, testLabelResourceErr)
		err := RunIpBlockLabelAdd(cfg)
		assert.Error(t, err)
	})
}

func TestRunIpBlockLabelRemove(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgIpBlockId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelKey), testLabelResourceVar)
		rm.Label.EXPECT().IpBlockDelete(testLabelResourceVar, testLabelResourceVar).Return(nil, nil)
		err := RunIpBlockLabelRemove(cfg)
		assert.NoError(t, err)
	})
}

func TestRunIpBlockLabelRemoveErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgIpBlockId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelKey), testLabelResourceVar)
		rm.Label.EXPECT().IpBlockDelete(testLabelResourceVar, testLabelResourceVar).Return(nil, testLabelResourceErr)
		err := RunIpBlockLabelRemove(cfg)
		assert.Error(t, err)
	})
}

func TestRunSnapshotLabelsList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgSnapshotId), testLabelResourceVar)
		rm.Label.EXPECT().SnapshotList(testLabelResourceVar).Return(testLabelResources, nil, nil)
		err := RunSnapshotLabelsList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunSnapshotLabelsListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgSnapshotId), testLabelResourceVar)
		rm.Label.EXPECT().SnapshotList(testLabelResourceVar).Return(testLabelResources, nil, testLabelResourceErr)
		err := RunSnapshotLabelsList(cfg)
		assert.Error(t, err)
	})
}

func TestRunSnapshotLabelGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgSnapshotId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelKey), testLabelResourceVar)
		rm.Label.EXPECT().SnapshotGet(testLabelResourceVar, testLabelResourceVar).Return(&testLabelResourceRes, nil, nil)
		err := RunSnapshotLabelGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunSnapshotLabelGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgSnapshotId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelKey), testLabelResourceVar)
		rm.Label.EXPECT().SnapshotGet(testLabelResourceVar, testLabelResourceVar).Return(&testLabelResourceRes, nil, testLabelResourceErr)
		err := RunSnapshotLabelGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunSnapshotLabelAdd(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgSnapshotId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelKey), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelValue), testLabelResourceVar)
		rm.Label.EXPECT().SnapshotCreate(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar).Return(&testLabelResourceRes, nil, nil)
		err := RunSnapshotLabelAdd(cfg)
		assert.NoError(t, err)
	})
}

func TestRunSnapshotLabelAddErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgSnapshotId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelKey), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelValue), testLabelResourceVar)
		rm.Label.EXPECT().SnapshotCreate(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar).Return(&testLabelResourceRes, nil, testLabelResourceErr)
		err := RunSnapshotLabelAdd(cfg)
		assert.Error(t, err)
	})
}

func TestRunSnapshotLabelRemove(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgSnapshotId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelKey), testLabelResourceVar)
		rm.Label.EXPECT().SnapshotDelete(testLabelResourceVar, testLabelResourceVar).Return(nil, nil)
		err := RunSnapshotLabelRemove(cfg)
		assert.NoError(t, err)
	})
}

func TestRunSnapshotLabelRemoveErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgSnapshotId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelKey), testLabelResourceVar)
		rm.Label.EXPECT().SnapshotDelete(testLabelResourceVar, testLabelResourceVar).Return(nil, testLabelResourceErr)
		err := RunSnapshotLabelRemove(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerLabelsList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgDataCenterId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgServerId), testLabelResourceVar)
		rm.Label.EXPECT().ServerList(testLabelResourceVar, testLabelResourceVar).Return(testLabelResources, nil, nil)
		err := RunServerLabelsList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunServerLabelsListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgDataCenterId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgServerId), testLabelResourceVar)
		rm.Label.EXPECT().ServerList(testLabelResourceVar, testLabelResourceVar).Return(testLabelResources, nil, testLabelResourceErr)
		err := RunServerLabelsList(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerLabelGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgDataCenterId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgServerId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelKey), testLabelResourceVar)
		rm.Label.EXPECT().ServerGet(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar).Return(&testLabelResourceRes, nil, nil)
		err := RunServerLabelGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunServerLabelGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgDataCenterId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgServerId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelKey), testLabelResourceVar)
		rm.Label.EXPECT().ServerGet(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar).Return(&testLabelResourceRes, nil, testLabelResourceErr)
		err := RunServerLabelGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerLabelAdd(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgDataCenterId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgServerId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelKey), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelValue), testLabelResourceVar)
		rm.Label.EXPECT().ServerCreate(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar, testLabelResourceVar).Return(&testLabelResourceRes, nil, nil)
		err := RunServerLabelAdd(cfg)
		assert.NoError(t, err)
	})
}

func TestRunServerLabelAddErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgDataCenterId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgServerId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelKey), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelValue), testLabelResourceVar)
		rm.Label.EXPECT().ServerCreate(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar, testLabelResourceVar).Return(&testLabelResourceRes, nil, testLabelResourceErr)
		err := RunServerLabelAdd(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerLabelRemove(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgDataCenterId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgServerId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelKey), testLabelResourceVar)
		rm.Label.EXPECT().ServerDelete(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar).Return(nil, nil)
		err := RunServerLabelRemove(cfg)
		assert.NoError(t, err)
	})
}

func TestRunServerLabelRemoveErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgDataCenterId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgServerId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelKey), testLabelResourceVar)
		rm.Label.EXPECT().ServerDelete(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar).Return(nil, testLabelResourceErr)
		err := RunServerLabelRemove(cfg)
		assert.Error(t, err)
	})
}

func TestRunVolumeLabelsList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.GrandParentName, config.ArgDataCenterId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgVolumeId), testLabelResourceVar)
		rm.Label.EXPECT().VolumeList(testLabelResourceVar, testLabelResourceVar).Return(testLabelResources, nil, nil)
		err := RunVolumeLabelsList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunVolumeLabelsListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.GrandParentName, config.ArgDataCenterId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgVolumeId), testLabelResourceVar)
		rm.Label.EXPECT().VolumeList(testLabelResourceVar, testLabelResourceVar).Return(testLabelResources, nil, testLabelResourceErr)
		err := RunVolumeLabelsList(cfg)
		assert.Error(t, err)
	})
}

func TestRunVolumeLabelGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.GrandParentName, config.ArgDataCenterId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgVolumeId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelKey), testLabelResourceVar)
		rm.Label.EXPECT().VolumeGet(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar).Return(&testLabelResourceRes, nil, nil)
		err := RunVolumeLabelGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunVolumeLabelGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.GrandParentName, config.ArgDataCenterId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgVolumeId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelKey), testLabelResourceVar)
		rm.Label.EXPECT().VolumeGet(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar).Return(&testLabelResourceRes, nil, testLabelResourceErr)
		err := RunVolumeLabelGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunVolumeLabelAdd(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.GrandParentName, config.ArgDataCenterId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgVolumeId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelKey), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelValue), testLabelResourceVar)
		rm.Label.EXPECT().VolumeCreate(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar, testLabelResourceVar).Return(&testLabelResourceRes, nil, nil)
		err := RunVolumeLabelAdd(cfg)
		assert.NoError(t, err)
	})
}

func TestRunVolumeLabelAddErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.GrandParentName, config.ArgDataCenterId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgVolumeId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelKey), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelValue), testLabelResourceVar)
		rm.Label.EXPECT().VolumeCreate(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar, testLabelResourceVar).Return(&testLabelResourceRes, nil, testLabelResourceErr)
		err := RunVolumeLabelAdd(cfg)
		assert.Error(t, err)
	})
}

func TestRunVolumeLabelRemove(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.GrandParentName, config.ArgDataCenterId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgVolumeId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelKey), testLabelResourceVar)
		rm.Label.EXPECT().VolumeDelete(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar).Return(nil, nil)
		err := RunVolumeLabelRemove(cfg)
		assert.NoError(t, err)
	})
}

func TestRunVolumeLabelRemoveErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.GrandParentName, config.ArgDataCenterId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgVolumeId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelKey), testLabelResourceVar)
		rm.Label.EXPECT().VolumeDelete(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar).Return(nil, testLabelResourceErr)
		err := RunVolumeLabelRemove(cfg)
		assert.Error(t, err)
	})
}

func TestGetDatacenterLabelIds(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() { return }
	w := bufio.NewWriter(&b)
	viper.Set(config.ArgConfig, "../pkg/testdata/config.json")
	getDatacenterLabelIds(w, testLabelResourceVar)
	err := w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`401 Unauthorized`)
	assert.True(t, re.Match(b.Bytes()))
}

func TestGetIpBlockLabelIds(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() { return }
	w := bufio.NewWriter(&b)
	viper.Set(config.ArgConfig, "../pkg/testdata/config.json")
	getIPBlockLabelIds(w, testLabelResourceVar)
	err := w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`401 Unauthorized`)
	assert.True(t, re.Match(b.Bytes()))
}

func TestGetSnapshotLabelIds(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() { return }
	w := bufio.NewWriter(&b)
	viper.Set(config.ArgConfig, "../pkg/testdata/config.json")
	getSnapshotLabelIds(w, testLabelResourceVar)
	err := w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`401 Unauthorized`)
	assert.True(t, re.Match(b.Bytes()))
}

func TestGetServerLabelIds(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() { return }
	w := bufio.NewWriter(&b)
	viper.Set(config.ArgConfig, "../pkg/testdata/config.json")
	getServerLabelIds(w, testLabelResourceVar, testLabelResourceVar)
	err := w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`401 Unauthorized`)
	assert.True(t, re.Match(b.Bytes()))
}

func TestGetVolumeLabelIds(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() { return }
	w := bufio.NewWriter(&b)
	viper.Set(config.ArgConfig, "../pkg/testdata/config.json")
	getVolumeLabelIds(w, testLabelResourceVar, testLabelResourceVar)
	err := w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`401 Unauthorized`)
	assert.True(t, re.Match(b.Bytes()))
}
