package image

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestForceConfidentialImageProperties(t *testing.T) {
	p := resources.ImageProperties{}
	forceConfidentialImageProperties(&p)

	assert.Equal(t, "NONE", *p.CloudInit)
	assert.False(t, *p.CpuHotPlug)
	assert.False(t, *p.RamHotPlug)
	assert.False(t, *p.NicHotPlug)
	assert.False(t, *p.DiscVirtioHotPlug)
	assert.False(t, *p.DiscScsiHotPlug)
	assert.False(t, *p.CpuHotUnplug)
	assert.False(t, *p.RamHotUnplug)
	assert.False(t, *p.NicHotUnplug)
	assert.False(t, *p.DiscVirtioHotUnplug)
	assert.False(t, *p.DiscScsiHotUnplug)
	assert.False(t, *p.RequireLegacyBios)
}

func TestAllImageColsHasRequiredFeatures(t *testing.T) {
	var found bool
	for _, c := range allImageCols {
		if c.Name == "RequiredFeatures" {
			found = true
		}
	}
	assert.True(t, found, "allImageCols must expose a RequiredFeatures column")
}

// preRunImageUploadConfidential registers the flags PreRunImageUpload touches, applies the given
// setup, and returns the resulting error.
func runPreRunImageUpload(t *testing.T, setup func(cfg *core.PreCommandConfig)) error {
	t.Helper()
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	var got error
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)

		fs := cfg.Command.Command.Flags()
		fs.StringSlice(FlagImage, nil, "")
		fs.Bool(constants.FlagConfidential, false, "")
		fs.String(constants.FlagCloudInit, "", "")
		fs.Bool(cloudapiv6.ArgCpuHotPlug, false, "")
		fs.Bool(cloudapiv6.ArgRamHotPlug, false, "")
		fs.Bool(cloudapiv6.ArgNicHotPlug, false, "")
		fs.Bool(cloudapiv6.ArgDiscVirtioHotPlug, false, "")
		fs.Bool(cloudapiv6.ArgDiscScsiHotPlug, false, "")
		fs.Bool(cloudapiv6.ArgRequireLegacyBios, false, "")

		setup(cfg)
		got = PreRunImageUpload(cfg)
	})
	return got
}

func TestPreRunImageUploadConfidential(t *testing.T) {
	t.Run("rejects non-qcow2", func(t *testing.T) {
		err := runPreRunImageUpload(t, func(cfg *core.PreCommandConfig) {
			fs := cfg.Command.Command.Flags()
			_ = fs.Set(constants.FlagConfidential, "true")
			_ = fs.Set(FlagImage, "disk.iso")
		})
		assert.Error(t, err)
	})

	t.Run("rejects cloud-init V1", func(t *testing.T) {
		err := runPreRunImageUpload(t, func(cfg *core.PreCommandConfig) {
			fs := cfg.Command.Command.Flags()
			_ = fs.Set(constants.FlagConfidential, "true")
			_ = fs.Set(FlagImage, "disk.qcow2")
			_ = fs.Set(constants.FlagCloudInit, "V1")
		})
		assert.Error(t, err)
	})

	t.Run("rejects hot-plug enabled", func(t *testing.T) {
		err := runPreRunImageUpload(t, func(cfg *core.PreCommandConfig) {
			fs := cfg.Command.Command.Flags()
			_ = fs.Set(constants.FlagConfidential, "true")
			_ = fs.Set(FlagImage, "disk.qcow2")
			_ = fs.Set(cloudapiv6.ArgCpuHotPlug, "true")
		})
		assert.Error(t, err)
	})

	t.Run("rejects require-legacy-bios true", func(t *testing.T) {
		err := runPreRunImageUpload(t, func(cfg *core.PreCommandConfig) {
			fs := cfg.Command.Command.Flags()
			_ = fs.Set(constants.FlagConfidential, "true")
			_ = fs.Set(FlagImage, "disk.qcow2")
			_ = fs.Set(cloudapiv6.ArgRequireLegacyBios, "true")
		})
		assert.Error(t, err)
	})

	t.Run("accepts clean qcow2", func(t *testing.T) {
		err := runPreRunImageUpload(t, func(cfg *core.PreCommandConfig) {
			fs := cfg.Command.Command.Flags()
			_ = fs.Set(constants.FlagConfidential, "true")
			_ = fs.Set(FlagImage, "disk.qcow2")
		})
		assert.NoError(t, err)
	})

	t.Run("no restriction without --confidential", func(t *testing.T) {
		err := runPreRunImageUpload(t, func(cfg *core.PreCommandConfig) {
			_ = cfg.Command.Command.Flags().Set(FlagImage, "disk.iso")
		})
		assert.NoError(t, err)
	})
}
