package core

import (
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/pkg/functional"
)

// CommandBuilder contains information about
// the new Command to be used in Cobra Command
type CommandBuilder struct {
	// Command is a Wrapper around Cobra Command
	Command *Command

	// Command Levels
	// Namespace is the first level of the Command. e.g. [ionosctl] server
	Namespace string
	// Resource is the second level of the Command. e.g. [ionosctl server] volume
	Resource string
	// Verb is the 3rd level of the Command. e.g. [ionosctl server volume] attach
	Verb string

	// Short Description
	ShortDesc string
	// Long Description
	LongDesc string
	// Aliases
	Aliases []string
	// Example of Command run
	Example string

	// Command Run functions
	// to be executed
	PreCmdRun PreCommandRun
	CmdRun    CommandRun
	// Init Client Services
	InitClient bool
}

func (c *CommandBuilder) GetNS() string {
	return fmt.Sprintf("%s.%s.%s", c.Namespace, c.Resource, c.Verb)
}

// PreRunWithDeprecatedFlags is a decorator for using a command with deprecated flags.
// When the first (deprecated) flag in the Tuple is set, its value is aliased onto the
// second (canonical) flag so the command - which reads the canonical flag via
// c.Flags() - sees it. The two flags are expected to be of the same type; the value is
// shared directly, which is type-agnostic and avoids string round-tripping (e.g. a
// StringSlice's "[a,b]" form would not parse back correctly).
func PreRunWithDeprecatedFlags(f PreCommandRun, flags ...functional.Tuple[string]) PreCommandRun {
	return func(c *PreCommandConfig) error {
		fs := c.Command.Command.Flags()
		for _, f := range flags {
			src := fs.Lookup(f.First)
			dst := fs.Lookup(f.Second)
			if src != nil && dst != nil && src.Changed {
				dst.Value = src.Value
				dst.Changed = true
			}
		}
		return f(c)
	}
}
