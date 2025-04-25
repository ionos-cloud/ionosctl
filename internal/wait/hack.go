package wait

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// this file contains various 'hacks' for integrating global 'wait' flag

// AddTimeoutFlag adds --timeout/-t flag to every command.
// We are not adding a global flag, but instead iterating through each command and adding it separately
// because some commands already define '-t' shorthand.
// if the '-t' shorthand already exists, skip it.
func AddTimeoutFlag(root *cobra.Command) {
	for _, cmd := range root.Commands() {
		AddTimeoutFlag(cmd)

		if cmd.Flags().Lookup(constants.ArgTimeout) != nil {
			continue
		}

		if cmd.Flags().ShorthandLookup(constants.ArgTimeoutShort) == nil {
			cmd.Flags().IntP(constants.ArgTimeout, constants.ArgTimeoutShort, constants.DefaultTimeoutSeconds,
				"Timeout for waiting for resource to reach desired state")
		} else {
			cmd.Flags().Int(constants.ArgTimeout, constants.DefaultTimeoutSeconds,
				"Timeout for waiting for resource to reach desired state")
		}

		viper.BindPFlag(constants.ArgTimeout, cmd.Flags().Lookup(constants.ArgTimeout))
	}
}
