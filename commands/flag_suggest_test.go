package commands

import (
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestParseUnknownLongFlag(t *testing.T) {
	cases := map[string]string{
		"unknown flag: --datacentr":   "datacentr",
		"unknown flag: --quiet":       "quiet",
		"unknown shorthand flag: 'x'": "", // shorthands not suggested
		"some other error":            "",
	}
	for msg, want := range cases {
		if got := parseUnknownLongFlag(msg); got != want {
			t.Errorf("parseUnknownLongFlag(%q) = %q, want %q", msg, got, want)
		}
	}
}

func TestFlagDistance_StripsIDSuffix(t *testing.T) {
	// "datacentr" is 4 edits from "datacenter-id" but only 1 from "datacenter".
	if d := flagDistance("datacentr", "datacenter-id"); d != 1 {
		t.Errorf("flagDistance = %d, want 1 (suffix should be stripped)", d)
	}
}

func TestSuggestFlags(t *testing.T) {
	cmd := &cobra.Command{Use: "list"}
	cmd.Flags().String("datacenter-id", "", "")
	cmd.Flags().Bool("quiet", false, "")
	cmd.Flags().String("name", "", "")

	got := suggestFlags(cmd, "datacentr")
	if len(got) == 0 || got[0] != "datacenter-id" {
		t.Errorf("suggestFlags(datacentr) = %v, want datacenter-id first", got)
	}

	if got := suggestFlags(cmd, "zzzzzzzz"); len(got) != 0 {
		t.Errorf("suggestFlags(zzzzzzzz) = %v, want none", got)
	}
}

func TestSuggestingFlagErrorFunc(t *testing.T) {
	cmd := &cobra.Command{Use: "list"}
	cmd.Flags().String("datacenter-id", "", "")

	err := suggestingFlagErrorFunc(cmd, errUnknownFlag("datacentr"))
	if !strings.Contains(err.Error(), "Did you mean this?") ||
		!strings.Contains(err.Error(), "--datacenter-id") {
		t.Errorf("expected suggestion in error, got: %q", err.Error())
	}

	// No close match -> original error is returned unchanged.
	orig := errUnknownFlag("zzzzzzzz")
	if got := suggestingFlagErrorFunc(cmd, orig); got.Error() != orig.Error() {
		t.Errorf("expected unchanged error, got: %q", got.Error())
	}
}

type unknownFlagErr string

func (e unknownFlagErr) Error() string { return "unknown flag: --" + string(e) }

func errUnknownFlag(name string) error { return unknownFlagErr(name) }
