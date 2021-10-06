package cloudapi_dbaas_pgsql

import (
	"errors"
	"testing"

	"github.com/ionos-cloud/ionosctl/internal/core"
	"github.com/stretchr/testify/assert"
)

func TestDBaaSPgsqlCmd(t *testing.T) {
	var err error
	core.RootCmdTest.AddCommand(DBaaSPgsqlCmd())
	if ok := DBaaSPgsqlCmd().IsAvailableCommand(); !ok {
		err = errors.New("non-available cmd")
	}
	assert.NoError(t, err)
}
