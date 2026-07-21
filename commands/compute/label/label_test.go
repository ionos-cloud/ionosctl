package label

import (
	"bufio"
	"bytes"
	"errors"
	"testing"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"

	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
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

func TestLabelCmd(t *testing.T) {
	var err error
	core.RootCmdTest.AddCommand(LabelCmd())
	if ok := LabelCmd().IsAvailableCommand(); !ok {
		err = errors.New("non-available cmd")
	}
	assert.NoError(t, err)
}

func TestPreRunResourceTypeLabelKey(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		cfg.SetFlag(cloudapiv6.ArgResourceType, cloudapiv6.DatacenterResource)
		cfg.SetFlag(cloudapiv6.ArgDataCenterId, testLabelVar)
		cfg.SetFlag(cloudapiv6.ArgLabelKey, testLabelVar)
		viper.Set(constants.ArgQuiet, false)
		err := PreRunResourceTypeLabelKey(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunResourceTypeLabelKeyErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		err := PreRunResourceTypeLabelKey(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunResourceTypeLabelKeyValue(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		cfg.SetFlag(cloudapiv6.ArgResourceType, cloudapiv6.DatacenterResource)
		cfg.SetFlag(cloudapiv6.ArgDataCenterId, testLabelVar)
		cfg.SetFlag(cloudapiv6.ArgLabelKey, testLabelVar)
		cfg.SetFlag(cloudapiv6.ArgLabelValue, testLabelVar)
		viper.Set(constants.ArgQuiet, false)
		err := PreRunResourceTypeLabelKeyValue(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunResourceTypeLabelKeyValueErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		err := PreRunResourceTypeLabelKeyValue(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunResourceTypeLabelKeyValueResourceErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		cfg.SetFlag(cloudapiv6.ArgResourceType, cloudapiv6.DatacenterResource)
		cfg.SetFlag(cloudapiv6.ArgLabelKey, "")
		viper.Set(constants.ArgQuiet, false)
		err := PreRunResourceTypeLabelKey(cfg)
		assert.Error(t, err)
	})
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		cfg.SetFlag(cloudapiv6.ArgResourceType, cloudapiv6.ServerResource)
		viper.Set(constants.ArgQuiet, false)
		err := PreRunResourceTypeLabelKey(cfg)
		assert.Error(t, err)
	})
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		cfg.SetFlag(cloudapiv6.ArgResourceType, cloudapiv6.VolumeResource)
		viper.Set(constants.ArgQuiet, false)
		err := PreRunResourceTypeLabelKey(cfg)
		assert.Error(t, err)
	})
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		cfg.SetFlag(cloudapiv6.ArgResourceType, cloudapiv6.IpBlockResource)
		viper.Set(constants.ArgQuiet, false)
		err := PreRunResourceTypeLabelKey(cfg)
		assert.Error(t, err)
	})
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		cfg.SetFlag(cloudapiv6.ArgResourceType, cloudapiv6.SnapshotResource)
		viper.Set(constants.ArgQuiet, false)
		err := PreRunResourceTypeLabelKey(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunLabelUrn(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		cfg.SetFlag(cloudapiv6.ArgLabelUrn, testLabelVar)
		viper.Set(constants.ArgQuiet, false)
		err := PreRunLabelUrn(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunLabelUrnErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		err := PreRunLabelUrn(cfg)
		assert.Error(t, err)
	})
}

func TestRunLabelList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		rm.CloudApiV6Mocks.Label.EXPECT().List().Return(testLabels, nil, nil)
		err := RunLabelList(cfg)
		assert.NoError(t, err)
	})
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		cfg.SetFlag(cloudapiv6.ArgResourceType, cloudapiv6.DatacenterResource)
		cfg.SetFlag(cloudapiv6.ArgDataCenterId, testLabelResourceVar)
		rm.CloudApiV6Mocks.Label.EXPECT().DatacenterList(testLabelResourceVar).Return(testLabelResources, nil, nil)
		err := RunLabelList(cfg)
		assert.NoError(t, err)
	})
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		cfg.SetFlag(cloudapiv6.ArgResourceType, cloudapiv6.IpBlockResource)
		cfg.SetFlag(cloudapiv6.ArgIpBlockId, testLabelResourceVar)
		rm.CloudApiV6Mocks.Label.EXPECT().IpBlockList(testLabelResourceVar).Return(testLabelResources, nil, nil)
		err := RunLabelList(cfg)
		assert.NoError(t, err)
	})
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		cfg.SetFlag(cloudapiv6.ArgResourceType, cloudapiv6.SnapshotResource)
		cfg.SetFlag(cloudapiv6.ArgSnapshotId, testLabelResourceVar)
		rm.CloudApiV6Mocks.Label.EXPECT().SnapshotList(testLabelResourceVar).Return(testLabelResources, nil, nil)
		err := RunLabelList(cfg)
		assert.NoError(t, err)
	})
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		cfg.SetFlag(cloudapiv6.ArgResourceType, cloudapiv6.ServerResource)
		cfg.SetFlag(cloudapiv6.ArgDataCenterId, testLabelResourceVar)
		cfg.SetFlag(cloudapiv6.ArgServerId, testLabelResourceVar)
		rm.CloudApiV6Mocks.Label.EXPECT().ServerList(testLabelResourceVar, testLabelResourceVar).Return(testLabelResources, nil, nil)
		err := RunLabelList(cfg)
		assert.NoError(t, err)
	})
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		cfg.SetFlag(cloudapiv6.ArgResourceType, cloudapiv6.VolumeResource)
		cfg.SetFlag(cloudapiv6.ArgDataCenterId, testLabelResourceVar)
		cfg.SetFlag(cloudapiv6.ArgVolumeId, testLabelResourceVar)
		rm.CloudApiV6Mocks.Label.EXPECT().VolumeList(testLabelResourceVar, testLabelResourceVar).Return(testLabelResources, nil, nil)
		err := RunLabelList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunLabelListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		rm.CloudApiV6Mocks.Label.EXPECT().List().Return(testLabels, nil, testLabelErr)
		err := RunLabelList(cfg)
		assert.Error(t, err)
		assert.True(t, err == testLabelErr)
	})
}

func TestRunLabelGetByUrn(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		cfg.SetFlag(cloudapiv6.ArgLabelUrn, testLabelVar)
		label := resources.Label{Label: testLabel}
		rm.CloudApiV6Mocks.Label.EXPECT().GetByUrn(testLabelVar).Return(&label, nil, nil)
		err := RunLabelGetByUrn(cfg)
		assert.NoError(t, err)
	})
}

func TestRunLabelGetByUrnErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		cfg.SetFlag(cloudapiv6.ArgLabelUrn, testLabelVar)
		label := resources.Label{Label: testLabel}
		rm.CloudApiV6Mocks.Label.EXPECT().GetByUrn(testLabelVar).Return(&label, nil, testLabelErr)
		err := RunLabelGetByUrn(cfg)
		assert.Error(t, err)
	})
}

func TestRunLabelGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		cfg.SetFlag(cloudapiv6.ArgResourceType, cloudapiv6.DatacenterResource)
		cfg.SetFlag(cloudapiv6.ArgDataCenterId, testLabelResourceVar)
		cfg.SetFlag(cloudapiv6.ArgLabelKey, testLabelResourceVar)
		rm.CloudApiV6Mocks.Label.EXPECT().DatacenterGet(testLabelResourceVar, testLabelResourceVar).Return(&testLabelResourceRes, nil, nil)
		err := RunLabelGet(cfg)
		assert.NoError(t, err)
	})
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		cfg.SetFlag(cloudapiv6.ArgResourceType, cloudapiv6.IpBlockResource)
		cfg.SetFlag(cloudapiv6.ArgIpBlockId, testLabelResourceVar)
		cfg.SetFlag(cloudapiv6.ArgLabelKey, testLabelResourceVar)
		rm.CloudApiV6Mocks.Label.EXPECT().IpBlockGet(testLabelResourceVar, testLabelResourceVar).Return(&testLabelResourceRes, nil, nil)
		err := RunLabelGet(cfg)
		assert.NoError(t, err)
	})
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		cfg.SetFlag(cloudapiv6.ArgResourceType, cloudapiv6.SnapshotResource)
		cfg.SetFlag(cloudapiv6.ArgSnapshotId, testLabelResourceVar)
		cfg.SetFlag(cloudapiv6.ArgLabelKey, testLabelResourceVar)
		rm.CloudApiV6Mocks.Label.EXPECT().SnapshotGet(testLabelResourceVar, testLabelResourceVar).Return(&testLabelResourceRes, nil, nil)
		err := RunLabelGet(cfg)
		assert.NoError(t, err)
	})
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		cfg.SetFlag(cloudapiv6.ArgResourceType, cloudapiv6.ServerResource)
		cfg.SetFlag(cloudapiv6.ArgDataCenterId, testLabelResourceVar)
		cfg.SetFlag(cloudapiv6.ArgServerId, testLabelResourceVar)
		cfg.SetFlag(cloudapiv6.ArgLabelKey, testLabelResourceVar)
		rm.CloudApiV6Mocks.Label.EXPECT().ServerGet(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar).Return(&testLabelResourceRes, nil, nil)
		err := RunLabelGet(cfg)
		assert.NoError(t, err)
	})
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		cfg.SetFlag(cloudapiv6.ArgResourceType, cloudapiv6.VolumeResource)
		cfg.SetFlag(cloudapiv6.ArgDataCenterId, testLabelResourceVar)
		cfg.SetFlag(cloudapiv6.ArgVolumeId, testLabelResourceVar)
		cfg.SetFlag(cloudapiv6.ArgLabelKey, testLabelResourceVar)
		rm.CloudApiV6Mocks.Label.EXPECT().VolumeGet(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar).Return(&testLabelResourceRes, nil, nil)
		err := RunLabelGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunLabelAdd(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		cfg.SetFlag(cloudapiv6.ArgResourceType, cloudapiv6.DatacenterResource)
		cfg.SetFlag(cloudapiv6.ArgDataCenterId, testLabelResourceVar)
		cfg.SetFlag(cloudapiv6.ArgLabelKey, testLabelResourceVar)
		cfg.SetFlag(cloudapiv6.ArgLabelValue, testLabelResourceVar)
		rm.CloudApiV6Mocks.Label.EXPECT().DatacenterCreate(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar).Return(&testLabelResourceRes, nil, nil)
		err := RunLabelAdd(cfg)
		assert.NoError(t, err)
	})
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		cfg.SetFlag(cloudapiv6.ArgResourceType, cloudapiv6.IpBlockResource)
		cfg.SetFlag(cloudapiv6.ArgIpBlockId, testLabelResourceVar)
		cfg.SetFlag(cloudapiv6.ArgLabelKey, testLabelResourceVar)
		cfg.SetFlag(cloudapiv6.ArgLabelValue, testLabelResourceVar)
		rm.CloudApiV6Mocks.Label.EXPECT().IpBlockCreate(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar).Return(&testLabelResourceRes, nil, nil)
		err := RunLabelAdd(cfg)
		assert.NoError(t, err)
	})
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		cfg.SetFlag(cloudapiv6.ArgResourceType, cloudapiv6.SnapshotResource)
		cfg.SetFlag(cloudapiv6.ArgSnapshotId, testLabelResourceVar)
		cfg.SetFlag(cloudapiv6.ArgLabelKey, testLabelResourceVar)
		cfg.SetFlag(cloudapiv6.ArgLabelValue, testLabelResourceVar)
		rm.CloudApiV6Mocks.Label.EXPECT().SnapshotCreate(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar).Return(&testLabelResourceRes, nil, nil)
		err := RunLabelAdd(cfg)
		assert.NoError(t, err)
	})
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		cfg.SetFlag(cloudapiv6.ArgResourceType, cloudapiv6.ServerResource)
		cfg.SetFlag(cloudapiv6.ArgDataCenterId, testLabelResourceVar)
		cfg.SetFlag(cloudapiv6.ArgServerId, testLabelResourceVar)
		cfg.SetFlag(cloudapiv6.ArgLabelKey, testLabelResourceVar)
		cfg.SetFlag(cloudapiv6.ArgLabelValue, testLabelResourceVar)
		rm.CloudApiV6Mocks.Label.EXPECT().ServerCreate(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar, testLabelResourceVar).Return(&testLabelResourceRes, nil, nil)
		err := RunLabelAdd(cfg)
		assert.NoError(t, err)
	})
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		cfg.SetFlag(cloudapiv6.ArgResourceType, cloudapiv6.VolumeResource)
		cfg.SetFlag(cloudapiv6.ArgDataCenterId, testLabelResourceVar)
		cfg.SetFlag(cloudapiv6.ArgVolumeId, testLabelResourceVar)
		cfg.SetFlag(cloudapiv6.ArgLabelKey, testLabelResourceVar)
		cfg.SetFlag(cloudapiv6.ArgLabelValue, testLabelResourceVar)
		rm.CloudApiV6Mocks.Label.EXPECT().VolumeCreate(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar, testLabelResourceVar).Return(&testLabelResourceRes, nil, nil)
		err := RunLabelAdd(cfg)
		assert.NoError(t, err)
	})
}

func TestRunLabelRemove(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		cfg.SetFlag(cloudapiv6.ArgResourceType, cloudapiv6.DatacenterResource)
		cfg.SetFlag(cloudapiv6.ArgDataCenterId, testLabelResourceVar)
		cfg.SetFlag(cloudapiv6.ArgLabelKey, testLabelResourceVar)
		rm.CloudApiV6Mocks.Label.EXPECT().DatacenterDelete(testLabelResourceVar, testLabelResourceVar).Return(nil, nil)
		err := RunLabelRemove(cfg)
		assert.NoError(t, err)
	})
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		cfg.SetFlag(cloudapiv6.ArgResourceType, cloudapiv6.IpBlockResource)
		cfg.SetFlag(cloudapiv6.ArgIpBlockId, testLabelResourceVar)
		cfg.SetFlag(cloudapiv6.ArgLabelKey, testLabelResourceVar)
		rm.CloudApiV6Mocks.Label.EXPECT().IpBlockDelete(testLabelResourceVar, testLabelResourceVar).Return(nil, nil)
		err := RunLabelRemove(cfg)
		assert.NoError(t, err)
	})
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		cfg.SetFlag(cloudapiv6.ArgResourceType, cloudapiv6.SnapshotResource)
		cfg.SetFlag(cloudapiv6.ArgSnapshotId, testLabelResourceVar)
		cfg.SetFlag(cloudapiv6.ArgLabelKey, testLabelResourceVar)
		rm.CloudApiV6Mocks.Label.EXPECT().SnapshotDelete(testLabelResourceVar, testLabelResourceVar).Return(nil, nil)
		err := RunLabelRemove(cfg)
		assert.NoError(t, err)
	})
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		cfg.SetFlag(cloudapiv6.ArgResourceType, cloudapiv6.ServerResource)
		cfg.SetFlag(cloudapiv6.ArgDataCenterId, testLabelResourceVar)
		cfg.SetFlag(cloudapiv6.ArgServerId, testLabelResourceVar)
		cfg.SetFlag(cloudapiv6.ArgLabelKey, testLabelResourceVar)
		rm.CloudApiV6Mocks.Label.EXPECT().ServerDelete(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar).Return(nil, nil)
		err := RunLabelRemove(cfg)
		assert.NoError(t, err)
	})
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		cfg.SetFlag(cloudapiv6.ArgResourceType, cloudapiv6.VolumeResource)
		cfg.SetFlag(cloudapiv6.ArgDataCenterId, testLabelResourceVar)
		cfg.SetFlag(cloudapiv6.ArgVolumeId, testLabelResourceVar)
		cfg.SetFlag(cloudapiv6.ArgLabelKey, testLabelResourceVar)
		rm.CloudApiV6Mocks.Label.EXPECT().VolumeDelete(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar).Return(nil, nil)
		err := RunLabelRemove(cfg)
		assert.NoError(t, err)
	})
}
