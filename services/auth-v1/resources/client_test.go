package resources

import (
	"testing"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"

	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
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

// TODO: Find a way to replace all of these duped funcs in the entirety of the codebase via a regex with NewClient
func getTestClient(t *testing.T) *client.Client {
	svc := client.NewClient("user", "pass", "", constants.DefaultApiURL)
	assert.NotNil(t, svc)
	assert.Equal(t, "user", svc.AuthClient.GetConfig().Username)
	assert.Equal(t, "pass", svc.AuthClient.GetConfig().Password)
	return svc
}
