package dbaas

import (
	"errors"
	"testing"

	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/stretchr/testify/assert"
)

func TestDataBaseServiceCmd(t *testing.T) {
	var err error
	core.RootCmdTest.AddCommand(DataBaseServiceCmd())
	if ok := DataBaseServiceCmd().IsAvailableCommand(); !ok {
		err = errors.New("non-available cmd")
	}
	assert.NoError(t, err)
}

func TestV2AliasesHidden(t *testing.T) {
	cmd := DataBaseServiceCmd().Command
	for _, tc := range []struct{ visible, hidden string }{
		{"postgres", "postgres-v2"},
		{"mariadb", "mariadb-v2"},
	} {
		v, _, err := cmd.Find([]string{tc.visible})
		assert.NoError(t, err)
		assert.Equal(t, tc.visible, v.Name())
		assert.False(t, v.Hidden)

		h, _, err := cmd.Find([]string{tc.hidden, "cluster", "list"})
		assert.NoError(t, err)
		assert.Equal(t, "list", h.Name())
		assert.True(t, h.Parent().Parent().Hidden)
	}
}
