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

	assert.NotEmpty(t, GoodUsername)
	assert.NotEmpty(t, GoodPassword)

	tok, _, err := client.Must(func(err error) {
		t.Fatalf(err.Error())
	}).AuthClient.TokensApi.TokensGenerate(context.Background()).Execute()
	assert.NoError(t, err)

	GoodToken = *tok.Token
}

func TestClientPkg(t *testing.T) {
	assert.NotEmpty(t, os.Getenv("IONOS_USERNAME"))
	assert.NotEmpty(t, os.Getenv("IONOS_PASSWORD"))

	pre(t)

	assert.NotEmpty(t, GoodUsername)
	assert.NotEmpty(t, GoodPassword)
	assert.NotEmpty(t, GoodToken)

	viper.Reset()
	os.Clearenv()

	testNewClient(t)
	testGetClient(t)
	testTestCreds(t)
}

func testTestCreds(t *testing.T) {
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

func testNewClient(t *testing.T) {
	t.Run("should return a valid client when token is provided", func(t *testing.T) {
		cl := client.NewClient("", "", "some-token", "")

		assert.NotNil(t, cl)
		assert.NoError(t, cl.TestCreds())
	})

	t.Run("should return a valid client when username and password are provided", func(t *testing.T) {
		cl := client.NewClient("username", "password", "", "")

		assert.NotNil(t, cl)
		assert.NoError(t, cl.TestCreds())

	})
}

func testGetClient(t *testing.T) {

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
