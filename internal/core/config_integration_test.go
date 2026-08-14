package core

import (
	"testing"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// withClientConfig temporarily installs cfg as the singleton client's config so
// findOverriddenURL can read config-file overrides during unit tests.
func withClientConfig(t *testing.T, cfg *fileconfiguration.FileConfig) {
	t.Helper()
	cl := client.Must(func(error) {})
	prev := cl.Config
	cl.Config = cfg
	t.Cleanup(func() { cl.Config = prev })
}

// osConfig builds an in-memory config exposing a single object-storage endpoint
// override for the given location spelling.
func osConfig(location, url string) *fileconfiguration.FileConfig {
	return &fileconfiguration.FileConfig{
		CurrentProfile: "p",
		Profiles:       []fileconfiguration.Profile{{Name: "p", Environment: "env"}},
		Environments: []fileconfiguration.Environment{{
			Name: "env",
			Products: []fileconfiguration.Product{{
				Name:      fileconfiguration.ObjectStorage,
				Endpoints: []fileconfiguration.Endpoint{{Location: location, Name: url}},
			}},
		}},
	}
}

// TestFindOverriddenURL_LocationFormatMismatch guards the object-storage config
// override fix: the config file may spell a region with slashes (eu/central/3)
// while the command flag uses dashes (eu-central-3), or vice versa. The override
// must resolve regardless of which convention each side uses.
func TestFindOverriddenURL_LocationFormatMismatch(t *testing.T) {
	const tmpl = constants.ObjectStorageApiRegionalURL
	const cfgURL = "https://cfg-override.example.com"

	t.Setenv(constants.EnvServerUrl, "")

	newCmd := func() *cobra.Command {
		cmd := &cobra.Command{Use: "get"}
		cmd.Flags().String(constants.ArgServerUrl, tmpl, "")
		return cmd
	}

	t.Run("config uses slashes, flag uses dashes", func(t *testing.T) {
		withClientConfig(t, osConfig("eu/central/3", cfgURL))
		got := findOverriddenURL(newCmd(), []string{fileconfiguration.ObjectStorage}, tmpl, "eu-central-3")
		if got != cfgURL {
			t.Errorf("url = %q, want config override %q", got, cfgURL)
		}
	})

	t.Run("config uses dashes, flag uses slashes", func(t *testing.T) {
		withClientConfig(t, osConfig("eu-central-3", cfgURL))
		got := findOverriddenURL(newCmd(), []string{fileconfiguration.ObjectStorage}, tmpl, "eu/central/3")
		if got != cfgURL {
			t.Errorf("url = %q, want config override %q", got, cfgURL)
		}
	})

	t.Run("no override falls back to normalized template", func(t *testing.T) {
		withClientConfig(t, nil)
		got := findOverriddenURL(newCmd(), []string{fileconfiguration.ObjectStorage}, tmpl, "eu/central/3")
		if got != "https://s3.eu-central-3.ionoscloud.com" {
			t.Errorf("url = %q, want per-location template", got)
		}
	})
}

// TestWithRegionalConfigOverride_DefaultLocationHonorsOverride guards the fix for
// single-resource commands run without --location: the override for the default
// (first allowed) region must be applied, not just the bare template URL.
func TestWithRegionalConfigOverride_DefaultLocationHonorsOverride(t *testing.T) {
	const cfgURL = "https://default-loc-override.example.com"
	// object-storage's first allowed location is eu-central-3; store it slash-form
	// to also exercise the format-mismatch path.
	withClientConfig(t, osConfig("eu/central/3", cfgURL))
	t.Setenv(constants.EnvServerUrl, "")

	c := &Command{Command: &cobra.Command{Use: "os"}}
	c = WithRegionalConfigOverride(c, []string{fileconfiguration.ObjectStorage},
		constants.ObjectStorageApiRegionalURL, constants.ObjectStorageLocations)

	prev := viper.GetString(constants.ArgServerUrl)
	t.Cleanup(func() { viper.Set(constants.ArgServerUrl, prev) })
	viper.Set(constants.ArgServerUrl, "")

	// No --location provided.
	if err := c.Command.PersistentPreRunE(c.Command, nil); err != nil {
		t.Fatalf("PersistentPreRunE: %v", err)
	}
	if got := viper.GetString(constants.ArgServerUrl); got != cfgURL {
		t.Errorf("server URL = %q, want default-location override %q", got, cfgURL)
	}
}
