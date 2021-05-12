package core

import (
	"fmt"

	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	multierror "go.uber.org/multierr"
)

type Command struct {
	// NS is the Global Namespace of the Command.
	NS      string
	Command *cobra.Command

	subCommands []*Command
}

func (c *Command) AddCommand(commands ...*Command) {
	c.subCommands = append(c.subCommands, commands...)
	for _, cmd := range commands {
		c.Command.AddCommand(cmd.Command)
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

func (c *Command) MarkFlagRequired(flagName string) error {
	return c.Command.MarkFlagRequired(flagName)
}

func (c *Command) AddStringFlag(name, shorthand, defaultValue, desc string) {
	flags := c.Command.Flags()
	if shorthand != "" {
		flags.StringP(name, shorthand, defaultValue, desc)
	} else {
		flags.String(name, defaultValue, desc)
	}
	if c.Command.Parent() != nil {
		viper.BindPFlag(GetFlagName(c.NS, name), c.Command.Flags().Lookup(name))
	} else {
		viper.BindPFlag(GetFlagName(c.NS, name), c.Command.Flags().Lookup(name))
	}
}

func (c *Command) AddStringSliceFlag(name, shorthand string, defaultValue []string, desc string) {
	flags := c.Command.Flags()
	if shorthand != "" {
		flags.StringSliceP(name, shorthand, defaultValue, desc)
	} else {
		flags.StringSlice(name, defaultValue, desc)
	}
	if c.Command.Parent() != nil {
		viper.BindPFlag(GetFlagName(c.NS, name), c.Command.Flags().Lookup(name))
	} else {
		viper.BindPFlag(GetFlagName(c.NS, name), c.Command.Flags().Lookup(name))
	}
}

func (c *Command) AddIntFlag(name, shorthand string, defaultValue int, desc string) {
	flags := c.Command.Flags()
	if shorthand != "" {
		flags.IntP(name, shorthand, defaultValue, desc)
	} else {
		flags.Int(name, defaultValue, desc)
	}
	if c.Command.Parent() != nil {
		viper.BindPFlag(GetFlagName(c.NS, name), c.Command.Flags().Lookup(name))
	} else {
		viper.BindPFlag(GetFlagName(c.NS, name), c.Command.Flags().Lookup(name))
	}
}

func (c *Command) AddFloat32Flag(name, shorthand string, defaultValue float32, desc string) {
	flags := c.Command.Flags()
	if shorthand != "" {
		flags.Float32P(name, shorthand, defaultValue, desc)
	} else {
		flags.Float32(name, defaultValue, desc)
	}
	if c.Command.Parent() != nil {
		viper.BindPFlag(GetFlagName(c.NS, name), c.Command.Flags().Lookup(name))
	} else {
		viper.BindPFlag(GetFlagName(c.NS, name), c.Command.Flags().Lookup(name))
	}
}

func (c *Command) AddBoolFlag(name, shorthand string, defaultValue bool, desc string) {
	flags := c.Command.Flags()
	if shorthand != "" {
		flags.BoolP(name, shorthand, defaultValue, desc)
	} else {
		flags.Bool(name, defaultValue, desc)
	}
	if c.Command.Parent() != nil {
		viper.BindPFlag(GetFlagName(c.NS, name), c.Command.Flags().Lookup(name))
	} else {
		viper.BindPFlag(GetFlagName(c.NS, name), c.Command.Flags().Lookup(name))
	}
}

func CheckRequiredGlobalFlags(cmdNs string, globalFlagsName ...string) error {
	var multiErr error
	for _, flagName := range globalFlagsName {
		if viper.GetString(GetGlobalFlagName(cmdNs, flagName)) == "" {
			multiErr = multierror.Append(multiErr, clierror.NewRequiredFlagErr(flagName))
		}
	}
	if multiErr != nil {
		return multiErr
	}
	return nil
}

func CheckRequiredFlags(ns string, localFlagsName ...string) error {
	var multiErr error
	for _, flagName := range localFlagsName {
		if viper.GetString(GetFlagName(ns, flagName)) == "" {
			multiErr = multierror.Append(multiErr, clierror.NewRequiredFlagErr(flagName))
		}
	}
	if multiErr != nil {
		return multiErr
	}
	return nil
}

func GetFlagName(ns, flagName string) string {
	return fmt.Sprintf("%s.%s", ns, flagName)
}

func GetGlobalFlagName(cmdNs, flagName string) string {
	return fmt.Sprintf("%s.%s", cmdNs, flagName)
}
