package client

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
	"github.com/spf13/viper"
)

// writeMinimalYAML writes a bare-bones valid FileConfig YAML to path.
func writeMinimalYAML(t *testing.T, path string) {
	const yamlContent = `version: 1.0
currentProfile: ""
profiles: []
environments: []
`
	if err := os.WriteFile(path, []byte(yamlContent), 0o600); err != nil {
		t.Fatalf("failed to write %q: %v", path, err)
	}
}

func TestRetrieveConfigFile_FlagOverridesEverything(t *testing.T) {
	// prepare a temp config file
	dir := t.TempDir()
	flagCfg := filepath.Join(dir, "flag.yaml")
	writeMinimalYAML(t, flagCfg)

	// set the --config flag
	viper.Set(constants.ArgConfig, flagCfg)
	// ensure env var cannot interfere
	t.Setenv(shared.IonosFilePathEnvVar, "")
	// ensure no default exists
	t.Setenv("HOME", dir) // no ~/.ionos/config

	src, err := retrieveConfigFile()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if src.Path != flagCfg {
		t.Errorf("expected Path %q, got %q", flagCfg, src.Path)
	}
	if src.Config == nil {
		t.Errorf("expected Config non-nil when flag file exists")
	}
}

func TestRetrieveConfigFile_EnvVarAsFallback(t *testing.T) {
	// clear any --config flag
	viper.Set(constants.ArgConfig, "")

	// prepare env-pointed config
	dir := t.TempDir()
	envCfg := filepath.Join(dir, "env.yaml")
	writeMinimalYAML(t, envCfg)

	// set IONOS_CONFIG_FILE
	t.Setenv(shared.IonosFilePathEnvVar, envCfg)
	// clear any SDK default
	t.Setenv("HOME", dir)

	src, err := retrieveConfigFile()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if src.Path != envCfg {
		t.Errorf("expected Path %q, got %q", envCfg, src.Path)
	}
	if src.Config == nil {
		t.Errorf("expected Config non-nil when env file exists")
	}
}

func TestRetrieveConfigFile_SDKDefault(t *testing.T) {
	// clear flag + env
	viper.Set(constants.ArgConfig, "")
	t.Setenv(shared.IonosFilePathEnvVar, "")

	// point HOME at a dir with a .ionos/config
	home := t.TempDir()
	dfltDir := filepath.Join(home, ".ionos")
	if err := os.MkdirAll(dfltDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	defaultCfg := filepath.Join(dfltDir, "config")
	writeMinimalYAML(t, defaultCfg)
	t.Setenv("HOME", home)

	src, err := retrieveConfigFile()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if src.Path != defaultCfg {
		t.Errorf("expected Path %q, got %q", defaultCfg, src.Path)
	}
	if src.Config == nil {
		t.Errorf("expected Config non-nil when default exists")
	}
}

func TestRetrieveConfigFile_NoneFound(t *testing.T) {
	// clear everything
	viper.Set(constants.ArgConfig, "")
	t.Setenv(shared.IonosFilePathEnvVar, "")
	// point HOME at empty dir
	home := t.TempDir()
	t.Setenv("HOME", home)

	src, err := retrieveConfigFile()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// when nothing is found, Path should be whatever viper.ArgConfig says (here, empty)
	if want, got := "", src.Path; want != got {
		t.Errorf("expected empty Path, got %q", got)
	}
	if src.Config != nil {
		t.Errorf("expected Config=nil when no file exists")
	}
}
