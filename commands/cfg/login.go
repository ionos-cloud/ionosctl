package cfg

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	configgen "github.com/ionos-cloud/ionosctl/v6/internal/config"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/confirm"
	"github.com/ionos-cloud/ionosctl/v6/pkg/pointer"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/viper"
	"golang.org/x/term"
	"gopkg.in/yaml.v3"
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
  1. Interactive mode: Prompts for username and password, and generates a token that will be saved in the config file.
  2. Use the '--user' and '--password' flags: Used to generate a token that will be saved in the config file.
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
  --version=1.1 --profile-name=my-custom-profile --environment=dev
`,
		PreCmdRun: func(c *core.PreCommandConfig) error {
			c.Command.Command.MarkFlagsMutuallyExclusive(constants.ArgToken, constants.ArgPassword)

			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {
			printExample, err := c.Command.Command.Flags().GetBool(FlagExample)
			if err != nil {
				return fmt.Errorf("could not get flag %s: %w", FlagExample, err)
			}

			configPath := viper.GetString(constants.ArgConfig)

			// if exists, prompt to overwrite with --force override
			if _, err := os.Stat(configPath); !printExample && err == nil {
				yes := confirm.FAsk(os.Stdin, fmt.Sprintf("Config file already exists at %s. Do you want to replace it", configPath), viper.GetBool(constants.ArgForce))
				if !yes {
					return fmt.Errorf(confirm.UserDenied)
				}
			}

			token := "<token>"
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
			version, err := c.Command.Command.Flags().GetFloat64(FlagSettingsVersion)
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
			if filterVisibility != "" {
				opts.Visibility = pointer.From(filterVisibility)
			}
			if filterGate != "" {
				opts.Gate = pointer.From(filterGate)
			}

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

			done := make(chan struct{})
			if !printExample {
				go spinner(c.Command.Command.ErrOrStderr(), done)
			}

			// generate config
			cfg, err := configgen.NewFromIndex(settings, opts)
			if err != nil {
				close(done)
				return fmt.Errorf("could not generate config: %w", err)
			}
			close(done)

			// marshal to YAML
			outBytes, err := yaml.Marshal(cfg)
			if err != nil {
				return fmt.Errorf("could not marshal config to YAML: %w", err)
			}

			if printExample {
				// just print the YAML to stdout
				if _, err := c.Command.Command.OutOrStdout().Write(outBytes); err != nil {
					return fmt.Errorf("could not write config to stdout: %w", err)
				}
				return nil
			}

			if err := os.MkdirAll(filepath.Dir(configPath), 0o700); err != nil {
				return fmt.Errorf("could not create config directory: %w", err)
			}

			// write the file with owner‑only permissions
			if err := os.WriteFile(configPath, outBytes, 0o600); err != nil {
				return fmt.Errorf("could not write config to file %s: %w", configPath, err)
			}

			fmt.Fprintf(c.Command.Command.OutOrStdout(), "Config file generated at %s\n", configPath)
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
		_, _ = fmt.Fprintln(c.Command.Command.OutOrStdout(), "Enter your username: ")
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
		_, _ = fmt.Fprintln(c.Command.Command.OutOrStdout(), "Enter your password: ")
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
	cmd.AddFloat64Flag(FlagSettingsVersion, "", 1.0, "Version of the config file to use")
}

func addFilterFlags(cmd *core.Command) {
	cmd.AddStringToStringFlag(FlagCustomNames, "", map[string]string{
		"authentication":            "auth",
		"certificatemanager":        fileconfiguration.Cert,
		"object-storage":            fileconfiguration.ObjectStorage,
		"object-storage-management": fileconfiguration.ObjectStorageManagement,
		"mongodb":                   fileconfiguration.Mongo,
		"postgresql":                fileconfiguration.PSQL,
		"in-memory-db":              fileconfiguration.InMemoryDB,
		"observability-monitoring":  fileconfiguration.Monitoring,
		"vmautoscaling":             fileconfiguration.Autoscaling,
	},
		"Define custom names for each spec")
	cmd.AddStringFlag(FlagFilterVersion, "", "", "Filter by major spec version (e.g. v1)")
	cmd.AddStringSliceFlag(FlagWhitelist, "", []string{}, "Comma-separated list of API names to include")
	cmd.AddStringSliceFlag(FlagBlacklist, "", []string{"object-storage-user-owned-buckets", "object-storage-contract-owned-buckets"}, "Comma-separated list of API names to exclude")
	cmd.AddStringFlag(FlagVisibility, "", "public", "(hidden) Filter by index visibility")
	cmd.AddStringFlag(FlagGate, "", "", "(hidden) Filter by release gate")

	_ = cmd.Command.Flags().MarkHidden(FlagVisibility)
	_ = cmd.Command.Flags().MarkHidden(FlagGate)
}

// spinner displays a loading spinner until the done channel is closed.
func spinner(out io.Writer, done <-chan struct{}) {
	spinChars := []rune{'|', '/', '-', '\\'}
	i := 0

	// In some cases, the generation takes a short amount of time, in which case don't pollute the output with a spinner right away
	time.Sleep(250 * time.Millisecond)
	// if done already closed, don't start the spinner
	select {
	case <-done:
		return
	default:
		// continue with spinner
	}

	for {
		select {
		case <-done:
			_, _ = fmt.Fprint(out, "\u001B[2K\r")
			return
		default:
			_, _ = fmt.Fprintf(out, "\u001B[2K%c\r", spinChars[i%len(spinChars)])
			time.Sleep(100 * time.Millisecond)
			i++
		}
	}
}
