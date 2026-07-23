package client

import (
	"os"
	"path/filepath"
	"testing"

	cfg "github.com/ionos-cloud/ionosctl/v6/internal/config"
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

// TestRetrieveConfigFile_EnvVarNotShadowedByDefaultFlag guards the fix for the
// --config flag's computed default silently shadowing IONOS_CONFIG_FILE. When
// --config still holds its default (user did not pass it explicitly) and the env
// var is set, the env var must win - even if the default config file exists.
func TestRetrieveConfigFile_EnvVarNotShadowedByDefaultFlag(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	t.Setenv("XDG_CONFIG_HOME", "")
	t.Setenv("APPDATA", filepath.Join(home, "AppData", "Roaming"))

	// Create the default config file (the --config flag's default target).
	defaultPath := cfg.DefaultConfigFilePath()
	if err := os.MkdirAll(filepath.Dir(defaultPath), 0o755); err != nil {
		t.Fatalf("mkdir default: %v", err)
	}
	writeMinimalYAML(t, defaultPath)

	// The flag holds its computed default (user did not pass --config).
	viper.Set(constants.ArgConfig, defaultPath)
	t.Cleanup(func() { viper.Set(constants.ArgConfig, "") })

	// IONOS_CONFIG_FILE points elsewhere and must take precedence.
	envCfg := filepath.Join(home, "env.yaml")
	writeMinimalYAML(t, envCfg)
	t.Setenv(shared.IonosFilePathEnvVar, envCfg)

	src, err := retrieveConfigFile()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if src.Path != envCfg {
		t.Errorf("env var shadowed by default flag: expected Path %q, got %q", envCfg, src.Path)
	}
}

// TestRetrieveConfigFile_ExplicitFlagWinsOverEnv ensures an explicitly-set
// --config (different from the default) still takes precedence over the env var.
func TestRetrieveConfigFile_ExplicitFlagWinsOverEnv(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	t.Setenv("XDG_CONFIG_HOME", "")
	t.Setenv("APPDATA", filepath.Join(home, "AppData", "Roaming"))

	flagCfg := filepath.Join(home, "flag.yaml")
	writeMinimalYAML(t, flagCfg)
	viper.Set(constants.ArgConfig, flagCfg)
	t.Cleanup(func() { viper.Set(constants.ArgConfig, "") })

	envCfg := filepath.Join(home, "env.yaml")
	writeMinimalYAML(t, envCfg)
	t.Setenv(shared.IonosFilePathEnvVar, envCfg)

	src, err := retrieveConfigFile()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if src.Path != flagCfg {
		t.Errorf("explicit --config should win over env var: expected %q, got %q", flagCfg, src.Path)
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
	t.Setenv("USERPROFILE", home)
	t.Setenv("APPDATA", filepath.Join(home, "AppData", "Roaming"))

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
	t.Setenv("USERPROFILE", home)
	t.Setenv("APPDATA", filepath.Join(home, "AppData", "Roaming"))

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
