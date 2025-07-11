package commands

import (
	"bufio"
	"bytes"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
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
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagResourceType), cloudapiv6.DatacenterResource)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testLabelVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagLabelKey), testLabelVar)
		viper.Set(constants.FlagQuiet, false)
		err := PreRunResourceTypeLabelKey(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunResourceTypeLabelKeyErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		err := PreRunResourceTypeLabelKey(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunResourceTypeLabelKeyValue(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagResourceType), cloudapiv6.DatacenterResource)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testLabelVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagLabelKey), testLabelVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagLabelValue), testLabelVar)
		viper.Set(constants.FlagQuiet, false)
		err := PreRunResourceTypeLabelKeyValue(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunResourceTypeLabelKeyValueErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		err := PreRunResourceTypeLabelKeyValue(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunResourceTypeLabelKeyValueResourceErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagResourceType), cloudapiv6.DatacenterResource)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagLabelKey), "")
		viper.Set(constants.FlagQuiet, false)
		err := PreRunResourceTypeLabelKey(cfg)
		assert.Error(t, err)
	})
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagResourceType), cloudapiv6.ServerResource)
		viper.Set(constants.FlagQuiet, false)
		err := PreRunResourceTypeLabelKey(cfg)
		assert.Error(t, err)
	})
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagResourceType), cloudapiv6.VolumeResource)
		viper.Set(constants.FlagQuiet, false)
		err := PreRunResourceTypeLabelKey(cfg)
		assert.Error(t, err)
	})
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagResourceType), cloudapiv6.IpBlockResource)
		viper.Set(constants.FlagQuiet, false)
		err := PreRunResourceTypeLabelKey(cfg)
		assert.Error(t, err)
	})
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagResourceType), cloudapiv6.SnapshotResource)
		viper.Set(constants.FlagQuiet, false)
		err := PreRunResourceTypeLabelKey(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunLabelUrn(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagLabelUrn), testLabelVar)
		viper.Set(constants.FlagQuiet, false)
		err := PreRunLabelUrn(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunLabelUrnErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		err := PreRunLabelUrn(cfg)
		assert.Error(t, err)
	})
}

func TestRunLabelList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		rm.CloudApiV6Mocks.Label.EXPECT().List(gomock.AssignableToTypeOf(testListQueryParam)).Return(testLabels, nil, nil)
		err := RunLabelList(cfg)
		assert.NoError(t, err)
	})
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagResourceType), cloudapiv6.DatacenterResource)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testLabelResourceVar)
		rm.CloudApiV6Mocks.Label.EXPECT().DatacenterList(gomock.AssignableToTypeOf(testListQueryParam), testLabelResourceVar).Return(testLabelResources, nil, nil)
		err := RunLabelList(cfg)
		assert.NoError(t, err)
	})
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagResourceType), cloudapiv6.IpBlockResource)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagIpBlockId), testLabelResourceVar)
		rm.CloudApiV6Mocks.Label.EXPECT().IpBlockList(gomock.AssignableToTypeOf(testListQueryParam), testLabelResourceVar).Return(testLabelResources, nil, nil)
		err := RunLabelList(cfg)
		assert.NoError(t, err)
	})
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagResourceType), cloudapiv6.SnapshotResource)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagSnapshotId), testLabelResourceVar)
		rm.CloudApiV6Mocks.Label.EXPECT().SnapshotList(gomock.AssignableToTypeOf(testListQueryParam), testLabelResourceVar).Return(testLabelResources, nil, nil)
		err := RunLabelList(cfg)
		assert.NoError(t, err)
	})
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagResourceType), cloudapiv6.ServerResource)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagServerId), testLabelResourceVar)
		rm.CloudApiV6Mocks.Label.EXPECT().ServerList(gomock.AssignableToTypeOf(testListQueryParam), testLabelResourceVar, testLabelResourceVar).Return(testLabelResources, nil, nil)
		err := RunLabelList(cfg)
		assert.NoError(t, err)
	})
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagResourceType), cloudapiv6.VolumeResource)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagVolumeId), testLabelResourceVar)
		rm.CloudApiV6Mocks.Label.EXPECT().VolumeList(gomock.AssignableToTypeOf(testListQueryParam), testLabelResourceVar, testLabelResourceVar).Return(testLabelResources, nil, nil)
		err := RunLabelList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunLabelListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		rm.CloudApiV6Mocks.Label.EXPECT().List(gomock.AssignableToTypeOf(testListQueryParam)).Return(testLabels, nil, testLabelErr)
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
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagLabelUrn), testLabelVar)
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
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagLabelUrn), testLabelVar)
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
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagResourceType), cloudapiv6.DatacenterResource)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagLabelKey), testLabelResourceVar)
		rm.CloudApiV6Mocks.Label.EXPECT().DatacenterGet(testLabelResourceVar, testLabelResourceVar).Return(&testLabelResourceRes, nil, nil)
		err := RunLabelGet(cfg)
		assert.NoError(t, err)
	})
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagResourceType), cloudapiv6.IpBlockResource)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagIpBlockId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagLabelKey), testLabelResourceVar)
		rm.CloudApiV6Mocks.Label.EXPECT().IpBlockGet(testLabelResourceVar, testLabelResourceVar).Return(&testLabelResourceRes, nil, nil)
		err := RunLabelGet(cfg)
		assert.NoError(t, err)
	})
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagResourceType), cloudapiv6.SnapshotResource)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagSnapshotId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagLabelKey), testLabelResourceVar)
		rm.CloudApiV6Mocks.Label.EXPECT().SnapshotGet(testLabelResourceVar, testLabelResourceVar).Return(&testLabelResourceRes, nil, nil)
		err := RunLabelGet(cfg)
		assert.NoError(t, err)
	})
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagResourceType), cloudapiv6.ServerResource)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagServerId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagLabelKey), testLabelResourceVar)
		rm.CloudApiV6Mocks.Label.EXPECT().ServerGet(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar).Return(&testLabelResourceRes, nil, nil)
		err := RunLabelGet(cfg)
		assert.NoError(t, err)
	})
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagResourceType), cloudapiv6.VolumeResource)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagVolumeId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagLabelKey), testLabelResourceVar)
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
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagResourceType), cloudapiv6.DatacenterResource)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagLabelKey), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagLabelValue), testLabelResourceVar)
		rm.CloudApiV6Mocks.Label.EXPECT().DatacenterCreate(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar).Return(&testLabelResourceRes, nil, nil)
		err := RunLabelAdd(cfg)
		assert.NoError(t, err)
	})
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagResourceType), cloudapiv6.IpBlockResource)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagIpBlockId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagLabelKey), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagLabelValue), testLabelResourceVar)
		rm.CloudApiV6Mocks.Label.EXPECT().IpBlockCreate(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar).Return(&testLabelResourceRes, nil, nil)
		err := RunLabelAdd(cfg)
		assert.NoError(t, err)
	})
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagResourceType), cloudapiv6.SnapshotResource)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagSnapshotId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagLabelKey), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagLabelValue), testLabelResourceVar)
		rm.CloudApiV6Mocks.Label.EXPECT().SnapshotCreate(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar).Return(&testLabelResourceRes, nil, nil)
		err := RunLabelAdd(cfg)
		assert.NoError(t, err)
	})
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagResourceType), cloudapiv6.ServerResource)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagServerId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagLabelKey), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagLabelValue), testLabelResourceVar)
		rm.CloudApiV6Mocks.Label.EXPECT().ServerCreate(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar, testLabelResourceVar).Return(&testLabelResourceRes, nil, nil)
		err := RunLabelAdd(cfg)
		assert.NoError(t, err)
	})
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagResourceType), cloudapiv6.VolumeResource)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagVolumeId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagLabelKey), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagLabelValue), testLabelResourceVar)
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
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagResourceType), cloudapiv6.DatacenterResource)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagLabelKey), testLabelResourceVar)
		rm.CloudApiV6Mocks.Label.EXPECT().DatacenterDelete(testLabelResourceVar, testLabelResourceVar).Return(nil, nil)
		err := RunLabelRemove(cfg)
		assert.NoError(t, err)
	})
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagResourceType), cloudapiv6.IpBlockResource)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagIpBlockId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagLabelKey), testLabelResourceVar)
		rm.CloudApiV6Mocks.Label.EXPECT().IpBlockDelete(testLabelResourceVar, testLabelResourceVar).Return(nil, nil)
		err := RunLabelRemove(cfg)
		assert.NoError(t, err)
	})
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagResourceType), cloudapiv6.SnapshotResource)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagSnapshotId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagLabelKey), testLabelResourceVar)
		rm.CloudApiV6Mocks.Label.EXPECT().SnapshotDelete(testLabelResourceVar, testLabelResourceVar).Return(nil, nil)
		err := RunLabelRemove(cfg)
		assert.NoError(t, err)
	})
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagResourceType), cloudapiv6.ServerResource)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagServerId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagLabelKey), testLabelResourceVar)
		rm.CloudApiV6Mocks.Label.EXPECT().ServerDelete(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar).Return(nil, nil)
		err := RunLabelRemove(cfg)
		assert.NoError(t, err)
	})
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagResourceType), cloudapiv6.VolumeResource)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagVolumeId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagLabelKey), testLabelResourceVar)
		rm.CloudApiV6Mocks.Label.EXPECT().VolumeDelete(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar).Return(nil, nil)
		err := RunLabelRemove(cfg)
		assert.NoError(t, err)
	})
}
