//go:build integration
// +build integration

package authv1_test

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	authv1 "github.com/ionos-cloud/ionosctl/v6/commands/auth-v1"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	authservices "github.com/ionos-cloud/ionosctl/v6/services/auth-v1"
	sdkgoauth "github.com/ionos-cloud/sdk-go-auth"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	cl                   *client.Client
	tokFirstCreationTime time.Time
	token                sdkgoauth.Token
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

	return err
}

func testCreateToken(t *testing.T) {
	var err error

	c := authv1.TokenPostCmd()
	tokFirstCreationTime = time.Now().In(time.UTC)
	err = c.Command.Execute()
	assert.NoError(t, err)

	tokens, _, err := client.Must().AuthClient.TokensApi.TokensGet(context.Background()).Execute()
	assert.NoError(t, err)

	allTokens, ok := tokens.GetTokensOk()
	assert.NotEmpty(t, ok)
	assert.NotEmpty(t, *allTokens)

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
			token = tok
		}
	}

	tokenContent, _, err = cl.AuthClient.TokensApi.TokensGenerate(context.Background()).Execute()
	assert.NoError(t, err)
	assert.NotEmpty(t, *tokenContent.Token)
}

func testListTokens(t *testing.T) {
	var err error

	c := authv1.TokenListCmd()
	err = c.Command.Execute()
	assert.NoError(t, err)
}

func testGetTokens(t *testing.T) {
	var err error

	c := authv1.TokenGetCmd()
	c.Command.Flags().Set(authservices.ArgTokenId, *token.Id)

	err = c.Command.Execute()
	assert.NoError(t, err)

	c = authv1.TokenGetCmd()
	c.Command.Flags().Set(authservices.ArgToken, *tokenContent.Token)

	err = c.Command.Execute()
	assert.NoError(t, err)
}

func testParseToken(t *testing.T) {
	var err error

	c := authv1.TokenParseCmd()
	c.Command.Flags().Set(authservices.ArgToken, *tokenContent.Token)

	err = c.Command.Execute()
	assert.NoError(t, err)

	c.Command.Flags().Set(authservices.ArgPrivileges, "true")
	err = c.Command.Execute()
	assert.NoError(t, err)
}

func testDeleteTokens(t *testing.T) {
	var err error
	viper.Set(constants.ArgForce, true)

	c := authv1.TokenDeleteCmd()
	c.Command.Flags().Set(authservices.ArgTokenId, *token.Id)

	err = c.Command.Execute()
	assert.NoError(t, err)

	c = authv1.TokenDeleteCmd()
	c.Command.Flags().Set(authservices.ArgToken, *tokenContent.Token)

	err = c.Command.Execute()
	assert.NoError(t, err)
}
