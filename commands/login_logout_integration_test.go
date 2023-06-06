//go:build integration
// +build integration

package commands_test

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/ionos-cloud/ionosctl/v6/commands"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/pkg/config"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/stretchr/testify/assert"
)

var (
	GoodUsername = ""
	GoodPassword = ""
	GoodToken    = ""
)

func TestAuthCmds(t *testing.T) {
	//if err := setup(); err != nil {
	//	t.Fatalf("Failed setting up auth tests: %s", err)
	//}
	//teardown()

	assert.NotEmpty(t, GoodUsername)
	assert.NotEmpty(t, GoodPassword)
	assert.NotEmpty(t, GoodToken)

	t.Parallel()

	_ = commands.LogoutCmd()

	t.Run("login test interactive input", func(t *testing.T) {
		login := commands.LoginCmd()

		out := &bytes.Buffer{}
		in := bytes.NewBufferString("MockUsername\nMockPassword\n")
		login.Command.SetOut(out)
		login.Command.SetIn(in)
		err := login.Command.Execute()
		assert.Contains(t, out.String(), "Enter your username: ")
		assert.Contains(t, out.String(), "Enter your password: ")
		// I tried mocking the terminal reader, but sadly can't inject it without major code rewrite
		assert.ErrorContains(t, err, "the set input does not have a file descriptor (is it set to a terminal?)")
	})

	t.Run("login test user interactively, password as flag", func(t *testing.T) {
		login := commands.LoginCmd()

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
	})

	t.Run("login test user, password as flag", func(t *testing.T) {
		login := commands.LoginCmd()

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
	})

	t.Run("login test user, password as flag", func(t *testing.T) {
		login := commands.LoginCmd()

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
	})

	t.Run("login test token as flag", func(t *testing.T) {
		login := commands.LoginCmd()

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
	})

	t.Run("login test user, pass, token", func(t *testing.T) {
		login := commands.LoginCmd()

		login.Command.Flags().Set(constants.ArgUser, GoodUsername)
		login.Command.Flags().Set(constants.ArgPassword, GoodPassword)
		login.Command.Flags().Set(constants.ArgToken, GoodToken)

		err := login.Command.Execute()
		assert.ErrorContains(t, err, "use either --user and/or --password, either --token")

	})

	toks, _, err := client.Must(func(err error) {
		panic(err)
	}).AuthClient.TokensApi.TokensGet(context.Background()).Execute()
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

	tok, _, err := client.Must(func(err error) {
		panic(err)
	}).AuthClient.TokensApi.TokensGenerate(context.Background()).Execute()

	GoodToken = *tok.Token

	return err

	// TODO: Mark tok generation time
}

func teardown() {

	// TODO: Delete tokens generated since setup
}
