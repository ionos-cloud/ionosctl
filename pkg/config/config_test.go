package config_test

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/ionos-cloud/ionosctl/v6/pkg/config"
	"github.com/stretchr/testify/assert"

	"github.com/spf13/viper"

	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
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
			errMsg:   "no such file or directory",
		},
		{
			name:     "should return an error when the file has invalid permissions",
			filename: "bad-permissions.json",
			data:     `{"key":"value"}`,
			perm:     0700,
			wantErr:  true,
			errMsg:   "expected 600, got 700",
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
