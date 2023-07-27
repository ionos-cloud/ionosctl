//go:build integration
// +build integration

package token_test

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/ionos-cloud/ionosctl/v6/commands/token"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	authservice "github.com/ionos-cloud/ionosctl/v6/services/auth-v1"
	sdkgoauth "github.com/ionos-cloud/sdk-go-auth"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	cl                   *client.Client
	tokFirstCreationTime time.Time
	testToken            sdkgoauth.Token
	tokenContent         sdkgoauth.Jwt
)

func TestTokenCommands(t *testing.T) {
	err := setup()
	if err != nil {
		assert.FailNow(t, err.Error())
	}

	viper.Set(constants.ArgOutput, "text")

	testCreateToken(t)
	testListTokens(t)
	testGetTokens(t)
	testParseToken(t)
	testDeleteTokens(t)
}

func setup() error {
	username := os.Getenv("IONOS_USERNAME")
	password := os.Getenv("IONOS_PASSWORD")

	if username == "" || password == "" {
		return fmt.Errorf("empty user/password")
	}

	var err error

	cl, err = client.NewTestClient(username, password, "", "")
	if err != nil {
		return err
	}

	tokenContent, _, err = cl.AuthClient.TokensApi.TokensGenerate(context.Background()).Execute()
	time.Sleep(2 * time.Second)

	return err
}

func testCreateToken(t *testing.T) {
	var err error

	tokFirstCreationTime = time.Now().In(time.UTC)
	viper.Set(constants.ArgQuiet, true)

	c := token.TokenPostCmd()
	err = c.Command.Execute()
	assert.NoError(t, err)

	time.Sleep(2 * time.Second)

	tokens, _, err := client.Must().AuthClient.TokensApi.TokensGet(context.Background()).Execute()
	assert.NoError(t, err)

	allTokens, ok := tokens.GetTokensOk()
	assert.NotEmpty(t, ok)
	assert.NotEmpty(t, *allTokens)

	var foundTokenViaSdk *sdkgoauth.Token
	foundTokenViaSdk = nil

	for _, tok := range *allTokens {
		strDate, ok := strings.CutSuffix(*tok.CreatedDate, "[UTC]")
		if !ok {
			panic("they changed the time format: no more [UTC] suffix")
		}

		date, err := time.Parse(time.RFC3339, strDate)
		if err != nil {
			panic(fmt.Errorf("they changed the date format: %w", err))
		}

		if date.After(tokFirstCreationTime) {
			temp := tok
			foundTokenViaSdk = &temp
		}
	}
	assert.NotNil(t, foundTokenViaSdk)

	testToken = *foundTokenViaSdk

	viper.Reset()
	viper.Set(constants.ArgOutput, "text")
}

func testListTokens(t *testing.T) {
	var err error

	c := token.TokenListCmd()
	err = c.Command.Execute()
	assert.NoError(t, err)
}

func testGetTokens(t *testing.T) {
	var err error

	c := token.TokenGetCmd()
	c.Command.Flags().Set(authservice.ArgTokenId, *testToken.Id)

	err = c.Command.Execute()
	assert.NoError(t, err)

	c = token.TokenGetCmd()
	c.Command.Flags().Set(authservice.ArgToken, *tokenContent.Token)

	err = c.Command.Execute()
	assert.NoError(t, err)
}

func testParseToken(t *testing.T) {
	var err error

	c := token.TokenParseCmd()
	c.Command.Flags().Set(authservice.ArgToken, *tokenContent.Token)

	err = c.Command.Execute()
	assert.NoError(t, err)

	c.Command.Flags().Set(authservice.ArgPrivileges, "true")
	err = c.Command.Execute()
	assert.NoError(t, err)
}

func testDeleteTokens(t *testing.T) {
	var err error
	viper.Set(constants.ArgForce, true)

	c := token.TokenDeleteCmd()
	c.Command.Flags().Set(authservice.ArgTokenId, *testToken.Id)

	err = c.Command.Execute()
	assert.NoError(t, err)

	c = token.TokenDeleteCmd()
	c.Command.Flags().Set(authservice.ArgToken, *tokenContent.Token)

	err = c.Command.Execute()
	assert.NoError(t, err)
}
