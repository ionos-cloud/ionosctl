package server

import (
	"bufio"
	"bytes"
	"strconv"
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

func TestAllServerColsHasEnabledFeatures(t *testing.T) {
	var found bool
	for _, c := range AllServerCols {
		if c.Name == "EnabledFeatures" {
			found = true
		}
	}
	assert.True(t, found, "AllServerCols must expose an EnabledFeatures column")
}

// getNewServer must not set cores/cpuFamily for a Confidential VM — both are derived from the image.
func TestGetNewServerConfidential(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagType), testServerEnterpriseType)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagConfidential), true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagRam), strconv.Itoa(int(ram)))
		// Set cores/cpu-family in viper to prove getNewServer IGNORES them when confidential.
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagCores), cores)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagCpuFamily), testServerVar)

		srv, err := getNewServer(cfg)
		assert.NoError(t, err)
		assert.Equal(t, testServerEnterpriseType, *srv.Properties.Type)
		assert.Nil(t, srv.Properties.Cores, "cores must be derived from the image")
		assert.Nil(t, srv.Properties.CpuFamily, "cpuFamily must be derived from the image")
	})
}

// getNewDAS builds the confidential boot volume with the requested storage type and size.
func TestGetNewDASConfidential(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagType), testServerEnterpriseType)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagConfidential), true)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagStorageType), "SSD")
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgSize), "20")
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgImageId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgVolumeName), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgBus), "VIRTIO")

		vol, err := getNewDAS(cfg)
		assert.NoError(t, err)
		assert.Equal(t, "SSD", *vol.Properties.Type)
		assert.Equal(t, float32(20), *vol.Properties.Size)
		assert.Equal(t, testServerVar, *vol.Properties.Image)
	})
}

func TestPreRunServerCreateConfidentialErrors(t *testing.T) {
	tests := []struct {
		name         string
		serverType   string
		setCores     bool
		setCpuFamily bool
		setImage     bool
	}{
		{name: "non-enterprise type", serverType: serverVCPUType, setImage: true},
		{name: "cores set", serverType: serverEnterpriseType, setCores: true, setImage: true},
		{name: "cpu-family set", serverType: serverEnterpriseType, setCpuFamily: true, setImage: true},
		{name: "missing image", serverType: serverEnterpriseType},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var b bytes.Buffer
			w := bufio.NewWriter(&b)
			core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
				viper.Reset()
				viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)

				fs := cfg.Command.Command.Flags()
				fs.Int(constants.FlagCores, 0, "")
				fs.String(constants.FlagCpuFamily, "", "")
				fs.String(cloudapiv6.ArgImageId, "", "")

				viper.Set(core.GetFlagName(cfg.NS, constants.FlagConfidential), true)
				viper.Set(core.GetFlagName(cfg.NS, constants.FlagType), tt.serverType)

				if tt.setCores {
					_ = fs.Set(constants.FlagCores, "4")
				}
				if tt.setCpuFamily {
					_ = fs.Set(constants.FlagCpuFamily, "INTEL_ICELAKE")
				}
				if tt.setImage {
					_ = fs.Set(cloudapiv6.ArgImageId, testServerVar)
				}

				err := PreRunServerCreate(cfg)
				assert.Error(t, err)
			})
		})
	}
}

// RunServerCreate must attach a boot volume built from the confidential image, and must not set
// cores/cpuFamily (the API derives them from the image).
func TestRunServerCreateConfidentialAttachesVolume(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagType), testServerEnterpriseType)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagConfidential), true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgImageId), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgVolumeName), testServerVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagStorageType), "HDD")
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgSize), "10")
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagRam), strconv.Itoa(int(ram)))
		viper.Set(constants.ArgWait, false)

		rm.CloudApiV6Mocks.Server.EXPECT().Create(testServerVar, gomock.Any()).DoAndReturn(
			func(dcId string, input resources.Server) (*resources.Server, *resources.Response, error) {
				assert.Nil(t, input.Properties.Cores, "cores must not be set for confidential")
				assert.Nil(t, input.Properties.CpuFamily, "cpuFamily must not be set for confidential")
				if assert.NotNil(t, input.Entities) && assert.NotNil(t, input.Entities.Volumes) {
					items := *input.Entities.Volumes.Items
					assert.Len(t, items, 1)
					assert.Equal(t, testServerVar, *items[0].Properties.Image)
				}
				return &resources.Server{Server: ionoscloud.Server{Properties: input.Properties}}, &testutil.TestResponse, nil
			})

		err := RunServerCreate(cfg)
		assert.NoError(t, err)
	})
}
