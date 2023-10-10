package shellpath

import (
	"strings"
)

// EscapePathForShell escapes special characters in paths to make them safe for Unix-like processes.
func EscapePathForShell(path string) string {
	// Note: Order matters, as characters like backslash must be escaped first.
	replacements := []struct {
		original string
		escaped  string
	}{
		{`\`, `\\`},
		{` `, `\ `},
		{`'`, `'\''`},
		{`"`, `\"`},
		{`$`, `\$`},
		{`!`, `\!`},
		{`&`, `\&`},
		{`;`, `\;`},
		{`|`, `\|`},
		{`(`, `\(`},
		{`)`, `\)`},
		{`{`, `\{`},
		{`}`, `\}`},
		{`*`, `\*`},
		{`?`, `\?`},
		{`[`, `\[`},
		{`]`, `\]`},
	}

	for _, r := range replacements {
		path = strings.ReplaceAll(path, r.original, r.escaped)
	}

	return path
}
