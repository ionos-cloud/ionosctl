package commands

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/ionos-cloud/ionosctl/pkg/builder"
	"github.com/stretchr/testify/assert"
)

func TestRunVersion(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		err := RunVersion(cfg)
		assert.NoError(t, err)
	})
}
