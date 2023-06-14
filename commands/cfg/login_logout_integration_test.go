//go:build integration
// +build integration

package cfg_test

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/ionos-cloud/ionosctl/v6/commands/cfg"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/pkg/config"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/stretchr/testify/assert"
)

var (
	GoodUsername    = ""
	GoodPassword    = ""
	GoodToken       = ""
	cl              *client.Client
	tokCreationTime time.Time
)

func TestAuthCmds(t *testing.T) {
	if err := setup(); err != nil {
		t.Fatalf("Failed setting up auth tests: %s", err)
	}
	t.Cleanup(teardown)

	assert.NotEmpty(t, GoodUsername)
	assert.NotEmpty(t, GoodPassword)
	assert.NotEmpty(t, GoodToken)

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

		login.Command.Flags().Set(constants.ArgUser, GoodUsername)
		login.Command.Flags().Set(constants.ArgPassword, GoodPassword)

		out := &bytes.Buffer{}
		login.Command.SetOut(out)
		err := login.Command.Execute()
		assert.Contains(t, out.String(), "Authentication successful")
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

	t.Run("login test token as flag", func(t *testing.T) {
		login := cfg.LoginCmd()
		logout := cfg.LogoutCmd()

		login.Command.Flags().Set(constants.ArgToken, GoodToken)

		out := &bytes.Buffer{}
		login.Command.SetOut(out)
		err := login.Command.Execute()
		assert.Contains(t, out.String(), "Authentication successful")
		assert.NoError(t, err)

		// Read the configuration after login and assert that the token is valid
		cfg, err := config.Read()
		assert.NoError(t, err)
		assert.NoError(t, client.TestCreds("", "", cfg[constants.CfgToken]))

		// In the case token is provided by user via --token ; then the saved cfg file token should be identical to the provided one
		assert.Equal(t, cfg[constants.CfgToken], GoodToken)

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

	t.Run("Pre-june config file logout - Username and password removed from config file", func(t *testing.T) {
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

	t.Run("cfg location cmd returns valid location", func(t *testing.T) {
		cfgLocCmd := cfg.CfgLocationCmd()
		out := &bytes.Buffer{}
		cfgLocCmd.Command.SetOut(out)
		err := cfgLocCmd.Command.Execute()
		assert.NoError(t, err)

		assert.Equal(t, config.GetConfigFile(), out.String())

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

}
