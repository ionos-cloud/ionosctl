package commands

import (
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	configgen "github.com/ionos-cloud/ionosctl/v6/pkg/cfggen"
	"github.com/spf13/cobra"
)

func GenCfgCmd() *core.Command {
	var (
		version    string
		whitelist  []string
		blacklist  []string
		visibility string
		gate       string

		mapCustomNames map[string]string
	)

	cmd := &cobra.Command{
		Use:   "cfggen",
		Short: "Generate sample endpoints YAML config",
		Long: `Generate a YAML file aggregating all product endpoint information
from the public OpenAPI index. This command prints the config to stdout.

You can filter by version or specific API names.
`,
		Example: `
# Generate all v1 public GA endpoints
ionosctl endpoints generate --version=v1

# Include only vpn and psql APIs, exclude billing
ionosctl endpoints generate --version=v1 \
  --whitelist=vpn,psql --blacklist=billing
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// build filter options
			opts := configgen.FilterOptions{
				CustomNames: mapCustomNames,
			}

			// apply version filter if provided
			if version != "" {
				opts.Version = &version
			}

			// always apply hidden filters (defaults set above)
			opts.Visibility = &visibility
			opts.Gate = &gate

			// apply whitelist only if flag passed
			if len(whitelist) > 0 {
				opts.Whitelist = make(map[string]bool)
				for _, name := range whitelist {
					opts.Whitelist[name] = true
				}
			}
			// apply blacklist only if flag passed
			if len(blacklist) > 0 {
				opts.Blacklist = make(map[string]bool)
				for _, name := range blacklist {
					opts.Blacklist[name] = true
				}
			}

			// generate config
			out, err := configgen.GenerateConfig(opts)
			if err != nil {
				return fmt.Errorf("could not generate config: %w", err)
			}

			// print to stdout
			_, err = cmd.OutOrStdout().Write(out)
			return err
		},
	}

	// public flags
	f := cmd.Flags()

	// override default spec names with our product names on sdk-go-bundle
	f.StringToStringVar(&mapCustomNames, "custom-names",
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

	f.StringVar(&version, "version", "", "Filter by spec version (e.g. v1)")
	f.StringSliceVar(&whitelist, "whitelist", nil, "Comma-separated list of API names to include")
	f.StringSliceVar(&blacklist, "blacklist", nil, "Comma-separated list of API names to exclude")

	// hidden flags with defaults
	f.StringVar(&visibility, "visibility", "public", "(hidden) Filter by index visibility")
	f.StringVar(&gate, "gate", "General-Availability", "(hidden) Filter by release gate")

	_ = f.MarkHidden("visibility")
	_ = f.MarkHidden("gate")

	cmd.SilenceUsage = true
	cmd.Flags().SortFlags = false

	return &core.Command{Command: cmd}
}
