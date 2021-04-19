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

func TestPreRunDcIdLabelKeyValidate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgDataCenterId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelKey), testLabelResourceVar)
		viper.Set(config.ArgQuiet, false)
		err := PreRunDcIdLabelKeyValidate(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunDcIdLabelKeyValidateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgDataCenterId), "")
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelKey), "")
		viper.Set(config.ArgQuiet, false)
		err := PreRunDcIdLabelKeyValidate(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunDcIdLabelKeyValueValidate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgDataCenterId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelKey), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelValue), testLabelResourceVar)
		viper.Set(config.ArgQuiet, false)
		err := PreRunDcIdLabelKeyValueValidate(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunDcIdLabelKeyValueValidateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgDataCenterId), "")
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelKey), "")
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelValue), "")
		viper.Set(config.ArgQuiet, false)
		err := PreRunDcIdLabelKeyValueValidate(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunIpBlockIdLabelKeyValidate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgIpBlockId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelKey), testLabelResourceVar)
		viper.Set(config.ArgQuiet, false)
		err := PreRunIpBlockIdLabelKeyValidate(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunIpBlockIdLabelKeyValidateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgIpBlockId), "")
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelKey), "")
		viper.Set(config.ArgQuiet, false)
		err := PreRunIpBlockIdLabelKeyValidate(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunIpBlockIdLabelKeyValueValidate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgIpBlockId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelKey), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelValue), testLabelResourceVar)
		viper.Set(config.ArgQuiet, false)
		err := PreRunIpBlockIdLabelKeyValueValidate(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunIpBlockIdLabelKeyValueValidateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgIpBlockId), "")
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelKey), "")
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelValue), "")
		viper.Set(config.ArgQuiet, false)
		err := PreRunIpBlockIdLabelKeyValueValidate(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunSnapshotIdLabelKeyValidate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgSnapshotId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelKey), testLabelResourceVar)
		viper.Set(config.ArgQuiet, false)
		err := PreRunSnapshotIdLabelKeyValidate(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunSnapshotIdLabelKeyValidateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgSnapshotId), "")
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelKey), "")
		viper.Set(config.ArgQuiet, false)
		err := PreRunSnapshotIdLabelKeyValidate(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunSnapshotIdLabelKeyValueValidate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgSnapshotId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelKey), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelValue), testLabelResourceVar)
		viper.Set(config.ArgQuiet, false)
		err := PreRunSnapshotIdLabelKeyValueValidate(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunSnapshotIdLabelKeyValueValidateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgSnapshotId), "")
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelKey), "")
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelValue), "")
		viper.Set(config.ArgQuiet, false)
		err := PreRunSnapshotIdLabelKeyValueValidate(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunGlobalDcIdServerLabelKeyValidate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgServerId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelKey), testLabelResourceVar)
		viper.Set(config.ArgQuiet, false)
		err := PreRunGlobalDcIdServerLabelKeyValidate(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunGlobalDcIdServerLabelKeyValidateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), "")
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgServerId), "")
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelKey), "")
		err := PreRunGlobalDcIdServerLabelKeyValidate(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunGlobalDcIdServerLabelKeyValidateGlobalErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), "")
		err := PreRunGlobalDcIdServerLabelKeyValidate(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunGlobalDcIdServerLabelKeyValueValidateGlobalErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), "")
		err := PreRunGlobalDcIdServerLabelKeyValueValidate(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunGlobalDcIdServerLabelKeyValueValidate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgServerId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelKey), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelValue), testLabelResourceVar)
		viper.Set(config.ArgQuiet, false)
		err := PreRunGlobalDcIdServerLabelKeyValueValidate(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunGlobalDcIdServerLabelKeyValueValidateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgServerId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelKey), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelValue), "")
		viper.Set(config.ArgQuiet, false)
		err := PreRunGlobalDcIdServerLabelKeyValueValidate(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunGlobalDcIdVolumeLabelKeyValidate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgVolumeId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelKey), testLabelResourceVar)
		viper.Set(config.ArgQuiet, false)
		err := PreRunGlobalDcIdVolumeLabelKeyValidate(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunGlobalDcIdVolumeLabelKeyValidateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgVolumeId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelKey), "")
		viper.Set(config.ArgQuiet, false)
		err := PreRunGlobalDcIdVolumeLabelKeyValidate(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunGlobalDcIdVolumeLabelKeyValidateGlobalErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), "")
		viper.Set(config.ArgQuiet, false)
		err := PreRunGlobalDcIdVolumeLabelKeyValidate(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunGlobalDcIdVolumeLabelKeyValueValidate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgVolumeId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelKey), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelValue), testLabelResourceVar)
		viper.Set(config.ArgQuiet, false)
		err := PreRunGlobalDcIdVolumeLabelKeyValueValidate(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunGlobalDcIdVolumeLabelKeyValueValidateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgVolumeId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelKey), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelValue), "")
		viper.Set(config.ArgQuiet, false)
		err := PreRunGlobalDcIdVolumeLabelKeyValueValidate(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunGlobalDcIdVolumeLabelKeyValueValidateGlobalErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), "")
		viper.Set(config.ArgQuiet, false)
		err := PreRunGlobalDcIdVolumeLabelKeyValueValidate(cfg)
		assert.Error(t, err)
	})
}

func TestRunDatacenterListLabels(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgDataCenterId), testLabelResourceVar)
		rm.Label.EXPECT().DatacenterList(testLabelResourceVar).Return(testLabelResources, nil, nil)
		err := RunDataCenterListLabels(cfg)
		assert.NoError(t, err)
	})
}

func TestRunDatacenterListLabelsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgDataCenterId), testLabelResourceVar)
		rm.Label.EXPECT().DatacenterList(testLabelResourceVar).Return(testLabelResources, nil, testLabelResourceErr)
		err := RunDataCenterListLabels(cfg)
		assert.Error(t, err)
	})
}

func TestRunDatacenterGetLabel(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgDataCenterId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelKey), testLabelResourceVar)
		rm.Label.EXPECT().DatacenterGet(testLabelResourceVar, testLabelResourceVar).Return(&testLabelResourceRes, nil, nil)
		err := RunDataCenterGetLabel(cfg)
		assert.NoError(t, err)
	})
}

func TestRunDatacenterGetLabelErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgDataCenterId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelKey), testLabelResourceVar)
		rm.Label.EXPECT().DatacenterGet(testLabelResourceVar, testLabelResourceVar).Return(&testLabelResourceRes, nil, testLabelResourceErr)
		err := RunDataCenterGetLabel(cfg)
		assert.Error(t, err)
	})
}

func TestRunDatacenterAddLabel(t *testing.T) {
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
		err := RunDataCenterAddLabel(cfg)
		assert.NoError(t, err)
	})
}

func TestRunDatacenterAddLabelErr(t *testing.T) {
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
		err := RunDataCenterAddLabel(cfg)
		assert.Error(t, err)
	})
}

func TestRunDatacenterRemoveLabel(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgDataCenterId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelKey), testLabelResourceVar)
		rm.Label.EXPECT().DatacenterDelete(testLabelResourceVar, testLabelResourceVar).Return(nil, nil)
		err := RunDataCenterRemoveLabel(cfg)
		assert.NoError(t, err)
	})
}

func TestRunDatacenterRemoveLabelErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgDataCenterId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelKey), testLabelResourceVar)
		rm.Label.EXPECT().DatacenterDelete(testLabelResourceVar, testLabelResourceVar).Return(nil, testLabelResourceErr)
		err := RunDataCenterRemoveLabel(cfg)
		assert.Error(t, err)
	})
}

func TestRunIpBlockListLabels(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgIpBlockId), testLabelResourceVar)
		rm.Label.EXPECT().IpBlockList(testLabelResourceVar).Return(testLabelResources, nil, nil)
		err := RunIpBlockListLabels(cfg)
		assert.NoError(t, err)
	})
}

func TestRunIpBlockListLabelsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgIpBlockId), testLabelResourceVar)
		rm.Label.EXPECT().IpBlockList(testLabelResourceVar).Return(testLabelResources, nil, testLabelResourceErr)
		err := RunIpBlockListLabels(cfg)
		assert.Error(t, err)
	})
}

func TestRunIpBlockGetLabel(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgIpBlockId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelKey), testLabelResourceVar)
		rm.Label.EXPECT().IpBlockGet(testLabelResourceVar, testLabelResourceVar).Return(&testLabelResourceRes, nil, nil)
		err := RunIpBlockGetLabel(cfg)
		assert.NoError(t, err)
	})
}

func TestRunIpBlockGetLabelErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgIpBlockId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelKey), testLabelResourceVar)
		rm.Label.EXPECT().IpBlockGet(testLabelResourceVar, testLabelResourceVar).Return(&testLabelResourceRes, nil, testLabelResourceErr)
		err := RunIpBlockGetLabel(cfg)
		assert.Error(t, err)
	})
}

func TestRunIpBlockAddLabel(t *testing.T) {
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
		err := RunIpBlockAddLabel(cfg)
		assert.NoError(t, err)
	})
}

func TestRunIpBlockAddLabelErr(t *testing.T) {
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
		err := RunIpBlockAddLabel(cfg)
		assert.Error(t, err)
	})
}

func TestRunIpBlockRemoveLabel(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgIpBlockId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelKey), testLabelResourceVar)
		rm.Label.EXPECT().IpBlockDelete(testLabelResourceVar, testLabelResourceVar).Return(nil, nil)
		err := RunIpBlockRemoveLabel(cfg)
		assert.NoError(t, err)
	})
}

func TestRunIpBlockRemoveLabelErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgIpBlockId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelKey), testLabelResourceVar)
		rm.Label.EXPECT().IpBlockDelete(testLabelResourceVar, testLabelResourceVar).Return(nil, testLabelResourceErr)
		err := RunIpBlockRemoveLabel(cfg)
		assert.Error(t, err)
	})
}

func TestRunSnapshotListLabels(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgSnapshotId), testLabelResourceVar)
		rm.Label.EXPECT().SnapshotList(testLabelResourceVar).Return(testLabelResources, nil, nil)
		err := RunSnapshotListLabels(cfg)
		assert.NoError(t, err)
	})
}

func TestRunSnapshotListLabelsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgSnapshotId), testLabelResourceVar)
		rm.Label.EXPECT().SnapshotList(testLabelResourceVar).Return(testLabelResources, nil, testLabelResourceErr)
		err := RunSnapshotListLabels(cfg)
		assert.Error(t, err)
	})
}

func TestRunSnapshotGetLabel(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgSnapshotId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelKey), testLabelResourceVar)
		rm.Label.EXPECT().SnapshotGet(testLabelResourceVar, testLabelResourceVar).Return(&testLabelResourceRes, nil, nil)
		err := RunSnapshotGetLabel(cfg)
		assert.NoError(t, err)
	})
}

func TestRunSnapshotGetLabelErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgSnapshotId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelKey), testLabelResourceVar)
		rm.Label.EXPECT().SnapshotGet(testLabelResourceVar, testLabelResourceVar).Return(&testLabelResourceRes, nil, testLabelResourceErr)
		err := RunSnapshotGetLabel(cfg)
		assert.Error(t, err)
	})
}

func TestRunSnapshotAddLabel(t *testing.T) {
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
		err := RunSnapshotAddLabel(cfg)
		assert.NoError(t, err)
	})
}

func TestRunSnapshotAddLabelErr(t *testing.T) {
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
		err := RunSnapshotAddLabel(cfg)
		assert.Error(t, err)
	})
}

func TestRunSnapshotRemoveLabel(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgSnapshotId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelKey), testLabelResourceVar)
		rm.Label.EXPECT().SnapshotDelete(testLabelResourceVar, testLabelResourceVar).Return(nil, nil)
		err := RunSnapshotRemoveLabel(cfg)
		assert.NoError(t, err)
	})
}

func TestRunSnapshotRemoveLabelErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgSnapshotId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelKey), testLabelResourceVar)
		rm.Label.EXPECT().SnapshotDelete(testLabelResourceVar, testLabelResourceVar).Return(nil, testLabelResourceErr)
		err := RunSnapshotRemoveLabel(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerListLabels(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgServerId), testLabelResourceVar)
		rm.Label.EXPECT().ServerList(testLabelResourceVar, testLabelResourceVar).Return(testLabelResources, nil, nil)
		err := RunServerListLabels(cfg)
		assert.NoError(t, err)
	})
}

func TestRunServerListLabelsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgServerId), testLabelResourceVar)
		rm.Label.EXPECT().ServerList(testLabelResourceVar, testLabelResourceVar).Return(testLabelResources, nil, testLabelResourceErr)
		err := RunServerListLabels(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerGetLabel(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgServerId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelKey), testLabelResourceVar)
		rm.Label.EXPECT().ServerGet(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar).Return(&testLabelResourceRes, nil, nil)
		err := RunServerGetLabel(cfg)
		assert.NoError(t, err)
	})
}

func TestRunServerGetLabelErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgServerId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelKey), testLabelResourceVar)
		rm.Label.EXPECT().ServerGet(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar).Return(&testLabelResourceRes, nil, testLabelResourceErr)
		err := RunServerGetLabel(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerAddLabel(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgServerId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelKey), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelValue), testLabelResourceVar)
		rm.Label.EXPECT().ServerCreate(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar, testLabelResourceVar).Return(&testLabelResourceRes, nil, nil)
		err := RunServerAddLabel(cfg)
		assert.NoError(t, err)
	})
}

func TestRunServerAddLabelErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgServerId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelKey), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelValue), testLabelResourceVar)
		rm.Label.EXPECT().ServerCreate(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar, testLabelResourceVar).Return(&testLabelResourceRes, nil, testLabelResourceErr)
		err := RunServerAddLabel(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerRemoveLabel(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgServerId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelKey), testLabelResourceVar)
		rm.Label.EXPECT().ServerDelete(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar).Return(nil, nil)
		err := RunServerRemoveLabel(cfg)
		assert.NoError(t, err)
	})
}

func TestRunServerRemoveLabelErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgServerId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelKey), testLabelResourceVar)
		rm.Label.EXPECT().ServerDelete(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar).Return(nil, testLabelResourceErr)
		err := RunServerRemoveLabel(cfg)
		assert.Error(t, err)
	})
}

func TestRunVolumeListLabels(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgVolumeId), testLabelResourceVar)
		rm.Label.EXPECT().VolumeList(testLabelResourceVar, testLabelResourceVar).Return(testLabelResources, nil, nil)
		err := RunVolumeListLabels(cfg)
		assert.NoError(t, err)
	})
}

func TestRunVolumeListLabelsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgVolumeId), testLabelResourceVar)
		rm.Label.EXPECT().VolumeList(testLabelResourceVar, testLabelResourceVar).Return(testLabelResources, nil, testLabelResourceErr)
		err := RunVolumeListLabels(cfg)
		assert.Error(t, err)
	})
}

func TestRunVolumeGetLabel(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgVolumeId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelKey), testLabelResourceVar)
		rm.Label.EXPECT().VolumeGet(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar).Return(&testLabelResourceRes, nil, nil)
		err := RunVolumeGetLabel(cfg)
		assert.NoError(t, err)
	})
}

func TestRunVolumeGetLabelErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgVolumeId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelKey), testLabelResourceVar)
		rm.Label.EXPECT().VolumeGet(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar).Return(&testLabelResourceRes, nil, testLabelResourceErr)
		err := RunVolumeGetLabel(cfg)
		assert.Error(t, err)
	})
}

func TestRunVolumeAddLabel(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgVolumeId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelKey), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelValue), testLabelResourceVar)
		rm.Label.EXPECT().VolumeCreate(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar, testLabelResourceVar).Return(&testLabelResourceRes, nil, nil)
		err := RunVolumeAddLabel(cfg)
		assert.NoError(t, err)
	})
}

func TestRunVolumeAddLabelErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgVolumeId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelKey), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelValue), testLabelResourceVar)
		rm.Label.EXPECT().VolumeCreate(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar, testLabelResourceVar).Return(&testLabelResourceRes, nil, testLabelResourceErr)
		err := RunVolumeAddLabel(cfg)
		assert.Error(t, err)
	})
}

func TestRunVolumeRemoveLabel(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgVolumeId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelKey), testLabelResourceVar)
		rm.Label.EXPECT().VolumeDelete(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar).Return(nil, nil)
		err := RunVolumeRemoveLabel(cfg)
		assert.NoError(t, err)
	})
}

func TestRunVolumeRemoveLabelErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgDataCenterId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgVolumeId), testLabelResourceVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgLabelKey), testLabelResourceVar)
		rm.Label.EXPECT().VolumeDelete(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar).Return(nil, testLabelResourceErr)
		err := RunVolumeRemoveLabel(cfg)
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

func TestGetIPBlockLabelIds(t *testing.T) {
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
