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
//   - The first location in 'allowedLocations' is used as the default for non-list commands.
//   - List commands using [CommandConfig.ListAllLocations] query all locations when '--location' is not set.
func WithRegionalConfigOverride(c *Command, productNames []string, templateFallbackURL string, allowedLocations []string) *Command {
	if len(productNames) == 0 {
		panic(fmt.Errorf("no productNames provided for %s", c.Command.Name()))
	}
	if len(allowedLocations) == 0 {
		panic(fmt.Errorf("no allowedLocations provided for %s", c.Command.Name()))
	}

	// Store regional metadata for child commands (e.g., ListAllLocations)
	if c.Command.Annotations == nil {
		c.Command.Annotations = map[string]string{}
	}
	c.Command.Annotations[AnnotationLocations] = strings.Join(allowedLocations, ",")
	c.Command.Annotations[AnnotationTemplateURL] = templateFallbackURL
	c.Command.Annotations[AnnotationProductNames] = strings.Join(productNames, ",")

	// Add the server URL flag
	c.Command.PersistentFlags().StringP(
		constants.ArgServerUrl, constants.ArgServerUrlShort, templateFallbackURL,
		fmt.Sprintf("Override default host URL. If contains placeholder, location will be embedded. "+
			"Preferred over the config file override '%s' and env var '%s'", productNames[0], constants.EnvServerUrl),
	)

	// Add the location flag. The default is left empty on multi-location APIs so
	// cobra does not print a misleading "(default "<first-loc>")": when unset, list
	// commands query all locations and single-resource commands require --location.
	// Single-location APIs keep the sole location as an explicit default, which is
	// accurate there. The empty default is resolved to allowedLocations[0] in
	// PersistentPreRunE below (via findOverriddenURL returning "").
	locationDefault := ""
	if len(allowedLocations) == 1 {
		locationDefault = allowedLocations[0]
	}
	c.Command.PersistentFlags().StringP(
		constants.FlagLocation, constants.FlagLocationShort, locationDefault,
		"Location of the resource to operate on. When unset, list commands query all locations. "+
			"Can be one of: "+strings.Join(allowedLocations, ", "),
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

		// Single-resource commands operate on one location. When --location is
		// unset, resolve the default (first allowed location) so the config-file
		// override for that region is honored, not just the bare template URL.
		// List commands ignore this viper value; they resolve a URL per location
		// via ListAllLocations instead.
		lookupLocation := location
		if lookupLocation == "" {
			lookupLocation = allowedLocations[0]
		}

		url := findOverriddenURL(cmd, productNames, templateFallbackURL, lookupLocation)
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
		url := findOverriddenURL(cmd, productNames, fallbackURL, "")
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

// locationVariants returns the candidate spellings of a location to look up in
// the config file. IONOS regions are written with slashes in some contexts
// (e.g. "de/fra", "eu/central/3") and with dashes in others (e.g.
// "eu-central-3"). The command flags and the config file do not always use the
// same convention, so try the location as given plus its slash/dash variants.
func locationVariants(location string) []string {
	if location == "" {
		return []string{""}
	}
	variants := []string{location}
	if slash := strings.ReplaceAll(location, "-", "/"); slash != location {
		variants = append(variants, slash)
	}
	if dash := strings.ReplaceAll(location, "/", "-"); dash != location {
		variants = append(variants, dash)
	}
	return variants
}

func findOverriddenURL(cmd *cobra.Command, productNames []string, fallbackURL, location string) string {
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
	cl, _ := client.Get()
	if cl != nil && cl.Config != nil {
		for _, prod := range productNames {
			for _, loc := range locationVariants(location) {
				if override := cl.Config.GetOverride(prod, loc); override != nil {
					return override.Name
				}
			}
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
