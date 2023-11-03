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
