package core

import (
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
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

// WithRegionalFlags adds regional flag support to a command, allowing users to specify a location or override the server URL.
// To use this function, wrap the root command of your API and specify the baseURL and allowed locations.
//
// Example:
//
// ```
//
//	func DNSCommand() *core.Command {
//		cmd := &core.Command{
//			Command: &cobra.Command{
//				Use:              "dns",
//				Short:            "The sub-commands of the 'dns' resource help automate DNS Zone and Record management",
//				TraverseChildren: true,
//			},
//		}
//
//		// Add regional flags
//		return core.WithRegionalFlags(cmd, "https://dns.%s.ionos.com", []string{"de/fra", "de/txl"})
//	}
//
// ```
//
//   - 'baseURL': The base URL for the API, with an optional '%s' placeholder for the location (e.g., '"https://dns.%s.ionos.com"').
//   - 'allowedLocations': A slice of allowed locations (e.g., '[]string{"de/fra", "de/txl"}'). These will populate the '--location' flag completion.
//
// # Notes
//
//   - The '--server-url' flag allows users to override the API host URL manually.
//   - The '--location' flag allows users to specify a region, which replaces the '%s' placeholder in the 'baseURL'.
//   - If '--location' is used and is valid (from 'allowedLocations'), the 'baseURL' is formatted with the normalized location.
//   - If '--location' is invalid or unsupported, a warning is logged, but the constructed URL is still attempted.
//   - If 'allowedLocations' is empty, the function panics, as this is considered a programming error.
//   - If an unsupported location is provided, a warning is logged:
//     'WARN: <location> is an invalid location. Valid locations are: <allowedLocations>'
//   - This also marks '--api-url' and '--location' flags as mutually exclusive.
//   - The first location in 'allowedLocations' is used as the default URL if no location is provided.
func WithRegionalFlags(c *Command, baseURL string, allowedLocations []string) *Command {
	// Generate the default URL using the first provided location, if available
	var defaultUrl string
	if len(allowedLocations) > 0 {
		// if baseURL does not contain a placeholder, throw a panic, as this is a programming error
		if !strings.Contains(baseURL, "%s") {
			panic(fmt.Errorf("baseURL %s does not contain a placeholder for location", baseURL))
		}

		defaultLocation := allowedLocations[0]
		defaultUrl = fmt.Sprintf(baseURL, strings.ReplaceAll(defaultLocation, "/", "-"))
	} else {
		// If no locations are provided, panic, this is a programming error
		panic(fmt.Errorf("no allowedLocations provided for %s", c.Command.Name()))
	}

	// Add the server URL flag
	c.Command.PersistentFlags().StringP(
		constants.ArgServerUrl, constants.ArgServerUrlShort, defaultUrl, "Override default host URL",
	)
	viper.BindPFlag(constants.ArgServerUrl, c.Command.PersistentFlags().Lookup(constants.ArgServerUrl))

	// Add the location flag
	c.Command.PersistentFlags().StringP(
		constants.FlagLocation, constants.FlagLocationShort, "", "Location of the resource to operate on. Can be one of: "+strings.Join(allowedLocations, ", "),
	)
	viper.BindPFlag(constants.FlagLocation, c.Command.PersistentFlags().Lookup(constants.FlagLocation))
	c.Command.RegisterFlagCompletionFunc(constants.FlagLocation,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return allowedLocations, cobra.ShellCompDirectiveNoFileComp
		},
	)

	// Wrap the pre-run logic to handle mutually exclusive flags
	originalPreRun := c.Command.PersistentPreRunE
	c.Command.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if originalPreRun != nil {
			if err := originalPreRun(cmd, args); err != nil {
				return err
			}
		}

		// Mark the flags as mutually exclusive
		c.Command.MarkFlagsMutuallyExclusive(constants.ArgServerUrl, constants.FlagLocation)

		// Handle the location flag
		if location, _ := cmd.Flags().GetString(constants.FlagLocation); location != "" {
			normalizedLoc := strings.ReplaceAll(location, "/", "-")
			if strings.Contains(baseURL, "%s") {
				viper.Set(constants.ArgServerUrl, fmt.Sprintf(baseURL, normalizedLoc))
			} else {
				// Log a warning if trying to use location with a non-placeholder URL
				fmt.Fprintf(c.Command.ErrOrStderr(), jsontabwriter.GenerateLogOutput(
					"WARN: Ignoring location %s because this API does not support region-based URLs", location))
			}
		}

		return nil
	}

	return c
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

// AddSetFlag adds a String slice flag with support for limitation to certain values in a slice.
// It also adds completions for those limited values,
// on top of throwing an error if the flag value isn't found among the marked valid values
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
