package cfg

import (
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	configgen "github.com/ionos-cloud/ionosctl/v6/pkg/cfggen"
	"github.com/spf13/cobra"
)

func GenCfgCmd() *core.Command {
	var (
		printExample bool

		filterVersion    string
		filterWhitelist  []string
		filterBlacklist  []string
		filterVisibility string
		filterGate       string

		mapCustomNames map[string]string
	)

	cmd := &cobra.Command{
		Use:   "login",
		Short: "Use credentials to generate a config file in `ionosctl cfg location`, or use '--example' to generate a sample endpoints YAML config",
		Long: `Generate a YAML file aggregating all product endpoint information at 'ionosctl cfg location'
using the public OpenAPI index.

If using '--example', this command prints the config to stdout.

You can filter by version or specific API names.

There are three ways you can authenticate with the IONOS Cloud APIs:
  1. Interactive mode: Just type 'ionosctl login' and you'll be prompted to enter your username and password.
  2. Use the '--user' and '--password' flags: Enter your credentials in the command.
  3. Use the '--token' flag: Provide an authentication token.
Notes:
  - If using '--token', you can skip verifying the used token with '--skip-verify'
  - If using '--example', the authentication step is skipped
`,
		Example: `
# Print an example YAML configuration file to stdout
ionosctl config login --example

# Login interactively, and generate a YAML config file with filters, to 'ionosctl config location'
ionosctl endpoints generate --version=v1 \
  --whitelist=vpn,psql --blacklist=billing
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			token := "<token>"
			if !printExample {
				token = login(cmd)
			}
			_ = token

			// build filter options
			opts := configgen.Filters{
				CustomNames: mapCustomNames,
			}

			settings := configgen.ProfileSettings{}

			// apply version filter if provided
			if filterVersion != "" {
				opts.Version = &filterVersion
			}

			// always apply hidden filters (defaults set above)
			opts.Visibility = &filterVisibility
			opts.Gate = &filterGate

			// apply whitelist only if flag passed
			if len(filterWhitelist) > 0 {
				opts.Whitelist = make(map[string]bool)
				for _, name := range filterWhitelist {
					opts.Whitelist[name] = true
				}
			}
			// apply blacklist only if flag passed
			if len(filterBlacklist) > 0 {
				opts.Blacklist = make(map[string]bool)
				for _, name := range filterBlacklist {
					opts.Blacklist[name] = true
				}
			}

			// generate config
			out, err := configgen.GenerateConfig(settings, opts)
			if err != nil {
				return fmt.Errorf("could not generate config: %w", err)
			}

			// _, err = cmd.OutOrStdout().Write(out)
			// return err
			if printExample {
				_, err = cmd.OutOrStdout().Write(out)
				if err != nil {
					return fmt.Errorf("could not write config to stdout: %w", err)
				}

				return nil // stop here
			}

			// write config to file
			err = cfg.WriteYAML()
			if err != nil {
				return fmt.Errorf("could not write config to file: %w", err)
			}

			fmt.Fprintf(c.Command.Command.OutOrStdout(), "Config file generated at %s\n", configgen.Location())
			return nil
		},
	}

	f := cmd.Flags()
	f.StringVar(&version, "version", "", "Filter by spec version (e.g. v1)")
	f.StringSliceVar(&whitelist, "whitelist", nil, "Comma-separated list of API names to include")
	f.StringSliceVar(&blacklist, "blacklist", nil, "Comma-separated list of API names to exclude")

	f.BoolVar(&printExample, "example", false, "Print an example YAML config file to stdout and skip authentication step")

	// hidden flags with defaults
	f.StringVar(&visibility, "visibility", "public", "(hidden) Filter by index visibility")
	f.StringVar(&gate, "gate", "General-Availability", "(hidden) Filter by release gate")
	_ = f.MarkHidden("visibility")
	_ = f.MarkHidden("gate")

	return &core.Command{Command: cmd}
}

func login() string {
	return ""
}

var (
	FlagCustomNames   = "custom-names"
	FlagFilterVersion = "filter-version"
	FlagWhitelist     = "whitelist"
	FlagBlacklist     = "blacklist"
	FlagExample       = "example"
	FlagVisibility    = "filter-visibility"
	FlagGate          = "filter-gate"

	FlagSettingsVersion = "version"
	FlagSettingsProfile = "profile-name"
	FlagSettingsEnv     = "environment"
)

func addLoginFlags(cmd *core.Command) {
	cmd.AddBoolFlag(FlagExample, "", false, "Print an example YAML config file to stdout and skip authentication step")

	cmd.AddStringFlag(constants.ArgUser, "", "", "Username to authenticate with. Will be used to generate a token")
	cmd.AddStringFlag(constants.ArgPassword, constants.ArgPasswordShort, "", "Password to authenticate with. Will be used to generate a token")
	cmd.AddStringFlag(constants.ArgToken, constants.ArgTokenShort, "", "Token to authenticate with. If used, will be saved to the config file without generating a new token. Note: mutually exclusive with --user and --password")
	cmd.AddBoolFlag(constants.FlagSkipVerify, "", false, "Forcefully write the provided token to the config file without verifying if it is valid. Note: --token is required")
}

func addProfileFlags(cmd *core.Command) {
	cmd.AddStringFlag(FlagSettingsProfile, "", "user", "Name of the profile to use")
	cmd.AddStringFlag(FlagSettingsEnv, "", "prod", "Environment to use")
	cmd.AddStringFlag(FlagSettingsVersion, "", "1.0", "Version of the config file to use")
}

func addFilterFlags(cmd *core.Command) {
	cmd.AddStringToStringFlag(FlagCustomNames, "", map[string]string{
		"apigateway":                "apigateway",
		"authentication":            "auth",
		"certificatemanager":        "cert",
		"cloud":                     "compute",
		"object‑storage":            "objectstorage",
		"object‑storage‑management": "objectstoragemanagement",
		"mongodb":                   "mongo",
		"postgresql":                "psql",
		"mariadb":                   "mariadb",
		//
		// These are currently the same as the spec name
		// but we can override them here if needed
		// "cdn":                       "cdn",
		// "containerregistry":         "containerregistry",
		// "dataplatform":              "dataplatform",
		// "dns":                       "dns",
		// "kafka":                     "kafka",
		// "logging":                   "logging",
		// "monitoring":                "monitoring",
		// "nfs":                       "nfs",
		// "vmautoscaling":             "vmautoscaling",
		// "vpn":                       "vpn",
	},
		"Define custom names for each spec")
	cmd.AddStringFlag(FlagFilterVersion, "", "", "Filter by spec version (e.g. v1)")
	cmd.AddStringSliceFlag(FlagWhitelist, "", []string{}, "Comma-separated list of API names to include")
	cmd.AddStringSliceFlag(FlagBlacklist, "", []string{}, "Comma-separated list of API names to exclude")
	cmd.AddStringFlag(FlagVisibility, "", "public", "(hidden) Filter by index visibility")
	cmd.AddStringFlag(FlagGate, "", "General-Availability", "(hidden) Filter by release gate")

	_ = cmd.Command.Flags().MarkHidden(FlagVisibility)
	_ = cmd.Command.Flags().MarkHidden(FlagGate)
}
