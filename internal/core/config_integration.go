package core

import (
	"fmt"
	"os"
	"strings"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// WithRegionalConfigOverride adds regional flag support to a command, allowing users to specify a location or override the server URL.
// To use this function, wrap the root command of your API and specify the baseURL and allowed locations.
//
// Example:
//
//		func DNSCommand() *core.Command {
//			cmd := &core.Command{
//				Command: &cobra.Command{
//					Use:              "dns",
//					Short:            "The sub-commands of the 'dns' resource help automate DNS Zone and Record management",
//					TraverseChildren: true,
//				},
//			}
//
//			// Add regional flags
//			return core.WithRegionalConfigOverride(cmd, []string{fileconfiguration.Cloud, "compute"}, "https://dns.%s.ionos.com", []string{"de/fra", "de/txl"})
//		}
//
//	  - 'baseURL': The base URL for the API, with an optional '%s' placeholder for the location (e.g., '"https://dns.%s.ionos.com"').
//	  - 'allowedLocations': A slice of allowed locations (e.g., '[]string{"de/fra", "de/txl"}'). These will populate the '--location' flag completion.
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
func WithRegionalConfigOverride(c *Command, productNames []string, templateFallbackURL string, allowedLocations []string) *Command {
	if len(productNames) == 0 {
		panic(fmt.Errorf("no productNames provided for %s", c.Command.Name()))
	}
	if len(allowedLocations) == 0 {
		panic(fmt.Errorf("no allowedLocations provided for %s", c.Command.Name()))
	}

	// Add the server URL flag
	c.Command.PersistentFlags().StringP(
		constants.ArgServerUrl, constants.ArgServerUrlShort, templateFallbackURL,
		fmt.Sprintf("Override default host URL. If contains placeholder, location will be embedded. "+
			"Preferred over the config file override '%s' and env var '%s'", productNames[0], constants.EnvServerUrl),
	)

	// Add the location flag
	c.Command.PersistentFlags().StringP(
		constants.FlagLocation, constants.FlagLocationShort, allowedLocations[0],
		"Location of the resource to operate on. Can be one of: "+strings.Join(allowedLocations, ", "),
	)
	viper.BindPFlag(constants.FlagLocation, c.Command.PersistentFlags().Lookup(constants.FlagLocation))
	c.Command.RegisterFlagCompletionFunc(constants.FlagLocation,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return allowedLocations, cobra.ShellCompDirectiveNoFileComp
		},
	)

	originalPreRun := c.Command.PersistentPreRunE
	c.Command.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		c.Command.MarkFlagsMutuallyExclusive(constants.ArgServerUrl, constants.FlagLocation)
		location, _ := cmd.Flags().GetString(constants.FlagLocation)

		url := findOverridenURL(cmd, productNames, templateFallbackURL, location)
		if url == "" {
			url = fmt.Sprintf(templateFallbackURL, strings.ReplaceAll(allowedLocations[0], "/", "-"))
		}
		viper.Set(constants.ArgServerUrl, url)

		if originalPreRun != nil {
			return originalPreRun(cmd, args)
		}
		return nil
	}

	return c
}

func WithConfigOverride(c *Command, productNames []string, fallbackURL string) *Command {
	if len(productNames) == 0 {
		panic(fmt.Errorf("no productNames provided for %s", c.Command.Name()))
	}
	if fallbackURL == "" {
		fallbackURL = constants.DefaultApiURL
	}

	c.Command.PersistentFlags().StringP(
		constants.ArgServerUrl, constants.ArgServerUrlShort, fallbackURL,
		fmt.Sprintf("Override default host URL. Preferred over the config file override '%s' and env var '%s'", productNames[0], constants.EnvServerUrl),
	)

	originalPreRun := c.Command.PersistentPreRunE
	c.Command.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		url := findOverridenURL(cmd, productNames, fallbackURL, "")
		if url == "" {
			url = fallbackURL
		}
		viper.Set(constants.ArgServerUrl, url)

		if originalPreRun != nil {
			return originalPreRun(cmd, args)
		}
		return nil
	}

	return c
}

func findOverridenURL(cmd *cobra.Command, productNames []string, fallbackURL, location string) string {
	// Check if the --server-url flag is set
	if cmd.Flags().Changed(constants.ArgServerUrl) {
		serverURL, _ := cmd.Flags().GetString(constants.ArgServerUrl)
		viper.Set(constants.ArgServerUrl, serverURL)
		return serverURL
	}

	// If IONOS_API_URL is set, use it as the server URL
	if envURL := os.Getenv(constants.EnvServerUrl); envURL != "" {
		viper.Set(constants.EnvServerUrl, envURL)
		return envURL
	}

	// return override from config file if available
	for _, prod := range productNames {
		cl, _ := client.Get()
		if override := cl.Config.GetOverride(prod, location); override != nil {
			return override.Name
		}
	}

	// otherwise, return the fallback URL
	if location != "" {
		if !strings.Contains(fallbackURL, "%s") {
			return fallbackURL
		}
		normalizedLocation := strings.ReplaceAll(location, "/", "-")
		return fmt.Sprintf(fallbackURL, normalizedLocation)
	}

	return ""
}
