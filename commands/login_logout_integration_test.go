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
	if err := setup(); err != nil {
		t.Fatalf("Failed setting up auth tests: %s", err)
	}
	t.Cleanup(teardown)

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
	})
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
}

func teardown() {
}
