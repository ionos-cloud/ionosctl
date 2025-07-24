package client

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	cfg "github.com/ionos-cloud/ionosctl/v6/internal/config"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

// ConfigSource holds a loaded FileConfig (or nil) plus the path it came from.
type ConfigSource struct {
	Config *fileconfiguration.FileConfig
	Path   string
}

// retrieveConfigFile returns the first successful config or a fatal error.
// If no config is found, ConfigSource.Config will be nil and Path will be the
// default flag value (for backwards‐compat).
func retrieveConfigFile() (ConfigSource, error) {
	loaders := []func() (ConfigSource, error){
		loadFromFlag,
		loadFromEnvVar,
		loadFromJSONMigration,
		loadFromSDKDefault,
	}

	for _, load := range loaders {
		src, err := load()
		if err != nil {
			// I/O or parse error: stop immediately
			return ConfigSource{}, err
		}
		if src.Config != nil {
			// found a valid config; return it
			return src, nil
		}
		// not found -> try next
	}

	fmt.Println("none found: printing default")
	// none found -> return nil config, default flag path
	return ConfigSource{nil, viper.GetString(constants.ArgConfig)}, nil
}

// loadFromFlag tries --config; if empty or missing -> (nil, ""), nil;
// on I/O/parse error -> err; on success -> (*FileConfig, path), nil.
func loadFromFlag() (ConfigSource, error) {
	path := viper.GetString(constants.ArgConfig)
	if path == "" {
		return ConfigSource{}, nil
	}
	return tryLoad(path, "--config")
}

// loadFromEnvVar tries IONOS_CONFIG_FILE; same semantics as loadFromFlag.
func loadFromEnvVar() (ConfigSource, error) {
	path := os.Getenv(shared.IonosFilePathEnvVar)
	if path == "" {
		return ConfigSource{}, nil
	}
	return tryLoad(path, fmt.Sprintf("env %s", shared.IonosFilePathEnvVar))
}

// loadFromJSONMigration migrates legacy config.json (if present next to --config path).
// On success returns the new YAML config; on any error returns err.
func loadFromJSONMigration() (ConfigSource, error) {
	yamlPath := viper.GetString(constants.ArgConfig)
	if yamlPath == "" {
		return ConfigSource{}, nil
	}
	jsonPath := filepath.Join(filepath.Dir(yamlPath), "config.json")
	if _, err := os.Stat(jsonPath); err != nil {
		return ConfigSource{}, nil
	}

	migrated, err := cfg.MigrateFromJSON(jsonPath)
	if err != nil {
		return ConfigSource{}, fmt.Errorf("failed migrating %q to YAML: %w", jsonPath, err)
	}
	if migrated == nil {
		return ConfigSource{}, nil
	}

	out, _ := yaml.Marshal(migrated)
	if err := os.WriteFile(yamlPath, out, 0o600); err != nil {
		fmt.Fprintf(os.Stderr,
			"Warning: could not write migrated config to %s: %v\n",
			yamlPath, err)
	}
	return ConfigSource{migrated, yamlPath}, nil
}

// loadFromSDKDefault tries the SDK default path (~/.ionos/config).
func loadFromSDKDefault() (ConfigSource, error) {
	defaultPath, err := fileconfiguration.DefaultConfigFileName()
	if err != nil {
		return ConfigSource{}, fmt.Errorf("could not determine default config path: %w", err)
	}
	return tryLoad(defaultPath, "SDK default")
}

// tryLoad loads a FileConfig from the given path.
// If the file does not exist, it returns a nil FileConfig and the path.
// If the file exists but cannot be loaded, it returns an error.
// If the file is loaded successfully, it returns the FileConfig and the path.
func tryLoad(path, sourceDesc string) (ConfigSource, error) {
	cfg, err := fileconfiguration.New(path)
	if err != nil {
		if strings.Contains(err.Error(), "does not exist") {
			return ConfigSource{nil, path}, nil
		}
		return ConfigSource{}, fmt.Errorf("%s: failed loading %q: %w", sourceDesc, path, err)
	}
	return ConfigSource{cfg, path}, nil
}
