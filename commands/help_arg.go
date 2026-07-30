package commands

import (
	"github.com/spf13/cobra"
)

// handleTrailingHelp makes "ionosctl <command> help" behave like
// "ionosctl <command> --help".
//
// Cobra only treats "help" as help when it is the first token (the root `help`
// command). A trailing "help" after a runnable command is otherwise consumed as
// a positional argument, so "ionosctl server create help" tried to run create
// and failed with a confusing "missing required flags" error. This detects the
// exact "runnable command followed by a lone 'help'" case and prints that
// command's help instead.
//
// It returns true when it handled the invocation, in which case the caller
// should skip normal execution.
func handleTrailingHelp(root *cobra.Command, args []string) bool {
	if len(args) == 0 {
		return false
	}

	cmd, rest, err := root.Find(args)
	if err != nil || cmd == nil {
		return false
	}

	// Only intercept when the sole leftover token is "help" and the resolved
	// command actually runs (i.e. help would otherwise be a stray positional).
	if len(rest) == 1 && rest[0] == "help" && cmd.Runnable() {
		_ = cmd.Help()
		return true
	}

	return false
}
