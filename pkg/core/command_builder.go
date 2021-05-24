package core

import (
	"fmt"
)

// CommandBuilder contains information about
// the new Command to be used in Cobra Command
type CommandBuilder struct {
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
