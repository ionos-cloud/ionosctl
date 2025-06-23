package cfg

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
	"gopkg.in/yaml.v3"
)

func LogoutCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Verb:      "logout",
		ShortDesc: "Convenience command for removing config file credentials",
		LongDesc: fmt.Sprintf(`This command is a 'Quality of Life' command which will parse your config file for fields that contain sensitive data.
If any such fields are found, their values will be replaced with an empty string.

%s`, constants.DescAuthenticationOrder),
		Example:   "ionosctl logout",
		PreCmdRun: core.NoPreRun,
		CmdRun: func(c *core.CommandConfig) error {
			if client.Must().Config == nil {
				return fmt.Errorf("no config file found, nothing to logout from")
			}

			cfg := client.Must().Config

			// for each profile, remove the sensitive fields
			for i := range cfg.Profiles {
				cfg.Profiles[i].Credentials = shared.Credentials{}
			}

			if err := os.MkdirAll(filepath.Dir(client.Must().ConfigPath), 0o700); err != nil {
				return fmt.Errorf("could not create config directory: %w", err)
			}

			// marshal to YAML
			outBytes, err := yaml.Marshal(client.Must().Config)
			if err != nil {
				return fmt.Errorf("could not marshal config to YAML: %w", err)
			}

			// write the file with owner‑only permissions
			if err := os.WriteFile(client.Must().ConfigPath, outBytes, 0o600); err != nil {
				return fmt.Errorf("could not write config to file %s: %w", client.Must().ConfigPath, err)
			}

			fmt.Fprintf(c.Command.Command.OutOrStdout(), "Removed credentials from %s\n", client.Must().ConfigPath)

			return nil
		},
		InitClient: false,
	})

	return cmd
}
