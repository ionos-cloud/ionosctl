package resources

import (
	"testing"

	"github.com/ionos-cloud/ionosctl/v6/pkg/config"

	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/stretchr/testify/assert"
)

func TestNewClientService(t *testing.T) {
	t.Run("no credentials", func(t *testing.T) {
		svc, err := config.NewTestClient("", "", "", "")
		assert.Nil(t, svc)
		assert.EqualError(t, err, "username, password or token incorrect")

		svc, err = config.NewTestClient("needspassword", "", "", "url")
		assert.Nil(t, svc)
		assert.EqualError(t, err, "username, password or token incorrect")
	})

	t.Run("success", func(t *testing.T) {
		svc, err := config.NewTestClient("", "", "token", "url")
		assert.NotNil(t, svc)
		assert.NoError(t, err)
		assert.Equal(t, "token", svc.PostgresClient.GetConfig().Token)

		svc, err = config.NewTestClient("user", "pass", "", "url")
		assert.NotNil(t, svc)
		assert.NoError(t, err)
		assert.Equal(t, "user", svc.PostgresClient.GetConfig().Username)
		assert.Equal(t, "pass", svc.PostgresClient.GetConfig().Password)
	})
}

func getTestClient(t *testing.T) *config.Client {
	svc, err := config.NewTestClient("user", "pass", "", constants.DefaultApiURL)
	assert.NotNil(t, svc)
	assert.NoError(t, err)
	assert.Equal(t, "user", svc.PostgresClient.GetConfig().Username)
	assert.Equal(t, "pass", svc.PostgresClient.GetConfig().Password)
	return svc
}
