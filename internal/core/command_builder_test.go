package core

import (
	"testing"

	"github.com/ionos-cloud/ionosctl/v6/pkg/functional"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

// Regression: the deprecated->canonical flag alias must be visible to commands that
// read the canonical flag via c.Flags() (cobra), not viper. See PreRunWithDeprecatedFlags.
func TestPreRunWithDeprecatedFlags_AliasesValueToCanonicalFlag(t *testing.T) {
	const deprecated, canonical = "rename-images", "image-alias"

	newCfg := func(setDeprecated bool) *PreCommandConfig {
		cmd := &cobra.Command{Use: testConst}
		cmd.Flags().StringSlice(deprecated, nil, "")
		cmd.Flags().StringSlice(canonical, nil, "")
		if setDeprecated {
			_ = cmd.Flags().Set(deprecated, "a")
			_ = cmd.Flags().Set(deprecated, "b")
		}
		return &PreCommandConfig{Command: &Command{Command: cmd}}
	}

	tuple := functional.Tuple[string]{First: deprecated, Second: canonical}

	t.Run("deprecated set -> canonical sees value", func(t *testing.T) {
		c := newCfg(true)
		err := PreRunWithDeprecatedFlags(NoPreRun, tuple)(c)
		assert.NoError(t, err)
		assert.True(t, c.Flags().Changed(canonical))
		assert.Equal(t, []string{"a", "b"}, c.Flags().StringSlice(canonical))
	})

	t.Run("deprecated unset -> canonical untouched", func(t *testing.T) {
		c := newCfg(false)
		err := PreRunWithDeprecatedFlags(NoPreRun, tuple)(c)
		assert.NoError(t, err)
		assert.False(t, c.Flags().Changed(canonical))
	})
}
