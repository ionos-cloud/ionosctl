package wait

import (
	"time"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/spf13/cobra"
)

// this file contains various 'hacks' for integrating global 'wait' flag

// AddTimeoutFlag adds --timeout/-t flag to every command.
// We are not adding a global flag, but instead iterating through each command and adding it separately
// because some commands already define '-t' shorthand.
// if the '-t' shorthand already exists, skip it.
func AddTimeoutFlag(root *cobra.Command) {
	for _, cmd := range root.Commands() {
		AddTimeoutFlag(cmd)

		if cmd.Flags().Lookup("timeout") != nil {
			continue
		}

		if cmd.Flags().ShorthandLookup("t") == nil {
			cmd.Flags().DurationP("timeout", "t", time.Duration(constants.DefaultTimeoutSeconds),
				"Timeout for waiting for resource to reach desired state")
		} else {
			cmd.Flags().Duration("timeout", time.Duration(constants.DefaultTimeoutSeconds),
				"Timeout for waiting for resource to reach desired state")
		}
	}
}
