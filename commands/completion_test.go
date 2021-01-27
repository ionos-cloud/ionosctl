package commands

import (
	"bufio"
	"bytes"
	"github.com/ionos-cloud/ionosctl/pkg/builder"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRunCompletionBash(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		cfg.Printer.Stdout = w
		err := RunCompletionBash(cfg)
		assert.NoError(t, err)
	})
}

func TestRunCompletionZsh(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		cfg.Printer.Stdout = w
		err := RunCompletionZsh(cfg)
		assert.NoError(t, err)
	})
}
