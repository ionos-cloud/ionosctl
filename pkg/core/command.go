package core

import (
	"fmt"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Command struct {
	// NS is the Global Namespace of the Command
	NS      string
	Command *cobra.Command

	subCommands []*Command
}

func (c *Command) AddCommand(commands ...*Command) {
	c.subCommands = append(c.subCommands, commands...)
	for _, cmd := range commands {
		if cmd != nil && c.Command != nil {
			c.Command.AddCommand(cmd.Command)
		}
	}
}

func (c *Command) Name() string {
	if c != nil && c.Command != nil {
		return c.Command.Name()
	} else {
		return ""
	}
}

func (c *Command) CommandPath() string {
	return c.Command.CommandPath()
}

func (c *Command) GlobalFlags() *flag.FlagSet {
	return c.Command.PersistentFlags()
}

func (c *Command) SubCommands() []*Command {
	return c.subCommands
}

func (c *Command) IsAvailableCommand() bool {
	if c != nil && c.Command != nil {
		return c.Command.IsAvailableCommand()
	} else {
		return false
	}
}

func (c *Command) AddStringFlag(name, shorthand, defaultValue, desc string, optionFunc ...FlagOptionFunc) {
	flags := c.Command.Flags()
	if shorthand != "" {
		flags.StringP(name, shorthand, defaultValue, desc)
	} else {
		flags.String(name, defaultValue, desc)
	}
	viper.BindPFlag(GetFlagName(c.NS, name), c.Command.Flags().Lookup(name))

	// Add Option to Flag
	for _, option := range optionFunc {
		option(c, name)
	}
}

func (c *Command) AddStringToStringFlag(name, shorthand string, defaultValue map[string]string, desc string, optionFunc ...FlagOptionFunc) {
	flags := c.Command.Flags()
	if shorthand != "" {
		flags.StringToStringP(name, shorthand, defaultValue, desc)
	} else {
		flags.StringToString(name, defaultValue, desc)
	}
	viper.BindPFlag(GetFlagName(c.NS, name), c.Command.Flags().Lookup(name))

	// Add Option to Flag
	for _, option := range optionFunc {
		option(c, name)
	}
}

func (c *Command) SetFlagAnnotation(name, key string, values ...string) {
	flags := c.Command.Flags()
	flags.SetAnnotation(name, key, values)
}

func (c *Command) GetAnnotations() map[string]string {
	return c.Command.Annotations
}

func (c *Command) GetAnnotationsByKey(key string) string {
	if c != nil && c.Command != nil {
		return c.Command.Annotations[key]
	} else {
		return ""
	}
}

func (c *Command) AddStringSliceFlag(name, shorthand string, defaultValue []string, desc string, optionFunc ...FlagOptionFunc) {
	flags := c.Command.Flags()
	if shorthand != "" {
		flags.StringSliceP(name, shorthand, defaultValue, desc)
	} else {
		flags.StringSlice(name, defaultValue, desc)
	}
	viper.BindPFlag(GetFlagName(c.NS, name), c.Command.Flags().Lookup(name))

	// Add Option to Flag
	for _, option := range optionFunc {
		option(c, name)
	}
}

func (c *Command) AddIntSliceFlag(name, shorthand string, defaultValue []int, desc string, optionFunc ...FlagOptionFunc) {
	flags := c.Command.Flags()
	if shorthand != "" {
		flags.IntSliceP(name, shorthand, defaultValue, desc)
	} else {
		flags.IntSlice(name, defaultValue, desc)
	}
	viper.BindPFlag(GetFlagName(c.NS, name), c.Command.Flags().Lookup(name))

	// Add Option to Flag
	for _, option := range optionFunc {
		option(c, name)
	}
}

func (c *Command) AddIntFlag(name, shorthand string, defaultValue int, desc string, optionFunc ...FlagOptionFunc) {
	flags := c.Command.Flags()
	if shorthand != "" {
		flags.IntP(name, shorthand, defaultValue, desc)
	} else {
		flags.Int(name, defaultValue, desc)
	}
	viper.BindPFlag(GetFlagName(c.NS, name), c.Command.Flags().Lookup(name))

	// Add Option to Flag
	for _, option := range optionFunc {
		option(c, name)
	}
}

func (c *Command) AddInt32Flag(name, shorthand string, defaultValue int32, desc string, optionFunc ...FlagOptionFunc) {
	flags := c.Command.Flags()
	if shorthand != "" {
		flags.Int32P(name, shorthand, defaultValue, desc)
	} else {
		flags.Int32(name, defaultValue, desc)
	}
	viper.BindPFlag(GetFlagName(c.NS, name), c.Command.Flags().Lookup(name))

	// Add Option to Flag
	for _, option := range optionFunc {
		option(c, name)
	}
}

func (c *Command) AddFloat32Flag(name, shorthand string, defaultValue float32, desc string, optionFunc ...FlagOptionFunc) {
	flags := c.Command.Flags()
	if shorthand != "" {
		flags.Float32P(name, shorthand, defaultValue, desc)
	} else {
		flags.Float32(name, defaultValue, desc)
	}
	viper.BindPFlag(GetFlagName(c.NS, name), c.Command.Flags().Lookup(name))

	// Add Option to Flag
	for _, option := range optionFunc {
		option(c, name)
	}
}

func (c *Command) AddBoolFlag(name, shorthand string, defaultValue bool, desc string, optionFunc ...FlagOptionFunc) {
	flags := c.Command.Flags()
	if shorthand != "" {
		flags.BoolP(name, shorthand, defaultValue, desc)
	} else {
		flags.Bool(name, defaultValue, desc)
	}
	viper.BindPFlag(GetFlagName(c.NS, name), c.Command.Flags().Lookup(name))

	// Add Option to Flag
	for _, option := range optionFunc {
		option(c, name)
	}
}

func GetFlagName(ns, flagName string) string {
	return fmt.Sprintf("%s.%s", ns, flagName)
}

// GetGlobalFlagName returns a string of cmdName on which the
// Flag is defined as a Global Flag concatenated with the name
// of the Flag.
//
// For example: in a `ionosctl namespace resource verb` command
// structure, if the flag is inherited from the namespace level
// at the verb level, the cmdName will be c.Namespace,
// If the Global Flag is defined at resource level, the cmdName
// will be c.Resource.
func GetGlobalFlagName(cmdName, flagName string) string {
	return fmt.Sprintf("%s.%s", cmdName, flagName)
}
