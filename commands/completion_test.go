package commands

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/stretchr/testify/assert"
)

func TestRunCompletionBash(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		cfg.Printer.SetStderr(w)
		err := RunCompletionBash(cfg)
		assert.NoError(t, err)
	})
}

func TestRunCompletionZsh(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		cfg.Printer.SetStderr(w)
		err := RunCompletionZsh(cfg)
		assert.NoError(t, err)
	})
}

func TestRunCompletionFish(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		cfg.Printer.SetStderr(w)
		err := RunCompletionFish(cfg)
		assert.NoError(t, err)
	})
}

func TestRunCompletionPowerShell(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		cfg.Printer.SetStderr(w)
		err := RunCompletionPowerShell(cfg)
		assert.NoError(t, err)
	})
}
