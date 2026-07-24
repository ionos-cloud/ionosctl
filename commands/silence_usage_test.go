package commands

import (
	"bytes"
	"errors"
	"strings"
	"testing"

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
