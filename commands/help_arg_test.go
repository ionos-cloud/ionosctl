package commands

import (
	"testing"

	"github.com/spf13/cobra"
)

func helpArgTree() *cobra.Command {
	root := &cobra.Command{Use: "root", TraverseChildren: true}
	// Global value flag + bool flag, mirroring real persistent flags.
	root.PersistentFlags().StringP("output", "o", "text", "output format")
	root.PersistentFlags().BoolP("force", "f", false, "force")

	server := &cobra.Command{Use: "server"} // grouping command
	create := &cobra.Command{Use: "create", RunE: func(*cobra.Command, []string) error { return nil }}
	create.Flags().String("name", "", "name")
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

func TestTrailingHelp_WithLocalValueFlagBeforeHelp(t *testing.T) {
	root := helpArgTree()
	// "server create --name foo help" -> rest ["--name","foo","help"].
	if !handleTrailingHelp(root, []string{"server", "create", "--name", "foo", "help"}) {
		t.Error("expected trailing 'help' after flag value to be handled")
	}
}

func TestTrailingHelp_WithGlobalValueFlagBeforeHelp(t *testing.T) {
	root := helpArgTree()
	// Global (inherited) value flag consumes its value, help is still trailing.
	if !handleTrailingHelp(root, []string{"server", "create", "-o", "json", "help"}) {
		t.Error("expected trailing 'help' after global flag value to be handled")
	}
}

func TestTrailingHelp_WithBoolFlagBeforeHelp(t *testing.T) {
	root := helpArgTree()
	// Bool flag does not consume the next arg, so "help" is a trailing request.
	if !handleTrailingHelp(root, []string{"server", "create", "--force", "help"}) {
		t.Error("expected trailing 'help' after bool flag to be handled")
	}
}

func TestTrailingHelp_WithInlineFlagValueBeforeHelp(t *testing.T) {
	root := helpArgTree()
	// "--name=foo" supplies its value inline; "help" is a trailing request.
	if !handleTrailingHelp(root, []string{"server", "create", "--name=foo", "help"}) {
		t.Error("expected trailing 'help' after inline flag value to be handled")
	}
}

func TestTrailingHelp_HelpAfterValueFlagHandled(t *testing.T) {
	root := helpArgTree()
	// A trailing "help" always wins, even right after a value flag: to pass the
	// literal "help" as a value, don't put it last (use "--name=help").
	if !handleTrailingHelp(root, []string{"server", "create", "--name", "help"}) {
		t.Error("expected trailing 'help' after a value flag to be handled")
	}
}

func TestTrailingHelp_HelpAfterShorthandValueHandled(t *testing.T) {
	root := helpArgTree()
	if !handleTrailingHelp(root, []string{"server", "create", "-o", "help"}) {
		t.Error("expected trailing 'help' after a shorthand value to be handled")
	}
}

func TestTrailingHelp_NotHandledWithoutHelp(t *testing.T) {
	root := helpArgTree()
	if handleTrailingHelp(root, []string{"server", "create"}) {
		t.Error("expected no handling when 'help' is absent")
	}
}

func TestTrailingHelp_NotHandledWhenHelpNotLast(t *testing.T) {
	root := helpArgTree()
	// "help" is not the trailing token -> leave it to normal execution.
	if handleTrailingHelp(root, []string{"server", "create", "help", "foo"}) {
		t.Error("expected no handling when 'help' is not the trailing arg")
	}
}

func TestTrailingHelp_NotHandledForNonRunnableParent(t *testing.T) {
	root := helpArgTree()
	// "server help": server is a non-runnable grouping command; leave it to the
	// did-you-mean / help wiring rather than intercepting here.
	if handleTrailingHelp(root, []string{"server", "help"}) {
		t.Error("expected no handling for a non-runnable parent command")
	}
}

func TestTrailingHelp_EmptyArgs(t *testing.T) {
	root := helpArgTree()
	if handleTrailingHelp(root, nil) {
		t.Error("expected no handling for empty args")
	}
}
