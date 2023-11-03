package postgres

import (
	"errors"
	"testing"

	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/stretchr/testify/assert"
)

func TestDBaaSPgsqlCmd(t *testing.T) {
	var err error
	core.RootCmdTest.AddCommand(DBaaSPostgresCmd())
	if ok := DBaaSPostgresCmd().IsAvailableCommand(); !ok {
		err = errors.New("non-available cmd")
	}
	assert.NoError(t, err)
}
