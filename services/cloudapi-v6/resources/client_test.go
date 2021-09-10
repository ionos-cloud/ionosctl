package resources

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClientService(t *testing.T) {
	t.Run("no credentials", func(t *testing.T) {
		svc, err := NewClientService("", "", "", "")
		assert.Nil(t, svc)
		assert.EqualError(t, err, "host-url incorrect")

		svc, err = NewClientService("needspassword", "", "", "url")
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
	})
}
