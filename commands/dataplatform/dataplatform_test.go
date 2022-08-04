package dataplatform

import (
	"errors"
	"testing"

	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/stretchr/testify/assert"
)

func TestDataPlatformServiceCmd(t *testing.T) {
	var err error
	core.RootCmdTest.AddCommand(DataPlatformServiceCmd())
	if ok := DataPlatformServiceCmd().IsAvailableCommand(); !ok {
		err = errors.New("non-available cmd")
	}
	assert.NoError(t, err)
}
