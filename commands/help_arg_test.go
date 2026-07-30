package commands

import (
	"testing"

	"github.com/spf13/cobra"
)

func helpArgTree() *cobra.Command {
	root := &cobra.Command{Use: "root", TraverseChildren: true}
	server := &cobra.Command{Use: "server"} // grouping command
	create := &cobra.Command{Use: "create", RunE: func(*cobra.Command, []string) error { return nil }}
	server.AddCommand(create)
	root.AddCommand(server)
	return root
}

func TestTrailingHelp_OnRunnableLeaf(t *testing.T) {
	root := helpArgTree()
	if !handleTrailingHelp(root, []string{"server", "create", "help"}) {
		t.Error("expected trailing 'help' after runnable leaf to be handled")
	}
}

func TestTrailingHelp_NotHandledWithoutHelp(t *testing.T) {
	root := helpArgTree()
	if handleTrailingHelp(root, []string{"server", "create"}) {
		t.Error("expected no handling when 'help' is absent")
	}
}

func TestTrailingHelp_NotHandledForExtraArgs(t *testing.T) {
	root := helpArgTree()
	// "help" is not the sole leftover token -> leave it to normal execution.
	if handleTrailingHelp(root, []string{"server", "create", "foo", "help"}) {
		t.Error("expected no handling when 'help' is not the lone leftover arg")
	}
}

func TestTrailingHelp_EmptyArgs(t *testing.T) {
	root := helpArgTree()
	if handleTrailingHelp(root, nil) {
		t.Error("expected no handling for empty args")
	}
}
