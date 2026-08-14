package commands

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/spf13/cobra"
)

// TestRootSilencesUsage guards the wiring: the root command must have
// SilenceUsage set so errors don't dump the whole flag list.
func TestRootSilencesUsage(t *testing.T) {
	if !GetRootCmd().Command.SilenceUsage {
		t.Error("root command SilenceUsage should be true")
	}
}

// TestSilenceUsageIsInherited confirms Cobra suppresses the usage dump for a
// subcommand error when only the root has SilenceUsage set.
func TestSilenceUsageIsInherited(t *testing.T) {
	root := &cobra.Command{Use: "root", SilenceUsage: true, SilenceErrors: true}
	sub := &cobra.Command{
		Use:   "boom",
		RunE:  func(*cobra.Command, []string) error { return errors.New("kaboom") },
		Short: "always fails",
	}
	root.AddCommand(sub)

	var out bytes.Buffer
	root.SetOut(&out)
	root.SetErr(&out)
	root.SetArgs([]string{"boom"})

	if err := root.Execute(); err == nil {
		t.Fatal("expected error")
	}
	if strings.Contains(out.String(), "Usage:") {
		t.Errorf("usage should be suppressed, got:\n%s", out.String())
	}
}

// TestRootTraverseChildren guards the flag that the unknown-subcommand
// suggestion wiring (internal/ux) depends on: without TraverseChildren, Cobra
// resolves subcommands differently and the did-you-mean path never triggers.
func TestRootTraverseChildren(t *testing.T) {
	if !GetRootCmd().Command.TraverseChildren {
		t.Error("root command TraverseChildren should be true")
	}
}

// TestRootVersionSet confirms the version string is wired into Cobra so
// "--version" and the version template render something.
func TestRootVersionSet(t *testing.T) {
	if GetRootCmd().Command.Version == "" {
		t.Error("root command Version should be set")
	}
}

// TestRootHelpCommand confirms the custom help command is wired in. Cobra
// stores a SetHelpCommand'd command out of the normal subcommand list, so this
// checks the package-level var the init wiring installs.
func TestRootHelpCommand(t *testing.T) {
	if helpCommand == nil || !strings.HasPrefix(helpCommand.Use, "help") {
		t.Errorf("custom help command not wired, got %+v", helpCommand)
	}
}

// TestRootPersistentFlags guards the global flags every command relies on.
func TestRootPersistentFlags(t *testing.T) {
	pf := GetRootCmd().Command.PersistentFlags()
	for _, name := range []string{
		constants.ArgOutput,
		constants.ArgForce,
		constants.ArgQuiet,
		constants.ArgVerbose,
		constants.ArgConfig,
	} {
		if pf.Lookup(name) == nil {
			t.Errorf("root persistent flag %q missing", name)
		}
	}
}
