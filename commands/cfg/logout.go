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
	"github.com/ionos-cloud/sdk-go-bundle/shared"
	"gopkg.in/yaml.v3"
)

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

	cmd.AddBoolFlag("only-purge-old", "", false,
		"Skip YAML logout and only purge legacy config.json")

	return cmd
}

// promptAndDelete handles the common prompt + delete pattern.
func promptAndDelete(path, desc string, c *core.CommandConfig) {
	fmt.Fprintf(c.Command.Command.OutOrStdout(),
		"⚠️  Detected %s at '%s'\n", desc, path)
	if confirm.FAsk(
		c.Command.Command.InOrStdin(),
		fmt.Sprintf("⚠️  Delete %s '%s'", desc, path),
		viper.GetBool(constants.ArgForce),
	) {
		if err := os.Remove(path); err != nil {
			fmt.Fprintf(c.Command.Command.ErrOrStderr(),
				"Warning: failed to delete '%s': %v\n", path, err)
		} else {
			fmt.Fprintf(c.Command.Command.OutOrStdout(),
				"Deleted %s: '%s'\n", desc, path)
		}
	}
}

// handleJSONConfig checks if --config is a JSON file and, if so, prompts & deletes it.
func handleJSONConfig(c *core.CommandConfig) error {
	cfgPath := viper.GetString(constants.ArgConfig)
	if filepath.Ext(cfgPath) != ".json" {
		return nil
	}
	promptAndDelete(cfgPath, "JSON config", c)
	return nil
}

// maybeDeleteOldConfig looks for a legacy config.json next to your YAML config,
// inspects it for userdata fields, and prompts & deletes if present.
func maybeDeleteOldConfig(c *core.CommandConfig, cl *client.Client) {
	// primary: same dir as YAML
	dir := filepath.Dir(cl.ConfigPath)
	jsonCfg := filepath.Join(dir, "config.json")

	// fallback: directory of --config
	if _, err := os.Stat(jsonCfg); os.IsNotExist(err) {
		jsonCfg = filepath.Join(filepath.Dir(viper.GetString(constants.ArgConfig)), "config.json")
	}

	if _, err := os.Stat(jsonCfg); os.IsNotExist(err) {
		return
	}

	raw, err := os.ReadFile(jsonCfg)
	if err != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(),
			"Warning: found legacy %s but failed to read: %v\n", jsonCfg, err)
		return
	}

	var old map[string]interface{}
	if err := json.Unmarshal(raw, &old); err != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(),
			"Warning: found legacy %s but failed to parse: %v\n", jsonCfg, err)
		return
	}

	// detect any userdata.* keys
	want := []string{"userdata.token", "userdata.name", "userdata.password", "userdata.api-url"}
	var has []string
	for _, k := range want {
		if _, ok := old[k]; ok {
			has = append(has, k)
		}
	}
	if len(has) == 0 {
		return
	}

	promptAndDelete(jsonCfg, fmt.Sprintf("legacy config.json containing %v", has), c)
}
