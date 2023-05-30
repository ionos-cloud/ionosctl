package client_test

import (
	"context"
	"os"
	"testing"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
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
	if tk := os.Getenv("IONOS_TOKEN"); tk != "" {
		GoodToken = tk
	} else {
		pre(t)
	}

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
	t.Parallel()

	t.Run("Client Get works, user & pass env", func(t *testing.T) {
		cl, err := client.Get()
		assert.NoError(t, err)
		_, _, err = cl.CloudClient.DataCentersApi.DatacentersGet(context.Background()).Limit(1).Depth(0).Execute()
		assert.NoError(t, err)
	})

	t.Run("Client Get works, token", func(t *testing.T) {
		viper.Reset()
		os.Clearenv()

		viper.Set("IONOS_TOKEN", GoodToken)

		cl, err := client.Get()
		assert.NoError(t, err)
		_, _, err = cl.CloudClient.DataCentersApi.DatacentersGet(context.Background()).Limit(1).Depth(0).Execute()
		assert.NoError(t, err)
	})

	t.Run("Client Get fails 1", func(t *testing.T) {
		viper.Reset()
		os.Clearenv()

		viper.Set("IONOS_TOKEN", "foobar")

		cl, err := client.Get()
		assert.NoError(t, err)
		_, _, err = cl.CloudClient.DataCentersApi.DatacentersGet(context.Background()).Limit(1).Depth(0).Execute()
		assert.NoError(t, err)
	})

	t.Run("Client Get fails 2", func(t *testing.T) {
		viper.Reset()
		os.Clearenv()

		viper.Set("IONOS_USERNAME", "foo")
		viper.Set("IONOS_PASSWORD", GoodPassword)

		cl, err := client.Get()
		assert.NoError(t, err)
		_, _, err = cl.CloudClient.DataCentersApi.DatacentersGet(context.Background()).Limit(1).Depth(0).Execute()
		assert.NoError(t, err)
	})

	t.Run("Client Get fails, bad token priority", func(t *testing.T) {
		viper.Set("IONOS_TOKEN", "foobar") // Bad token still has priority
		viper.Set("IONOS_USERNAME", GoodUsername)
		viper.Set("IONOS_PASSWORD", GoodPassword)

		cl, err := client.Get()
		assert.NoError(t, err)
		_, _, err = cl.CloudClient.DataCentersApi.DatacentersGet(context.Background()).Limit(1).Depth(0).Execute()
		assert.NoError(t, err)
	})

}
