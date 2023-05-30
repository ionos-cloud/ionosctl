package client_test

import (
	"context"
	"os"
	"testing"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	GoodUsername = ""
	GoodPassword = ""
	GoodToken    = ""
)

func pre(t *testing.T) {
	GoodUsername = os.Getenv("IONOS_USERNAME")
	GoodPassword = os.Getenv("IONOS_PASSWORD")

	tok, _, err := client.Must().AuthClient.TokensApi.TokensGenerate(context.Background()).Execute()
	assert.NoError(t, err)

	GoodToken = *tok.Token
}

func TestClientPkg(t *testing.T) {
	pre(t)

	viper.Reset()
	os.Clearenv()

	TestNewClient(t)
	TestGet(t)
	TestTestCreds(t)
}

func TestTestCreds(t *testing.T) {
	t.Parallel()

	t.Run("empty creds", func(t *testing.T) {
		err := client.TestCreds("", "", "")
		assert.ErrorContains(t, err, "empty")
	})

	t.Run("good user & pass", func(t *testing.T) {
		err := client.TestCreds(GoodUsername, GoodPassword, "")
		assert.NoError(t, err)
	})

	t.Run("good token 1", func(t *testing.T) {
		err := client.TestCreds("foo", "bar", GoodToken)
		assert.NoError(t, err)
	})

	t.Run("good token 2", func(t *testing.T) {
		err := client.TestCreds("", "", GoodToken)
		assert.NoError(t, err)
	})

	t.Run("good user & pass & token", func(t *testing.T) {
		err := client.TestCreds(GoodUsername, GoodPassword, GoodToken)
		assert.NoError(t, err)
	})

	t.Run("bad creds 1", func(t *testing.T) {
		err := client.TestCreds("foo", "bar", "tok")
		assert.ErrorContains(t, err, "credentials test failed")
	})

	t.Run("bad creds 2", func(t *testing.T) {
		err := client.TestCreds("foo", "bar", "")
		assert.ErrorContains(t, err, "credentials test failed")
	})

	t.Run("bad creds 3", func(t *testing.T) {
		err := client.TestCreds("", "", "tok")
		assert.ErrorContains(t, err, "credentials test failed")
	})

	t.Run("bad creds 4", func(t *testing.T) {
		err := client.TestCreds("foo", "", "")
		assert.ErrorContains(t, err, "empty")
	})
}

func TestNewClient(t *testing.T) {
	t.Run("should return an error when both token and username/password are empty", func(t *testing.T) {
		cl, err := client.NewClient("", "", "", "")

		assert.Error(t, err)
		assert.ErrorContains(t, err, "empty")
		assert.Nil(t, cl)
	})

	t.Run("should return a valid client when token is provided", func(t *testing.T) {
		cl, err := client.NewClient("", "", "some-token", "")

		assert.NoError(t, err)
		assert.NotNil(t, cl)
	})

	t.Run("should return a valid client when username and password are provided", func(t *testing.T) {
		cl, err := client.NewClient("username", "password", "", "")

		assert.NoError(t, err)
		assert.NotNil(t, cl)
	})
}

func TestGet(t *testing.T) {
	t.Run("Client Get works fine", func(t *testing.T) {
		cl, err := client.Get()
		assert.NoError(t, err)
		_, _, err = cl.CloudClient.DataCentersApi.DatacentersGet(context.Background()).Limit(1).Depth(0).Execute()
		assert.NoError(t, err)
	})

	t.Run("should return a client when token is provided", func(t *testing.T) {
		viper.Set(constants.ArgToken, "some-token")

		cl, err := client.Get()

		assert.NoError(t, err)
		assert.NotNil(t, cl)

		viper.Reset()
	})

}
