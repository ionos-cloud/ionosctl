package cfg

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	configgen "github.com/ionos-cloud/ionosctl/v6/pkg/cfggen"
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

	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "config",
		Resource:  "config",
		Verb:      "login",
		ShortDesc: "Use credentials to generate a config file in `ionosctl cfg location`, or use '--example' to generate a sample endpoints YAML config",
		LongDesc: `Generate a YAML file aggregating all product endpoint information at 'ionosctl cfg location'
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
			if !printExample {
				token = getToken()
			}
			_ = token

			// build filter options
			opts := configgen.Filters{
				CustomNames: mapCustomNames,
			}

			name, _ := cmd.Flags().GetString("profile")
			env, _ := cmd.Flags().GetString("environment")
			version, _ := cmd.Flags().GetString("version")
			settings := configgen.ProfileSettings{
				Token:       token,
				ProfileName: name,
				Environment: env,
				Version:     version,
			}

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
			}

			// else,write to config file
			return nil
		},
	})

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	return cmd
}

func getToken() string {
	return ""
}

func addFilterFlags() {
	// override default spec names with our product names on sdk-go-bundle
	f.StringToString("custom-names",
		map[string]string{
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

	f.String("filter-version", "", "Filter by spec version (e.g. v1)")
	f.StringSlice("whitelist", nil, "Comma-separated list of API names to include")
	f.StringSlice("blacklist", nil, "Comma-separated list of API names to exclude")

	f.Bool("example", false, "Print an example YAML config file to stdout and skip authentication step")

	// hidden flags with defaults
	f.String("visibility", "public", "(hidden) Filter by index visibility")
	f.String("gate", "General-Availability", "(hidden) Filter by release gate")

	_ = f.MarkHidden("visibility")
	_ = f.MarkHidden("gate")
}
