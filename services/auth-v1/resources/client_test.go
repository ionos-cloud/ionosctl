package resources

import (
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"testing"

	"github.com/ionos-cloud/ionosctl/pkg/constants"
	"github.com/stretchr/testify/assert"
)

func TestNewClientService(t *testing.T) {
	t.Run("no credentials", func(t *testing.T) {
		svc, err := NewClientService("needspassword", "", "", "url")
		assert.Nil(t, svc)
		assert.EqualError(t, err, "username, password or token incorrect")
	})

	t.Run("success", func(t *testing.T) {
		svc, err := NewClientService("", "", "token", "url")
		assert.NotNil(t, svc)
		assert.NoError(t, err)
		assert.Equal(t, "token", svc.GetConfig().Token)

		svc, err = NewClientService("user", "pass", "", "url")
		assert.NotNil(t, svc)
		assert.NoError(t, err)
		assert.Equal(t, "user", svc.GetConfig().Username)
		assert.Equal(t, "pass", svc.GetConfig().Password)

		svc, err = NewClientService("user", "pass", "", "")
		assert.NotNil(t, svc)
		assert.NoError(t, err)
		assert.Equal(t, "user", svc.GetConfig().Username)
		assert.Equal(t, "pass", svc.GetConfig().Password)
	})
}

func getTestClient(t *testing.T) *config.Client {
	svc, err := config.NewTestClient("user", "pass", "", constants.DefaultApiURL)
	assert.NotNil(t, svc)
	assert.NoError(t, err)
	assert.Equal(t, "user", svc.AuthClient.GetConfig().Username)
	assert.Equal(t, "pass", svc.AuthClient.GetConfig().Password)
	return svc
}
