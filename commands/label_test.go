package commands

import (
	"bufio"
	"bytes"
	"errors"
	"testing"

	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	testLabel = ionoscloud.Label{
		Id: &testLabelVar,
		Properties: &ionoscloud.LabelProperties{
			Key:          &testLabelVar,
			Value:        &testLabelVar,
			ResourceId:   &testLabelVar,
			ResourceType: &testLabelVar,
		},
	}
	testLabels = resources.Labels{
		Labels: ionoscloud.Labels{
			Id:    &testLabelVar,
			Items: &[]ionoscloud.Label{testLabel},
		},
	}
	testLabelVar = "test-label"
	testLabelErr = errors.New("label test error")
)

func TestPreRunGlobalResourceTypeLabelKey(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgType), config.DatacenterResource)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgDataCenterId), testLabelVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLabelKey), testLabelVar)
		viper.Set(config.ArgQuiet, false)
		err := PreRunGlobalResourceTypeLabelKey(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunGlobalResourceTypeLabelKeyErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgType), "")
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLabelKey), "")
		viper.Set(config.ArgQuiet, false)
		err := PreRunGlobalResourceTypeLabelKey(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunGlobalResourceTypeLabelKeyValue(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgType), config.DatacenterResource)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgDataCenterId), testLabelVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLabelKey), testLabelVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLabelValue), testLabelVar)
		viper.Set(config.ArgQuiet, false)
		err := PreRunGlobalResourceTypeLabelKeyValue(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunGlobalResourceTypeLabelKeyValueErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgType), "")
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLabelKey), "")
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLabelValue), "")
		viper.Set(config.ArgQuiet, false)
		err := PreRunGlobalResourceTypeLabelKeyValue(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunGlobalResourceTypeLabelKeyResourceErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgType), config.DatacenterResource)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgDataCenterId), "")
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLabelKey), "")
		viper.Set(config.ArgQuiet, false)
		err := PreRunGlobalResourceTypeLabelKey(cfg)
		assert.Error(t, err)
	})
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgType), config.ServerResource)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgDataCenterId), "")
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgServerId), "")
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLabelKey), "")
		viper.Set(config.ArgQuiet, false)
		err := PreRunGlobalResourceTypeLabelKey(cfg)
		assert.Error(t, err)
	})
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgType), config.VolumeResource)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgDataCenterId), "")
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgVolumeId), "")
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLabelKey), "")
		viper.Set(config.ArgQuiet, false)
		err := PreRunGlobalResourceTypeLabelKey(cfg)
		assert.Error(t, err)
	})
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgType), config.IpBlockResource)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgIpBlockId), "")
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLabelKey), "")
		viper.Set(config.ArgQuiet, false)
		err := PreRunGlobalResourceTypeLabelKey(cfg)
		assert.Error(t, err)
	})
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgType), config.SnapshotResource)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgSnapshotId), "")
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLabelKey), "")
		viper.Set(config.ArgQuiet, false)
		err := PreRunGlobalResourceTypeLabelKey(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunLabelUrn(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLabelUrn), testLabelVar)
		viper.Set(config.ArgQuiet, false)
		err := PreRunLabelUrn(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunLabelUrnErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLabelUrn), "")
		viper.Set(config.ArgQuiet, false)
		err := PreRunLabelUrn(cfg)
		assert.Error(t, err)
	})
}

func TestRunLabelList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		rm.Label.EXPECT().List().Return(testLabels, nil, nil)
		err := RunLabelList(cfg)
		assert.NoError(t, err)
	})
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgType), config.DatacenterResource)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgDataCenterId), testLabelResourceVar)
		rm.Label.EXPECT().DatacenterList(testLabelResourceVar).Return(testLabelResources, nil, nil)
		err := RunLabelList(cfg)
		assert.NoError(t, err)
	})
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgType), config.IpBlockResource)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgIpBlockId), testLabelResourceVar)
		rm.Label.EXPECT().IpBlockList(testLabelResourceVar).Return(testLabelResources, nil, nil)
		err := RunLabelList(cfg)
		assert.NoError(t, err)
	})
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgType), config.SnapshotResource)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgSnapshotId), testLabelResourceVar)
		rm.Label.EXPECT().SnapshotList(testLabelResourceVar).Return(testLabelResources, nil, nil)
		err := RunLabelList(cfg)
		assert.NoError(t, err)
	})
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgType), config.ServerResource)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgDataCenterId), testLabelResourceVar)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgServerId), testLabelResourceVar)
		rm.Label.EXPECT().ServerList(testLabelResourceVar, testLabelResourceVar).Return(testLabelResources, nil, nil)
		err := RunLabelList(cfg)
		assert.NoError(t, err)
	})
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgType), config.VolumeResource)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgDataCenterId), testLabelResourceVar)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgVolumeId), testLabelResourceVar)
		rm.Label.EXPECT().VolumeList(testLabelResourceVar, testLabelResourceVar).Return(testLabelResources, nil, nil)
		err := RunLabelList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunLabelListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		rm.Label.EXPECT().List().Return(testLabels, nil, testLabelErr)
		err := RunLabelList(cfg)
		assert.Error(t, err)
		assert.True(t, err == testLabelErr)
	})
}

func TestRunLabelGetByUrn(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLabelUrn), testLabelVar)
		label := resources.Label{Label: testLabel}
		rm.Label.EXPECT().GetByUrn(testLabelVar).Return(&label, nil, nil)
		err := RunLabelGetByUrn(cfg)
		assert.NoError(t, err)
	})
}

func TestRunLabelGetByUrnErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLabelUrn), testLabelVar)
		label := resources.Label{Label: testLabel}
		rm.Label.EXPECT().GetByUrn(testLabelVar).Return(&label, nil, testLabelErr)
		err := RunLabelGetByUrn(cfg)
		assert.Error(t, err)
	})
}

func TestRunLabelGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgType), config.DatacenterResource)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgDataCenterId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLabelKey), testLabelResourceVar)
		rm.Label.EXPECT().DatacenterGet(testLabelResourceVar, testLabelResourceVar).Return(&testLabelResourceRes, nil, nil)
		err := RunLabelGet(cfg)
		assert.NoError(t, err)
	})
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgType), config.IpBlockResource)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgIpBlockId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLabelKey), testLabelResourceVar)
		rm.Label.EXPECT().IpBlockGet(testLabelResourceVar, testLabelResourceVar).Return(&testLabelResourceRes, nil, nil)
		err := RunLabelGet(cfg)
		assert.NoError(t, err)
	})
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgType), config.SnapshotResource)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgSnapshotId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLabelKey), testLabelResourceVar)
		rm.Label.EXPECT().SnapshotGet(testLabelResourceVar, testLabelResourceVar).Return(&testLabelResourceRes, nil, nil)
		err := RunLabelGet(cfg)
		assert.NoError(t, err)
	})
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgType), config.ServerResource)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgDataCenterId), testLabelResourceVar)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgServerId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLabelKey), testLabelResourceVar)
		rm.Label.EXPECT().ServerGet(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar).Return(&testLabelResourceRes, nil, nil)
		err := RunLabelGet(cfg)
		assert.NoError(t, err)
	})
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgType), config.VolumeResource)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgDataCenterId), testLabelResourceVar)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgVolumeId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLabelKey), testLabelResourceVar)
		rm.Label.EXPECT().VolumeGet(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar).Return(&testLabelResourceRes, nil, nil)
		err := RunLabelGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunLabelAdd(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgType), config.DatacenterResource)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgDataCenterId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLabelKey), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLabelValue), testLabelResourceVar)
		rm.Label.EXPECT().DatacenterCreate(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar).Return(&testLabelResourceRes, nil, nil)
		err := RunLabelAdd(cfg)
		assert.NoError(t, err)
	})
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgType), config.IpBlockResource)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgIpBlockId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLabelKey), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLabelValue), testLabelResourceVar)
		rm.Label.EXPECT().IpBlockCreate(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar).Return(&testLabelResourceRes, nil, nil)
		err := RunLabelAdd(cfg)
		assert.NoError(t, err)
	})
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgType), config.SnapshotResource)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgSnapshotId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLabelKey), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLabelValue), testLabelResourceVar)
		rm.Label.EXPECT().SnapshotCreate(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar).Return(&testLabelResourceRes, nil, nil)
		err := RunLabelAdd(cfg)
		assert.NoError(t, err)
	})
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgType), config.ServerResource)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgDataCenterId), testLabelResourceVar)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgServerId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLabelKey), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLabelValue), testLabelResourceVar)
		rm.Label.EXPECT().ServerCreate(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar, testLabelResourceVar).Return(&testLabelResourceRes, nil, nil)
		err := RunLabelAdd(cfg)
		assert.NoError(t, err)
	})
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgType), config.VolumeResource)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgDataCenterId), testLabelResourceVar)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgVolumeId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLabelKey), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLabelValue), testLabelResourceVar)
		rm.Label.EXPECT().VolumeCreate(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar, testLabelResourceVar).Return(&testLabelResourceRes, nil, nil)
		err := RunLabelAdd(cfg)
		assert.NoError(t, err)
	})
}

func TestRunLabelRemove(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgType), config.DatacenterResource)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgDataCenterId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLabelKey), testLabelResourceVar)
		rm.Label.EXPECT().DatacenterDelete(testLabelResourceVar, testLabelResourceVar).Return(nil, nil)
		err := RunLabelRemove(cfg)
		assert.NoError(t, err)
	})
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgType), config.IpBlockResource)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgIpBlockId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLabelKey), testLabelResourceVar)
		rm.Label.EXPECT().IpBlockDelete(testLabelResourceVar, testLabelResourceVar).Return(nil, nil)
		err := RunLabelRemove(cfg)
		assert.NoError(t, err)
	})
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgType), config.SnapshotResource)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgSnapshotId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLabelKey), testLabelResourceVar)
		rm.Label.EXPECT().SnapshotDelete(testLabelResourceVar, testLabelResourceVar).Return(nil, nil)
		err := RunLabelRemove(cfg)
		assert.NoError(t, err)
	})
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgType), config.ServerResource)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgDataCenterId), testLabelResourceVar)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgServerId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLabelKey), testLabelResourceVar)
		rm.Label.EXPECT().ServerDelete(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar).Return(nil, nil)
		err := RunLabelRemove(cfg)
		assert.NoError(t, err)
	})
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgType), config.VolumeResource)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgDataCenterId), testLabelResourceVar)
		viper.Set(core.GetGlobalFlagName(cfg.Resource, config.ArgVolumeId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLabelKey), testLabelResourceVar)
		rm.Label.EXPECT().VolumeDelete(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar).Return(nil, nil)
		err := RunLabelRemove(cfg)
		assert.NoError(t, err)
	})
}
