package config

import (
	"os"
	"path/filepath"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/die"

	"github.com/spf13/viper"
)

// GetConfigFilePath sanitizes the --config flag input and returns the path to the config file.
// If none set, it returns the default config path.
func GetConfigFilePath() string {
	path := filepath.Join(getConfigHomeDir(), constants.DefaultConfigFileName)
	if fn := constants.ArgConfig; viper.IsSet(fn) {
		path = viper.GetString(fn)
	}

	// We don't perform an `isAbs` check before turning it into an absolute path
	// because it internally has this check and will perform filepath.Clean on it if so
	// which is a great thing to have (sanitizes the path for multiple separators, path name elements, etc.)
	absPath, err := filepath.Abs(path)
	if err != nil {
		// just use the given provided by the user if err. Read and Write can still handle relative paths,
		// the only downside is annoyance for the user of not having his pwd prepended to `ionosctl location`
		return path
	}

	// Always prefer returning an absolute, cleaned path if possible.
	return absPath
}

func getConfigHomeDir() string {
	configPath, err := os.UserConfigDir()
	if err != nil {
		die.Die("is $HOME defined? couldn't get config dir: " + err.Error())
	}
	return filepath.Join(configPath, "ionosctl")
}
