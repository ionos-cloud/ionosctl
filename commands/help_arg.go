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
// "runnable command followed by a trailing 'help'" case and prints that
// command's help instead.
//
// root.Find strips the command-path tokens but keeps flags and their values, so
// the leftover args ("rest") look like ["--name", "foo", "help"] for
// "server create --name foo help". We therefore key off the *last* leftover
// token rather than requiring "help" to be the sole leftover. A trailing "help"
// always wins: to pass the literal string "help" as a value put it anywhere but
// last (e.g. "--name=help" or "--name help --enabled").
//
// It returns true when it handled the invocation, in which case the caller
// should skip normal execution.
func handleTrailingHelp(root *cobra.Command, args []string) bool {
	if len(args) == 0 {
		return false
	}

	cmd, rest, err := root.Find(args)
	if err != nil || cmd == nil || !cmd.Runnable() {
		return false
	}

	if len(rest) == 0 || rest[len(rest)-1] != "help" {
		return false
	}

	_ = cmd.Help()
	return true
}
