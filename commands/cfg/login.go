package cfg

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	configgen "github.com/ionos-cloud/ionosctl/v6/pkg/cfggen"
	"github.com/ionos-cloud/ionosctl/v6/pkg/pointer"
	"golang.org/x/term"
)

func Login() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "config",
		Resource:  "config",
		Verb:      "login",
		ShortDesc: "Use credentials to generate a config file in `ionosctl cfg location`, or use '--example' to generate a sample endpoints YAML config",
		LongDesc: `Generate a YAML file aggregating all product endpoint information at 'ionosctl cfg location' using the public OpenAPI index.

If using '--example', this command prints the config to stdout without any authentication step.

You can filter by version (--filter-version), whitelist (--whitelist) or blacklist (--blacklist) specific APIs,
and customize the names of the APIs in the config file using --custom-names.

There are three ways you can authenticate with the IONOS Cloud APIs:
  1. Interactive mode: Just type 'ionosctl login' and you'll be prompted to enter your username and password.
  2. Use the '--user' and '--password' flags: Enter your credentials in the command.
  3. Use the '--token' flag: Provide an authentication token.
Notes:
  - If using '--example', the authentication step is skipped
`,
		Example: `
# Print an example YAML configuration file to stdout
ionosctl config login --example

# Login interactively, and generate a YAML config file with filters, to 'ionosctl config location'
ionosctl endpoints generate --filter-version=v1 \
  --whitelist=vpn,psql --blacklist=billing

# Specify a token, a config version, a custom profile name, and a custom environment
ionosctl config login --token $IONOS_TOKEN \
  --version=v1 --profile=my-custom-profile --environment=dev
`,
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {
			token := "<token>"
			fmt.Println("Getting from ", core.GetFlagName(c.NS, FlagExample))

			printExample, err := c.Command.Command.Flags().GetBool(FlagExample)
			if err != nil {
				return fmt.Errorf("could not get flag %s: %w", FlagExample, err)
			}

			if !printExample {
				var err error
				token, err = getToken(c)
				if err != nil {
					return fmt.Errorf("could not retrieve token: %w", err)
				}
			}

			// build filter options
			customNames, err := c.Command.Command.Flags().GetStringToString(FlagCustomNames)
			if err != nil {
				return fmt.Errorf("could not get flag %s: %w", FlagCustomNames, err)
			}
			opts := configgen.Filters{
				CustomNames: customNames,
			}

			profileName, err := c.Command.Command.Flags().GetString(FlagSettingsProfile)
			if err != nil {
				return fmt.Errorf("could not get flag %s: %w", FlagSettingsProfile, err)
			}
			env, err := c.Command.Command.Flags().GetString(FlagSettingsEnv)
			if err != nil {
				return fmt.Errorf("could not get flag %s: %w", FlagSettingsEnv, err)
			}
			version, err := c.Command.Command.Flags().GetString(FlagSettingsVersion)
			if err != nil {
				return fmt.Errorf("could not get flag %s: %w", FlagSettingsVersion, err)
			}
			settings := configgen.ProfileSettings{
				Token:       token,
				ProfileName: profileName,
				Environment: env,
				Version:     version,
			}

			// apply version filter if provided
			filterVersion, err := c.Command.Command.Flags().GetString(FlagFilterVersion)
			if err != nil && filterVersion != "" {
				opts.Version = pointer.From(filterVersion)
			}

			// always apply hidden filters (defaults set above)
			filterVisibility, _ := c.Command.Command.Flags().GetString(FlagVisibility)
			filterGate, _ := c.Command.Command.Flags().GetString(FlagGate)
			opts.Visibility = pointer.From(filterVisibility)
			opts.Gate = pointer.From(filterGate)

			// apply whitelist only if flag passed
			filterWhitelist, _ := c.Command.Command.Flags().GetStringSlice(FlagWhitelist)
			if len(filterWhitelist) > 0 {
				opts.Whitelist = make(map[string]bool)
				for _, name := range filterWhitelist {
					opts.Whitelist[name] = true
				}
			}
			// apply blacklist only if flag passed
			filterBlacklist, _ := c.Command.Command.Flags().GetStringSlice(FlagBlacklist)
			if len(filterBlacklist) > 0 {
				opts.Blacklist = make(map[string]bool)
				for _, name := range filterBlacklist {
					opts.Blacklist[name] = true
				}
			}

			// generate config
			cfg, err := configgen.GenerateConfig(settings, opts)
			if err != nil {
				return fmt.Errorf("could not generate config: %w", err)
			}

			if printExample {
				bytes, err := cfg.ToBytesYAML()
				if err != nil {
					return fmt.Errorf("could not convert config to bytes: %w", err)
				}
				_, err = c.Command.Command.OutOrStdout().Write(bytes)
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
	})

	addProfileFlags(cmd)
	addLoginFlags(cmd)
	addFilterFlags(cmd)

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	return cmd
}

func getToken(c *core.CommandConfig) (string, error) {
	token, err := c.Command.Command.Flags().GetString(constants.ArgToken)
	if err != nil {
		return "", fmt.Errorf("could not get flag %s: %w", constants.ArgToken, err)
	}
	if token != "" {
		return token, nil
	}

	username, _ := c.Command.Command.Flags().GetString(constants.ArgUser)
	if username == "" {
		fmt.Fprintln(c.Command.Command.OutOrStdout(), "Enter your username: ")
		reader := bufio.NewReader(c.Command.Command.InOrStdin())
		var err error
		username, err = reader.ReadString('\n')
		if err != nil {
			return "", fmt.Errorf("failed reading username from set reader")
		}
		username = strings.TrimSpace(username) // remove trailing newline
	}

	password, _ := c.Command.Command.Flags().GetString(constants.ArgPassword)
	if password == "" {
		fmt.Fprintln(c.Command.Command.OutOrStdout(), "Enter your password: ")
		if file, ok := c.Command.Command.InOrStdin().(*os.File); ok {
			bytePassword, err := term.ReadPassword(int(file.Fd()))
			if err != nil {
				return "", fmt.Errorf("failed securely reading password from set file descriptor")
			}
			password = string(bytePassword)
		} else {
			return "", fmt.Errorf("the set input does not have a file descriptor (is it set to a terminal?)")
		}
	}

	apiToken, _, err := client.NewClient(username, password, "", "").
		AuthClient.TokensApi.TokensGenerate(context.Background()).Execute()
	if err != nil {
		return "", fmt.Errorf("failed using username and password to generate a token: %w", err)
	}

	return *apiToken.Token, nil
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

	// cant use viper here, because it would also look at USER env var value
	cmd.Command.Flags().StringP(constants.ArgUser, "", "", "Username to authenticate with. Will be used to generate a token")
	cmd.AddStringFlag(constants.ArgPassword, constants.ArgPasswordShort, "", "Password to authenticate with. Will be used to generate a token")
	cmd.AddStringFlag(constants.ArgToken, constants.ArgTokenShort, "", "Token to authenticate with. If used, will be saved directly to the config file. Note: mutually exclusive with --user and --password")
	cmd.AddBoolFlag(constants.FlagSkipVerify, "", false, "Forcefully write the provided token to the config file without verifying if it is valid. Note: --token is required")
}

func addProfileFlags(cmd *core.Command) {
	cmd.AddStringFlag(FlagSettingsProfile, "", "user", "Name of the profile to use")
	cmd.AddStringFlag(FlagSettingsEnv, "", "prod", "Environment to use")
	cmd.AddStringFlag(FlagSettingsVersion, "", "1.0", "Version of the config file to use")
}

func addFilterFlags(cmd *core.Command) {
	cmd.AddStringToStringFlag(FlagCustomNames, "", map[string]string{
		"authentication":            "auth",
		"certificatemanager":        "cert",
		"cloud":                     "compute",
		"object‑storage":            "objectstorage",
		"object‑storage‑management": "objectstoragemanagement",
		"mongodb":                   "mongo",
		"postgresql":                "psql",

		//
		// These are currently the same as the spec name
		// but we can override them here if needed
		//
		// "cdn":                       "cdn",
		// "apigateway":                "apigateway",
		// "mariadb":                   "mariadb",
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
