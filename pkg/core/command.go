package core

import (
	"fmt"
	"net"
	"strings"
	"time"

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

func (c *Command) AddIpSliceFlag(name, shorthand string, defaultValue []net.IP, desc string, optionFunc ...FlagOptionFunc) {
	flags := c.Command.Flags()
	if shorthand != "" {
		flags.IPSliceP(name, shorthand, defaultValue, desc)
	} else {
		flags.IPSlice(name, defaultValue, desc)
	}
	_ = viper.BindPFlag(GetFlagName(c.NS, name), c.Command.Flags().Lookup(name))

	// Add Option to Flag
	for _, option := range optionFunc {
		option(c, name)
	}
}

func (c *Command) AddIpFlag(name, shorthand string, defaultValue net.IP, desc string, optionFunc ...FlagOptionFunc) {
	flags := c.Command.Flags()
	if shorthand != "" {
		flags.IPP(name, shorthand, defaultValue, desc)
	} else {
		flags.IP(name, defaultValue, desc)
	}
	_ = viper.BindPFlag(GetFlagName(c.NS, name), c.Command.Flags().Lookup(name))

	// Add Option to Flag
	for _, option := range optionFunc {
		option(c, name)
	}
}

func (c *Command) AddUUIDFlag(name, shorthand, defaultValue, desc string, optionFunc ...FlagOptionFunc) {
	flags := c.Command.Flags()
	if shorthand != "" {
		flags.VarP(newUuidFlag(defaultValue), name, shorthand, desc)
	} else {
		flags.Var(newUuidFlag(defaultValue), name, desc)
	}
	viper.BindPFlag(GetFlagName(c.NS, name), c.Command.Flags().Lookup(name))

	// Add Option to Flag
	for _, option := range optionFunc {
		option(c, name)
	}
}

func (c *Command) AddSetFlag(name, shorthand, defaultValue string, allowed []string, desc string, optionFunc ...FlagOptionFunc) {
	flags := c.Command.Flags()
	desc += fmt.Sprintf(". Can be one of: %s", strings.Join(allowed, ", "))
	if shorthand != "" {
		flags.VarP(newSetFlag(defaultValue, allowed), name, shorthand, desc)
	} else {
		flags.Var(newSetFlag(defaultValue, allowed), name, desc)
	}
	viper.BindPFlag(GetFlagName(c.NS, name), c.Command.Flags().Lookup(name))

	_ = c.Command.RegisterFlagCompletionFunc(name, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allowed, cobra.ShellCompDirectiveNoFileComp
	})

	// Add Option to Flag
	for _, option := range optionFunc {
		option(c, name)
	}
}

func (c *Command) AddStringVarFlag(address *string, name, shorthand, value, desc string, optionFunc ...FlagOptionFunc) {
	flags := c.Command.Flags()
	if shorthand != "" {
		flags.StringVarP(address, name, shorthand, value, desc)
	} else {
		flags.StringVar(address, name, value, desc)
	}
	viper.BindPFlag(GetFlagName(c.NS, name), c.Command.Flags().Lookup(name))

	// Add Option to Flag
	for _, option := range optionFunc {
		option(c, name)
	}
}

func (c *Command) AddDurationFlag(name, shorthand string, defaultValue time.Duration, desc string, optionFunc ...FlagOptionFunc) {
	flags := c.Command.Flags()
	if shorthand != "" {
		flags.DurationP(name, shorthand, defaultValue, desc)
	} else {
		flags.Duration(name, defaultValue, desc)
	}
	viper.BindPFlag(GetFlagName(c.NS, name), c.Command.Flags().Lookup(name))

	// Add Option to Flag
	for _, option := range optionFunc {
		option(c, name)
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

func (c *Command) AddStringToStringVarFlag(v *map[string]string, name, shorthand string, defaultValue map[string]string, desc string, optionFunc ...FlagOptionFunc) {
	flags := c.Command.Flags()
	if shorthand != "" {
		flags.StringToStringVarP(v, name, shorthand, defaultValue, desc)
	} else {
		flags.StringToStringVar(v, name, defaultValue, desc)
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

func (c *Command) AddStringSliceVarFlag(address *[]string, name, shorthand string, defaultValue []string, desc string, optionFunc ...FlagOptionFunc) {
	flags := c.Command.Flags()
	if shorthand != "" {
		flags.StringSliceVarP(address, name, shorthand, defaultValue, desc)
	} else {
		flags.StringSliceVar(address, name, defaultValue, desc)
	}
	viper.BindPFlag(GetFlagName(c.NS, name), c.Command.Flags().Lookup(name))

	// Add Option to Flag
	for _, option := range optionFunc {
		option(c, name)
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

func (c *Command) AddInt32VarFlag(address *int32, name, shorthand string, defaultValue int32, desc string, optionFunc ...FlagOptionFunc) {
	flags := c.Command.Flags()
	if shorthand != "" {
		flags.Int32VarP(address, name, shorthand, defaultValue, desc)
	} else {
		flags.Int32Var(address, name, defaultValue, desc)
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
