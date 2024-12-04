package config_test

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"runtime"
	"testing"
	"time"

	"github.com/ionos-cloud/ionosctl/v6/internal/config"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/stretchr/testify/assert"

	"github.com/spf13/viper"
)

func TestRead(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		data     string
		perm     os.FileMode
		wantErr  bool
		errMsg   string
	}{
		{
			name:     "should read the configuration file successfully",
			filename: "config.json",
			data:     `{"key":"value", "key2": "value2", "key3": "value3"}`,
			perm:     0600,
			wantErr:  false,
		},
		{
			name:     "should return an error when the file does not exist",
			filename: "non-existent-file",
			wantErr:  true,
			errMsg:   "failed getting config file info",
		},
		{
			name:     "should return an error when the file has invalid permissions",
			filename: "bad-permissions.json",
			data:     `{"key":"value"}`,
			perm:     0777,
			wantErr:  true,
			errMsg:   "expected 600, got 777",
		},
		{
			name:     "should return an error when the file has invalid json",
			filename: "bad-json.json",
			data:     `{"key":`,
			perm:     0600,
			wantErr:  true,
			errMsg:   "failed unmarshalling config file data: unexpected end of JSON input",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.data != "" {
				tmpfile, err := ioutil.TempFile("", tt.filename)
				assert.NoError(t, err)

				defer os.Remove(tmpfile.Name())

				_, err = tmpfile.Write([]byte(tt.data))
				assert.NoError(t, err)

				if runtime.GOOS == "windows" && tt.perm != 0600 {
					// If using Windows, skip any tests related to invalid permissions.
					// Refer to os.Chmod documentation: On Windows, can only set the "read" bit of the permissions.
					// This would lead to the test 'should_return_an_error_when_the_file_has_invalid_permissions' breaking.
					t.SkipNow()
				}

				err = tmpfile.Chmod(tt.perm)
				assert.NoError(t, err)

				err = tmpfile.Close()
				assert.NoError(t, err)

				viper.Set(constants.ArgConfig, tmpfile.Name())
			} else {
				viper.Set(constants.ArgConfig, tt.filename)
			}

			cfg, err := config.Read()

			if tt.wantErr {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tt.errMsg)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, "value", cfg["key"])
				assert.Equal(t, "value2", cfg["key2"])
				assert.Equal(t, "value3", cfg["key3"])
			}
		})
	}
}

func TestWrite(t *testing.T) {
	tests := []struct {
		name    string
		data    map[string]string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "should write the configuration file successfully",
			data:    map[string]string{"key": "value", "key2": "value2", "key3": "value3"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpfile, err := ioutil.TempFile("", "config.json")
			assert.NoError(t, err)

			defer os.Remove(tmpfile.Name())

			if err := tmpfile.Close(); err != nil {
				t.Fatal(err)
			}

			viper.Set(constants.ArgConfig, tmpfile.Name())

			err = config.Write(tt.data)

			if tt.wantErr {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.errMsg)
				return
			}

			assert.NoError(t, err)

			// Validate using ioutil.ReadFile
			data, err := ioutil.ReadFile(tmpfile.Name())
			assert.NoError(t, err)

			var cfg map[string]string
			err = json.Unmarshal(data, &cfg)
			assert.NoError(t, err)

			assert.Equal(t, tt.data, cfg)

			// Validate using config.Read
			cfg, err = config.Read()
			assert.NoError(t, err)

			assert.Equal(t, tt.data, cfg)
		})
	}
}

func TestConcurrency(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", "config.json")
	assert.NoError(t, err)

	defer os.Remove(tmpfile.Name())

	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	viper.Set(constants.ArgConfig, tmpfile.Name())

	data1 := map[string]string{"key1": "value1"}
	data2 := map[string]string{"key2": "value2"}

	go func() {
		for i := 0; i < 10; i++ {
			err := config.Write(data1)
			assert.NoError(t, err)
		}
	}()

	go func() {
		for i := 0; i < 10; i++ {
			err := config.Write(data2)
			assert.NoError(t, err)
		}
	}()

	time.Sleep(1 * time.Second)

	cfg, err := config.Read()
	assert.NoError(t, err)

	// we cannot predict which write will happen last
	possibleResults := []map[string]string{data1, data2}
	assert.Contains(t, possibleResults, cfg)
}

func TestGetServerUrl(t *testing.T) {
	tests := []struct {
		name              string
		flagVal           string
		envVal            string
		cfgVal            string
		expectedServerUrl string
	}{
		{
			name:              "Flag value is used and different from default",
			flagVal:           "http://flag.url",
			envVal:            "http://env.url",
			cfgVal:            "http://cfg.url",
			expectedServerUrl: "http://flag.url",
		},
		{
			name:              "Flag value is DNS Default, return flag value",
			flagVal:           "dns.de-fra.ionos.com",
			envVal:            "http://env.url",
			cfgVal:            "http://cfg.url",
			expectedServerUrl: "dns.de-fra.ionos.com",
		},
		{
			name:              "Flag value is DNS default, all other empty, return flag value",
			flagVal:           "dns.de-fra.ionos.com",
			envVal:            "",
			cfgVal:            "",
			expectedServerUrl: "dns.de-fra.ionos.com",
		},
		{
			name:              "Flag value is empty, env and cfg set, return env value",
			flagVal:           "",
			envVal:            "dns.de-fra.ionos.com",
			cfgVal:            "dns.de-txl.ionos.com",
			expectedServerUrl: "dns.de-fra.ionos.com",
		},
		{
			name:              "Explicit flag URL is returned",
			flagVal:           "http://explicit-url.com",
			envVal:            "",
			cfgVal:            "",
			expectedServerUrl: "http://explicit-url.com",
		},
		{
			name:              "Explicit flag URL is prefered over explicit env var",
			flagVal:           "http://explicit-url.com",
			envVal:            "http://env.url",
			cfgVal:            "",
			expectedServerUrl: "http://explicit-url.com",
		},
		{
			name:              "Default API Url explicitly set is preferred over explicit env var",
			flagVal:           constants.DefaultApiURL,
			envVal:            "http://env.url",
			cfgVal:            "",
			expectedServerUrl: constants.DefaultApiURL,
		},
		{
			name:              "CFG value is preferred over defaults",
			flagVal:           "",
			envVal:            "",
			cfgVal:            "cfg-url",
			expectedServerUrl: "cfg-url",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock viper values
			viper.Set(constants.ArgServerUrl, tt.flagVal)
			viper.Set(constants.EnvServerUrl, tt.envVal)
			viper.Set(constants.CfgServerUrl, tt.cfgVal)

			got := config.GetServerUrl()
			if got != tt.expectedServerUrl {
				t.Errorf("Expected %s but got %s", tt.expectedServerUrl, got)
			}
		})
	}
}
