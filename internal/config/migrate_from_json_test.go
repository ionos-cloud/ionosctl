package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMigrateFromJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")

	// 1. missing file
	fc, err := MigrateFromJSON(path)
	assert.Error(t, err) // file-not-found
	assert.Nil(t, fc)

	// 2. invalid JSON
	os.WriteFile(path, []byte("{bad-json"), 0o600)
	fc, err = MigrateFromJSON(path)
	assert.Error(t, err)

	// 3. empty creds -> nil,nil
	os.WriteFile(path, []byte(`{"foo":"bar"}`), 0o600)
	fc, err = MigrateFromJSON(path)
	assert.NoError(t, err)
	assert.Nil(t, fc)

	// 4. token only
	os.WriteFile(path, []byte(`{"userdata.token":"T"}`), 0o600)
	fc, err = MigrateFromJSON(path)
	assert.NoError(t, err)
	assert.NotNil(t, fc)
	assert.Equal(t, "T", fc.Profiles[0].Credentials.Token)

	// 5. user+pass only
	os.WriteFile(path, []byte(`{"userdata.name":"U","userdata.password":"P"}`), 0o600)
	fc, err = MigrateFromJSON(path)
	assert.NoError(t, err)
	assert.NotNil(t, fc)
	creds := fc.Profiles[0].Credentials
	assert.Equal(t, "U", creds.Username)
	assert.Equal(t, "P", creds.Password)

	// 6. all three
	os.WriteFile(path, []byte(`{"userdata.token":"X","userdata.name":"Y","userdata.password":"Z"}`), 0o600)
	fc, err = MigrateFromJSON(path)
	assert.NoError(t, err)
	assert.Equal(t, "X", fc.Profiles[0].Credentials.Token)
	assert.Equal(t, "Y", fc.Profiles[0].Credentials.Username)
	assert.Equal(t, "Z", fc.Profiles[0].Credentials.Password)
}
