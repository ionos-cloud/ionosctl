package cloudapi_v5

import (
	"bufio"
	"bytes"
	"errors"
	"testing"

	"github.com/ionos-cloud/ionosctl/internal/config"
	"github.com/ionos-cloud/ionosctl/internal/core"
	cloudapiv5 "github.com/ionos-cloud/ionosctl/services/cloudapi-v5"
	"github.com/ionos-cloud/ionosctl/services/cloudapi-v5/resources"
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

func TestRunDataCenterLabelsList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgDataCenterId), testLabelResourceVar)
		rm.CloudApiV5Mocks.Label.EXPECT().DatacenterList(testLabelResourceVar).Return(testLabelResources, &testResponse, nil)
		err := RunDataCenterLabelsList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunDataCenterLabelsListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgDataCenterId), testLabelResourceVar)
		rm.CloudApiV5Mocks.Label.EXPECT().DatacenterList(testLabelResourceVar).Return(testLabelResources, nil, testLabelResourceErr)
		err := RunDataCenterLabelsList(cfg)
		assert.Error(t, err)
	})
}

func TestRunDataCenterLabelGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgDataCenterId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgLabelKey), testLabelResourceVar)
		rm.CloudApiV5Mocks.Label.EXPECT().DatacenterGet(testLabelResourceVar, testLabelResourceVar).Return(&testLabelResourceRes, &testResponse, nil)
		err := RunDataCenterLabelGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunDataCenterLabelGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgDataCenterId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgLabelKey), testLabelResourceVar)
		rm.CloudApiV5Mocks.Label.EXPECT().DatacenterGet(testLabelResourceVar, testLabelResourceVar).Return(&testLabelResourceRes, nil, testLabelResourceErr)
		err := RunDataCenterLabelGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunDatacenterLabelAdd(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgDataCenterId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgLabelKey), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgLabelValue), testLabelResourceVar)
		rm.CloudApiV5Mocks.Label.EXPECT().DatacenterCreate(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar).Return(&testLabelResourceRes, &testResponse, nil)
		err := RunDataCenterLabelAdd(cfg)
		assert.NoError(t, err)
	})
}

func TestRunDatacenterLabelAddErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgDataCenterId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgLabelKey), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgLabelValue), testLabelResourceVar)
		rm.CloudApiV5Mocks.Label.EXPECT().DatacenterCreate(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar).Return(&testLabelResourceRes, nil, testLabelResourceErr)
		err := RunDataCenterLabelAdd(cfg)
		assert.Error(t, err)
	})
}

func TestRunDatacenterLabelRemove(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgDataCenterId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgLabelKey), testLabelResourceVar)
		rm.CloudApiV5Mocks.Label.EXPECT().DatacenterDelete(testLabelResourceVar, testLabelResourceVar).Return(&testResponse, nil)
		err := RunDataCenterLabelRemove(cfg)
		assert.NoError(t, err)
	})
}

func TestRunDatacenterLabelRemoveErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgDataCenterId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgLabelKey), testLabelResourceVar)
		rm.CloudApiV5Mocks.Label.EXPECT().DatacenterDelete(testLabelResourceVar, testLabelResourceVar).Return(nil, testLabelResourceErr)
		err := RunDataCenterLabelRemove(cfg)
		assert.Error(t, err)
	})
}

func TestRunIpBlockLabelsList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgIpBlockId), testLabelResourceVar)
		rm.CloudApiV5Mocks.Label.EXPECT().IpBlockList(testLabelResourceVar).Return(testLabelResources, &testResponse, nil)
		err := RunIpBlockLabelsList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunIpBlockLabelsListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgIpBlockId), testLabelResourceVar)
		rm.CloudApiV5Mocks.Label.EXPECT().IpBlockList(testLabelResourceVar).Return(testLabelResources, nil, testLabelResourceErr)
		err := RunIpBlockLabelsList(cfg)
		assert.Error(t, err)
	})
}

func TestRunIpBlockLabelGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgIpBlockId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgLabelKey), testLabelResourceVar)
		rm.CloudApiV5Mocks.Label.EXPECT().IpBlockGet(testLabelResourceVar, testLabelResourceVar).Return(&testLabelResourceRes, &testResponse, nil)
		err := RunIpBlockLabelGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunIpBlockLabelGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgIpBlockId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgLabelKey), testLabelResourceVar)
		rm.CloudApiV5Mocks.Label.EXPECT().IpBlockGet(testLabelResourceVar, testLabelResourceVar).Return(&testLabelResourceRes, nil, testLabelResourceErr)
		err := RunIpBlockLabelGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunIpBlockLabelAdd(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgIpBlockId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgLabelKey), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgLabelValue), testLabelResourceVar)
		rm.CloudApiV5Mocks.Label.EXPECT().IpBlockCreate(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar).Return(&testLabelResourceRes, &testResponse, nil)
		err := RunIpBlockLabelAdd(cfg)
		assert.NoError(t, err)
	})
}

func TestRunIpBlockLabelAddErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgIpBlockId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgLabelKey), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgLabelValue), testLabelResourceVar)
		rm.CloudApiV5Mocks.Label.EXPECT().IpBlockCreate(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar).Return(&testLabelResourceRes, nil, testLabelResourceErr)
		err := RunIpBlockLabelAdd(cfg)
		assert.Error(t, err)
	})
}

func TestRunIpBlockLabelRemove(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgIpBlockId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgLabelKey), testLabelResourceVar)
		rm.CloudApiV5Mocks.Label.EXPECT().IpBlockDelete(testLabelResourceVar, testLabelResourceVar).Return(&testResponse, nil)
		err := RunIpBlockLabelRemove(cfg)
		assert.NoError(t, err)
	})
}

func TestRunIpBlockLabelRemoveErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgIpBlockId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgLabelKey), testLabelResourceVar)
		rm.CloudApiV5Mocks.Label.EXPECT().IpBlockDelete(testLabelResourceVar, testLabelResourceVar).Return(nil, testLabelResourceErr)
		err := RunIpBlockLabelRemove(cfg)
		assert.Error(t, err)
	})
}

func TestRunSnapshotLabelsList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgSnapshotId), testLabelResourceVar)
		rm.CloudApiV5Mocks.Label.EXPECT().SnapshotList(testLabelResourceVar).Return(testLabelResources, &testResponse, nil)
		err := RunSnapshotLabelsList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunSnapshotLabelsListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgSnapshotId), testLabelResourceVar)
		rm.CloudApiV5Mocks.Label.EXPECT().SnapshotList(testLabelResourceVar).Return(testLabelResources, nil, testLabelResourceErr)
		err := RunSnapshotLabelsList(cfg)
		assert.Error(t, err)
	})
}

func TestRunSnapshotLabelGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgSnapshotId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgLabelKey), testLabelResourceVar)
		rm.CloudApiV5Mocks.Label.EXPECT().SnapshotGet(testLabelResourceVar, testLabelResourceVar).Return(&testLabelResourceRes, &testResponse, nil)
		err := RunSnapshotLabelGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunSnapshotLabelGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgSnapshotId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgLabelKey), testLabelResourceVar)
		rm.CloudApiV5Mocks.Label.EXPECT().SnapshotGet(testLabelResourceVar, testLabelResourceVar).Return(&testLabelResourceRes, nil, testLabelResourceErr)
		err := RunSnapshotLabelGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunSnapshotLabelAdd(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgSnapshotId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgLabelKey), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgLabelValue), testLabelResourceVar)
		rm.CloudApiV5Mocks.Label.EXPECT().SnapshotCreate(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar).Return(&testLabelResourceRes, &testResponse, nil)
		err := RunSnapshotLabelAdd(cfg)
		assert.NoError(t, err)
	})
}

func TestRunSnapshotLabelAddErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgSnapshotId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgLabelKey), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgLabelValue), testLabelResourceVar)
		rm.CloudApiV5Mocks.Label.EXPECT().SnapshotCreate(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar).Return(&testLabelResourceRes, nil, testLabelResourceErr)
		err := RunSnapshotLabelAdd(cfg)
		assert.Error(t, err)
	})
}

func TestRunSnapshotLabelRemove(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgSnapshotId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgLabelKey), testLabelResourceVar)
		rm.CloudApiV5Mocks.Label.EXPECT().SnapshotDelete(testLabelResourceVar, testLabelResourceVar).Return(&testResponse, nil)
		err := RunSnapshotLabelRemove(cfg)
		assert.NoError(t, err)
	})
}

func TestRunSnapshotLabelRemoveErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgSnapshotId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgLabelKey), testLabelResourceVar)
		rm.CloudApiV5Mocks.Label.EXPECT().SnapshotDelete(testLabelResourceVar, testLabelResourceVar).Return(nil, testLabelResourceErr)
		err := RunSnapshotLabelRemove(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerLabelsList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgDataCenterId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgServerId), testLabelResourceVar)
		rm.CloudApiV5Mocks.Label.EXPECT().ServerList(testLabelResourceVar, testLabelResourceVar).Return(testLabelResources, &testResponse, nil)
		err := RunServerLabelsList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunServerLabelsListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgDataCenterId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgServerId), testLabelResourceVar)
		rm.CloudApiV5Mocks.Label.EXPECT().ServerList(testLabelResourceVar, testLabelResourceVar).Return(testLabelResources, nil, testLabelResourceErr)
		err := RunServerLabelsList(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerLabelGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgDataCenterId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgServerId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgLabelKey), testLabelResourceVar)
		rm.CloudApiV5Mocks.Label.EXPECT().ServerGet(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar).Return(&testLabelResourceRes, &testResponse, nil)
		err := RunServerLabelGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunServerLabelGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgDataCenterId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgServerId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgLabelKey), testLabelResourceVar)
		rm.CloudApiV5Mocks.Label.EXPECT().ServerGet(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar).Return(&testLabelResourceRes, nil, testLabelResourceErr)
		err := RunServerLabelGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerLabelAdd(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgDataCenterId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgServerId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgLabelKey), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgLabelValue), testLabelResourceVar)
		rm.CloudApiV5Mocks.Label.EXPECT().ServerCreate(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar, testLabelResourceVar).
			Return(&testLabelResourceRes, &testResponse, nil)
		err := RunServerLabelAdd(cfg)
		assert.NoError(t, err)
	})
}

func TestRunServerLabelAddErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgDataCenterId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgServerId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgLabelKey), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgLabelValue), testLabelResourceVar)
		rm.CloudApiV5Mocks.Label.EXPECT().ServerCreate(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar, testLabelResourceVar).Return(&testLabelResourceRes, nil, testLabelResourceErr)
		err := RunServerLabelAdd(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerLabelRemove(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgDataCenterId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgServerId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgLabelKey), testLabelResourceVar)
		rm.CloudApiV5Mocks.Label.EXPECT().ServerDelete(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar).Return(&testResponse, nil)
		err := RunServerLabelRemove(cfg)
		assert.NoError(t, err)
	})
}

func TestRunServerLabelRemoveErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgDataCenterId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgServerId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgLabelKey), testLabelResourceVar)
		rm.CloudApiV5Mocks.Label.EXPECT().ServerDelete(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar).Return(nil, testLabelResourceErr)
		err := RunServerLabelRemove(cfg)
		assert.Error(t, err)
	})
}

func TestRunVolumeLabelsList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgDataCenterId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgVolumeId), testLabelResourceVar)
		rm.CloudApiV5Mocks.Label.EXPECT().VolumeList(testLabelResourceVar, testLabelResourceVar).Return(testLabelResources, &testResponse, nil)
		err := RunVolumeLabelsList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunVolumeLabelsListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgDataCenterId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgVolumeId), testLabelResourceVar)
		rm.CloudApiV5Mocks.Label.EXPECT().VolumeList(testLabelResourceVar, testLabelResourceVar).Return(testLabelResources, nil, testLabelResourceErr)
		err := RunVolumeLabelsList(cfg)
		assert.Error(t, err)
	})
}

func TestRunVolumeLabelGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgDataCenterId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgVolumeId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgLabelKey), testLabelResourceVar)
		rm.CloudApiV5Mocks.Label.EXPECT().VolumeGet(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar).Return(&testLabelResourceRes, &testResponse, nil)
		err := RunVolumeLabelGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunVolumeLabelGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgDataCenterId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgVolumeId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgLabelKey), testLabelResourceVar)
		rm.CloudApiV5Mocks.Label.EXPECT().VolumeGet(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar).Return(&testLabelResourceRes, nil, testLabelResourceErr)
		err := RunVolumeLabelGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunVolumeLabelAdd(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgDataCenterId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgVolumeId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgLabelKey), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgLabelValue), testLabelResourceVar)
		rm.CloudApiV5Mocks.Label.EXPECT().VolumeCreate(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar, testLabelResourceVar).
			Return(&testLabelResourceRes, &testResponse, nil)
		err := RunVolumeLabelAdd(cfg)
		assert.NoError(t, err)
	})
}

func TestRunVolumeLabelAddErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgDataCenterId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgVolumeId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgLabelKey), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgLabelValue), testLabelResourceVar)
		rm.CloudApiV5Mocks.Label.EXPECT().VolumeCreate(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar, testLabelResourceVar).Return(&testLabelResourceRes, nil, testLabelResourceErr)
		err := RunVolumeLabelAdd(cfg)
		assert.Error(t, err)
	})
}

func TestRunVolumeLabelRemove(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgDataCenterId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgVolumeId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgLabelKey), testLabelResourceVar)
		rm.CloudApiV5Mocks.Label.EXPECT().VolumeDelete(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar).Return(&testResponse, nil)
		err := RunVolumeLabelRemove(cfg)
		assert.NoError(t, err)
	})
}

func TestRunVolumeLabelRemoveErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgDataCenterId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgVolumeId), testLabelResourceVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgLabelKey), testLabelResourceVar)
		rm.CloudApiV5Mocks.Label.EXPECT().VolumeDelete(testLabelResourceVar, testLabelResourceVar, testLabelResourceVar).Return(nil, testLabelResourceErr)
		err := RunVolumeLabelRemove(cfg)
		assert.Error(t, err)
	})
}
