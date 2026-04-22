package comptplus

import (
	"strings"

	"github.com/elk-language/go-prompt"
	istrings "github.com/elk-language/go-prompt/strings"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// CobraLexer provides syntax highlighting for cobra command input.
// It colorizes commands, flags, flag values, and unknown tokens by
// walking the cobra command tree on each keystroke.
type CobraLexer struct {
	root   *cobra.Command
	tokens []prompt.Token
	pos    int

	// Colors are configurable. Zero values fall back to defaults.
	CommandColor Color
	FlagColor    Color
	ValueColor   Color
	ErrorColor   Color
}

// Color wraps prompt.Color so consumers don't need to import go-prompt directly.
type Color = prompt.Color

// NewCobraLexer creates a lexer that highlights input based on the given cobra command tree.
func NewCobraLexer(root *cobra.Command) *CobraLexer {
	return &CobraLexer{
		root:         root,
		CommandColor: prompt.Green,
		FlagColor:    prompt.Cyan,
		ValueColor:   prompt.Yellow,
		ErrorColor:   prompt.Red,
	}
}

func (l *CobraLexer) Init(input string) {
	l.tokens = l.tokenize(input)
	l.pos = 0
}

func (l *CobraLexer) Next() (prompt.Token, bool) {
	if l.pos >= len(l.tokens) {
		return nil, false
	}
	t := l.tokens[l.pos]
	l.pos++
	return t, true
}

// span represents a word's position and text in the input.
type span struct {
	start int
	end   int
	text  string
}

// tokenize walks the input, splits into whitespace and word spans,
// then classifies words using cobra command tree context.
func (l *CobraLexer) tokenize(input string) []prompt.Token {
	if len(input) == 0 {
		return nil
	}

	// Split input into spans preserving byte positions
	spans := splitSpans(input)
	if len(spans) == 0 {
		return nil
	}

	// Resolve the command path: walk the cobra tree for as many leading words as match
	cmd := l.root
	commandSpanCount := 0
	for _, s := range spans {
		if s.text == "" { // whitespace span
			continue
		}
		found := false
		for _, sub := range cmd.Commands() {
			if sub.Name() == s.text || containsAlias(sub, s.text) {
				cmd = sub
				commandSpanCount++
				found = true
				break
			}
		}
		if !found {
			break
		}
	}

	// Now classify each span
	// NOTE: go-prompt's renderer uses LastByteIndex as inclusive (it does input[first:last+1]),
	// so we pass end-1 as the last byte index.
	var tokens []prompt.Token
	wordIndex := 0     // counts non-whitespace spans
	prevWord := ""     // the previous non-whitespace word

	for _, s := range spans {
		lastByte := istrings.ByteNumber(s.end - 1)

		// Whitespace — default color
		if s.text == "" {
			tokens = append(tokens, prompt.NewSimpleToken(
				istrings.ByteNumber(s.start),
				lastByte,
			))
			continue
		}

		color := l.classifyWord(s.text, wordIndex, commandSpanCount, prevWord, cmd)
		tokens = append(tokens, prompt.NewSimpleToken(
			istrings.ByteNumber(s.start),
			lastByte,
			prompt.SimpleTokenWithColor(color),
		))

		prevWord = s.text
		wordIndex++
	}

	return tokens
}

// classifyWord determines the color for a word token.
func (l *CobraLexer) classifyWord(word string, wordIndex, commandSpanCount int, prevWord string, cmd *cobra.Command) prompt.Color {
	// Part of the command path
	if wordIndex < commandSpanCount {
		return l.CommandColor
	}

	// Flag (starts with -)
	if strings.HasPrefix(word, "-") {
		flagName := extractFlagName(word)
		if lookupFlag(cmd, flagName) != nil {
			return l.FlagColor
		}
		return l.ErrorColor
	}

	// Check if this is a flag value (previous word was a non-bool flag)
	if prevWord != "" && strings.HasPrefix(prevWord, "-") {
		flagName := extractFlagName(prevWord)
		if f := lookupFlag(cmd, flagName); f != nil && f.Value.Type() != "bool" {
			return l.ValueColor
		}
	}

	// Positional argument — default color
	return prompt.DefaultColor
}

// splitSpans splits input into alternating whitespace and word spans,
// preserving byte offsets. Handles single/double quotes and backslash escapes.
func splitSpans(input string) []span {
	var spans []span
	i := 0

	for i < len(input) {
		// Whitespace span
		if input[i] == ' ' || input[i] == '\t' {
			start := i
			for i < len(input) && (input[i] == ' ' || input[i] == '\t') {
				i++
			}
			spans = append(spans, span{start: start, end: i, text: ""})
			continue
		}

		// Word span (respects quotes and escapes)
		start := i
		inSingle := false
		inDouble := false
		for i < len(input) {
			if input[i] == '\\' && !inSingle && i+1 < len(input) {
				i += 2
				continue
			}
			if input[i] == '\'' && !inDouble {
				inSingle = !inSingle
				i++
				continue
			}
			if input[i] == '"' && !inSingle {
				inDouble = !inDouble
				i++
				continue
			}
			if !inSingle && !inDouble && (input[i] == ' ' || input[i] == '\t') {
				break
			}
			i++
		}
		spans = append(spans, span{start: start, end: i, text: input[start:i]})
	}

	return spans
}

// lookupFlag finds a flag by name (long or shorthand) on the command, including inherited flags.
func lookupFlag(cmd *cobra.Command, name string) *pflag.Flag {
	if len(name) == 0 {
		return nil
	}
	if len(name) == 1 {
		if f := cmd.Flags().ShorthandLookup(name); f != nil {
			return f
		}
		return cmd.InheritedFlags().ShorthandLookup(name)
	}
	if f := cmd.Flags().Lookup(name); f != nil {
		return f
	}
	return cmd.InheritedFlags().Lookup(name)
}

func extractFlagName(word string) string {
	name := strings.TrimLeft(word, "-")
	if eqIdx := strings.IndexByte(name, '='); eqIdx >= 0 {
		name = name[:eqIdx]
	}
	return name
}

func containsAlias(cmd *cobra.Command, name string) bool {
	for _, a := range cmd.Aliases {
		if a == name {
			return true
		}
	}
	return false
}
