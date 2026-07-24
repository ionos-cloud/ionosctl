package commands

import (
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
)

// flagSuggestMaxDistance is the maximum Levenshtein distance for a flag name to
// be offered as a "did you mean" suggestion.
const flagSuggestMaxDistance = 2

// suggestingFlagErrorFunc wraps flag-parse errors with a "Did you mean this?"
// hint when the offending flag is close to a known flag. pflag/Cobra do not
// suggest flags out of the box (unlike subcommands), so `--datacentr` used to
// produce only `unknown flag: --datacentr` followed by the whole flag list.
//
// Registered once on the root command; Cobra's FlagErrorFunc is inherited by
// every subcommand, and `cmd` here is the command that failed to parse.
func suggestingFlagErrorFunc(cmd *cobra.Command, err error) error {
	typed := parseUnknownLongFlag(err.Error())
	if typed == "" {
		return err
	}

	suggestions := suggestFlags(cmd, typed)
	if len(suggestions) == 0 {
		return err
	}

	var sb strings.Builder
	sb.WriteString(err.Error())
	sb.WriteString("\n\nDid you mean this?\n")
	for _, s := range suggestions {
		fmt.Fprintf(&sb, "\t--%s\n", s)
	}
	sb.WriteString("\nRun with --help to see all available flags.")
	return fmt.Errorf("%s", sb.String())
}

// parseUnknownLongFlag extracts the flag name from pflag's
// "unknown flag: --name" message. Returns "" for anything else (e.g. unknown
// shorthand flags, which are too short to suggest usefully).
func parseUnknownLongFlag(msg string) string {
	const prefix = "unknown flag: --"
	if !strings.HasPrefix(msg, prefix) {
		return ""
	}
	return strings.TrimSpace(strings.TrimPrefix(msg, prefix))
}

// suggestFlags returns up to three known flag names (local + inherited) whose
// distance to typed is within flagSuggestMaxDistance, closest first.
func suggestFlags(cmd *cobra.Command, typed string) []string {
	type cand struct {
		name string
		dist int
	}
	var cands []cand
	seen := map[string]bool{}

	consider := func(f *flag.Flag) {
		if f.Hidden || seen[f.Name] {
			return
		}
		seen[f.Name] = true
		if d := flagDistance(typed, f.Name); d <= flagSuggestMaxDistance {
			cands = append(cands, cand{f.Name, d})
		}
	}
	cmd.Flags().VisitAll(consider)
	cmd.InheritedFlags().VisitAll(consider)

	sort.Slice(cands, func(i, j int) bool {
		if cands[i].dist != cands[j].dist {
			return cands[i].dist < cands[j].dist
		}
		return cands[i].name < cands[j].name
	})

	out := make([]string, 0, 3)
	for _, c := range cands {
		out = append(out, c.name)
		if len(out) == 3 {
			break
		}
	}
	return out
}

// flagDistance is the Levenshtein distance between typed and a flag name. Many
// ionosctl flags end in "-id"/"-ids", which users routinely omit, so the name
// is also compared with that suffix stripped and the smaller distance wins.
func flagDistance(typed, name string) int {
	d := levenshtein(typed, name)
	for _, suf := range []string{"-ids", "-id"} {
		if base, ok := strings.CutSuffix(name, suf); ok {
			if dd := levenshtein(typed, base); dd < d {
				d = dd
			}
			break
		}
	}
	return d
}

// levenshtein computes the edit distance between two strings.
func levenshtein(a, b string) int {
	ra, rb := []rune(a), []rune(b)
	prev := make([]int, len(rb)+1)
	curr := make([]int, len(rb)+1)
	for j := range prev {
		prev[j] = j
	}
	for i := 1; i <= len(ra); i++ {
		curr[0] = i
		for j := 1; j <= len(rb); j++ {
			cost := 1
			if ra[i-1] == rb[j-1] {
				cost = 0
			}
			curr[j] = min(prev[j]+1, curr[j-1]+1, prev[j-1]+cost)
		}
		prev, curr = curr, prev
	}
	return prev[len(rb)]
}
