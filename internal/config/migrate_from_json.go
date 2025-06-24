package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ionos-cloud/sdk-go-bundle/shared"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
)

// MigrateFromJSON looks for an old JSON file at jsonPath, and if it defines
// any of userdata.token, userdata.name, userdata.password, returns
// a minimal FileConfig with exactly those credentials (v1.0, single
// profile/environment). If none found, or on error, returns (nil, err).
func MigrateFromJSON(jsonPath string) (*fileconfiguration.FileConfig, error) {
	raw, err := os.ReadFile(jsonPath)
	if err != nil {
		return nil, err
	}
	var old map[string]string
	if err := json.Unmarshal(raw, &old); err != nil {
		return nil, err
	}

	creds := shared.Credentials{
		Token:    old["userdata.token"],
		Username: old["userdata.name"],
		Password: old["userdata.password"],
	}
	// bail if there’s really nothing to carry over
	if creds.Token == "" && (creds.Username == "" || creds.Password == "") {
		return nil, nil
	}

	fmt.Fprintf(os.Stderr, "WARNING: - Migrating from legacy JSON config at '%s'.\n", jsonPath)
	fmt.Fprintf(os.Stderr, "         - The JSON config file is now deprecated.\n")
	fmt.Fprintf(os.Stderr, "         - It is recommended to re-generate your config file using 'ionosctl login'\n")

	fc := &fileconfiguration.FileConfig{
		Version:        fileconfiguration.Version(1.0),
		CurrentProfile: "user",
		Profiles: []fileconfiguration.Profile{
			{
				Name:        "user",
				Environment: "prod",
				Credentials: creds,
			},
		},
		Environments: []fileconfiguration.Environment{
			{
				Name:     "prod",
				Products: nil,
			},
		},
	}
	return fc, nil
}
