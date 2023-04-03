package client

import (
	"os"
	"testing"

	sdk "github.com/ionos-cloud/sdk-go-bundle/shared"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestGetClient(t *testing.T) {
	type args struct {
		name    string
		pwd     string
		token   string
		hostUrl string
	}
	tests := []struct {
		name    string
		runs    int
		args    args
		want    *Client
		wantErr bool
	}{
		{"MissingCredentials", 1, args{"", "", "", ""}, nil, true},
		{"MultipleGetClients", 4, args{"user", "pass", "token", "url"}, instance, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Clearenv()
			viper.Reset()

			assert.NoError(t, os.Setenv(sdk.IonosUsernameEnvVar, tt.args.name))
			assert.NoError(t, os.Setenv(sdk.IonosPasswordEnvVar, tt.args.pwd))
			assert.NoError(t, os.Setenv(sdk.IonosTokenEnvVar, tt.args.token))
			assert.NoError(t, os.Setenv(sdk.IonosApiUrlEnvVar, tt.args.hostUrl))

			for i := 0; i < tt.runs; i++ {
				client, err := Get()
				if !tt.wantErr && err != nil {
					t.Errorf("Did not expect error: %v", err)
				}
				assert.Equalf(t, tt.want, client, "newClient(%v, %v, %v, %v)", tt.args.name, tt.args.pwd, tt.args.token, tt.args.hostUrl)
			}
		})
	}
}
