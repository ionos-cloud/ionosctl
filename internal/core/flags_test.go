package core

import (
	"testing"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

// c.Flags() must see global persistent flags inherited from a parent command (e.g.
// the root's --timeout), not just flags defined locally. cobra merges parent
// persistent flags into cmd.Flags() during ParseFlags, which runs before PreRun/Run,
// so reads inside command handlers resolve inherited globals correctly.
func TestFlagsReadsInheritedGlobalFlag(t *testing.T) {
	newTree := func() (*cobra.Command, *CommandConfig) {
		root := &cobra.Command{Use: "root"}
		root.PersistentFlags().Int(constants.ArgTimeout, constants.DefaultTimeoutSeconds, "")
		leaf := &cobra.Command{Use: "leaf", Run: func(*cobra.Command, []string) {}}
		root.AddCommand(leaf)
		return root, &CommandConfig{Command: &Command{Command: leaf}}
	}

	t.Run("explicit value", func(t *testing.T) {
		root, c := newTree()
		root.SetArgs([]string{"leaf", "--" + constants.ArgTimeout, "42"})
		assert.NoError(t, root.Execute())
		assert.Equal(t, 42, c.Flags().Int(constants.ArgTimeout))
	})

	t.Run("falls back to inherited default when unset", func(t *testing.T) {
		root, c := newTree()
		root.SetArgs([]string{"leaf"})
		assert.NoError(t, root.Execute())
		assert.Equal(t, constants.DefaultTimeoutSeconds, c.Flags().Int(constants.ArgTimeout))
	})
}

// c.Flags() getters coerce a flag's value to the requested type the way viper did,
// e.g. reading an int-typed flag as a string.
func TestFlagsCoercesAcrossTypes(t *testing.T) {
	cmd := &cobra.Command{Use: "x"}
	cmd.Flags().Int("port", 0, "")
	_ = cmd.Flags().Set("port", "8080")
	c := &CommandConfig{Command: &Command{Command: cmd}}

	assert.Equal(t, "8080", c.Flags().String("port")) // int flag read as string
	assert.Equal(t, 8080, c.Flags().Int("port"))      // native typed read
	assert.Equal(t, "", c.Flags().String("missing"))  // undefined -> zero value
}
