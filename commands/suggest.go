package commands

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/spf13/cobra"
)

// enableUnknownSubcommandSuggestions walks the command tree and makes every
// non-runnable parent command (i.e. a command that only groups subcommands)
// report an actionable error when it receives an unknown subcommand.
//
// Without this, a typo like "ionosctl server craete" silently prints the parent
// help and exits 0, because the root command uses TraverseChildren: Cobra's
// Traverse() returns the parent with the unknown token as a leftover arg and no
// error, so a non-runnable parent just falls through to Help(). By attaching a
// RunE we force the error path and surface Cobra's built-in "Did you mean this?"
// suggestions.
func enableUnknownSubcommandSuggestions(cmd *cobra.Command) {
	for _, sub := range cmd.Commands() {
		enableUnknownSubcommandSuggestions(sub)
	}

	// Only patch pure grouping commands. Anything already runnable, or with no
	// children, is left untouched.
	if cmd.Runnable() || !cmd.HasSubCommands() {
		return
	}

	// Mark it so doc generation still treats it as a non-runnable parent even
	// though we attach a RunE below (otherwise every group gets its own page).
	if cmd.Annotations == nil {
		cmd.Annotations = map[string]string{}
	}
	cmd.Annotations[core.GroupCommandAnnotation] = "true"

	cmd.RunE = func(c *cobra.Command, args []string) error {
		if len(args) == 0 {
			// No subcommand given: preserve the previous behaviour of printing help.
			return c.Help()
		}

		// Silence the usage dump; the message below is self-contained.
		c.SilenceUsage = true
		return errors.New(unknownSubcommandError(c, args[0]))
	}
}

// unknownSubcommandError builds an "unknown command" message mirroring Cobra's
// own wording, including "Did you mean this?" suggestions when any are close.
func unknownSubcommandError(c *cobra.Command, arg string) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "unknown command %q for %q", arg, c.CommandPath())
	if suggestions := suggestSubcommands(c, arg); len(suggestions) > 0 {
		sb.WriteString("\n\nDid you mean this?\n")
		for _, s := range suggestions {
			fmt.Fprintf(&sb, "\t%v\n", s)
		}
	}
	fmt.Fprintf(&sb, "\nRun '%v --help' for usage.", c.CommandPath())
	return sb.String()
}

// suggestSubcommandsMaxDistance is the maximum Levenshtein distance for a
// subcommand name/alias to be offered as a suggestion.
const suggestSubcommandsMaxDistance = 2

// suggestSubcommands returns the names of child commands close to the typed
// token. Unlike Cobra's SuggestionsFor it also considers hidden commands, so
// backward-compat aliases (e.g. the root-level "server" alias of
// "compute server") are still suggested on a typo.
//
// Only the closest matches are returned: for a short typo like "dnn" a plain
// distance-2 threshold matches nearly everything (cdn, lan, man, ...), so we
// keep just the candidates at the smallest distance found ("dns" at distance
// 1). A prefix match counts as distance 0 (strongest signal).
func suggestSubcommands(c *cobra.Command, typed string) []string {
	lowerTyped := strings.ToLower(typed)
	best := map[string]int{} // command name -> smallest match distance
	consider := func(cmdName, matchName string) {
		d := levenshtein(lowerTyped, strings.ToLower(matchName))
		if strings.HasPrefix(strings.ToLower(matchName), lowerTyped) {
			d = 0
		}
		if d > suggestSubcommandsMaxDistance {
			return
		}
		if prev, ok := best[cmdName]; !ok || d < prev {
			best[cmdName] = d
		}
	}

	for _, cmd := range c.Commands() {
		// Match on the command name and any of its aliases, including hidden
		// commands (which Cobra's own suggester skips).
		for _, name := range append([]string{cmd.Name()}, cmd.Aliases...) {
			consider(cmd.Name(), name)
		}
		// Honour explicit SuggestFor hints as well.
		for _, sf := range cmd.SuggestFor {
			if strings.EqualFold(typed, sf) {
				best[cmd.Name()] = 0
			}
		}
	}

	if len(best) == 0 {
		return nil
	}

	// Keep only the closest bucket, so distant coincidental matches drop out.
	minDist := suggestSubcommandsMaxDistance + 1
	for _, d := range best {
		if d < minDist {
			minDist = d
		}
	}
	var suggestions []string
	for name, d := range best {
		if d == minDist {
			suggestions = append(suggestions, name)
		}
	}
	sort.Strings(suggestions)
	return suggestions
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
