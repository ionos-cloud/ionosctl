//go:build integration
// +build integration

package commands_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/spf13/pflag"

	"github.com/cilium/fake"
	"github.com/ionos-cloud/ionosctl/v6/commands"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	GoodUsername = os.Getenv("IONOS_USERNAME")
	GoodPassword = os.Getenv("IONOS_PASSWORD")
	GoodToken    = ""
)

func pre(t *testing.T) {
	cl, err := client.NewClient(GoodUsername, GoodPassword, "", "")
	if err != nil {
		t.Skipf(fmt.Errorf("failed setting up login tests: %w. Make sure IONOS_USERNAME and IONOS_PASSWORD are set", err).Error())
	}
	jwt, _, err := cl.AuthClient.TokensApi.TokensGenerate(context.Background()).Execute()
	if err != nil {
		t.Skipf(fmt.Errorf("failed setting up login tests: %w. Make sure IONOS_USERNAME and IONOS_PASSWORD are set", err).Error())
	}
	GoodToken = *jwt.Token

	panic("pre")

	os.Clearenv()
	viper.Reset()
	viper.Set(constants.ArgOutput, "text")
}

func TestLoginFlagsUserPassGood(t *testing.T) {
	pre(t)

	c := commands.LoginCmd()
	c.Command.Flags().Set(constants.ArgUser, GoodUsername)
	c.Command.Flags().Set(constants.ArgPassword, GoodPassword)
	err := fmt.Errorf("")
	c.Command.Flags().Visit(func(flag *pflag.Flag) {
		err = errors.Join(err, fmt.Errorf("%s : %s", flag.Name, flag.Value))
	})
	//err := c.Command.Execute()
	panic(err)
	assert.ErrorContains(t, err, "failed using username and password to generate a token: 401 Unauthorized")
}

func TestLoginInteractive(t *testing.T) {
	pre(t)
	// TODO: setting output of cobra command doesnt change anything. cant test output
	// TODO: gotta do weird workarounds for changing the input
}

func TestLoginFlagsUserPass401(t *testing.T) {
	pre(t)

	c := commands.LoginCmd()
	c.Command.Flags().Set(constants.ArgUser, fake.Adjective())
	c.Command.Flags().Set(constants.ArgPassword, fake.Adjective())
	err := c.Command.Execute()
	assert.ErrorContains(t, err, "failed using username and password to generate a token: 401 Unauthorized")
}

func TestLoginFlagsToken401(t *testing.T) {
	pre(t)

	c := commands.LoginCmd()
	c.Command.Flags().Set(constants.ArgToken, fake.AlphaNum(32))
	err := c.Command.Execute()
	assert.ErrorContains(t, err, "401 Unauthorized")
}
