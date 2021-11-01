package pg

import (
	"errors"
	"testing"

	"github.com/ionos-cloud/ionosctl/internal/core"
	"github.com/stretchr/testify/assert"
)

func TestDBaaSPgsqlCmd(t *testing.T) {
	var err error
	core.RootCmdTest.AddCommand(DBaaSPgCmd())
	if ok := DBaaSPgCmd().IsAvailableCommand(); !ok {
		err = errors.New("non-available cmd")
	}
	assert.NoError(t, err)
}
