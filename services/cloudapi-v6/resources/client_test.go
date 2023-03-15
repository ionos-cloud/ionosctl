package resources

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"testing"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"

	"github.com/stretchr/testify/assert"
)

func TestNewClientService(t *testing.T) {
	t.Run("no credentials", func(t *testing.T) {
		svc, err := client.NewTestClient("needspassword", "", "", "url")
		assert.Nil(t, svc)
		assert.EqualError(t, err, "username, password or token incorrect")
	})

	t.Run("success", func(t *testing.T) {
		svc, err := client.NewTestClient("", "", "token", "url")
		assert.NotNil(t, svc)
		assert.NoError(t, err)
		assert.Equal(t, "token", svc.CloudClient.GetConfig().Token)

		svc, err = client.NewTestClient("user", "pass", "", "url")
		assert.NotNil(t, svc)
		assert.NoError(t, err)
		assert.Equal(t, "user", svc.CloudClient.GetConfig().Username)
		assert.Equal(t, "pass", svc.CloudClient.GetConfig().Password)

		svc, err = client.NewTestClient("user", "pass", "", "")
		assert.NotNil(t, svc)
		assert.NoError(t, err)
		assert.Equal(t, "user", svc.CloudClient.GetConfig().Username)
		assert.Equal(t, "pass", svc.CloudClient.GetConfig().Password)
	})
}
