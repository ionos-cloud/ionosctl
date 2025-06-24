package cfg

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/confirm"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

// LogoutCmd returns the `logout` command, now with an extra post-check for config.json
func LogoutCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Verb:      "logout",
		ShortDesc: "Remove credentials from your YAML config (and purge old JSON)",
		LongDesc: fmt.Sprintf(`This 'Quality of Life' command will:

  1. Clear out any sensitive fields in your YAML config.
  2. Afterwards, detect and optionally delete any legacy config.json alongside it.

You can skip the YAML logout and **only** purge the old JSON with:

    ionosctl logout --only-purge-old

%s`, constants.DescAuthenticationOrder),
		Example:   "ionosctl logout\nionosctl logout --only-purge-old",
		PreCmdRun: core.NoPreRun,

		CmdRun: func(c *core.CommandConfig) error {
			cl, _ := client.Get()

			// ensure there's a YAML config to act on
			if cl == nil || cl.Config == nil {
				return fmt.Errorf("no config file found, nothing to logout from")
			}

			// perform the actual logout: clear credentials
			cfg := cl.Config
			for i := range cfg.Profiles {
				cfg.Profiles[i].Credentials = shared.Credentials{}
			}

			// ensure directory and write updated YAML
			if err := os.MkdirAll(
				filepath.Dir(cl.ConfigPath),
				0o700,
			); err != nil {
				return fmt.Errorf("could not create config directory: %w", err)
			}
			outBytes, err := yaml.Marshal(cl.Config)
			if err != nil {
				return fmt.Errorf("could not marshal config to YAML: %w", err)
			}
			if err := os.WriteFile(
				cl.ConfigPath,
				outBytes,
				0o600,
			); err != nil {
				return fmt.Errorf(
					"could not write config to file %s: %w",
					cl.ConfigPath,
					err,
				)
			}

			fmt.Fprintf(
				c.Command.Command.OutOrStdout(),
				"Removed credentials from %s but kept URL overrides\n",
				cl.ConfigPath,
			)

			// clean up any legacy JSON config if exits
			maybeDeleteOldConfig(c, cl)
			return nil
		},
		InitClient: false,
	})

	cmd.AddBoolFlag("only-purge-old", "", false,
		"Skip YAML logout and only purge legacy config.json")

	return cmd
}

// maybeDeleteOldConfig looks for a legacy config.json alongside the new YAML config,
// warns the user if it contains any userdata.* keys, and (only if they consent or --force)
// deletes it.  If they decline, we simply leave it in place.
func maybeDeleteOldConfig(c *core.CommandConfig, cl *client.Client) {
	dir := filepath.Dir(cl.ConfigPath)
	jsonCfgPath := filepath.Join(dir, "config.json")

	// bail early if there's no old JSON file at all
	if _, err := os.Stat(jsonCfgPath); err != nil {
		return
	}

	// read and detect any userdata.* keys
	raw, err := os.ReadFile(jsonCfgPath)
	if err != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(),
			"Warning: found legacy config.json but failed to read it: %v\n", err)
		return
	}
	var oldCfg map[string]interface{}
	if err := json.Unmarshal(raw, &oldCfg); err != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(),
			"Warning: found legacy config.json but failed to parse it: %v\n", err)
		return
	}

	keys := []string{
		"userdata.token",
		"userdata.name",
		"userdata.password",
		"userdata.api-url",
	}
	var found []string
	for _, k := range keys {
		if _, ok := oldCfg[k]; ok {
			found = append(found, k)
		}
	}
	if len(found) == 0 {
		return // nothing meaningful to delete
	}

	// warn & prompt
	fmt.Fprintf(
		c.Command.Command.OutOrStdout(),
		"⚠️  Detected legacy config.json at %s containing fields: %v\n",
		jsonCfgPath, found,
	)
	if confirm.FAsk(
		c.Command.Command.InOrStdin(),
		fmt.Sprintf("⚠️  Delete legacy %s", jsonCfgPath),
		viper.GetBool(constants.ArgForce),
	) {
		if err := os.Remove(jsonCfgPath); err != nil {
			fmt.Fprintf(c.Command.Command.ErrOrStderr(),
				"Warning: failed to delete %s: %v\n", jsonCfgPath, err)
		} else {
			fmt.Fprintf(
				c.Command.Command.OutOrStdout(),
				"Deleted legacy config.json: %s\n",
				jsonCfgPath,
			)
		}
	}
}
