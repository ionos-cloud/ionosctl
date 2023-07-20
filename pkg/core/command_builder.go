package core

import (
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/functional"
	"github.com/spf13/viper"
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

// PreRunWithDeprecatedFlags is a decorator for using a command with deprecated flags
// The value of the first flag in the Tuple is set as the value of the second flag of the Tuple
func PreRunWithDeprecatedFlags(f PreCommandRun, flags ...functional.Tuple[string]) PreCommandRun {
	return func(c *PreCommandConfig) error {
		for _, f := range flags {
			if fn := GetFlagName(c.NS, f.First); viper.IsSet(fn) {
				val := viper.Get(fn)
				viper.Set(GetFlagName(c.NS, f.Second), val)
			}
		}
		return f(c)
	}
}
