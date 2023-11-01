//go:build integration
// +build integration

package cfg_test

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/cilium/fake"
	"github.com/ionos-cloud/ionosctl/v6/commands/cfg"
	"github.com/ionos-cloud/ionosctl/v6/internal/config"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/spf13/viper"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/stretchr/testify/assert"
)

var (
	GoodUsername    = ""
	GoodPassword    = ""
	GoodToken       = ""
	cl              *client.Client
	tokCreationTime time.Time

	testDir string
)

func TestAuthCmds(t *testing.T) {
	if err := setup(); err != nil {
		t.Fatalf("Failed setting up auth tests: %s", err)
	}
	t.Cleanup(teardown)

	assert.NotEmpty(t, GoodUsername)
	assert.NotEmpty(t, GoodPassword)
	assert.NotEmpty(t, GoodToken)

	viper.Set(constants.ArgOutput, "text")
	viper.Set(constants.ArgForce, "true")

	_, filename, _, _ := runtime.Caller(0)
	testDir = strings.Split(filepath.Dir(filename), "/ionosctl/")[0] + "/ionosctl/temp-login-tests"
	err := os.MkdirAll(testDir, 0777)
	if err != nil {
		t.Fatalf("failed creating config dir: %s", err.Error())
	}
	viper.Set(constants.ArgConfig, testDir+"/"+fake.Adjective()+".json")
	err = config.Write(map[string]string{
		constants.CfgUsername:  "UsernameHere",
		constants.CfgPassword:  "PasswordHere",
		constants.CfgToken:     "TokenHere",
		constants.CfgServerUrl: "sample-url.com",
	})
	if err != nil {
		t.Fatalf("failed setting up a test config file: %s", err.Error())
	}

	t.Parallel()

	t.Run("login test user interactively, password as flag - valid token saved to config", func(t *testing.T) {
		login := cfg.LoginCmd()
		logout := cfg.LogoutCmd()

		login.Command.Flags().Set(constants.ArgPassword, GoodPassword)

		out := &bytes.Buffer{}
		in := bytes.NewBufferString(GoodUsername + "\n")
		login.Command.SetOut(out)
		login.Command.SetIn(in)

		// Exec login
		err := login.Command.Execute()
		assert.Contains(t, out.String(), "Enter your username: ")
		assert.NoError(t, err)

		// Read the configuration after login and assert that the token is valid
		cfg, err := config.Read()
		assert.NoError(t, err)

		assert.NoError(t, client.TestCreds("", "", cfg[constants.CfgToken]))
		assert.Empty(t, cfg[constants.ArgServerUrl])
		assert.Empty(t, cfg[constants.CfgUsername])
		assert.Empty(t, cfg[constants.CfgPassword])

		// Exec logout
		out = &bytes.Buffer{}
		logout.Command.SetOut(out)
		err = logout.Command.Execute()

		// Assert that logout was successful and the token is removed from the configuration
		assert.Contains(t, out.String(), "De-authentication successful")
		assert.NoError(t, err)
		cfg, err = config.Read()
		assert.NoError(t, err)
		assert.Empty(t, cfg[constants.CfgToken])
	})

	t.Run("login test user, password as flag - valid token saved to config", func(t *testing.T) {
		login := cfg.LoginCmd()
		logout := cfg.LogoutCmd()

		viper.Set(core.GetFlagName(login.NS, constants.ArgUser), GoodUsername)
		viper.Set(core.GetFlagName(login.NS, constants.ArgPassword), GoodPassword)

		out := &bytes.Buffer{}
		login.Command.SetOut(out)
		err := login.Command.Execute()
		assert.Contains(t, out.String(), "Config file updated successfully")
		assert.NoError(t, err)

		// Read the configuration after login and assert that the token is valid
		cfg, err := config.Read()
		assert.NoError(t, err)
		assert.NoError(t, client.TestCreds("", "", cfg[constants.CfgToken]))

		out = &bytes.Buffer{}
		logout.Command.SetOut(out)
		err = logout.Command.Execute()
		assert.Contains(t, out.String(), "De-authentication successful")
		assert.NoError(t, err)
		cfg, err = config.Read()
		assert.NoError(t, err)
		assert.Empty(t, cfg[constants.CfgToken])
	})

	t.Run("invalid user, password as flag - valid token saved to config - should still error out", func(t *testing.T) {
		login := cfg.LoginCmd()
		logout := cfg.LogoutCmd()

		viper.Set(core.GetFlagName(login.NS, constants.ArgUser), fake.Adjective())
		viper.Set(core.GetFlagName(login.NS, constants.ArgPassword), fake.Adjective())

		config.Write(map[string]string{
			constants.CfgToken: GoodToken,
		})

		// Read the configuration after login and assert that the token is valid
		cfg, err := config.Read()
		assert.NoError(t, err)
		assert.NoError(t, client.TestCreds("", "", cfg[constants.CfgToken]))

		out := &bytes.Buffer{}
		login.Command.SetOut(out)
		err = login.Command.Execute()
		assert.ErrorContains(t, err, "401 Unauthorized")

		// Recreate file that is replaced by above login
		config.Write(map[string]string{
			constants.CfgToken: GoodToken,
		})

		out = &bytes.Buffer{}
		logout.Command.SetOut(out)
		err = logout.Command.Execute()
		assert.Contains(t, out.String(), "De-authentication successful")
		assert.NoError(t, err)
		cfg, err = config.Read()
		assert.NoError(t, err)
		assert.Empty(t, cfg[constants.CfgToken])
	})

	t.Run("non existant config file should be created", func(t *testing.T) {
		login := cfg.LoginCmd()
		logout := cfg.LogoutCmd()

		viper.Set(core.GetFlagName(login.NS, constants.ArgUser), GoodUsername)
		viper.Set(core.GetFlagName(login.NS, constants.ArgPassword), GoodPassword)
		newCfgFilePath := testDir + "/" + fake.Noun() + ".json"
		viper.Set(constants.ArgConfig, newCfgFilePath)

		out := &bytes.Buffer{}
		login.Command.SetOut(out)
		err := login.Command.Execute()
		assert.NoError(t, err)

		// Read the configuration after login and assert that the token is valid
		cfg, err := config.Read()
		assert.NoError(t, err)
		assert.NoError(t, client.TestCreds("", "", cfg[constants.CfgToken]))

		out = &bytes.Buffer{}
		logout.Command.SetOut(out)
		err = logout.Command.Execute()
		assert.Contains(t, out.String(), "De-authentication successful")
		assert.NoError(t, err)
		cfg, err = config.Read()
		assert.NoError(t, err)
		assert.Empty(t, cfg[constants.CfgToken])

		// err = os.Remove(newCfgFilePath)
		// if err != nil {
		// 	t.Fatalf("note: failed to remove file %s: %s, remove it manually\n", newCfgFilePath, err.Error())
		// }
	})

	t.Run("Pre-rework config file logout - Username and password removed from config file", func(t *testing.T) {
		logout := cfg.LogoutCmd()

		before := map[string]string{
			constants.CfgUsername:  "UsernameHere",
			constants.CfgPassword:  "PasswordHere",
			constants.CfgToken:     "TokenHere",
			constants.CfgServerUrl: "dont-kill-me.com",
		}
		config.Write(before)

		// Read the configuration - is it what we expect ?
		cfg, err := config.Read()
		assert.NoError(t, err)
		assert.Equal(t, before, cfg)

		out := &bytes.Buffer{}
		logout.Command.SetOut(out)
		err = logout.Command.Execute()

		// Assert that logout was successful and the username, password, token removed from cfg
		assert.Contains(t, out.String(), "De-authentication successful")
		assert.NoError(t, err)
		after, err := config.Read()
		assert.NoError(t, err)
		assert.Empty(t, after[constants.CfgToken])
		assert.Empty(t, after[constants.CfgPassword])
		assert.Empty(t, after[constants.CfgUsername])
		assert.Equal(t, before[constants.CfgServerUrl], after[constants.CfgServerUrl])
	})

	// cfg location tests
	t.Run("cfg location cmd returns valid location", func(t *testing.T) {
		cfgLocCmd := cfg.LocationCmd()
		out := &bytes.Buffer{}
		cfgLocCmd.Command.SetOut(out)
		err := cfgLocCmd.Command.Execute()
		assert.NoError(t, err)

		assert.Equal(t, config.GetConfigFilePath(), out.String())
	})

	// cfg whoami tests
	t.Run("cfg whoami returns current user - Env Username & Password", func(t *testing.T) {
		cmd := cfg.WhoamiCmd()

		viper.Set(core.GetFlagName(cmd.NS, constants.ArgUser), GoodUsername)
		viper.Set(core.GetFlagName(cmd.NS, constants.ArgPassword), GoodPassword)

		out := &bytes.Buffer{}
		cmd.Command.SetOut(out)
		err := cmd.Command.Execute()
		assert.NoError(t, err)

		assert.Contains(t, out.String(), GoodUsername)
	})

}

func setup() error {
	GoodUsername = os.Getenv("IONOS_USERNAME")
	GoodPassword = os.Getenv("IONOS_PASSWORD")

	if GoodUsername == "" || GoodPassword == "" {
		return fmt.Errorf("empty user/pass")
	}

	cl = client.NewClient(GoodUsername, GoodPassword, "", "")
	tok, _, err := cl.AuthClient.TokensApi.TokensGenerate(context.Background()).Execute()

	if err != nil {
		return err
	}

	if tok.Token == nil {
		return fmt.Errorf("tok is nil")
	}

	GoodToken = *tok.Token
	tokCreationTime = time.Now().In(time.UTC).Add(-1 * time.Minute)

	return nil
}

func teardown() {
	toks, _, err := cl.AuthClient.TokensApi.TokensGet(context.Background()).Execute()

	if err != nil {
		panic(err)
	}

	// Delete tokens generated since setup
	for _, t := range *toks.Tokens {
		strDate, ok := strings.CutSuffix(*t.CreatedDate, "[UTC]")
		if !ok {
			panic("they changed the date format: no more [UTC] suffix")
		}
		date, err := time.Parse(time.RFC3339, strDate)
		if err != nil {
			panic(fmt.Errorf("they changed the date format: %w", err))
		}

		// Delete the token if it was created after setup
		if date.After(tokCreationTime) {
			_, _, err := cl.AuthClient.TokensApi.TokensDeleteById(context.Background(), *t.Id).Execute()

			if err != nil {
				panic(err)
			}
		}
	}

	err = os.RemoveAll(testDir)
	if err != nil {
		fmt.Printf("failed cleaning up %s: %s", testDir, err.Error())
	}

}
