package volume

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/ionos-cloud/ionosctl/v6/commands/compute/testutil"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	sizeVolume        = float32(12)
	zoneVolume        = "ZONE_1"
	testServerVar     = "test-server"
	testVolumeVar     = "test-volume"
	testResourceVar   = "test-resource"
	testVolumeBoolVar = false
	testVolumeErr     = errors.New("volume test: error occurred")
	v                 = ionoscloud.Volume{
		Id: &testVolumeVar,
		Properties: &ionoscloud.VolumeProperties{
			Name:                &testVolumeVar,
			Size:                &sizeVolume,
			LicenceType:         &testVolumeVar,
			Type:                &testVolumeVar,
			Bus:                 &testVolumeVar,
			Image:               &testVolumeVar,
			ImageAlias:          &testVolumeVar,
			AvailabilityZone:    &zoneVolume,
			BackupunitId:        &testVolumeVar,
			UserData:            &testVolumeVar,
			CpuHotPlug:          &testVolumeBoolVar,
			RamHotPlug:          &testVolumeBoolVar,
			NicHotPlug:          &testVolumeBoolVar,
			NicHotUnplug:        &testVolumeBoolVar,
			DiscVirtioHotPlug:   &testVolumeBoolVar,
			DiscVirtioHotUnplug: &testVolumeBoolVar,
			BootServer:          &testVolumeVar,
		},
		Metadata: &ionoscloud.DatacenterElementMetadata{
			State: &testVolumeVar,
		},
	}
	serverVolume = ionoscloud.Volume{
		Id: &testServerVar,
		Properties: &ionoscloud.VolumeProperties{
			Name:                &testVolumeVar,
			Size:                &sizeVolume,
			LicenceType:         &testVolumeVar,
			Type:                &testVolumeVar,
			Bus:                 &testVolumeVar,
			Image:               &testVolumeVar,
			ImageAlias:          &testVolumeVar,
			AvailabilityZone:    &zoneVolume,
			BackupunitId:        &testVolumeVar,
			UserData:            &testVolumeVar,
			CpuHotPlug:          &testVolumeBoolVar,
			RamHotPlug:          &testVolumeBoolVar,
			NicHotPlug:          &testVolumeBoolVar,
			NicHotUnplug:        &testVolumeBoolVar,
			DiscVirtioHotPlug:   &testVolumeBoolVar,
			DiscVirtioHotUnplug: &testVolumeBoolVar,
			BootServer:          &testVolumeVar,
		},
		Metadata: &ionoscloud.DatacenterElementMetadata{
			State: &testVolumeVar,
		},
	}
	vsAttached = resources.AttachedVolumes{
		AttachedVolumes: ionoscloud.AttachedVolumes{
			Id:    &testVolumeVar,
			Items: &[]ionoscloud.Volume{v},
		},
	}
	vsAttachedList = resources.AttachedVolumes{
		AttachedVolumes: ionoscloud.AttachedVolumes{
			Id:    &testVolumeVar,
			Items: &[]ionoscloud.Volume{serverVolume, serverVolume},
		},
	}
)

func TestPreRunDcServerIdsRequiredFlagsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		err := PreRunDcServerIds(cfg)
		assert.Error(t, err)
	})
}

// PreRunDcServerIds checks required datacenter and server ID flags.
func PreRunDcServerIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgDataCenterId, cloudapiv6.ArgServerId)
}

func TestPreRunDcServerVolumeIds(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgVolumeId), testServerVar)
		err := PreRunDcServerVolumeIds(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunServerVolumeList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar)
		err := PreRunServerVolumeList(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunServerVolumeListFilters(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar)
		cfg.Command.Command.Flags().Set(constants.FlagFilters, fmt.Sprintf("createdby=%s", testutil.TestQueryParamVar))
		err := PreRunServerVolumeList(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunServerVolumeListFiltersErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar)
		cfg.Command.Command.Flags().Set(constants.FlagFilters, fmt.Sprintf("%s=%s", testutil.TestQueryParamVar, testutil.TestQueryParamVar))
		err := PreRunServerVolumeList(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunDcServerVolumeIdsRequiredFlagsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		err := PreRunDcServerVolumeIds(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerVolumeAttach(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgVolumeId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.Server.EXPECT().AttachVolume(testServerVar, testServerVar, testServerVar).Return(&resources.Volume{Volume: v}, nil, nil)
		err := RunServerVolumeAttach(cfg)
		assert.NoError(t, err)
	})
}

func TestRunServerVolumeAttachErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgVolumeId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.Server.EXPECT().AttachVolume(testServerVar, testServerVar, testServerVar).Return(&resources.Volume{Volume: v}, nil, testVolumeErr)
		err := RunServerVolumeAttach(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerVolumeAttachWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgVolumeId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), true)
		rm.CloudApiV6Mocks.Server.EXPECT().AttachVolume(testServerVar, testServerVar, testServerVar).Return(&resources.Volume{Volume: v}, &testutil.TestResponse, nil)
		rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testutil.TestRequestIdVar).Return(&testutil.TestRequestStatus, nil, testutil.TestRequestErr)
		err := RunServerVolumeAttach(cfg)
		assert.Error(t, err)
	})
}

func TestServerVolumeList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	var funcToTest = RunServerVolumesList
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		var tests = []core.TestCase{
			{
				Name: "server volume list",
				Args: []core.FlagValuePair{
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar},
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar},
				},
				Calls: func(...*gomock.Call) {
					gomock.InOrder(
						rm.CloudApiV6Mocks.Server.EXPECT().ListVolumes(gomock.AssignableToTypeOf(testResourceVar), gomock.AssignableToTypeOf(testResourceVar)).Return(vsAttached, nil, nil),
					)
				},
				ExpectedErr: false,
			},
			{
				Name: "server volume list error",
				Args: []core.FlagValuePair{
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar},
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar},
				},
				Calls: func(...*gomock.Call) {
					gomock.InOrder(
						rm.CloudApiV6Mocks.Server.EXPECT().ListVolumes(gomock.AssignableToTypeOf(testResourceVar), gomock.AssignableToTypeOf(testResourceVar)).Return(vsAttached, nil, testVolumeErr),
					)
				},
				ExpectedErr: true,
			},
		}
		core.ExecuteTestCases(t, funcToTest, tests, cfg)
	})
}

func TestServerVolumeGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	var funcToTest = RunServerVolumeGet
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		var tests = []core.TestCase{
			{
				Name: "server volume get",
				Args: []core.FlagValuePair{
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar},
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar},
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgVolumeId), testServerVar},
				},
				Calls: func(...*gomock.Call) {
					gomock.InOrder(
						rm.CloudApiV6Mocks.Server.EXPECT().GetVolume(testServerVar, testServerVar, testServerVar).Return(&resources.Volume{Volume: v}, nil, nil),
					)
				},
				ExpectedErr: false,
			},
			{
				Name: "server volume get error",
				Args: []core.FlagValuePair{
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar},
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar},
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgVolumeId), testServerVar},
				},
				Calls: func(...*gomock.Call) {
					gomock.InOrder(
						rm.CloudApiV6Mocks.Server.EXPECT().GetVolume(testServerVar, testServerVar, testServerVar).Return(&resources.Volume{Volume: v}, nil, testVolumeErr),
					)
				},
				ExpectedErr: true,
			},
		}
		core.ExecuteTestCases(t, funcToTest, tests, cfg)
	})
}

func TestServerVolumeDetach(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	var funcToTest = RunServerVolumeDetach
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		var tests = []core.TestCase{
			{
				Name: "server volume detach",
				Args: []core.FlagValuePair{
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar},
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar},
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgVolumeId), testServerVar},
					{core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false},
					{constants.ArgForce, true},
				},
				Calls: func(...*gomock.Call) {
					gomock.InOrder(
						rm.CloudApiV6Mocks.Server.EXPECT().DetachVolume(testServerVar, testServerVar, testServerVar).Return(nil, nil),
					)
				},
				ExpectedErr: false,
			},
			{
				Name: "server volume detach error",
				Args: []core.FlagValuePair{
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar},
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar},
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgVolumeId), testServerVar},
					{core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false},
					{constants.ArgForce, true},
				},
				Calls: func(...*gomock.Call) {
					gomock.InOrder(
						rm.CloudApiV6Mocks.Server.EXPECT().DetachVolume(testServerVar, testServerVar, testServerVar).Return(&testutil.TestResponse, testVolumeErr),
					)
				},
				ExpectedErr: true,
			},
			{
				Name: "server volume detach all",
				Args: []core.FlagValuePair{
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar},
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar},
					{core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false},
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true},
					{constants.ArgForce, true},
				},
				Calls: func(...*gomock.Call) {
					gomock.InOrder(
						rm.CloudApiV6Mocks.Server.EXPECT().ListVolumes(testServerVar, testServerVar).Return(vsAttachedList, nil, nil),
						rm.CloudApiV6Mocks.Server.EXPECT().DetachVolume(testServerVar, testServerVar, testServerVar).Return(&testutil.TestResponse, nil),
						rm.CloudApiV6Mocks.Server.EXPECT().DetachVolume(testServerVar, testServerVar, testServerVar).Return(&testutil.TestResponse, nil),
					)
				},
				ExpectedErr: false,
			},
			{
				Name: "server volume detach all (error)",
				Args: []core.FlagValuePair{
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar},
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar},
					{core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false},
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true},
					{constants.ArgForce, true},
				},
				Calls: func(...*gomock.Call) {
					rm.CloudApiV6Mocks.Server.EXPECT().ListVolumes(testServerVar, testServerVar).Return(vsAttachedList, nil, nil)
					rm.CloudApiV6Mocks.Server.EXPECT().DetachVolume(testServerVar, testServerVar, testServerVar).Return(&testutil.TestResponse, testVolumeErr)
					rm.CloudApiV6Mocks.Server.EXPECT().DetachVolume(testServerVar, testServerVar, testServerVar).Return(&testutil.TestResponse, nil)
				},
				ExpectedErr: true,
			},
			{
				Name: "server volume detach all (list error)",
				Args: []core.FlagValuePair{
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar},
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar},
					{core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false},
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true},
					{constants.ArgForce, true},
				},
				Calls: func(...*gomock.Call) {
					gomock.InOrder(
						rm.CloudApiV6Mocks.Server.EXPECT().ListVolumes(
							testServerVar,
							testServerVar,
						).Return(
							vsAttachedList,
							nil,
							testVolumeErr,
						),
					)
				},
				ExpectedErr: true,
			},
			{
				Name: "server volume detach all (wrong items error)",
				Args: []core.FlagValuePair{
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar},
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar},
					{core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false},
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true},
					{constants.ArgForce, true},
				},
				Calls: func(...*gomock.Call) {
					gomock.InOrder(
						rm.CloudApiV6Mocks.Server.EXPECT().ListVolumes(
							testServerVar,
							testServerVar,
						).Return(
							vsAttachedList,
							nil,
							testVolumeErr,
						),
					)
				},
				ExpectedErr: true,
			},
			{
				Name: "server volume detach all (length error)",
				Args: []core.FlagValuePair{
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar},
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar},
					{core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false},
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true},
					{constants.ArgForce, true},
				},
				Calls: func(...*gomock.Call) {
					gomock.InOrder(
						rm.CloudApiV6Mocks.Server.EXPECT().ListVolumes(
							testServerVar,
							testServerVar,
						).Return(
							resources.AttachedVolumes{AttachedVolumes: ionoscloud.AttachedVolumes{Items: &[]ionoscloud.Volume{}}},
							nil,
							nil,
						),
					)
				},
				ExpectedErr: true,
			},
			{
				Name: "server volume detach (user confirm)",
				Args: []core.FlagValuePair{
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar},
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar},
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgVolumeId), testServerVar},
					{core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false},
					{constants.ArgForce, false},
				},
				Calls: func(...*gomock.Call) {
					gomock.InOrder(
						rm.CloudApiV6Mocks.Server.EXPECT().DetachVolume(
							testServerVar,
							testServerVar,
							testServerVar,
						).Return(
							nil,
							nil,
						),
					)
				},
				UserInput:   bytes.NewReader([]byte("YES\n")),
				ExpectedErr: false,
			},
			{
				Name: "server volume detach (wait error)",
				Args: []core.FlagValuePair{
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar},
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar},
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgVolumeId), testServerVar},
					{core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), true},
					{constants.ArgForce, true},
				},
				Calls: func(...*gomock.Call) {
					gomock.InOrder(
						rm.CloudApiV6Mocks.Server.EXPECT().DetachVolume(testServerVar, testServerVar, testServerVar).Return(&testutil.TestResponse, nil),
						rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testutil.TestRequestIdVar).Return(&testutil.TestRequestStatus, nil, testutil.TestRequestErr),
					)
				},
				ExpectedErr: true,
			},
			{
				Name: "server volume detach (user confirm error)",
				Args: []core.FlagValuePair{
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar},
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testServerVar},
					{core.GetFlagName(cfg.NS, cloudapiv6.ArgVolumeId), testServerVar},
					{core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false},
					{constants.ArgForce, false},
				},
				Calls: func(...*gomock.Call) {
					gomock.InOrder()
				},
				UserInput:   bytes.NewReader([]byte("\n")),
				ExpectedErr: true,
			},
		}
		core.ExecuteTestCases(t, funcToTest, tests, cfg)
	})
}
