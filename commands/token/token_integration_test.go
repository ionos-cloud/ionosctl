//go:build integration
// +build integration

package token_test

import (
	"bytes"
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
	jwt                  sdkgoauth.Jwt
	tokenContent         string
)

func TestTokenCommands(t *testing.T) {
	err := setup()
	if err != nil {
		assert.FailNow(t, err.Error())
	}

	viper.Set(constants.ArgOutput, "text")

	t.Run("create token", testCreateToken)
	t.Run("list tokens", testListTokens)
	t.Run("get tokens by id and by token content", testGetTokens)
	t.Run("parse token content", testParseToken)
	t.Run("delete tokens by id and by token content", testDeleteTokens)
}

func setup() error {
	username := os.Getenv("IONOS_USERNAME")
	password := os.Getenv("IONOS_PASSWORD")

	if username == "" || password == "" {
		return fmt.Errorf("empty user/password")
	}

	var err error

	cl = client.NewClient(username, password, "", "")

	return err
}

func testCreateToken(t *testing.T) {
	var err error

	viper.Set(constants.ArgQuiet, true)

	tokFirstCreationTime = time.Now().In(time.UTC)
	c := token.TokenPostCmd()
	err = c.Command.Execute()
	assert.NoError(t, err)

	time.Sleep(5 * time.Second)

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
	if foundTokenViaSdk == nil {
		assert.FailNow(t, "created token could not be found")
	}

	testToken = *foundTokenViaSdk

	viper.Reset()
	viper.Set(constants.ArgOutput, "text")

	buff := bytes.NewBuffer([]byte{})
	c = token.TokenPostCmd()
	c.Command.SetOut(buff)
	err = c.Command.Execute()
	assert.NoError(t, err)

	tokenContent = buff.String()
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
	c.Command.Flags().Set(authservice.ArgToken, tokenContent)

	err = c.Command.Execute()
	assert.NoError(t, err)
}

func testParseToken(t *testing.T) {
	var err error

	c := token.TokenParseCmd()
	c.Command.Flags().Set(authservice.ArgToken, tokenContent)

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
	c.Command.Flags().Set(authservice.ArgToken, tokenContent)

	err = c.Command.Execute()
	assert.NoError(t, err)
}
