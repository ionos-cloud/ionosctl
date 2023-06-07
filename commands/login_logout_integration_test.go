//go:build integration
// +build integration

package commands_test

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/ionos-cloud/ionosctl/v6/commands"
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
	teardown()

	assert.NotEmpty(t, GoodUsername)
	assert.NotEmpty(t, GoodPassword)
	assert.NotEmpty(t, GoodToken)

	t.Parallel()

	t.Run("login test user interactively, password as flag - valid token saved to config", func(t *testing.T) {
		login := commands.LoginCmd()
		logout := commands.LogoutCmd()

		login.Command.Flags().Set(constants.ArgPassword, GoodPassword)

		out := &bytes.Buffer{}
		in := bytes.NewBufferString(GoodUsername + "\n")
		login.Command.SetOut(out)
		login.Command.SetIn(in)
		err := login.Command.Execute()
		assert.Contains(t, out.String(), "Enter your username: ")
		assert.NoError(t, err)

		cfg, err := config.Read()
		assert.NoError(t, err)
		assert.NoError(t, client.TestCreds("", "", cfg[constants.CfgToken]))
		assert.Empty(t, cfg[constants.ArgServerUrl])
		assert.Empty(t, cfg[constants.CfgUsername])
		assert.Empty(t, cfg[constants.CfgPassword])

		out = &bytes.Buffer{}
		logout.Command.SetOut(out)
		err = logout.Command.Execute()
		assert.Contains(t, out.String(), "De-authentication successful")
		assert.NoError(t, err)
		cfg, err = config.Read()
		assert.NoError(t, err)
		assert.Empty(t, cfg[constants.CfgToken])
	})

	t.Run("login test user, password as flag - valid token saved to config", func(t *testing.T) {
		login := commands.LoginCmd()
		logout := commands.LogoutCmd()

		login.Command.Flags().Set(constants.ArgUser, GoodUsername)
		login.Command.Flags().Set(constants.ArgPassword, GoodPassword)

		out := &bytes.Buffer{}
		login.Command.SetOut(out)
		err := login.Command.Execute()
		assert.Contains(t, out.String(), "Authentication successful")
		assert.NoError(t, err)

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
		login := commands.LoginCmd()
		logout := commands.LogoutCmd()

		login.Command.Flags().Set(constants.ArgToken, GoodToken)

		out := &bytes.Buffer{}
		login.Command.SetOut(out)
		err := login.Command.Execute()
		assert.Contains(t, out.String(), "Authentication successful")
		assert.NoError(t, err)

		cfg, err := config.Read()
		assert.NoError(t, err)
		assert.NoError(t, client.TestCreds("", "", cfg[constants.CfgToken]))

		assert.Equal(t, cfg[constants.CfgToken], GoodToken)

		out = &bytes.Buffer{}
		logout.Command.SetOut(out)
		err = logout.Command.Execute()
		assert.Contains(t, out.String(), "De-authentication successful")
		assert.NoError(t, err)
		cfg, err = config.Read()
		assert.NoError(t, err)
		assert.Empty(t, cfg[constants.CfgToken])
	})

	toks, _, err := cl.AuthClient.TokensApi.TokensGet(context.Background()).Execute()
	if err != nil {
		return
	}

	msg := ""

	for _, t := range *toks.Tokens {
		date := t.CreatedDate
		msg += fmt.Sprintf("Tok %s created at %s\n", *t.Id, *date)
	}

	panic(msg)
}

func setup() error {
	GoodUsername = os.Getenv("IONOS_USERNAME")
	GoodPassword = os.Getenv("IONOS_PASSWORD")

	if GoodUsername == "" || GoodPassword == "" {
		return fmt.Errorf("empty user/pass")
	}

	cl, _ = client.NewClient(GoodUsername, GoodPassword, "", "")
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
