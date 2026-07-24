package commands

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func newTestTree() *cobra.Command {
	root := &cobra.Command{Use: "root"}
	parent := &cobra.Command{Use: "server"} // grouping command, not runnable
	create := &cobra.Command{Use: "create", RunE: func(*cobra.Command, []string) error { return nil }}
	parent.AddCommand(create)
	root.AddCommand(parent)
	enableUnknownSubcommandSuggestions(root)
	return root
}

func TestUnknownSubcommand_Suggests(t *testing.T) {
	root := newTestTree()
	root.SetArgs([]string{"server", "craete"})
	root.SetOut(&bytes.Buffer{})
	root.SetErr(&bytes.Buffer{})

	err := root.Execute()
	if err == nil {
		t.Fatal("expected error for unknown subcommand, got nil")
	}
	if !strings.Contains(err.Error(), `unknown command "craete"`) {
		t.Errorf("missing unknown-command text: %q", err.Error())
	}
	if !strings.Contains(err.Error(), "Did you mean this?") || !strings.Contains(err.Error(), "create") {
		t.Errorf("missing suggestion for 'create': %q", err.Error())
	}
}

func TestParentWithNoArgs_ShowsHelpNoError(t *testing.T) {
	root := newTestTree()
	root.SetArgs([]string{"server"})
	root.SetOut(&bytes.Buffer{})
	root.SetErr(&bytes.Buffer{})

	if err := root.Execute(); err != nil {
		t.Errorf("bare parent should not error, got: %v", err)
	}
}

func TestUnknownSubcommand_SuggestsHidden(t *testing.T) {
	root := &cobra.Command{Use: "root"}
	// A hidden backward-compat alias, like the root-level "server" alias.
	hidden := &cobra.Command{Use: "server", Hidden: true, RunE: func(*cobra.Command, []string) error { return nil }}
	root.AddCommand(hidden)

	got := suggestSubcommands(root, "serverr")
	if len(got) != 1 || got[0] != "server" {
		t.Errorf("expected hidden 'server' to be suggested, got %v", got)
	}
}

func TestSuggestSubcommands_OnlyClosestBucket(t *testing.T) {
	root := &cobra.Command{Use: "root"}
	for _, name := range []string{"cdn", "dbaas", "dns", "lan", "man", "vpn"} {
		root.AddCommand(&cobra.Command{Use: name, RunE: func(*cobra.Command, []string) error { return nil }})
	}

	// "dnn" is distance 1 from "dns" but distance 2 from cdn/lan/man/vpn.
	// Only the closest (dns) should be suggested.
	got := suggestSubcommands(root, "dnn")
	if len(got) != 1 || got[0] != "dns" {
		t.Errorf("expected only 'dns', got %v", got)
	}
}

func TestRunnableLeaf_NotPatched(t *testing.T) {
	root := newTestTree()
	// The runnable leaf must keep its own RunE, not be overwritten.
	root.SetArgs([]string{"server", "create"})
	root.SetOut(&bytes.Buffer{})
	root.SetErr(&bytes.Buffer{})

	if err := root.Execute(); err != nil {
		t.Errorf("runnable leaf should execute cleanly, got: %v", err)
	}
}
